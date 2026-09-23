# Walkthrough — H-10 SPR/PIDGEN subagent wrapper

## Changes

- Added the exact `<Γ>` / `SIGMA/1: PIDGEN/text` / `</Γ>` envelope and its policy-boundary meaning to the subagent and agy contracts.
- Added parent-side fail-closed validation requirements and kept the five handoff fields in order.
- Updated the agy router example and all nine agent TOMLs.

## Verification

- `harness-lint.exe -path rules/MANUAL.lrf -level 2 -json`: pass.
- `harness-lint.exe -path rules/agy.lrf -level 2 -json`: pass.
- `harness-lint.exe -path agents/SUBAGENTS.lrf -level 1 -json`: pass.
- Parsed all nine `agents/*.toml` files with `tomllib`: pass.
- `git diff --check` on tracked agent and router edits: pass.
- No test suite run. SUBAGENTS LRF Level 2 remains blocked by its pre-existing `subagent.mode` effect mismatch at line 51.