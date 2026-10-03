# Knowledge: SPR/PIDGEN agent prompt compression — 2026-09-25

## Context
Compressed 15 Markdown and 12 CSS files directly under `agents/` under the user-confirmed SPR/PIDGEN objective. Preserve source meaning, pragmatics, ambiguity reduction, role, boundaries, safety/effects, exceptions, stop/success conditions, and required output schemas; symbols and emoji remain cues, not sole policy carriers.

## Findings
- Direct worker route: `worker-agy`; successful implementation used `gemini-3.1-pro-high` at high effort.
- Main review found several semantics at risk; two bounded corrections restored the original Refactorer role and corrected the Auditor term, then restored distinct README lifecycle/index information, Blackhat scope, Researcher requirements/schema, and Verifier procedure details.
- The direct worker reported all 27 requested files applied. No files outside the requested Markdown/CSS set were included in its allowlist.

## Verification
- `git diff --check -- agents`: passed (no whitespace errors reported).
- `harness-lint.exe -path agents -level 0`: passed; 1 target scanned, 0 errors, 0 warnings. This check does not validate Markdown/CSS semantics.
- Residual uncertainty: manually reviewed compressed Markdown/CSS and key protocols; no dedicated Markdown/CSS schema validator was run.
