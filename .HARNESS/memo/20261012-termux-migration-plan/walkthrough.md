# walkthrough / 20261012-termux-migration-plan

依頼: issues.mdに移行計画を記載し、別モデルのサブエージェントを利用、PowerShell手順をbashに読み替える。
変更: ~/.harness/issues.mdにH-09/H-10を追記。SSOT実体、削除なしのステージング、バックアップ/復元、crabcode優先、bash版検証、短いhandoffと1 writerを定義。
検証: 見出し各1件、bash版preflightと独立レビュー未成立の明記、末尾空白なし。別モデル起動2件は認証制限で失敗。移行・設定切替・VCS操作なし。
残余: 実効設定パス・移行fixture・実装・別モデル疎通は未検証。

更新 (2026-10-12 07:32): 設定探索順序を公式docsで確定（layer2=`~/.config/crabcode`、XDG未設定）。`~/.harness/crabcode` を実体として `~/.config/crabcode` を symlink 化。バックアップ+manifest+復元試験を `~/backup/harness-20261012-072322` に実施、切替前後ハッシュ一致、新instance起動 exit 0。openai/gpt-6-luna を `--model` 明示で疎通・worker委譲に成功。`scripts/invoke-rule-preflight.sh`（bash版preflight）を実装させ、正常系PASS/異常系6種の非0終了を主担当がfixtureで検証。残: 他 `.ps1` 15ファイルのbash移植、誤リンク・既存ディレクトリ・空白パスのfixture。

更新 (2026-10-12 10:58): 残作業は完了済みであるため H-09 手順4・6 を [x] に更新。クローズ検証として bash -n 9/件全通過、sync-harness.test.sh 14アサーション PASS、--consolidate 再実行冪等（exit 0）、preflight exit 0、~/.harness と repo/config/.HARNESS の diff ゼロを確認し、issues.md へクローズ検証を追記した。
