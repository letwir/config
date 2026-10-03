# Walkthrough — H-12 nested subagent handoff probe

## Overview

Tested whether a dispatched Codex subagent could dispatch one child using the current SUBAGENTS contract. Root-to-parent dispatch succeeded. The parent read the routed rules and prepared an exact child wrapper, but the runtime rejected non-root child dispatch with a root-thread-only boundary. No child executed, so no independent child result is claimed.

## Contract update

- `subagent.invocation-codec` now applies to each authorized handoff author, including explicitly authorized nested subagents; persona CSS remains presentation-only.
- `subagent.nested-handoff` requires task authorization, `agent.mode` permission, depth at most two, and a runtime delegation tool available to the current actor.
- Root-only or unavailable delegation returns FAILED/ESCALATE with the prepared prompt and blocker. Thread creation is not a proxy for child dispatch.
- `subagent.etl-context` binds both rules to the exact wrapper and ordered handoff fields.

## Verification

- `harness-lint.exe -path agents/SUBAGENTS.lrf -level 1`: PASS.
- Level 2 reports only the existing `subagent.mode` effect mismatch at line 53.
- Targeted `git diff --check`: PASS.
- Runtime evidence: root-to-parent succeeded; parent-to-child was rejected as root-only. No child result exists.

## Uncertainty

Nested Codex delegation is not available from the tested non-root runtime. This probe does not establish whether another runtime exposes nested dispatch.
