import { tool } from "@opencode-ai/plugin"
import { spawn } from "node:child_process"
import { mkdtemp, rm } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join } from "node:path"

export const AGGRESSIVE_ROSTER = [
  { model: "gemini-3.7-flash-high", effort: "high" },
  { model: "gemini-3.1-pro-high", effort: "high" },
  { model: "claude-sonnet-4-6", effort: "high" },
  { model: "claude-opus-4-6-thinking", effort: "high" },
  { model: "gpt-oss-120b-medium", effort: "medium" },
] as const

const MAX_CONCURRENCY = 4
const DISCOVERY_TIMEOUT_MS = 30_000
const GLOBAL_DEADLINE_MS = 360_000
const MAX_INPUT_BYTES = 12_000
const MAX_PROCESS_OUTPUT_BYTES = 1_000_000
const SECRET = /-----BEGIN [^-\r\n]*PRIVATE KEY-----|(?:^|[,{\s])(api[_-]?key|access[_-]?token|refresh[_-]?token|password|passwd|secret|authorization|cookie)\s*[=:]\s*["']?[^\s,"'}]+/im
const WINDOWS_ABSOLUTE = /(?:^|[\s("'`])(?:[a-z]:[\\/]|\\\\)/i
const POSIX_ABSOLUTE = /(?:^|[\s("'`])\/(?!\/)[^\s]*/
const PARENT_TRAVERSAL = /(?:^|[\s\\/])\.\.(?:[\\/]|$)/
const FILE_URL = /file:\/\//i

const PROVIDER_SCHEMA = {
  type: "object",
  additionalProperties: false,
  required: ["status", "summary", "findings"],
  properties: {
    status: { type: "string", enum: ["SUCCESS", "FAILED", "ADVICE"] },
    summary: { type: "string", maxLength: 4000 },
    findings: { type: "array", maxItems: 12, items: { type: "string", maxLength: 1000 } },
  },
} as const
const PROVIDER_SCHEMA_JSON = JSON.stringify(PROVIDER_SCHEMA)

type ProcessResult = { code: number; stdout: string }
type Runner = (command: string, args: string[], options: { cwd: string; timeoutMs: number; signal?: AbortSignal }) => Promise<ProcessResult>
type Dependencies = {
  runProcess: Runner
  makeTemp: () => Promise<string>
  removeTemp: (path: string) => Promise<void>
  now: () => number
}
type ProviderStatus = "SUCCESS" | "FAILED" | "ADVICE"
type AggregateStatus = ProviderStatus | "UNAVAILABLE"
type CouncilResult = { model: string; status: AggregateStatus; summary: string; findings: string[] }

async function terminateTree(child: ReturnType<typeof spawn>) {
  if (child.exitCode !== null || !child.pid) return
  if (process.platform !== "win32") {
    child.kill("SIGKILL")
    return
  }
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
}

async function runBoundedProcess(
  command: string,
  args: string[],
  options: { cwd: string; timeoutMs: number; signal?: AbortSignal },
): Promise<ProcessResult> {
  return await new Promise<ProcessResult>((resolve, reject) => {
    const child = spawn(command, args, { cwd: options.cwd, shell: false, windowsHide: true, stdio: ["ignore", "pipe", "pipe"] })
    const stdout: Buffer[] = []
    let outputBytes = 0
    let settled = false
    let stopping = false
    const finish = (error?: Error, result?: ProcessResult) => {
      if (settled) return
      settled = true
      clearTimeout(timer)
      options.signal?.removeEventListener("abort", abort)
      if (error) reject(error)
      else resolve(result!)
    }
    const stop = async (reason: string) => {
      if (stopping || settled) return
      stopping = true
      await terminateTree(child)
      finish(new Error(reason))
    }
    const abort = () => void stop("agy process aborted")
    const timer = setTimeout(() => void stop("agy process timeout"), Math.max(1, options.timeoutMs))
    options.signal?.addEventListener("abort", abort, { once: true })
    if (options.signal?.aborted) void stop("agy process aborted")
    const count = (chunk: Buffer) => {
      outputBytes += chunk.length
      if (outputBytes > MAX_PROCESS_OUTPUT_BYTES) void stop("agy output limit exceeded")
    }
    child.stdout?.on("data", (chunk: Buffer) => {
      count(chunk)
      if (!stopping) stdout.push(chunk)
    })
    child.stderr?.on("data", count)
    child.once("error", () => finish(new Error("agy process could not start")))
    child.once("close", (code) => {
      if (!stopping) finish(undefined, { code: code ?? -1, stdout: Buffer.concat(stdout).toString("utf8") })
    })
  })
}

const defaultDependencies: Dependencies = {
  runProcess: runBoundedProcess,
  makeTemp: () => mkdtemp(join(tmpdir(), "agy-council-")),
  removeTemp: (path) => rm(path, { recursive: true, force: true }),
  now: () => Date.now(),
}

function unsafeUrl(value: string) {
  for (const match of value.matchAll(/https?:\/\/[^\s)\]}>]+/gi)) {
    try {
      const url = new URL(match[0])
      if (url.username || url.password || url.search || url.hash) return true
    } catch {
      return true
    }
  }
  return false
}

function protectedText(value: string) {
  return SECRET.test(value) || WINDOWS_ABSOLUTE.test(value) || POSIX_ABSOLUTE.test(value) || PARENT_TRAVERSAL.test(value) || FILE_URL.test(value) || unsafeUrl(value)
}

export function parseModels(stdout: string): string[] {
  const roster = new Set(AGGRESSIVE_ROSTER.map((entry) => entry.model))
  return stdout.split(/\r?\n/).map((line) => line.split("\t", 1)[0].trim()).filter((model) => roster.has(model as (typeof AGGRESSIVE_ROSTER)[number]["model"]))
}

function parseProviderResponse(stdout: string): { status: ProviderStatus; summary: string; findings: string[] } {
  let envelope: any
  try { envelope = JSON.parse(stdout) } catch { throw new Error("invalid agy envelope") }
  if (!envelope || typeof envelope !== "object" || Array.isArray(envelope) || envelope.status !== "SUCCESS") throw new Error("agy envelope failed")
  const carriers = ["structured_output", "response"].filter((key) => Object.hasOwn(envelope, key))
  if (carriers.length !== 1 || Object.keys(envelope).sort().join(",") !== ["status", carriers[0]].sort().join(",")) throw new Error("unexpected agy envelope fields")
  let value = envelope[carriers[0]]
  if (typeof value === "string") {
    try { value = JSON.parse(value) } catch { throw new Error("invalid agy payload") }
  }
  if (!value || typeof value !== "object" || Array.isArray(value)) throw new Error("invalid agy payload")
  if (Object.keys(value).sort().join(",") !== "findings,status,summary") throw new Error("unexpected agy payload fields")
  if (!["SUCCESS", "FAILED", "ADVICE"].includes(value.status)) throw new Error("invalid agy status")
  if (typeof value.summary !== "string" || value.summary.length > 4000 || protectedText(value.summary)) throw new Error("protected agy summary")
  if (!Array.isArray(value.findings) || value.findings.length > 12 || value.findings.some((item: unknown) => typeof item !== "string" || item.length > 1000 || protectedText(item))) throw new Error("protected agy findings")
  return value as { status: ProviderStatus; summary: string; findings: string[] }
}

function fixedResult(model: string, status: AggregateStatus, summary: string): CouncilResult {
  return { model, status, summary, findings: [] }
}

let councilActive = false

async function executeUnlocked(args: any, context: any, dependencies: Dependencies): Promise<string> {
  if (typeof args.task !== "string" || args.task.length === 0 || Buffer.byteLength(args.task, "utf8") > MAX_INPUT_BYTES || protectedText(args.task)) {
    return JSON.stringify({ status: "FAILED", results: [], unavailable_models: [], error: "invalid or protected task" })
  }
  const roster = [...AGGRESSIVE_ROSTER]
  if (args.primary_model !== undefined) {
    if (!roster.some((entry) => entry.model === args.primary_model)) return JSON.stringify({ status: "FAILED", results: [], unavailable_models: [], error: "primary model is outside the council roster" })
    roster.sort((left, right) => Number(right.model === args.primary_model) - Number(left.model === args.primary_model))
  }
  const deadlineMs = Math.min(Math.max(Number(args.deadline_ms) || GLOBAL_DEADLINE_MS, 1_000), GLOBAL_DEADLINE_MS)
  const deadlineAt = dependencies.now() + deadlineMs
  const controller = new AbortController()
  const callerAbort = () => controller.abort()
  context?.abort?.addEventListener?.("abort", callerAbort, { once: true })
  const deadlineTimer = setTimeout(() => controller.abort(), deadlineMs)
  try {
    let discoveryTemp = ""
    let discovered: string[] = []
    try {
      discoveryTemp = await dependencies.makeTemp()
      const remaining = deadlineAt - dependencies.now()
      if (remaining <= 0) throw new Error("deadline")
      const discovery = await dependencies.runProcess("agy.exe", ["models"], { cwd: discoveryTemp, timeoutMs: Math.min(DISCOVERY_TIMEOUT_MS, remaining), signal: controller.signal })
      if (discovery.code !== 0) throw new Error("discovery")
      discovered = parseModels(discovery.stdout)
    } catch {
      return JSON.stringify({ status: "FAILED", results: roster.map((entry) => fixedResult(entry.model, "UNAVAILABLE", "model unavailable")), unavailable_models: roster.map((entry) => entry.model), error: "model discovery unavailable" })
    } finally {
      if (discoveryTemp) {
        try { await dependencies.removeTemp(discoveryTemp) } catch {
          return JSON.stringify({ status: "FAILED", results: roster.map((entry) => fixedResult(entry.model, "UNAVAILABLE", "model unavailable")), unavailable_models: roster.map((entry) => entry.model), error: "temporary cleanup failed" })
        }
      }
    }

    const live = roster.filter((entry) => discovered.includes(entry.model))
    const results: CouncilResult[] = roster.filter((entry) => !discovered.includes(entry.model)).map((entry) => fixedResult(entry.model, "UNAVAILABLE", "model unavailable"))
    let cursor = 0
    const consult = async () => {
      while (true) {
        const entry = live[cursor++]
        if (!entry) return
        if (controller.signal.aborted || dependencies.now() >= deadlineAt) {
          results.push(fixedResult(entry.model, "FAILED", "global deadline exceeded"))
          continue
        }
        let providerTemp = ""
        let result: CouncilResult
        try {
          providerTemp = await dependencies.makeTemp()
          const remaining = deadlineAt - dependencies.now()
          if (remaining <= 0) throw new Error("deadline")
          const prompt = [
            "Role: read-only council analyst",
            `Task: ${args.task}`,
            "Return only the required JSON object.",
            "Do not request, reveal, or infer credentials, local paths, workspace files, raw logs, or private environment values.",
            "Do not edit files, publish, deploy, send messages, or perform VCS operations.",
          ].join("\n")
          const argv = ["--model", entry.model, "--effort", entry.effort, "--mode", "plan", "--output-format", "json", "--print-timeout", "5m", "--sandbox", "--new-project", "--disable-slash-commands", `--json-schema=${PROVIDER_SCHEMA_JSON}`, `--prompt=${prompt}`]
          const response = await dependencies.runProcess("agy.exe", argv, { cwd: providerTemp, timeoutMs: Math.min(300_000, remaining), signal: controller.signal })
          if (response.code !== 0) throw new Error("provider")
          const parsed = parseProviderResponse(response.stdout)
          result = { model: entry.model, ...parsed }
        } catch {
          result = fixedResult(entry.model, "FAILED", "consultation failed")
        }
        if (providerTemp) {
          try { await dependencies.removeTemp(providerTemp) } catch { result = fixedResult(entry.model, "FAILED", "temporary cleanup failed") }
        }
        results.push(result)
      }
    }
    await Promise.allSettled(Array.from({ length: Math.min(MAX_CONCURRENCY, live.length) }, consult))
    for (const entry of roster) if (!results.some((result) => result.model === entry.model)) results.push(fixedResult(entry.model, "FAILED", "global deadline exceeded"))
    results.sort((left, right) => roster.findIndex((entry) => entry.model === left.model) - roster.findIndex((entry) => entry.model === right.model))
    const successCount = results.filter((result) => result.status === "SUCCESS").length
    const status = successCount === roster.length ? "SUCCESS" : successCount === 0 ? "FAILED" : "ADVICE"
    return JSON.stringify({ status, primary_model: args.primary_model ?? null, results, unavailable_models: results.filter((result) => result.status === "UNAVAILABLE").map((result) => result.model) })
  } finally {
    clearTimeout(deadlineTimer)
    context?.abort?.removeEventListener?.("abort", callerAbort)
    controller.abort()
  }
}

export async function executeAgyCouncil(args: any, context: any, overrides: Partial<Dependencies> = {}): Promise<string> {
  if (councilActive) return JSON.stringify({ status: "FAILED", results: [], unavailable_models: [], error: "another council invocation is active" })
  councilActive = true
  try { return await executeUnlocked(args, context, { ...defaultDependencies, ...overrides }) } finally { councilActive = false }
}

export default tool({
  description: "Ask up to five live external models for a bounded read-only consultation. One approval may consume quota for five model calls.",
  args: {
    task: tool.schema.string().min(1).max(MAX_INPUT_BYTES).describe("Bounded consultation text only; no secrets, absolute paths, file contents, or credential-bearing URLs."),
    primary_model: tool.schema.enum(AGGRESSIVE_ROSTER.map((entry) => entry.model) as [string, ...string[]]).optional().describe("Optional model to present first; every live roster model is still consulted."),
    deadline_ms: tool.schema.number().int().min(1_000).max(GLOBAL_DEADLINE_MS).optional().describe("Overall deadline including discovery and queue time."),
  },
  async execute(args, context) { return await executeAgyCouncil(args, context) },
})
