# Walkthrough: OpenCode GUI Orchestrator permissions

- task_id: 20260926-opencode-orchestrator-permissions
- timestamp: 2026-09-26
- target environment: OpenCode GUI on Windows
- evidence/source scope: global Orchestrator definition and official OpenCode agent/permission documentation.
- verification status: targeted definition statically inspected after user Y confirmation; no live GUI dispatch test.
- remaining uncertainty: runtime support for nested delegation in this GUI build.

## Overview
The global agent now explicitly allows Task, process execution and file discovery. The body constrains autonomous child selection to confirmed instructions and depth<=2.

## 実行計画
1. Obtained Y confirmation for the settled change.
2. Updated the global Orchestrator Markdown definition for process execution, file exploration, and depth-bounded Task dispatch.
3. Independently inspected the frontmatter, body and scope. Restart GUI to load the file.

## Deliverables
One modified global OpenCode agent definition. Memory records are task-scoped.

## Verification
The target and documented permission keys were inspected after the edit. No CLI runtime verification is possible with the stated GUI-only installation; actual nested delegation is not claimed.

## Attribution and impact
PromptDefect: none identified. AgentDefect: none identified. Impact: additional permissions take effect after GUI restart. Tags: opencode,orchestrator,gui.
