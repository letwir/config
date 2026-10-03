# Diary: OpenCode GUI Orchestrator permissions

### 2026-09-26
- task_id: 20260926-opencode-orchestrator-permissions
- target environment: OpenCode GUI; evidence/source scope: global agent file, official docs, harness policy.
- Hypothesis: enabling the agent's `bash`, read/discovery, and `task` permissions is the relevant configuration change.
- Tried: inspected current agent and official permissions/agents documentation; obtained Y confirmation; changed the single global agent definition and independently read the result.
- Rejected: assuming a CLI installation or claiming nested Task dispatch without a runtime check.
- Uncertainty: GUI runtime nesting capability; the first confirmation prompt was unanswered, and the user later answered Y.
- Attribution: PromptDefect=none identified; AgentDefect=none identified.
- Search: memory precedent query returned no useful record.
- Correction: user clarified that the installed OpenCode is GUI-only; replaced CLI validation with static checks.
- Emotion: neutral. Thoughts: use static validation rather than a missing CLI.
- Plan: completed targeted configuration change; verify in GUI after restart.
- Verified knowledge: `bash: deny` presently blocks processes, while `task: allow` is already present.
- Method: bounded read and official documentation review.
- Verification status: target read and static structure checked; no live GUI test. Remaining uncertainty: runtime nested dispatch.
- Completion evidence: target now allows file exploration, external process execution, and Task with depth-bounded prompt.
- Impact: GUI Orchestrator permission change requires restart. Tags: opencode,gui,permissions.
