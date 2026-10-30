# Walkthrough: 160 subagent capacity

1. Read pi-subagents configuration and nested-agent documentation; confirmed global/per-run concurrency, session/run spawn budgets, and depth controls.
2. Raised `~/.pi/agent/extensions/subagent/config.json` limits to 160 and kept depth at 2.
3. Added `~/.pi/agent/agents/orchestrator.md` with child-safe nested orchestration, a 40-launch per-coordinator policy, and no Astra route. Connected `~/.pi/agent/POLYGLYPH.md` to the Harness codec after the first child stopped during bootstrap.
4. Raised applicable Harness routing/workflow caps to 160; documented up to four root Orchestrators and 40 launches each, bounded by the shared 160 total. Left the Astra escalation route unchanged.
5. Ran Harness preflight and verified its receipt (PASS); parsed pi settings/subagent JSON and verified Harness symlinks (PASS). Confirmed pi had already loaded the old 64 per-run cap in this process; restart is required to apply 160.
6. Did not launch a 160-child workload. An Orchestrator audit returned text but the runtime rejected the result for missing command-run evidence; no accepted independent review is claimed.
