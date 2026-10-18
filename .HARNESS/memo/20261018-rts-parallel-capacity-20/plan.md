# Implementation Plan: RTS 20-descendant capacity

## Decision

Implement the user's first alternative: allow up to 20 direct child workers and at most 20 active descendants total per canonical RTS Runtime instance. A Runtime instance is keyed by the exact server identity plus canonical worktree; all root sessions within that same instance share one budget. Different worktrees or server identities use separate Runtime instances and separate budgets. Root children and grandchildren in one instance share the same active-descendant budget. Do not create a separate four-orchestrator × twenty-grandchildren tier.

## Preserved boundaries

- Keep maximum hierarchy root→child→grandchild (depth 2); no great-grandchildren.
- Keep only researcher/proposer eligible to dispatch read-only grandchildren; preserve inherited exact scope.
- Keep worker-only write ownership, leases, one-writer-per-file, current task authorization, preflight, wrapper, cancel semantics, and fail-closed behavior.
- Keep per-Runtime budget shared across every root session attached to the same exact server identity and canonical worktree; different server identities/worktrees have independent Runtime budgets. Failed/blocked/uncertain jobs continue to consume capacity until confirmed stop, as today.
- Current slots may be lower due to live runtime/provider limits. Setting 20 is an admission ceiling, not a promise of 20 simultaneous model completions. No live-provider capacity claim.
- The current scheduler is immediate admission/rejection, not a queue. Admission is shared first-come across root sessions in the same Runtime (same server identity + canonical worktree); one such root may consume all 20 slots, and other roots in that Runtime receive an explicit capacity-exhausted failure. No fair-share guarantee, reservation between roots, or preemption is added. Different Runtime instances have independent budgets. This is a documented limitation rather than an acceptance claim.
- No VCS write, dependencies, external side effects, production service changes, or credential usage.

## Exact implementation scope

1. `agents/AGENT_ROUTER.md`: update concurrency guidance from six to max twenty when runtime slots are available; state process/worktree-wide shared budget and immediate reject-on-saturation/no-fair-share behavior.
2. `agents/SUBAGENTS.lrf`: change Orchestrator fanout from `min(6,...)` to `min(20,...)`, while preserving wait-all, depth, safety, and saturation behavior.
3. `opencode/plugins/orchestration.ts`: define named `MAX_ACTIVE_DESCENDANTS=20`, set root direct-child cap to 20, make all nested jobs share the same 20 active slots, remove the two-reserved-slots policy, and return a precise capacity-exhausted error. Do not otherwise alter scheduler behavior.
4. `opencode/tests/orchestration.test.mjs`: replace old six/four reservation test with deterministic tests for 20 root direct children allowed, 21st rejected, nested descendants consume the same Runtime budget, distinct root sessions in one Runtime share and cannot bypass the cap, separate Runtime instances have independent budgets, reservations consume slots before launch, confirmed completion/stop releases slots, blocked/uncertain jobs retain capacity, and existing depth/scope/write invariants remain.
5. `opencode/README.md`: document twenty shared descendant slots per server/worktree Runtime, up to twenty root children, first-come cross-root-session admission without fair-share or queueing, immediate rejection at saturation, reservations consuming slots before launch, and grandchildren consuming the same Runtime budget; separate Runtime instances have separate budgets; retain stated limitation that mocks do not prove live runtime/provider capacity.

## Verification

- Focused orchestration tests, if existing runtime test dependencies are available without changing dependencies; report missing prerequisite rather than installing dependencies.
- TypeScript/JavaScript syntax check for the plugin where supported.
- `git diff --check` and exact allowlist diff review.
- New tests must exercise 20, 21, nested budget consumption, release after confirmed termination, retained capacity while blocked/uncertain, and depth 3 rejection.
- Run LRF preflight with code/subagent/worker/doc tags and verify receipt.
- An independent diff-first verifier checks that no stale limit (6/4) remains in active rules/docs/runtime/tests and no unrelated safety boundary changed.
