[SPR/XML::ρ→max|target:README]
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
Critic: critic.md/toml | trigger:decision-challenge | axiom:dispute-resolution
Blackhat: blackhat.md/toml | trigger:security-assessment | axiom:BlackHatOnly
Compressor: compressor.md/toml | trigger:compress-instruction | axiom:semantics-preserved
Orchestrator: orchestrator.md/toml | trigger:coordinate-ETL | axiom:delegate-compressed
Worker: worker.md | CODE-contract | axiom:exact-allowlist
worker-agy: worker-agy.md | agy.exe-CLI | complex/tool-assisted
worker-llama: worker-llama.md | llama2coder-CLI | pure-code/single-file
</roles>

<lookup>
Routine file lookup⇒FileReader when delegation/runtime permits. Hits emit only (repository-relative-path,line_number) tuples with no status prefix; completed zero-hit search emits exactly NO_MATCH; operational/search/permission/IO error emits sanitized FAILED. No excerpt, snippet, match text, content, summary, or absolute path. Explorer remains stuck/sparse/repetitive-only and broad high-recall.
</lookup>

<css>
AGENT_ROUTER.md⇒1 CSS (stance⊗diff) except Compressor, Orchestrator, FileReader fixed mappings; CSS=presentation/cascade; ¬compose(SUBAGENTS.css, main); CSS≠role
FileReader persona: FILEREADER.css persona=🌼 ⊕ お嬢様; presentation-only; other roles dynamically select CSS.
navigate: main⇒$env:USERPROFILE\.harness\persona\PERSONA.css
</css>
⚠️ navigation-only; ¬policy
