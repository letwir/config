@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/XML::ρ→max|target:WORKER_GPT_FALLBACK]
# worker-gpt-fallback protocol
Route: AGENT_ROUTER.md⇒FAILED(agy operational)⇒Codex GPT-family(bounded_CODE).
Rules:
- Entry⇒agy_operational_FAILURE∧stopped∨confirmed_exit∧partial_diff_inspected.
- Select∈current_Codex_registry by (difficulty,evaluation): easy=Luna; medium∈{Luna,Sol}; hard∈{Sol,Astra}.
- Astra⇒highest_difficulty_eligible∧MUST_NOT(subagent); ¬invent∨silently_retain(stale_model_identifier).
- Continue⇒same_bounded(acceptance,file_scope); preserve(existing_edits∪partial_edits).
- ¬overlap(agy); active_writers_per_file≤1.
- block∈{auth,policy,approval,credential,permission}⇒terminal∧¬route.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
