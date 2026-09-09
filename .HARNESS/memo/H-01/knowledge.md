# Knowledge Base

## Context

対象: `scripts/sync-harness.ps1` / 検証日: 2026-09-08 / supersedes: none

## Findings

health checkは全チェックを終えた後に集約結果を返す。現在の実環境は14チェック成功。異常fixtureでは非0終了を返す。

## Morphism

`Obj(SSOTと6共有先) -> Mor(health checks) -> Obj(Results, FailureCount, Healthy, process exit)`
