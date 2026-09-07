---
name: create-skill
description: Create Codex Agent Skills. Use when authoring a new skill or asking about SKILL.md structure.
---
[SPR/PIDGIN::ρ→max] 🛠️create 📂location 📜structure
🛠️ 作成 create ⊢ Spec ∧ Scope ⇒ Agent Skill生成 | `<skill-name>/SKILL.md`（YAML frontmatter + 本文）
📂 配置 location ⊢ Personal ∨ Project ⇒ 個人用 `~/.cursor/skills/` または プロジェクト用 `.cursor/skills/`
📜 構造 structure ⊢ Layout ⇒ ディレクトリ構成（`SKILL.md`, `reference.md`, `examples.md`, `scripts/`）

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-skill
R|use.create|task:create-skill|MUST|LW_SCOPE|target=skills/<name>/SKILL.md; yaml-frontmatter=name,description
R|use.content|task:create-skill|MUST|RO_LOCAL|user-verbatim=preserve-strictly; paraphrase-user-copy=forbidden
R|guard.discovery|task:create-skill|MUST|RO_LOCAL|description=trigger-scenario-rich; name=kebab-case
