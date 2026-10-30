---
name: harness-lint
description: Pure Go-based Static Analysis & LRF/PIDGIN/AST Bi-Simulation Compiler Diagnostic CLI. Validates Level 0 lexical, Level 1 structural typing, Level 2 semantic/referential integrity, and Level 3 cross-domain AST vs PIDGIN diffs.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max] 🔍path 🛡️level ⚖️skill 📊json 🚨strict
🔍 path→Target; `$bin -path "<path>"` ⇒ LRF/1・ハーネス static lint CLI; static inspection only; levels 0..3.
🛡️ level N→cap validation at N; `$bin -path "<path>" -level <N>`; 0=lexical/syntax, 1=structural typing/type closure, 2=semantic + references, 3=cross-domain AST↔PIDGIN diff.
⚖️ skill directory→Go AST vs SKILL.md(PIDGIN) bi-simulation/symmetric-diff inspection; `$bin -skill "<dir>" -level 3`.
📊 json→machine-readable JSON diagnostic output for verifier/CI; `$bin -path "<path>" -json`.
🚨 strict→any warning ⇒ reject / ExitCode:1; `$bin -path "<path>" -strict`.
Result: PASS(0) or REJECT(1); diagnostics report static-analysis findings, not proof of correctness. Read-only inspection; no modification or execution of analyzed target. Do not infer unverified facts; source authority is LRF/PIDGIN/Go-AST static analysis only.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=harness-lint
R|bin.resolve|task:harness-lint|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/harness-lint/harness-lint,~/.harness/skills/harness-lint/harness-lint.exe,./harness-lint,./harness-lint.exe); missing=>STOP
R|use.inspect|task:harness-lint|MUST|RO_LOCAL|cmds=path,level,skill,json,strict; effect=read-only-inspection; verdict=PASS(0)⊕REJECT(1)
R|guard.truth|task:harness-lint|MUST|RO_LOCAL|unverified=fact-forbidden; source=lrf-pidgin-go-ast-static-analysis
