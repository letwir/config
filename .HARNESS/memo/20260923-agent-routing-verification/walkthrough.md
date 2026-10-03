# walkthrough | 20260923-agent-routing-verification

## Result
- `main.seq` calls `../agents/AGENT_ROUTER.md` only when delegation is eligible; it does not select a role or model.
- Router table defines role triggers, model profiles, one flexible CSS choice per invocation, and a TOML only for Codex custom-agent routes.
- Explorer handles local structure/symbol/dependency discovery; Researcher handles current/external claims and source-backed evidence. `agent.mode` gates Codex delegation; otherwise the root works locally.
- `worker.md` defines the common PIDGEN contract. `worker-agy.md` and `worker-llama.md` define separate execution paths: bounded edit flow vs pure-code stdout reviewed/applied by main. No implicit cross-backend fallback after launch.
- TOML review removed vague proactive/parallel triggers and made Refactorer read-only. Auditor stale TOML imports removed; hard-coded CSS imports and role-to-CSS bindings removed from protocols.

## Checks
- `python.exe` with `tomllib`: 8 TOMLs parsed, all names unique, all sandbox modes read-only.
- `harness-lint.exe -path agents -level 2 -json`: no diagnostics.
- `harness-lint.exe -path etl -level 2 -json`: no diagnostics.
- `harness-lint.exe -path rules -level 2 -json`: one existing unrelated `rules/documentation.lrf:11` `L2:EFFECT_MISMATCH` (`rewrite` under `RO_LOCAL`).
- `git -c core.whitespace=cr-at-eol diff --check` on the routing scope: clean.
- AGY execution completed successfully; main inspected actual files and performed final local review.
