import { tool } from "@opencode-ai/plugin"
import orchestration from "../plugins/orchestration.ts"
import { realpath } from "node:fs/promises"

// Discovery adapters share the exact process registry with the bootstrap plugin.
function proxy(name: string, args: any) {
  return tool({ description: `RTS ${name}; shared canonical scheduler, queued natural-turn observation only.`, args, async execute(input, context) {
    const root = (await realpath(context.worktree)).toLowerCase()
    const candidates = await Promise.all([...orchestration.testing.registry.values()])
    const runtimes = candidates.filter((r) => r.root.toLowerCase() === root && !r.disposed)
    if (runtimes.length !== 1) return JSON.stringify({ status: "FAILED", error: "canonical scheduler missing or server identity ambiguous" })
    return runtimes[0].hooks().tool![`orchestration_${name}`].execute(input, context)
  } })
}
export const dispatch = proxy("dispatch", { role: tool.schema.string(), handoff: tool.schema.string(), scope_files: tool.schema.array(tool.schema.string()), write_files: tool.schema.array(tool.schema.string()).default([]), timeout_ms: tool.schema.number().optional() })
export const status = proxy("status", { task_id: tool.schema.string().optional() })
export const result = proxy("result", { task_id: tool.schema.string() })
export const cancel = proxy("cancel", { task_id: tool.schema.string() })
