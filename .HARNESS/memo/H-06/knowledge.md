# H-06 Knowledge Base

## Context

File distribution does not establish that a binary starts, its help contract matches, a representative operation succeeds, or a client discovers the skill.

## Findings

- Executable health, catalog observation, and representative operation need separate fields because each can have a different result.
- Path checks must resolve reparse points before containment checks.
- A timeout must terminate the process tree; killing only the wrapper can leave the actual client running.
- Output contracts are checked against raw process output before bounded redaction, while receipts persist only bounded redacted text and its hash.
- Codex help proves the wrapper starts but does not provide deterministic skill-catalog evidence, so its catalog state remains `Unverified`.

## Morphism

`Obj(exact manifest + installed tools/clients) -> Mor(resolve, identify, invoke bounded, classify, record) -> Obj(schema-valid runtime receipt + exit code)`

Verified: 2026-09-09. Supersedes no prior knowledge.
