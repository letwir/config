---
name: harness-lint
description: Pure Go-based Static Analysis & LRF/PIDGIN/AST Bi-Simulation Compiler Diagnostic CLI. Validates Level 0 lexical, Level 1 structural typing, Level 2 semantic/referential integrity, and Level 3 cross-domain AST vs PIDGIN diffs.
---
[SPR/PIDGIN::ρ→max] 🔍path 🛡️level ⚖️skill 📊json 🚨strict
🔍 走査 path ⊢ Path ⇒ LRF/1・ハーネス静的解析実行 (Level 0〜3) | `& $bin -path "<path>"`
🛡️ 階層 level ⊢ Level ⇒ 指定階層まで検証(0:構文, 1:型閉包, 2:意味論, 3:AST差分) | `& $bin -path "<path>" -level <N>`
⚖️ 差分 skill ⊢ SkillDir ⇒ Go AST ↔ SKILL.md(PIDGIN) 双模倣対称差分検証 | `& $bin -skill "<dir>" -level 3`
📊 出力 json ⊢ Target ⇒ Verifier/CI連携機械可読JSON診断出力 | `& $bin -path "<path>" -json`
🚨 厳格 strict ⊢ Target ⇒ 警告1件でも即死判定(ExitCode:1) | `& $bin -path "<path>" -strict`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=harness-lint
R|bin.resolve|task:harness-lint|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/harness-lint.exe,$env:USERPROFILE/.harness/skills/harness-lint/harness-lint.exe); missing=>STOP
R|use.inspect|task:harness-lint|MUST|RO_LOCAL|cmds=path,level,skill,json,strict; effect=read-only-inspection; verdict=PASS(0)⊕REJECT(1)
R|guard.truth|task:harness-lint|MUST|RO_LOCAL|unverified=fact-forbidden; source=lrf-pidgin-go-ast-static-analysis
