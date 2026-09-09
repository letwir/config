# Knowledge Base

## Context

An oracle is not execution evidence when no consumer compares it with the active rules and workflow.

## Findings

- `harness-lint` does not consume `evaluation/oracle.*`, `result.schema.json`, or `etl/main.seq`; using it here would not prove the H-05 contract.
- Natural-language requests do not have an authoritative classifier in the current LRF. The evaluator therefore accepts `proposed_effect` as structured fixture input and never reclassifies request prose.
- The active policy and workflow sources are hashed into each receipt together with the evaluator and schemas.

## Morphism

`Obj(structured case + live LRF/main.seq) -> Mor(parse, derive, compare, record) -> Obj(schema-valid receipt + exit code)`

Verified: 2026-09-09. Supersedes no prior knowledge.