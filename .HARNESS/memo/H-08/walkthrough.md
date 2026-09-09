# H-08 Walkthrough

## Overview

H-08 adds a bounded `LOCAL_REVIEW` workflow for small read-only repository checks. It runs after precedent lookup and ends directly after a human-readable review report.

## Deliverables

- Three local-review LRF contracts for eligibility, commands, and escalation.
- An early `LOCAL_REVIEW` / `FULL_ETL` / `STOP` orchestration branch.
- A fixed-command PowerShell collector and strict JSON Schema.
- Derived policy-evaluator routing and expanded oracle coverage.
- Documentation for receipt meaning, exit codes, supported shells, and terminal boundaries.

## Verification

- Collector Pester: pwsh 6/6; Windows PowerShell 6/6.
- Policy evaluator Pester: pwsh 7/7; Windows PowerShell 7/7.
- Live two-target collector: exit 0, PASS, schema valid, `git diff --check` exit 0.
- 30-case oracle: local review, code/public/external/live/VCS expansion, malformed targets, credential forbid, and retry-4 classified as expected.
- `evaluation/oracle.json` and generated run receipts validate against their schemas.
- Independent Verifier: PASS after correction of whitespace and credential decision preservation.
- No VCS write was performed.

Verified: 2026-09-10.
