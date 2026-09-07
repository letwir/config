---
name: migrate-to-skills
description: Convert intelligent rules (.cursor/rules/*.mdc) and slash commands (.cursor/commands/*.md) to Agent Skills format (.cursor/skills/).
disable-model-invocation: true
---
[SPR/PIDGIN::ρ→max] 🚚migrate 📜rule2skill ⚡cmd2skill 🔄undo
🚚 移行 migrate ⊢ Rules ∨ Commands ⇒ Agent Skills形式へ一括移行 | `.cursor/skills/<name>/SKILL.md`
📜 ルール rule2skill ⊢ .mdc ⇒ Frontmatter改編（name/description付与、globs削除、本文完全保持）
⚡ コマンド cmd2skill ⊢ .md ⇒ `disable-model-invocation: true` 付与、本文完全保持
🔄 復元 undo ⊢ Backup ⇒ 移行前バックアップからのロールバック復元

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=migrate-to-skills
R|use.migrate|task:migrate-to-skills|MUST|LW_SCOPE|preserve-body=verbatim-strict; do-not-reformat=true
R|use.rule|task:migrate-to-skills&case:rule|MUST|LW_SCOPE|condition=has-description-and-no-globs-and-not-alwaysApply
R|use.command|task:migrate-to-skills&case:command|MUST|LW_SCOPE|inject-frontmatter=disable-model-invocation:true
R|guard.destruct|task:migrate-to-skills|MUST|LW_SCOPE|delete-original=only-after-verified; undo-guidance=required
