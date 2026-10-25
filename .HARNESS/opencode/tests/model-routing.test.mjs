import assert from "node:assert/strict"
import { readFile } from "node:fs/promises"
import test from "node:test"
import {
  ASTRA_MODEL,
  HARD_SOL_MODEL,
  configuredAgentModel,
  modelMatches,
  parseWorkerStatus,
  samePathSet,
  selectWorkerProfile,
  uniqueCurrentAssistant,
  selectAttackerProfile,
  reviewIdentity,
  hasReasoningInsufficiency,
  assertAttackerEscalation,
} from "../tools/model_routing.ts"

const mutationAndChildTools = ["edit", "write", "write_files", "apply_patch", "task", "terminal_session", "bash_output", "bash_kill", "bash_restart"]

// Crabcode uses ordered last-match-wins wildcard rules, not first-match ACLs.
function filePermission(rules, path) {
  let result
  for (const [pattern, action] of Object.entries(rules)) {
    const regex = pattern.replace(/[.+^${}()|[\]\\]/g, "\\$&").replaceAll("*", ".*").replaceAll("?", ".")
    if (new RegExp(`^${regex}$`).test(path)) result = action
  }
  return result
}

function assertSourceReviewRules(permission) {
  for (const tool of ["read", "grep"]) {
    const rules = permission[tool]
    const entries = Object.entries(rules)
    assert.deepEqual(entries[0], ["*", "allow"], `${tool}: broad source allow first`)
    assert.deepEqual(entries.slice(-2), [[".env.example", "allow"], ["*/.env.example", "allow"]], `${tool}: only env example exceptions last`)
    assert.ok(entries.length > 3)
    for (const [pattern, action] of entries.slice(1, -2)) assert.equal(action, "deny", `${tool}: ${pattern}`)
    for (const pattern of ["*.env", "*.env.*", "*credential*", "*secret*", "*auth*", "*token*", "*.key", "*.pem", "*.p12", "*.pfx", "*.crt", "*.cer", "*.der"]) {
      assert.equal(rules[pattern], "deny", `${tool}: missing sensitive pattern ${pattern}`)
      assert.ok(Object.keys(rules).indexOf(pattern) > Object.keys(rules).indexOf("*"))
    }
    for (const path of ["src/parser.ts", "/checkout/src/parser.ts", "README.md"]) assert.equal(filePermission(rules, path), "allow", `${tool}: ${path}`)
    for (const path of [".env", "app.env", "/checkout/.env.local", ".env.example.local", "app.env.example", "credentials.json", "secrets.yaml", "auth.json", "token.txt", "password.txt", "private.key", "cert.pem", "bundle.p12", "bundle.pfx", "cert.crt", "cert.cer", "cert.der", "store.jks", "store.keystore", "bundle.pkcs12", "private.pkcs8", "private.p8", "store.kdbx", "private-key.json", "private_key.json", "id_ed25519", ".ssh/config", ".aws/config", ".kube/config", "/home/test/.ssh/config", "/home/test/.aws/config", "/home/test/.kube/config", ".npmrc", ".netrc", ".git-credentials"]) {
      assert.equal(filePermission(rules, path), "deny", `${tool}: ${path}`)
    }
    for (const path of [".env.example", "/checkout/.env.example", "/checkout/auth/.env.example"]) assert.equal(filePermission(rules, path), "allow", `${tool}: ${path}`)
  }
  assert.deepEqual(permission.read, permission.grep)
}

test("attacker keeps client-specific permissions, fixed routes and the OpenCode blackhat alias", async () => {
  const configs = await Promise.all(["../../crabcode/crabcode.jsonc", "../opencode.jsonc"].map(async (file) =>
    JSON.parse(await readFile(new URL(file, import.meta.url), "utf8"))))
  const inspectOnly = { "*": "deny", read: "allow", grep: "allow", edit: "deny", write: "deny", write_files: "deny", apply_patch: "deny", task: "deny" }
  for (const [client, config] of configs.entries()) {
    const profile = config.agent.attacker
    assert.equal(profile.mode, "subagent")
    assert.equal(profile.prompt, "{file:../agents/blackhat.md}")
    if (client === 1) assert.deepEqual(profile.permission, { ...inspectOnly, read: config.permission.read })
    assert.equal(profile.model, "openai/gpt-6-luna")
    assert.equal(profile.reasoningEffort, "low")
    for (const depth of ["low", "medium", "hard"]) {
      const selected = selectAttackerProfile(depth, depth === "hard" ? "sol-review" : undefined)
      const tier = config.agent[selected.agent]
      assert.equal(tier.mode, "subagent")
      assert.equal(tier.model, `${selected.providerID}/${selected.modelID}`)
      assert.equal(tier.reasoningEffort, selected.reasoningEffort)
      assert.equal(tier.prompt, profile.prompt)
      assert.deepEqual(tier.permission, profile.permission)
    }
    if (client === 0) {
      assert.deepEqual(config.disabled_providers, ["opencode"])
      for (const name of ["attacker-low", "attacker-medium", "attacker-hard", "attacker"]) {
        const permission = config.agent[name].permission
        assert.equal(permission["*"], "deny")
        for (const tool of ["list", "glob", "external_directory", "bash"]) assert.equal(permission[tool], "allow", `${name}: ${tool} must be unrestricted`)
        for (const tool of mutationAndChildTools) assert.equal(permission[tool], "deny", `${name}: ${tool}`)
        assertSourceReviewRules(permission)
      }
    }
    assert.match(profile.description, /requirements.*detailed designs.*CODE.*Verifier/)
  }
  assert.deepEqual(configs[1].agent.blackhat.permission, configs[1].agent.attacker.permission)
  assert.equal(configs[1].agent.blackhat.prompt, configs[1].agent.attacker.prompt)
  assert.equal(configs[1].agent.blackhat.mode, "subagent")
  assert.equal(configs[1].agent.blackhat.model, "openai/gpt-6-luna")
  assert.equal(configs[1].agent.blackhat.reasoningEffort, "low")
  const contract = await readFile(new URL("../../agents/blackhat.md", import.meta.url), "utf8")
  assert.match(contract, /^description:.*Crabcode bounded reproduction.*annotation.*actual Researcher CVE triage.*explicit parent proceed confirmation.*OpenCode remains read-only under its config\./m)
  assert.match(contract, /^name: attacker$/m)
  for (const field of ["severity", "surface", "evidence", "assumptions_preconditions", "impact", "recommended_smallest_defensive_correction", "confidence"]) {
    assert.match(contract, new RegExp(`\\b${field}:`))
  }
  assert.match(contract, /NO_FINDINGS only when scope was actually assessed/)
  assert.match(contract, /ESCALATE\+INCOMPLETE/)
  assert.match(contract, /FAILED\+INCOMPLETE/)
  const toml = await readFile(new URL("../../agents/blackhat.toml", import.meta.url), "utf8")
  assert.match(toml, /^name = "attacker"$/m)
  assert.match(toml, /^sandbox_mode = "read-only"$/m)
  assert.doesNotMatch(toml, /^model\s*=/m)
})

test("Crabcode Researcher is read-only public CVE triage without shell, children or a model pin", async () => {
  const config = JSON.parse(await readFile(new URL("../../crabcode/crabcode.jsonc", import.meta.url), "utf8"))
  const researcher = config.agent.researcher
  assert.equal(researcher.mode, "subagent")
  assert.equal(researcher.prompt, "{file:../agents/researcher.md}")
  assert.equal(Object.hasOwn(researcher, "model"), false)
  assert.equal(Object.hasOwn(researcher, "reasoningEffort"), false)
  assert.equal(researcher.permission["*"], "deny")
  for (const tool of ["list", "glob", "websearch", "webfetch"]) assert.equal(researcher.permission[tool], "allow")
  for (const tool of [...mutationAndChildTools, "bash"]) assert.equal(researcher.permission[tool], "deny")
  assertSourceReviewRules(researcher.permission)
  assert.match(researcher.description, /CVE_CANDIDATE.*NOT_CVE_CANDIDATE.*INSUFFICIENT_EVIDENCE/)
  assert.match(researcher.description, /Never reproduce, submit CVE, contact CNA\/vendor/)
})

test("annotated verification requires actual Researcher triage and post-triage parent proceed, never submission approval", async () => {
  const contract = await readFile(new URL("../../agents/blackhat.md", import.meta.url), "utf8")
  assert.match(contract, /discovery -> annotation -> Researcher triage -> parent confirms\/no-go -> only then bounded repro -> Attacker report -> separate Verifier/)
  assert.match(contract, /numbered preflight annotation/)
  for (const field of ["finding ID", "exact target and revision", "evidence/source lines", "suspected weakness and preconditions", "expected signal", "impact and blast radius", "data classification", "exact commands/actions", "containment/isolation", "cleanup/rollback", "writes/network/credential inventory", "applicable current-task authorization"]) assert.ok(contract.includes(field), field)
  assert.match(contract, /CVE_CANDIDATE \| NOT_CVE_CANDIDATE \| INSUFFICIENT_EVIDENCE/)
  for (const field of ["known-CVE/advisory dedup", "affected versions", "severity rationale", "eligibility/recordability evidence", "missing info", "cited public sources"]) assert.ok(contract.includes(field), field)
  assert.match(contract, /no reproduction until report includes annotation AND actual Researcher result AND parent explicitly confirms proceed/)
  assert.match(contract, /INSUFFICIENT_EVIDENCE => stop\/no reproduction/)
  assert.match(contract, /NOT_CVE_CANDIDATE => do not market\/report as CVE.*parent confirmation.*explicitly named owned\/authorized target/)
  assert.match(contract, /CVE_CANDIDATE => verification only after parent confirmation/)
  assert.match(contract, /No public CVE submission\/disclosure\/CNA contact without separate explicit approval/)
  assert.match(contract, /Public read-only CVE review has no public write effect/)
  assert.match(contract, /NOT a hard sandbox/)
  assert.match(contract, /Discovery and ordinary review are read-only; Crabcode Attacker Bash can run explicitly gated bounded verification/)
  const subagents = await readFile(new URL("../../agents/SUBAGENTS.lrf", import.meta.url), "utf8")
  assert.match(subagents, /^R\|subagent\.blackhat\|.*no promise of vulnerability discovery; bounded exploit reproduction only under subagent\.attacker-verification and the blackhat contract; no zero-day guarantee/m)
  assert.doesNotMatch(subagents, /^\s+R\|/m)
  assert.doesNotMatch(subagents, /no exploit-test or zero-day-discovery promises/)
  assert.match(contract, /Attacker reports, never fixes/)
  assert.match(contract, /Attacker task remains denied\/leaf-only.*parent, which dispatches/)
  assert.match(contract, /Missing\/unavailable Researcher⇒stop\/no reproduction/)
  assert.match(contract, /exact current-task target AND effect approval; no scanning beyond it/)
  assert.match(contract, /prohibited_regardless: credentials\/secrets access or exfiltration; stealth\/persistence; destructive data changes; unauthorized escalation; uncontrolled propagation/)
})

test("attacker gates cover proposals and security design/diff without replacing independent Verifier", async () => {
  const research = await readFile(new URL("../../etl/research.seq", import.meta.url), "utf8")
  const change = await readFile(new URL("../../etl/change.seq", import.meta.url), "utf8")
  assert.match(research, /requirements_or_detailed_design_proposals_in_scope.*BEFORE_finalization\/approval/)
  assert.match(research, /non_security_without_proposals⇒no_needless_Attacker/)
  assert.match(research, /record_missing_independent_assessment/)
  const codeGate = change.indexOf("PASS ∧ required_user_authorization ∧ explicit_effect_approvals ⇒ CODE")
  const designGate = change.indexOf("Attacker(approved_design_and_security_assumptions)_before_CODE")
  assert.ok(designGate >= 0 && codeGate > designGate)
  assert.match(change, /Attacker\(final_integrated_diff.*BEFORE_final_independent_Verifier/)
  assert.match(change, /independent_Verifier_is_neither_implementation_author_nor_Attacker_assessor/)
  assert.match(change, /all_review_failures_share_correction_cycles<=3_across_proposal\/design\/implementation_reviews/)
  assert.match(change, /never_claim_Attacker_ran/)
})

test("routed policies and client documentation retain the complete forensic gate sequence", async () => {
  for (const file of ["../../agents/SUBAGENTS.lrf", "../../rules/MANUAL.lrf", "../../rules/engineering.lrf", "../../agents/AGENT_ROUTER.md", "../../etl/change.seq"]) {
    const policy = await readFile(new URL(file, import.meta.url), "utf8")
    assert.match(policy, /discovery -> annotation -> Researcher triage -> parent confirms\/no-go -> only then bounded repro -> Attacker report -> separate Verifier/, file)
    assert.match(policy, /annotation AND actual Researcher result AND explicit post-triage parent proceed confirmation/, file)
    assert.match(policy, /INSUFFICIENT_EVIDENCE or unavailable Researcher\/runtime=>stop\/no reproduction/, file)
    assert.match(policy, /public CVE submission\/disclosure\/CNA contact requires separate explicit approval/, file)
  }
  const readme = await readFile(new URL("../README.md", import.meta.url), "utf8")
  assert.match(readme, /Discovery → annotation → Researcher triage → parent confirms\/no-go → only then bounded repro → Attacker report → separate Verifier/)
  assert.match(readme, /OpenCode runtime\/config permissions are unchanged/)
  assert.match(readme, /OpenCode configuration and sandbox remain read-only, with no bash-based reproduction/)
  assert.match(readme, /Bounded reproduction is allowed only after the annotation, actual Researcher CVE triage, and explicit post-triage parent proceed gate/)
  assert.match(readme, /no exploit or zero-day success is promised/)
  assert.doesNotMatch(readme, /No exploit tests or zero-day discovery are promised/)
  assert.match(readme, /not a hard sandbox/)
  assert.match(readme, /static-only read-only review; no edits\/reproduction\/external actions\/children/)
  for (const [file, record] of [["../../agents/SUBAGENTS.lrf", "subagent.attacker-verification"], ["../../rules/MANUAL.lrf", "agent.attacker-verification"], ["../../rules/engineering.lrf", "workflow.attacker-verification"]]) {
    const policy = await readFile(new URL(file, import.meta.url), "utf8")
    assert.match(policy, new RegExp(`R\\|${record}\\|[^\\n]*\\|LW_SCOPE\\|`))
  }
})

test("worker difficulty maps to fixed static profiles and requested efforts", () => {
  assert.deepEqual(selectWorkerProfile("easy"), { agent: "worker-easy", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "low" })

  assert.deepEqual(selectWorkerProfile("medium"), { agent: "worker-medium", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "medium" })
  assert.deepEqual(selectWorkerProfile("hard"), { agent: "worker-hard", providerID: "openai", modelID: "gpt-6.1-sol", reasoningEffort: "high" })
})


test("Crabcode and OpenCode static profile configs match the selected route and shared prompt", async () => {
  const crabcode = JSON.parse(await readFile(new URL("../../crabcode/crabcode.jsonc", import.meta.url), "utf8"))
  const opencode = JSON.parse(await readFile(new URL("../opencode.jsonc", import.meta.url), "utf8"))
  assert.deepEqual(Object.keys(crabcode.provider), ["tailnet-llama"])
  assert.equal(Object.hasOwn(crabcode, "model"), false)
  assert.equal(Object.hasOwn(opencode, "model"), false)
  const expected = {
    "worker-easy": ["openai/gpt-6-luna", "low"],
    "worker-medium": ["openai/gpt-6-luna", "medium"],
    "worker-hard": ["openai/gpt-6.1-sol", "high"],
    "worker-astra": ["openai/gpt-6-astra", undefined],
  }
  for (const [name, [model, effort]] of Object.entries(expected)) {
    for (const agents of [crabcode.agent, opencode.agent]) {
      const profile = agents[name]
      assert.equal(profile.mode, "subagent")
      assert.equal(profile.model, model)
      assert.equal(profile.reasoningEffort, effort)
      assert.equal(profile.permission.task, "deny")
    }
  }
  const crabcodePrompt = crabcode.agent["worker-easy"].prompt
  const opencodePrompt = opencode.agent["worker-easy"].prompt
  assert.equal(crabcodePrompt, "{file:./prompts/worker-prompt.md}")
  assert.equal(opencodePrompt, "{file:../crabcode/prompts/worker-prompt.md}")
  const crabcodePromptPath = new URL(`../../crabcode/${crabcodePrompt.slice("{file:./".length, -1)}`, import.meta.url)
  const opencodePromptPath = new URL("../../crabcode/prompts/worker-prompt.md", import.meta.url)
  const prompt = await readFile(crabcodePromptPath, "utf8")
  assert.equal(await readFile(opencodePromptPath, "utf8"), prompt)
  assert.match(prompt, /^# Bounded CODE Worker Prompt/m)
  assert.match(prompt, /RTS_STATUS: SUCCESS.*RTS_STATUS: FAILED.*RTS_STATUS: ESCALATE/s)
})

test("Astra is available only for explicit hard-task escalation and uses provider default effort", () => {
  assert.deepEqual(selectWorkerProfile("hard", "sol-task-1"), { agent: "worker-astra", providerID: "openai", modelID: "gpt-6-astra" })
  assert.equal(Object.hasOwn(selectWorkerProfile("hard", "sol-task-1"), "reasoningEffort"), false)
  assert.deepEqual(ASTRA_MODEL, { providerID: "openai", modelID: "gpt-6-astra" })
  assert.deepEqual(HARD_SOL_MODEL, { providerID: "openai", modelID: "gpt-6.1-sol" })
  assert.throws(() => selectWorkerProfile("medium", "sol-task-1"), /hard difficulty/)
  assert.throws(() => selectWorkerProfile("hard", ""), /source task/)
})

test("difficulty and terminal status parsing fail closed", () => {
  assert.throws(() => selectWorkerProfile("unknown"), /difficulty/)
  assert.equal(parseWorkerStatus("RTS_STATUS: SUCCESS\nDone"), "SUCCESS")
  assert.equal(parseWorkerStatus("RTS_STATUS: FAILED\nReason"), "FAILED")
  assert.equal(parseWorkerStatus("RTS_STATUS: ESCALATE\nReason"), "ESCALATE")
  for (const text of ["Done", "prefix\nRTS_STATUS: FAILED", "RTS_STATUS: FAILED\nRTS_STATUS: SUCCESS", "RTS_STATUS: FAILED extra", " RTS_STATUS: FAILED", undefined]) {
    assert.equal(parseWorkerStatus(text), "UNKNOWN")
  }
})

test("configured agent model metadata accepts exact string/object forms only", () => {
  assert.deepEqual(configuredAgentModel({ model: "openai/gpt-6-luna" }), { providerID: "openai", modelID: "gpt-6-luna" })
  assert.deepEqual(configuredAgentModel({ model: { providerID: "openai", modelID: "gpt-6.1-sol" } }), { providerID: "openai", modelID: "gpt-6.1-sol" })
  for (const agent of [{}, { model: "gpt-6-luna" }, { model: "openai/" }, { model: { providerID: "openai" } }]) assert.equal(configuredAgentModel(agent), undefined)
  assert.equal(modelMatches({ providerID: "openai", modelID: "gpt-6-luna" }, { providerID: "openai", modelID: "gpt-6-luna" }), true)
  assert.equal(modelMatches({ providerID: "openai", modelID: "gpt-6-astra" }, { providerID: "openai", modelID: "gpt-6-luna" }), false)
})

test("parent model identity requires one exact correlated assistant message", () => {
  const message = { info: { id: "parent-message", role: "assistant", providerID: "openai", modelID: "gpt-6-luna" } }
  assert.deepEqual(uniqueCurrentAssistant([message], "parent-message"), message.info)
  assert.equal(uniqueCurrentAssistant([message], "stale-message"), undefined)
  assert.equal(uniqueCurrentAssistant([message, message], "parent-message"), undefined)
  assert.equal(uniqueCurrentAssistant([{ info: { ...message.info, modelID: "" } }], "parent-message"), undefined)
  assert.equal(uniqueCurrentAssistant([], "parent-message"), undefined)
})

test("scope/write sets compare order-independently and reject aliases duplicated after canonicalization", () => {
  assert.equal(samePathSet(["/repo/a", "/repo/b"], ["/repo/b", "/repo/a"]), true)
  assert.equal(samePathSet(["/repo/a"], ["/repo/b"]), false)
  assert.equal(samePathSet(["/repo/a", "/repo/a"], ["/repo/a", "/repo/a"]), false)
})

const insufficiency = (patch = {}) => "RTS_STATUS: ESCALATE\nRTS_ESCALATION: " + JSON.stringify({
  cause: "reasoning_insufficiency", bounded_assessment_completed: true,
  further_reasoning_needed: true, blocked_by: [], unresolved_analysis: "Two bounded path assumptions remain inconsistent", ...patch,
}) + "\nassessed_scope: proposal@rev-1, trust-boundary; findings retained"
const solReview = () => ({ parent: "parent", root: "root", depth: 1, role: "attacker", profile: "attacker-medium",
  reviewDepth: "medium", review: reviewIdentity("proposal", "rev-1", "trust-boundary"),
  actualModel: HARD_SOL_MODEL, routeStatus: "matched", stopped: true, resultLoaded: true, state: "completed",
  resultStatus: "ESCALATE", result: insufficiency(), scope: ["/repo/proposal.md"], writes: [],
})
const hardRequest = () => ({ parent: "parent", root: "root", depth: 1, scope: ["/repo/proposal.md"], writes: [],
  review: reviewIdentity("proposal", "rev-1", "trust-boundary"), approved: true,
})

test("review profiles are independent of worker difficulty and hard never routes directly", () => {
  assert.deepEqual(selectAttackerProfile(), { agent: "attacker-low", providerID: "openai", modelID: "gpt-6-luna", reasoningEffort: "low" })
  assert.deepEqual(selectAttackerProfile("medium"), { agent: "attacker-medium", ...HARD_SOL_MODEL, reasoningEffort: "high" })
  assert.deepEqual(selectAttackerProfile("hard", "sol-review"), { agent: "attacker-hard", ...ASTRA_MODEL })
  for (const depth of ["easy", "unknown", null]) assert.throws(() => selectAttackerProfile(depth))
  for (const id of [undefined, "", " "]) assert.throws(() => selectAttackerProfile("hard", id))
  assert.throws(() => selectAttackerProfile("low", "source"))
  assert.throws(() => selectAttackerProfile("medium", "source"))
})

test("review identity rejects absent, multiline, unbounded and control-delimited fields", () => {
  for (const value of [undefined, "", " ", "a\nb", "a\rb", "a\0b", "a\tb", "a\x7fb", "x".repeat(1001)]) {
    for (let index = 0; index < 3; index++) {
      const fields = ["proposal", "rev-1", "lens"]; fields[index] = value
      assert.throws(() => reviewIdentity(...fields))
    }
  }
})

test("only exact complete blocker-free reasoning evidence is eligible", () => {
  assert.equal(hasReasoningInsufficiency(insufficiency()), true)
  for (const patch of [{ cause: "permission" }, { blocked_by: ["auth"] }, { bounded_assessment_completed: false },
    { further_reasoning_needed: false }, { unresolved_analysis: " " }, { extra: true }]) {
    assert.equal(hasReasoningInsufficiency(insufficiency(patch)), false)
  }
  for (const text of [undefined, "RTS_STATUS: ESCALATE", insufficiency().replace("ESCALATE", "FAILED"),
    "prefix\n" + insufficiency(), insufficiency() + "\nRTS_STATUS: ESCALATE",
    insufficiency() + "\nRTS_ESCALATION: {}", "RTS_STATUS: ESCALATE\nRTS_ESCALATION: {bad}",
    insufficiency() + "x".repeat(64000)]) assert.equal(hasReasoningInsufficiency(text), false)
})

test("attacker Astra requires correlated actual completed Sol, approval, exact identity/scope and no replay", () => {
  assert.doesNotThrow(() => assertAttackerEscalation(solReview(), hardRequest()))
  assert.throws(() => assertAttackerEscalation(undefined, hardRequest()))
  for (const patch of [{ parent: "other" }, { root: "other" }, { depth: 2 }, { role: "worker" },
    { profile: "attacker-low" }, { reviewDepth: "low" }, { actualModel: ASTRA_MODEL }, { actualModel: undefined },
    { routeStatus: "configured" }, { stopped: false }, { resultLoaded: false }, { state: "cancelled" },
    { resultStatus: "FAILED" }, { result: insufficiency({ blocked_by: ["approval"] }) },
    { error: "IO failed" }, { escalationReservedBy: "already-used" }, { writes: ["/repo/proposal.md"] },
    { scope: ["/repo/other.md"] }, { review: undefined },
    ...["target", "revision", "lens"].map((key) => ({ review: { ...solReview().review, [key]: "different" } })),
  ]) assert.throws(() => assertAttackerEscalation({ ...solReview(), ...patch }, hardRequest()))
  for (const patch of [{ approved: false }, { approved: undefined }, { writes: ["/repo/proposal.md"] }, { scope: [] }, { review: undefined }]) {
    assert.throws(() => assertAttackerEscalation(solReview(), { ...hardRequest(), ...patch }))
  }
})

test("PG handoff example retains exact envelope and ordered labels with recoverable value anchors", async () => {
  const worker = await readFile(new URL("../../agents/worker.md", import.meta.url), "utf8")
  const example = worker.match(/<Γ>\nSIGMA\/1: PIDGEN\/text\nRole: worker;[\s\S]*?\n<\/Γ>/)?.[0]
  assert.ok(example)
  assert.deepEqual([...example.matchAll(/^(Role|Target|Acceptance|Scope|Known facts):/gm)].map((m) => m[1]), ["Role", "Target", "Acceptance", "Scope", "Known facts"])
  assert.match(example, /\nSIGMA\/1\n<\/Γ>$/)
  for (const literal of ["source task ID=task-parser-17", "revision=rev-42", "model=openai/gpt-6.1-sol", "status=plan-approved", "no external/VCS writes", "🧪=local regression check"]) assert.ok(example.includes(literal))
})
