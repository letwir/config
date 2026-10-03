# Knowledge Base - H-09-contradictions-20260914

## Context
- Task: H-09-contradictions-20260914; scope: C:/Users/letwir/.harness active rules, split ETL and validation scripts.
- User decisions on 2026-09-14: agy is mandatory for implementation; ingest task memory at every stop; Astra must not be a subagent. No Codex subagents were invoked.

## Findings
- Verified from live files: BOOTSTRAP had no L edges to MANUAL/LOAD; the preflight traversed only its root. Explicit L records restore machine traversal.
- Verified: main.seq had @entry ROUTE while PRECEDENT was required; the evaluator still expected the removed monolithic layout.
- Verified: unconditional agy implementation conflicted with opt-in routing; unconditional memory ingest conflicted with pending-approval completion branches.
- Revised scope: standing authorization covers bounded CODE agy invocation and task report ingestion only. Arbitrary live writes, VCS, publication and external messaging remain gated.
- Static .seq contracts are agent instructions, not an executable orchestration engine. Static checks do not establish end-to-end runtime dispatch.

## Morphism
Task evidence -> bounded policy edits and agy implementation -> real diff, rule traversal and boundary checks -> task memory receipts.

## Sources and uncertainty
Sources: rules/BOOTSTRAP.lrf, rules/engineering.lrf, rules/documentation.lrf, etl/*.seq and scripts/invoke-policy-evaluation.ps1 inspected on 2026-09-14. Final tests and receipts are recorded in walkthrough.md and record.json. New-session instruction loading remains unverified.

Tags: harness, contradiction, agy, memory, preflight, split-etl
