---
name: split-to-prs
description: Split current work into small reviewable PRs. Use when the user asks to split a chat, set of changes, branch, or PR.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] ✂️split 💾snapshot 🌿branch 🚀submit
✂️ 分割: Trigger=明示的依頼(split chat/change-set/branch/PR). Input={DirtyDiff,Intent,base,ownership}; inspect diff-stat＋CODEOWNERS; propose reviewer-interest slices before any action. Output=reviewable slice plan only. `git diff --stat <base>`
💾 Snapshot: Uncommitted changes⇒安全退避; 作業tree変更禁止 until snapshot. `git update-ref refs/backup/pre-split-$(date +%s) $(git stash create pre-split)`
🌿 ApprovedSlice⇒dedicated branch＋only specified files staged. `git checkout -b <branch> <base>` ∧ `git add <files>`
🚀 Push ∧ OpenPR⇒create PR, report URL list. `gh pr create --title "..." --body "..."`
Scope: planning is not authorization. Do not create branches, snapshot, stage, push, or create PR absent applicable explicit authorization; no VCS permission is inferred. Destructive git operations forbidden: `reset--hard`, `clean-fdx`, force-push, history rewrite without approval. Staging, if authorized, explicit files or hunks only; `git add -all` forbidden.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=split-to-prs
R|use.inspect|task:split-to-prs|MUST|RO_LOCAL|cmds=diff,ownership; scan=CODEOWNERS; propose-before-action=true
R|use.snapshot|task:split-to-prs|MUST|LW_SCOPE|create-ref=refs/backup/pre-split; working-tree-modify=forbid-until-snapshot
R|guard.git-destructive|task:split-to-prs|MUST|VCS_WRITE|forbid=reset--hard,clean-fdx,force-push,history-rewrite-without-approval
R|guard.stage|task:split-to-prs|MUST|LW_SCOPE|stage=explicit-files-or-hunks-only; git-add-all=forbidden
