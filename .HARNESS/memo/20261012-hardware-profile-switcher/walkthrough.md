# walkthrough / 20261012-hardware-profile-switcher

依頼: Pwsh・WindowsとBashの環境をHARDWARE.mdで切り替えできる記法にする。

変更内容:
1. `~/.harness/HARDWARE.md`: 冒頭に `active_profile` セレクター（`android-termux-bash` (既定), `windows-pwsh`, `linux-native-bash`）を導入し、プロファイル別のハードウェア・OS・シェル・パス規約・禁止事項を独立定義。さらにプロファイル別コマンド対照表（Cross-Profile Command Mapping）を新設。
2. `~/.harness/AGENTS.md`: `HARDWARE.md` の `active_profile` を参照し、宣言されたプロファイルのシェル・パス体系・禁止事項に適合するコマンドを実行するよう誘導。
3. `~/.harness/rules/BOOTSTRAP.lrf`: `R|host|*` を静的定義から `hardware=../HARDWARE.md; resolve=active_profile-in-HARDWARE.md; adapt-shell-and-paths=active_profile; default-profile=android-termux-bash` に更新。
4. `~/.harness/rules/README-forLLM.md`: `active_profile` の解決を指示し、プロファイル別のコマンド（Bash / PowerShell）を併記。
5. `~/.harness/skills/default-settings/SKILL.md` (および `SKILL_default.md`), `skills/agy-subagent-router/SKILL.md`, `skills/economics-api/USAGE.md`, `skills/rikei-search/USAGE.md`, `method.md`: `active_profile` に連動して Bash と PowerShell の記法を切り替える方式へ改修。
6. `~/.harness/issues.md`: H-12 として環境プロファイル切り替え記法の導入を記録。

検証:
- `bash scripts/invoke-rule-preflight.sh --task change`: 終了コード 0、全 LRF 構文検査 PASS（selected records 42件、total_bytes 16554）。
