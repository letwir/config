---
name: proof-checker
description: AST-based Category Theory & CODE_RULE.md Mathematical Soundness CI Gate CLI. Verifies ¬SilentSwallow, EarlyReturn, defer RAII, IO/Pure separation, and Concurrency safety across Go/Rust codebases.
---
[SPR/PIDGIN::ρ→max] 🛡️check 🚨strict 🧪vet 📊json
🛡️ 検証 check ⊢ Dir ⇒ AST健全性・圏論不変条件検証(¬SilentSwallow,EarlyReturn,RAII,IO/Pure) | `& $bin -path "<dir>"`
🚨 厳格 strict ⊢ Dir ⇒ 警告全REJECT厳格CIプルーフゲート | `& $bin -path "<dir>" -strict`
🧪 検査 vet ⊢ Dir ⇒ go vet 統合検証 | `& $bin -path "<dir>" -strict -vet`
📊 出力 json ⊢ Dir ⇒ CI/Verifier連携構造化JSON出力 | `& $bin -path "<dir>" -json`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=proof-checker
R|bin.resolve|task:proof-checker|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/proof-checker.exe,$env:USERPROFILE/.harness/skills/proof-checker/proof-checker.exe); missing=>STOP
R|use.inspect|task:proof-checker|MUST|RO_LOCAL|cmds=check,strict,vet,json; effect=read-only-inspection; verdict=PASS(0)⊕REJECT(1)
R|guard.truth|task:proof-checker|MUST|RO_LOCAL|unverified=fact-forbidden; source=go-ast-static-analysis
