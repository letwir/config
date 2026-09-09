# Knowledge Base

## Context

対象: rule distribution / 検証日: 2026-09-09 / supersedes: none

## Findings

現在の有効entry chainはHarness AGENTS→canonical、Codex AGENTS→Codex rule、Gemini GEMINI→canonical、OpenCode config→AGENTS→OpenCode rule。Codex/OpenCode差分はmanifestでrecord ID単位に許容される。

## Morphism

`Obj(Manifest, ClientEntries, LRFs) -> Mor(ReadOnlyDoctor) -> Obj(ClassifiedDriftReceipt)`
