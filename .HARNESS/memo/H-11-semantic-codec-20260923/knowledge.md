# Knowledge — H-11 semantic codec instructions

- The codec is an invocation-prompt construction rule for every authorized handoff author in `agents/SUBAGENTS.lrf`; it is bound to `subagent.etl-context`. Persona CSS remains a thin presentation-only layer.
- The objective minimizes token count subject to approximate semantic and pragmatic preservation and ambiguity reduction; vocabulary combines five-word language, up to thirty symbols, and emoji, with relation symbols and affect emoji or scalar cues.
- `H-readability=∅` is explicitly defined as excluding human readability from the optimization objective, not as permission to lose semantic precision.
- The rule applies to invocation payloads, not persona CSS or subagent answer style. Existing LRF records continue to own scope, effects, safety, stop conditions, and handoff requirements.
- A live root-to-parent subagent probe confirmed wrapper preparation, but the child dispatch was rejected from the non-root runtime context; no child executed and no independent child result exists.
- Nested delegation now requires task authorization, `agent.mode` permission, depth at most two, and an available runtime tool. Root-only/unavailable delegation must be reported as FAILED/ESCALATE without thread-creation substitution or fabricated verification.
- Checks: subagent LRF Level 1 lint and targeted `git diff --check` passed. Level 2 reports a pre-existing mismatch in `subagent.mode`.
