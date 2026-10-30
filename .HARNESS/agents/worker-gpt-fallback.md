@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

Representation: POLYGLYPH.md; ρ→max; target=WORKER_GPT_FALLBACK.
# worker-gpt-fallback protocol
Route: AGENT_ROUTER.md⇒FAILED(agy operational)⇒Codex GPT-family(bounded_CODE).
Rules:
- Entry⇒agy_operational_FAILURE∧stopped∨confirmed_exit∧partial_diff_inspected.
- agy operational fallback worker profiles: easy=worker-easy(openai/gpt-6-luna,low); medium=worker-medium(openai/gpt-6-luna,medium); hard=worker-hard(openai/gpt-6.1-sol,high); verify actual terminal model.
- Astra is NEVER selected for agy failure. Only explicit `escalation_of` may select worker-astra after a completed same-scope hard-Sol result is verified and marked FAILED/ESCALATE and the parent inspected the result/diff. Astra is subagent-only, leaf-only, never main.
- Continue⇒same_bounded(acceptance,file_scope); preserve(existing_edits∪partial_edits).
- ¬overlap(agy); active_writers_per_file≤1.
- block∈{auth,policy,approval,credential,permission}⇒terminal∧¬route.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
