@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

Representation: POLYGLYPH.md; ρ→max; target=README.
# Agent and worker routing
AGENT_ROUTER.md: loaded_by(etl/main.seq); chooses(role, model, CSS, TOML); LRF=auth(permissions, effects, safety)

<lifecycle>
plan-audit: Proposer(draft) → Auditor(challenge) → Critic(dispute)
implement-verify: Worker(CODE) → Verifier(evidence) → Refactorer(repair)
stages: 探索(explore) → 計画(plan) → 実装(implement) → Verify → 総括(summary)
gate-outcomes: SUCCESS|FAILED|ADVICE|ESCALATE
</lifecycle>

<roles>
FileReader: filereader.md/toml | trigger:ordinary-repository-file-lookup | axiom:hits-only-(repository-relative-path,line_number); zero=NO_MATCH; errors=FAILED
Explorer: explorer.md/toml | trigger:stuck/sparse-investigation | axiom:broad-high-recall
Researcher: researcher.md/toml | trigger:scoped-research | axiom:fact-based
Proposer: proposer.md/toml | trigger:bounded-plan | axiom:clear-boundary
Auditor: auditor.md/toml | trigger:plan-audit | axiom:independent-challenge
Verifier: verifier.md/toml | trigger:change-verification | axiom:PoC/BFP/AF
Verifier-blackhat: [verifier-blackhat.md](verifier-blackhat.md)/[toml](verifier-blackhat.toml) | trigger:security-sensitive-or-explicit-security-review | axiom:read-only-independent-final-diff-first-then-Attacker-report-audit; DISCOVERY_COMPLETE before report access; RETRY_REQUIRED names missed coverage and may advise attacker-low/medium, Main alone dispatches within existing max-3; no child delegation, exploit execution or automatic Astra | flow:[security-review.seq](../etl/security-review.seq) from change.seq; profiles require actual runtime support, configuration is not a dispatch claim
Critic: critic.md/toml | trigger:decision-challenge | axiom:dispute-resolution
Blackhat: blackhat.md/toml | trigger:security-assessment | axiom:BlackHatOnly
Compressor: compressor.md/toml | trigger:compress-instruction | axiom:semantics-preserved
Orchestrator: orchestrator.md/toml | trigger:coordinate-ETL | axiom:delegate-compressed
Worker: worker.md | CODE-contract | axiom:exact-allowlist
worker-agy: worker-agy.md | agy.exe-CLI | complex/tool-assisted
worker-llama: worker-llama.md | llama2coder-CLI | pure-code/single-file
</roles>

<worker-model-routing>
CODE `role=worker` requires explicit difficulty and selects the matching static profile: easy=`worker-easy` (GPT-6 Luna/low), medium=`worker-medium` (GPT-6 Luna/medium), hard=`worker-hard` (GPT-6.1 Sol/high). Verify the child assistant's actual provider/model before claiming the route. Crabcode uses native Task with `subagent_type` set to the profile; it ignores OpenCode plugins and custom tools. OpenCode may use the RTS `difficulty`/`escalation_of` tool only when that plugin is loaded. `worker-astra` is a subagent profile for parent-requested retry only after a completed, effective hard-Sol result returns `RTS_STATUS: FAILED|ESCALATE` and the parent reviews its result/partial diff. Crabcode native Task does not enforce a one-shot token, so the parent must perform only one reviewed retry; OpenCode RTS enforces a one-shot reservation. Astra is a leaf subagent, never main, and cannot dispatch children. No automatic upgrade or Astra fallback for agy failures.
</worker-model-routing>

<lookup>
Routine file lookup⇒FileReader when delegation/runtime permits. Hits emit only (repository-relative-path,line_number) tuples with no status prefix; completed zero-hit search emits exactly NO_MATCH; operational/search/permission/IO error emits sanitized FAILED. No excerpt, snippet, match text, content, summary, or absolute path. Explorer remains stuck/sparse/repetitive-only and broad high-recall.
</lookup>

<css>
AGENT_ROUTER.md⇒1 CSS (stance⊗diff) except Compressor, Orchestrator, FileReader fixed mappings; CSS=presentation/cascade; ¬compose(SUBAGENTS.css, main); CSS≠role
FileReader persona: FILEREADER.css persona=🌼 ⊕ お嬢様; presentation-only; other roles dynamically select CSS.
navigate: main⇒$env:USERPROFILE\.harness\persona\PERSONA.css
</css>
⚠️ navigation-only; ¬policy
