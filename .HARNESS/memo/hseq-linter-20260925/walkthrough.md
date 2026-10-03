# Walkthrough: HSEQ lint support

## Change

- Added `.seq` single-file and recursive directory dispatch to `harness-lint`, preserving `.lrf` handling and existing diagnostic formats.
- Added HSEQ structural checks and the first-line dialect marker `@dialect harness-seq/1` to all five workflow sequence files.
- Aligned `change.seq` entry with its first section (`PLAN`). Documented the dialect and lint invocation in the skill and LLM navigation.
- Built the installed `harness-lint.exe` from `mcp2go/src/harness-lint`.

## Evidence

- `go build -o C:\Users\letwir\.harness\skills\harness-lint\harness-lint.exe .` — exit 0.
- `harness-lint.exe -path C:\Users\letwir\.harness\etl -strict -json` — exit 0, output `null` (no diagnostics).
- `harness-lint.exe -path C:\Users\letwir\.harness\rules\BOOTSTRAP.lrf -level 0 -json` — exit 0, output `null` (no diagnostics).
- `git -c core.whitespace=cr-at-eol diff --check -- etl/main.seq skills/harness-lint/SKILL.md` — exit 0.
- The first AGY worker completed its Go edits and formatting. The second AGY invocation returned API `UNAVAILABLE`/503 after writing its allowlisted changes; the actual files were inspected. No alternate backend was used.

## Limits

No tests were added or run. The lint checks header and section structure; it does not parse the complete HSEQ expression language or prove sequence execution. No runtime consumer was identified or exercised. Existing unrelated dirty files were preserved.
