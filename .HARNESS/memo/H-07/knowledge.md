# H-07 Knowledge Base

## Context

Task-local evidence and external memory operations have different effect boundaries and must remain usable and auditable independently.

## Findings

- `memo/<task_id>/` was already the established location for walkthrough, knowledge, diary, and evaluation artifacts in H-01 through H-06.
- A local `record.json` can link stable artifact hashes and final external receipts without performing any external operation itself.
- `N_A` means no attempt was made; `FAILED` means an attempted operation failed; `SUCCESS` requires a bounded receipt ID.
- `memory_report_complete` is true only when both the ingest and evaluation receipts are `SUCCESS`; an optional diary mirror never changes that result.
- Case-sensitive task identity, reparse rejection, exclusive locking, and atomic replacement are required to keep the index tied to one exact task root.

## Morphism

`Obj(task-local artifacts + final boundary receipts) -> Mor(validate, snapshot, hash, atomically index) -> Obj(schema-valid record.json + explicit completion state)`

Verified: 2026-09-09. Supersedes the H-07 split-record uncertainty in `issues.md`.