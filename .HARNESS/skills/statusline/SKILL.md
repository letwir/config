---
name: statusline
description: Configure a custom status line in the CLI. Use when configuring CLI status bar, prompt footer, or adding session context above the prompt.
---
[SPR/PIDGIN::ρ→max] 📊config 📥payload 🧪test
📊 設定 config ⊢ ConfigJson ⇒ CLIステータスライン定義 | `~/.cursor/cli-config.json` の `statusLine`（type=command, command, padding）
📥 ペイロード payload ⊢ StdinJSON ⇒ stdinセッション情報（session, model, context_window, cwd, vim, worktree）受信
🧪 テスト test ⊢ MockJSON ⇒ モックJSONパイプによるスクリプト動作検証 | `echo '{"model":{"display_name":"Opus"}}' | ./statusline.sh`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=statusline
R|use.config|task:statusline|MUST|LW_SCOPE|target=~/.cursor/cli-config.json; key=statusLine; type=command
R|use.render|task:statusline|MUST|RO_LOCAL|stdout=ansi-color-supported; multi-line=row-separated; token-cost=zero
R|guard.timeout|task:statusline|MUST|RO_LOCAL|timeoutMs<=2000; debounce>=300ms; abort-on-new-update=true
