# LLM-only DOC_RULE implementation walkthrough

## Overview

Added a scoped token-minimization rule for files explicitly designated LLM-only.

## Deliverables

- `rules/documentation.lrf`: new `doc.llm-only-codec` and audience-aware SKILL contract.

## Verification

- `harness-lint.exe -path rules/documentation.lrf -level 1`: passed with zero errors.
- `git diff --check -- rules/documentation.lrf`: passed.
- Level 2 reported one existing diagnostic at `rules/documentation.lrf:11`; unrelated to this change.

## Constraints

No commit or push was part of this task. Existing unrelated working-tree changes remain untouched.
