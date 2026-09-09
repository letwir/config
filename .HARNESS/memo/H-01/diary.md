# Diary

### 2026-09-08

Hypothesis: 表示だけのエラーを集約結果とprocess exitへ変換すればCI境界で判定できる。
Tried: injectable UserHomeとtemp junction fixture。
Rejected: 実ユーザーの共有先を壊して異常系を作る方法。
Uncertainty: H-02以降。
Attribution: 初回healthy fixtureで単一DirectoryInfoのCountを見落としたAgentDefect。テストで検出し修正。PromptDefectなし。
Search: llm-mem先例なし。
Correction: `@(Get-ChildItem ...)`でcollection境界を固定。
Impact: H-01の完了前に解消、残存影響なし。
