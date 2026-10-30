# Knowledge Base — pi-config-link-install-20261026

## Context
A repository-distributed Pi user configuration can be installed by linking selected entries into `~/.pi/agent`; replacing or linking the whole agent tree risks overwriting extensions and runtime state. The distribution belongs under `.HARNESS/.pi/` and sources shared Harness resources from its parent `.HARNESS` directory, not from a pre-existing user `~/.harness` tree.

## Findings
The installation map links `settings.json`, only `extensions/subagent`, only `agents/pi`, individual shared agent contracts, and the relevant shared `AGENTS.md`, `POLYGLYPH.md`, `HARDWARE.md`, rules, skills, prompts, persona, ETL, scripts, evaluation, and memo resources. Pi package declarations are retained in settings but packages are not installed by the scripts. Authentication, models-store, sessions, binaries, npm installation state, and runtime data are excluded.

## Method
Preflight every source and destination conflict before mutation; refuse conflicts by default; allow a user-selected unique sibling backup; preserve backups and report rollback errors; create and verify actual symbolic links; test only against disposable temporary directories. Keep Bash and PowerShell installer behavior aligned.

## Risks
Windows may require Developer Mode or elevation for symbolic links. Native Linux/Termux are not verified by Windows/Git Bash tests. `defaultProjectTrust: always` is retained from existing settings and carries risk for untrusted projects.
