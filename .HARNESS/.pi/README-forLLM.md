@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# Pi configuration distribution
Target=LLM ops. Codec=POLYGLYPH. HumanReadability=nonobjective. Semantics∧constraints≫token reduction. PG≠authority.

## Sources and link map
SourceRoot=parent `.HARNESS` of this `.pi`; `~/.harness`≠source. InstallTarget=`~/.pi/agent`; override=`--agent-dir` / `-AgentDir`.
Link(entries; per-entry):
- `agent/settings.json` → `settings.json`
- `agent/extensions/subagent` → `extensions/subagent`; ¬link(`extensions` whole)
- `agent/agents/pi` → `agents/pi`; ¬link(`agents` whole)
- shared `agents/worker.md`, `agents/WORKER.css`, `agents/SUBAGENTS.lrf`, `agents/AGENT_ROUTER.md` → same paths under destination `agents/`
- shared `AGENTS.md`, `POLYGLYPH.md`, `HARDWARE.md`, `rules`, `skills`, `persona`, `etl`, `scripts`, `evaluation`, `memo` → same paths under destination
- shared `crabcode/prompts` → destination `prompts`; ¬destination `crabcode/prompts`
Pi-subagents: `defaultReads`→chain cwd≠definition directory. Profiles use `~/...` references to links. Fixed routes={easy:`openai-codex/gpt-6-luna`/low,medium=same model/medium,hard:`openai-codex/gpt-6.1-sol`/high}. `agentScanDirs`=`~/.pi/agent/agents/pi`. Custom `AgentDir`⇒create equivalent `~` links or deliberately adjust config/profile references.

## Installer contract
Bash: `bash install-pi.sh [--agent-dir DIR] [--dry-run] [--backup] [--help]`.
PowerShell: `install-pi.ps1 [-AgentDir DIR] [-DryRun] [-Backup] [-Help]`.
Preflight(all sources,destination conflicts,parent components,source/destination overlap)→before writes. Equal/ancestor/descendant/symlink-parent overlap⇒refuse. DryRun⇒¬create/modify destination. Entries classify={correct|wrong|dangling symlink|regular file}; correct→idempotent. Default conflict→refuse atomically. `--backup` / `-Backup`→move conflicts to unique sibling backups; overwrite original=false. Link type=real symbolic link only. Bash validates `-L`, rejects copied files. Failure→remove created links+restore backup iff destination unoccupied; rollback error⇒explicit report. Download/package-manager/system installation/runtime/auth/session/model-data changes=false.
Windows symlink permission⇒Developer Mode or elevation may be needed. Live installation∉validation. Shared links may depend on `~/.harness`; provision separately; never substitute it as distribution source. Package declarations retained≠installed; after human review and only if requested, separately reconcile via `pi update --extensions`. Existing `defaultProjectTrust: always` retained; untrusted projects⇒risk.

## Validation
Commands: `tests/test-install-pi.sh`; `pwsh.exe -NoProfile -File tests/test-install-pi.ps1`. Report actual exit codes/platform. Coverage={fresh install,rerun,dry-run no writes,conflict atomicity,regular/wrong/dangling backups,missing source,overlap,failure rollback}. Native Linux/Termux or Windows support claim⇒require test on named platform. State review limitations.
