import assert from "node:assert/strict"
import test from "node:test"
import { readFile, realpath } from "node:fs/promises"
import bootstrap from "../plugins/bootstrap.ts"
import orchestration from "../plugins/orchestration.ts"
import oldPlugin from "../../../.opencode/plugins/agy-worker-plugin.ts"
import inertPlugin from "../plugins/agy-worker-plugin.ts"
import agy from "../tools/agy_worker.ts"

test("legacy tools preserve named and default export identities; original three tests unchanged", async () => {
  for (const name of ["agy_worker", "codex_worker", "agy_council"]) {
    const canonical = await import(`../tools/${name}.ts`)
    const adapter = await import(`../../../.opencode/tools/${name}.ts`)
    assert.deepEqual(Object.keys(adapter), Object.keys(canonical))
    for (const k of Object.keys(canonical)) assert.equal(adapter[k], canonical[k])
    const testName = name.replaceAll("_", "-")
    assert.equal(await readFile(new URL(`./${testName}.test.mjs`, import.meta.url), "utf8"), await readFile(new URL(`../../../.opencode/tests/${testName}.test.mjs`, import.meta.url), "utf8"))
  }
})

test("bootstrap overlap registers only once; config preserves provider/global shell/model", async () => {
  const root = await realpath(new URL("../", import.meta.url))
  const input = { worktree: root, directory: root, serverUrl: new URL("http://localhost:9944"), client: {} }
  const hooks = await bootstrap(input)
  const provider = { untouched: { options: { baseURL: "http://localhost:1234" }, models: { live: { name: "live" } } } }
  const cfg = { provider, shell: "pwsh", model: "live/configured", plugin: ["existing-auth"], agent: { orchestrator: { mode: "subagent" } } }
  await hooks.config(cfg)
  assert.equal(cfg.provider, provider); assert.equal(cfg.shell, "pwsh"); assert.equal(cfg.model, "live/configured"); assert.deepEqual(cfg.plugin, ["existing-auth"])
  assert.equal(cfg.default_agent, "orchestrator"); assert.equal(cfg.agent.orchestrator.mode, "primary"); assert.equal(cfg.subagent_depth, 2)
  assert.equal(cfg.permission.agy_worker_current_task, "ask"); assert.equal(cfg.permission.task, "deny")
  assert.deepEqual(await oldPlugin(input), {}); assert.deepEqual(await orchestration(input), {}); assert.deepEqual(await inertPlugin(input), {})
  await hooks.dispose()
})

test("orchestration discovery first does not lose bootstrap config/tools or duplicate hooks", async () => {
  const root = await realpath(new URL("../", import.meta.url))
  const input = { worktree: root, directory: root, serverUrl: new URL("http://localhost:9945"), client: {} }
  const first = await orchestration(input), next = await bootstrap(input)
  assert.equal(typeof first.event, "function"); assert.equal(next.event, undefined)
  assert.equal(typeof first["chat.params"], "function")
  assert.equal(typeof next.config, "function"); assert.equal(next.tool.agy_worker, agy)
  assert.deepEqual(await bootstrap(input), {})
  await first.dispose()
})

test("agy actual tool boundary requires context.ask; supplied boolean cannot infer approval", async () => {
  let asks = 0
  await assert.rejects(() => agy.execute({ task: "bounded", allowed_files: ["a.txt"], approved: true }, { ask: async (request) => { asks++; assert.equal(request.permission, "agy_worker_current_task"); assert.deepEqual(request.always, []); throw new Error("approval denied") } }), /approval denied/)
  assert.equal(asks, 1)
})

test("canonical and legacy launchers/configs connect without global role drift", async () => {
  const globalConfig = await readFile(new URL("../../../.config/opencode/opencode.jsonc", import.meta.url), "utf8")
  assert.ok(globalConfig.includes(".harness/opencode/plugins/bootstrap.ts"))
  const adapter = await readFile(new URL("../../../.opencode/Invoke-OpenCode.ps1", import.meta.url), "utf8")
  assert.ok(adapter.includes("../.harness/opencode/Invoke-OpenCode.ps1"))
  const launcher = await readFile(new URL("../Invoke-OpenCode.ps1", import.meta.url), "utf8")
  assert.ok(launcher.includes("1.18.25")); assert.ok(launcher.includes("finally")); assert.ok(launcher.includes("$env:OPENCODE_CONFIG_DIR = $previousConfigDir"))
  const role = await readFile(new URL("../../../.config/opencode/agents/orchestrator.md", import.meta.url), "utf8")
  assert.ok(role.includes("mode: primary")); assert.ok(!role.includes("Y/N and waits"))
})
