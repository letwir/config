# Walkthrough — Compressor/Orchestrator runtime registry — 2026-09-25

## Result

Configuration changes completed; live dispatch requires a fresh OpenCode/Codex session and remains unverified in this session.

## Changes

- Registered `compressor` and `orchestrator` with `agents.<name>.config_file` in the Codex user config, resolving through `.codex/agents` → `.harness/agents`.
- Added global OpenCode agent definitions in `.config/opencode/agents/compressor.md` and `orchestrator.md` so the current OpenCode task-agent registry can load the role names on restart.
- Corrected both role TOMLs to `sandbox_mode = "read-only"` and role-specific developer instructions.
- Refined `compressor.css` to Finnish slender stone-mill girl with timid words/strong convictions/taciturn archetype; refined `orchestrator.css` to Austrian city-girl violinist/いいんちょ and kept the wavelength catchphrase.
- Added `🏫🇯🇵⊕Ojou-sama×Chaotic` to `WORKER.css` and updated its voice cue.
- Updated router and narrow MANUAL/SUBAGENTS LRF exceptions for the two role-bound CSS files.

## Checks

- Python TOML parse and role config path resolution: PASS.
- `codex doctor --json --summary`: config parse OK.
- MANUAL and SUBAGENTS LRF level 1: PASS.
- Scoped `git diff --check`: no diagnostics.

## Runtime evidence and limit

The running `functions.task` registry rejected the role names before reload. A separate Codex CLI smoke run did not return child IDs, so there is no verified child response. OpenCode config is startup-loaded; quit/restart OpenCode, then invoke `compressor` and `orchestrator` in a fresh session to complete the runtime test.
