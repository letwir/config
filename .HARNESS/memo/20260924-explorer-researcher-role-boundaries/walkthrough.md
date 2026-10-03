# walkthrough | 20260924-explorer-researcher-role-boundaries

## Role boundaries
- Explorer: invoked only when direct investigation is stuck or results are sparse/repetitive. Search broadly for analogs and leads without topic-similarity or source-precision filters. Label source, relation to target, confidence, and verification status. Do not recommend adoption or imply validation.
- Researcher: stay inside current task knowledge needs. Gather applicable documentation/references, OSS implementations/tests, and technical blogs. Compare candidate methods with fit, prerequisites, usage, side effects/risks, and source quality/confidence. No unrelated analog hunt.
- Router, TOMLs, Markdown role protocols, MANUAL/SUBAGENTS LRF, README index, and `etl/research.seq` report schema now agree.

## Checks
- `harness-lint.exe -path agents -level 2 -json`: no diagnostics.
- `harness-lint.exe -path etl -level 2 -json`: no diagnostics.
- `harness-lint.exe -path rules -level 2 -json`: no diagnostics.
- Python `tomllib`: 8 agent TOMLs parse and have unique names.
- `git -c core.whitespace=cr-at-eol diff --check` for the scoped files: clean.
- Runtime routing was not exercised.
