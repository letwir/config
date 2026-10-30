# knowledge / 20261025-crabcode-config-links

- topic: Global Crabcode configuration links into the Harness SSOT
- timestamp: 2026-10-25 13:23 local
- target environment: Windows, `C:\Users\letwir`, Crabcode 0.0.14; `XDG_CONFIG_HOME` unset
- evidence/source scope: `.config/crabcode`, `.harness` layout, local Crabcode 0.0.14 resolver/skill source, Harness H-09 and local preflight; no model/provider invocation
- verified facts:
  1. Created a normal directory at `~/.config/crabcode` (not a directory symlink) and added 17 NTFS symbolic links into `~/.harness`.
  2. File links: `AGENTS.md`, `POLYGLYPH.md`, `HARDWARE.md`, `crabcode.jsonc`, `decisions.md`, `issues.md`, `knowledge.md`, `method.md`.
  3. Directory links: `agents/`, `etl/`, `evaluation/`, `memo/`, `persona/`, `prompts/`, `rules/`, `scripts/`, `skills/`.
  4. `crabcode.jsonc` points to the existing canonical `~/.harness/crabcode/crabcode.jsonc`; `prompts/` points to its existing sibling prompt assets. The file config references `../agents/*.md`; existing `~/.config/agents` was left untouched.
  5. PowerShell `Get-ChildItem` showed the 17 entries as symbolic links with the intended `.harness` targets. `Test-Path` returned True for the global config, `rules/BOOTSTRAP.lrf`, `prompts/worker-prompt.md`, `skills/`, and `agents/SUBAGENTS.lrf`.
  6. `invoke-rule-preflight.ps1 -Task change` returned PASS.
- remaining uncertainty: Crabcode was not started with a live model request, so runtime effective-config diagnostics and provider/model access were not tested. Restart Crabcode for the new global config directory to be read.
- verification status: filesystem links and key target paths verified; live application startup not tested; no existing OpenCode config or `~/.config/agents` entries changed.
