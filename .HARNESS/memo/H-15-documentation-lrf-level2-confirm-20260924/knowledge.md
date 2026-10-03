# Knowledge Base

## Context

The user reported one remaining Level 2 effect mismatch at `rules/documentation.lrf:11`.

## Findings

- Current line 11 is `doc.secret-found`, with effect `RO_LOCAL` and payload `STOP and report` only.
- VCS-affecting cleanup conditions are isolated on line 12 under `VCS_WRITE`.
- Current Level 2 output reports zero errors and zero warnings.

## Morphism

`secret-found → STOP+report (RO_LOCAL)`; `history rewrite or tracked-file removal → current-task authorization (VCS_WRITE)`.

## Evidence

Harness-lint Level 1 and Level 2 both pass. The current turn made no source changes.
