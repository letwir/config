---
name: babysit
description: Keep a PR merge-ready by triaging comments, resolving clear conflicts, and fixing CI in a loop.
---
[SPR/PIDGIN::ρ→max] 👶babysit ⚔️conflict 💬triage 🧪ci 🚀ready
👶 監視 babysit ⊢ PR ⇒ CI・コメント・競合を自動修復しマージ可能状態へ導くループ | `gh pr status` ∧ `gh pr checks`
⚔️ 衝突 conflict ⊢ BaseBranch ⇒ マージ競合の解消（意図保持、競合深刻時は即中断） | `git merge <base>`
💬 コメント triage ⊢ Comments ⇒ 未解決コメント・Bugbotの検証と対応修正 | `gh pr view --comments`
🧪 修正 ci ⊢ WorkflowLogs ⇒ PR起因のCI失敗修正・再検証ループ | `gh run view <run-id> --log-failed`
🚀 準備 ready ⊢ GreenChecks ⇒ 全CIパス＋全コメント対応済みのマージReady確認 | `gh pr view`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=babysit
R|use.inspect|task:babysit|MUST|RO_PUBLIC|cmds=status,comments,ci; gh-api=minimal-payload
R|use.edit|task:babysit&case:conflict,ci-fix|MUST|LW_SCOPE|preserve-intent=true; unrelated-code=untouched; ci-workflow-bypass=forbidden
R|guard.vcs|task:babysit&case:push,merge|MUST|VCS_WRITE|require-explicit-approval-if-pushing; force-push=forbidden-unless-explicit
R|guard.truth|task:babysit|MUST|RO_LOCAL|bugbot-validation=rigorous; disagree-or-uncertain=report-to-user
