---
name: statusline
description: Configure a custom status line in the CLI. Use when configuring CLI status bar, prompt footer, or adding session context above the prompt.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

🔎 CLI statusline設定依頼 ⇒ ~/.cursor/cli-config.json の statusLine に限定して扱う; setup/change trigger の追加条件=未定義⇒既存範囲維持、確認要。
⊢ 許容設定: key=statusLine; type=command; values=command,padding（具体値=未指定）。⛔ 他設定・user settings変更は対象外。
📥 stdin JSON→session, model, context_window, cwd, vim, worktree を入力文脈として受信。
→ stdout: ANSI色対応; 複数行=行分割; token-cost=0。MockJSONで動作検査: `echo '{"model":{"display_name":"Opus"}}' | ./statusline.sh`
👤 user-choice trigger/選択肢、および rollback 手順・復元対象=不明⇒創作せず保持; 設定変更前に確認が必要。
✅ timeoutMs<=2000; debounce>=300ms; 新更新時abort=true。

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=statusline
R|use.config|task:statusline|MUST|LW_SCOPE|target=~/.cursor/cli-config.json; key=statusLine; type=command
R|use.render|task:statusline|MUST|RO_LOCAL|stdout=ansi-color-supported; multi-line=row-separated; token-cost=zero
R|guard.timeout|task:statusline|MUST|RO_LOCAL|timeoutMs<=2000; debounce>=300ms; abort-on-new-update=true
