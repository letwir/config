# Walkthrough: Fallback Route Clarification

## Overview
- Date: 2026-09-25
- Task ID: `20260925-fallback-route-clarification`
- Request: Clarify the worker fallback route and whether AGY can restart file application at a target other than `agents/AGENT_ROUTER.md`.
- Outcome: No dispatch. The existing AGY attempt launched and failed during transactional workspace apply with `EPERM`; current routing rules make post-launch failures terminal and prohibit backend switching or retry/fallback to route around them. Starting the file list elsewhere would still retry the failed run, so it was not attempted.

## Verified Rule
- `rules/agy.lrf`: fallback is limited to a missing executable or OS/CLI startup failure before a model request, with one bounded launch retry; post-launch backend failures are terminal.
- `agents/AGENT_ROUTER.md`: do not switch backends after a launched run fails; retry/fallback is limited to pre-launch failure under the AGY contract.

## Residual Risk
- The previous apply failure leaves current agent-file state requiring inspection before any separately authorized future task.
