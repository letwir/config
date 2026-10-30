# Task report — subagent-polygliph-audit-20261026

## Result
Audit verdict=PARTIAL, not all role instructions align with the newly adopted codec rule. Direct conflicts remain in `.HARNESS/agents/compressor.md:14-17` and `.HARNESS/agents/filereader.md:16`: both prescribe the retired fixed multilingual/token-count optimizer, uniqueness formula and Unicode-NFKC tie-break. These body rules conflict with `.HARNESS/rules/documentation.lrf:9` and the router policy.

## Positive evidence
`AGENT_ROUTER.md:1-2,45` declares POLYGLYPH and scopes it to LLM/subagent prompts and payload values. `worker.md:15,29` states the codec and wrapper boundary; `orchestrator.md` uses compact structured syntax. Pi worker profiles point to shared `../worker.md` and use compact route/guard notation.

## Protected metadata
Role YAML/frontmatter `description` strings are natural-language metadata. The active codec rule explicitly requires exact frontmatter preservation, so these are not treated as body-codec violations in this audit.

## Scope and changes
Read-only audit of the 15 Harness role Markdown files and 3 Pi worker profiles, plus router/handoff rules; searched for retired optimizer signatures. No source/prompt files changed. Finding is bounded to confirmed conflicts; no full semantic decoding proof claimed.

## Follow-up
The two contradictory role bodies require an authorized plan/scope update before any edit; this request asked for confirmation only. No secret material recorded.

## Memory ETL
Each of `knowledge.md`, `walkthrough.md`, and `diary.md` was submitted once. All CLI calls exited 0 without a verifiable receipt ID; `llm-mem eval` exited 1. No retry attempted; details withheld. Shared receipt/evaluation investigation is tracked by blocker #9.
