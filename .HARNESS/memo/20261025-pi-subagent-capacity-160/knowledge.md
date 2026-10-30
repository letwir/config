# Knowledge: pi-subagents 160-capacity configuration

## Verified
- `~/.pi/agent/extensions/subagent/config.json` now sets `globalConcurrencyLimit`, `maxSubagentSpawnsPerSession`, `maxSubagentSpawnsPerRun`, and `parallel.maxTasks/concurrency` to 160; `maxSubagentDepth` remains 2.
- Harness routing/workflow policy now permits a shared maximum of 160 descendants per Runtime and up to four root Orchestrator children, each with a 40-launch policy cap. Coordinators count in the shared total, so the four per-child maxima are not additive.
- `~/.pi/agent/agents/orchestrator.md` defines a nested Orchestrator profile. pi-subagents has no separate per-agent numeric hard limit; the 40 limit is enforced by the role/policy and per-workflow arguments, while the configured 160 ceilings are runtime settings.
- Astra routing in Harness was not changed. The Orchestrator profile does not allow Astra as a child agent.
- A first child attempt stopped before audit because `~/.pi/agent/POLYGLYPH.md` was missing. The link was added. A second child returned an audit, but pi-subagents rejected its result because its acceptance contract required command-run evidence. Its returned text was treated as a lead and independently checked; it is not recorded as a successful independent review.
- Harness LRF preflight and receipt verification passed after the policy changes; JSON and symlink checks passed.

## Limitation
pi-subagents' `globalConcurrencyLimit` is per run, and its spawn budget is per session/run, not an OS/server-wide cross-session arbiter. Multiple independent Pi root sessions may exceed 160 in aggregate. Restart Pi for updated runtime settings to load.
