[SPR/XML::ρ→max|target:AGENT_ROUTER|legibility:LLM≫human]
# Agent Router
<router auth="rules/BOOTSTRAP.lrf, rules/MANUAL.lrf, agents/SUBAGENTS.lrf">
<routes>
FileReader: ordinary repository-file lookup⇒bounded_read-only|codex_luna/terra|filereader.md+filereader.toml|FILEREADER.css
Explorer: stuck∨sparse/repetitive⇒broad_high-recall|codex_luna/terra|explorer.md+explorer.toml
Researcher: docs/methods/usage⇒scoped|codex_terra/sol|researcher.md+researcher.toml
Proposer: ambiguity/plan⇒bounded_plan|codex_luna/sol|proposer.md+proposer.toml
Auditor: high-risk/disagreement⇒challenge|codex_terra/sol|auditor.md+auditor.toml
Verifier: evidence/check⇒independent|codex_terra/sol|verifier.md+verifier.toml
Refactorer: verifier_reject⇒repair_plan|codex_luna/sol|refactorer.md+refactorer.toml
Critic: dispute⇒converge|codex_terra/sol|critic.md+critic.toml
Blackhat: security⇒assessment|codex_terra/sol|blackhat.md+blackhat.toml
Compressor: human_plan_settled∧¬unresolved_human_decision⇒compress_immediately_before_execution|codex_sol_high|COMPRESSOR.css|compressor.md+compressor.toml
Orchestrator: ETL⇒coordinate+delegate_compressed|codex_luna_high|ORCHESTRATOR.css|orchestrator.md+orchestrator.toml
worker-gpt: eligible CODE⇒Codex_GPT_family(accept-edits)|current_registry+context_evidence|worker.md
worker-llama: eligible single-file CODE⇒llama2coder(stdout)|check_local_model+context_evidence|worker.md+worker-llama.md
worker-agy: separate_explicit_current_task_request⇒agy.exe(mode:plan|accept-edits)|check_agy.exe_models|worker.md+worker-agy.md
worker-gpt-fallback: legacy agy.failure compatibility⇒Codex_GPT_family(accept-edits)|serial_same_scope+current_registry|worker.md+worker-gpt-fallback.md
</routes>

<constraints>
Hierarchy: root owns dispatch; root(0)→child(1)→grandchild(2); ¬great-grandchild; ¬codex⇒root_local
Diff: uncertainty⊕scope⊕risk⊕verify_burden ↦ {easy,med,hard} ↦ model
CSS: 1 persona_CSS/invocation (stance⊗difficulty) ∨ {COMPRESSOR|ORCHESTRATOR|FILEREADER}; FileReader fixed=FILEREADER.css(persona=🌼 ⊕ お嬢様); all other role CSS dynamic; ¬main_persona; ¬SUBAGENTS.css; CSS=presentation-only
FileReader: ordinary_file_lookup⇒FileReader when delegation/runtime permits; hits⇒only (repository-relative-path,line_number) tuples with no status prefix ∧ ¬excerpt∧¬snippet∧¬match_text∧¬content∧¬summary∧¬absolute_path; zero_hits_after_completed_search⇒exactly NO_MATCH; operational/search/permission/IO error⇒sanitized FAILED. Stuck/sparse/repetitive investigation⇒Explorer broad high-recall; preserve ¬topic/source precision filters.
Backend: ≤1 active worker/workstream; select among eligible currently available workers by smallest comparable measured context consumption for the bounded task (same accounting, workload and scope); context window capacity ≠ consumed context; if no comparable measurements use evaluation/worker-context-priority.json dated nonnormative user heuristic in listed order, skipping unavailable/ineligible models; no invented figures; update measurements only from attributable official comparable sources when found; agy opt-in only, never default; llama=single-file_stdout; operational failure⇒stop_or_confirm_exit+inspect_partial_diff⇒serial_same_scope_reselection; Astra⇒leaf_no_subagent_dispatch; ¬fallback(auth|policy|approval|credential|permission)
Isolation: Researcher/Explorer ¬CODE; Worker ¬external/OOB_write
Handoff: rules/MANUAL.lrf envelope; scoped_facts_only
Concurrency: ≤6 independent children concurrently when runtime slots permit; unlimited total invocations; 1_writer/file; ordered dependencies remain serial
</constraints>

<backend_bounds>
worker-agy: separate_explicit_current_task_request_only; agy.exe_models_only; accept-edits=explicit_scope; 1_retry_pre-launch; any operational FAILED⇒stop_or_confirm_exit+inspect_partial_diff⇒serial_same_scope_reselection; ¬auth_bypass; ¬policy_bypass; non_recursive
worker-gpt: current_registry_only; exact_scope; Astra=leaf_no_subagent_dispatch
worker-gpt-fallback: legacy compatibility after explicitly requested agy operational failure only; current_registry_only; same_scope; no concurrent agy; Astra=leaf_no_subagent_dispatch
worker-llama: stdout_pure_code; main_applies_diff; ¬direct_pipe
Both: return(files, checks, fails, risk); main_owns_verification
</backend_bounds>
</router>
