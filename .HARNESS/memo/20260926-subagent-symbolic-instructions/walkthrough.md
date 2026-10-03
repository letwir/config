# Walkthrough: symbolic subagent instructions

## 実行計画

1. Compile the active bootstrap rules.
2. Inspect the applicable subagent handoff and symbolic-codec contracts.
3. Preserve required handoff fields and wrapper constraints while recording the requested task-payload style.
4. Report the one unresolved boundary between the payload suffix and wrapper close.

## Verification

- Read `rules/BOOTSTRAP.lrf`, `rules/MANUAL.lrf`, `rules/LOAD.lrf`, `agents/SUBAGENTS.lrf`, and `agents/AGENT_ROUTER.md`.
- No subagent invocation or workspace policy edit occurred.
- Remaining uncertainty: whether `final_line≡SIGMA/1` applies to the compressed payload or to the whole wrapped invocation.
