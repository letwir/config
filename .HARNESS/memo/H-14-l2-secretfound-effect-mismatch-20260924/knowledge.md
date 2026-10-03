# Knowledge Base

## Context

Level 2 harness-lint reported `EFFECT_MISMATCH` at `rules/documentation.lrf:11`: the phrase `history rewrite` in a `RO_LOCAL` payload implied mutation.

## Findings

- `doc.secret-found` combines a read-only stop/report response with conditions about history rewrite and tracked-file removal.
- The policy requires explicit current-task authorization for VCS-affecting cleanup.
- Separating the read-only response from the VCS authorization clause aligns each record's effect with its actions.

## Morphism

`secret-found → STOP+report (RO_LOCAL) ⊕ cleanup involving history rewrite or tracked-file removal → current-task authorization (VCS_WRITE)`.

## Evidence

Level 1 and Level 2 lint both report zero errors and zero warnings; `git diff --check -- rules/documentation.lrf` passes.
