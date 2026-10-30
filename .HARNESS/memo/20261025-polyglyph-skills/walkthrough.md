# Walkthrough: 20261025-polyglyph-skills

## Overview
Completed on 2026-10-25. User requested reusing converted `.harness/skills` content, converting legacy bodies otherwise, applying canonical and repository copies, and continuing until verified success.

## Deliverables
- 27 canonical skills: AGY body genuinely recoded in Japanese/English/logical POLYGLYPH; 26 already-converted bodies reused with stale descriptive labels repaired.
- 13 matching repository SKILL copies synchronized while preserving each original YAML frontmatter.
- Nine source-specific reference unions preserve documented options, output fields, conceptual models and monitoring recipes; installed support is not newly asserted.
- Repository-root `POLYGLYPH.md` copied byte-identically for nearest-ancestor decoder resolution.
- Existing fixed `SIGMA/1: PIDGEN/text` examples intentionally retained because they are opaque execution-protocol literals required by shared rules, not active prose representation labels.
- Original snapshots and before/after hashes retained under this task directory. Existing unrelated repository changes were not targeted. No product Go source, binary, shared rule or agent-profile changes; no commit/push, AGY execution or external API call.

## Verification
- `pwsh.exe -NoProfile -File .../validate-final.ps1`: exit 0; 27 canonical skills + 13 repo copies pass strict UTF-8/no-BOM/bare-CR, exact canonical expected diffs, original frontmatter, opaque fences/LRF, local Markdown links, six-field LRF rows, copied-body equality, immutable snapshot and manifest hashes; codec hash matches.
- CR-aware scoped `git diff --check`: exit 0. Default Git whitespace settings flag preserved mixed CRLF as trailing whitespace; original LRF bytes were deliberately retained, not normalized.
- Independent native delegate executed validator: exit 0; authored invalid-UTF8 and deliberate-LRF-mutation negative fixtures: exit 0.
- Independent report `independent-review-final.md`: Verdict PASS, with sampled semantic-review limits. Earlier independent rejection was addressed by genuine AGY prose recoding and executable review tooling.
- Native runtime provider/model verified: openai-codex/gpt-6-luna, medium. Some inferred native attestation ledgers remain rejected because children omitted structured acceptance-report fences; these are distinguished from observed process exits and parent acceptance. No explicit gate failure is presented as a pass.

## Risks and limits
Independent semantic inspection sampled AGY and economics-api, not a formal proof of every sentence. Parent compared all unique local originals and specified the nine reference unions. LRF six-field checks and byte preservation are not a full grammar/type proof. Harness-lint/proof-checker execution unavailable/not applicable here; product source compilation and live skill integration tests were not performed. Legacy reference flags and recipes remain documentation, not installed-runtime guarantees or new permissions.

## Prompt and agent feedback
No supported user-prompt defect. Operational failures were agent-side: shell quoting, tool capability selection, premature partial delivery, unsupported language selection, byte normalization and initially omitted local references. Recorded attribution 0% prompt / 100% agent applies only to observed task friction, not a calibrated user/model score.

Tags: POLYGLYPH, skill-migration, canonical-reuse, verified, prompt-defect, agent-defect
