---
name: update-cli-config
description: View and modify Codex CLI configuration settings in ~/.cursor/cli-config.json.
metadata: { surfaces: [cli] }
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

⚙️ CLI config ⊢ request inspect∨modify⇒inspect `~/.cursor/cli-config.json`; JSON.
✎ Modify ⊢ explicit request ∧ applicable approval ⇒ only `permissions.allow`, `permissions.deny`, `approvalMode` (`allowlist`∨`unrestricted`), `editor.vimMode`; no other fields; no unapproved configuration changes.
🚫 `version`, `model`, `selectedModel`, `modelParameters`, `privacyCache`, `authInfo` ⇒ MUST_NOT modify manually.
✅ Validate JSON after approved edit; invalid⇒STOP, report; do not proceed as if changed. Restart required after write.
#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=update-cli-config
R|use.read|task:update-cli-config|MUST|RO_LOCAL|target=~/.cursor/cli-config.json; parse=json
R|use.write|task:update-cli-config|MUST|LW_SCOPE|target=~/.cursor/cli-config.json; restart-required=true
R|guard.protect|task:update-cli-config|MUST_NOT|LW_SCOPE|modify=version,model,selectedModel,modelParameters,privacyCache,authInfo
