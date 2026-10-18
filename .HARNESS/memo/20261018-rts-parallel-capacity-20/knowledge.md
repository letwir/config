# Knowledge / 20261018-rts-parallel-capacity-20

- topic: RTS parallel descendant capacity
- verified-date: 2026-10-18
- scope: `.HARNESS/agents` and `.HARNESS/opencode` RTS admission/routing
- facts:
  1. Before this change, the RTS scheduler admitted at most six active descendants per Runtime and at most four direct root children, reserving two slots for nested work.
  2. Runtime instances are keyed by server identity plus canonical worktree. Root sessions in the same Runtime share one in-memory `State`; different server/worktree identities use separate Runtime states.
  3. The implemented policy allows at most twenty active descendants total per Runtime, including grandchildren, with at most twenty direct root children. Slots are first-come, immediate-reject, no queue/fair-share/preemption; blocked/uncertain jobs continue consuming slots until confirmed stop.
  4. Depth remains root→child→grandchild only. Only researcher/proposer children may dispatch inherited-scope read-only grandchildren; workers/verifiers remain leaves. Existing write leases and scope guards remain.
  5. Static syntax and LRF preflight passed. The orchestration suite passed 17/17 using a temporary Node loader stub for the absent `@opencode-ai/plugin` package. Normal `node --test` cannot resolve that package in this checkout. No dependency was installed and no live runtime/provider capacity was tested.
- source: `.HARNESS/opencode/plugins/orchestration.ts`, `.HARNESS/opencode/tests/orchestration.test.mjs`, router/LRF/README, verified 2026-10-18
