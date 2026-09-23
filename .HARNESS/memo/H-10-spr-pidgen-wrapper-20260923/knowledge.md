# Knowledge — H-10 SPR/PIDGEN subagent wrapper

- The invocation envelope is `<Γ>`, followed by `SIGMA/1: PIDGEN/text`, a task-specific payload, and `</Γ>`.
- Codex subagent handoffs retain the ordered fields Role, Target, Acceptance, Scope, Known facts; known facts must be verified.
- The header identifies the text protocol and loaded policy. It grants no authority and does not replace bootstrap plus routed LRF loading.
- The contract is defined in `rules/MANUAL.lrf`, `agents/SUBAGENTS.lrf`, and `rules/agy.lrf`; the agy guide and all nine role TOMLs reflect it.
- Validation: MANUAL and agy LRF Level 2 pass; SUBAGENTS LRF Level 1 passes; all nine TOMLs parse; `git diff --check` passes. SUBAGENTS LRF Level 2 still reports an effect mismatch at its pre-existing `subagent.mode` record, line 51.