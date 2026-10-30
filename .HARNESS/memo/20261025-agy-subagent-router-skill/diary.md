### 2026-10-25 — AGY router skill refresh

- task_id: `20261025-agy-subagent-router-skill`
- timestamp_utc: `2026-10-25T06:32:20.951707+00:00`
- target_environment: Windows PowerShell; local Harness root
- evidence/source_scope: AGY skill and governing Harness LRF/router sources; installed CLI help and read-only model inventory; target validator, parser, fixture, relative links, independent review
- verification_status: Target skill and preservation fixture checks PASS; independent Verifier PASS_WITH_ADVISORIES; full Polyglyph CLI FAIL due manifest/inventory mismatch
- remaining_uncertainty: No Git diff available for `.harness`; effective child models remain resolution_pending; no AGY model was selected and no AGY delegation was executed.

#### Hypothesis
The existing AGY router skill could be made clearer and executable across supported host profiles while preserving its opt-in trigger and the global effect/permission policy.

#### Tried
Compared the skill with `agy.lrf`, `AGENT_ROUTER.md`, `worker-agy.md`, `SUBAGENTS.lrf`, `HARDWARE.md`, and `documentation.lrf`. Read installed CLI help and ran the read-only `models` inventory to confirm its current tab-separated ID/name row shape. Revised the skill body, kept frontmatter and LRF authority block intact, and updated the user-approved preservation fixture record.

#### Rejected
Did not invoke AGY for this task, select a model, expand permissions, modify routing policies, or update unrelated inventory entries. Did not treat JSON output or a child success message as proof of execution success.

#### Uncertainty
The full Polyglyph CLI exits 1 at `manifest/inventory mismatch`; the normalized inventory difference is fixture-only `skills/.system/plugin-creator/SKILL.md` versus inventory-only `agents/CALL.md`. The target-specific signature, JSON, and links still pass. `.harness` is not a Git repository, and subagent effective model metadata was not returned.

#### Attribution
- PromptDefect: none identified; the user’s scope was clear. The preservation-fixture extension was explicitly requested and approved after the initial target-only repair proved insufficient for frozen code-fence hashes.
- AgentDefect: several worker iterations were needed to align command resolution, output display, model-row validation, and fixture JSON punctuation; parent-side parsing caught and corrected the missing separator. One verifier’s first hash invocation had a quoting failure; a later independent review completed the checks. All effective subagent model fields are recorded as `resolution_pending` rather than inferred.

#### Search
Local Harness sources and installed CLI only. `agy.exe --help` and `agy.exe models` were read-only inspections; no external delegation, network write, database access, credential handling, or VCS operation occurred.

#### Correction
Clarified explicit opt-in; corrected current CLI effort/mode/output flags; added the exact handoff ending; made examples host-profile-aware; added a parseable model inventory gate and manual exact-ID selection; aligned sandbox/auth/failure/fallback/final verification guidance to the routed LRF; added local links; refreshed the target preservation hashes.

#### Emotion
No user emotion inferred or recorded.

#### Thoughts
The target is now locally valid and source-aligned. The remaining lint failure is global inventory drift, not a target signature failure; changing unrelated inventory or fixtures would exceed the approved scope.
