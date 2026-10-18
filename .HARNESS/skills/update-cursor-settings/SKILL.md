---
name: update-cursor-settings
description: Modify Codex/VSCode user settings in settings.json (themes, font size, format on save).
metadata: { surfaces: [ide] }
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[PG] Trigger: user requests update of Cursor/Codex editor settings. Apply only the settings explicitly requested; supported keys: `workbench.colorTheme`, `editor.fontSize`, `editor.tabSize`, `editor.formatOnSave`, `files.autoSave`, `files.exclude`. Choose one valid target scope: user/global `settings.json` or workspace `.vscode/settings.json`; do not conflate them. User-settings paths: Windows `%APPDATA%\Codex\User\settings.json`; macOS `~/Library/Application Support/Codex/User/settings.json`. Parse settings as JSONC (comments permitted), preserve existing content, and write with 2-space indentation. Validate the resulting JSONC and confirm requested keys/values and target scope before notifying that reload is needed. Ambiguous scope, requested setting/value, or validation failure ⇒ do not write; retain existing settings and ask for clarification. [/PG]

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=update-cursor-settings
R|use.read|task:update-cursor-settings|MUST|RO_LOCAL|target=settings.json; parse=jsonc; handle-comments=true
R|use.write|task:update-cursor-settings|MUST|LW_SCOPE|indent=2-spaces; preserve-existing=true; reload-notify=true
R|guard.scope|task:update-cursor-settings|MUST|RO_LOCAL|target-separation: user-settings(global) vs workspace-settings(.vscode/settings.json)
