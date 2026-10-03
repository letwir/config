# LLM-only DOC_RULE codec

## Context and plan

The user requested that explicitly LLM-only files use a DOC_RULE expressing token-minimal semantic and pragmatic preservation, ambiguity reduction, a compact language/symbol/emoji vocabulary, and preservation of intent, modality, time, polarity, and register. The user clarified this belongs in the instruction rules and subagent-facing instruction path, with CSS remaining only a thin persona layer. The current task edits the canonical documentation LRF rule set.

## Verified changes

- Added `doc.llm-only-codec` to `rules/documentation.lrf`.
- The rule only applies when the user or an existing file contract explicitly marks the target LLM-only. Mixed, human-facing, or unclear audiences do not inherit the codec.
- It preserves semantic/pragmatic fidelity, ambiguity reduction, authority, effects, safety, negation, exceptions, stop/success conditions, and exact schemas. It states that symbols and emoji are cues, never sole policy carriers.
- Updated the SKILL documentation contract to route explicitly LLM-only files through this rule and keep human-facing or mixed files human-auditable.
- Level 1 harness-lint passed; `git diff --check` passed. Level 2 still reports a pre-existing `doc.secret-found` diagnostic at line 11; it concerns existing LRF text and was not changed in this task.
- The working tree contains unrelated existing edits; they were preserved. No commit or push was requested in this task.

## Outcome

The DOC_RULE now expresses the requested codec as an audience-scoped optimization rule, while LRF retains authority for policy and protocol meaning.
