# Walkthrough: 接続パスの {HarnessRoot} 相対化と拡張子解決

## Overview
パス移植性の依頼に対し、実装を行わず現状棚卸しと設計方針の確定、および issues.md への H-14 起票までを実施した読み取り中心＋1 ファイル変更の作業。

## Deliverables
- `issues.md` に **H-14 / P1 / 接続パスの {HarnessRoot} 相対化と active_profile による拡張子解決** を追加。確認事実、設計方針 3 点（{HarnessRoot} 相対／両実装維持＋トークン解決／両ツリー同一内容）、Subtasks 7 項目を明記。
- 同一内容を `repo/config/.HARNESS/issues.md` へ反映。

## Verification
- `diff -rq ~/.harness ~/repo/config/.HARNESS` → exit 0（両ツリー一致）。
- `bash scripts/invoke-rule-preflight.sh --task change` → exit 0（全 LRF 構文検査 PASS）。
- 現環境の実体確認: `codex`・`agy`・`llm-mem` 有効、`harness-lint`・`proof-checker`・`sec-forensics` バイナリおよび `pwsh` は不在（`.ps1` 側の実行検証は現環境では不可）。

## Next Step
H-14 の Subtasks 1 番（`HARDWARE.md` への `path_resolution` 追加）から CODE 段階を開始。preflight の L 参照検証が `{HarnessRoot}` トークンを解決できるかを同じ作業で確認する。
