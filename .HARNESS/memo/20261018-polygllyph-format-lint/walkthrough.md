# Walkthrough / 20261018-polygllyph-format-lint

## Overview

**Incomplete.** The requested format migration and lint were not delivered because required worker execution was unavailable.

## Partial artifacts

- `.HARNESS/POLYGLYPH.md`: worker added `MAY` interop text and §8 decoder/embedding text. Current section order is §8 before §7, so this specification edit is incomplete.
- `.HARNESS/scripts/tests/fixtures/polyglglyph-targets.json`: untracked, misspelled file; contains 59 Git-root-prefixed paths rather than the settled harness-root-relative manifest shape.
- `.HARNESS/scripts/tests/fixtures/polyglglyph-preservation.json`: untracked, misspelled file; main's independent comparison against baseline `6a8889674bd01d867e9374430dbab46ebd3115b8` found fence-signature mismatches on 8 files.

## Not delivered

- Conversion of `.HARNESS/AGENTS.md`, agent prompts, skill prompts, LLM navigation, or OpenCode agent prompt.
- `.HARNESS/scripts/lint-polyglyph.py`, unit tests, verified manifest/baseline fixture, representative semantic review.
- Any successful full-corpus lint or format verification.

## Checks and evidence

- Baseline inventory count: 59 tracked candidate paths.
- Independent raw-byte hash verification of partial preservation fixture: 51/59 matched; 8 mismatched, so fixture rejected.
- `git status`: partial artifacts above plus pre-existing llm-memory changes and earlier task reports; target prompt corpus otherwise unchanged.
- `git diff --check`: exit 0, but tracked diff is empty and this does not validate untracked partial files.
- No unit tests, py_compile, LRF preflight after changes, or semantic review were run.
- No Git write occurred.

## Residual risk

Do not use the partial fixture or claim POLYGLYPH compliance. Section ordering is invalid against the plan, and no lint currently exists. A worker-capable run should regenerate and independently verify the baseline before converting any target, then finish lint/tests and perform semantic review. Preserve unrelated worktree changes.
