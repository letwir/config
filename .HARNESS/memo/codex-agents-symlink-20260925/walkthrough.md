# Walkthrough - .codex AGENTS symlink to .harness

## 概要
旦那様の依頼に従い、`.codex\AGENTS.md` を `.harness\AGENTS.md` へのシンボリックリンクに変更しました。変更前の内容は日付付きバックアップに保存しました。

## 成果物一覧
- `C:\Users\letwir\.codex\AGENTS.md` -> `C:\Users\letwir\.harness\AGENTS.md`
- `C:\Users\letwir\.codex\AGENTS.md.pre-harness-symlink-20260925.bak`

## 検証結果
リンク種別が `SymbolicLink` で、解決先が `.harness\AGENTS.md` と一致することを確認しました。バックアップの SHA-256 は変更前と一致しました。Fresh Codex session での instruction reload は未確認です。
