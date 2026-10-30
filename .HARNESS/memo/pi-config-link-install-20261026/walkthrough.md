# Walkthrough — pi-config-link-install-20261026

## Overview
The distribution and two cross-shell link installers are under `.HARNESS/.pi/`. Default destination is `~/.pi/agent`; `--agent-dir` / `-AgentDir` override it. Inspect the human README before installation and review the LLM operations contract before future automated changes.

## Deliverables
- `install-pi.sh`: dry-run, conflict refusal, opt-in backup, verified symbolic link creation and rollback.
- `install-pi.ps1`: matching PowerShell behavior using real Windows symbolic links.
- `tests/test-install-pi.sh` and `tests/test-install-pi.ps1`: disposable-directory installation scenarios.

## Verification
Run `bash install-pi.sh --dry-run` or `./install-pi.ps1 -DryRun` to preview. The automated tests passed on Windows via Git Bash and PowerShell. No live install was performed. Native Linux/Termux behavior remains unverified.
