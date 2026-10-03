# Walkthrough — H-11 semantic codec instructions

## Overview

Added the requested compact semantic-encoding rule to every authorized dispatcher’s invocation payload contract in the canonical subagent LRF. The codec does not belong in persona CSS. A live parent subagent prepared the required child wrapper and fields, but the runtime rejected child dispatch from the non-root context; no child ran. The LRF now defines this failure boundary and prohibits claiming an unreturned child result.

## Deliverables

- `agents/SUBAGENTS.lrf`: bound `subagent.invocation-codec` to every authorized handoff author, added a nested-handoff gate, and bound both rules to `subagent.etl-context`.
- `persona/PERSONA.css`: retains only the thin presentation persona; the temporary codec addition was removed.

## Verification

- `harness-lint.exe -path agents/SUBAGENTS.lrf -level 1`: PASS.
- Level 2 lint reports only the existing `subagent.mode` effect mismatch; that unrelated rule was left unchanged.
- `git diff --check -- agents/SUBAGENTS.lrf persona/PERSONA.css`: PASS.
- Nested delegation probe: root-to-parent succeeded; parent prepared a child prompt but could not dispatch because the runtime restricted the attempt to the root thread. No child result was available.
- No test suite run.

## Remaining uncertainty

Nested child execution remains unavailable from the tested subagent context; only root-to-subagent dispatch was confirmed.
