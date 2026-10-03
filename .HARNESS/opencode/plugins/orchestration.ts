import { tool } from "@opencode-ai/plugin"
import type { PluginInput, Hooks } from "@opencode-ai/plugin"
import { randomBytes, createHash } from "node:crypto"
import { lstat, realpath, readFile, rm } from "node:fs/promises"
import { spawn } from "node:child_process"
import path from "node:path"
import { fileURLToPath } from "node:url"
import { containsCredentialLikeContent } from "../tools/agy_worker.ts"

const harness = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../..")
const roles = Object.freeze(["explorer", "researcher", "proposer", "auditor", "verifier", "refactorer", "critic", "blackhat", "compressor", "worker"])
const delegators = new Set(["researcher", "proposer"])
const key = (value: string) => path.resolve(value).toLowerCase()
const inside = (root: string, target: string) => key(target) === key(root) || key(target).startsWith(key(root) + path.sep)
const messageID = () => `msg_${Date.now().toString(16).padStart(12, "0")}${randomBytes(7).toString("hex")}`
type Job = { id: string; parent: string; root: string; depth: number; role: string; scope: readonly string[]; writes: readonly string[]; leases: string[]; session?: string; promptID: string; state: string; error?: string; result?: string; resultLoaded?: boolean; deadline: number; stopped: boolean; cancelRequested: boolean; pendingLaunch: boolean; abortConfirmed?: boolean }

// Domain state is synchronous. All reservations happen before control yields to IO.
class State {
  jobs = new Map<string, Job>()
  sessions = new Map<string, Job>()
  leases = new Map<string, string>()
  queues = new Map<string, Map<string, { task_id: string; state: string }>>()
  reserve(parent: string, role: string, scope: string[], writes: string[], id: string, promptID: string, now: number, timeout: number) {
    const ancestor = this.sessions.get(parent)
    if (ancestor && (ancestor.stopped || ancestor.cancelRequested || ancestor.depth >= 2 || !delegators.has(ancestor.role))) throw new Error("parent is not an eligible live delegator")
    if (!roles.includes(role)) throw new Error("role is outside explicit allowlist")
    if (writes.length && role !== "worker") throw new Error("only worker may own writes")
    if (ancestor && (writes.length || scope.some((p) => !ancestor.scope.includes(p)))) throw new Error("grandchild scope must be inherited read-only")
    const active = [...this.jobs.values()].filter((j) => !j.stopped)
    if (active.length >= 6 || (!ancestor && active.filter((j) => j.depth === 1).length >= 4)) throw new Error("descendant slots exhausted (two reserved for grandchildren)")
    for (const p of writes) if (this.leases.has(p)) throw new Error("overlapping file ownership")
    const j: Job = { id, parent, root: ancestor?.root ?? parent, depth: (ancestor?.depth ?? 0) + 1, role, scope: Object.freeze([...scope]), writes: Object.freeze([...writes]), leases: [...writes], promptID, state: "reserved", deadline: now + timeout, stopped: false, cancelRequested: false, pendingLaunch: true }
    this.jobs.set(id, j)
    for (const p of writes) this.leases.set(p, id)
    return j
  }
  canonicalize(j: Job, scope: string[], writes: string[]) {
    for (const p of writes) if (this.leases.has(p) && this.leases.get(p) !== j.id) throw new Error("realpath alias overlaps another worker")
    const ancestor = this.sessions.get(j.parent)
    if (ancestor && scope.some((p) => !ancestor.scope.includes(p))) throw new Error("canonical scope escapes inherited scope")
    for (const p of j.leases) if (this.leases.get(p) === j.id) this.leases.delete(p)
    j.scope = Object.freeze(scope); j.writes = Object.freeze(writes); j.leases = [...writes]
    for (const p of writes) this.leases.set(p, j.id)
  }
  notify(j: Job) {
    const q = this.queues.get(j.parent) ?? new Map()
    q.set(j.id, { task_id: j.id, state: j.state }); this.queues.set(j.parent, q)
  }
  stop(j: Job, state: string) {
    if (j.stopped) return
    j.state = state; j.stopped = true
    for (const p of j.leases) if (this.leases.get(p) === j.id) this.leases.delete(p)
    this.notify(j)
  }
  observe(j: Job, info: any, parts: any[] = []) {
    if (j.stopped || info?.sessionID !== j.session || info?.role !== "assistant" || info?.parentID !== j.promptID) return
    if (!info.time?.completed || (!info.error && info.finish !== "stop")) return
    const text = parts.filter((p) => p.type === "text" && typeof p.text === "string").map((p) => p.text).join("\n").slice(0, 64_000)
    j.result = containsCredentialLikeContent(text) ? "FAILED: protected returned evidence rejected" : text
    if (info.error) j.error = "assigned prompt assistant returned an error"
    this.stop(j, info.error ? "failed" : "completed")
  }
  descendants(session: string): Job[] {
    const direct = [...this.jobs.values()].filter((j) => j.parent === session)
    return direct.flatMap((j) => [...(j.session ? this.descendants(j.session) : []), j])
  }
}

function wrapper(text: string) {
  if (typeof text !== "string" || text.length > 20_000 || !/^<Γ>\r?\nSIGMA\/1: PIDGEN\/text\r?\nRole: [\s\S]+\r?\nTarget: [\s\S]+\r?\nAcceptance: [\s\S]+\r?\nScope: [\s\S]+\r?\nKnown facts: [\s\S]+\r?\nSIGMA\/1\r?\n<\/Γ>$/.test(text)) throw new Error("invalid exact SIGMA wrapper/ordered fields")
  if ((text.match(/^Role:|^Target:|^Acceptance:|^Scope:|^Known facts:/gm) ?? []).length !== 5) throw new Error("duplicate handoff fields")
}

async function canonical(root: string, value: string) {
  if (typeof value !== "string" || !value || /[\x00-\x1f*?]/.test(value) || value.replace(/^[a-z]:[\\/]/i, "").includes(":")) throw new Error("invalid file path")
  const absolute = path.resolve(root, value)
  if (!inside(root, absolute) || key(root) === key(absolute)) throw new Error("path escapes worktree or is a directory")
  let cursor = absolute, suffix: string[] = []
  while (true) {
    try {
      const info = await lstat(cursor)
      if (cursor === absolute && !info.isFile() && !info.isSymbolicLink()) throw new Error("non-file target")
      break
    } catch (e: any) {
      if (e.code !== "ENOENT") throw e
      suffix.unshift(path.basename(cursor)); const parent = path.dirname(cursor)
      if (parent === cursor) throw new Error("no existing path ancestor")
      cursor = parent
    }
  }
  const resolved = path.join(await realpath(cursor), ...suffix)
  if (!suffix.length && !(await lstat(resolved)).isFile()) throw new Error("realpath target is not a file")
  if (resolved.split(/[\\/]/).some((segment) => /(^\.env(?:\..*)?$|credential|secret|token|\.pem$|\.key$|^id_rsa|^id_ed25519|\.pfx$|\.p12$|\.kdbx$)/i.test(segment))) throw new Error("protected realpath scope")
  if (!inside(await realpath(root), resolved)) throw new Error("realpath escapes worktree")
  return key(resolved)
}

// Deliberately accept only the installed native shapes; unknown operations are denied.
function operationFiles(name: string, args: any): { files: string[]; write: boolean } {
  if (!args || typeof args !== "object" || Array.isArray(args)) throw new Error("malformed tool arguments")
  const exact = (required: string[], optional: string[] = []) => {
    if (required.some((p) => typeof args[p] !== "string") || Object.keys(args).some((p) => ![...required, ...optional].includes(p))) throw new Error("unrecognized/malformed operation")
    if (Object.hasOwn(args, "replaceAll") && typeof args.replaceAll !== "boolean" || ["offset", "limit"].some((p) => Object.hasOwn(args, p) && (!Number.isInteger(args[p]) || args[p] < 0)) || Object.hasOwn(args, "include") && typeof args.include !== "string") throw new Error("malformed optional arguments")
  }
  if (name === "read") { exact(["filePath"], ["offset", "limit"]); return { files: [args.filePath], write: false } }
  if (name === "grep") { exact(["pattern", "path"], ["include"]); return { files: [args.path], write: false } }
  if (name === "edit") { exact(["filePath", "oldString", "newString"], ["replaceAll"]); return { files: [args.filePath], write: true } }
  if (name === "write") { exact(["filePath", "content"]); return { files: [args.filePath], write: true } }
  if (name !== "apply_patch") throw new Error("managed sessions deny native task, process, external and unknown tools")
  exact(["patchText"])
  const lines = args.patchText.replaceAll("\r\n", "\n").trimEnd().split("\n")
  if (lines.shift() !== "*** Begin Patch" || lines.pop() !== "*** End Patch") throw new Error("malformed patch envelope")
  const files: string[] = []; let action = "", moved = false, body = false
  for (const line of lines) {
    const header = /^\*\*\* (Add|Update|Delete) File: (.+)$/.exec(line)
    if (header) { files.push(header[2]); action = header[1]; moved = false; body = false; continue }
    const move = /^\*\*\* Move to: (.+)$/.exec(line)
    if (move && action === "Update" && !moved && !body) { files.push(move[1]); moved = true; continue }
    if (!action || action === "Delete" || line.startsWith("***") && line !== "*** End of File") throw new Error("unrecognized patch operation")
    if (action === "Add" ? !line.startsWith("+") : !/^(?:@@(?: .*?)?|[ +\-].*|\*\*\* End of File)$/.test(line)) throw new Error("malformed patch body")
    body = true
  }
  if (!files.length) throw new Error("empty patch")
  return { files, write: true }
}

function envelope(response: any, operation: string, requireData = true) {
  if (!response || response.error || response.response?.ok === false) throw new Error(`${operation}: explicit HTTP error envelope`)
  if (requireData && response.data === undefined) throw new Error(`${operation}: missing response data`)
  return response.data
}

async function preflight(role: string) {
  const receipt = path.join("A:/TMP/opencode", `rts-${randomBytes(12).toString("hex")}.json`)
  const script = path.join(harness, "scripts/invoke-rule-preflight.ps1")
  const run = (verify = false) => new Promise<void>((resolve, reject) => {
    const command = verify ? `& '${script.replaceAll("'", "''")}' -VerifyReceipt '${receipt}'` : `& '${script.replaceAll("'", "''")}' -Task change -Tag @('code','subagent','worker') -OutFile '${receipt}'`
    const child = spawn("pwsh.exe", ["-NoProfile", "-NonInteractive", "-Command", command], { stdio: "ignore", windowsHide: true, shell: false })
    const timer = setTimeout(() => { child.kill(); reject(new Error("local preflight timeout")) }, 30_000)
    child.once("error", () => { clearTimeout(timer); reject(new Error("local preflight unavailable")) })
    child.once("exit", (code) => { clearTimeout(timer); code === 0 ? resolve() : reject(new Error("local preflight failed")) })
  })
  try {
    await run(); await run(true)
    const r = JSON.parse(await readFile(receipt, "utf8"))
    if (r.schema !== "lrf-preflight/v1" || r.status !== "PASS" || key(r.root) !== key(path.join(harness, "rules/BOOTSTRAP.lrf")) || !Array.isArray(r.files) || !r.files.length || !["code", "worker", "subagent"].every((tag) => r.tags.includes(tag))) throw new Error("invalid preflight receipt")
    for (const file of r.files) {
      if (!inside(harness, file.path) || !/^[a-f0-9]{64}$/i.test(file.sha256)) throw new Error("escaped/malformed receipt reference")
      const actual = createHash("sha256").update(await readFile(file.path)).digest("hex")
      if (actual !== file.sha256.toLowerCase()) throw new Error("stale receipt hash")
    }
    // Canonical bounded role reads follow hash verification, never guessed model selection.
    await readFile(path.join(harness, "agents", `${role}.md`), "utf8")
  } finally { await rm(receipt, { force: true }) }
}

class Runtime {
  state = new State()
  disposed = false
  registered = false
  reconciling = false
  pollCursor = 0
  rootAgents = new Map<string, string>()
  errors: string[] = []
  timer: ReturnType<typeof setInterval>
  pending = new Set<Promise<any>>()
  input: PluginInput
  root: string
  check: typeof preflight
  constructor(input: PluginInput, root: string, check = preflight, interval = 2000) {
    this.input = input; this.root = root; this.check = check
    this.timer = setInterval(() => this.background(this.reconcile()), interval); this.timer.unref()
  }
  background(p: Promise<any>) {
    this.pending.add(p)
    void p.catch((e) => { this.errors.push(String(e.message).slice(0, 300)); this.errors = this.errors.slice(-20) }).finally(() => this.pending.delete(p))
  }
  async call(operation: string, options: any, requireData = true) {
    const controller = new AbortController()
    let timer: ReturnType<typeof setTimeout>
    try {
      const timeout = new Promise<never>((_, reject) => { timer = setTimeout(() => { controller.abort(); reject(new Error(`${operation}: transport timeout; acceptance ambiguous`)) }, 15_000) })
      const namespace = operation === "agents" ? this.input.client.app : this.input.client.session
      return await Promise.race([namespace[operation]({ ...options, query: { ...options.query, directory: this.root }, signal: controller.signal }).then((r: any) => envelope(r, operation, requireData)), timeout])
    } finally { clearTimeout(timer!) }
  }
  dispatch(args: any, context: any) {
    if (this.disposed) throw new Error("scheduler disposed")
    wrapper(args.handoff)
    if (containsCredentialLikeContent(args.handoff)) throw new Error("protected handoff text")
    if (args.timeout_ms !== undefined && (!Number.isInteger(args.timeout_ms) || args.timeout_ms < 1000 || args.timeout_ms > 900000)) throw new Error("invalid timeout")
    if (!Array.isArray(args.scope_files) || !Array.isArray(args.write_files) || args.scope_files.length > 40 || args.write_files.length > 40) throw new Error("invalid bounded scope")
    const lexical = (values: string[]) => values.map((p) => {
      if (typeof p !== "string" || !p || !inside(this.root, path.resolve(this.root, p))) throw new Error("invalid scope")
      return key(path.resolve(this.root, p))
    })
    const scope = lexical(args.scope_files), writes = lexical(args.write_files)
    if (new Set(scope).size !== scope.length || new Set(writes).size !== writes.length || writes.some((p) => !scope.includes(p))) throw new Error("duplicate/uninherited write scope")
    const parentJob = this.state.sessions.get(context.sessionID)
    if (!parentJob && context.agent !== "orchestrator") throw new Error("only root orchestrator can initiate dispatch")
    const j = this.state.reserve(context.sessionID, args.role, scope, writes, `rts_${randomBytes(12).toString("hex")}`, messageID(), Date.now(), Math.min(Math.max(args.timeout_ms ?? 300_000, 1000), 900_000))
    this.background(this.launch(j, structuredClone(args)))
    return { task_id: j.id, state: j.state, wake_mode: "queued-next-natural-parent-turn" }
  }
  async launch(j: Job, args: any) {
    try {
      await this.check(j.role)
      const parent = await this.call("get", { path: { id: j.parent } })
      if (parent.id !== j.parent || (!this.state.sessions.has(j.parent) && parent.parentID)) throw new Error("unknown parent worker/ancestry")
      if (this.state.sessions.has(j.parent)) {
        const p = this.state.sessions.get(j.parent)!
        if (parent.parentID !== p.parent || p.stopped || p.cancelRequested) throw new Error("SDK ancestry mismatch or stopped parent")
      }
      const agents = await this.call("agents", {})
      if (!Array.isArray(agents)) throw new Error("app.agents malformed data")
      const agent = agents.find((a: any) => a.name === j.role && a.mode !== "primary" && !a.disable)
      if (!agent) throw new Error("requested agent unavailable in live app.agents")
      const scope = await Promise.all(args.scope_files.map((p: string) => canonical(this.root, p)))
      const writes = await Promise.all(args.write_files.map((p: string) => canonical(this.root, p)))
      if (new Set(scope).size !== scope.length || new Set(writes).size !== writes.length) throw new Error("case/realpath aliases in scope")
      this.state.canonicalize(j, scope, writes)
      if (j.cancelRequested || this.disposed) { this.state.stop(j, "cancelled"); return }
      j.state = "creating"
      const created = await this.call("create", { body: { parentID: j.parent, title: `RTS ${j.role} ${j.id}` } })
      if (typeof created.id !== "string") throw new Error("session.create missing ID; acceptance ambiguous")
      j.session = created.id; this.state.sessions.set(created.id, j)
      const actual = await this.call("get", { path: { id: created.id } })
      if (actual.id !== created.id || actual.parentID !== j.parent) throw new Error("SDK child identity/parentID mismatch")
      if (j.cancelRequested || this.disposed) { await this.abortJob(j); return }
      // Re-verify immediately before dispatch, after SDK/filesystem awaits.
      await this.check(j.role)
      if (j.cancelRequested || this.disposed) { await this.abortJob(j); return }
      j.state = "submitting"
      await this.call("promptAsync", { path: { id: j.session }, body: { messageID: j.promptID, agent: j.role, system: `Managed RTS depth ${j.depth}. Exact inherited file scope: ${JSON.stringify(j.scope)}. Write ownership: ${JSON.stringify(j.writes)}. Native task, process and external tools are denied. Only read-only researcher/proposer at depth 1 may use orchestration_dispatch for inherited read-only grandchildren. Parent handles build/test. Report evidence/failures; never infer effects from permissions.`, parts: [{ type: "text", text: args.handoff }] } }, false)
      if (j.state === "submitting" && !j.cancelRequested) j.state = "running"
    } catch (e: any) {
      j.error = String(e.message).slice(0, 500)
      if (j.stopped) this.state.notify(j)
      else if (["creating", "submitting"].includes(j.state) || j.session) { j.state = "blocked"; this.state.notify(j) }
      else this.state.stop(j, j.cancelRequested ? "cancelled" : "failed")
    } finally { j.pendingLaunch = false; if (j.cancelRequested && !j.stopped && j.session) await this.abortJob(j) }
  }
  async reconcile() {
    if (this.reconciling || this.disposed) return
    this.reconciling = true
    try {
      const candidates = [...this.state.jobs.values()].filter((j) => !j.stopped || !j.resultLoaded && j.session)
      const start = candidates.length ? this.pollCursor % candidates.length : 0
      const bounded = [...candidates.slice(start), ...candidates.slice(0, start)].slice(0, 12)
      this.pollCursor = start + bounded.length
      for (const j of bounded) {
        if (!j.stopped && Date.now() >= j.deadline) await this.cancel(j.id, j.parent)
        if (!j.session) continue
        // Poll bounded recent messages, not idle-as-success. Correlation survives reordered events.
        const messages = await this.call("messages", { path: { id: j.session }, query: { directory: this.root, limit: 100 } })
        if (!Array.isArray(messages)) throw new Error("session.messages malformed data")
        for (const m of messages) this.state.observe(j, m.info, m.parts)
        if (j.stopped && !j.resultLoaded) {
          const final = messages.find((m: any) => m.info.parentID === j.promptID && m.info.time?.completed && (m.info.finish === "stop" || m.info.error))
          if (final) {
            const text = final.parts.filter((p: any) => p.type === "text").map((p: any) => p.text).join("\n").slice(0, 64_000)
            j.result = containsCredentialLikeContent(text) ? "FAILED: protected returned evidence rejected" : text
            j.resultLoaded = true
          }
          else if (j.state === "cancelled") j.resultLoaded = true
        }
      }
    } finally { this.reconciling = false }
  }
  async event(event: any) {
    if (event.type === "server.instance.disposed") { await this.dispose(); return }
    if (event.type === "message.updated") {
      const info = event.properties?.info, j = this.state.sessions.get(info?.sessionID)
      if (j) { this.state.observe(j, info); if (j.stopped) this.background(this.reconcile()) }
    }
    const sid = event.properties?.sessionID, j = this.state.sessions.get(sid)
    if (!j || j.stopped) return
    if (event.type === "permission.updated") j.state = "waiting_permission"
    if (event.type === "permission.replied" && j.state === "waiting_permission") j.state = "running"
    if (event.type === "session.error") { j.state = "blocked"; j.error = "session error; await correlated message/confirmed stop"; this.state.notify(j) }
    if (event.type === "session.idle" || event.type === "session.status") this.background(this.reconcile())
  }
  async abortJob(j: Job) {
    if (j.stopped || j.pendingLaunch && !j.session) return
    if (!j.session) { j.state = "blocked"; j.error = "no confirmed session ID; reservation retained"; this.state.notify(j); return }
    j.state = "cancelling"
    try {
      const stopped = await this.call("abort", { path: { id: j.session } })
      if (stopped !== true) throw new Error("abort did not confirm stop")
      j.abortConfirmed = true
      const statuses = await this.call("status", {})
      if (statuses[j.session]?.type !== "idle") throw new Error("abort accepted but idle stop not confirmed")
      this.state.stop(j, "cancelled")
    } catch (e: any) { if (!j.stopped) j.state = "blocked"; j.error = String(e.message); this.state.notify(j) }
  }
  authorized(id: string, session: string) {
    const j = this.state.jobs.get(id)
    if (!j || !(j.parent === session || j.root === session)) throw new Error("task is not owned by caller")
    return j
  }
  async cancel(id: string, session: string) {
    const j = this.authorized(id, session)
    const jobs = [...(j.session ? this.state.descendants(j.session) : []), j]
    for (const item of jobs) item.cancelRequested = true
    for (const item of jobs) if (!item.pendingLaunch) await this.abortJob(item)
    return this.view(j)
  }
  view(j: Job) { return { task_id: j.id, session_id: j.session, parent_id: j.parent, depth: j.depth, role: j.role, state: j.state, confirmed_stopped: j.stopped, error: j.error } }
  async before(input: any, args: any) {
    const j = this.state.sessions.get(input.sessionID)
    if (["orchestration_status", "orchestration_result", "orchestration_cancel"].includes(input.tool)) return
    if (input.tool === "orchestration_dispatch" && j && delegators.has(j.role) && j.depth < 2 && !j.stopped && !j.cancelRequested) return
    if (j) {
      if (j.stopped || j.cancelRequested || ["blocked", "cancelling"].includes(j.state)) throw new Error("managed session is stopped/blocked")
      const op = operationFiles(input.tool, args)
      for (const file of op.files) {
        const p = await canonical(this.root, file)
        if (!j.scope.includes(p) || op.write && (!j.writes.includes(p) || this.state.leases.get(p) !== j.id)) throw new Error("file outside immutable session scope/lease")
      }
      return
    }
    // Unrecognized child sessions (including native task) never get a write bypass.
    const actual = await this.call("get", { path: { id: input.sessionID } })
    if (actual.id !== input.sessionID || actual.parentID || this.rootAgents.get(input.sessionID) !== "orchestrator") throw new Error("unknown worker/root identity: fail closed")
    if (input.tool === "task") throw new Error("native task bypass denied; use RTS dispatch")
    if (this.state.leases.size && !["read", "glob", "grep", "list", "orchestration_dispatch", "edit", "write", "apply_patch"].includes(input.tool)) throw new Error("opaque root operation denied while worker leases are active")
    if (["edit", "write", "apply_patch"].includes(input.tool)) {
      if ([...this.state.jobs.values()].some((job) => job.pendingLaunch && job.writes.length)) throw new Error("root writes paused until canonical file reservations settle")
      const op = operationFiles(input.tool, args)
      for (const file of op.files) if (this.state.leases.has(await canonical(this.root, file))) throw new Error("root cannot write leased worker file")
    }
  }
  async dispose() {
    if (this.disposed) return
    this.disposed = true; clearInterval(this.timer)
    for (const j of this.state.jobs.values()) if (!j.stopped) j.cancelRequested = true
    await Promise.allSettled([...this.pending])
    for (const j of [...this.state.jobs.values()].sort((a, b) => b.depth - a.depth)) if (!j.stopped) await this.abortJob(j)
  }
  hooks(): Hooks {
    const encode = async (fn: () => any) => { try { return JSON.stringify(await fn()) } catch (e: any) { return JSON.stringify({ status: "FAILED", error: String(e.message) }) } }
    const id = { task_id: tool.schema.string() }
    return {
      tool: {
        orchestration_dispatch: tool({ description: "Reserve a bounded child task and return its ID immediately. Observe asynchronously; no automatic parent prompts.", args: { role: tool.schema.enum(roles as any), handoff: tool.schema.string(), scope_files: tool.schema.array(tool.schema.string()).max(40), write_files: tool.schema.array(tool.schema.string()).max(40).default([]), timeout_ms: tool.schema.number().int().min(1000).max(900000).default(300000) }, execute: async (args, context) => encode(() => this.dispatch(args, context)) }),
        orchestration_status: tool({ description: "Observe owned task state and deduplicated queued notifications (does not wait for completion).", args: { task_id: tool.schema.string().optional() }, execute: async (args, c) => encode(() => ({ tasks: args.task_id ? [this.view(this.authorized(args.task_id, c.sessionID))] : [...this.state.jobs.values()].filter((j) => j.parent === c.sessionID || j.root === c.sessionID).map((j) => this.view(j)), notifications: [...(this.state.queues.get(c.sessionID)?.values() ?? [])], observer_errors: this.errors })) }),
        orchestration_result: tool({ description: "Read correlated final result snapshot; pending is not successful evidence. Wait for every requested task before final synthesis.", args: id, execute: async (args, c) => encode(() => { const j = this.authorized(args.task_id, c.sessionID); this.background(this.reconcile()); return { ...this.view(j), evidence_ready: j.resultLoaded === true, result: j.stopped && j.resultLoaded ? j.result : undefined } }) }),
        orchestration_cancel: tool({ description: "Cancel task subtree; retain leases on abort/stop uncertainty.", args: id, execute: async (args, c) => encode(() => this.cancel(args.task_id, c.sessionID)) }),
      },
      event: async ({ event }) => this.event(event),
      "chat.message": async (input, output) => { const agent = input.agent ?? output.message.agent; if (agent) this.rootAgents.set(input.sessionID, agent) },
      "chat.params": async (input) => { this.rootAgents.set(input.sessionID, input.agent) },
      "tool.execute.before": async (input, output) => this.before(input, output.args),
      "experimental.chat.system.transform": async (input, output) => {
        if (!input.sessionID) return
        const q = this.state.queues.get(input.sessionID)
        if (q?.size) { output.system.push(`RTS queued notifications (check result tools; wait-all before synthesis): ${JSON.stringify([...q.values()])}`); q.clear() }
      },
      dispose: async () => this.dispose(),
    }
  }
}

const registrySymbol = Symbol.for("harness.rts.orchestrators.v1")
const registry: Map<string, Promise<Runtime>> = (globalThis as any)[registrySymbol] ??= new Map()
async function orchestration(input: PluginInput): Promise<Hooks> {
  const root = await realpath(input.worktree || input.directory)
  const identity = `${input.serverUrl.toString()}|${key(root)}`
  let promise = registry.get(identity)
  if (!promise) { promise = Promise.resolve(new Runtime(input, root)); registry.set(identity, promise) }
  const runtime = await promise
  if (runtime.registered) return {}
  runtime.registered = true
  return runtime.hooks()
}
// One plugin export only: auto-discovery must not interpret test helpers as plugins.
orchestration.testing = { State, Runtime, wrapper, canonical, operationFiles, envelope, registry, preflight }
export default orchestration
