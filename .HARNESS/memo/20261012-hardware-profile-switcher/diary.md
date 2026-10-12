# diary / 20261012-hardware-profile-switcher

- timestamp: 2026-10-12T08:52:30+09:00
- task: HARDWARE.md における環境プロファイル切り替え記法の導入
- request-evidence: ユーザー指示「今回の変更、Pwsh・WindowsとBashの環境をHARDWARE.mdで切り替えできる記法にできないかな」
- action:
  1. `HARDWARE.md` に `active_profile` セレクター（`android-termux-bash`, `windows-pwsh`, `linux-native-bash`）、各プロファイル定義、およびプロファイル別コマンド対照表（Cross-Profile Command Mapping）を導入。
  2. `AGENTS.md`, `rules/BOOTSTRAP.lrf`, `rules/README-forLLM.md`, `method.md`, 各スキル（`default-settings`, `agy-subagent-router`, `economics-api`, `rikei-search`）を `active_profile` に連動してシェル・パス記法を切り替える記法へ更新。
  3. `issues.md` に H-12 として完了記録。
- result: `bash scripts/invoke-rule-preflight.sh --task change` exit 0 (PASS, 42 records, 16554 bytes). 全 LRF ルールおよびドキュメントで切り替え可能記法が正常に機能することを確認。
- friction: なし。
- attribution: PromptDefect 0%; AgentDefect 0%
- impact: `HARDWARE.md` の `active_profile` を変更するだけで、Windows/PowerShell と Android-Termux/Bash（および Linux/Bash）の環境追従をハーネス全体で即座に切り替えられる柔軟な構成が実現された。
- feedback: なし。指示の意図に完璧に沿った設計を実装。
- rewritten-request: 現依頼のまま
