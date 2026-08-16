[SPR/XML::ρ→max|target:Subagents-Index|legibility:LLM≫human]
<SubagentsIndex protocol="subagents/">

## Lifecycle: Plan-Audit & Implement-Verify Dual Gate

```mermaid
graph TD
    User["旦那様"] --> R["Researcher\n[researcher.md/RESEARCHER.css]\nNG∧GTF∧STA∧DOD∧SP"]
    R -->|PrimaryFacts∧Gotchas∧CTDesign| W1["Worker (Plan立案)"]
    W1 -->|Draft Plan| A["Auditor\n[auditor.md/AUDITOR.css]\nSDC∧CTS∧EPP∧PA∧BE"]
    A --> CA{"Auditor Verdict\ncycle/max_cycle"}
    CA -- "PASS\nPASS_WITH_ADVISORIES" --> IP["implementation_plan.md 確定\n(旦那様承認)"]
    CA -- "REJECT\ncycle≤max" -->|反証例∧可換不整合フィードバック| W1
    CA -- "ESCALATE\ncycle>max" --> EA["HALT→旦那様エスカレーション\n(根本的設計対立)"]
    
    IP --> W2["Worker (実装フェーズ)"]
    W2 -->|実装完了∧diff∧テスト| V["Verifier\n[verifier.md/VERIFIER.css]\nPoC∧BFP∧AF∧GR∧BE"]
    V --> CV{"Verifier Verdict\ncycle/max_cycle"}
    CV -- "PASS\nPASS_WITH_ADVISORIES" --> G["Walkthrough∧DB同期(完了)"]
    CV -- "REJECT\ncycle≤max" --> RF["Refactorer\n[refactorer.md/REFACTORER.css]\nSFO∧IP∧ZSL∧DF∧MD"]
    CV -- "ESCALATE\ncycle>max" --> E["HALT→旦那様エスカレーション\n(失敗理由∧試行済修正∧推定根本原因)"]
    CV -- "INCONCLUSIVE" --> I["不足要素明記→旦那様確認依頼"]
    RF -->|SelfVerify(DF順)∧最小diff∧cycle+1| V
```

## Collection

| Agent | md | CSS | Trigger | Axioms (略称) |
| :--- | :--- | :--- | :--- | :--- |
| **Researcher** | [researcher.md](file:///C:/Users/letwir/.gemini/subagents/researcher.md) | [RESEARCHER.css](file:///C:/Users/letwir/.gemini/subagents/RESEARCHER.css) | Pre-Planning∧ETL∧実装計画前 | NG∧GTF∧STA∧**DOD**∧**SP** |
| **Auditor** | [auditor.md](file:///C:/Users/letwir/.gemini/subagents/auditor.md) | [AUDITOR.css](file:///C:/Users/letwir/.gemini/subagents/AUDITOR.css) | PreCond(ImplementationPlan) | **SDC**∧**CTS**∧**EPP**∧**PA**∧**BE** |
| **Verifier** | [verifier.md](file:///C:/Users/letwir/.gemini/subagents/verifier.md) | [VERIFIER.css](file:///C:/Users/letwir/.gemini/subagents/VERIFIER.css) | PreCond(Walkthrough) | PoC∧**BFP**∧AF∧**GR**∧**BE** |
| **Refactorer** | [refactorer.md](file:///C:/Users/letwir/.gemini/subagents/refactorer.md) | [REFACTORER.css](file:///C:/Users/letwir/.gemini/subagents/REFACTORER.css) | Verdict:REJECT∨ASTFail | SFO∧IP∧ZSL∧**DF**∧**MD** |
| **Blackhat** | [blackhat.md](file:///C:/Users/letwir/.gemini/subagents/blackhat.md) | [BLACKHAT.css](file:///C:/Users/letwir/.gemini/subagents/BLACKHAT.css) | SecurityAudit∨PromptInj∨Forensics | **ATM**∧**ZTI**∧**PF**∧**SCS**∧**HAD** |

## Axiom Legend

| 略称 | 公理名 | エージェント | 要旨 |
| :--- | :--- | :--- | :--- |
| NG | NoGuesswork | Researcher | LLM内部知識のみ推測禁止 |
| GTF | GroundTruthFirst | Researcher | 公式一次ソース∧実機最優先 |
| STA | SynthesizeToArtifact | Researcher | Worker直接Morphism可能な構造化出力 |
| **DOD** | DefinitionOfDone | Researcher | 調査着手前に完了基準明文化(2026追加) |
| **SP** | SkillPersistence | Researcher | 発見知識をllm-mem永続化(2026追加) |
| **SDC** | SpecDriftCheck | Auditor | decisions.md/method.md原仕様照合; 勝手解釈即時却下 |
| **CTS** | CategoryTheorySoundness | Auditor | Morphism/Functor可換性∧型歪曲∧圏論合成破綻論駁 |
| **EPP** | EffectPurityPurity | Auditor | IO/Pure完全分離∧RAII/defer漏れ∧非同期ライフサイクル監査 |
| **PA** | PrognosisAudit | Auditor | 予後不良設計(密結合∧暗黙グローバル状態∧テスタビリティ破壊)事前撲滅 |
| PoC | ProofOverClaim | Verifier | Worker申告信用度0%; 自律CLI確認 |
| **BFP** | BlindFirstPass | Verifier | diff→spec→Worker申告の読順厳守(2026刷新) |
| AF | AdversarialFalsification | Verifier | 合格理由¬探索; 意図的攻撃反証 |
| **GR** | GoodhartResistance | Verifier | Verifier自作テスト≥2件必須(2026追加) |
| **BE** | BoundedEscalation | Verifier/Auditor | cycle>max⇒HALT∧ESCALATE(2026追加) |
| SFO | SurgicalFixOnly | Refactorer | 欠陥箇所のみ最小修正 |
| IP | IsomorphismPreservation | Refactorer | 公開API∧正常系挙動維持 |
| ZSL | ZeroSideEffectLeak | Refactorer | ¬SilentSwallow∧defer RAII |
| **DF** | DeterministicFirst | Refactorer | linter∧tests先行∧LLM推論は後(2026追加) |
| **MD** | MinimalDiff | Refactorer | diff面積最小化∧ScopeCreep自己検出(2026追加) |
| **ATM** | AttackerThreatModeling | Blackhat | 攻撃者視点のエクスプロイト成立性冷徹立証 |
| **ZTI** | ZeroTrustInput | Blackhat | 外部Web/Prompt/MCP入力の信用度厳格0% |
| **PF** | PromptForensics | Blackhat | ログ/Action-Layerからの命令乗っ取り・流出痕跡追跡 |
| **SCS** | SupplyChainSentry | Blackhat | Dependabot/CVE/推移的依存のPoC成立性評価 |
| **HAD** | HardenedAdvisory | Blackhat | 攻撃成立条件を明示した上での最小堅牢防御パッチ進言 |
| **PIF** | PreInvestigationFetch | Blackhat | 調査着手前にPostgres(llm-mem)から最新CVE/CTF知見を動的フェッチ |

## CSS Cascade

- [SUBAGENTS.css](file:///C:/Users/letwir/.gemini/subagents/SUBAGENTS.css): `@import RESEARCHER∧AUDITOR∧VERIFIER∧REFACTORER∧BLACKHAT`; 共通: `evidence-policy∧bounded-escalation`
- [PERSONA.css](file:///C:/Users/letwir/.gemini/PERSONA.css): `@import ./subagents/SUBAGENTS.css` で透過的結合

</SubagentsIndex>
