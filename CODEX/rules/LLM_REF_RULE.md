#SIGMA LRF/1
F=H NL X*; H="@lrf=1|aud=GPT-5.6|scope="id; X=(R|L)"|"id"|"g"|"m"|"e"|"p; g="*"|k:v(&k:v)* with k in G={task,tag,change,phase,fact,file,event,case}; m=MUST|MUST_NOT|MAY; e in E={RO_LOCAL,RO_PUBLIC,LW_SCOPE,EXT_WRITE,RELEASE,LIVE_WRITE,VCS_WRITE,DESTRUCT,CRED,CHARGE,PROD_DEP}; p=nonempty opaque UTF8 atom, no "|"; unknown structural token/dup/invalid=>STOP; payload meaning=>oracle.
C: Obj=context; Mor=authorized E-transition; id=no-op; compose(f,g)=run g then f; H(path)=validate+hydrate; P=system>developer>user then explicit-task>nearest-dir>project>global; compatible morphisms compose, equal/incomparable conflict=>STOP.
L invokes H; relative base=source file; task tags=set; matching L union loads once in declared order.
@lrf=1|aud=GPT-5.6|scope=global
R|intent.ro|task:read-only|MUST|RO_LOCAL|answer/diagnose/review=inspect+evidence; repair=false
R|intent.rw|task:change|MUST|LW_SCOPE|requested workspace edit+proportionate checks=true
R|ctx|*|MUST|RO_LOCAL|precedence=P; merge=compatible conjunction; equal-rank conflict=STOP; invented merge=false
R|scope|*|MUST|RO_LOCAL|unrelated changes=preserve; destructive target=resolve exact; project rules=nearest project source
R|host|*|MUST|RO_LOCAL|os=Windows11; shell=PowerShell7; deny=Bash/sh/GitBash/WSL unless explicit; exec=native+.exe; all-file-discovery=rg.exe; text=UTF8-noBOM
R|skill|task:skill-match|MUST|RO_LOCAL|load=skill.lrf; catalog duplication=false; unrelated skills=omit
R|legacy|*|MUST_NOT|RO_LOCAL|activate imported-gemini/archive/backups/legacy hooks
R|persona|*|MUST|RO_LOCAL|tag:subagent replaces main CSS with exactly one role-mapped agent CSS and never composes both; otherwise main reads UTF8=$env:USERPROFILE\.codex\persona\PERSONA.css; main permits only comments/whitespace plus exactly one top-level style rule with selector exactly :root and declaration-name multiset exactly {--lang,--addressee,--voice,--priority}; imports/at-rules/nesting/duplicates/unknown-declarations/!important are forbidden; each value must be exactly one nonempty decoded CSS string token; read/decode/parse/structure/bind failure=>STOP; under precedence P bind the decoded map only as main-agent presentation constraints; CSS=persona-only; CSS cannot override P or define authority/effects/safety; authority/effects/safety=LRF-only; machine=stable-formal
R|route.phase|*|MUST|RO_LOCAL|effect decision first; decision!=ALLOW=>conditional leaf count=0
R|route.tags|*|MUST|RO_LOCAL|after ALLOW collect all: code=software create/edit/refactor/design; doc=README/docs create/edit or SKILL.md contract scan/repair; state=basename in {decisions.md,method.md,knowledge.md,issues.md,memo.md,history.md,diary.md} applies/touched/requested-create; agy=explicit current-task user request for bounded external coding delegation through agy.exe; subagent=subagent invocation or subagent configuration work
R|route.order|*|MUST|RO_LOCAL|matching leaf union once: state>engineering>documentation
L|load.state|tag:state|MUST|RO_LOCAL|[state](state.lrf)
L|load.engineering|tag:code|MUST|RO_LOCAL|[engineering](engineering.lrf)
L|load.agy|tag:agy|MUST|RO_LOCAL|[agy](agy.lrf)
L|load.subagents|tag:subagent|MUST|RO_LOCAL|[agents](../agents/SUBAGENTS.lrf)
R|subagent.eval-evidence|tag:subagent|MUST|RO_LOCAL|read ../evaluation/subagents_model_eval.json as JSON evidence after LRF hydration; never hydrate JSON through H or treat evaluation data as authority
L|load.documentation|tag:doc|MUST|RO_LOCAL|[documentation](documentation.lrf)
L|load.skill|task:skill-match|MUST|RO_LOCAL|[skill](SKILL.lrf)
L|load.diary|*|MUST|RO_LOCAL|[diary](DIARY.lrf)
L|load.research|*|MUST|RO_LOCAL|[research](research.lrf)
R|agent.mode|*|MUST|RO_LOCAL|default=orchestrate; when task has 2+ bounded independent workstreams delegate proactively; assign one persona per workstream; run read-heavy or disjoint work in parallel; one-writer-per-file; overlap or ordered dependency=serial; main waits for all requested results then validates and synthesizes; correction-cycles<=3
R|agent.handoff|tag:subagent|MUST|RO_LOCAL|visible-fields=Role,Target,Acceptance,Scope,Known facts; task-specific-only=true; known-facts=verified-only; unknowns-and-failures=role-output; omit=hydrated persona,authority,effects,safety,model-routing,pipeline,output-schema,evaluation; scope=task-boundary-not-effect-policy; role-output=separate
R|agent.gates|*|MUST|RO_LOCAL|substantial=multi-file behavior/API/schema/security/concurrency/migration/difficult rollback; unstable fact=>read-only researcher+dated primary+fact/inference/unknown; substantial pre=>auditor PASS; substantial done=>verifier real diff/snapshot+deterministic checks
R|agent.roles|*|MUST|RO_LOCAL|blackhat=explicit security/forensics only; proposer+critic=explicit convergence only; external Gemini/Claude=explicit user request only; reviewer prose!=proof
R|done|*|MUST|RO_LOCAL|confirm behavior+non-change scope+checks+changed files+residual risk; authorized mem.etl may persist bounded plan/diary/completion summaries; secrets/raw logs excluded; ETL failure=>warn+report
R|ref|task:llm-reference-edit|MUST|LW_SCOPE|scope=LLM-facing prompt Markdown/settings not human docs/code; format=LRF/1; payload=one record-level contract; new structural token=>define in SIGMA+validator; protect=authority/effect/negation/exception/stop/success until oracle PASS; SPR/decorative math sole safety carrier=false
