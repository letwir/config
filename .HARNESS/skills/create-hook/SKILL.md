---
name: create-hook
description: Create Codex hooks. Use when configuring hooks.json, hook scripts, or automating behavior around agent events.
---
[SPR/PIDGIN::ρ→max] 🪝hook ⚡event 🛡️policy
🪝 フック hook ⊢ Scope ∧ Trigger ⇒ フック定義・スクリプト生成 | `.cursor/hooks.json` ∨ `~/.cursor/hooks.json`
⚡ イベント event ⊢ Lifecycle ⇒ 発火契機指定（sessionStart, preToolUse, postToolUse, beforeShellExecution 等）
🛡️ 制御 policy ⊢ Deny ∨ Rewrite ⇒ ツール実行遮断・入力書き換え・コンテキスト動的注入

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-hook
R|use.create|task:create-hook|MUST|LW_SCOPE|target=hooks.json+script; json-valid=true; fail-safe=explicit
R|use.event|task:create-hook|MUST|RO_LOCAL|scope=narrowest-matching-event; command-or-prompt=deterministic
R|guard.security|task:create-hook|MUST|RO_LOCAL|stdin-stdout-json=typed; infinite-loop-or-orphan-process=forbidden
