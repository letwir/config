# Verified knowledge — llm-mem live validation

Verified 2026-10-25 from the current source and live CLI checks.

- `eval` input is a strict JSON contract: `axes` is a nonempty object mapping axis names to numeric scores in `[0,1]`; attribution uses `prompt_defects` and `agent_defects` arrays, with optional ratios. Unknown fields and trailing JSON values are rejected by `decodeEvaluationInput` in `llm-memory/evaluation.go`; focused tests cover the boundary.
- The Windows active executable resolved to `~/.harness/skills/llm-memory/llm-mem.exe` (the first existing candidate under the skill's resolution order). It was rebuilt offline from `C:/Users/letwir/repo/repo/llm-memory` without embedding a DB URL and installed with a matching SHA-256 check.
- Runtime DB access used the existing `LLM_MEMORY_DB_URL` environment variable; its value was neither displayed nor recorded. Gemini keys were removed from the ingest process environment, so report ingestion used the local fallback and made no Gemini API call.
- History search for `llm-mem-repair-20261025` returned four existing memory records. Before this validation, `evals -key llm-mem-repair` returned zero records.
- Survey response `1046ee13-e160-4fb6-9ec1-08fca3b09d37` was submitted and read back as `survey_observation`, form `harness-feedback/v1`, with seven answers and task ID `llm-mem-live-validation-20261025`. The respondent is explicitly the assistant, not the user.
- The active `~/.harness/skills/llm-memory/SKILL.md` describes the evaluation schema only as `axes{obj}` and `attribution{prompt,agent}`; the source `USAGE.md` provides the more specific numeric-map and field-name contract. Treat the shorthand as insufficient for generating eval JSON.

No secrets or raw environment values are included here.
