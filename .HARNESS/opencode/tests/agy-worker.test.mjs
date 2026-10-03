import assert from "node:assert/strict"
import { createHash } from "node:crypto"
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises"
import os from "node:os"
import path from "node:path"
import test from "node:test"

import {
  applyTransaction,
  containsCredentialLikeContent,
  executeAgyWorker,
  normalizeFileList,
  selectRoute,
} from "../tools/agy_worker.ts"

test("selectRoute fails closed when the required live model is absent", () => {
  assert.deepEqual(selectRoute("hard", ["gemini-3.1-pro-high"]), {
    model: "gemini-3.1-pro-high",
    effort: "high",
  })
  assert.throws(() => selectRoute("medium", ["gemini-3.7-flash-low"]), /required model unavailable/)
})

test("credential-like content is rejected", () => {
  assert.equal(containsCredentialLikeContent("const value = 42"), false)
  assert.equal(containsCredentialLikeContent("api_key = 'not-for-transmission'"), true)
  assert.equal(containsCredentialLikeContent("-----BEGIN PRIVATE KEY-----"), true)
})

test("normalizeFileList rejects protected and case-insensitive duplicate paths", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "agy-worker-test-"))
  try {
    await writeFile(path.join(root, "safe.txt"), "safe")
    await assert.rejects(() => normalizeFileList(root, [".env"], { allowMissing: true }), /protected path/)
    await assert.rejects(
      () => normalizeFileList(root, ["safe.txt", "SAFE.txt"], { allowMissing: true }),
      /case-insensitive duplicate/,
    )
  } finally {
    await rm(root, { recursive: true, force: true })
  }
})

test("applyTransaction writes only an allowed file and stops on collision", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "agy-worker-test-"))
  const target = path.join(root, "target.txt")
  try {
    await writeFile(target, "before")
    const allowed = await normalizeFileList(root, ["target.txt"], { allowMissing: true })
    const baseline = new Map([[
      allowed[0].key,
      { exists: true, hash: createHash("sha256").update("before").digest("hex"), content: "before" },
    ]])
    await applyTransaction([{ path: "target.txt", content: "after" }], allowed, baseline)
    assert.equal(await readFile(target, "utf8"), "after")

    await assert.rejects(
      () => applyTransaction([{ path: "target.txt", content: "again" }], allowed, baseline),
      /concurrent change detected/,
    )
  } finally {
    await rm(root, { recursive: true, force: true })
  }
})

test("dry-run validates a mocked agy response without writing the workspace", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "agy-worker-test-"))
  const target = path.join(root, "target.txt")
  try {
    await writeFile(target, "before")
    let call = 0
    const dependencies = {
      async runProcess() {
        call += 1
        if (call === 1) return { code: 0, stdout: "gemini-3.7-flash-low\tGemini 3.7 Flash (Low)\n" }
        return {
          code: 0,
          stdout: JSON.stringify({
            status: "SUCCESS",
            response: JSON.stringify({ status: "SUCCESS", summary: "changed", files: [{ path: "target.txt", content: "after" }] }),
          }),
        }
      },
    }
    const result = await executeAgyWorker(
      { task: "Replace the text", difficulty: "easy", allowed_files: ["target.txt"], context_files: [], dry_run: true },
      { directory: root, abort: new AbortController().signal },
      dependencies,
    )
    assert.match(result, /^SUCCESS: dry-run/)
    assert.equal(await readFile(target, "utf8"), "before")
  } finally {
    await rm(root, { recursive: true, force: true })
  }
})
