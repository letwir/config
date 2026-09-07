---
name: vm-monitor
description: VictoriaMetrics cluster monitoring, MetricsQL/PromQL instant/range query execution, TSDB health verification, and intelligent spare compute resource capacity analysis.
---
[SPR/PIDGIN::ρ→max] 🩺status 📊metrics ⚡query 📈query_range 🖥️nodes 💡suggest
🩺 疎通 status ⊢ ∅ ⇒ TSDB健全性 ∧ 系列数 ∧ メトリクス概要 | `& $bin -status` ∨ `& $bin -url <URL> -status`
📊 系列 metrics ⊢ Keyword ⇒ 登録メトリクス系列一覧・検索 | `& $bin -metrics` ∨ `& $bin -search "<keyword>"`
⚡ 即時 query ⊢ MetricsQL ⇒ 即時ベクトル評価・PromQL実行 | `& $bin -query "<expr>" [-format json]`
📈 範囲 query_range ⊢ MetricsQL⊗Start⊗End⊗Step ⇒ 範囲クエリ実行 | `& $bin -query-range "<expr>" -start "<start>" -end "<end>" -step "<step>"`
🖥️ ノード nodes ⊢ ∅ ⇒ ノード別リソース集計(CPU/RAM/Disk/Load1) | `& $bin -nodes` ∨ `& $bin -top`
💡 推奨 suggest ⊢ ∅ ⇒ 余剰コンピュート容量スコアリング & ワークロード配置提案 | `& $bin -suggest` ∨ `& $bin -capacity`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=vm-monitor
R|bin.resolve|task:vm-monitor|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/vm-mon.exe,$env:USERPROFILE/.harness/skills/vm-monitor/vm-mon.exe); missing=>STOP
R|use.inspect|task:vm-monitor|MUST|RO_LOCAL|cmds=status,metrics,query,query_range,nodes,suggest; effect=read-only-inspection; default-url=opt($env:VM_ADDR,"http://docker.tigris-tailor.ts.net:8428")
R|guard.truth|task:vm-monitor|MUST|RO_LOCAL|unverified=fact-forbidden; source=victoriametrics-tsdb
