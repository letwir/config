# Diary — 20260926-bootstrap-readme-relative-path

- Timestamp: 2026-09-26
- Task: Make the bootstrap navigation reference explicit and repository-relative.
- Request evidence: User clarified that the README is under `rules/` and approved the exact one-line change (Y).
- Action: Updated only the `bootstrap.next` target in `rules/BOOTSTRAP.lrf` to `rules/README-forLLM.md`; ran strict level-1 harness-lint and checked that the target exists.
- Result: Passed. The requested path is explicit; unrelated workspace changes were preserved.
- Friction: Initial inspection looked for the navigation file at the workspace root. The user supplied its `rules/` location; this was an agent lookup error, not an ambiguous request.
- Attribution: PromptDefect 0%; AgentDefect 100% (the initial path lookup did not resolve relative to the stated rules directory). Estimates reflect this bounded task.
- Impact: One extra clarification/turn before the intended edit.
- Feedback: No prompt change needed. Internal correction: when a bootstrap navigation filename is absent at root, inspect its configured rules directory before declaring it missing.
- Rewritten request: In `rules/BOOTSTRAP.lrf`, change the README navigation target to `rules/README-forLLM.md` and preserve other content.
