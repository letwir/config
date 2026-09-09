# Diary

### 2026-09-08

Hypothesis: 実コピー前の仮想状態と事前backupで競合を決定的かつ復元可能にできる。
Tried: Fail/PreferSource/PreferDestination、複数source、partial、restore、binary fixture。
Rejected: コピー順だけで暗黙上書きする従来方式。
Uncertainty: 実環境のConsolidate更新は未実行。
Attribution: 初回実装の検証漏れはAgentDefect、依頼は明確。
Correction: exact binary source、downstream stop、manifest provenance、source tracking、backup順序を追加。
Impact: 完了前に解消。
