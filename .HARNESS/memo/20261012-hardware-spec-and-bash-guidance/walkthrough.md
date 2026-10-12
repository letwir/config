# walkthrough / 20261012-hardware-spec-and-bash-guidance

依頼: この環境を .harness/HARDWARE.md に記録しつつ、それを元に各地のコマンド実行方法を HARDWARE.md を参照して環境に合わせたシェルを使うように誘導する。

変更内容:
1. `~/.harness/HARDWARE.md`: ホストハードウェア仕様（Sony XQ-GE44, Snapdragon 8 Gen 3, 16GB RAM, 449GB storage）、OS/カーネル（Android 16, Linux 6.12.38-android16, aarch64）、Termuxランタイム、シェル規約（既定/必須 Bash、PowerShell/cmd厳格禁止、POSIXパス、tmpdir注意）、ネイティブツールチェイン、旧PowerShell→Termux Bash移行対応表を正本（SSOT）として新規作成。
2. `~/.harness/AGENTS.md` (および `$HOME/AGENTS.md`): ホスト環境とコマンド実行について `HARDWARE.md` を参照し環境適応シェル（bash）を使用する旨、PowerShell/cmd禁止を明記。
3. `~/.harness/rules/BOOTSTRAP.lrf`: `R|host|*` に `hardware=../HARDWARE.md` を追加・更新し、コマンド実行時の遵守を規定。
4. `~/.harness/rules/README-forLLM.md`: LLM向けナビゲーション例を PowerShell / Windows パスから Bash (`rg`, `invoke-rule-preflight.sh`) に更新し、`HARDWARE.md` 参照指針を明記。
5. `~/.harness/rules/MANUAL.lrf` & `engineering.lrf`: `$env:...` 表記を `$HOME`, `$LLM_MEMORY_BIN` に是正。
6. `~/.harness/skills/default-settings/SKILL.md` (および `SKILL_default.md`), `skills/agy-subagent-router/SKILL.md`, `skills/economics-api/USAGE.md`, `skills/rikei-search/USAGE.md`, `method.md`: シェル指定・コマンド例を Bash / `HARDWARE.md` 準拠に更新。
7. `~/.harness/issues.md`: H-11 としてホスト仕様正本化と各地コマンドの Bash 誘導完了を記録。

検証:
- `bash scripts/invoke-rule-preflight.sh --task change`: 終了コード 0、全 LRF 構文検査 PASS（selected records 42件、total_bytes 16607）。
- ネイティブツール（`llm-mem status`, `crabcode --version`, `rg`）の疎通確認。
