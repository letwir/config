# Knowledge Base: OpenCode GUI Orchestrator permissions

- task_id: 20260926-opencode-orchestrator-permissions
- timestamp: 2026-09-26
- target environment: OpenCode GUI on Windows
- evidence/source scope: global `~/.config/opencode/agents/orchestrator.md`, OpenCode official agents and permissions documentation (2026-09-25), harness policy files.
- verification status: Y confirmed; global agent edited and independently inspected; live GUI invocation not tested.
- remaining uncertainty: whether this GUI build exposes nested Task calls at runtime; no CLI is installed.

## Context
The global Orchestrator agent is a subagent with `edit: deny`, `bash: deny`, and `task: allow`. Its prompt enforces depth 0–2 and checks runtime capability.

## Findings
OpenCode documents `bash` for external process execution, `read`/`glob`/`grep`/`list` for file discovery, and `task` for subagent invocation. Agent-specific permissions can be specified in Markdown frontmatter; Task permission is matched against the subagent type. External-directory access is a separate permission.

## Morphism
Confirmed scope and documented permission keys → narrowly updated agent permissions and delegation instruction → independent static inspection; runtime nesting remains to be verified in the GUI.

## Attribution
- Plan: after Y confirmation, edit only the global agent definition and inspect frontmatter/body and scope.
- Verified knowledge: current permission block and official permission keys.
- Method: direct file inspection plus official documentation.
- Prompt defect: none identified. Agent defect: none identified.
- Completion evidence: target frontmatter now explicitly has `read`, `glob`, `grep`, `list`, `bash`, and `task` allowed with `edit` denied; body retains LRF gates and depth<=2.
- Impact: OpenCode GUI Orchestrator permissions broadened on next restart; nested dispatch capability remains unproven.
- Tags: opencode,gui,orchestrator,permissions,nested-task
