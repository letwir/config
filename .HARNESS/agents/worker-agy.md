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
- GPT_mapping: easy=Luna; medium∈{Luna,Sol}; hard∈{Sol,Astra}; registry=current_only; Astra=highest_difficulty_eligible∧leaf_only∧¬subagent_dispatch.
- block∈{auth,policy,approval,credential,permission}⇒terminal∧¬fallback.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
