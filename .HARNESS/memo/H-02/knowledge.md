# Knowledge Base

## Context

対象: `scripts/sync-harness.ps1` / 検証日: 2026-09-08 / supersedes: none

## Findings

`-WhatIf`はSSOT directoryが存在しない場合を含め、同期fixtureのpath、type、hash、junction targetを変更しない。

## Morphism

`Obj(FilesystemSnapshot) -> Mor(sync preview) -> Obj(IdenticalFilesystemSnapshot)`
