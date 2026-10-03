# Walkthrough: FileReader and multilingual codec update

task_id: `20260927-subagent-codec-filereader-implementation`; timestamp: `2026-09-27T13:52:35+09:00`; target_environment: Windows 11, PowerShell 7, `C:\Users\letwir\.harness`; evidence_scope: active LRF, agent role/routing files, HSEQ local and change sequences, bounded preflight; verification_status: PASS after corrections; remaining_uncertainty: supplied W3Techs percentage source not independently verified.

## 実行計画

1. Compile the active LRF/workflow and compare the current FileReader and codec contracts with the user-provided specification.
2. Preserve existing partial FileReader files, clarify score-primary ordering with the user, and audit the 12-file implementation plan.
3. Use one selected CODE worker; reconcile codec, FileReader status outputs, persona binding, Explorer boundary and the settled-payload pause rules.
4. Run strict LRF/HSEQ/TOML/preflight checks and independently verify the final snapshot.
5. Apply verifier corrections to ADV gates, rerun focused verification, and record the final result.

## Verification

- `harness-lint.exe -path rules/documentation.lrf -level 0 -strict` → exit 0.
- `harness-lint.exe -path rules/MANUAL.lrf -level 0 -strict` → exit 0.
- `harness-lint.exe -path agents/SUBAGENTS.lrf -level 0 -strict` → exit 0.
- `harness-lint.exe -path etl/research-local.seq -strict` → exit 0.
- `harness-lint.exe -path etl/change.seq -strict` → exit 0.
- `python.exe` `tomllib` parsing of `agents/filereader.toml` and `agents/compressor.toml` → PASS.
- `invoke-rule-preflight.ps1 -Task change -Tag code,worker,doc,subagent` and receipt hash verification → PASS.
- Targeted trailing-whitespace check across all 12 allowed files → none found.
- Independent review: initial PASS gates found two ADV wording gaps in Compressor; both were corrected. Focused final verifier confirmed ADV.1 and ADV.3 match `agents/SUBAGENTS.lrf`.
- Broad `git diff --check` exited 2 on pre-existing unrelated `.system/openai-docs` trailing whitespace; targeted changed-file checks did not find trailing whitespace.
- Changed paths: `rules/documentation.lrf`, `rules/MANUAL.lrf`, `agents/SUBAGENTS.lrf`, `agents/compressor.md`, `agents/compressor.toml`, `agents/AGENT_ROUTER.md`, `agents/README.md`, `agents/filereader.md`, `agents/filereader.toml`, `agents/FILEREADER.css`, `etl/research-local.seq`, `etl/change.seq`.
- Residual risk: several targets were already modified/untracked, preventing a clean Git baseline diff. No unrelated changes were intentionally made.
