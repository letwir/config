---
name: paddock-continue
description: Resume explicitly requested work in the PADDOCK repository from its durable Issue, Git, and verification state. The experimental prompt.md auto-selection mode requires that prompt to be explicitly invoked; this is not a scheduled loop.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# PADDOCK continuation

## Purpose and activation

Use this skill only when the user explicitly asks to continue, resume, or restart work in the PADDOCK repository. A general question about PADDOCK is not an execution request. This skill is a convenience procedure, not a PADDOCK authority source: PADDOCK's `decisions.md` ACCEPTED register and routed LRF remain authoritative inside that repository. The skill does not modify or supersede PADDOCK policy.

The ordinary CONTINUE/1 path reconciles state and handles one bounded Issue. The automatic design-choice and multi-Issue loop in PADDOCK's experimental `prompt.md` are enabled only when the user explicitly invokes `prompt.md` for the current task. Merely invoking this skill or saying “continue” does not activate that delegation. Prompt invocation remains task-local and grants no effect or persistent authority.

This is not a timer, recurring scheduler, or the `.harness` proposition-convergence `loop-agent`. Do not change the separate `.harness` `loop` or `loop-agent` skills as part of PADDOCK continuation.

## Read-only preflight and reconciliation

1. Confirm the exact PADDOCK Git root; require its `AGENTS.md`, `rules/BOOTSTRAP.lrf`, `decisions.md`, `method.md`, `THREAT_MODEL.md`, and `Issues.md`. Read `AGENTS.md`, compile its bootstrap and matching LRF records, then read the design and state sources in their prescribed order. Treat `docs/handover.md` as a stale, non-authoritative snapshot.
2. Inspect actual HEAD, branch/upstream/divergence, index, worktree diff, recent history, active Issue, dependencies, and relevant verification Receipts. Reconcile the durable Cursor against those observations; a Cursor value or checkbox is a hint, not proof. Do not overwrite a discrepancy silently.
3. Run PADDOCK's required `scripts/lint_design.py` and `scripts/continue_plan.py` using the repository's documented command and the active host profile. A missing script, lint/planner error, invalid Cursor, or unreadable required source is a stop condition; do not install dependencies or skip the failed gate silently.
4. Treat `continue_plan.py` output only as a candidate `{issue_id, reason}`. It does not execute an Issue or establish READY/DONE, receipt validity, dependency completion, priority, or every safety gate. Independently reconcile its candidate with observed files and Receipts.
5. Apply `method.md` selection rules: resume a valid active WIP/IMPLEMENTED/TESTED Issue only after reconciling its partial work and dependencies; otherwise select only work whose dependencies are DONE by observed evidence and whose blockers/gates pass. Give Safety/OPEN blockers the specified precedence, then use P0 → P1 → P2, the small-Issue dependency DAG, and blockers. If Cursor, planner, actual state, or authority disagree, stop and report the discrepancy rather than silently selecting or rewriting state.

## Task-local design-choice mode

Read and apply this section only when the user explicitly invokes PADDOCK `prompt.md` for the current task. That prompt is experimental and non-canonical; its delegation expires with the task.

Before CODE, when a material design choice is unspecified:

- Enumerate exactly three distinct, honest, feasible options when three exist. If fewer than three are feasible, explain why and stop under the governing choice rule; never invent an option.
- Record the Main recommendation and ask a distinct, read-only evaluator, routed through the currently available Harness process. The evaluator may compare evidence and feasibility but must not implement, test, access DB/live state/credentials/network, or perform VCS effects.
- A candidate is strictly superior only when it is feasible in the actual repository, satisfies ACCEPTED decisions and Issue acceptance, preserves safeguards, adds no unapproved scope/effect, is evidence-supported, and is materially better on a stated criterion. Select the best-supported superior candidate; if none qualifies, retain the recommendation.
- Stop and ask if ACCEPTED authority is missing or conflicting, the evaluator is unavailable or materially disagrees, feasibility/evidence is uncertain, rankings are tied or incomparable, a new material choice appears, or safe implementation requires unapproved scope/effects. Do not majority-vote or fabricate certainty.

This selection delegation never changes PADDOCK's ACCEPTED decisions, grants implementation/effect permission, or persists beyond this explicitly invoked task.

## Bounded Issue execution and continuation loop

Work serially on one small Issue at a time in this skill. Follow PADDOCK's exact acceptance criteria and the currently applicable execution/verification gates; do not infer permission from this skill or from `prompt.md`. Keep PADDOCK changes inside the confirmed PADDOCK root. Do not write to `.HARNESS` or its parent repository while acting in PADDOCK.

A continuation request does not itself approve DB, network, credential, production, dependency-install, memory, VCS, or other external effects. Stop at an approval boundary and request the exact authorization required. Preserve unrelated work and do not reset, overwrite, or silently repair user changes.

Use the PADDOCK state sequence without skipping: TODO → WIP → IMPLEMENTED → TESTED → DONE. Only mark DONE after the Runner actually executes the acceptance checks and records their exit codes, target revision, and relevant hashes; an independent Verifier or trusted CI reviews the final change; and Main reconciles evidence, findings, all acceptance conditions, and residual risks. A Worker/Verifier report or planner output alone is not evidence. Disclose when checks or CI were not run.

Without explicit `prompt.md` invocation, stop after reconciling one Issue's result; do not automatically start another Issue. With explicit prompt invocation, proceed to the next READY Issue only after the previous Issue is truly DONE and the prompt's post-DONE gates have succeeded, in order: fresh current-task LIVE_WRITE approval and successful memory ingest; then fresh exact-scope VCS_WRITE approval, exact approved changes only, and verified VCS result. Then refresh `docs/handover.md` by whole-file replacement, update Cursor/generated index from observed state, rerun the required lint/planner, and select the next READY Issue. Missing approval, failed operation, unresolved review, or uncertain result stops the loop. Never infer future permission from an earlier approval.

## Stop and report

Stop for an unknown/wrong root, missing or conflicting policy, failed lint/planner, unverified dependency or Receipt, Cursor/Git mismatch, blocker, material uncertainty, unavailable required evaluator, unapproved effect, or failed verification. Preserve the worktree and report observed facts, unknowns, the exact blocker, and one safe next action. Never claim PADDOCK completion, CI attestation, or successful continuation without the corresponding observation.

#SIGMA LRF/1
@lrf=1|aud=LLM|scope=paddock-continue
R|trigger|task:paddock-continue|MUST|RO_LOCAL|activate=explicit-current-user-PADDOCK-continuation; auto-design-and-multi-issue-loop=explicit-current-task-prompt.md-only; scheduler-or-loop-agent=false
R|authority|task:paddock-continue|MUST|RO_LOCAL|PADDOCK=decisions.md-ACCEPTED+routed-LRF; skill=convenience-only; ¬override; PADDOCK-writes=confirmed-root-only; ¬.HARNESS-or-parent
R|preflight|task:paddock-continue|MUST|RO_LOCAL|read=AGENTS+bootstrap+decisions+method+THREAT_MODEL+Issues; observe=Git+worktree+Receipts; reconcile=Cursor; conflict-or-unknown=>STOP
R|planner|task:paddock-continue|MUST|RO_LOCAL|continue_plan.py=候補のみ; ¬execution∨READY∨DONE-proof; independently-check=dependencies+blockers+evidence+Safety/OPEN-priority+P0-P1-P2+DAG
R|design.auto|task:paddock-continue&event:explicit-prompt-invocation|MUST|RO_LOCAL|3-honest-feasible-options-when-available; distinct-read-only-evaluator; strictly-superior=>select; else=retain-recommendation; uncertainty/tie/disagreement/authority-or-effect-gap=>STOP
R|execution|task:paddock-continue|MUST|LW_SCOPE|ordinary-mode=one-serial-bounded-Issue; prompt-mode=next-only-after-DONE+memory-approval-and-success+exact-VCS-approval-and-verified-result; no-future-permission
R|verification|task:paddock-continue|MUST|RO_LOCAL|DONE-requires=Runner-observation+exit/hash/revision+independent-review-or-trusted-CI+Main-reconciliation+residual-risk; reports-alone!=evidence
R|effects|task:paddock-continue|MUST|RO_LOCAL|continuation-or-design-choice!=effect-approval; DB/network/credential/production/dependency/memory/VCS gates=independent-current-policy; unknown-or-denied=>STOP
