# Walkthrough - H-09-contradictions-20260914

## Overview
Completed the active harness contradiction repair on 2026-09-15. The user selected mandatory agy implementation, mandatory task-memory ingestion at each stop, and Astra restricted to the main agent. Codex subagents were not invoked; direct agy CLI performed the evaluator implementation and corrections.

## Plan and decisions
Inspect the active graph and preserve initial changes; align policy and sequence contracts; delegate bounded evaluator work through agy; verify actual diffs and failure boundaries; ingest the three generated reports and evaluation with actual receipts. Requested policy/docs were edited by the main agent, which also owned final validation.

## Deliverables
- BOOTSTRAP now links MANUAL and LOAD explicitly, defines precedence/effect vocabulary, requires agy for CODE, and requires bounded task-memory ingestion with task-specific titles.
- AGENTS, MANUAL, engineering, documentation, research, diary/state, agy and subagent contracts agree on those choices. Configuration-only subagent tags do not invoke roles or replace the main persona.
- main.seq starts at PRECEDENT. Local research and bounded local review have defined routes; read-only review returns to FINISH without implementation. Memory failure is recorded as FAILED rather than an approval-pending success substitute.
- The evaluator validates the split layout, exact source/target load contracts, effect vocabulary, UTF-8 decoding and exact sequence entries. It records hashes for MANUAL/LOAD and each leaf. A PowerShell boolean-precedence defect that excluded specific MUST_NOT guards was also repaired.
- 22 existing files changed relative to the starting snapshot. Exact paths and hashes: changed-files.json. Starting evidence: baseline.json and starting-status.txt. Existing unrelated changes were not staged or committed.

## Verification
- Focused Pester: 17 passed, 0 failed (tests.json). Fixtures exercise Windows PowerShell 5.1; live execution uses PowerShell 7.
- Live oracle: 30 cases matched all expected decisions, with 0 oracle mismatches. The retry-4 case correctly produces retry_exhausted and process exit 1.
- The 29 non-exhaustion cases return exit 0, update the latest receipt, and pass evaluation/run.schema.json validation.
- Combined rule preflight: 10 files and 148 selected records; receipt hash verification PASS.
- Changed LRF files passed structural lint; modified skills preserve their frontmatter and valid six-field contracts. Active rule links and PowerShell parsing passed. git diff --check returned 0.
- Knowledge ingestion: e864ae8d-5991-4506-a71a-d33d7624c05e. Diary ingestion: 6cc64760-6b89-4048-bc83-ca4e4d2f649c. All final ingestion/evaluation results are recorded in record.json and individual receipt files.

## Corrections and limitations
Initial agy outputs required main-agent corrections and retesting; tool SUCCESS was not treated as proof. Generic document titles initially collided with the active-title unique index (SQLSTATE 23505); task-specific titles resolved the collision without altering existing memories. See diary.md for AgentDefect attribution.

The .seq checks validate static contracts; they do not implement a sequence executor or prove instruction loading in a fresh session. No independent Codex subagent audit was performed. No VCS write, deployment, arbitrary DB write or external diary mirror was requested or performed.

Tags: harness, agy, memory-ingest, rule-preflight, policy-evaluation, contradictions
