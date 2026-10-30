# Knowledge Base: AGY Router Skill Refresh

- task_id: `20261025-agy-subagent-router-skill`
- timestamp_utc: `2026-10-25T06:32:20.951707+00:00`
- target_environment: Windows PowerShell profile; Harness root `C:/Users/letwir/.harness`
- evidence/source_scope: `skills/agy-subagent-router/SKILL.md`; `rules/agy.lrf`; `rules/documentation.lrf`; `agents/AGENT_ROUTER.md`; `agents/worker-agy.md`; `agents/SUBAGENTS.lrf`; `HARDWARE.md`; installed `agy.exe --help` and read-only `agy.exe models`; Polyglyph parser and fixtures
- verification_status: Skill validator PASS; target parser/frozen signature/JSON/link checks PASS; independent Verifier PASS_WITH_ADVISORIES; full Polyglyph CLI FAIL on manifest/inventory mismatch
- remaining_uncertainty: `.harness` has no Git metadata, so no Git diff/history attestation. Effective models for the 12 subagent calls are `resolution_pending`; runtime models were not present in task results. Full project lint remains blocked by an inventory discrepancy: fixture-only `skills/.system/plugin-creator/SKILL.md`, inventory-only `agents/CALL.md`.

## Context Findings

- The skill is opt-in for an explicit current-task request to route bounded CODE through AGY when `AGENT_ROUTER.md` selects `worker-agy`; invoking the skill alone does not grant task authorization or effects.
- Canonical rules live in `rules/agy.lrf`, `agents/AGENT_ROUTER.md`, `agents/worker-agy.md`, `agents/SUBAGENTS.lrf`, and `HARDWARE.md`. The skill is operational guidance, not an authority override.
- Installed `agy.exe --help` reports `--effort` values `low|medium|high|max`, mode `plan|accept-edits`, output `text|json|stream-json`, and flags `--print`, `--sandbox`, `--add-dir`. Read-only `models` output was observed as tab-separated model ID and display name rows; no model was selected and no delegation was run.
- Handoffs require `<Γ>`, `SIGMA/1: PIDGEN/text`, exact field order `Role`, `Target`, `Acceptance`, `Scope`, `Known facts`, terminal payload line `SIGMA/1`, then `</Γ>`.
- Auth/policy/approval/credential/permission blocks are terminal. AGY operational fallback must remain serial, same-scope, and parent-selected; no concurrent fallback, Astra fallback, or recursion.

## Morphism

- Input: Explicit current-task request to route a bounded CODE task through AGY, with worker-agy selected by Harness routing.
- Operation: Resolve binary by active profile → inspect current model inventory → construct exact bounded handoff → select plan/accept-edits/read-only mode and sandbox → execute only after applicable authorization → parent confirms exit, inspects actual diff, runs checks, and reconciles acceptance.
- Result evidence: updated `skills/agy-subagent-router/SKILL.md`, preserved in its Polyglyph signature fixture, with verified frontmatter, LRF parse, relative links, CLI-facing examples, stop gates, and parent-owned completion criteria.
- Guard: no skill-inferred permission, stale model IDs, secrets, dangerous permission bypass, recursive delegation, external writes, or VCS operations.
