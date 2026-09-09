# Walkthrough

## Overview

H-05 adds a reproducible policy and workflow evaluation command. The evaluator treats `evaluation/oracle.json` as inert input and comparison data, then derives decisions from the live LRF sources and retry behavior from `etl/main.seq`.

## Deliverables

- Structured case and run schemas.
- A fail-closed PowerShell evaluator with atomic receipts.
- Twenty-one live cases covering effects, authorization, routing, conflicts, and retries.
- An isolated Pester suite and one documented rerun command.

## Verification

- Pester: 5 passed, 0 failed.
- Live receipt: 21 cases; only `retry-4` is nonzero by design.
- Receipt schema: valid.
- Independent verifier: PASS.
- Verified: 2026-09-09.