# Knowledge — POLYGLYPH policy migration

## Facts
`documentation.lrf` now requires POLYGLYPH for positively designated LLM-facing artifacts and marks human readability nonobjective. Legacy fixed-language/tokenizer/NFKC optimizer was replaced; mixed-audience SKILL rule aligned. Router, worker, orchestrator, and Pi README-forLLM were updated.

## Invariants
Meaning+Constraints+DecisionBoundary preserved. Exact frontmatter, bootstrap, code/fences, operational identifiers, schemas, labels/order, wrappers and terminal markers remain immutable. `SIGMA/1: PIDGEN/text` is outer protocol; payload prose uses PG. No authority/effect/permission changes.

## Verification
Doc preflight and protected-marker/fence/literal/encoding checks passed. Attacker plan review v3 returned NO_FINDINGS; independent post-edit Verifier unavailable, so main verification only.

## Open questions
Other noncentral role prompts remained outside the prior bounded migration; subsequent audit found two stale optimizer clauses, handled in separate task `subagent-codec-fix-20261026`.
