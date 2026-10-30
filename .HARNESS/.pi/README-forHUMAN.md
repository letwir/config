# Pi 設定の配布・導入

このフォルダーは、Pi のユーザー設定のうち再利用可能なものを配布します。設定とエージェントはこのリポジトリを参照するシンボリックリンクとして導入します。**実行前に内容を確認してください。**

## 導入

Bash:

```bash
bash install-pi.sh --dry-run
bash install-pi.sh
```

PowerShell 7:

```powershell
.\install-pi.ps1 -DryRun
.\install-pi.ps1
```

既定のリンク先は `~/.pi/agent` です。別の場所には Bash の `--agent-dir DIR`、PowerShell の `-AgentDir DIR` を指定できます。既存ファイル等が衝突した場合、既定では**何も変更せず中止**します。`--backup` / `-Backup` を明示すると、衝突対象を同じ親の一意な `.backup.*` 名へ移動してからリンクします。バックアップは上書き・削除しません。再実行は正しいリンクを保つよう設計されています。

導入物は `agent/settings.json` と選択した拡張・agent/resource エントリへの個別リンクです。`extensions` 全体、`agents` 全体は置き換えません。PowerShell版はジャンクションではなく実シンボリックリンクを作成します。Windowsでは開発者モードまたは管理者権限が必要な場合があります。

## 配布・依存関係

共有 Harness ファイルは、この配布物の親 `.HARNESS` を参照します。共有 contracts は `agent/agents/` に個別リンクし、prompt は `agent/prompts` にリンクします。`~/.harness` はソースではありません。共有ルールやプロンプトが `~/.harness` 等のユーザー環境資源を参照する場合、その依存先は別途用意してください。導入スクリプトはダウンロード、npm/システムパッケージ導入、Pi 実行ファイルや認証情報の移動をしません。

Pi の package 宣言は `agent/settings.json` に保持しています。確認後、必要な場合に限って別途 `pi update --extensions` を実行してください。実行中の Pi が設定を書き換える可能性があります。`defaultProjectTrust: always` は既存値を維持しています。信頼できないプロジェクトでは危険になり得るため、ユーザー自身で判断してください。

除外対象: Pi バイナリ、npm 実体、認証情報、モデルストア、セッションや実行時データ。導入・テストで普段使いの `~/.pi` を変更しないでください。
