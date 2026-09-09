# H-07 Walkthrough

## Overview

H-07 makes `memo/<task_id>/` the canonical workspace record root and adds one local `record.json` index that links the fixed task artifacts to final external operation receipts.

## Deliverables

- A local-only atomic task-record writer and JSON Schema.
- A canonical task-local diary contract with an optional separately authorized external mirror.
- FINISH ordering that records separate ingest, evaluation, and mirror outcomes only after they are final.
- Closed `SUCCESS`, `N_A`, and `FAILED` receipt semantics with bounded reason codes.

## Verification

- Pester under pwsh: 8 passed, 0 failed.
- Pester under Windows PowerShell 5.1: 8 passed, 0 failed.
- Focused DIARY LRF lint: exit 0; proof-checker rule tree: PASS.
- Independent verifier: PASS for schema, SUCCESS, lock timeout preservation, case mismatch, and incoherent receipt rejection.
- Existing H-01 through H-06 artifacts were not rewritten.
- Verified: 2026-09-09.