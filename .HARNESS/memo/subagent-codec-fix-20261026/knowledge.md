# Knowledge — subagent codec rule correction

## Facts
The obsolete optimizer was removed from `agents/compressor.md` and `agents/filereader.md`. Both now specify `Codec=POLYGLYPH`, `HumanReadability=nonobjective`, max recoverable SemanticDensity constrained by contextual equivalence and preserved safety/authority/output invariants.

## Rules
Preserve exact SIGMA/1 handoff wrapper, ordered fields, agent frontmatter, output/status schema, paths/privacy limits, effects and permissions. Ambiguity/unknown⇒preserve+STOP+ask. Symbols/emoji remain contextual cues, never sole carriers for safety/authority.

## Verification
Doc preflight PASS; retired optimizer scan across Harness role Markdown PASS; protected-contract checks, UTF-8/whitespace and `git diff --check` PASS.

## Open questions
No implementation question remains. Independent reviewer not used; main performed bounded checks.

## Failures
Memory ETL remains unverified pending receipts; no success claim.
