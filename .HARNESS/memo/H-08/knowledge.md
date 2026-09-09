# H-08 Knowledge Base

## Context

Small, read-only repository reviews need a bounded path that preserves policy classification without invoking the full research, implementation, and recording workflow.

## Findings

- `LOCAL_REVIEW` is eligible only for `RO_LOCAL`, a single bounded fixed-string query, and 1–32 exact relative file targets.
- The collector directly invokes only fixed-string `rg`, target-limited `git status --short`, and target-limited `git diff --check`.
- Command output is bounded in memory and persisted only as SHA-256, command identity, and exit code; search matches and Git output are absent from the receipt.
- An `rg` no-match exit of 1 is valid evidence. A required command failure yields `FAIL`; malformed or escaping input yields `STOP`.
- Edit, code/config/schema, public, external, live, VCS, credential, mixed, or expanded scope returns to `FULL_ETL`. Unknown effects and malformed contracts stop.
- The local-review terminal ends after `ReviewReport → Human`; it cannot reach FINISH, task records, diaries, llm-memory, or VCS writes.

## Morphism

`Obj(RO_LOCAL + fixed query + finite targets) -> Mor(validate, collect fixed evidence, hash) -> Obj(schema-valid receipt + ReviewReport -> End)`

Verified: 2026-09-10. Resolves H-08 in `issues.md`.
