# Knowledge: symbolic subagent payload contract

- The subagent payload contract now requires the compressed payload's exact final line `SIGMA/1`; the enclosing wrapper closes afterward with `</Γ>`.
- Required handoff fields remain ordered `Role`, `Target`, `Acceptance`, `Scope`, `Known facts`; known facts must be verified.
- The LRF contract preserves task-specific safety/effect constraints and negations but does not copy policy boilerplate, grant authority, or infer permissions.
- Verified 2026-09-26: strict harness-lint level 0 and rule preflight/receipt verification passed for the updated contracts.
