---
name: create-hook
description: Create Codex hooks. Use when configuring hooks.json, hook scripts, or automating behavior around agent events.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true
[SPR/PIDGIN::ρ→max] 🪝hook ⚡event 🛡️policy
🔎 HookCreation(task) ∧ Scope ∧ Trigger ⊢ define_hook+generate_script → `.cursor/hooks.json` ∨ `~/.cursor/hooks.json`
⚡ LifecycleEvent ∈ {sessionStart, preToolUse, postToolUse, beforeShellExecution, ...} ⊢ select_trigger → hook_fires_at_matching_event
🛡️ HookPolicy ∈ {Deny, Rewrite} → block_tool_execution ∨ rewrite_input ∨ inject_context_dynamically

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=create-hook
R|use.create|task:create-hook|MUST|LW_SCOPE|target=hooks.json+script; json-valid=true; fail-safe=explicit
R|use.event|task:create-hook|MUST|RO_LOCAL|scope=narrowest-matching-event; command-or-prompt=deterministic
R|guard.security|task:create-hook|MUST|RO_LOCAL|stdin-stdout-json=typed; infinite-loop-or-orphan-process=forbidden
