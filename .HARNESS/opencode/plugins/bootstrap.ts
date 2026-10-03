import { readFile, realpath } from "node:fs/promises"
import { fileURLToPath } from "node:url"
import path from "node:path"
import type { Plugin } from "@opencode-ai/plugin"
import orchestration from "./orchestration.ts"
import agyWorker from "../tools/agy_worker.ts"
import codexWorker from "../tools/codex_worker.ts"
import agyCouncil from "../tools/agy_council.ts"

const configured = Symbol.for("harness.rts.bootstrap.configured.v1")
const claims: Set<string> = (globalThis as any)[configured] ??= new Set()
const configPath = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../opencode.jsonc")

const bootstrap: Plugin = async (input) => {
  const hooks = await orchestration(input)
  const identity = `${input.serverUrl}|${(await realpath(input.worktree || input.directory)).toLowerCase()}`
  if (claims.has(identity)) return hooks
  claims.add(identity)
  return {
    ...hooks,
    tool: { ...hooks.tool, agy_worker: agyWorker, codex_worker: codexWorker, agy_council: agyCouncil },
    config: async (cfg) => {
      // Canonical file is strict JSON despite .jsonc extension. Never load a config/plugin recursively.
      const canonical = JSON.parse(await readFile(configPath, "utf8"))
      for (const name of ["default_agent", "subagent_depth", "agent", "share", "snapshot", "instructions", "skills", "watcher", "permission"]) {
        if (Object.hasOwn(canonical, name)) cfg[name] = structuredClone(canonical[name])
      }
      const role = await readFile(path.resolve(path.dirname(configPath), "agents/orchestrator.md"), "utf8")
      cfg.agent!.orchestrator.prompt = role.replace(/^---\r?\n[\s\S]*?\r?\n---\r?\n/, "")
      // provider, credentials, shell, model and unrelated global configuration remain untouched.
    },
  }
}
export default bootstrap
