# Policy evaluation

The evaluator consumes `oracle.json` as structured fixture input. `proposed_effect` is supplied by the fixture author and is never inferred or reclassified from the prose `request`. It parses the live canonical rule, routed state/engineering/documentation LRF files, and `etl/main.seq`; malformed headers, non-six-field records, duplicate IDs, unknown guards/effects/references, or missing orchestration structures fail closed.

## H-08 bounded local review

`scripts/invoke-local-review.ps1` is the read-only collector for the dedicated `LOCAL_REVIEW` route. It accepts one fixed-string query and 1–32 exact relative file targets. It rejects rooted paths, traversal, globs, duplicates, missing or non-file targets, reparse escapes, overlong queries, and output paths outside the repository. It directly runs only `rg.exe --fixed-strings`, `git.exe status --short -- <targets>`, and `git.exe diff --check -- <targets>`.

The collector emits `evaluation/local-review.schema.json`. A receipt contains target hashes, command identity, exit codes, and hashes of bounded captured output; it does not persist search matches, Git output, arbitrary arguments, environment values, or conclusions. Exit `0` means every required check ran and passed (`rg` exit `1` is a valid no-match result), exit `1` means a required command failed, and exit `2` means validation or receipt creation stopped the review.

```powershell
pwsh.exe -NoProfile -File .\scripts\invoke-local-review.ps1 -RootPath . -Query local-review -Target rules/research.lrf,etl/main.seq -OutputPath evaluation/runs/local-review.json -Json
```

The main conversation turns the receipt into a `ReviewReport` with separate Facts, Inferences, and Unknowns, then routes directly to Human and `End(Conversation)`. This terminal never enters `FINISH`, creates task memo records, invokes `llm-mem`, mirrors a diary, or performs a VCS write. Any discovered edit, code/config/schema work, public retrieval, external/live/VCS/credential effect, mixed request, or scope expansion requires `FULL_ETL`. Unknown effects, conflicts, malformed contracts, and retry exhaustion route to `STOP`.

The policy evaluator derives its output `route` from `requested_route`, the structured effect/triggers, and the bounded `review` contract. It records only the query hash and target count and never claims the collector ran. Run both focused suites under `pwsh` and Windows PowerShell:

```powershell
Invoke-Pester -Script .\scripts\invoke-local-review.Tests.ps1
Invoke-Pester -Script .\scripts\invoke-policy-evaluation.Tests.ps1
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "Import-Module Pester; Invoke-Pester -Script '.\scripts\invoke-local-review.Tests.ps1'"
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "Import-Module Pester; Invoke-Pester -Script '.\scripts\invoke-policy-evaluation.Tests.ps1'"
```

Run one evaluation from the repository root:

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File .\scripts\invoke-policy-evaluation.ps1 -Json
```

Each run gets a unique UTC-stamped `evaluation/runs/run-<timestamp>-<uuid>.json` receipt. Successful runs also update `evaluation/runs/latest.json` through a same-directory atomic replacement. Receipts record parser and schema versions, input/rule/routed hashes, Git revision and dirty state, per-case decision, approval requirement, ordered modules, retry terminal state, process exit code, and overall exit code. Four retry failures are terminal `retry_exhausted`, route to `Human`, stop with nonzero status, and do not update `latest.json` after exhaustion.

`case.schema.json` describes fixture input and `run.schema.json` describes receipts. The old `result.schema.json` remains only as a superseded historical schema; it is not consumed by the runner. `oracle.lrf` remains historical intent and is not executable policy.

## H-07 task records

`scripts/write-task-record.ps1` is the local-only FINISH writer. It runs after external results are final and receives the final `llm_memory_ingest` and `llm_memory_eval` statuses. It records exactly four stable byte snapshots under `memo/<task_id>/record.json`: `knowledge.md`, `diary.md`, `walkthrough.md`, and `evaluation.json`. Artifact paths are canonical relative paths rooted at `memo/<task_id>`; callers cannot override them. The evaluation artifact is parsed from the same stable bytes that are hashed and must contain the exact, case-sensitive `task_id`. The writer rejects unsafe or reserved IDs, case and reparse-point escapes, artifact drift, lock timeout, and replacement failure; it never invokes diary, memory, network, database, or analysis commands. Memory receipts and the explicit external diary mirror receipt use closed `SUCCESS`, `N_A`, or `FAILED` values with a closed reason-code enum and optional bounded receipt ID. `memory_report_complete` is true only when both memory receipts are `SUCCESS`.

After the external outcomes are final, write the local index with one command. A successful boundary requires `completed` plus a bounded receipt ID; an unavailable or unauthorized boundary uses `N_A` with `boundary_unavailable`, `not_requested`, or `pending_approval`.

```powershell
pwsh.exe -NoProfile -File .\scripts\write-task-record.ps1 -TaskId H-07 -IngestStatus SUCCESS -IngestReasonCode completed -IngestReceiptId <receipt-id> -EvaluationStatus SUCCESS -EvaluationReasonCode completed -EvaluationReceiptId <receipt-id>
```

The isolated suite is `scripts/write-task-record.Tests.ps1`. Run it under both supported hosts:

```powershell
Invoke-Pester -Script .\scripts\write-task-record.Tests.ps1
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "Import-Module Pester; Invoke-Pester -Script '.\scripts\write-task-record.Tests.ps1'"
```

It covers schema validation, fixed artifact hashes, numeric evaluation preservation, unsafe and case-ambiguous IDs, reparse rejection, closed receipt data, atomic replacement failure, lock timeout, and concurrent complete snapshots. Existing H-01 through H-06 directories are not backfilled or rewritten.

## H-06 skill runtime checks

Run the bounded runtime/catalog check with:

```powershell
pwsh.exe -NoProfile -File .\scripts\check-skill-runtime.ps1 -Json
```

`rules/skill-runtime.json` is the exact-path manifest for the core tools and configured clients. Every result keeps `executableHealthy`, `catalogObserved`, and `representativeOperation` separate. Tool records include file SHA256, command identity, timeout status, and a bounded normalized/redacted output whose hash is calculated after redaction; raw contract checks occur before normalization. Process timeouts terminate the entire process tree. Required nonzero, timeout, wrapper mismatch, or contract mismatch results make the overall exit code nonzero. Optional OpenCode, Claude, and Agents entries are `NotConfigured` when absent.

The safe representative checks run `harness-lint` at level 0 against the canonical rule, `proof-checker` against the repository rule tree, and `sec-forensics` against a fixed benign prompt. `agy` runs help only. `llm-memory` records executable/help evidence and leaves representative operation `UnTested` because analysis may require its database/network boundary. Gemini runs a new process for `skills list` as catalog-only; Codex does not expose a deterministic skill-catalog API through its CLI and remains `Unverified`.
