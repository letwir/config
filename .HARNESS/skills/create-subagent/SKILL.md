---
name: create-subagent
description: Create custom subagents for specialized AI tasks (code reviewers, debuggers, domain assistants).
disable-model-invocation: true
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max] 🤖agent 📂scope 🎯dispatch
🤖agent ⊢ Name ∧ Prompt → custom-specialist `<name>.md` {name,description,system prompt}; examples={code reviewer,debugger,domain assistant}
📂scope ⊢ target∈{Project,User} → Project:`.cursor/agents/` (priority↑) | User:`~/.cursor/agents/`
🎯dispatch ⊢ Description → explicit, concrete delegation-trigger text (AI can decide when to delegate)

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-subagent
R|use.create|task:create-subagent|MUST|LW_SCOPE|target=agents/<name>.md; format=yaml-frontmatter+markdown
R|use.naming|task:create-subagent|MUST|RO_LOCAL|name=lowercase-hyphens-only; description=concrete-trigger-defined
R|guard.prompt|task:create-subagent|MUST|RO_LOCAL|prompt=role-boundary-bounded; context-isolation=max
