# Knowledge Base

## Context

対象: `scripts/sync-harness.ps1` / 検証日: 2026-09-08 / supersedes: none

## Findings

skill統合は既定Failで競合時zero-write。更新を選ぶ場合だけ明示policyを要求し、manifestとbackupから適用結果・復元を追跡できる。

## Morphism

`Obj(SourceSnapshots, SSOT) -> Mor(Preflight, Policy, Apply) -> Obj(Manifest, Backup, RestorableSSOT)`
