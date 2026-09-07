---
name: create-rule
description: Create Codex rules for persistent AI guidance. Use when setting up project conventions or .cursor/rules/*.mdc.
---
[SPR/PIDGIN::ρ→max] 📜rule 🎯scope 📐standards
📜 ルール rule ⊢ Globs ∧ Standards ⇒ プロジェクト規約定義 | `.cursor/rules/<name>.mdc`
🎯 適用 scope ⊢ Always ∨ Glob ⇒ 常時適用（`alwaysApply: true`）またはファイルパターン一致（`globs: "**/*.ts"`）
📐 規約 standards ⊢ StyleGuideline ⇒ 言語規約・設計パターン・必須規約の明確化

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-rule
R|use.create|task:create-rule|MUST|LW_SCOPE|target=.cursor/rules/*.mdc; format=yaml-frontmatter+markdown
R|use.scope|task:create-rule|MUST|RO_LOCAL|clarify=file-patterns-or-always-apply; ask-if-underspecified=true
R|guard.content|task:create-rule|MUST|RO_LOCAL|concrete-examples=preferred; vague-instruction=forbidden
