---
name: create-subagent
description: Create custom subagents for specialized AI tasks (code reviewers, debuggers, domain assistants).
disable-model-invocation: true
---
[SPR/PIDGIN::ρ→max] 🤖agent 📂scope 🎯dispatch
🤖 サブエージェント agent ⊢ Name ∧ Prompt ⇒ 特化型サブエージェント作成 | `<name>.md`（name, description, system prompt）
📂 スコープ scope ⊢ Project ∨ User ⇒ プロジェクト用 `.cursor/agents/`（高優先） ∨ ユーザー用 `~/.cursor/agents/`
🎯 委譲 dispatch ⊢ Description ⇒ 委譲トリガーを明確化（AIが委譲判断できる明確な説明文）

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-subagent
R|use.create|task:create-subagent|MUST|LW_SCOPE|target=agents/<name>.md; format=yaml-frontmatter+markdown
R|use.naming|task:create-subagent|MUST|RO_LOCAL|name=lowercase-hyphens-only; description=concrete-trigger-defined
R|guard.prompt|task:create-subagent|MUST|RO_LOCAL|prompt=role-boundary-bounded; context-isolation=max
