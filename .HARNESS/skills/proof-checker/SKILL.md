---
name: proof-checker
description: AST-based Category Theory & CODE_RULE.md Mathematical Soundness CI Gate CLI. Verifies ¬SilentSwallow, EarlyReturn, defer RAII, IO/Pure separation, and Concurrency safety across Go/Rust codebases.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

Obj(Dir: Go∪Rust source context; AST) → Mor(proof-checker CLI: static AST/code checks) → Obj(verdict∧optional JSON evidence).
🔎 trigger: inspect a directory for code properties below. ⊢ purpose=Category Theory & CODE_RULE.md Mathematical Soundness CI Gate; analysis source=go-ast-static-analysis; results are static-analysis findings, not proof of authority or runtime behavior.

Checks: ¬SilentSwallow ∧ EarlyReturn ∧ defer RAII ∧ IO/Pure separation ∧ Concurrency safety; AST健全性・圏論不変条件検証.

Invocation (preserve command/flags):
🛡️ check ⊢ Dir ⇒ AST/code-property inspection | `$bin -path "<dir>"`
🚨 strict ⊢ Dir ⇒ warnings are REJECTed; strict CI gate | `$bin -path "<dir>" -strict`
🧪 vet ⊢ Dir ⇒ go vet integrated check; strict mode | `$bin -path "<dir>" -strict -vet`
📊 json ⊢ Dir ⇒ structured JSON output for CI/Verifier integration | `$bin -path "<dir>" -json`

LRF defines the executable lookup, read-only scope, command set, and verdicts. ⛔ missing executable⇒STOP; unverified findings must not be stated as fact. ✅ reported verdict: PASS(0) ⊕ REJECT(1). JSON selection specifies output format, not broader authority.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=proof-checker
R|bin.resolve|task:proof-checker|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/proof-checker/proof-checker,~/.harness/skills/proof-checker/proof-checker.exe,./proof-checker,./proof-checker.exe); missing=>STOP
R|use.inspect|task:proof-checker|MUST|RO_LOCAL|cmds=check,strict,vet,json; effect=read-only-inspection; verdict=PASS(0)⊕REJECT(1)
R|guard.truth|task:proof-checker|MUST|RO_LOCAL|unverified=fact-forbidden; source=go-ast-static-analysis
