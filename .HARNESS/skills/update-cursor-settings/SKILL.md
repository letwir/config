---
name: update-cursor-settings
description: Modify Codex/VSCode user settings in settings.json (themes, font size, format on save).
metadata: { surfaces: [ide] }
---
[SPR/PIDGIN::ρ→max] ⚙️settings 🎨theme 📐editor 📁files
⚙️ 設定 settings ⊢ OSPath ⇒ Codex/VSCode設定の読込・更新 | Win: `%APPDATA%\Codex\User\settings.json` / Mac: `~/Library/Application Support/Codex/User/settings.json`
🎨 テーマ theme ⊢ Workbench ⇒ テーマ・配色設定の変更 | `workbench.colorTheme: "..."`
📐 エディタ editor ⊢ EditorConfig ⇒ フォント・インデント・保存時フォーマット設定 | `editor.fontSize`, `editor.tabSize`, `editor.formatOnSave`
📁 ファイル files ⊢ FilesConfig ⇒ 自動保存・除外設定 | `files.autoSave`, `files.exclude`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=update-cursor-settings
R|use.read|task:update-cursor-settings|MUST|RO_LOCAL|target=settings.json; parse=jsonc; handle-comments=true
R|use.write|task:update-cursor-settings|MUST|LW_SCOPE|indent=2-spaces; preserve-existing=true; reload-notify=true
R|guard.scope|task:update-cursor-settings|MUST|RO_LOCAL|target-separation: user-settings(global) vs workspace-settings(.vscode/settings.json)
