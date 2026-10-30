# Rule — LLM-facing POLYGLYPH representation

## Scope
Target=artifact positively designated for LLM use by current request, filename, role/location, or existing contract (e.g. `README-forLLM`, agent/subagent prompt, handoff). Human-only/operator docs without such designation are excluded. Mixed readers≠human-readability requirement when LLM designation exists.

## Requirements
Explicit LLM target⇒POLYGLYPH; HumanReadability=nonobjective. Objective=max recoverable SemanticDensity subject to Meaning∧Constraints∧DecisionBoundary invariants. Preserve polarity, modality, time, exception, uncertainty, authority, effect, safety, scope, acceptance, stop/success. Representation≠permission/authority. Replaced legacy language/tokenizer/NFKC tie-break optimizer. `SKILL.md` explicit-designation branch aligned.

## Procedure
Encode eligible prose and task-value payloads only. Preserve frontmatter/bootstrap/meta, every wrapper/schema/field order/terminal marker, LRF rows, opaque fences/code, inline literals, links/filenames, and all operational identifiers anywhere. Ambiguity/unknown⇒preserve+STOP+ask. Handoff `SIGMA/1: PIDGEN/text` remains protocol literal; payload=PG.

## Exceptions
Unmarked human/operator material⇒existing human-auditable contract. Machine schema, code, commands, paths, model IDs, flags, keys, selectors, status/effect tokens and exact wrapper literals are not rewritten as prose. No authority/effect/provider change.

## Verification
Preflight change+doc PASS; compare protected bytes/markers/fences/identifiers; ensure old optimizer absent and both explicit LLM/human-only audience branches remain distinguishable; report review/runtime limitations truthfully.
