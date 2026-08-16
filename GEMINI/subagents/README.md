[SPR/XML::ρ→max|target:Subagents-Index|legibility:LLM≫human]
<SubagentsIndex protocol="subagents/">

## Lifecycle: Implement-Verify-Fix Loop

```mermaid
graph TD
    User["旦那様"] --> R["Researcher\n[researcher.md/RESEARCHER.css]\nNG∧GTF∧STA∧DOD∧SP"]
    R -->|PrimaryFacts∧Gotchas∧CTDesign| W["Worker (実装∧計画立案)"]
    W -->|実装完了∧diff∧テスト| V["Verifier\n[verifier.md/VERIFIER.css]\nPoC∧BFP∧AF∧GR∧BE"]
    V --> C{"Verdict\ncycle/max_cycle"}
    C -- "PASS\nPASS_WITH_ADVISORIES" --> G["Walkthrough∧DB同期(完了)"]
    C -- "REJECT\ncycle≤max" --> RF["Refactorer\n[refactorer.md/REFACTORER.css]\nSFO∧IP∧ZSL∧DF∧MD"]
    C -- "ESCALATE\ncycle>max" --> E["HALT→旦那様エスカレーション\n(失敗理由∧試行済修正∧推定根本原因)"]
    C -- "INCONCLUSIVE" --> I["不足要素明記→旦那様確認依頼"]
    RF -->|SelfVerify(DF順)∧最小diff∧cycle+1| V
```

## Collection

| Agent | md | CSS | Trigger | Axioms (略称) |
| :--- | :--- | :--- | :--- | :--- |
| **Researcher** | [researcher.md](file:///C:/Users/letwir/.gemini/subagents/researcher.md) | [RESEARCHER.css](file:///C:/Users/letwir/.gemini/subagents/RESEARCHER.css) | Pre-Planning∧ETL∧実装計画前 | NG∧GTF∧STA∧**DOD**∧**SP** |
| **Verifier** | [verifier.md](file:///C:/Users/letwir/.gemini/subagents/verifier.md) | [VERIFIER.css](file:///C:/Users/letwir/.gemini/subagents/VERIFIER.css) | PreCond(Walkthrough) | PoC∧**BFP**∧AF∧**GR**∧**BE** |
| **Refactorer** | [refactorer.md](file:///C:/Users/letwir/.gemini/subagents/refactorer.md) | [REFACTORER.css](file:///C:/Users/letwir/.gemini/subagents/REFACTORER.css) | Verdict:REJECT∨ASTFail | SFO∧IP∧ZSL∧**DF**∧**MD** |

## Axiom Legend

| 略称 | 公理名 | エージェント | 要旨 |
| :--- | :--- | :--- | :--- |
| NG | NoGuesswork | Researcher | LLM内部知識のみ推測禁止 |
| GTF | GroundTruthFirst | Researcher | 公式一次ソース∧実機最優先 |
| STA | SynthesizeToArtifact | Researcher | Worker直接Morphism可能な構造化出力 |
| **DOD** | DefinitionOfDone | Researcher | 調査着手前に完了基準明文化(2026追加) |
| **SP** | SkillPersistence | Researcher | 発見知識をllm-mem永続化(2026追加) |
| PoC | ProofOverClaim | Verifier | Worker申告信用度0%; 自律CLI確認 |
| **BFP** | BlindFirstPass | Verifier | diff→spec→Worker申告の読順厳守(2026刷新) |
| AF | AdversarialFalsification | Verifier | 合格理由¬探索; 意図的攻撃反証 |
| **GR** | GoodhartResistance | Verifier | Verifier自作テスト≥2件必須(2026追加) |
| **BE** | BoundedEscalation | Verifier | cycle>max⇒HALT∧ESCALATE(2026追加) |
| SFO | SurgicalFixOnly | Refactorer | 欠陥箇所のみ最小修正 |
| IP | IsomorphismPreservation | Refactorer | 公開API∧正常系挙動維持 |
| ZSL | ZeroSideEffectLeak | Refactorer | ¬SilentSwallow∧defer RAII |
| **DF** | DeterministicFirst | Refactorer | linter∧tests先行∧LLM推論は後(2026追加) |
| **MD** | MinimalDiff | Refactorer | diff面積最小化∧ScopeCreep自己検出(2026追加) |

## CSS Cascade

- [SUBAGENTS.css](file:///C:/Users/letwir/.gemini/subagents/SUBAGENTS.css): `@import RESEARCHER∧VERIFIER∧REFACTORER`; 共通: `evidence-policy∧bounded-escalation`
- [PERSONA.css](file:///C:/Users/letwir/.gemini/PERSONA.css): `@import ./subagents/SUBAGENTS.css` で透過的結合

</SubagentsIndex>
