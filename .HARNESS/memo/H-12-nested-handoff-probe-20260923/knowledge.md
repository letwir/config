# Knowledge — H-12 nested subagent handoff probe

- `SUBAGENTS.lrf` is the source for subagent handoff format and role routing. Persona CSS is only presentation.
- The root agent successfully invoked a parent subagent. The parent prepared a child prompt using the exact `<Γ>` wrapper and fields, but the runtime rejected child dispatch from the non-root context as root-thread-only.
- No nested child ran and no child findings exist. Record this as failed/unavailable, never as successful verification.
- Nested delegation requires current-task authorization, `agent.mode` permission, depth at most two, and a delegation tool exposed to the current actor. If unavailable, return FAILED/ESCALATE with the prepared prompt and blocker.
- Validation: SUBAGENTS LRF Level 1 passed; Level 2 retains one unrelated existing `subagent.mode` effect mismatch; targeted diff check passed.
