# Task report — llm-mem live validation

## Result

- Deployed the strict evaluation-input decoder to the canonical Windows executable `~/.harness/skills/llm-memory/llm-mem.exe`; staged and installed SHA-256 matched.
- The repaired CLI rejected malformed nested axes and unknown attribution fields before attempting DB access.
- Submitted one assistant-authored `survey_observation` to `harness-feedback/v1`, then read it back successfully (7 answers; response ID `1046ee13-e160-4fb6-9ec1-08fca3b09d37`).
- Read-only history search for `llm-mem-repair-20261025` returned 4 records. Before adding this validation's evaluation, `evals -key llm-mem-repair` returned 0 records.
- Runtime DB access used the configured environment variable without displaying its value. Ingested each task Markdown once with Gemini keys removed from the child process: knowledge `580c789a-bc5c-458e-b074-8c935d946667`, diary `fc93f28a-e98d-43de-9915-d434531b76b6`, walkthrough `0f43ed3b-cb0a-46c2-9f8c-6031f231d81a`. No Gemini API call was made.
- Submitted the task evaluation once; receipt `8c9b0a4b-cd09-45b3-8b69-b5b1695e7e7f`. Readback returned one evaluation with PromptDefect 25% / AgentDefect 75%. Current task search returned three records; the pre-existing history search returned four records.
- `analyze -file diary.md -json -suggest=false` displayed one task entry and average PromptDefect 25% / AgentDefect 75%. Its keyword-frequency counters are separate from the explicit ratios (Prompt=1, Agent=0 in this diary), so use the evaluation record for the task-scoped attribution counts.

## Task-scoped Defect attribution

Self-assessed current task-history ratio: **PromptDefect 25% / AgentDefect 75%**. The prompt-side item is the active skill's underspecified eval shorthand. The agent-side items are the prior malformed axes input, receipt parsing mismatch, and initial harness-root misresolution. This is not a global defect rate.

## Verification

- Focused offline Go tests passed.
- `go build ./...` passed offline.
- `git diff --check` passed.
- Strict-schema executable smoke, installed binary hash check, survey template, DB history search, survey submit, and survey readback passed.

## Boundaries

No raw environment values were displayed or written to these artifacts. Gemini API calls were not used. The skill document in `~/.harness` remains underspecified and was not modified. The `LLM_MEMORY_BIN` environment path is set but its binary hash differs from the canonical resolved `~/.harness` executable; it was left unchanged, and the canonical binary was invoked directly per the skill resolution order.
