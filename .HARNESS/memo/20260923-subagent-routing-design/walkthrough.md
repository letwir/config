# walkthrough | 20260923-subagent-routing-design

## 実行計画
1. Read current SIGMA and matching subagent instructions.
2. Inspect current main/change/research ETL and the existing worktree diff.
3. Review recent nested-handoff evidence and design a single authoritative role selector.
4. State implementation scope, acceptance gates, and unresolved decisions without editing operational files.

## 設計案
- Preserve `etl/main.seq` as the overall ETL. Reduce only its subagent-specific logic to calls into a canonical selector.
- Add `agents/selection.lrf` as the sole normative role-selection matrix; optionally add `agents/selection.md` as a non-authoritative explanation. Load it through `rules/LOAD.lrf` and reference the same source from `agents/README.md`.
- Root/main agent alone selects and dispatches roles. No role may spawn a child; this matches the tested H-12 runtime boundary.
- Explorer: local repository/layout/source inventory. Researcher: current public/primary facts, API/docs, and evidence synthesis. Either alone by task need; both only for independent workstreams.
- Proposer: ambiguous/high-impact design draft. Auditor: gated independent plan review. Verifier: gated independent verification. Critic/Blackhat remain explicit-convergence/security-only. No routine fan-out.
- Create `agents/worker.md` as SRP/PIDGEN protocol for one bounded agy CODE task, with exact wrapper, ordered visible fields, input contract, exclusions, and compact result schema. Under the current policy, CODE execution remains direct agy from main.
- Full TOML pass removes `proactively`, aligns every description with selection triggers, eliminates repeated policy text where canonical references suffice, sets least-privilege sandbox per role, and removes fixed routing overrides unless explicitly required.
- Resolve `worker.toml`: retire it from Codex selection or repurpose as a read-only agy-task packager; do not leave an active Codex code-writing Worker alongside mandatory agy.

## Acceptance gates for implementation
- Every eligible role has exactly one selector predicate, protocol reference, existing CSS/TOML/Markdown requirements, trigger, and exclusion; no overlap except explicitly independent parallel work.
- `main.seq` contains no duplicate role-choice table and delegates selection to the canonical selector.
- All selected references resolve; no missing role protocol; no `proactively` language remains in agent TOMLs.
- Every TOML parses; harness-lint validates changed LRF; targeted route fixtures cover read-only/local research, CODE, high-risk plan audit, verification, explicit critic/security, and no-delegation cases.
- Fresh-session role discovery/runtime dispatch is checked separately; static validation alone is not called runtime success.
- Implementation preserves all pre-existing dirty changes outside the agreed allowlist.

## Known constraint
The current working tree already has extensive user changes in `etl/main.seq`, `agents/SUBAGENTS.lrf`, all nine agent TOMLs, rules, skills, and evaluation files. Future edits must patch the current working-tree versions surgically.
