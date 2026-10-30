# Knowledge — Subagent POLYGLYPH compliance audit

## Facts
Inventory=15 `.HARNESS/agents/*.md` role instructions + 3 `.HARNESS/.pi/agent/agents/pi/*.md` profiles. Router/worker/orchestrator declare PG for LLM-target prose and preserve protocol/schema boundaries. Confirmed conflicting obsolete optimizer remains in `agents/compressor.md:14-17` and `agents/filereader.md:16`: fixed language set, tokenizer-based token savings, uniqueness formula, NFKC tie-break.

## Rules
`rules/documentation.lrf:9` applies to role/location-designated LLM instructions and replaces the legacy codec; `AGENT_ROUTER.md:45` repeats this. `frontmatter` is exact-protected; natural-language agent descriptions there are not body violations under current contract. `SIGMA/1: PIDGEN/text` remains immutable wrapper, not payload codec.

## Constraints
Audit request=read-only; no instruction source edited. Search confirms direct contradiction only in the two cited roles; it does not certify every sentence in all prompt bodies as optimally encoded.

## Open questions
Whether to expand the prior bounded migration plan and update `compressor.md` and `filereader.md` requires a new implementation request/plan.

## Failures
None in local search. Earlier llm-mem ingest/eval receipt issue remains separate and unresolved.
