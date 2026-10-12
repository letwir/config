# Harness · SPR/PIDGEN

Σ := rules/BOOTSTRAP.lrf; authority := Σ + routed LRF; 本書 := nav ≠ policy; host/cmd := ../HARDWARE.md (resolve active_profile: android-termux-bash -> bash; windows-pwsh -> pwsh).
For Codex subagents the payload must retain exactly the visible fields Role, Target, Acceptance, Scope, Known facts in order. The header identifies the text protocol and currently loaded policy only; it grants no authority and does not replace reading/applying SIGMA/1 plus routed LRF.
🧭 READ+COMPILE(Σ) → classify(task,effect) → 🔎match → hydrate(rules/{MANUAL,LOAD}.lrf, matching) → etl/main.seq.
⛔ missing|unreadable|invalid|unknown-structure|equal-rank-conflict(Σ/link) ⇒ STOP+ask.
🔎 rg-first → hit±context; 全文catalog=NO; unrelated=NO; shell := resolve active_profile in HARDWARE.md.

```bash
# active_profile: android-termux-bash (or linux-native-bash)
# (when active_profile is windows-pwsh: use rg.exe and backslashes rules\*.lrf)
rg -n "^[RL]\|[^|]*\.<topic>" rules/*.lrf
rg -n "^[RL]\|[^|]*\|[^|]*<guard>:[^|]*" rules/*.lrf
rg -n "\|MUST_NOT\|" rules/*.lrf
rg -n -e "\|LW_SCOPE\|" -e "\|LIVE_WRITE\|" -e "\|VCS_WRITE\|" rules/*.lrf
rg -n "^[RL]\|.*<keyword>" rules/*.lrf
```

🧬 etl/: main.seq → research.seq → change.seq → finish.seq.
dialect := `@dialect harness-seq/1`; preflight := `active_profile == android-termux-bash ? bash scripts/invoke-rule-preflight.sh --task <task> : .\scripts\invoke-rule-preflight.ps1 -Task <task>`.
research-local.seq iff local/private evidence or explicit(local/EPUB research); read := selected ∪ required-downstream.
順序 := PRECEDENT → RESEARCH → PLAN → PRE_VERIFY → CODE → POST_VERIFY → FINISH; read-only/LOCAL_REVIEW ⇒ FINISH without change stages.
🧾 state lookup:

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

✅ verify := behavior ∧ non-change-scope ∧ checks ∧ changed-files ∧ residual-risk ∧ links/filenames ∧ exact-diff.
🔐 secrets-output=NEVER; effect-approval := LRF; docs⇒permission=FALSE.
