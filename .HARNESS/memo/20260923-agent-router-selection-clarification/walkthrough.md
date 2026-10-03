# walkthrough | 20260923-agent-router-selection-clarification

## Updated design
- Do not add `agents/selection.lrf` as a startup-key or role-selection file.
- `agents/AGENT_ROUTER.md` contains role/model/TOML selection constraints and chooses one base persona CSS according to task needs and difficulty; CSS is not rigidly mapped to the agent name.
- `etl/main.seq` invokes AGENT_ROUTER.md only at agent dispatch points; workflow stages remain in main.seq.
- `SUBAGENTS.lrf` and MANUAL LRF retain security, authority, effects, handoff, and CSS-composition constraints. They explicitly load and bind the router for selection, and remove conflicting fixed role-to-CSS rules.
- Keep current CODE execution boundary (direct agy) unless user changes it; resolve worker.toml against that boundary before implementation.
- No operational edits or tests were performed.
