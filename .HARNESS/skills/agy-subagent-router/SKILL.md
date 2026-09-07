---
name: agy-subagent-router
description: Route a bounded external coding task through agy.exe with current model discovery, bounded permissions, reasoning effort, and compact structured output.
---
[SPR/PIDGIN::ρ→max] 📋models 🤖delegate ⚡exec 📊eval
📋 モデル models ⊢ ∅ ⇒ 利用可能モデル一覧取得・確証 | `agy.exe models`
🤖 委譲 delegate ⊢ Task⊗Model ⇒ 外部エージェントハンドオフ作成・計画実行 | `agy.exe --print --model <id> --effort <low|med|high> --mode plan --output-format json --prompt "<task>"`
⚡ 実行 exec ⊢ Task⊗Files ⇒ 排他的ファイルスコープ外科的編集実行 | `agy.exe --print --model <id> --mode accept-edits --output-format json --prompt "<task>"`
📊 評価 eval ⊢ RunResult ⇒ 結果検証 & subagents_model_eval.json 記録 | SUCCESS ⊕ FAILED ⊕ ADVICE ⊕ ESCALATE

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=agy-subagent-router
R|bin.resolve|task:agy|MUST|RO_LOCAL|bin=first-existing(agy.exe,$env:USERPROFILE/bin/agy.exe); missing=>STOP
R|use.inspect|task:agy|MUST|RO_LOCAL|cmds=models; effect=read-only-inspection
R|use.delegate|task:agy|MUST|LIVE_WRITE|cmds=delegate,exec; pre=explicit-user-request-for-external-delegation; approval=current-task
R|guard.perm|task:agy|MUST_NOT|LW_SCOPE|flag=--dangerously-skip-permissions; external-authority-expansion=forbidden
R|guard.cred|task:agy|MUST_NOT|CRED|deny=credentials,tokens,raw-secret-logs; allow=task-prompts,model-ids
R|guard.truth|task:agy|MUST|RO_LOCAL|unverified=fact-forbidden; main-agent-diff-inspection=mandatory
