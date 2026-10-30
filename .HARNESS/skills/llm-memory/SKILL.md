---
name: llm-memory
description: Multi-client Bi-Temporal & multi-level (L0-L3) agent memory and knowledge graph CLI client backed by PostgreSQL. Works from Windows & Debian CLI.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] CLI=llm-mem; `$bin`=resolved executable. Commands are categorized by the exact contract below; command examples preserve their literal flags.
🔎 Inspect (RO): status, clients, stock, search, semantic, evals, graph list, survey show/template/list/rank, analyze ⇒ read-only evidence/output; no memory mutation.
📥 Mutate (approval required): ingest, add, compact, supersede, embed, graph node/edge ⇒ live memory/database changes. Execute only with current-task approval; ingest only verified evidence. Walkthrough: ingest a file with category walkthrough. Research: ingest title/text with category knowledge.
🧾 Receipt evaluation: eval is a live write, not inspection. Requires memo/<task_id>/{knowledge,diary,walkthrough}.md and task_id, comparison_key, axes{obj}, attribution{prompt,agent}; bypasses compaction and graph; current-task approval required. evals only reads history.
📮 Survey: show/template/list/rank inspect. A receipt cue with ≥1 of DEFECT, FRICTION, REQUEST may optionally prompt a survey template keyed by form; submit is a live write and requires current-task approval. Do not submit by implication.
🛡️ Guard: never expose credentials, raw database URLs, raw secret logs, or unmasked environment values. Task-scoped paths, tests, and memory IDs may be used. Unverified claims are forbidden; label observations as survey_observation; state remaining uncertainty explicitly.
⛔ Resolution/failure: resolve `$bin` as first existing candidate in this exact order: `~/.harness/skills/llm-memory/llm-mem`, `~/.local/bin/llm-mem`, `~/.harness/skills/llm-memory/llm-mem.exe`, `$LLM_MEMORY_BIN`. If none exists, STOP. No command may be run without the governing approval; no failed or unverified operation is a successful receipt.
📡 Connectivity inspection: `$bin status` ∧ `$bin clients`
📥 Ingest (mutating): `$bin ingest -file <f> -cat <c> [-survey-signals <K:n>]` ∨ `$bin ingest -title <t> -text <s>`
📝 Add (mutating) / stock (inspect): `$bin add -title <T> -content <L0> [-cat <c>] [-l1 <l1>] [-l2 <l2>] [-tags <t>]` ∧ `$bin stock [-cat <c>] [-json]`
🔍 Search (inspect), semantic (inspect), embed (mutating): `$bin search -q "<q>" [-level <0..3>]` ∧ `$bin semantic -q "<q>"` ∧ `$bin embed -id <UUID>`
🔄 Compact/supersede (mutating): `$bin compact -limit <N>` ∧ `$bin supersede -id <UUID> -title <T> -content <L0>`
🩺 Analyze (inspect): `$bin analyze [-file <f>] [-suggest] [-json]`
⚖️ Eval (mutating) / evals (inspect): `$bin eval -file <f> [-json]` ∧ `$bin evals -key <k>`
🕸️ Graph node/edge (mutating), list (inspect): `$bin graph node -name <n>` ∧ `$bin graph edge -src <s> -tgt <t>` ∧ `$bin graph list`
📮 Survey subcommands: `$bin survey show|template|submit|list|classify|rank` (submit mutates; classify is not in the read-only inspection contract; do not infer its effect).
Capability model: BiTemporal(ValidTime×TxTime), progressive compaction(L0→L1→L2→L3), self-edit outcomes{ADD,UPDATE/Supersede,NOOP,DEPRECATE}, graph triples; canonical executable-resolution and approval contract controls. Reference commands: ingest -title <title> -text <text> -cat <cat>; search -tag <tag> -level 2; graph edge -src <S> -tgt <T> -rel <R>. Walkthrough/report ingest hooks are approval-gated references, never automatic arbitrary live writes. eval/evals remain as documented canonically.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=llm-memory
R|bin.resolve|task:llm-memory|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/llm-memory/llm-mem,~/.local/bin/llm-mem,~/.harness/skills/llm-memory/llm-mem.exe,$LLM_MEMORY_BIN); missing=>STOP
R|use.inspect|task:llm-memory|MUST|RO_LOCAL|cmds=status,clients,stock,search,semantic,evals,graph.list,survey.show,survey.template,survey.list,survey.rank,analyze; effect=read-only-inspection; output=text-or-json
R|use.mutate|task:llm-memory|MUST|LIVE_WRITE|cmds=ingest,add,compact,supersede,embed,graph.node,graph.edge; pre=verified-task-evidence; approval=current-task
R|use.pipeline|task:llm-memory|MUST|LIVE_WRITE|walkthrough=ingest-file-walkthrough-cat-walkthrough; research=ingest-title-text-cat-knowledge; approval=current-task
R|gate.eval|task:llm-memory|MUST|LIVE_WRITE|cmd=eval; pre=memo/<task_id>/{knowledge,diary,walkthrough}.md; input=task_id,comparison_key,axes{obj},attribution{prompt,agent}; bypass=compaction+graph; approval=current-task
R|gate.survey|task:llm-memory&event:task-receipt|MUST|RO_LOCAL|signals=DEFECT,FRICTION,REQUEST; threshold=1; cue=receipt-style-optional-survey; template=survey-template-form-key; submit=LIVE_WRITE&approval
R|guard.cred|task:llm-memory|MUST_NOT|CRED|deny=credentials,raw-db-url,raw-secret-logs,unmasked-env; allow=task-scoped-paths,tests,memory-ids
R|guard.truth|task:llm-memory|MUST|RO_LOCAL|unverified=fact-forbidden; observation-not-fact=survey_observation; remaining-uncertainty=explicit
