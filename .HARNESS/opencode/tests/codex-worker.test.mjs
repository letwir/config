import assert from "node:assert/strict"
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import os from "node:os"
import path from "node:path"
import test from "node:test"
import { containsCredentialLikeContent, executeCodexWorker, normalizeFileList, selectRoute } from "../tools/codex_worker.ts"

test("routes difficulty to the current Codex models", () => {
  assert.deepEqual(selectRoute("easy"), { model: "gpt-5.6-luna", effort: "low" })
  assert.deepEqual(selectRoute("medium"), { model: "gpt-5.6-terra", effort: "medium" })
  assert.deepEqual(selectRoute("hard"), { model: "gpt-5.6-sol", effort: "xhigh" })
})

test("rejects protected, traversal, binary, and credential-like context", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "codex-worker-test-"))
  try {
    await writeFile(path.join(root, "safe.txt"), "safe")
    await writeFile(path.join(root, "binary.bin"), Buffer.from([0, 1, 2]))
    assert.equal(containsCredentialLikeContent("api_key='x'"), true)
    await assert.rejects(() => normalizeFileList(root, [".env"]), /protected path/)
    await assert.rejects(() => normalizeFileList(root, ["..\\outside.txt"]), /path escapes/)
    const result = await executeCodexWorker({ task: "inspect", difficulty: "easy", context_files: ["binary.bin"] }, { directory: root, abort: new AbortController().signal }, { resolveExecutable: async () => "codex.exe", runProcess: async () => ({ code: 0, stdout: "{}" }) })
    assert.match(result, /binary or oversized/)
  } finally { await rm(root, { recursive: true, force: true }) }
})

test("uses injected Codex process and leaves workspace unchanged", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "codex-worker-test-")); const file = path.join(root, "context.txt")
  try {
    await writeFile(file, "context")
    const calls = []
    const result = await executeCodexWorker({ task: "advise", difficulty: "medium", context_files: ["context.txt"] }, { directory: root, abort: new AbortController().signal }, { resolveExecutable: async () => "C:\\Codex\\codex.exe", runProcess: async (...args) => { calls.push(args); return { code: 0, stdout: JSON.stringify({ status: "ADVICE", summary: "review carefully", files: [] }) } } })
    assert.equal(result, "ADVICE: review carefully"); assert.equal(calls.length, 1); assert.equal(calls[0][0], "C:\\Codex\\codex.exe"); assert.equal(calls[0][1][0], "exec"); assert.equal(calls[0][1].includes("--sandbox"), true); assert.equal(calls[0][1].includes('model_reasoning_effort="medium"'), true); assert.equal(await readFile(file, "utf8"), "context")
  } finally { await rm(root, { recursive: true, force: true }) }
})
