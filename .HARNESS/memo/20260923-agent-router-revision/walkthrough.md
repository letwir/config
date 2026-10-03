# walkthrough | 20260923-agent-router-revision

## 実行計画
Incorporate the user's router-file preference into the current design while preserving the LRF authority and direct-agy CODE boundaries.

## 更新設計
- `agents/AGENT_ROUTER.md` is the single selection manifest: trigger, role, model/reasoning profile, base CSS, TOML, protocol, scope, and exclusions.
- `etl/main.seq` retains overall ETL stages and reduces its agent-specific behavior to invoking this router.
- `agents/SUBAGENTS.lrf` defines how the Markdown manifest is validated and consumed as configuration data; LRF remains authoritative for effects, safety, authorization, and stop behavior.
- Existing `SUBAGENTS.lrf` role-choice rules should be removed or reduced to execution invariants to avoid duplicate selectors.
- Keep model selection dynamic; do not pin models in role TOMLs. TOML descriptions contain concrete triggers, exclusions, and purpose, never `proactively`.
- CODE remains direct agy under current policy. Before implementation, decide whether `worker.toml` is retired/disabled or repurposed as a read-only task-packet role; it cannot remain a Codex code-writing route in parallel.
- No operational configuration or policy files changed in this design turn.
