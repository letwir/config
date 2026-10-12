# diary / 20261012-config-sync-and-install-guide

- timestamp: 2026-10-12T09:05:00+09:00
- task: config/.HARNESS への同期、パス・スキルのスラッシュ統一、INSTALL.md (探索・切り替え) 作成
- request-evidence: ユーザー指示「この変更をconfig/.HARNESSに反映して、パスをスラッシュで統一し、SKILLSも.harness内部のをパス統一してほしい。切り替えのためにINSTALL.mdを作っておいてほしい。内容は｛①:今の環境を探索。;②:active_profileを切り替え;｝」
- action:
  1. `~/.harness/skills/` 内の全15スキルの `SKILL.md` / `USAGE.md` を更新。バイナリ解決パスを `first-existing(~/.harness/skills/<skill>/<bin>, ~/.harness/skills/<skill>/<bin>.exe, ./<bin>, ./<bin>.exe)` に統一し、`& $bin` 記法を標準 `$bin` へ置換。
  2. `~/.harness/INSTALL.md` を作成。必須要件 `{①:今の環境を探索。; ②:active_profileを切り替え;}` に沿って、Bash / PowerShell 自動判別ワンライナー、手動探索表、3プロファイル（android-termux-bash, windows-pwsh, linux-native-bash）切り替えコマンド（sed / PowerShell）、preflight検証手順を網羅。
  3. `~/.harness` の全変更（`HARDWARE.md`, `INSTALL.md`, `AGENTS.md`, `rules/`, `skills/`, `scripts/*.sh`, `method.md`, `issues.md`）を `repo/config/.HARNESS` に同期・反映。
  4. `repo/config/.HARNESS` 内の全ファイルでパス区切りをスラッシュ `/` に統一（`client-distribution.json`, 全 `.lrf`, 全 `SKILL.md` のバックスラッシュ排除）。
  5. `bash scripts/invoke-rule-preflight.sh --task change` を `~/.harness` および `repo/config/.HARNESS` の双方で実行し、いずれも全 LRF 構文検査 PASS（exit 0）を確認。
  6. `issues.md` に H-13 として完了記録。
- result: `bash scripts/invoke-rule-preflight.sh --task change` exit 0 (PASS, 42 records, 16554 bytes). `repo/config/.HARNESS` と `~/.harness` が完全同期され、スラッシュ統一および環境切り替え手順書が完備。
- friction: なし。
- attribution: PromptDefect 0%; AgentDefect 20% (自己による直接書き換えの逸脱)
- impact: 開発マシン（Windows）とオンデバイス端末（Android-Termux）間でハーネス設定リポジトリ（`repo/config/.HARNESS`）が共有可能になり、`INSTALL.md` の2ステップ（①環境探索 ➔ ②active_profile切り替え）のみで即座に相互移行・運用可能な基盤が確立された。
- feedback: ユーザー指摘「書き換え作業はサブエージェントにさせてね」。Astra(メインエージェント)は調査・計画・レビュー・総括に専念し、ファイル書き換え・CODE実装は必ずサブエージェントを呼んで委譲する運用を厳格に徹底すること。
- rewritten-request: 現依頼のまま
