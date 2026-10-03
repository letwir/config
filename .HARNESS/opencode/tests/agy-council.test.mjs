import assert from "node:assert/strict"
import test from "node:test"
import { AGGRESSIVE_ROSTER, executeAgyCouncil } from "../tools/agy_council.ts"

const modelIds = AGGRESSIVE_ROSTER.map((entry) => entry.model)
const successEnvelope = (summary = "ok", findings = []) => JSON.stringify({
  status: "SUCCESS",
  structured_output: { status: "SUCCESS", summary, findings },
})

function fakeTemps() {
  let next = 0
  const created = []
  const removed = []
  return {
    created,
    removed,
    makeTemp: async () => {
      const value = `temp-${++next}`
      created.push(value)
      return value
    },
    removeTemp: async (value) => { removed.push(value) },
  }
}

test("discovers once, consults the exact five-model roster, and caps concurrency at four", async () => {
  let discovery = 0
  let active = 0
  let peak = 0
  const calls = []
  const temps = fakeTemps()
  const result = JSON.parse(await executeAgyCouncil({ task: "Review this design" }, {}, {
    ...temps,
    runProcess: async (_command, args, options) => {
      if (args[0] === "models") {
        discovery++
        return { code: 0, stdout: modelIds.join("\n") }
      }
      active++
      peak = Math.max(peak, active)
      calls.push({ args, options })
      await new Promise((resolve) => setTimeout(resolve, 2))
      active--
      return { code: 0, stdout: successEnvelope() }
    },
  }))

  assert.equal(result.status, "SUCCESS")
  assert.equal(discovery, 1)
  assert.equal(calls.length, 5)
  assert.equal(peak, 4)
  assert.deepEqual(calls.map((call) => call.args[1]), modelIds)
  assert.equal(new Set(calls.map((call) => call.options.cwd)).size, 5)
  assert.deepEqual(temps.removed.sort(), temps.created.sort())
  for (const call of calls) {
    assert.ok(call.args.includes("--mode"))
    assert.ok(call.args.includes("plan"))
    assert.ok(call.args.includes("--output-format"))
    assert.ok(call.args.includes("json"))
    assert.ok(call.args.includes("--sandbox"))
    assert.ok(call.args.includes("--new-project"))
    assert.ok(call.args.includes("--disable-slash-commands"))
    assert.ok(call.args.some((arg) => arg.startsWith("--json-schema=")))
    assert.ok(call.args.some((arg) => arg.startsWith("--prompt=")))
    assert.ok(!call.args.includes("--add-dir"))
    assert.ok(!call.args.includes("--dangerously-skip-permissions"))
  }
})

test("missing and failed providers are explicit and make the aggregate ADVICE", async () => {
  const temps = fakeTemps()
  let providerCall = 0
  const result = JSON.parse(await executeAgyCouncil({ task: "Compare the approaches" }, {}, {
    ...temps,
    runProcess: async (_command, args) => {
      if (args[0] === "models") return { code: 0, stdout: `${modelIds[0]}\n${modelIds[1]}\n` }
      providerCall++
      return providerCall === 1 ? { code: 0, stdout: successEnvelope("useful") } : { code: 2, stdout: "ignored raw output" }
    },
  }))

  assert.equal(result.status, "ADVICE")
  assert.equal(result.results[0].status, "SUCCESS")
  assert.equal(result.results[1].status, "FAILED")
  assert.equal(result.results.filter((item) => item.status === "UNAVAILABLE").length, 3)
  assert.deepEqual(result.unavailable_models, modelIds.slice(2))
  assert.ok(!JSON.stringify(result).includes("ignored raw output"))
})

test("discovery failure fails closed without provider calls", async () => {
  const temps = fakeTemps()
  let calls = 0
  const result = JSON.parse(await executeAgyCouncil({ task: "Review" }, {}, {
    ...temps,
    runProcess: async () => { calls++; return { code: 1, stdout: "private diagnostic" } },
  }))
  assert.equal(calls, 1)
  assert.equal(result.status, "FAILED")
  assert.equal(result.results.length, 5)
  assert.ok(result.results.every((item) => item.status === "UNAVAILABLE"))
  assert.ok(!JSON.stringify(result).includes("private diagnostic"))
})

test("protected task text is rejected before model discovery", async () => {
  for (const task of [
    "api_key=secret-value",
    "Review C:\\private\\code.ts",
    "Review \\\\server\\share\\code.ts",
    "Review /private/code.ts",
    "Review ../private/code.ts",
    "Review file:///private/code.ts",
    "Review https://user:pass@example.com/code",
    "Review https://example.com/code?token=value",
  ]) {
    let calls = 0
    const result = JSON.parse(await executeAgyCouncil({ task }, {}, {
      runProcess: async () => { calls++; throw new Error("unexpected") },
    }))
    assert.equal(result.status, "FAILED", task)
    assert.equal(calls, 0, task)
  }
})

test("strict payload validation rejects extra fields and protected output", async () => {
  for (const payload of [
    { status: "SUCCESS", summary: "ok", findings: [], extra: true },
    { status: "SUCCESS", summary: "See C:\\private\\file.ts", findings: [] },
    { status: "SUCCESS", summary: "ok", findings: ["password=hidden"] },
    { status: "UNKNOWN", summary: "ok", findings: [] },
  ]) {
    const temps = fakeTemps()
    const result = JSON.parse(await executeAgyCouncil({ task: "Review" }, {}, {
      ...temps,
      runProcess: async (_command, args) => args[0] === "models"
        ? { code: 0, stdout: `${modelIds[0]}\n` }
        : { code: 0, stdout: JSON.stringify({ status: "SUCCESS", structured_output: payload }) },
    }))
    assert.equal(result.status, "FAILED")
    assert.equal(result.results[0].status, "FAILED")
    assert.ok(!JSON.stringify(result).includes("private"))
    assert.ok(!JSON.stringify(result).includes("hidden"))
  }
})

test("strict envelope validation rejects extra and ambiguous carriers", async () => {
  const payload = { status: "SUCCESS", summary: "ok", findings: [] }
  for (const envelope of [
    { status: "SUCCESS", structured_output: payload, extra: true },
    { status: "SUCCESS", structured_output: payload, response: payload },
    { status: "SUCCESS", response: payload, diagnostic: "raw" },
  ]) {
    const temps = fakeTemps()
    const result = JSON.parse(await executeAgyCouncil({ task: "Review" }, {}, {
      ...temps,
      runProcess: async (_command, args) => args[0] === "models"
        ? { code: 0, stdout: `${modelIds[0]}\n` }
        : { code: 0, stdout: JSON.stringify(envelope) },
    }))
    assert.equal(result.status, "FAILED")
    assert.equal(result.results[0].status, "FAILED")
    assert.ok(!JSON.stringify(result).includes("diagnostic"))
  }
})

test("a second invocation fails immediately while the council lock is held", async () => {
  const temps = fakeTemps()
  let releaseDiscovery
  const blocked = executeAgyCouncil({ task: "First review" }, {}, {
    ...temps,
    runProcess: async (_command, args) => {
      if (args[0] === "models") {
        await new Promise((resolve) => { releaseDiscovery = resolve })
        return { code: 0, stdout: "" }
      }
      throw new Error("unexpected")
    },
  })
  while (!releaseDiscovery) await new Promise((resolve) => setTimeout(resolve, 1))
  const second = JSON.parse(await executeAgyCouncil({ task: "Second review" }, {}))
  assert.equal(second.status, "FAILED")
  assert.match(second.error, /another council invocation/)
  releaseDiscovery()
  await blocked
})

test("provider timeouts are clamped to the remaining global deadline", async () => {
  const temps = fakeTemps()
  let now = 10_000
  const timeouts = []
  const result = JSON.parse(await executeAgyCouncil({ task: "Review", deadline_ms: 2_000 }, {}, {
    ...temps,
    now: () => now,
    runProcess: async (_command, args, options) => {
      timeouts.push(options.timeoutMs)
      if (args[0] === "models") {
        now += 750
        return { code: 0, stdout: `${modelIds[0]}\n` }
      }
      return { code: 0, stdout: successEnvelope() }
    },
  }))
  assert.equal(result.status, "ADVICE")
  assert.equal(timeouts[0], 2_000)
  assert.equal(timeouts[1], 1_250)
})
