# Walkthrough

## Overview

2026-09-09、H-04としてrule distribution manifestとread-only doctorを追加した。

## Deliverables

- `rules/client-distribution.json`: canonical SHA256、entry chain、expected hashes、allowed LRF record IDs。
- `scripts/check-rule-distribution.ps1`: manifest/chain/hash/LRF payload検証とJSON receipt。
- `scripts/check-rule-distribution.Tests.ps1`: Exact、Allowed、Unexpected、payload drift、EntryMismatch、malformed/duplicate、optional、invalid manifest。

## Verification

- Live doctor exit 0: harness-project Exact、Codex AllowedDifference、Gemini Exact、OpenCode AllowedDifference、Agents/Claude NotConfigured。
- Codex AGENTS、Gemini GEMINI、OpenCode config→AGENTS→ruleの参照文字列を検査。
- Pester 3.4: Windows PowerShell 8/8、pwsh 8/8。
- Parser、manifest JSON、`git diff --check`: PASS。
- 独立Verifier: PASS。

## Uncertainty

AgentsとClaudeはrule entry未構成のため任意NotConfigured。実クライアントプロセスの新規起動はH-06の対象。

## Attribution

初回実装はSIGMA prelude誤判定、entry参照未検証、同一ID payload driftの見逃し、Windows PowerShell encoding非互換があった。いずれもAgentDefectとして完了前に修正。PromptDefectなし。当該手戻りはPromptDefect 0%、AgentDefect 100%。

Tags: harness, H-04, rules, manifest, doctor, drift
