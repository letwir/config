# Harness · SPR/PIDGEN

Σ := rules/BOOTSTRAP.lrf; authority := Σ + routed LRF; 本書 := nav ≠ policy.
For Codex subagents the payload must retain exactly the visible fields Role, Target, Acceptance, Scope, Known facts in order. The header identifies the text protocol and currently loaded policy only; it grants no authority and does not replace reading/applying SIGMA/1 plus routed LRF.
🧭 READ+COMPILE(Σ) → classify(task,effect) → 🔎match → hydrate(rules/{MANUAL,LOAD}.lrf, matching) → etl/main.seq.
⛔ missing|unreadable|invalid|unknown-structure|equal-rank-conflict(Σ/link) ⇒ STOP+ask.
🔎 rg-first → hit±context; 全文catalog=NO; unrelated=NO.

```powershell
rg.exe -n "^[RL]\|[^|]*\.<topic>" rules\*.lrf
rg.exe -n "^[RL]\|[^|]*\|[^|]*<guard>:[^|]*" rules\*.lrf
rg.exe -n "\|MUST_NOT\|" rules\*.lrf
rg.exe -n -e "\|LW_SCOPE\|" -e "\|LIVE_WRITE\|" -e "\|VCS_WRITE\|" rules\*.lrf
rg.exe -n "^[RL]\|.*<keyword>" rules\*.lrf
```

🧬 etl/: main.seq → research.seq → change.seq → finish.seq.
dialect := `@dialect harness-seq/1`; lint := `& $env:USERPROFILE\.harness\skills\harness-lint\harness-lint.exe -path etl` (accepts `.seq` files and directories).
research-local.seq iff local/private evidence or explicit(local/EPUB research); read := selected ∪ required-downstream.
順序 := PRECEDENT → RESEARCH → PLAN → PRE_VERIFY → CODE → POST_VERIFY → FINISH; read-only/LOCAL_REVIEW ⇒ FINISH without change stages.
🧾 state lookup:

```powershell
rg.exe -n "^## " decisions.md
rg.exe -n -e "<api id=" -e "<term>" -A 8 knowledge.md
rg.exe -n "\[[ -~*x]\]" -C 3 issues.md
rg.exe -n "^### " diary.md
rg.exe -n "^[@🔎🔭🧠🧐⚒🏁🧾]" etl\*.seq
& $env:USERPROFILE\.harness\skills\harness-lint\harness-lint.exe -path etl
```

✅ verify := behavior ∧ non-change-scope ∧ checks ∧ changed-files ∧ residual-risk ∧ links/filenames ∧ exact-diff.
🔐 secrets-output=NEVER; effect-approval := LRF; docs⇒permission=FALSE.
