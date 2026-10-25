import assert from "node:assert/strict"
import test from "node:test"
import path from "node:path"
import { mkdtemp, mkdir, writeFile, rm, symlink, realpath } from "node:fs/promises"
import { tmpdir } from "node:os"
import orchestration from "../plugins/orchestration.ts"
import { dispatch as dispatchAdapter } from "../tools/orchestration.ts"

const { State, Runtime, wrapper, operationFiles, envelope, canonical, roleContract } = orchestration.testing
const handoff = "<Γ>\nSIGMA/1: PIDGEN/text\nRole: researcher\nTarget: inspect\nAcceptance: evidence\nScope: a.txt\nKnown facts: N_A\nSIGMA/1\n</Γ>"
const reserve = (s, parent, role, id, writes = [], scope = []) => s.reserve(parent, role, scope, writes, id, `msg_${id}`, 0, 900000)
const args = (role = "researcher", writes = [], difficulty = "easy") => ({ role, ...(role === "worker" ? { difficulty } : {}), handoff, scope_files: ["a.txt", "b.txt"], write_files: writes })

test("plugin and discovery adapter both expose difficulty and one-shot escalation fields", async (t) => {
  const f = await fixture(t)
  const pluginArgs = f.runtime.hooks().tool.orchestration_dispatch.args
  for (const args of [pluginArgs, dispatchAdapter.args]) {
    assert.ok(Object.hasOwn(args, "difficulty"))
    assert.ok(Object.hasOwn(args, "escalation_of"))
  }
  // Legacy discovery adapter is unchanged; Attacker must use the plugin tool.
  for (const key of Object.keys(dispatchAdapter.args)) assert.ok(Object.hasOwn(pluginArgs, key))
  for (const key of ["review_depth", "review_target", "review_revision", "review_lens", "review_approved"]) assert.ok(Object.hasOwn(pluginArgs, key))
})
const rootContext = { sessionID: "root", agent: "orchestrator", messageID: "root-message" }
const tick = () => new Promise((resolve) => setTimeout(resolve, 1))
async function settle(runtime) { for (let i = 0; i < 100 && runtime.pending.size; i++) await tick(); assert.equal(runtime.pending.size, 0) }

async function fixture(t, overrides = {}) {
  const root = await mkdtemp(path.join(process.env.TMPDIR ?? tmpdir(), "opencode-rts-test-"))
  await writeFile(path.join(root, "a.txt"), "a"); await writeFile(path.join(root, "b.txt"), "b")
  let n = 0
  const sessions = new Map([["root", { id: "root" }]])
  const messages = new Map([["root", [{ info: { id: "root-message", sessionID: "root", role: "assistant", providerID: "openai", modelID: "gpt-6-luna" }, parts: [] }]]])
  const calls = [], aborts = []
  const input = { worktree: root, directory: root, serverUrl: new URL("http://localhost:1234"), client: {
     app: { agents: async () => ({ data: [
       ...["researcher", "proposer", "worker", "verifier", "auditor", "attacker", "blackhat"].map((name) => ({ name, mode: "subagent" })),
       { name: "attacker-low", mode: "subagent", model: "openai/gpt-6-luna", reasoningEffort: "low" },
       { name: "attacker-medium", mode: "subagent", model: "openai/gpt-6.1-sol", reasoningEffort: "high" },
       { name: "attacker-hard", mode: "subagent", model: "openai/gpt-6-astra" },
       { name: "worker-easy", mode: "subagent", model: { providerID: "openai", modelID: "gpt-6-luna" } },
       { name: "worker-medium", mode: "subagent", model: { providerID: "openai", modelID: "gpt-6-luna" } },
       { name: "worker-hard", mode: "subagent", model: { providerID: "openai", modelID: "gpt-6.1-sol" } },
       { name: "worker-astra", mode: "subagent", hidden: true, model: { providerID: "openai", modelID: "gpt-6-astra" } },
     ] }) },
    session: {
      get: async (o) => ({ data: sessions.get(o.path.id) }),
      create: async (o) => { const s = { id: `ses_${++n}`, parentID: o.body.parentID }; sessions.set(s.id, s); calls.push(["create", o]); return { data: s } },
      promptAsync: async (o) => { calls.push(["promptAsync", o]); return { data: undefined, response: { ok: true } } },
      prompt: async () => { throw new Error("NEVER prompt parent") },
      messages: async (o) => ({ data: messages.get(o.path.id) ?? [] }),
      abort: async (o) => { aborts.push(o.path.id); return { data: true } },
      status: async () => ({ data: Object.fromEntries([...sessions.keys()].map((id) => [id, { type: "idle" }])) }),
      ...overrides,
    },
  } }
  let checks = 0
  const runtime = new Runtime(input, root, async () => { checks++ }, 60000)
  runtime.rootAgents.set("root", "orchestrator")
  t.after(async () => { await runtime.dispose(); await rm(root, { recursive: true, force: true }) })
  return { root, runtime, input, sessions, messages, calls, aborts, checks: () => checks }
}

function final(j, extra = {}) { return { id: `msg_final_${j.id}`, sessionID: j.session, parentID: j.promptID, role: "assistant", providerID: j.requestedModel?.providerID ?? "openai", modelID: j.requestedModel?.modelID ?? "gpt-6-luna", time: { created: 1, completed: 2 }, finish: "stop", ...extra } }

test("attacker and blackhat dispatch as read-only leaves using the retained shared contract", async (t) => {
  const f = await fixture(t)
  for (const role of ["attacker", "blackhat"]) {
    assert.equal(roleContract(role), "blackhat.md")
    const reviewArgs = { ...args(role), review_target: "proposal", review_revision: "rev-1", review_lens: "trust-boundary", handoff: handoff.replace("Role: researcher", `Role: ${role}`) }
    assert.throws(() => f.runtime.dispatch({ ...reviewArgs, write_files: ["a.txt"] }, rootContext), /only worker may own writes/)
    assert.throws(() => f.runtime.dispatch({ ...reviewArgs, difficulty: "hard" }, rootContext), /only for CODE worker/)
    assert.throws(() => f.runtime.dispatch({ ...reviewArgs, escalation_of: "source" }, rootContext), /review_depth=hard/)
    const id = f.runtime.dispatch(reviewArgs, rootContext).task_id
    await settle(f.runtime)
    const job = f.runtime.state.jobs.get(id)
    assert.equal(job.role, "attacker")
    assert.equal(job.profile, "attacker-low")
    assert.deepEqual(job.requestedModel, { providerID: "openai", modelID: "gpt-6-luna" })
    assert.deepEqual(job.writes, [])
    assert.equal(f.calls.find(([name, o]) => name === "promptAsync" && o.path.id === job.session)[1].body.agent, "attacker-low")
    const before = (tool, args) => f.runtime.before({ tool, sessionID: job.session }, args)
    await before("read", { filePath: "a.txt" })
    await assert.rejects(() => before("read", { filePath: "outside.txt" }))
    for (const [tool, input] of [
      ["edit", { filePath: "a.txt", oldString: "a", newString: "x" }],
      ["write", { filePath: "a.txt", content: "x" }],
      ["write_files", { files: [{ file_path: "a.txt", content: "x" }] }],
      ["apply_patch", { patchText: "*** Begin Patch\n*** Update File: a.txt\n@@\n-a\n+x\n*** End Patch" }],
      ["bash", { command: "anything" }], ["webfetch", { url: "https://example.com" }],
      ["task", {}], ["orchestration_dispatch", {}],
    ]) await assert.rejects(() => before(tool, input))
    assert.throws(() => f.runtime.dispatch(args("verifier"), { sessionID: job.session, agent: role, messageID: "review-message" }), /eligible/)
  }
  assert.equal(roleContract("verifier"), "verifier.md")
})

test("exact wrapper and explicit HTTP envelope validation fail closed", () => {
  wrapper(handoff)
  for (const bad of [handoff + "\n", handoff.replace("Known facts:", "Facts:"), handoff.replace("\nSIGMA/1\n</Γ>", "\n</Γ>"), handoff.replace("Scope: a.txt", "Scope: a.txt\nScope: b.txt")]) assert.throws(() => wrapper(bad))
  for (const r of [undefined, { error: { code: 400 } }, { response: { ok: false } }, {}]) assert.throws(() => envelope(r, "op"))
  assert.equal(envelope({ response: { ok: true } }, "promptAsync", false), undefined)
})

test("Runtime-wide capacity allows 20 descendants, rejects 21, and shares budget across roots/nesting", () => {
  const s = new State()
  for (let i = 0; i < 20; i++) reserve(s, "root", "researcher", `c${i}`)
  assert.equal(s.jobs.size, 20)
  assert.throws(() => reserve(s, "second-root", "researcher", "overflow"), /capacity exhausted/)
  // A separate Runtime owns an independent State and budget.
  const independent = new State()
  reserve(independent, "root", "researcher", "independent")

  const nested = new State()
  const first = reserve(nested, "root-a", "researcher", "root-child", [], ["a"])
  first.session = "child-a"; nested.sessions.set(first.session, first)
  const grand = reserve(nested, "child-a", "researcher", "grandchild", [], ["a"])
  grand.session = "grandchild-session"; nested.sessions.set(grand.session, grand)
  for (let i = 2; i < 20; i++) reserve(nested, "root-b", "researcher", `other-root-${i}`)
  assert.equal(nested.jobs.size, 20)
  assert.throws(() => reserve(nested, "root-b", "researcher", "twenty-first"), /capacity exhausted/)
  assert.equal(grand.depth, 2)
  assert.throws(() => reserve(nested, "grandchild-session", "researcher", "depth3"), /eligible/)
  assert.throws(() => reserve(nested, "child-a", "worker", "write-grand", ["a"], ["a"]), /read-only/)
  assert.throws(() => reserve(nested, "child-a", "researcher", "escape", [], ["b"]), /inherited/)

  const direct = new State()
  for (let i = 0; i < 20; i++) reserve(direct, "root", "researcher", `direct-${i}`)
  assert.throws(() => reserve(direct, "root", "researcher", "direct-21"), /capacity exhausted/)
  const own = new State(); reserve(own, "root", "worker", "w1", ["a"], ["a"])
  assert.throws(() => reserve(own, "root", "worker", "w2", ["a"], ["a"]), /overlapping/)
  assert.throws(() => reserve(own, "root", "auditor", "w3", ["b"], ["b"]), /only worker/)
  const w = own.jobs.get("w1"); w.session = "w"; own.sessions.set("w", w)
  assert.throws(() => reserve(own, "w", "researcher", "leaf"), /eligible/)
})

test("separate Runtime instances retain independent State capacity", async (t) => {
  const a = await fixture(t), b = await fixture(t)
  assert.notEqual(a.runtime.state, b.runtime.state)
  for (let i = 0; i < 20; i++) reserve(a.runtime.state, "root", "researcher", `a${i}`)
  reserve(b.runtime.state, "root", "researcher", "b0")
  assert.equal(a.runtime.state.jobs.size, 20)
  assert.equal(b.runtime.state.jobs.size, 1)
})

test("reservations retain capacity while blocked and release it only after confirmed stop", () => {
  const state = new State()
  const first = reserve(state, "root-a", "researcher", "reserved")
  for (let i = 1; i < 20; i++) reserve(state, "root-b", "researcher", `fill-${i}`)
  assert.equal(first.pendingLaunch, true)
  assert.throws(() => reserve(state, "root-b", "researcher", "overflow"), /capacity exhausted/)

  first.state = "blocked"
  assert.equal(first.stopped, false)
  assert.throws(() => reserve(state, "root-b", "researcher", "still-full"), /capacity exhausted/)

  state.stop(first, "completed")
  reserve(state, "root-b", "researcher", "slot-reused")
  assert.equal([...state.jobs.values()].filter((job) => !job.stopped).length, 20)
})

test("canonical Runtime registry shares state only for the same server and worktree", async (t) => {
  const a = await fixture(t), b = await fixture(t)
  const owned = new Set()
  t.after(async () => { await Promise.all([...owned].map((runtime) => runtime.dispose())) })
  const openRuntime = async (f, serverUrl = f.input.serverUrl) => {
    const input = { ...f.input, serverUrl }
    await orchestration(input)
    const resolved = await realpath(f.root)
    const canonicalRoot = process.platform === "win32" ? resolved.toLowerCase() : resolved
    const identity = `${serverUrl.toString()}|${canonicalRoot}`
    const promise = orchestration.testing.registry.get(identity)
    assert.ok(promise)
    const runtime = await promise
    owned.add(runtime)
    return runtime
  }

  const first = await openRuntime(a)
  assert.equal(await openRuntime(a), first)
  const otherServer = await openRuntime(a, new URL("http://localhost:5678"))
  const otherWorktree = await openRuntime(b)
  assert.notEqual(otherServer, first)
  assert.notEqual(otherWorktree, first)
  assert.notEqual(otherWorktree, otherServer)
})

test("dispatch ID is immediate; real SDK parentID links grandchildren; preflight before each submission", async (t) => {
  let release
  const f = await fixture(t, { promptAsync: async (o) => { f.calls.push(["promptAsync", o]); await new Promise((r) => { release = r }); return { response: { ok: true } } } })
  const response = f.runtime.dispatch(args(), rootContext)
  assert.equal(typeof response.task_id, "string"); assert.equal(response.state, "reserved")
  assert.equal(f.calls.length, 0)
  while (!release) await tick()
  const child = f.runtime.state.jobs.get(response.task_id)
  assert.equal(child.parent, "root"); assert.equal(f.sessions.get(child.session).parentID, "root")
  release(); await settle(f.runtime)
  f.input.client.session.promptAsync = async (o) => { f.calls.push(["promptAsync", o]); return { response: { ok: true } } }
   const childMessageID = "researcher-parent-message"
   f.messages.set(child.session, [{ info: { id: childMessageID, sessionID: child.session, role: "assistant", providerID: "openai", modelID: "gpt-6-luna" }, parts: [] }])
   const grandID = f.runtime.dispatch(args("verifier"), { sessionID: child.session, agent: "researcher", messageID: childMessageID }).task_id
  await settle(f.runtime)
  const grand = f.runtime.state.jobs.get(grandID)
  assert.equal(grand.depth, 2); assert.equal(f.sessions.get(grand.session).parentID, child.session)
  assert.equal(f.calls.find(([name, o]) => name === "promptAsync" && o.path.id === grand.session)[1].body.messageID, grand.promptID)
  assert.ok(f.checks() >= 4)
  assert.ok(f.calls.filter(([name]) => name === "promptAsync").every(([, o]) => o.path.id !== "root" && !Object.hasOwn(o.body, "noReply")))
   assert.throws(() => f.runtime.dispatch(args(), { sessionID: grand.session, agent: "verifier", messageID: "grand-message" }), /eligible/)
})

test("idle/intermediate/unrelated messages are not success; immediate/reordered events correlate final evidence and dedup queue", async (t) => {
  const f = await fixture(t)
  const task = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id
  await settle(f.runtime); const j = f.runtime.state.jobs.get(task)
  await f.runtime.event({ type: "session.idle", properties: { sessionID: j.session } }); await settle(f.runtime)
  assert.equal(j.stopped, false)
  for (const info of [final(j, { finish: "tool-calls" }), final(j, { parentID: "other" }), final(j, { time: { created: 1 } })]) await f.runtime.event({ type: "message.updated", properties: { info } })
  assert.equal(j.stopped, false); assert.equal(f.runtime.state.leases.size, 1)
  await f.runtime.event({ type: "permission.updated", properties: { sessionID: j.session } }); assert.equal(j.state, "waiting_permission")
  await f.runtime.event({ type: "permission.replied", properties: { sessionID: j.session } }); assert.equal(j.state, "running")
  f.messages.set(j.session, [{ info: final(j), parts: [{ type: "text", text: "verified evidence" }] }])
  await f.runtime.event({ type: "message.updated", properties: { info: final(j) } }); await settle(f.runtime)
  await f.runtime.event({ type: "message.updated", properties: { info: final(j) } })
  await f.runtime.event({ type: "permission.updated", properties: { sessionID: j.session } })
  assert.equal(j.state, "completed"); assert.equal(j.result, "verified evidence"); assert.equal(f.runtime.state.leases.size, 0)
  assert.equal(f.runtime.state.queues.get("root").size, 1)
  const hooks = f.runtime.hooks(), out = { system: [] }
  await hooks["experimental.chat.system.transform"]({ sessionID: "root" }, out)
  assert.equal(out.system.length, 1); assert.equal(f.runtime.state.queues.get("root").size, 0)
  await hooks["experimental.chat.system.transform"]({ sessionID: "root" }, out); assert.equal(out.system.length, 1)
})

test("event during submission is never overwritten by acceptance or transport ambiguity", async (t) => {
  const f = await fixture(t)
  f.input.client.session.promptAsync = async (o) => {
    const j = f.runtime.state.sessions.get(o.path.id)
    await f.runtime.event({ type: "message.updated", properties: { info: final(j) } })
    throw new Error("connection lost after acceptance")
  }
  const id = f.runtime.dispatch(args(), rootContext).task_id
  await settle(f.runtime)
  assert.equal(f.runtime.state.jobs.get(id).state, "completed")
})

test("polling gets final errors and transport ambiguity never re-prompts; leases retained", async (t) => {
  const f = await fixture(t, { promptAsync: async (o) => { f.calls.push(["promptAsync", o]); throw new Error("ambiguous transport") } })
  const id = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id
  await settle(f.runtime); const j = f.runtime.state.jobs.get(id)
  assert.equal(j.state, "blocked"); assert.equal(f.runtime.state.leases.size, 1)
  await f.runtime.reconcile(); await f.runtime.reconcile()
  assert.equal(f.calls.filter(([name]) => name === "promptAsync").length, 1)
  f.messages.set(j.session, [{ info: final(j, { finish: undefined, error: { name: "UnknownError" } }), parts: [{ type: "text", text: "error evidence" }] }])
  await f.runtime.reconcile(); assert.equal(j.state, "failed"); assert.equal(j.result, "error evidence")
  assert.equal(f.runtime.state.leases.size, 0)
})

test("scope denial covers cross-worker edits, root writes, every patch path, malformed/opaque tools", async (t) => {
  const f = await fixture(t)
  const id = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id
  await settle(f.runtime); const j = f.runtime.state.jobs.get(id)
  const before = (name, a, sessionID = j.session) => f.runtime.before({ tool: name, sessionID }, a)
  await before("edit", { filePath: "a.txt", oldString: "a", newString: "z" })
  if (process.platform === "win32") await before("edit", { filePath: "A.txt", oldString: "a", newString: "z" })
  else await assert.rejects(() => before("edit", { filePath: "A.txt", oldString: "a", newString: "z" }), /scope/)
  await before("read", { filePath: "b.txt" })
  for (const [name, a] of [["write", { filePath: "b.txt", content: "no" }], ["edit", { filePath: "a.txt", oldString: "a", newString: "z", otherFiles: ["b.txt"] }], ["bash", { command: "test" }], ["task", {}], ["agy_worker", {}], ["webfetch", { url: "https://example.com" }], ["glob", {}], ["delete", { filePath: "a.txt" }]]) await assert.rejects(() => before(name, a))
  const multi = "*** Begin Patch\n*** Update File: a.txt\n@@\n-a\n+z\n*** Delete File: b.txt\n*** End Patch"
  await assert.rejects(() => before("apply_patch", { patchText: multi }), /scope/)
  await assert.rejects(() => before("apply_patch", { patchText: "*** Begin Patch\n*** Update File: a.txt\n*** Move to: b.txt\n@@\n-a\n+z\n*** End Patch" }), /scope/)
  await before("apply_patch", { patchText: "*** Begin Patch\n*** Delete File: a.txt\n*** End Patch" })
  await assert.rejects(() => before("write", { filePath: "a.txt", content: "root overwrite" }, "root"), /leased/)
  f.sessions.set("unknown", { id: "unknown", parentID: "root" })
  await assert.rejects(() => before("write", { filePath: "a.txt", content: "x" }, "unknown"), /unknown worker/)
  f.sessions.set("unknown-root-worker", { id: "unknown-root-worker" })
  await assert.rejects(() => before("write", { filePath: "b.txt", content: "x" }, "unknown-root-worker"), /identity/)
  assert.throws(() => f.runtime.dispatch(args("worker", [process.platform === "win32" ? "A.txt" : "a.txt"]), rootContext), /overlapping/)
})

test("junction aliases and junction changes are resolved on every use", async (t) => {
  const f = await fixture(t)
  const one = path.join(f.root, "one"), two = path.join(f.root, "two"), alias = path.join(f.root, "alias")
  await mkdir(one); await mkdir(two); await writeFile(path.join(one, "file.txt"), "one"); await writeFile(path.join(two, "file.txt"), "two")
  await symlink(one, alias, "junction")
  assert.equal(await canonical(f.root, "alias/file.txt"), await canonical(f.root, "one/file.txt"))
  const spec = { role: "worker", difficulty: "easy", handoff, scope_files: ["alias/file.txt"], write_files: ["alias/file.txt"] }
  const id = f.runtime.dispatch(spec, rootContext).task_id; await settle(f.runtime)
  const j = f.runtime.state.jobs.get(id)
  await f.runtime.before({ sessionID: j.session, tool: "write" }, { filePath: "one/file.txt", content: "ok" })
  await rm(alias); await symlink(two, alias, "junction")
  await assert.rejects(() => f.runtime.before({ sessionID: j.session, tool: "write" }, { filePath: "alias/file.txt", content: "no" }), /scope/)
  await assert.rejects(() => canonical(f.root, "../escape.txt"), /escapes/)
  await assert.rejects(() => canonical(f.root, "alias"), /not a file/)
  await assert.rejects(() => canonical(f.root, path.join(f.root, "a.txt:stream")), /invalid file path/)
})

test("multi-file/move parser rejects malformed/unknown operations", () => {
  assert.deepEqual(operationFiles("apply_patch", { patchText: "*** Begin Patch\n*** Update File: a\n*** Move to: b\n@@\n-x\n+y\n*** Add File: c\n+new\n*** Delete File: d\n*** End Patch" }).files, ["a", "b", "c", "d"])
  for (const patchText of ["bad", "*** Begin Patch\n*** Rename File: a\n*** End Patch", "*** Begin Patch\n*** Delete File: a\n+unexpected\n*** End Patch", "*** Begin Patch\n*** Add File: a\nnot-prefixed\n*** End Patch"]) assert.throws(() => operationFiles("apply_patch", { patchText }))
})

test("cancel subtree aborts deepest first and abort failure retains worker leases", async (t) => {
  const f = await fixture(t)
  const cID = f.runtime.dispatch(args(), rootContext).task_id; await settle(f.runtime)
  const c = f.runtime.state.jobs.get(cID)
   const childMessageID = "cancel-child-message"
   f.messages.set(c.session, [{ info: { id: childMessageID, sessionID: c.session, role: "assistant", providerID: "openai", modelID: "gpt-6-luna" }, parts: [] }])
   const gID = f.runtime.dispatch(args("verifier"), { sessionID: c.session, agent: "researcher", messageID: childMessageID }).task_id; await settle(f.runtime)
  const g = f.runtime.state.jobs.get(gID)
  await f.runtime.cancel(cID, "root"); assert.deepEqual(f.aborts, [g.session, c.session]); assert.equal(g.state, "cancelled"); assert.equal(c.state, "cancelled")
  const wID = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id; await settle(f.runtime)
  const w = f.runtime.state.jobs.get(wID)
  f.input.client.session.abort = async () => ({ error: { message: "failed abort" } })
  await f.runtime.cancel(wID, "root"); assert.equal(w.state, "blocked"); assert.equal(w.stopped, false); assert.equal(f.runtime.state.leases.size, 1)
  f.input.client.session.abort = async () => ({ data: true })
  f.input.client.session.status = async () => ({ data: { [w.session]: { type: "busy" } } })
  await f.runtime.cancel(wID, "root"); assert.equal(w.stopped, false); assert.equal(f.runtime.state.leases.size, 1)
})

test("create ambiguity keeps reservations; unavailable agent/ancestry/preflight fails before prompt", async (t) => {
  const f = await fixture(t, { create: async () => { throw new Error("create acceptance ambiguous") } })
  const id = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id; await settle(f.runtime)
  assert.equal(f.runtime.state.jobs.get(id).state, "blocked"); assert.equal(f.runtime.state.leases.size, 1)
  await f.runtime.cancel(id, "root"); assert.equal(f.runtime.state.jobs.get(id).stopped, false)
  const g = await fixture(t)
  g.input.client.app.agents = async () => ({ data: [] })
  const bad = g.runtime.dispatch(args(), rootContext).task_id; await settle(g.runtime)
  assert.equal(g.runtime.state.jobs.get(bad).state, "failed"); assert.equal(g.calls.length, 0)
  g.runtime.check = async () => { throw new Error("stale receipt hash") }
  const stale = g.runtime.dispatch(args(), rootContext).task_id; await settle(g.runtime)
  assert.match(g.runtime.state.jobs.get(stale).error, /stale/)
})

test("periodic reconciliation, timeout and disposal cleanup do not prompt parents", async (t) => {
  const f = await fixture(t)
  const id = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id; await settle(f.runtime)
  const j = f.runtime.state.jobs.get(id); j.deadline = Date.now() - 1
  await f.runtime.reconcile(); assert.equal(j.state, "cancelled")
  await f.runtime.dispose(); assert.equal(f.runtime.disposed, true); assert.equal(f.runtime.timer._destroyed, true)
  assert.throws(() => f.runtime.dispatch(args(), rootContext), /disposed/)
})

test("permission/error events during acceptance are preserved and queued result tools stay nonblocking", async (t) => {
  const f = await fixture(t)
  f.input.client.session.promptAsync = async (o) => {
    await f.runtime.event({ type: "permission.updated", properties: { sessionID: o.path.id } })
    return { response: { ok: true } }
  }
  const id = f.runtime.dispatch(args(), rootContext).task_id; await settle(f.runtime)
  assert.equal(f.runtime.state.jobs.get(id).state, "waiting_permission")
  let release
  f.input.client.session.messages = async () => { await new Promise((r) => { release = r }); return { data: [] } }
  const out = JSON.parse(await f.runtime.hooks().tool.orchestration_result.execute({ task_id: id }, rootContext))
  assert.equal(out.evidence_ready, false); assert.equal(out.state, "waiting_permission")
  assert.ok(release); release(); await settle(f.runtime)
})

test("root writes are denied until alias reservations become canonical", async (t) => {
  const f = await fixture(t)
  let release
  f.runtime.check = async () => { await new Promise((r) => { release = r }) }
  const id = f.runtime.dispatch(args("worker", ["a.txt"]), rootContext).task_id
  await assert.rejects(() => f.runtime.before({ tool: "write", sessionID: "root" }, { filePath: "b.txt", content: "no early alias race" }), /canonical file reservations/)
  f.runtime.check = async () => {}
  release(); await settle(f.runtime)
  assert.equal(f.runtime.state.jobs.get(id).state, "running")
})

test("worker difficulty selects the fixed static profile and records the requested route", async (t) => {
  for (const [difficulty, profile, modelID, effort] of [
    ["easy", "worker-easy", "gpt-6-luna", "low"],
    ["medium", "worker-medium", "gpt-6-luna", "medium"],
    ["hard", "worker-hard", "gpt-6.1-sol", "high"],
  ]) {
    const f = await fixture(t)
    const id = f.runtime.dispatch(args("worker", [], difficulty), rootContext).task_id
    await settle(f.runtime)
    const j = f.runtime.state.jobs.get(id)
    assert.equal(j.profile, profile)
    assert.equal(j.requestedModel.modelID, modelID)
    assert.equal(j.requestedEffort, effort)
    assert.equal(j.routeStatus, "configured")
    const submitted = f.calls.find(([name]) => name === "promptAsync")?.[1]
    assert.equal(submitted.body.agent, profile)
  }
})

test("worker difficulty is required and non-worker roles cannot request model routes", () => {
  const f = { runtime: new Runtime({ worktree: "/tmp", directory: "/tmp", serverUrl: new URL("http://localhost:1234"), client: {} }, "/tmp", async () => {}, 60000) }
  try {
    assert.throws(() => f.runtime.dispatch({ ...args("worker"), difficulty: undefined }, rootContext), /difficulty/)
    assert.throws(() => f.runtime.dispatch({ ...args("researcher"), difficulty: "hard" }, rootContext), /only for CODE worker/)
    assert.throws(() => f.runtime.dispatch({ ...args("researcher"), escalation_of: "source" }, rootContext), /only for CODE worker/)
  } finally { void f.runtime.dispose() }
})

test("missing or mismatched worker profile fails before session creation", async (t) => {
  const missing = await fixture(t)
  missing.input.client.app.agents = async () => ({ data: [] })
  const missingID = missing.runtime.dispatch(args("worker"), rootContext).task_id
  await settle(missing.runtime)
  assert.equal(missing.runtime.state.jobs.get(missingID).state, "failed")
  assert.equal(missing.calls.filter(([name]) => name === "create").length, 0)

  const mismatch = await fixture(t)
  mismatch.input.client.app.agents = async () => ({ data: [{ name: "worker-easy", mode: "subagent", model: { providerID: "tailnet-llama", modelID: "qwen3.8-flash-next-orca-iq3_m" } }] })
  const mismatchID = mismatch.runtime.dispatch(args("worker"), rootContext).task_id
  await settle(mismatch.runtime)
  const j = mismatch.runtime.state.jobs.get(mismatchID)
  assert.equal(j.state, "failed")
  assert.equal(j.routeStatus, "unavailable")
  assert.equal(mismatch.calls.filter(([name]) => name === "create").length, 0)

  for (const data of [
    [{ name: "worker-easy", mode: "primary", model: { providerID: "openai", modelID: "gpt-6-luna" } }],
    [
      { name: "worker-easy", mode: "subagent", model: { providerID: "openai", modelID: "gpt-6-luna" } },
      { name: "worker-easy", mode: "subagent", model: { providerID: "openai", modelID: "gpt-6-luna" } },
    ],
  ]) {
    const f = await fixture(t)
    f.input.client.app.agents = async () => ({ data })
    const id = f.runtime.dispatch(args("worker"), rootContext).task_id
    await settle(f.runtime)
    assert.equal(f.runtime.state.jobs.get(id).state, "failed")
    assert.equal(f.calls.filter(([name]) => name === "create").length, 0)
  }
})

test("root/parent model gate requires the exact current message and rejects Astra before create", async (t) => {
  const missing = await fixture(t)
  assert.throws(() => missing.runtime.dispatch(args("worker"), { ...rootContext, messageID: "" }), /identity/)
  const missingID = missing.runtime.dispatch(args("worker"), { ...rootContext, messageID: "stale-message" }).task_id
  await settle(missing.runtime)
  assert.equal(missing.runtime.state.jobs.get(missingID).state, "failed")
  assert.equal(missing.calls.filter(([name]) => name === "create").length, 0)

  const astra = await fixture(t)
  astra.messages.set("root", [{ info: { id: rootContext.messageID, sessionID: "root", role: "assistant", providerID: "openai", modelID: "gpt-6-astra" }, parts: [] }])
  const astraID = astra.runtime.dispatch(args("worker"), rootContext).task_id
  await settle(astra.runtime)
  assert.equal(astra.runtime.state.jobs.get(astraID).state, "failed")
  assert.equal(astra.calls.filter(([name]) => name === "create").length, 0)
})

test("effective model mismatch is surfaced and cannot authorize Astra escalation", async (t) => {
  const f = await fixture(t)
  const sourceID = f.runtime.dispatch(args("worker", [], "hard"), rootContext).task_id
  await settle(f.runtime)
  const source = f.runtime.state.jobs.get(sourceID)
  f.messages.set(source.session, [{ info: final(source, { modelID: "gpt-6-luna" }), parts: [{ type: "text", text: "RTS_STATUS: FAILED\nSol route not used" }] }])
  await f.runtime.reconcile()
  assert.equal(source.resultLoaded, true)
  assert.equal(source.routeStatus, "mismatch")
  assert.equal(source.resultStatus, "FAILED")
  assert.throws(() => f.runtime.dispatch({ ...args("worker", [], "hard"), escalation_of: sourceID }, rootContext), /not eligible/)
})

test("one-shot Astra escalation requires reviewed source evidence, same scope, and cannot nest", async (t) => {
  const f = await fixture(t)
  const sourceID = f.runtime.dispatch(args("worker", ["a.txt"], "hard"), rootContext).task_id
  await settle(f.runtime)
  const source = f.runtime.state.jobs.get(sourceID)
  f.messages.set(source.session, [{ info: final(source), parts: [{ type: "text", text: "RTS_STATUS: ESCALATE\nHard Sol completed but did not satisfy acceptance." }] }])
  await f.runtime.reconcile()
  assert.equal(source.resultLoaded, true)
  assert.equal(source.routeStatus, "matched")
  assert.equal(source.resultStatus, "ESCALATE")

  const retryID = f.runtime.dispatch({ ...args("worker", ["a.txt"], "hard"), scope_files: ["b.txt", "a.txt"], escalation_of: sourceID }, rootContext).task_id
  await settle(f.runtime)
  const retry = f.runtime.state.jobs.get(retryID)
  assert.equal(retry.profile, "worker-astra")
  assert.equal(retry.routeStatus, "configured")
  assert.equal(source.escalationReservedBy, retryID)
  assert.equal(f.calls.filter(([name, call]) => name === "promptAsync" && call.body.agent === "worker-astra").length, 1)
  assert.throws(() => f.runtime.dispatch({ ...args("worker", ["a.txt"], "hard"), escalation_of: sourceID }, rootContext), /not eligible/)

  const childMessageID = "astra-child-message"
  f.messages.set(retry.session, [{ info: { id: childMessageID, sessionID: retry.session, role: "assistant", providerID: "openai", modelID: "gpt-6-astra" }, parts: [] }])
  assert.throws(() => f.runtime.dispatch(args("researcher"), { sessionID: retry.session, agent: "worker-astra", messageID: childMessageID }), /leaf-only/)
})

test("pre-create escalation failures release reservation; ambiguous create retains it", async (t) => {
  const mismatch = await fixture(t)
  const sourceID = mismatch.runtime.dispatch(args("worker", ["a.txt"], "hard"), rootContext).task_id
  await settle(mismatch.runtime)
  const source = mismatch.runtime.state.jobs.get(sourceID)
  mismatch.messages.set(source.session, [{ info: final(source), parts: [{ type: "text", text: "RTS_STATUS: FAILED\nNo acceptance" }] }])
  await mismatch.runtime.reconcile()
  mismatch.input.client.app.agents = async () => ({ data: [{ name: "worker-astra", mode: "subagent", model: { providerID: "openai", modelID: "wrong-model" } }] })
  const failedID = mismatch.runtime.dispatch({ ...args("worker", ["a.txt"], "hard"), escalation_of: sourceID }, rootContext).task_id
  await settle(mismatch.runtime)
  assert.equal(mismatch.runtime.state.jobs.get(failedID).state, "failed")
  assert.equal(source.escalationReservedBy, undefined)
  assert.equal(mismatch.calls.filter(([name]) => name === "create").length, 1)

  const ambiguous = await fixture(t)
  const ambiguousSourceID = ambiguous.runtime.dispatch(args("worker", ["a.txt"], "hard"), rootContext).task_id
  await settle(ambiguous.runtime)
  const ambiguousSource = ambiguous.runtime.state.jobs.get(ambiguousSourceID)
  ambiguous.messages.set(ambiguousSource.session, [{ info: final(ambiguousSource), parts: [{ type: "text", text: "RTS_STATUS: FAILED\nNo acceptance" }] }])
  await ambiguous.runtime.reconcile()
  ambiguous.input.client.session.create = async () => { throw new Error("create acceptance ambiguous") }
  const ambiguousID = ambiguous.runtime.dispatch({ ...args("worker", ["a.txt"], "hard"), escalation_of: ambiguousSourceID }, rootContext).task_id
  await settle(ambiguous.runtime)
  assert.equal(ambiguous.runtime.state.jobs.get(ambiguousID).state, "blocked")
  assert.equal(ambiguousSource.escalationReservedBy, ambiguousID)
})

const reviewArgs = (depth = "medium", extra = {}) => ({ ...args("attacker"), review_depth: depth,
  review_target: "proposal", review_revision: "rev-1", review_lens: "trust-boundary", ...extra })
const reasoningReport = "RTS_STATUS: ESCALATE\nRTS_ESCALATION: " + JSON.stringify({ cause: "reasoning_insufficiency",
  bounded_assessment_completed: true, further_reasoning_needed: true, blocked_by: [], unresolved_analysis: "Unresolved bounded path ordering" })
async function reviewedSol(f, report = reasoningReport, info = {}) {
  const id = f.runtime.dispatch(reviewArgs(), rootContext).task_id
  await settle(f.runtime)
  const job = f.runtime.state.jobs.get(id)
  f.messages.set(job.session, [{ info: final(job, info), parts: [{ type: "text", text: report }] }])
  await f.runtime.reconcile()
  return job
}

test("Sol attacker completion is correlated and only approved same-lens one-shot Astra can follow", async (t) => {
  const f = await fixture(t), source = await reviewedSol(f)
  assert.equal(source.state, "completed")
  assert.equal(source.resultStatus, "ESCALATE")
  assert.equal(source.routeStatus, "matched")
  assert.equal(source.resultLoaded, true)
  const hard = reviewArgs("hard", { escalation_of: source.id, review_approved: true })
  for (const patch of [{ review_approved: false }, { review_revision: "rev-2" }, { review_lens: "failure-path" },
    { review_target: "other" }, { scope_files: ["a.txt"] }]) assert.throws(() => f.runtime.dispatch({ ...hard, ...patch }, rootContext))
  const id = f.runtime.dispatch(hard, rootContext).task_id
  assert.throws(() => f.runtime.dispatch(hard, rootContext), /not eligible/)
  await settle(f.runtime)
  const retry = f.runtime.state.jobs.get(id)
  assert.equal(retry.profile, "attacker-hard")
  assert.equal(retry.requestedEffort, undefined)
  assert.equal(retry.state, "running")
  assert.equal(source.escalationReservedBy, retry.id)
  assert.equal(retry.artifactDigest, source.artifactDigest)
  const body = f.calls.find(([op, call]) => op === "promptAsync" && call.path.id === retry.session)[1].body
  assert.match(body.system, new RegExp(source.id))
  assert.throws(() => f.runtime.dispatch(args("researcher"), { sessionID: retry.session, agent: "attacker-hard", messageID: "child" }), /leaf-only/)
  f.messages.set(retry.session, [{ info: final(retry, { modelID: "gpt-6.1-sol" }), parts: [{ type: "text", text: "RTS_STATUS: SUCCESS" }] }])
  await f.runtime.reconcile()
  assert.equal(retry.routeStatus, "mismatch")
  assert.equal(retry.state, "failed")
})

test("attacker rejects direct Astra, missing identity, blockers and wrong actual source model", async (t) => {
  const f = await fixture(t)
  assert.throws(() => f.runtime.dispatch(reviewArgs("hard"), rootContext), /escalation_of/)
  assert.throws(() => f.runtime.dispatch(reviewArgs("low", { review_revision: undefined }), rootContext), /target\/revision\/lens/)
  for (const [report, info] of [["RTS_STATUS: FAILED\npermission denied", {}], ["RTS_STATUS: ESCALATE\nmissing evidence", {}],
    [reasoningReport, { modelID: "gpt-6-luna" }], [reasoningReport, { parentID: "uncorrelated" }]]) {
    const source = await reviewedSol(f, report, info)
    assert.throws(() => f.runtime.dispatch(reviewArgs("hard", { escalation_of: source.id, review_approved: true }), rootContext), /not eligible/)
  }
})

test("attacker revision bytes and route effort checked before create, failed retry stays consumed", async (t) => {
  const f = await fixture(t), source = await reviewedSol(f)
  await writeFile(path.join(f.root, "a.txt"), "changed revision")
  const hard = reviewArgs("hard", { escalation_of: source.id, review_approved: true })
  const id = f.runtime.dispatch(hard, rootContext).task_id
  await settle(f.runtime)
  assert.match(f.runtime.state.jobs.get(id).error, /artifact revision changed/)
  assert.equal(f.calls.filter(([op]) => op === "create").length, 1)
  assert.equal(source.escalationReservedBy, id)
  assert.throws(() => f.runtime.dispatch(hard, rootContext), /not eligible/)
  for (const agent of [{ name: "attacker-medium", mode: "subagent", model: "openai/gpt-6.1-sol", reasoningEffort: "low" },
    { name: "attacker-medium", mode: "subagent", model: "openai/gpt-6-luna", reasoningEffort: "high" }]) {
    const g = await fixture(t)
    g.input.client.app.agents = async () => ({ data: [agent] })
    const failed = g.runtime.dispatch(reviewArgs(), rootContext).task_id
    await settle(g.runtime)
    assert.equal(g.runtime.state.jobs.get(failed).state, "failed")
    assert.equal(g.calls.length, 0)
  }
})

test("independent lenses share revision, run together and queue results without advancing parent", async (t) => {
  const f = await fixture(t)
  const ids = ["trust-boundary", "failure-path"].map((lens) => f.runtime.dispatch(reviewArgs("low", { review_lens: lens }), rootContext).task_id)
  await settle(f.runtime)
  const jobs = ids.map((id) => f.runtime.state.jobs.get(id))
  assert.equal(jobs[0].artifactDigest, jobs[1].artifactDigest)
  assert.ok(jobs.every((j) => j.state === "running"))
  for (const j of jobs) f.messages.set(j.session, [{ info: final(j), parts: [{ type: "text", text: "RTS_STATUS: SUCCESS\nNO_FINDINGS within assessed lens" }] }])
  await f.runtime.reconcile()
  assert.ok(jobs.every((j) => j.resultLoaded && j.state === "completed"))
  assert.equal(f.runtime.state.queues.get("root").size, 2)
  assert.equal(f.calls.filter(([op]) => op === "promptAsync").length, 2)
})

test("Astra attacker configured model and provider-default effort are mandatory before create", async (t) => {
  for (const patch of [{ model: "openai/gpt-6.1-sol" }, { reasoningEffort: "high" }, { mode: "primary" }]) {
    const f = await fixture(t), source = await reviewedSol(f)
    f.input.client.app.agents = async () => ({ data: [{ name: "attacker-hard", mode: "subagent", model: "openai/gpt-6-astra", ...patch }] })
    const id = f.runtime.dispatch(reviewArgs("hard", { escalation_of: source.id, review_approved: true }), rootContext).task_id
    await settle(f.runtime)
    assert.equal(f.runtime.state.jobs.get(id).state, "failed")
    assert.equal(source.escalationReservedBy, id)
    assert.equal(f.calls.filter(([op]) => op === "create").length, 1) // Sol source only.
  }
})
