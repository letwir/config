# Walkthrough

## Overview

2026-09-08、H-03としてskill統合をpreflight、明示policy、manifest、backup、restoreを持つ処理へ変更した。

## Deliverables

- ConflictPolicy: Fail（既定）、PreferSource、PreferDestination。
- relative path / type / SHA256による事前差分、固定source順、複数競合のPreviousSource記録。
- version 1 manifest: RunId、Status、Completed、Error、before/after、operation provenance。
- SSOT変更前backupと、drift拒否・Force対応のRestoreManifest。
- reparse/path containment、TOCTOU、partial failure、binary seedのguard。

## Verification

- PowerShell parser、`git diff --check`: PASS。
- Pester: 15 passed / 0 failed。
- Failはexit 1かつSSOT/artifact/downstream不変。
- PreferSourceの高優先`.codex`採用とPreviousSource `.agents`を確認。
- drift restoreはexit 1、Forceはexit 0。hidden/overwrite/add/delete復元を確認。
- partialはexit 1、Status=Partial、Completed=false。
- injected binary seedはSHA256一致。
- 独立Verifier: PASS。H-01/H-02の先行9テストもPASS。

## Uncertainty

H-04以降は未着手。実環境でConsolidateによる更新は実行していない。

## Attribution

初回Worker実装でbinary path、downstream停止、manifest provenanceをテストから漏らしたAgentDefectを主エージェントのdiff確認で発見し修復。PromptDefectなし。当該手戻りはPromptDefect 0%、AgentDefect 100%。

Tags: harness, H-03, conflict, manifest, backup, restore
