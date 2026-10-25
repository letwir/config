@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

Σ=rules/BOOTSTRAP.lrf; authority=Σ+routed-LRF; 本書=nav≠policy; host/cmd=HARDWARE.md relative to HarnessRoot, never process CWD (active_profile).
Subagent handoff: visible fields exactly Role,Target,Acceptance,Scope,Known facts in order; header identifies text protocol/policy only, grants no authority; read/apply SIGMA/1+routed LRF.
READ+COMPILE(Σ)→classify→match→hydrate(matching MANUAL+LOAD)→etl/main.seq. Missing|unreadable|invalid|unknown structure|equal-rank conflict⇒STOP+ask.
rg-first; hit±context; full catalog=NO; unrelated=NO. Choose shell from active_profile.
```bash
# active_profile: android-termux-bash (or linux-native-bash)
# (when active_profile is windows-pwsh: use rg.exe and backslashes rules\*.lrf)
rg -n "^[RL]\|[^|]*\.<topic>" rules/*.lrf
rg -n "^[RL]\|[^|]*\|[^|]*<guard>:[^|]*" rules/*.lrf
rg -n "\|MUST_NOT\|" rules/*.lrf
rg -n -e "\|LW_SCOPE\|" -e "\|LIVE_WRITE\|" -e "\|VCS_WRITE\|" rules/*.lrf
rg -n "^[RL]\|.*<keyword>" rules/*.lrf
```

```bash
# active_profile: android-termux-bash
rg -n "^## " decisions.md
rg -n -e "<api id=" -e "<term>" -A 8 knowledge.md
rg -n "\[[ -~*x]\]" -C 3 issues.md
rg -n "^### " diary.md
rg -n "^[@🔎🔭🧠🧐⚒🏁🧾]" etl/*.seq
bash scripts/invoke-rule-preflight.sh --task change
# (when windows-pwsh: & .\scripts\invoke-rule-preflight.ps1 -Task change)
```

ETL: main.seq→research.seq→change.seq→finish.seq; dialect=@dialect harness-seq/1; research-local iff local/private evidence or explicit local/EPUB research; read=selected∪required-downstream.
Order PRECEDENT→RESEARCH→PLAN→PRE_VERIFY→CODE→POST_VERIFY→FINISH; read-only/LOCAL_REVIEW⇒FINISH without change stages.
State lookup: decisions, knowledge, issues, diary, etl/*.seq.
Verify=behavior∧non-change-scope∧checks∧changed-files∧residual-risk∧links/filenames∧exact-diff. secrets-output=NEVER; effect-approval=LRF; docs⇒permission=FALSE.
