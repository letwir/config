# Walkthrough: symbolic subagent payload implementation

## 実行計画

1. Compile bootstrap and routed subagent/documentation rules; inspect the current handoff contracts and local precedent.
2. Audit a three-file change covering the normative LRF handoff rule, codec contract, and worker wrapper example.
3. Dispatch one context-selected GPT-family CODE worker with exact file scope.
4. Run strict LRF lint and rule-preflight receipt verification; inspect actual target contents and adversarial handoff boundaries.
5. Correct the verifier's task-scope versus permission-policy ambiguity and rerun checks.
6. Record the result and residual workspace baseline limitation.

## Verification

- `skills/harness-lint/harness-lint.exe -path rules/MANUAL.lrf -level 0 -strict` → exit 0.
- `skills/harness-lint/harness-lint.exe -path agents/SUBAGENTS.lrf -level 0 -strict` → exit 0.
- `scripts/invoke-rule-preflight.ps1 -Task change -Tag code,worker,doc,subagent` → PASS; receipt hash verification → PASS.
- An intermediate lint run failed on unescaped `|` in an LRF payload; escaping the literal resolved it. The final strict lint/preflight snapshot passed.
- Verifier outcome: PASS_WITH_ADVISORIES; its advisory about task-specific effect constraints versus policy grants was applied to the LRF rule and worker example. Main-agent review confirmed the final payload/wrapper ordering and fixed field order.
- Changed target files: `rules/MANUAL.lrf`, `agents/SUBAGENTS.lrf`, `agents/worker.md`.
- Residual risk: broad pre-existing uncommitted/untracked state prevents comparison against a clean Git baseline; no unrelated changes were modified intentionally.
