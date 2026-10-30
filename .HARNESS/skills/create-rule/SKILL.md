---
name: create-rule
description: Create Codex rules for persistent AI guidance. Use when setting up project conventions or .cursor/rules/*.mdc.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true
[POLYGLYPH::ρ→max] 📜rule 🎯scope 📐standards
🔎 RuleAuthoring(task) ∧ Globs ∧ Standards ⊢ define_project_guidance → `.cursor/rules/<name>.mdc`
🎯 Scope ∈ {Always(`alwaysApply: true`), matching_files(`globs: "**/*.ts"`)} → specify_when_rule_applies
📐 Standards ∈ {style_guidelines, design_patterns, required_conventions} → make_project_expectations_explicit

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-rule
R|use.create|task:create-rule|MUST|LW_SCOPE|target=.cursor/rules/*.mdc; format=yaml-frontmatter+markdown
R|use.scope|task:create-rule|MUST|RO_LOCAL|clarify=file-patterns-or-always-apply; ask-if-underspecified=true
R|guard.content|task:create-rule|MUST|RO_LOCAL|concrete-examples=preferred; vague-instruction=forbidden
