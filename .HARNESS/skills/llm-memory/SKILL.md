---
name: llm-memory
description: Multi-client Bi-Temporal & multi-level (L0-L3) agent memory and knowledge graph CLI client backed by PostgreSQL. Works from Windows & Debian CLI.
---
[SPR/PIDGIN::ρ→max] 📡status|clients 📥ingest 📝add|stock 🔍search|semantic|embed 🔄compact|supersede 🩺analyze ⚖️eval|evals 🕸️graph 📮survey
📡 疎通 status|clients ⊢ ∅ ⇒ DB接続確認 ∧ 端末一覧 | `& $bin status` ∧ `& $bin clients`
📥 取込 ingest ⊢ File∨Text ⇒ JITMIND自己編集 ∧ 多段縮約(L0-L3) ∧ グラフ抽出 | `& $bin ingest -file <f> -cat <c> [-survey-signals <K:n>]` ∨ `& $bin ingest -title <t> -text <s>`
📝 記憶 add|stock ⊢ Title⊗L0 ⇒ 新規多段登録(L0~L3,type,scope,rationale,evidence) ∧ 概観一覧(L2→L1→L0) | `& $bin add -title <T> -content <L0> [-cat <c>] [-l1 <l1>] [-l2 <l2>] [-tags <t>]` ∧ `& $bin stock [-cat <c>] [-json]`
🔍 探索 search|semantic|embed ⊢ Query∨UUID ⇒ キーワード/タグ部分一致(level 0..3) ∧ Gemini埋め込み意味検索 ∧ 特徴量保存 | `& $bin search -q "<q>" [-level <0..3>]` ∧ `& $bin semantic -q "<q>"` ∧ `& $bin embed -id <UUID>`
🔄 更新 compact|supersede ⊢ Limit∨OldUUID ⇒ 未要約バッチProgressiveCompaction ∧ JITMIND無効化&後続新記憶登録 | `& $bin compact -limit <N>` ∧ `& $bin supersede -id <UUID> -title <T> -content <L0>`
🩺 因果 analyze ⊢ Diary ⇒ 因果帰属分析(PromptDefect vs AgentDefect) ∧ 最適化Diff提示 | `& $bin analyze [-file <f>] [-suggest] [-json]`
⚖️ 評価 eval|evals ⊢ EvalJSON ⇒ タスク完了評価記録(LLM抽出bypass直接登録) ∧ 履歴参照 | `& $bin eval -file <f> [-json]` ∧ `& $bin evals -key <k>`
🕸️ 知識 graph ⊢ Node∨Edge ⇒ グラフ三つ組操作(node登録 ∧ edge関係線 ∧ list一覧) | `& $bin graph node -name <n>` ∧ `& $bin graph edge -src <s> -tgt <t>` ∧ `& $bin graph list`
📮 意見 survey ⊢ Op ⇒ ハーネス改善アンケート(template雛形 ∧ submit送信 ∧ list回答 ∧ classify分類 ∧ rank順位) | `& $bin survey template|submit|list|classify|rank`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=llm-memory
R|bin.resolve|task:llm-memory|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/llm-mem.exe,$env:USERPROFILE/.harness/skills/llm-memory/llm-mem.exe,A:/Users/letwir/repo/llm-memory/llm-mem.exe,$env:LLM_MEMORY_BIN); missing=>STOP
R|use.inspect|task:llm-memory|MUST|RO_LOCAL|cmds=status,clients,stock,search,semantic,evals,graph.list,survey.template,survey.list,survey.rank,analyze; effect=read-only-inspection; output=text-or-json
R|use.mutate|task:llm-memory|MUST|LIVE_WRITE|cmds=ingest,add,compact,supersede,embed,graph.node,graph.edge; pre=verified-task-evidence; approval=current-task
R|use.pipeline|task:llm-memory|MUST|LIVE_WRITE|walkthrough=ingest-file-walkthrough-cat-walkthrough; research=ingest-title-text-cat-knowledge; approval=current-task
R|gate.eval|task:llm-memory|MUST|LIVE_WRITE|cmd=eval; pre=memo/<task_id>/{knowledge,diary,walkthrough}.md; input=task_id,comparison_key,axes{obj},attribution{prompt,agent}; bypass=compaction+graph; approval=current-task
R|gate.survey|task:llm-memory&event:task-receipt|MUST|RO_LOCAL|signals=DEFECT,FRICTION,REQUEST; threshold=1; cue=receipt-style-optional-survey; template=survey-template-form-key; submit=LIVE_WRITE&approval
R|guard.cred|task:llm-memory|MUST_NOT|CRED|deny=credentials,raw-db-url,raw-secret-logs,unmasked-env; allow=task-scoped-paths,tests,memory-ids
R|guard.truth|task:llm-memory|MUST|RO_LOCAL|unverified=fact-forbidden; observation-not-fact=survey_observation; remaining-uncertainty=explicit
