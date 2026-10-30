---
name: migrate-to-skills
description: Convert intelligent rules (.cursor/rules/*.mdc) and slash commands (.cursor/commands/*.md) to Agent Skills format (.cursor/skills/).
disable-model-invocation: true
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max|task:migrate-to-skills|map:input→operation→result]
🔎 Trigger: explicit migration request for `.cursor/rules/*.mdc`, `.cursor/commands/*.md`, or a requested rollback from migration backups.
⊢ Scope: only those requested source rule/command files and their corresponding `.cursor/skills/<name>/SKILL.md` outputs; no unrelated files.
📜 Rule: `.mdc` ⇒ add `name`/`description`, remove `globs`; condition=has-description-and-no-globs-and-not-alwaysApply; preserve body verbatim, without reformatting.
⚡ Command: `.md` ⇒ add `disable-model-invocation: true`; preserve body verbatim, without reformatting.
🔄 Rollback: restore from migration-before backup; delete originals only after output verification; retain required undo guidance.
⛔ Failure/ambiguity: do not guess, broaden scope, or delete unverified originals; stop and report the unresolved input or failed verification.
✅ Success: verified output at `.cursor/skills/<name>/SKILL.md` with required frontmatter and unchanged body; rollback restores the pre-migration state.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=migrate-to-skills
R|use.migrate|task:migrate-to-skills|MUST|LW_SCOPE|preserve-body=verbatim-strict; do-not-reformat=true
R|use.rule|task:migrate-to-skills&case:rule|MUST|LW_SCOPE|condition=has-description-and-no-globs-and-not-alwaysApply
R|use.command|task:migrate-to-skills&case:command|MUST|LW_SCOPE|inject-frontmatter=disable-model-invocation:true
R|guard.destruct|task:migrate-to-skills|MUST|LW_SCOPE|delete-original=only-after-verified; undo-guidance=required
