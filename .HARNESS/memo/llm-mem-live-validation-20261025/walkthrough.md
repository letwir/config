# Walkthrough — llm-mem live validation

1. Resolved the executable by the documented first-existing order. The canonical `.exe` candidate existed; no higher-priority extensionless or local-bin candidate existed.
2. Built from `C:/Users/letwir/repo/repo/llm-memory` with `GOPROXY=off`, `GOSUMDB=off`, and `GOTOOLCHAIN=local`. Removed `LLM_MEMORY_BUILD_DB_URL` from the build process so no DB URL was embedded.
3. Smoke-tested the new binary before replacement: `survey template harness-feedback` returned a valid seven-answer JSON template; malformed nested axes/unknown attribution keys failed in JSON parsing before DB access.
4. Replaced `C:/Users/letwir/.harness/skills/llm-memory/llm-mem.exe`; installed SHA-256 matched the staged executable. The installed template again parsed and contained seven questions.
5. Used the runtime `LLM_MEMORY_DB_URL` environment setting without printing its value. Read-only `search -q llm-mem-repair-20261025 -level 2 -limit 10 -json` returned four historical records. Read-only `evals -key llm-mem-repair -limit 20 -json` returned zero records before this validation's evaluation was recorded.
6. Submitted one `harness-feedback/v1` response using the CLI-generated template. It is explicitly attributed to `Pi assistant (agent observation)`, has seven answers, and is tagged as `survey_observation`; readback confirmed response ID `1046ee13-e160-4fb6-9ec1-08fca3b09d37` and task ID `llm-mem-live-validation-20261025`.
7. For task-report ingestion, invoke `llm-mem ingest -file <artifact> -cat <knowledge|diary|walkthrough>` exactly once per file. Remove `GEMINI_API_KEY` and `GEMINI_GROUNDING_API_KEY` only in that process to force the local extraction fallback and avoid external API charges. No `status` command was used because its implementation auto-registers the client.
8. Record one valid `evaluation.json` with numeric axes and `attribution.prompt_defects`/`agent_defects`, ingest each Markdown report once, then submit the evaluation once. Read back with `evals -key llm-mem-repair`; do not retry a write if its outcome is uncertain.

Checks already passed: focused offline Go tests, `go build ./...`, `git diff --check`, strict JSON smoke, installed hash check, template generation, DB history search, survey submit and readback.
