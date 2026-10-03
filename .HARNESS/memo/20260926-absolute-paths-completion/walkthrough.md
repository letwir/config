# Walkthrough: Absolute-Path Cleanup Completion

## Overview
Generalized user-home and fixed Windows installation paths in maintained prompts, configuration guidance, and Chrome discovery.

## Deliverables
- Updated nine agent TOML instructions, `agents/SUBAGENTS.lrf`, `agents/verifier.md`, `issues.md`, `scripts/sync-harness.ps1`, both default-settings skill files, and `skills/llm-memory/SKILL.md`.
- Replaced fixed Chrome install paths with environment-derived candidates in `skills/chrome-cdp-search/scripts/chrome_cdp_search.mjs`.
- Preserved exact-command rules and historical evidence files.

## Verification
- PowerShell parser: PASS.
- `node --check skills/chrome-cdp-search/scripts/chrome_cdp_search.mjs`: PASS.
- TOML parser: PASS (9 agent files).
- `harness-lint.exe -path agents\SUBAGENTS.lrf -level 0 -strict`: PASS.
- Targeted `git diff --check`: PASS.
- Search of the approved targets for user-specific home paths and fixed Program Files literals: no matches.
- Repository-wide `git diff --check` still reports existing trailing whitespace in unrelated `skills/.system/openai-docs/` files.

## Residual Risk
`LLM_MEMORY_REPO` must point to the repository root for that optional executable fallback. Exact-command rules and historical records still contain absolute paths by design.
