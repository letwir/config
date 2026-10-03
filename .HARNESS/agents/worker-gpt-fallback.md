[SPR/XML::ρ→max|target:WORKER_GPT_FALLBACK]
# worker-gpt-fallback protocol
Route: AGENT_ROUTER.md⇒FAILED(agy operational)⇒Codex GPT-family(bounded_CODE).
Rules:
- Start only after the agy process is stopped or confirmed exited and its partial diff is inspected.
- Select from the current Codex registry by task difficulty and evaluation evidence: easy=Luna; medium=Luna or Sol; hard=Sol or Astra.
- Astra is eligible at the highest task difficulty on this fallback route and MUST_NOT invoke any subagent; never invent or silently retain a stale model identifier.
- Continue only the same bounded acceptance and file scope; preserve existing and partial edits.
- Never overlap with agy; one active writer per file.
- Auth/policy/approval/credential/permission blocks are terminal and must not reach this route.
- Return(changed_files, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬alter(handoff, acceptance, authority, effects).
