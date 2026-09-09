# H-06 Walkthrough

## Overview

H-06 adds one bounded command that checks exact tool paths, executable identity, safe representative operations, and configured client catalog visibility without merging those states into one success claim.

## Deliverables

- A manifest for five core tools, two required clients, and three optional clients.
- A fail-closed PowerShell runner with final-path containment, bounded process execution, redaction, and atomic receipts.
- A receipt schema and seven isolated failure and lifecycle tests.
- A documented live rerun command and status semantics.

## Verification

- Pester: 7 passed, 0 failed.
- Live receipt: exit 0; schema valid; latest updated atomically.
- Core tools: 5 Healthy; safe representative checks: 3 Passed.
- Clients: Gemini Observed; Codex Unverified; optional three NotConfigured.
- Independent verifier: PASS, including independent process-tree and reparse escape probes.
- Verified: 2026-09-09.
