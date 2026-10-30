# Diary: Harness/pi parallel-capacity change

- Timestamp: 2026-10-25 (workspace date convention)
- Task: Raise Harness/pi-subagents concurrency configuration to 160 total with nested Orchestrator policy of 40 launches each; retain Astra restrictions.
- Request evidence: User explicitly asked for total 160, child Orchestrators with up to 40 launches each, and unchanged Astra limits.
- Action/result: Updated pi-subagents global/session/run/parallel settings to 160, added a pi Orchestrator agent profile, raised Harness routing limits to 160, set four root Orchestrators and a 40-launch policy cap per coordinator, and connected the missing Pi POLYGLYPH codec link. Preflight/receipt and JSON/link checks passed.
- Friction: One Orchestrator child stopped before audit on the missing codec; retry returned findings but was rejected by the runtime acceptance gate for missing command-run evidence. No child result is claimed as a successful accepted review. Runtime numeric per-agent 40 enforcement and cross-session global enforcement are not provided by the installed pi-subagents settings.
- Attribution/impact: Configuration/policy changes only; no 160-agent workload was launched. Runtime settings require Pi restart.
- Prompt/agent defect estimate: Prompt 0% (desired ceilings and Astra constraint were explicit); agent/tooling friction was missing codec link and acceptance contract mismatch.
- Residual: Aggregate across concurrent root Pi sessions is policy-only, not enforced by one shared runtime quota.
