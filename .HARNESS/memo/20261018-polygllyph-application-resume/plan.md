# Plan / 20261018-polygllyph-application-resume

## Scope

Apply the user-provided POLYGLYPH codec to the 59 `.HARNESS` instruction/contract Markdown targets from `scripts/tests/fixtures/polyglyph-targets.json`. Preserve human-facing README/INSTALL/HARDWARE/state/reference docs, source/config data, code and binaries. Preserve the already-applied 20-descendant RTS capacity change and unrelated llm-memory binary changes.

## Invariants

- Every target has the exact PG metadata preamble after immutable YAML frontmatter.
- PG is representation only; SIGMA/1 plus routed LRF remain sole authority. No permission/effect inference.
- Preserve `MAY`, polarity/modality, exception, scope, temporal/threshold, evidence, uncertainty, stop/success and decision boundaries.
- Preserve frontmatter, fences, embedded LRF blocks and runtime schemas byte-exact against the fixed baseline commit.
- LRF block may be followed by PG body; preserve all marker/header/records/separator bytes.
- Structural lint is not semantic proof; use bounded human semantic review for representative high-risk contracts.

## Final acceptance evidence

- 59 targets, 59/59 frozen-region signature matches, 59/59 current PG structural parses.
- `python3 -m unittest discover -s .HARNESS/scripts/tests -p test_lint_polyglyph.py -v`: 10 tests pass.
- `python3 -m py_compile` for linter and tests: PASS.
- `python3 .HARNESS/scripts/lint-polyglyph.py --root .HARNESS`: PASS; incomplete targets=0.
- `git -c core.whitespace=cr-at-eol diff --check`: PASS.
- Independent representative semantic review: five cases PASS; this is not a formal proof or claim of universal decoding.
- No VCS write, dependency install, provider call, live service action or credential use.
