# Walkthrough / 20261018-polygllyph-application-resume

## Result

Completed the resumed POLYGLYPH application to the 59-file manifest. The linter is read-only and explicitly reports structural validation only, not semantic proof.

## Implementation

- `.HARNESS/POLYGLYPH.md`: ordered §0–§8, retained MAY, defined exact preamble, decoding ambiguity/stop behavior, CommonMark fence constraints, and embedded LRF/1 block boundaries/resumption.
- `.HARNESS/scripts/lint-polyglyph.py`: stdlib-only checks for exact inventory, encoding, markers, metadata, frozen regions, fences, LRF blocks/records, SIGMA serialization, and Γ handoff.
- `.HARNESS/scripts/tests/test_lint_polyglyph.py`: 10 positive/negative tests including malformed fences, LRF block embedded before resumed body, EOF blank bytes, SIGMA symbol and handoff order.
- `.HARNESS/scripts/tests/fixtures/polyglyph-targets.json` and `polyglyph-preservation.json`: exact 59-file inventory and frozen signatures from baseline `6a8889674bd01d867e9374430dbab46ebd3115b8`.
- 59 manifest target documents updated; pre-existing `agents/AGENT_ROUTER.md` 20-descendant limit preserved.
- `.HARNESS/memo/20261018-polygllyph-format-lint/semantic-review.md` contains detailed representative source/target comparisons; verifier and Critic reviews are recorded in the task transcript.

## Verification

- `python3 -m unittest discover -s .HARNESS/scripts/tests -p test_lint_polyglyph.py -v` — 10 passed.
- `python3 -m py_compile .HARNESS/scripts/lint-polyglyph.py .HARNESS/scripts/tests/test_lint_polyglyph.py` — PASS.
- `python3 .HARNESS/scripts/lint-polyglyph.py --root .HARNESS` — PASS, incomplete targets=0.
- Independent raw Git baseline comparison — 59/59 frozen signatures match; all 59 current targets parse.
- `git -c core.whitespace=cr-at-eol diff --check` — PASS.
- LRF preflight and receipt verification for change/code/doc/state/worker — PASS (9 routed files).
- Representative independent semantic review — PASS: worker handoff, RO search, LIVE_WRITE memory, plugin/config exceptions, OpenCode orchestrator.

## Residual risk and preservation

The linter establishes structure and frozen byte regions only, not semantic equivalence. The five representative reviews provide bounded evidence but are not a formal proof or universal decoder guarantee. Cross-model agreement remains empirical. Rikei-search body was already a compact symbolic mapping and was left unchanged apart from the standard bootstrap metadata. Existing llm-memory binary modifications, prior task records, and all unrelated state remain untouched. No VCS write occurred.
