import { tool } from "@opencode-ai/plugin"
import { spawn } from "node:child_process"
import { createHash } from "node:crypto"
import {
  constants,
  lstat,
  mkdir,
  mkdtemp,
  open,
  readFile,
  realpath,
  rename,
  rm,
  stat,
  unlink,
  writeFile,
} from "node:fs/promises"
import os from "node:os"
import path from "node:path"

const MAX_FILES = 40
const MAX_FILE_BYTES = 512 * 1024
const MAX_TOTAL_BYTES = 2 * 1024 * 1024
const MAX_OUTPUT_BYTES = 3 * 1024 * 1024
const GENERATION_TIMEOUT_MS = 310_000
const MODELS_TIMEOUT_MS = 30_000

const ROUTES = {
  easy: { model: "gemini-3.7-flash-low", effort: "low" },
  medium: { model: "gemini-3.1-pro-low", effort: "medium" },
  hard: { model: "gemini-3.1-pro-high", effort: "high" },
} as const

const PROTECTED_SEGMENT = /(^\.env(?:\..*)?$|credential|secret|token|\.pem$|\.key$|^id_rsa|^id_ed25519|\.pfx$|\.p12$|\.kdbx$)/i
const CREDENTIAL_CONTENT = /-----BEGIN [^-\r\n]*PRIVATE KEY-----|(?:^|[,{\s])(api[_-]?key|access[_-]?token|refresh[_-]?token|password|passwd|secret|authorization|cookie)\s*[=:]\s*["']?[^\s,"'}]+/im

const RESULT_SCHEMA = JSON.stringify({
  type: "object",
  additionalProperties: false,
  required: ["status", "summary", "files"],
  properties: {
    status: { type: "string", enum: ["SUCCESS", "FAILED", "ADVICE", "ESCALATE"] },
    summary: { type: "string", maxLength: 4000 },
    files: {
      type: "array",
      maxItems: MAX_FILES,
      items: {
        type: "object",
        additionalProperties: false,
        required: ["path", "content"],
        properties: {
          path: { type: "string", maxLength: 500 },
          content: { type: "string", maxLength: MAX_FILE_BYTES },
        },
      },
    },
  },
})

let invocationActive = false

function digest(value: string | Buffer) {
  return createHash("sha256").update(value).digest("hex")
}

function inside(root: string, target: string) {
  const base = path.resolve(root).toLowerCase()
  const candidate = path.resolve(target).toLowerCase()
  return candidate === base || candidate.startsWith(base + path.sep)
}

function protectedPath(relative: string) {
  return relative.split(/[\\/]+/).some((segment) => PROTECTED_SEGMENT.test(segment))
}

export function containsCredentialLikeContent(content: string) {
  return CREDENTIAL_CONTENT.test(content)
}

export function selectRoute(difficulty: keyof typeof ROUTES, availableModels: Iterable<string>) {
  const route = ROUTES[difficulty]
  if (!route) throw new Error("unsupported difficulty")
  if (!new Set(availableModels).has(route.model)) throw new Error(`required model unavailable for ${difficulty}`)
  return route
}

async function exists(target: string) {
  try {
    await stat(target)
    return true
  } catch (error: any) {
    if (error?.code === "ENOENT") return false
    throw error
  }
}

async function assertNoReparseTraversal(workspace: string, target: string, targetMayBeMissing: boolean) {
  const workspaceReal = await realpath(workspace)
  const relative = path.relative(workspace, target)
  const segments = relative === "" ? [] : relative.split(path.sep)
  let cursor = workspace

  for (let index = 0; index < segments.length; index += 1) {
    cursor = path.join(cursor, segments[index])
    let info
    try {
      info = await lstat(cursor)
    } catch (error: any) {
      if (error?.code === "ENOENT" && targetMayBeMissing && index === segments.length - 1) break
      throw error
    }
    if (info.isSymbolicLink()) throw new Error(`reparse traversal rejected: ${relative}`)
  }

  let existing = target
  while (!(await exists(existing))) {
    const parent = path.dirname(existing)
    if (parent === existing) throw new Error(`no existing parent: ${relative}`)
    existing = parent
  }
  const existingReal = await realpath(existing)
  if (!inside(workspaceReal, existingReal)) throw new Error(`path escapes workspace: ${relative}`)
}

export async function normalizeFileList(workspace: string, values: string[], options: { allowMissing: boolean }) {
  if (!Array.isArray(values) || values.length > MAX_FILES) throw new Error("invalid file list size")
  const normalized = [] as Array<{ relative: string; absolute: string; key: string; exists: boolean }>
  const keys = new Set<string>()

  for (const value of values) {
    if (typeof value !== "string" || value.length === 0 || value.length > 500 || path.isAbsolute(value)) {
      throw new Error("file paths must be bounded relative paths")
    }
    const relative = path.normalize(value).replace(/^\.[\\/]/, "")
    const absolute = path.resolve(workspace, relative)
    if (!inside(workspace, absolute) || relative === ".." || relative.startsWith(`..${path.sep}`)) {
      throw new Error(`path escapes workspace: ${value}`)
    }
    if (protectedPath(relative)) throw new Error(`protected path rejected: ${relative}`)
    const key = relative.replaceAll("\\", "/").toLowerCase()
    if (keys.has(key)) throw new Error(`case-insensitive duplicate path: ${relative}`)
    keys.add(key)

    const present = await exists(absolute)
    if (!present && !options.allowMissing) throw new Error(`context file missing: ${relative}`)
    await assertNoReparseTraversal(workspace, absolute, options.allowMissing)
    if (present && !(await lstat(absolute)).isFile()) throw new Error(`non-file path rejected: ${relative}`)
    normalized.push({ relative, absolute, key, exists: present })
  }
  return normalized
}

async function readSafeText(file: { relative: string; absolute: string }) {
  const data = await readFile(file.absolute)
  if (data.length > MAX_FILE_BYTES) throw new Error(`file too large: ${file.relative}`)
  if (data.includes(0)) throw new Error(`binary file rejected: ${file.relative}`)
  let content: string
  try {
    content = new TextDecoder("utf-8", { fatal: true }).decode(data)
  } catch {
    throw new Error(`non-UTF-8 file rejected: ${file.relative}`)
  }
  if (containsCredentialLikeContent(content)) throw new Error(`credential-like content rejected: ${file.relative}`)
  return content
}

async function terminateTree(child: ReturnType<typeof spawn>) {
  if (child.exitCode !== null || !child.pid) return
  if (process.platform === "win32") {
    await new Promise<void>((resolve) => {
      const killer = spawn("taskkill.exe", ["/PID", String(child.pid), "/T", "/F"], {
        windowsHide: true,
        shell: false,
        stdio: "ignore",
      })
      killer.once("exit", () => resolve())
      killer.once("error", () => resolve())
      setTimeout(resolve, 10_000).unref()
    })
  } else {
    child.kill("SIGKILL")
  }
}

type ProcessResult = { code: number; stdout: string }

export async function runBoundedProcess(
  command: string,
  args: string[],
  options: { cwd: string; timeoutMs: number; signal?: AbortSignal },
): Promise<ProcessResult> {
  return await new Promise<ProcessResult>((resolve, reject) => {
    const child = spawn(command, args, { cwd: options.cwd, windowsHide: true, shell: false })
    child.stdin?.end()
    const stdout: Buffer[] = []
    let outputBytes = 0
    let settled = false
    let stopReason = ""

    const finishError = (message: string) => {
      if (settled) return
      settled = true
      reject(new Error(message))
    }
    const stop = async (reason: string) => {
      if (stopReason) return
      stopReason = reason
      await terminateTree(child)
      setTimeout(() => finishError(reason), 10_000).unref()
    }
    const timer = setTimeout(() => void stop("agy process timeout"), options.timeoutMs)
    const abort = () => void stop("agy process aborted")
    options.signal?.addEventListener("abort", abort, { once: true })
    if (options.signal?.aborted) void stop("agy process aborted")

    child.stdout?.on("data", (chunk: Buffer) => {
      outputBytes += chunk.length
      if (outputBytes > MAX_OUTPUT_BYTES) {
        void stop("agy output limit exceeded")
        return
      }
      stdout.push(chunk)
    })
    child.stderr?.on("data", (chunk: Buffer) => {
      outputBytes += chunk.length
      if (outputBytes > MAX_OUTPUT_BYTES) void stop("agy output limit exceeded")
    })
    child.once("error", () => finishError("agy process could not start"))
    child.once("close", (code) => {
      clearTimeout(timer)
      options.signal?.removeEventListener("abort", abort)
      if (settled) return
      settled = true
      if (stopReason) return reject(new Error(stopReason))
      resolve({ code: code ?? -1, stdout: Buffer.concat(stdout).toString("utf8") })
    })
  })
}

async function removeWithRetry(target: string, options: { recursive: boolean }) {
  let lastError: unknown
  for (let attempt = 0; attempt < 8; attempt += 1) {
    try {
      await rm(target, { recursive: options.recursive, force: true })
      return
    } catch (error) {
      lastError = error
      await new Promise((resolve) => setTimeout(resolve, 250 * (attempt + 1)))
    }
  }
  throw lastError
}

function parseModels(stdout: string) {
  return stdout
    .split(/\r?\n/)
    .map((line) => line.split("\t", 1)[0].trim())
    .filter((value) => /^[a-z0-9][a-z0-9.-]+$/i.test(value))
}

function extractStructuredResponse(stdout: string) {
  let envelope: any
  try {
    envelope = JSON.parse(stdout)
  } catch {
    throw new Error("agy envelope is not valid JSON")
  }
  if (envelope?.status !== "SUCCESS") throw new Error("agy returned a non-success envelope")
  let response = envelope.structured_output ?? envelope.response
  if (typeof response === "string") {
    try {
      response = JSON.parse(response)
    } catch {
      throw new Error("agy structured response is not valid JSON")
    }
  }
  if (!response || typeof response !== "object") throw new Error("agy response is not structured JSON")
  return response
}

function buildPrompt(args: any, allowed: Array<{ relative: string }>, context: Array<{ relative: string }>) {
  return [
    "Role: bounded external implementation worker",
    `Task: ${args.task}`,
    `Allowed output files: ${allowed.map((file) => file.relative).join(", ")}`,
    `Read-only context files: ${context.map((file) => file.relative).join(", ") || "none"}`,
    "Work only in this isolated staging directory.",
    "Return structured JSON matching the supplied schema.",
    "For every file to change or create, return its relative path and complete UTF-8 text content.",
    "Do not return or request files outside Allowed output files. Do not delete files.",
    "Do not include credentials, tokens, secrets, private keys, cookies, or authorization values.",
    "Do not run deployment, VCS write, publication, or external messaging operations.",
    "If the task cannot be completed from the staged files, return ESCALATE with no file outputs.",
  ].join("\n")
}

async function hashState(target: { absolute: string; exists: boolean }) {
  const present = await exists(target.absolute)
  return { exists: present, hash: present ? digest(await readFile(target.absolute)) : null }
}

async function replaceWithContent(target: string, content: string) {
  const temporary = `${target}.agy-worker-${process.pid}-${Date.now()}.tmp`
  await writeFile(temporary, content, { encoding: "utf8", flag: "wx" })
  try {
    await rename(temporary, target)
  } catch (error) {
    await rm(temporary, { force: true })
    throw error
  }
}

export async function applyTransaction(
  outputs: Array<{ path: string; content: string }>,
  allowed: Array<{ relative: string; absolute: string; key: string; exists: boolean }>,
  baseline: Map<string, { exists: boolean; hash: string | null; content: string | null }>,
  options: { dryRun?: boolean } = {},
) {
  const allowedByKey = new Map(allowed.map((file) => [file.key, file]))
  const seen = new Set<string>()
  let total = 0
  const prepared = [] as Array<{ file: (typeof allowed)[number]; content: string }>

  for (const output of outputs) {
    if (!output || typeof output.path !== "string" || typeof output.content !== "string") throw new Error("invalid agy file output")
    const key = path.normalize(output.path).replaceAll("\\", "/").toLowerCase()
    const file = allowedByKey.get(key)
    if (!file || seen.has(key)) throw new Error(`unallowed or duplicate agy output: ${output.path}`)
    if (containsCredentialLikeContent(output.content)) throw new Error(`credential-like generated content rejected: ${file.relative}`)
    total += Buffer.byteLength(output.content, "utf8")
    if (total > MAX_TOTAL_BYTES || Buffer.byteLength(output.content, "utf8") > MAX_FILE_BYTES) throw new Error("generated content size limit exceeded")
    seen.add(key)
    prepared.push({ file, content: output.content })
  }

  for (const item of prepared) {
    const before = baseline.get(item.file.key)!
    const current = await hashState(item.file)
    if (current.exists !== before.exists || current.hash !== before.hash) throw new Error(`concurrent change detected: ${item.file.relative}`)
  }
  if (options.dryRun) return prepared.map((item) => item.file.relative)

  const committed = [] as Array<{ file: (typeof allowed)[number]; postHash: string }>
  try {
    for (const item of prepared) {
      const before = baseline.get(item.file.key)!
      const current = await hashState(item.file)
      if (current.exists !== before.exists || current.hash !== before.hash) throw new Error(`mid-commit collision: ${item.file.relative}`)
      await replaceWithContent(item.file.absolute, item.content)
      committed.push({ file: item.file, postHash: digest(item.content) })
    }
  } catch (error: any) {
    const rollbackFailed: string[] = []
    for (const item of committed.reverse()) {
      const before = baseline.get(item.file.key)!
      const current = await hashState(item.file)
      if (!current.exists || current.hash !== item.postHash) {
        rollbackFailed.push(item.file.relative)
        continue
      }
      try {
        if (before.exists) await replaceWithContent(item.file.absolute, before.content!)
        else await unlink(item.file.absolute)
      } catch {
        rollbackFailed.push(item.file.relative)
      }
    }
    if (rollbackFailed.length) throw new Error(`PARTIAL_ROLLBACK_FAILED: ${rollbackFailed.join(", ")}`)
    throw error
  }
  return prepared.map((item) => item.file.relative)
}

export async function executeAgyWorker(args: any, context: any, dependencies = { runProcess: runBoundedProcess }) {
  if (invocationActive) return "FAILED: another agy_worker invocation is active"
  invocationActive = true
  let stage = ""
  let lockHandle: Awaited<ReturnType<typeof open>> | undefined
  let lockPath = ""

  try {
    if (typeof args.task !== "string" || args.task.length === 0 || args.task.length > 20_000) throw new Error("task must be bounded text")
    if (containsCredentialLikeContent(args.task)) throw new Error("credential-like task content rejected")
    const workspace = await realpath(context.directory)
    const allowed = await normalizeFileList(workspace, args.allowed_files, { allowMissing: true })
    if (allowed.length === 0) throw new Error("at least one allowed file is required")
    const contextFiles = await normalizeFileList(workspace, args.context_files ?? [], { allowMissing: false })

    lockPath = path.join(os.tmpdir(), `opencode-agy-worker-${digest(workspace.toLowerCase())}.lock`)
    try {
      lockHandle = await open(lockPath, constants.O_CREAT | constants.O_EXCL | constants.O_WRONLY)
      await lockHandle.writeFile(String(process.pid))
    } catch (error: any) {
      if (error?.code === "EEXIST") throw new Error("workspace agy_worker lock already exists")
      throw error
    }

    const baseline = new Map<string, { exists: boolean; hash: string | null; content: string | null }>()
    const inputContent = new Map<string, string>()
    let stagedBytes = 0
    stage = await mkdtemp(path.join(os.tmpdir(), "opencode-agy-worker-stage-"))
    const stagedByKey = new Set<string>()
    for (const file of [...allowed, ...contextFiles]) {
      if (!file.exists || stagedByKey.has(file.key)) continue
      const content = await readSafeText(file)
      stagedBytes += Buffer.byteLength(content, "utf8")
      if (stagedBytes > MAX_TOTAL_BYTES) throw new Error("staged input size limit exceeded")
      const destination = path.join(stage, file.relative)
      await mkdir(path.dirname(destination), { recursive: true })
      await writeFile(destination, content, "utf8")
      inputContent.set(file.key, content)
      stagedByKey.add(file.key)
    }
    for (const file of allowed) {
      const content = file.exists ? inputContent.get(file.key)! : null
      baseline.set(file.key, { exists: file.exists, hash: content === null ? null : digest(content), content })
    }

    const modelsResult = await dependencies.runProcess("agy.exe", ["models"], {
      cwd: stage,
      timeoutMs: MODELS_TIMEOUT_MS,
      signal: context.abort,
    })
    if (modelsResult.code !== 0) throw new Error("agy model discovery failed")
    const route = selectRoute(args.difficulty, parseModels(modelsResult.stdout))
    const prompt = buildPrompt(args, allowed, contextFiles)
    const generationResult = await dependencies.runProcess(
      "agy.exe",
      [
        "--model", route.model,
        "--effort", route.effort,
        "--mode", "plan",
        "--output-format", "json",
        "--print-timeout", "5m",
        "--sandbox",
        "--new-project",
        "--add-dir", stage,
        `--json-schema=${RESULT_SCHEMA}`,
        `--prompt=${prompt}`,
      ],
      { cwd: stage, timeoutMs: GENERATION_TIMEOUT_MS, signal: context.abort },
    )
    if (generationResult.code !== 0) throw new Error("agy generation failed")
    const result = extractStructuredResponse(generationResult.stdout)
    if (result.status !== "SUCCESS") {
      const summary = String(result.summary ?? "agy did not produce changes").slice(0, 4000)
      if (containsCredentialLikeContent(summary)) return "FAILED: agy returned protected content"
      return `${result.status}: ${summary}`
    }
    if (!Array.isArray(result.files)) throw new Error("agy files result is invalid")

    if (args.dry_run) {
      await applyTransaction(result.files, allowed, baseline, { dryRun: true })
      return `SUCCESS: dry-run validated ${result.files.length} file output(s); workspace unchanged`
    }

    const changed = await applyTransaction(result.files, allowed, baseline)
    return `SUCCESS: agy ${route.model}/${route.effort} applied ${changed.length} file(s): ${changed.join(", ")}`
  } catch (error: any) {
    return `FAILED: ${String(error?.message ?? "agy_worker failed").slice(0, 4000)}`
  } finally {
    let cleanupFailure = ""
    if (stage) {
      try {
        await removeWithRetry(stage, { recursive: true })
      } catch {
        cleanupFailure = `staging cleanup failed: ${path.basename(stage)}`
      }
    }
    if (lockHandle) {
      try {
        await lockHandle.close()
        await removeWithRetry(lockPath, { recursive: false })
      } catch {
        cleanupFailure ||= `lock cleanup failed: ${path.basename(lockPath)}`
      }
    }
    invocationActive = false
    if (cleanupFailure) return `FAILED: ${cleanupFailure}`
  }
}

export default tool({
  description: "Delegate one bounded text-file implementation task to agy.exe through isolated staging and an allowlisted transactional apply.",
  args: {
    task: tool.schema.string().min(1).max(20_000).describe("Exact implementation task and acceptance checks; never include secrets."),
    difficulty: tool.schema.enum(["easy", "medium", "hard"]).describe("Routing difficulty."),
    allowed_files: tool.schema.array(tool.schema.string().min(1).max(500)).min(1).max(MAX_FILES).describe("Relative text files agy may create or replace."),
    context_files: tool.schema.array(tool.schema.string().min(1).max(500)).max(MAX_FILES).default([]).describe("Relative read-only text files copied into staging."),
    dry_run: tool.schema.boolean().default(false).describe("Validate a real agy response without writing the workspace."),
  },
  async execute(args, context) {
    await context.ask({ permission: "agy_worker_current_task", patterns: [args.task], always: [], metadata: { allowed_files: args.allowed_files, opt_in: "Explicit current-task agy implementation approval required" } })
    return await executeAgyWorker(args, context)
  },
})
