# walkthrough / 20261025-crabcode-config-links

## 実行計画
1. Inspect current Crabcode/Harness layout, global config root, link state, and H-09 setup constraints.
2. Map only required config, instruction, rule, state, agent, skill, ETL, evaluation, script, persona, and prompt paths to `.harness`.
3. Create a regular `.config/crabcode` directory with non-conflicting symbolic links; preserve `.config/agents` and OpenCode config.
4. Verify symlink targets and Harness change preflight; avoid live provider requests.

## 実施内容
- Created `C:\Users\letwir\.config\crabcode` as a normal directory.
- Added 17 symbolic links:
  - Files: `AGENTS.md`, `POLYGLYPH.md`, `HARDWARE.md`, `crabcode.jsonc`, `decisions.md`, `issues.md`, `knowledge.md`, `method.md`.
  - Directories: `agents`, `etl`, `evaluation`, `memo`, `persona`, `prompts`, `rules`, `scripts`, `skills`.
- Every target points into `C:\Users\letwir\.harness`; `crabcode.jsonc` and `prompts` point into `.harness\crabcode`.
- Existing `C:\Users\letwir\.config\agents` remains unchanged; config prompt references `../agents/*.md` continue to resolve there.

## 検証
- `invoke-rule-preflight.ps1 -Task change`: PASS.
- PowerShell directory listing showed all 17 entries as symbolic links with `.harness` targets.
- `Test-Path`: Crabcode config, rules bootstrap, worker prompt, skills dir, and `agents/SUBAGENTS.lrf` all True.
- Crabcode was not run with a model/provider; no external network, charge, or runtime config test was performed. Restart Crabcode before judging live behavior.
- No source Harness file, OpenCode config, or `.config/agents` entry was modified; only the new `.config/crabcode` directory and links plus task-scoped records were created.
