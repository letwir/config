# walkthrough / 20261012-config-sync-and-install-guide

## 概要
ユーザーの依頼に基づき、以下を実施した：
1. `~/.harness/skills/` 内の全スキルのバイナリ解決パスおよびコマンド記法を `.harness` 内部統一形式・スラッシュ `/` に統一。
2. 移植・移行ガイド `INSTALL.md`（①:今の環境を探索。 ②:active_profileを切り替え）を作成。
3. `~/.harness` の全変更をリポジトリ原本 `repo/config/.HARNESS` に同期・反映。
4. `repo/config/.HARNESS` 内のパス区切りをスラッシュ `/` に統一。
5. 両ディレクトリで LRF プレフライト検査を実行し PASS を確認。

## 変更ファイル一覧
- `~/.harness/INSTALL.md` (新規)
- `repo/config/.HARNESS/HARDWARE.md` (新規反映)
- `repo/config/.HARNESS/INSTALL.md` (新規反映)
- `repo/config/.HARNESS/scripts/invoke-rule-preflight.sh` (新規反映)
- `repo/config/.HARNESS/scripts/sync-harness.sh` (新規反映)
- `repo/config/.HARNESS/scripts/sync-harness.test.sh` (新規反映)
- `repo/config/.HARNESS/skills/google-search/` (新規反映)
- `repo/config/.HARNESS/AGENTS.md` (同期)
- `repo/config/.HARNESS/issues.md` (同期)
- `repo/config/.HARNESS/method.md` (同期)
- `repo/config/.HARNESS/rules/` 配下全ルールファイル (同期・スラッシュ統一)
- `repo/config/.HARNESS/skills/` 配下全スキル定義 (同期・パス統一)

## 検証結果
- `bash ~/.harness/scripts/invoke-rule-preflight.sh --task change` -> `status: PASS`
- `bash ~/repo/config/.HARNESS/scripts/invoke-rule-preflight.sh --task change` -> `status: PASS`
- `git -C ~/repo/config status` -> クリーンなステージング前差分を確認（不要なバイナリ・バックアップなし）
