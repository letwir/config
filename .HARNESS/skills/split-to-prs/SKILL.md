---
name: split-to-prs
description: Split current work into small reviewable PRs. Use when the user asks to split a chat, set of changes, branch, or PR.
---
[SPR/PIDGIN::ρ→max] ✂️split 💾snapshot 🌿branch 🚀submit
✂️ 分割 split ⊢ DirtyDiff ∧ Intent ⇒ レビュアー・関心事単位のスライス計画提案 | `git diff --stat <base>`
💾 退避 snapshot ⊢ Uncommitted ⇒ 安全退避スナップショット作成 | `git update-ref refs/backup/pre-split-$(date +%s) $(git stash create pre-split)`
🌿 作成 branch ⊢ ApprovedSlice ⇒ 専用ブランチ作成＋指定ファイルのみステージ | `git checkout -b <branch> <base>` ∧ `git add <files>`
🚀 提出 submit ⊢ Push ∧ OpenPR ⇒ PR作成とURL一覧報告 | `gh pr create --title "..." --body "..."`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=split-to-prs
R|use.inspect|task:split-to-prs|MUST|RO_LOCAL|cmds=diff,ownership; scan=CODEOWNERS; propose-before-action=true
R|use.snapshot|task:split-to-prs|MUST|LW_SCOPE|create-ref=refs/backup/pre-split; working-tree-modify=forbid-until-snapshot
R|guard.git-destructive|task:split-to-prs|MUST|VCS_WRITE|forbid=reset--hard,clean-fdx,force-push,history-rewrite-without-approval
R|guard.stage|task:split-to-prs|MUST|LW_SCOPE|stage=explicit-files-or-hunks-only; git-add-all=forbidden
