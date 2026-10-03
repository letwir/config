# walkthrough | 20260923-harness-agent-routing-audit

## 実行計画
1. Load and classify the active LRF bootstrap.
2. Trace main.seq through research and change stages.
3. Compare subagent role policy, role files, and Codex registration/linkage.
4. Report verified facts, inferences, unknown runtime behavior, and bounded improvements.

## 結果
- Bootstrap/read-only classification succeeded; loaded matching subagent rules.
- Confirmed `.codex/agents` -> `.harness/agents` symbolic link and matching SHA-256 for `researcher.toml`.
- Confirmed Codex agent runtime enabled, concurrency cap 4; this proves configuration, not a live dispatch.
- Found README lifecycle mismatch and missing `worker.md`; current active LRF routes CODE to agy and planning to main.
- No workspace policy or runtime configuration edits; no tests executed.
- Memory precedent lookup returned `null` with process exit 0, so prior memory evidence was not available for this query.
