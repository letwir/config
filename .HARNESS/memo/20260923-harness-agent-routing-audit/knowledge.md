# knowledge | 20260923-harness-agent-routing-audit

## Facts
- `.codex/agents` is a directory symlink to `.harness/agents`; `researcher.toml` hashes matched at both paths. Codex config enables agents with a four-thread cap.
- `etl/main.seq` routes through PRECEDENT, then local review, local research, or RESEARCH. `research.seq` permits delegation only under agent.mode and with at least two independent workstreams or ambiguity/high impact.
- `etl/change.seq` assigns planning to main, verification to a permitted Verifier, and implementation to bounded agy. `SUBAGENTS.lrf` says Codex Worker cannot replace agy.
- `agents/README.md` describes a different legacy Worker-plan -> Auditor -> Worker-code -> Verifier -> Refactorer flow.
- `worker.toml` exists; `worker.md` does not. `WORKER.css` exists.

## Inference
- Agent definitions are discoverable by Codex but are not automatically activated by the main ETL. The explicit policy gates and runtime delegation capability decide whether a role branches off.
- The stale README and Worker entry are likely sources of role confusion. Explorer also has no explicit stage assignment in the current ETL rules.

## Unknowns
- This read-only review did not dispatch a role, so successful fresh-session invocation and role-specific instruction loading remain unverified.
- The llm-mem precedent search returned `null` with exit code 0; treat memory precedent as unavailable.
