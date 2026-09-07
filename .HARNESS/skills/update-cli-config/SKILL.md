---
name: update-cli-config
description: View and modify Codex CLI configuration settings in ~/.cursor/cli-config.json.
metadata: { surfaces: [cli] }
---
[SPR/PIDGIN::ρ→max] ⚙️config 🛡️perms 🎛️mode 🚫guard
⚙️ 設定 config ⊢ Path ⇒ CLI設定の閲覧・編集 | `~/.cursor/cli-config.json` ∨ `.cursor/cli.json`
🛡️ 権限 perms ⊢ Allow ∧ Deny ⇒ ツール実行権限（permissions.allow/deny）の更新
🎛️ モード mode ⊢ Approval ∧ Vim ⇒ `approvalMode`（allowlist/unrestricted）や `editor.vimMode` 切り替え
🚫 保護 guard ⊢ InternalFields ⇒ `version`, `model`, `authInfo`, `privacyCache` の手動変更を厳禁

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=update-cli-config
R|use.read|task:update-cli-config|MUST|RO_LOCAL|target=~/.cursor/cli-config.json; parse=json
R|use.write|task:update-cli-config|MUST|LW_SCOPE|target=~/.cursor/cli-config.json; restart-required=true
R|guard.protect|task:update-cli-config|MUST_NOT|LW_SCOPE|modify=version,model,selectedModel,modelParameters,privacyCache,authInfo
