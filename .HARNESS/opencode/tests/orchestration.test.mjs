import assert from "node:assert/strict"
import test from "node:test"
import path from "node:path"
import { mkdtemp, mkdir, writeFile, rm, symlink } from "node:fs/promises"
import orchestration from "../plugins/orchestration.ts"

const { State, Runtime, wrapper, operationFiles, envelope, canonical } = orchestration.testing
const handoff = "<Γ>\nSIGMA/1: PIDGEN/text\nRole: researcher\nTarget: inspect\nAcceptance: evidence\nScope: a.txt\nKnown facts: N_A\nSIGMA/1\n</Γ>"
const reserve = (s, parent, role, id, writes = [], scope = []) => s.reserve(parent, role, scope, writes, id, `msg_${id}`, 0, 900000)
const args = (role = "researcher", writes = []) => ({ role, handoff, scope_files: ["a.txt", "b.txt"], write_files: writes })
const rootContext = { sessionID: "root", agent: "orchestrator" }
const tick = () => new Promise((resolve) => setTimeout(resolve, 1))
async function settle(runtime) { for (let i = 0; i < 100 && runtime.pending.size; i++) await tick(); assert.equal(runtime.pending.size, 0) }

async function fixture(t, overrides = {}) {
  const root = await mkdtemp("A:/TMP/opencode/rts-test-")
  await writeFile(path.join(root, "a.txt"), "a"); await writeFile(path.join(root, "b.txt"), "b")
  let n = 0
  const sessions = new Map([["root", { id: "root" }]])
  const messages = new Map()
  const calls = [], aborts = []
  const input = { worktree: root, directory: root, serverUrl: new URL("http://localhost:1234"), client: {
    app: { agents: async () => ({ data: ["researcher", "proposer", "worker", "verifier", "auditor"].map((name) => ({ name, mode: "subagent" })) }) },
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

function final(j, extra = {}) { return { id: `msg_final_${j.id}`, sessionID: j.session, parentID: j.promptID, role: "assistant", time: { created: 1, completed: 2 }, finish: "stop", ...extra } }

test("exact wrapper and explicit HTTP envelope validation fail closed", () => {
  wrapper(handoff)
  for (const bad of [handoff + "\n", handoff.replace("Known facts:", "Facts:"), handoff.replace("\nSIGMA/1\n</Γ>", "\n</Γ>"), handoff.replace("Scope: a.txt", "Scope: a.txt\nScope: b.txt")]) assert.throws(() => wrapper(bad))
  for (const r of [undefined, { error: { code: 400 } }, { response: { ok: false } }, {}]) assert.throws(() => envelope(r, "op"))
  assert.equal(envelope({ response: { ok: true } }, "promptAsync", false), undefined)
})

test("pure reservations enforce max6 including grandchildren, root max4 and immutable ownership", () => {
  const s = new State()
  for (let i = 0; i < 4; i++) { const j = reserve(s, "root", "researcher", `c${i}`, [], ["a"]); j.session = `ses${i}`; s.sessions.set(j.session, j) }
  assert.throws(() => reserve(s, "root", "researcher", "fifth"), /slots/)
  const g = reserve(s, "ses0", "researcher", "grand", [], ["a"]); g.session = "grand-session"; s.sessions.set(g.session, g)
  assert.equal(g.depth, 2)
  assert.throws(() => reserve(s, "grand-session", "researcher", "depth3"), /eligible/)
  assert.throws(() => reserve(s, "ses1", "worker", "write-grand", ["a"], ["a"]), /read-only/)
  assert.throws(() => reserve(s, "ses1", "researcher", "escape", [], ["b"]), /inherited/)
  reserve(s, "ses1", "verifier", "sixth", [], ["a"])
  assert.throws(() => reserve(s, "ses2", "researcher", "seventh", [], ["a"]), /slots/)
  assert.throws(() => reserve(s, "root", "guessed", "bad"), /allowlist/)
  const own = new State(); reserve(own, "root", "worker", "w1", ["a"], ["a"])
  assert.throws(() => reserve(own, "root", "worker", "w2", ["a"], ["a"]), /overlapping/)
  assert.throws(() => reserve(own, "root", "auditor", "w3", ["b"], ["b"]), /only worker/)
  const w = own.jobs.get("w1"); w.session = "w"; own.sessions.set("w", w)
  assert.throws(() => reserve(own, "w", "researcher", "leaf"), /eligible/)
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
  const grandID = f.runtime.dispatch(args("verifier"), { sessionID: child.session, agent: "researcher" }).task_id
  await settle(f.runtime)
  const grand = f.runtime.state.jobs.get(grandID)
  assert.equal(grand.depth, 2); assert.equal(f.sessions.get(grand.session).parentID, child.session)
  assert.equal(f.calls.find(([name, o]) => name === "promptAsync" && o.path.id === grand.session)[1].body.messageID, grand.promptID)
  assert.ok(f.checks() >= 4)
  assert.ok(f.calls.filter(([name]) => name === "promptAsync").every(([, o]) => o.path.id !== "root" && !Object.hasOwn(o.body, "noReply")))
  assert.throws(() => f.runtime.dispatch(args(), { sessionID: grand.session, agent: "verifier" }), /eligible/)
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
  await before("edit", { filePath: "A.txt", oldString: "a", newString: "z" })
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
  assert.throws(() => f.runtime.dispatch(args("worker", ["A.txt"]), rootContext), /overlapping/)
})

test("junction aliases and junction changes are resolved on every use", async (t) => {
  const f = await fixture(t)
  const one = path.join(f.root, "one"), two = path.join(f.root, "two"), alias = path.join(f.root, "alias")
  await mkdir(one); await mkdir(two); await writeFile(path.join(one, "file.txt"), "one"); await writeFile(path.join(two, "file.txt"), "two")
  await symlink(one, alias, "junction")
  assert.equal(await canonical(f.root, "alias/file.txt"), await canonical(f.root, "one/file.txt"))
  const spec = { role: "worker", handoff, scope_files: ["alias/file.txt"], write_files: ["alias/file.txt"] }
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
  const gID = f.runtime.dispatch(args("verifier"), { sessionID: c.session, agent: "researcher" }).task_id; await settle(f.runtime)
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
