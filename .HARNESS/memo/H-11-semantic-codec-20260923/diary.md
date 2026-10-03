# Diary — H-11 semantic codec instructions

### 2026-09-23

- **Request:** Add the supplied semantic compression rule to intermediate interpretation and update subagent instructions.
- **Correction:** 旦那様 clarified that CSS is only a thin persona layer and the codec belongs in the main-to-subagent invocation prompt.
- **Action:** Removed the codec from `persona/PERSONA.css`; added `subagent.invocation-codec` to `agents/SUBAGENTS.lrf` and bound it to `subagent.etl-context` before dispatch. Preserved the handoff schema and policy boundaries.
- **Nested probe:** Root dispatched a parent subagent. It prepared a correctly wrapped child prompt, but could not invoke the child because nested collaboration dispatch was unavailable/root-only in that runtime context. No child ran. Added an explicit fail-closed nested-dispatch contract to `SUBAGENTS.lrf`.
- **Result:** Level 1 LRF lint and targeted whitespace/diff check passed. Level 2 lint reports only the existing effect mismatch in `subagent.mode`; it was left unchanged.
- **PromptDefect:** 0% — clarification resolved placement.
- **AgentDefect:** Initial placement in persona CSS misread the intended invocation point; corrected after clarification. The requested nested runtime execution could not complete because child dispatch was root-only/unavailable.
- **Uncertainty:** Runtime adoption was not tested.
- **Impact:** Main and subagent prompt guidance now use the same compact encoding objective.
