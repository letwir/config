[SPR/XML::ρ→max|target:WORKER_AGY]
# worker-agy protocol
Route: AGENT_ROUTER.md⇒agy.exe(bounded_CODE).
Rules:
- `agy.exe models`⇒use_returned_ID.
- mode=plan(design) | mode=accept-edits(explicit_scope).
- Restrict_workspace; ¬--dangerously-skip-permissions; ¬recursive_delegation.
- 1_retry_pre-launch; inspect_diff_before_fallback.
- Any operational agy FAILED: stop or confirm process exit, inspect partial diff, then route one GPT-family fallback reselected by task difficulty and current evaluation evidence for the same bounded scope; never run fallback concurrently.
- GPT fallback mapping: easy=Luna; medium=Luna or Sol; hard=Sol or Astra; current registry only; Astra is eligible at highest difficulty and runs leaf-only without subagent dispatch.
- Terminal without fallback: auth/policy/approval/credential/permission blocks.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
