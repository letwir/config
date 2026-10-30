---
name: create-skill
description: Create Codex Agent Skills. Use when authoring a new skill or asking about SKILL.md structure.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true
[POLYGLYPH::ρ→max] 🛠️create 📂location 📜structure
🔎 SkillAuthoring(task) ∧ Spec ∧ Scope ⊢ create_agent_skill → `<skill-name>/SKILL.md` (YAML frontmatter + 本文)
📂 Location ∈ {Personal:`~/.cursor/skills/`, Project:`.cursor/skills/`} → place_skill_in_selected_scope
📜 Layout ⊢ organize_skill_files ∈ {`SKILL.md`, `reference.md`, `examples.md`, `scripts/`}

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-skill
R|use.create|task:create-skill|MUST|LW_SCOPE|target=skills/<name>/SKILL.md; yaml-frontmatter=name,description
R|use.content|task:create-skill|MUST|RO_LOCAL|user-verbatim=preserve-strictly; paraphrase-user-copy=forbidden
R|guard.discovery|task:create-skill|MUST|RO_LOCAL|description=trigger-scenario-rich; name=kebab-case
