---
name: babysit
description: Keep a PR merge-ready by triaging comments, resolving clear conflicts, and fixing CI in a loop.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 👶babysit ⚔️conflict 💬triage 🧪ci 🚀ready
🔎 Trigger: requested babysit(PR) ∧ existing PR context ⊢ inspect status, checks, comments | `gh pr status` ∧ `gh pr checks`
⊢ Bound: maintain that PR only; iterate inspect→triage→eligible fix→verify. No merge/release or unrelated work.
⚔️ Conflict ⊢ merge BaseBranch only to resolve clear, understood conflicts while preserving intent; `git merge <base>`. Unclear/deep conflict or intent loss risk ⇒ ⛔ stop + report; do not guess.
💬 Comments ⊢ inspect unresolved review comments/Bugbot, validate each claim, fix only verified actionable PR issues | `gh pr view --comments`
🧪 CI ⊢ inspect PR-caused failures and logs; make minimal in-scope code fixes and re-run relevant checks | `gh run view <run-id> --log-failed`. Workflow bypass or unrelated edits forbidden.
✅ Success: checks green ∧ actionable comments addressed ∧ no unresolved blocker; report merge-ready status, evidence, fixes and remaining risks | `gh pr view`
⛔ Stop/escalate: ambiguous request/PR, uncertain Bugbot finding, unclear intent, serious conflict, non-PR/root cause, or unsafe/out-of-scope fix ⇒ stop that item, preserve state, report evidence and ask. Never claim success without verification.
Output: concise status→findings→actions/check results→remaining blockers/risks; distinguish verified facts from uncertainty. No push/merge unless separately explicitly authorized; apply routed LRF effects.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=babysit
R|use.inspect|task:babysit|MUST|RO_PUBLIC|cmds=status,comments,ci; gh-api=minimal-payload
R|use.edit|task:babysit&case:conflict,ci-fix|MUST|LW_SCOPE|preserve-intent=true; unrelated-code=untouched; ci-workflow-bypass=forbidden
R|guard.vcs|task:babysit&case:push,merge|MUST|VCS_WRITE|require-explicit-approval-if-pushing; force-push=forbidden-unless-explicit
R|guard.truth|task:babysit|MUST|RO_LOCAL|bugbot-validation=rigorous; disagree-or-uncertain=report-to-user
