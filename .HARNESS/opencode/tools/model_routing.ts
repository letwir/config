export const WORKER_DIFFICULTIES = Object.freeze(["easy", "medium", "hard"] as const)
export type WorkerDifficulty = typeof WORKER_DIFFICULTIES[number]

export type WorkerProfile = Readonly<{
  agent: string
  providerID: string
  modelID: string
  reasoningEffort?: "low" | "medium" | "high"
}>

const WORKER_ROUTES: Readonly<Record<WorkerDifficulty, WorkerProfile>> = Object.freeze({
  easy: Object.freeze({ agent: "worker-easy", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "low" }),
  medium: Object.freeze({ agent: "worker-medium", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "medium" }),
  hard: Object.freeze({ agent: "worker-hard", providerID: "openai", modelID: "gpt-6.1-sol", reasoningEffort: "high" }),
})

export const ASTRA_WORKER_PROFILE: WorkerProfile = Object.freeze({
  agent: "worker-astra",
  providerID: "openai",
  modelID: "gpt-6-astra",
})

export const ASTRA_MODEL = Object.freeze({ providerID: "openai", modelID: "gpt-6-astra" })
export const HARD_SOL_MODEL = Object.freeze({ providerID: "openai", modelID: "gpt-6.1-sol" })

export type ReviewDepth = "low" | "medium" | "hard"
export type ReviewIdentity = Readonly<{ target: string; revision: string; lens: string }>
export const normalizeReviewRole = (role: string) => role === "blackhat" ? "attacker" : role

// Review depth is intentionally independent of CODE difficulty.
export function selectAttackerProfile(depth: unknown = "low", escalationOf?: unknown): WorkerProfile {
  if (!["low", "medium", "hard"].includes(depth as string)) throw new Error("review_depth must be low, medium, or hard")
  if (depth === "hard") {
    if (typeof escalationOf !== "string" || !escalationOf.trim()) throw new Error("hard review requires escalation_of a completed Sol attacker")
    return Object.freeze({ agent: "attacker-hard", ...ASTRA_MODEL })
  }
  if (escalationOf !== undefined) throw new Error("attacker escalation requires review_depth=hard")
  return depth === "medium"
    ? Object.freeze({ agent: "attacker-medium", ...HARD_SOL_MODEL, reasoningEffort: "high" })
    : Object.freeze({ agent: "attacker-low", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "low" })
}

export function reviewIdentity(target: unknown, revision: unknown, lens: unknown): ReviewIdentity {
  const values = [target, revision, lens]
  if (values.some((v) => typeof v !== "string" || !v.trim() || v.length > 1000 || /[\x00-\x1f\x7f]/.test(v))) throw new Error("explicit bounded review target/revision/lens required")
  return Object.freeze({ target: target as string, revision: revision as string, lens: lens as string })
}

// A prose ESCALATE, partial assessment, or operational failure never authorizes Astra.
// The parent must additionally inspect the whole report and explicitly approve.
export function hasReasoningInsufficiency(text: unknown): boolean {
  if (typeof text !== "string" || text.length > 64_000 || parseWorkerStatus(text) !== "ESCALATE") return false
  const lines = text.replaceAll("\r\n", "\n").split("\n")
  if (lines.filter((line) => line.startsWith("RTS_ESCALATION:")).length !== 1 || !lines[1]?.startsWith("RTS_ESCALATION: ")) return false
  try {
    const evidence = JSON.parse(lines[1].slice("RTS_ESCALATION: ".length))
    return evidence !== null && typeof evidence === "object" &&
      Object.keys(evidence).sort().join(",") === "blocked_by,bounded_assessment_completed,cause,further_reasoning_needed,unresolved_analysis" &&
      evidence.cause === "reasoning_insufficiency" && evidence.bounded_assessment_completed === true &&
      evidence.further_reasoning_needed === true && Array.isArray(evidence.blocked_by) && evidence.blocked_by.length === 0 &&
      typeof evidence.unresolved_analysis === "string" && evidence.unresolved_analysis.trim().length > 0
  } catch { return false } // Invalid evidence is a denial, not an operational success.
}

type AttackerSource = {
  parent: string; root: string; depth: number; role: string; profile: string; reviewDepth?: ReviewDepth; review?: ReviewIdentity
  actualModel?: { providerID: string; modelID: string }; routeStatus: string; stopped: boolean; resultLoaded?: boolean
  state: string; resultStatus: string; result?: string; error?: string; escalationReservedBy?: string
  scope: readonly string[]; writes: readonly string[]
}
export function assertAttackerEscalation(source: AttackerSource | undefined, request: {
  parent: string; root: string; depth: number; scope: readonly string[]; writes: readonly string[]; review?: ReviewIdentity; approved?: boolean
}): void {
  if (!source || source.escalationReservedBy || source.error || source.parent !== request.parent || source.root !== request.root ||
    ![1, 2].includes(source.depth) || source.depth !== request.depth ||
    source.role !== "attacker" || source.profile !== "attacker-medium" || source.reviewDepth !== "medium" ||
    !modelMatches(source.actualModel, HARD_SOL_MODEL) || source.routeStatus !== "matched" || !source.stopped ||
    !source.resultLoaded || source.state !== "completed" || source.resultStatus !== "ESCALATE" ||
    !hasReasoningInsufficiency(source.result) || request.approved !== true || source.writes.length || request.writes.length ||
    !samePathSet(source.scope, request.scope) || !source.review || !request.review ||
    source.review.target !== request.review.target || source.review.revision !== request.review.revision || source.review.lens !== request.review.lens) {
    throw new Error("source result is not eligible for one-shot Sol attacker reasoning escalation")
  }
}

export function isWorkerDifficulty(value: unknown): value is WorkerDifficulty {
  return typeof value === "string" && (WORKER_DIFFICULTIES as readonly string[]).includes(value)
}

export function selectWorkerProfile(difficulty: unknown, escalationOf?: unknown): WorkerProfile {
  if (!isWorkerDifficulty(difficulty)) throw new Error("worker difficulty must be easy, medium, or hard")
  if (escalationOf !== undefined) {
    if (difficulty !== "hard") throw new Error("Astra escalation requires hard difficulty")
    if (typeof escalationOf !== "string" || !escalationOf.trim()) throw new Error("escalation_of must identify a source task")
    return ASTRA_WORKER_PROFILE
  }
  return WORKER_ROUTES[difficulty]
}

export function parseWorkerStatus(text: unknown): "SUCCESS" | "FAILED" | "ESCALATE" | "UNKNOWN" {
  if (typeof text !== "string") return "UNKNOWN"
  const lines = text.replaceAll("\r\n", "\n").split("\n")
  const markers = lines.filter((line) => /^RTS_STATUS:/.test(line))
  if (markers.length !== 1) return "UNKNOWN"
  const match = /^RTS_STATUS: (SUCCESS|FAILED|ESCALATE)$/.exec(lines[0] ?? "")
  if (!match) return "UNKNOWN"
  return match[1] as "SUCCESS" | "FAILED" | "ESCALATE"
}

export function configuredAgentModel(agent: any): { providerID: string; modelID: string } | undefined {
  const model = agent?.model
  if (typeof model === "string") {
    const slash = model.indexOf("/")
    if (slash <= 0 || slash === model.length - 1 || model.indexOf("/", slash + 1) !== -1) return undefined
    return { providerID: model.slice(0, slash), modelID: model.slice(slash + 1) }
  }
  if (model && typeof model === "object" && typeof model.providerID === "string" && model.providerID && typeof model.modelID === "string" && model.modelID) {
    return { providerID: model.providerID, modelID: model.modelID }
  }
  return undefined
}

export function modelMatches(actual: any, expected: { providerID: string; modelID: string }): boolean {
  return !!actual && actual.providerID === expected.providerID && actual.modelID === expected.modelID
}

export function uniqueCurrentAssistant(messages: unknown, messageID: unknown): any | undefined {
  if (!Array.isArray(messages) || typeof messageID !== "string" || !messageID) return undefined
  const matches = messages.filter((message: any) => message?.info?.id === messageID && message.info.role === "assistant")
  if (matches.length !== 1) return undefined
  const info = matches[0].info
  if (typeof info.providerID !== "string" || !info.providerID || typeof info.modelID !== "string" || !info.modelID) return undefined
  return info
}

export function samePathSet(a: readonly string[], b: readonly string[]): boolean {
  return a.length === b.length && new Set(a).size === a.length && new Set(b).size === b.length && a.every((value) => b.includes(value))
}
