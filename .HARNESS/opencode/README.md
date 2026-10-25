# Canonical OpenCode RTS integration

Configuration, tools, launcher and scheduler live in `.harness/opencode`. Installed `@opencode-ai/plugin` and root SDK **1.18.25** are reused through the approved `node_modules` junction to `.opencode/node_modules`; no dependency installation is needed. Global provider configuration, shell and credential storage are unchanged. Legacy tools export named/default adapters; legacy launcher forwards here. Quit and restart OpenCode after changing configuration.

## Run

```powershell
& "$env:USERPROFILE/.harness/opencode/Invoke-OpenCode.ps1"
# Tests (all process runners / SDK calls are mocked; no live model calls)
$env:TEMP = 'A:/TMP/opencode'; $env:TMP = $env:TEMP
node --experimental-strip-types --test opencode/tests/*.test.mjs
```

The launcher checks CLI 1.18.25, sets `OPENCODE_CONFIG` and `OPENCODE_CONFIG_DIR` for its invocation, and restores both. Bootstrap copies canonical behavior/permissions without recursion, preserving providers, model and shell. Discovery overlap shares a process singleton and registers hooks once.

### Global plugin auto-discovery

OpenCode automatically loads local plugin files from `~/.config/opencode/plugins/` on ordinary startup. The config repository shares the canonical implementation through these relative symlinks (targets are relative to each link's parent directory):

| Repository link | Relative target |
| --- | --- |
| `.config/opencode/plugins/bootstrap.ts` | `../../../.HARNESS/opencode/plugins/bootstrap.ts` |
| `.config/opencode/package.json` | `../../.HARNESS/opencode/package.json` |

The `package.json` link exposes the canonical plugin dependency declarations in OpenCode's config directory; it is not a second dependency manifest to maintain. OpenCode uses that manifest to install plugin dependencies at startup, so dependency availability and any installation effects must be handled separately from creating these links. Keep the plugin as a symlink, not a copied file: real-path module resolution keeps its sibling imports and `import.meta.url`-based canonical config lookup together under `.HARNESS/opencode` (the deployed runtime's `.harness/opencode`).

Do not also add a `plugin` array entry for this bootstrap: global directory auto-discovery replaces the old Windows-only `file:///C:/Users/letwir/.harness/opencode/plugins/bootstrap.ts` entry. Canonical `opencode.jsonc` uses the config-file-relative instruction path `../AGENTS.md`, resolving to the Harness root's `AGENTS.md` on Linux and Windows. Provider and model settings are unchanged.

The forward-slash relative targets are portable, but deployment must preserve the `.config`/Harness relative layout (including the repository's `.HARNESS` spelling, or consistently mapping it to runtime `.harness`). Windows Git checkouts require symlink capability (for example Developer Mode or appropriate privileges) and `core.symlinks=true`; a checkout containing plain text target files will not load the plugin. Refuse to overwrite existing link paths when installing. Restart OpenCode after setup; static checks do not prove live auto-loading. Crabcode intentionally ignores OpenCode `plugin` settings and does not load this plugin; use Crabcode's native `Task` profiles instead.

## Dispatch and observation

Use `orchestration_dispatch` with `role`, exact `handoff`, `scope_files`, `write_files` and optional bounded `timeout_ms`. For `role=worker`, `difficulty` is required (`easy|medium|hard`); Attacker uses separate review-depth profiles below; other specialist roles retain their existing routing. Worker profiles are static in both Crabcode/OpenCode configuration: easy=GPT-6 Luna/low, medium=GPT-6 Luna/medium, hard=GPT-6.1 Sol/high. The dispatcher verifies the configured profile before `session.create` and the correlated terminal assistant `providerID/modelID` afterward; a requested/configured model alone is not effective-model evidence.

**Crabcode runtime difference:** Crabcode 0.0.14 documents that OpenCode `plugin` and custom-tool settings are ignored. In Crabcode, use its native `Task` tool with `subagent_type` set to the selected static profile. The RTS plugin/`orchestration_dispatch` rules below apply only when this OpenCode plugin is actually loaded; do not claim Crabcode enforces the plugin's one-shot reservation.

### Attacker review specialist

Use canonical `orchestration_dispatch(role="attacker", review_depth="low|medium|hard", review_target="...", review_revision="...", review_lens="...", handoff=..., scope_files=[...], write_files=[])` from the **loaded OpenCode plugin**. Explicit nonempty target/revision/lens and a nonempty file scope are mandatory. Review depth is separate from CODE difficulty:

| Depth | Static profile | Configured model / effort | Purpose |
| --- | --- | --- | --- |
| low (default) | `attacker-low` | `openai/gpt-6-luna` / low | Classical defensive review |
| medium | `attacker-medium` | `openai/gpt-6.1-sol` / high | Deeper path and assumption analysis |
| hard | `attacker-hard` | `openai/gpt-6-astra` / provider default | Conditional one-shot reasoning escalation only |

`attacker` and the existing `blackhat` alias map to low. Crabcode uses native `Task(subagent_type="attacker-low|attacker-medium|attacker-hard")`; its parent must verify the actual returned model. Shared contract: [`blackhat.md`](../agents/blackhat.md); [`blackhat.toml`](../agents/blackhat.toml) stays read-only. No tier fixes product code or delegates children; optional bounded verification follows the gates below and the actual client's stricter runtime limits. OpenCode runtime/config permissions are unchanged: do not infer Crabcode's additional tools are available in OpenCode. Configuration is not live model-availability evidence. The legacy `tools/orchestration.ts` discovery adapter does not expose review identity/depth arguments: use the plugin tool directly for Attacker, not that adapter.

Hard review additionally requires `escalation_of=<source task ID>` and `review_approved=true`. The source must be a **completed same-parent/root attacker-medium actual Sol** result, loaded and correlated, with exact first line `RTS_STATUS: ESCALATE`, followed by `RTS_ESCALATION: ` and a JSON object containing exactly:

```json
{"cause":"reasoning_insufficiency","bounded_assessment_completed":true,"further_reasoning_needed":true,"blocked_by":[],"unresolved_analysis":"Concrete reasoning still unresolved after the bounded assessment"}
```

The parent inspects the whole report and explicitly approves the **same exact target, revision, lens and scope**. Operational, scope/evidence, model, auth, policy, approval, credential and permission blockers are `FAILED`, never reasoning insufficiency. OpenCode verifies source profile/result/model/ancestry/identity/scope and file-byte digest, consumes a one-shot reservation (including failed Attacker launches), checks the configured Astra route before create and the correlated effective model after completion, and forbids Astra children/root dispatch. Never invoke Astra directly. Crabcode parents must retain equivalent evidence, source ID, approval and one-shot records; **no plugin/runtime one-shot guarantee is claimed for Crabcode**. `worker-astra` CODE retry rules are unchanged.

ETL requires bounded review of requirements and detailed-design proposals **before finalization/approval** when those artifacts are in scope (including non-security proposals). Security-sensitive implementation also requires review of the **approved design before CODE** and the **integrated diff before the final independent Verifier**. Material revisions invalidate the affected assessment. Non-security tasks without proposals do not need this role. Attacker does not replace the Verifier: the latter must be distinct from both the implementation author and the attacker assessor. All review failures share the existing three-correction-cycle cap; runtime concurrency and CODE routing are unchanged. These are ETL/parent obligations, not an automatic scheduler guarantee. If delegation is unavailable or not permitted, record the missing independent assessment and label any main-agent check separately; never claim attacker ran.

The reviewer critiques completeness, contradictions, abuse cases, trust boundaries, failure paths, security assumptions and unintended effects without taking over design authorship. Output identifies the actual assessed scope/revision and gives each finding's **severity, affected requirement/design/code surface, evidence, assumptions/preconditions, impact, recommended smallest defensive correction and confidence**. `NO_FINDINGS` is valid only after a bounded assessment was actually completed, not as a universal safety claim. Unresolved scope/evidence or operational blockers yield `FAILED` with `INCOMPLETE`; only completed bounded reasoning insufficiency may yield structured `ESCALATE` with `INCOMPLETE`. This permission change concerns Crabcode only: OpenCode configuration and sandbox remain read-only, with no bash-based reproduction. Bounded reproduction is allowed only after the annotation, actual Researcher CVE triage, and explicit post-triage parent proceed gate for the exact target/revision/actions/effects; no exploit or zero-day success is promised. The parent owns corrections and reassessment before advancing the relevant gate. Configuration/unit checks do not constitute a live attacker invocation.

### Annotated forensic verification (Crabcode permissions differ)

**Discovery → annotation → Researcher triage → parent confirms/no-go → only then bounded repro → Attacker report → separate Verifier.** Discovery and proposal/design review remain read-only. Attacker reports, never fixes; the existing proposal/design gates, independent parallel lenses, wait-all behavior, Luna/Sol/Astra rules, CODE escalation, shared 20-descendant/direct-child caps and three-correction-cycle limit still apply. Astra remains read-only and leaf-only.

Crabcode's four Attacker profiles (`attacker-low`, `attacker-medium`, `attacker-hard`, alias `attacker`) allow `read`, `grep`, `list`, `glob`, `external_directory`, and **unrestricted `bash`** for authorized tests/static tools and bounded forensic verification. Direct `edit`, `write`, `write_files`, `apply_patch`, `task`, `terminal_session`, `bash_output`, `bash_kill` and `bash_restart` remain denied. **This is not a hard sandbox:** unrestricted shell can mutate files and use the network; LRF behavioral restrictions still apply and shell must not bypass them. Read/grep broad source allows come first, credential/sensitive-file denies next, and only `.env.example` exceptions last, per last-match-wins semantics. These exclusions apply behaviorally to shell too.

Before any reproduction, attach a **numbered preflight annotation**: finding ID; exact target/revision; evidence/source lines; suspected weakness/preconditions; expected signal; impact/blast radius; data classification; exact commands/actions; containment/isolation; cleanup/rollback; writes/network/credential inventory; applicable current-task authorization. Prefer a local disposable fixture/checkout. Test writes are limited to explicitly approved disposable verification artifacts, never product fixes. Any nonlocal/remote/third-party test needs exact current-task target **and effect approval**; never scan beyond it.

The parent asks the read-only **Researcher** for `CVE_CANDIDATE | NOT_CVE_CANDIDATE | INSUFFICIENT_EVIDENCE`, known-CVE/advisory dedup, affected versions, severity rationale, eligibility/recordability evidence, missing info and cited public sources. Crabcode registers `agent.researcher` with `{file:../agents/researcher.md}` and no model pin; local read/list/glob/grep plus public websearch/webfetch when available, no shell/mutation/child Task. For this request it returns CVE triage only—no reproduction, file persistence, CVE submission or CNA/vendor contact. Task delegation depends on actual runtime support; Attacker cannot dispatch Researcher itself (`task` denied), and the parent must return the actual correlated result. Missing Researcher/runtime support blocks reproduction; a main-only substitute is not this required gate.

No reproduction until the report includes **annotation AND actual Researcher result AND explicit parent proceed confirmation after both**, binding exact target/revision/commands/effects. Changed scope requires re-gating. `INSUFFICIENT_EVIDENCE` means stop/no reproduction. `NOT_CVE_CANDIDATE` means never market/report it as CVE; only bounded defensive reproduction on an explicitly named owned/authorized target after parent confirmation. `CVE_CANDIDATE` permits verification only after parent confirmation, not a CVE/zero-day claim. Always prohibit credentials/secrets access or exfiltration, stealth/persistence, destructive data changes, unauthorized escalation and uncontrolled propagation. Record actual commands, exit codes, signals, cleanup and residuals before the separate Verifier. Public read-only CVE review has no public write effect; **public CVE submission/disclosure/CNA contact requires separate explicit approval** and a separate workflow, never just triage or proceed approval. No exploits run as part of this configuration task; failed CLI/tool attempts are not successful scans or product-vulnerability evidence.

For CODE, `escalation_of` is accepted only with `role=worker,difficulty=hard` and only for a completed same-parent/root/source-role Sol job with verified effective model, loaded terminal result, and exact `RTS_STATUS: FAILED` or `RTS_STATUS: ESCALATE`. The parent must inspect the result and partial diff and explicitly request one same-scope Astra retry. The dispatcher reserves at most one escalation per source result. Astra uses provider-default reasoning, is a hidden subagent where supported, is never selected as the root/main agent, and cannot dispatch children. There is no automatic upgrade and no Astra route for AGY failures. Missing or ambiguous profile/parent-message identity fails closed before child session creation.

The handoff keeps its exact machine envelope, header, ordered field labels and terminal marker. Encode only field **values** as POLYGLYPH: choose English/Japanese/Chinese/Spanish vocabulary per concept, math for logical connections and context-defined emoji anchors, never decoration. Maximize recoverable SemanticDensity subject to Meaning ∧ Constraints ∧ DecisionBoundary; safety, authority, effects, negation, acceptance and source ID/revision/path/model/status never become emoji-only. PG is representation, not SIGMA/1 or authority. This illustrative handoff is not live evidence or permission:

```text
<Γ>
SIGMA/1: PIDGEN/text
Role: attacker; review_depth=medium; 🔎=read-only assessment
Target: docs/proposal.md revision=rev-7; lens=trust-boundary; 中文: 边界假设を検査
Acceptance: 🔎 evidence-backed findings∨NO_FINDINGS after actual bounded assessment; preservar severity∧source∧preconditions∧impact∧confidence; blockers⇒FAILED
Scope: only docs/proposal.md; static-only read-only review; no edits/reproduction/external actions/children. Any reproduction follows the separate annotation + Researcher triage + explicit parent proceed gate; OpenCode remains read-only with no bash repro.
Known facts: source task ID=review-17; revision=rev-7; path=docs/proposal.md; model=openai/gpt-6.1-sol; status=proposal-ready
SIGMA/1
</Γ>
```

At each required review stage use at least two genuinely independent, artifact-supported lenses in parallel when useful, otherwise one for a single independent question. Give reviewers the same exact revision/scope and distinct lenses, wait for all, correlate findings and synthesize before advancing. Reuse only exact unchanged revision/scope/lens assessments. Runtime limits remain 20 active descendants (including grandchildren), 20 direct root children and the shared three correction cycles.

While OpenCode reviewers run asynchronously, the parent does useful nonconflicting preparation (checklists, test planning, evidence organization), not gated proposal finalization/approval, CODE or final Verifier launch. Use only supported queued notifications/await mechanisms, never busy-poll or fabricate wakeups. Crabcode native Task may block synchronously: prepare before waiting and report that limitation rather than claiming parallel background parent work.

Reservation returns `task_id` immediately. Local preflight and receipt hash verification, live `app.agents` resolution, canonical scope checks and SDK ancestry verification happen before asynchronous submission. `session.create` uses the root SDK `{body:{parentID,title}}`, never v2 session permission arguments. The assigned `messageID` is correlated to assistant `parentID`. Only a completed final `finish: stop` or completed assistant error terminates the job; idle and intermediate tool steps do not. Events plus bounded periodic polling reconcile reordered/duplicate events and expose permission waits. Results are bounded text evidence, not an automatic declaration that acceptance passed.

`orchestration_status` shows owned state, queued notifications and observer errors. `orchestration_result` retrieves correlated final evidence without waiting for model completion. Await **all** requested final evidence before synthesis while doing only non-conflicting parent work. `orchestration_cancel` cancels the entire task subtree. Abort failure/uncertain stop keeps leases and slots; uncertain prompt acceptance is never automatically re-prompted. An ambiguous session-create response with no ID remains blocked until runtime disposal/manual operator recovery; never bypass it by resubmitting.

Root depth 0 → child 1 → grandchild 2. Root orchestrator can select explicit allowlisted roles. Only researcher/proposer children may delegate read-only grandchildren within inherited exact file scope; worker/verifier are leaves. Twenty active descendants total per Runtime (exact server identity + canonical worktree), including grandchildren; at most twenty direct root children. Root sessions in one Runtime share first-come admission: reservations consume slots before launch and saturation is rejected immediately. There is no queue, fair-share, or preemption; one root may consume the full budget. Distinct Runtime instances have independent budgets. This is an admission ceiling, not a guarantee of live runtime/provider capacity. Terminal jobs do not consume slots. Host-aware canonical realpath leases (case-folded on Windows; exact case on case-sensitive hosts) reject overlapping writes and root edits of leased files.

## Safety and limits

Managed sessions are enforced at `tool.execute.before`, with immutable exact scopes and all files of supported edit/write/patch/move/delete operations checked. Unsupported shapes, native task, opaque processes/external tools and unknown worker sessions fail closed. Reads are restricted to explicit files (`read` and file-scoped `grep`); broad glob/list operations are intentionally denied. Parent handles builds/tests. Root opaque operations are denied while worker leases are active. The existing agy tool stays opt-in: its actual tool boundary calls `context.ask` using an ask-only current-task permission, not a user boolean. The Codex worker remains read-only advisory, not CODE. Existing auth modes are unchanged.

**Queued wake mode:** completion notifications are deduplicated, available in status tools, then injected on the next natural parent model turn. No timer/event prompts or promptAsync calls to parents; noReply is not a safe wake admission proxy. No toast is required.

This is **not sandbox-grade protection**. A residual filesystem check-to-write race exists between the before-hook realpath check and native execution; hard-link aliases and hostile filesystem mutation are not a sandbox boundary. In-memory state has **no restart recovery/persistence**; after restart unknown descendant writes are denied, but historical job evidence and lease recovery are unavailable. SDK mock tests do not prove live GUI event timing, plugin loading, permission prompts, agent/model availability or nested model behavior. Live runtime verification must be done separately; do not claim it from mocks. Errors are surfaced rather than silently swallowed; disposal clears timers and attempts subtree abort, retaining uncertain reservations in the process registry.
