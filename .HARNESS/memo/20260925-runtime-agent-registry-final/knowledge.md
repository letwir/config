# Knowledge — Compressor/Orchestrator runtime registry — 2026-09-25

## Facts
- Codex custom role declarations are registered through `agents.<name>.config_file`; the existing `.codex/agents` directory symlink resolves to `.harness/agents`.
- Added `agents.compressor` and `agents.orchestrator` config-file references in `C:\Users\letwir\.codex\config.toml`.
- Added OpenCode global agent definitions under `C:\Users\letwir\.config\opencode\agents\` for both roles, with role-specific model/variant and permissions.
- Corrected both role TOMLs to the repository's top-level schema and read-only sandbox mode.
- Added role-bound Finnish Compressor and Austrian Orchestrator CSS. Updated Worker CSS with `🏫🇯🇵⊕Ojou-sama×Chaotic` and a matching voice description.
- Updated `AGENT_ROUTER.md`, `MANUAL.lrf`, and `SUBAGENTS.lrf` to register the role-specific styles while preserving presentation-only semantics and flexible selection for other roles.
- Python TOML parsing and role config path checks passed. Codex Doctor reports `config.toml` parse OK. MANUAL/SUBAGENTS LRF level 1 checks passed. `git diff --check` emitted no diagnostics on checked tracked targets.
- A native child-dispatch smoke test did not produce verifiable child IDs/results. `functions.task` still reported the roles unknown in the already-running session; OpenCode agent files are loaded at startup, so a restart is required before retesting.

## Residual constraint

Registration is configured, but live dispatch remains unverified until OpenCode/Codex reloads the configuration and a fresh session returns actual child responses.
