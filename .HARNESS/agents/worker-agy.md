@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/XML::ρ→max|target:WORKER_AGY]
# worker-agy protocol
Route: AGENT_ROUTER.md⇒agy.exe(bounded_CODE).
Rules:
- `agy.exe models`⇒use_returned_ID.
- mode=plan(design) | mode=accept-edits(explicit_scope).
- Restrict_workspace; ¬--dangerously-skip-permissions; ¬recursive_delegation.
- 1_retry_pre-launch; inspect_diff_before_fallback.
- agy_operational_FAILURE⇒stop∨confirm_exit∧inspect(partial_diff)⇒reselect_one(GPT-family; difficulty=current_task; evidence=current_evaluation; scope=same_bounded); ¬concurrent_fallback.
- GPT_mapping: for bounded worker CODE use easy=worker-easy(openai/gpt-6-luna,low), medium=worker-medium(openai/gpt-6-luna,medium), hard=worker-hard(openai/gpt-6.1-sol,high); verify terminal effective model. An agy operational failure NEVER selects Astra. Astra is eligible only by explicit one-shot escalation_of of a completed hard-Sol FAILED/ESCALATE result after parent review; Astra is subagent-only leaf, never main and never dispatches children.
- block∈{auth,policy,approval,credential,permission}⇒terminal∧¬fallback.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
