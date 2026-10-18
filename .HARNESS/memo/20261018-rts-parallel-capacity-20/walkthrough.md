# Walkthrough / 20261018-rts-parallel-capacity-20

## Overview

Implemented the first requested alternative: a 20-descendant RTS admission ceiling per Runtime, including up to 20 direct root children. Grandchildren consume the same total budget.

## Changed files

- `agents/AGENT_ROUTER.md` — active concurrency guidance now permits up to 20 independent children when runtime slots exist; documents shared per-Runtime pool and reject-on-saturation.
- `agents/SUBAGENTS.lrf` — Orchestrator fanout cap changed from `min(6, ...)` to `min(20, ...)`, with shared descendant budget and no fair-share/queue language.
- `opencode/plugins/orchestration.ts` — named caps at 20; admission counts all non-stopped jobs and no longer reserves two nested slots; capacity error includes configured cap.
- `opencode/tests/orchestration.test.mjs` — added 20/21, shared-root/nested, independent State, lifecycle release/retention and registry identity tests. Fixture temp root now uses `TMPDIR`/`os.tmpdir()` rather than Windows-only `A:/TMP`.
- `opencode/README.md` — documents per-server/worktree Runtime budget, immediate first-come rejection, no queue/fairness/preemption, and mock/live-runtime limits.

## Verification

- `node --check .HARNESS/opencode/plugins/orchestration.ts` — exit 0.
- `node --check .HARNESS/opencode/tests/orchestration.test.mjs` — exit 0.
- `git diff --check` on the five-file scope — exit 0.
- LRF change preflight with `code,subagent,worker,doc`; receipt verification — PASS, 8 policy files.
- Stale-cap search over implementation/rules/tests/README — no old six/four enforcement references.
- `node --experimental-loader <temporary mock> --test .HARNESS/opencode/tests/orchestration.test.mjs` — 17/17 passed. The temporary ESM loader only supplies the absent plugin's `tool`/schema builder interface and is removed after the run.
- Normal `node --test .HARNESS/opencode/tests/orchestration.test.mjs` — blocked before test execution: `ERR_MODULE_NOT_FOUND @opencode-ai/plugin`. No dependency installation attempted.
- Independent verifier — PASS after lifecycle/registry tests were added.

## Limits and residual risk

The 20 slots are an in-memory admission ceiling per canonical Runtime (server identity + worktree), not measured provider/runtime capacity. Root sessions in the same Runtime share first-come slots; one root may take all 20, there is no queue or fairness/preemption. Different Runtime identities have separate budgets. Real OpenCode model/session concurrency and external rate limits remain unverified. No VCS write occurred; pre-existing unrelated worktree changes remain untouched.
