---
name: vm-monitor
description: VictoriaMetrics cluster monitoring, MetricsQL/PromQL instant/range query execution, TSDB health verification, and intelligent spare compute resource capacity analysis.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🩺status 📊metrics ⚡query 📈query_range 🖥️nodes 💡suggest
🩺status ⊢ ∅ ⇒ inspect TSDB health + series count + metric overview; connectivity/health status only, not a write. | `$bin -status` ∨ `$bin -url <URL> -status`
📊metrics ⊢ Keyword ⇒ list/search registered metric series; keyword narrows search. | `$bin -metrics` ∨ `$bin -search "<keyword>"`
⚡query ⊢ MetricsQL/PromQL expr ⇒ instant-vector evaluation; result reflects queried expression/time, not an unqualified capacity guarantee. | `$bin -query "<expr>" [-format json]`
📈query_range ⊢ MetricsQL/PromQL expr ⊗ Start ⊗ End ⊗ Step ⇒ evaluate across requested interval/step; returned samples are time-bounded. | `$bin -query-range "<expr>" -start "<start>" -end "<end>" -step "<step>"`
🖥️nodes ⊢ ∅ ⇒ per-node aggregate CPU/RAM/Disk/Load1; `-top` is alternate node view. | `$bin -nodes` ∨ `$bin -top`
💡suggest ⊢ ∅ ⇒ score spare compute capacity + suggest workload placement; recommendation ≠ guaranteed available capacity or safe placement. | `$bin -suggest` ∨ `$bin -capacity`
Interpretation: query values/series describe observed VictoriaMetrics data; empty, stale, partial, or failed responses do not prove zero usage, health, or capacity. Distinguish query/transport failure from valid results; if source, freshness, units, scope, or result meaning is unclear, preserve uncertainty and do not infer.
Boundary: monitoring/query commands are read-only; do not modify metrics, nodes, workloads, or service configuration. Capacity suggestions are estimates, not authorization to allocate/deploy; no threshold, reserve, freshness window, or placement guarantee is specified here.
Reference queries (preserve literal documented expressions): CPU=`100 - (avg by (instance) (rate(node_cpu_seconds_total{mode!="idle"}[5m])) * 100)`; RAM=`node_memory_MemAvailable_bytes / (1024 * 1024 * 1024)`; disk=`node_filesystem_avail_bytes{mountpoint="/"} / (1024 * 1024 * 1024)`; load=`node_load1`; cache=`vm_cache_requests_total{type="metric_id"} > 0`; overload=`node_load1 > count by(instance)(node_cpu_seconds_total{mode='idle'})`. -nodes/-top and -suggest/-capacity may request -format json. -suggest -format json → CapacityRank==1 evidence → deployment recommendation only, never deployment authorization. Candidate workloads={LLM_Inference,BatchPipeline,Microservices,EdgeDaemons}; no query execution.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=vm-monitor
R|bin.resolve|task:vm-monitor|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/vm-monitor/vm-mon,~/.harness/skills/vm-monitor/vm-mon.exe,./vm-mon,./vm-mon.exe); missing=>STOP
R|use.inspect|task:vm-monitor|MUST|RO_LOCAL|cmds=status,metrics,query,query_range,nodes,suggest; effect=read-only-inspection; default-url=opt($VM_ADDR,"http://docker.tigris-tailor.ts.net:8428")
R|guard.truth|task:vm-monitor|MUST|RO_LOCAL|unverified=fact-forbidden; source=victoriametrics-tsdb
