# Walkthrough: AGY Router Skill Refresh

- task_id: `20261025-agy-subagent-router-skill`
- timestamp_utc: `2026-10-25T06:32:20.951707+00:00`
- target_environment: Windows PowerShell; local Harness root
- evidence/source_scope: AGY skill, current Harness routing/effect rules, installed CLI help/model output, target preservation fixture
- verification_status: Skill validator PASS; target JSON/registration/signature/links PASS; independent Verifier PASS_WITH_ADVISORIES; full Polyglyph CLI FAIL on manifest/inventory mismatch
- remaining_uncertainty: `.harness` has no Git metadata; child effective models resolution_pending; full inventory mismatch unresolved; AGY was not delegated to.

## Overview

Refined `skills/agy-subagent-router/SKILL.md` without changing its frontmatter or routed LRF policy. The user explicitly approved updating its existing Polyglyph preservation signature fixture after code examples changed.

## Deliverables

- `C:/Users/letwir/.harness/skills/agy-subagent-router/SKILL.md`
- The existing target record in `C:/Users/letwir/.harness/scripts/tests/fixtures/polyglyph-preservation.json`
- Required Harness evaluation record: `C:/Users/letwir/.harness/evaluation/subagents_model_eval.json` (12 child invocations; effective model/reasoning marked resolution_pending)

## 実行計画

1. Read the current Harness authority, AGY, documentation, and delegation contracts plus the existing skill.
2. Verify installed CLI option names and read-only current model inventory format; do not delegate AGY.
3. Have the plan independently challenged; refine the plan to keep the AGY flow explicitly authorized, bounded, and parent-verified.
4. Update only the skill body; preserve name, description, invocation default, and existing LRF policy.
5. With user approval, refresh only the target’s frozen preservation signature.
6. Run skill validator, Polyglyph structural/signature/link checks, full linter, and independent final review; report any global lint limitation.
7. Record subagent runs and task-close evidence with all unavailable model metadata marked pending.

## Verification

- `quick_validate.py` on `skills/agy-subagent-router`: PASS, “Skill is valid!”.
- Target parse and preservation assertion: PASS; the current target is registered and its frontmatter, three fenced blocks, and LRF signature match the fixture.
- JSON parse of preservation fixture: PASS.
- Relative links: PASS, all 7 resolve.
- Installed `agy.exe --help`: exit 0; confirmed current options, including effort values `low|medium|high|max`, modes, sandbox, add-dir, print and JSON output.
- Read-only `agy.exe models`: exit 0, 13 tab-separated ID/display-name rows; no model was selected or used.
- Independent Verifier: PASS_WITH_ADVISORIES; independently confirmed the target-specific checks and source alignment.
- Full `lint-polyglyph.py --targets skills/agy-subagent-router/SKILL.md`: FAIL, exit 1, `manifest/inventory mismatch`. Normalized set difference: fixture-only `skills/.system/plugin-creator/SKILL.md`; inventory-only `agents/CALL.md`. These unrelated entries were not changed.
- `.harness` is not a Git repository; no Git diff or commit evidence exists.

## Residual risk

The entire Harness inventory gate remains blocked by the above unrelated mismatch. No AGY delegation was run, so no model was selected and no runtime behavior of the example commands was executed. Subagent effective model and reasoning metadata were unavailable and remain resolution_pending.
