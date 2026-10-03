import { tool } from "@opencode-ai/plugin"
import { spawn } from "node:child_process"
import { createHash } from "node:crypto"
import { constants, lstat, mkdir, mkdtemp, open, readFile, readdir, realpath, rm, stat, writeFile } from "node:fs/promises"
import os from "node:os"
import path from "node:path"

const MAX_FILES = 40
const MAX_FILE_BYTES = 512 * 1024
const MAX_TOTAL_BYTES = 2 * 1024 * 1024
const MAX_OUTPUT_BYTES = 3 * 1024 * 1024
const EXEC_TIMEOUT_MS = 310_000
const ROUTES = {
  easy: { model: "gpt-5.6-luna", effort: "low" },
  medium: { model: "gpt-5.6-terra", effort: "medium" },
  hard: { model: "gpt-5.6-sol", effort: "xhigh" },
} as const
const PROTECTED = /(^\.env(?:\..*)?$|credential|secret|token|\.pem$|\.key$|^id_rsa|^id_ed25519|\.pfx$|\.p12$|\.kdbx$)/i
const CREDENTIAL = /-----BEGIN [^-\r\n]*PRIVATE KEY-----|(?:^|[,{\s])(api[_-]?key|access[_-]?token|refresh[_-]?token|password|passwd|secret|authorization|cookie)\s*[=:]\s*["']?[^\s,"'}]+/im
const RESULT_SCHEMA = JSON.stringify({
  type: "object", additionalProperties: false, required: ["status", "summary", "files"],
  properties: { status: { type: "string", enum: ["SUCCESS", "FAILED", "ADVICE", "ESCALATE"] }, summary: { type: "string", maxLength: 4000 }, files: { type: "array", maxItems: MAX_FILES, items: { type: "object", additionalProperties: false, required: ["path", "content"], properties: { path: { type: "string", maxLength: 500 }, content: { type: "string", maxLength: MAX_FILE_BYTES } } } } },
})
let invocationActive = false
const digest = (value: string | Buffer) => createHash("sha256").update(value).digest("hex")
const inside = (root: string, target: string) => { const r = path.resolve(root).toLowerCase(); const t = path.resolve(target).toLowerCase(); return t === r || t.startsWith(r + path.sep) }
export function containsCredentialLikeContent(content: string) { return CREDENTIAL.test(content) }
export function selectRoute(difficulty: keyof typeof ROUTES) { return ROUTES[difficulty] ?? (() => { throw new Error("unsupported difficulty") })() }

async function exists(target: string) { try { await stat(target); return true } catch (e: any) { if (e?.code === "ENOENT") return false; throw e } }
async function removeWithRetry(target: string, recursive: boolean) {
  let lastError: unknown
  for (let attempt = 0; attempt < 8; attempt += 1) {
    try { await rm(target, { recursive, force: true }); return }
    catch (error) { lastError = error; await new Promise((resolve) => setTimeout(resolve, 250 * (attempt + 1))) }
  }
  throw lastError
}
function protectedPath(relative: string) { return relative.split(/[\\/]+/).some((s) => PROTECTED.test(s)) }
async function assertNoReparseTraversal(workspace: string, target: string) {
  const root = await realpath(workspace); const rel = path.relative(workspace, target); let cursor = workspace
  for (const segment of rel ? rel.split(path.sep) : []) { cursor = path.join(cursor, segment); const info = await lstat(cursor); if (info.isSymbolicLink()) throw new Error(`reparse traversal rejected: ${rel}`) }
  let existing = target; while (!(await exists(existing))) { const parent = path.dirname(existing); if (parent === existing) throw new Error(`missing path: ${rel}`); existing = parent }
  if (!inside(root, await realpath(existing))) throw new Error(`path escapes workspace: ${rel}`)
}
export async function normalizeFileList(workspace: string, values: string[]) {
  if (!Array.isArray(values) || values.length > MAX_FILES) throw new Error("invalid file list size")
  const result: Array<{ relative: string; absolute: string; key: string }> = []; const keys = new Set<string>()
  for (const value of values) {
    if (typeof value !== "string" || !value || value.length > 500 || path.isAbsolute(value)) throw new Error("file paths must be bounded relative paths")
    const relative = path.normalize(value).replace(/^\.([\\/])/, ""); const absolute = path.resolve(workspace, relative)
    if (!inside(workspace, absolute) || relative === ".." || relative.startsWith(`..${path.sep}`)) throw new Error(`path escapes workspace: ${value}`)
    if (protectedPath(relative)) throw new Error(`protected path rejected: ${relative}`)
    const key = relative.replaceAll("\\", "/").toLowerCase(); if (keys.has(key)) throw new Error(`case-insensitive duplicate path: ${relative}`); keys.add(key)
    if (!(await exists(absolute))) throw new Error(`context file missing: ${relative}`); await assertNoReparseTraversal(workspace, absolute)
    if (!(await lstat(absolute)).isFile()) throw new Error(`non-file path rejected: ${relative}`)
    result.push({ relative, absolute, key })
  }
  return result
}
async function readSafeText(file: { relative: string; absolute: string }) {
  const data = await readFile(file.absolute); if (data.length > MAX_FILE_BYTES || data.includes(0)) throw new Error(`binary or oversized file rejected: ${file.relative}`)
  let text: string; try { text = new TextDecoder("utf-8", { fatal: true }).decode(data) } catch { throw new Error(`non-UTF-8 file rejected: ${file.relative}`) }
  if (containsCredentialLikeContent(text)) throw new Error(`credential-like content rejected: ${file.relative}`); return text
}
async function resolveCodexExecutable() {
  const local = process.env.LOCALAPPDATA; if (!local) throw new Error("Codex executable root unavailable")
  const root = path.join(local, "OpenAI", "Codex", "bin"); const entries = await readdir(root, { withFileTypes: true }); const candidates: Array<{ path: string; mtime: number }> = []
  for (const entry of entries.slice(0, 64)) { if (!entry.isDirectory()) continue; const candidate = path.join(root, entry.name, "codex.exe"); if (!(await exists(candidate))) continue; const resolved = await realpath(candidate); if (!inside(root, resolved)) throw new Error("Codex executable escapes installation root"); const info = await stat(resolved); candidates.push({ path: resolved, mtime: info.mtimeMs }) }
  if (!candidates.length) throw new Error("installed codex.exe not found"); candidates.sort((a, b) => b.mtime - a.mtime || a.path.localeCompare(b.path)); return candidates[0].path
}
async function terminate(child: ReturnType<typeof spawn>) { if (child.exitCode !== null || !child.pid) return; await new Promise<void>((resolve) => { const k = spawn("taskkill.exe", ["/PID", String(child.pid), "/T", "/F"], { windowsHide: true, shell: false, stdio: "ignore" }); k.once("close", () => resolve()); k.once("error", () => resolve()); setTimeout(resolve, 10_000).unref() }) }
export async function runBoundedProcess(command: string, args: string[], options: { cwd: string; timeoutMs: number; signal?: AbortSignal }) {
  return await new Promise<{ code: number; stdout: string }>((resolve, reject) => { const child = spawn(command, args, { cwd: options.cwd, windowsHide: true, shell: false }); child.stdin?.end(); const chunks: Buffer[] = []; let bytes = 0; let done = false; let reason = ""
    const fail = (message: string) => { if (!done) { done = true; reject(new Error(message)) } }; const stop = async (message: string) => { if (reason) return; reason = message; await terminate(child); fail(message) }; const timer = setTimeout(() => void stop("codex process timeout"), options.timeoutMs); const abort = () => void stop("codex process aborted"); options.signal?.addEventListener("abort", abort, { once: true })
    child.stdout?.on("data", (c: Buffer) => { bytes += c.length; if (bytes > MAX_OUTPUT_BYTES) void stop("codex output limit exceeded"); else chunks.push(c) }); child.stderr?.on("data", (c: Buffer) => { bytes += c.length; if (bytes > MAX_OUTPUT_BYTES) void stop("codex output limit exceeded") }); child.once("error", () => fail("codex process could not start")); child.once("close", (code) => { clearTimeout(timer); options.signal?.removeEventListener("abort", abort); if (done) return; done = true; resolve({ code: code ?? -1, stdout: Buffer.concat(chunks).toString("utf8") }) })
  })
}
function parseResult(stdout: string) { let envelope: any; try { envelope = JSON.parse(stdout) } catch { throw new Error("codex response is not valid JSON") }; const value = envelope?.structured_output ?? envelope?.response ?? envelope; if (typeof value !== "object" || !value) throw new Error("codex response is not structured JSON"); if (!["SUCCESS", "FAILED", "ADVICE", "ESCALATE"].includes(value.status) || typeof value.summary !== "string" || !Array.isArray(value.files)) throw new Error("codex response schema invalid"); return value }
export async function executeCodexWorker(args: any, context: any, dependencies: any = {}) {
  if (invocationActive) return "FAILED: another codex_worker invocation is active"; invocationActive = true; let stage = ""; let lock: any; let lockPath = ""; let response = "FAILED: codex_worker failed"
  try {
    if (typeof args.task !== "string" || !args.task || args.task.length > 20_000 || containsCredentialLikeContent(args.task)) throw new Error("task must be bounded and non-sensitive")
    const workspace = await realpath(context.directory); const files = await normalizeFileList(workspace, args.context_files ?? []); lockPath = path.join(os.tmpdir(), `opencode-codex-worker-${digest(workspace.toLowerCase())}.lock`); lock = await open(lockPath, constants.O_CREAT | constants.O_EXCL | constants.O_WRONLY); await lock.writeFile(String(process.pid)); stage = await mkdtemp(path.join(os.tmpdir(), "opencode-codex-worker-stage-")); let total = 0
    for (const file of files) { const content = await readSafeText(file); total += Buffer.byteLength(content); if (total > MAX_TOTAL_BYTES) throw new Error("staged input size limit exceeded"); const destination = path.join(stage, file.relative); await mkdir(path.dirname(destination), { recursive: true }); await writeFile(destination, content, "utf8") }
    const route = selectRoute(args.difficulty); const exe = dependencies.resolveExecutable ? await dependencies.resolveExecutable() : await resolveCodexExecutable(); const run = dependencies.runProcess ?? runBoundedProcess; const schemaPath = path.join(stage, "result-schema.json"); const lastMessagePath = path.join(stage, "last-message.json"); await writeFile(schemaPath, RESULT_SCHEMA, "utf8"); const prompt = `Provide read-only advisory analysis for this task: ${args.task}\nOnly inspect staged context files: ${files.map((f: any) => f.relative).join(", ") || "none"}. Do not modify files or run external actions. Return JSON matching the supplied schema.`
    const result = await run(exe, ["exec", "--model", route.model, "-c", `model_reasoning_effort=\"${route.effort}\"`, "--sandbox", "read-only", "--ephemeral", "--ignore-user-config", "--skip-git-repo-check", "--output-schema", schemaPath, "--output-last-message", lastMessagePath, "--", prompt], { cwd: stage, timeoutMs: EXEC_TIMEOUT_MS, signal: context.abort }); if (result.code !== 0) throw new Error("codex generation failed"); const output = await readFile(lastMessagePath, "utf8").catch(() => result.stdout); const parsed = parseResult(output); const summary = String(parsed.summary).slice(0, 4000); if (containsCredentialLikeContent(JSON.stringify(parsed))) throw new Error("codex returned protected content"); response = `${parsed.status}: ${summary}`
  } catch (error: any) { response = `FAILED: ${String(error?.message ?? "codex_worker failed").slice(0, 4000)}`
  } finally {
    let cleanupFailure = ""
    if (stage) { try { await removeWithRetry(stage, true) } catch { cleanupFailure = `staging cleanup failed: ${path.basename(stage)}` } }
    if (lock) {
      try { await lock.close(); await removeWithRetry(lockPath, false) }
      catch { cleanupFailure ||= `lock cleanup failed: ${path.basename(lockPath)}` }
    }
    invocationActive = false
    if (cleanupFailure) response = `FAILED: ${cleanupFailure}`
  }
  return response
}
export default tool({ description: "Ask-gated, read-only Codex advisory worker using isolated staged context.", args: { task: tool.schema.string().min(1).max(20_000), difficulty: tool.schema.enum(["easy", "medium", "hard"]), context_files: tool.schema.array(tool.schema.string().min(1).max(500)).max(MAX_FILES).default([]) }, async execute(args, context) { return await executeCodexWorker(args, context) } })
