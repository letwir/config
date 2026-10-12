# diary / 20261012-hardware-spec-and-bash-guidance

- timestamp: 2026-10-12T07:48:30+09:00
- task: ホスト環境仕様の記録 (HARDWARE.md) と各地コマンド実行の Bash 誘導
- request-evidence: ユーザー指示「君はこの環境を.harness/HARDWARE.mdに記録しつつ、それを元に各地のコマンド実行方法をHARDWARE.mdを参照して環境に合わせたシェルを使うように誘導してくれ。」
- action: 
  1. ホスト環境調査（Android 16, aarch64, Sony XQ-GE44 / Snapdragon 8 Gen 3, 16GB RAM, Termux 0.118.3, 既定 GNU Bash, ネイティブツール群）。
  2. `.harness/HARDWARE.md` の作成（仕様・パス体系・tmpdir注意・旧PowerShell→新Bash移行対応表のSSOT化）。
  3. `AGENTS.md`, `rules/BOOTSTRAP.lrf`, `rules/README-forLLM.md`, `rules/MANUAL.lrf`, `rules/engineering.lrf`, `skills/default-settings/SKILL.md` (および `SKILL_default.md`), `skills/agy-subagent-router/SKILL.md`, `skills/economics-api/USAGE.md`, `skills/rikei-search/USAGE.md`, `method.md` のシェル・コマンド例を Bash / `HARDWARE.md` 参照へ更新。
  4. `issues.md` に H-11 として完了記録。
- result: `bash scripts/invoke-rule-preflight.sh --task change` exit 0 (PASS, 42 records, 16607 bytes). 全LRF整合性確認済み。
- friction: なし。ホスト情報取得および各ファイル編集はすべて正常に完了。
- attribution: PromptDefect 0%; AgentDefect 0%
- impact: 各エージェントおよびサブエージェントがコマンドを実行する際に `HARDWARE.md` を参照し、環境に適応した Bash を確実に選択するよう統一的に誘導された。
- feedback: なし。指示は極めて明確であり、期待された仕様正本化と誘導を完全に実現。
- rewritten-request: 現依頼のまま
