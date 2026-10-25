---
name: proof-checker
description: AST-based Category Theory & CODE_RULE.md Mathematical Soundness CI Gate CLI. Verifies ¬SilentSwallow, EarlyReturn, defer RAII, IO/Pure separation, and Concurrency safety across Go/Rust codebases.
---

# proof-checker: 機械的・数学的健全性プルーフゲート CLI

[SPR/XML::ρ→max|legibility:LLM≫human|protocol:hydrate→exec]
<Γ id="skill:proof-checker">
header: AST ∧ CT_Soundness ∧ InvariantProof ∧ CIGate;
scope: "A:\Users\letwir\repo\mcp2go\src" ∨ "C:\Users\letwir\.gemini";
rules: {
  ERR-001: "¬SilentSwallow ∧ ¬_ = err ∧ ¬EmptyCatch",
  NEST-001: "EarlyReturn ∧ ¬NestedIfElse ∧ GuardClause(exit_fast) ∧ MaxDepth(3)",
  RES-001: "Acquire ⇒ defer Resource.Close() ∧ RAII",
  PUR-001: "Sep(IO_Effect, PureDomain) ∧ SideEffectFn[VerbPrefix] ∧ PureMorph[NounForm]",
  LEN-001: "fn[len>60L] ⇒ name_affix(Batch∨Mega∨Complex∨Handler)",
  CONC-001: "∀Async ⇒ Ctx{context.Context, Timeout} ∧ ¬Leak"
};
commands: {
  check_all: 'proof-checker.exe -path "A:\Users\letwir\repo\mcp2go\src"',
  strict_gate: 'proof-checker.exe -path "<dir>" -strict',
  full_ci: 'proof-checker.exe -path "<dir>" -strict -vet',
  json_audit: 'proof-checker.exe -path "<dir>" -json'
};
verdict: "PASS (Exit 0) ⊕ REJECT (Exit 1: Block CI/Walkthrough)";
</Γ>

## 概要 (Overview)
`proof-checker` は、`CODE_RULE.md` にて規定された圏論的射（Domain → Codomain）、副作用分離、エラーハンドリング完全性、リソース解放（RAII）、ネスト制限を Go AST（抽象構文木）レベルで機械的かつ超高速に静的検証・証明する CI Gate ツールです。

## 使用方法 (Usage)

### 1. ディレクトリ全体のプルーフ検証
```powershell
& "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\proof-checker\proof-checker.exe" -path "A:\Users\letwir\repo\mcp2go\src"
```

### 2. 厳格モード（警告もエラーとして扱い REJECT）
```powershell
& "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\proof-checker\proof-checker.exe" -path "." -strict
```

### 3. go vet 統合検証
```powershell
& "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\proof-checker\proof-checker.exe" -path "." -vet
```

### 4. JSON 出力（CI / Verifier 連携）
```powershell
& "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\proof-checker\proof-checker.exe" -path "." -json
```
