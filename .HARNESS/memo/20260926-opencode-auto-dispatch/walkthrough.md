# Walkthrough: OpenCode automatic subagent routing (2026-09-26)

## Overview
Task: register on-demand subagents for GUI OpenCode and remove mandatory agy-first CODE routing. User authorized direct Luna workers when nested delegation is unavailable, up to six simultaneous independent children, and no total invocation-count limit. No agy was invoked.

## Deliverables
- Added explicit orchestrator and specialist/worker registration with scoped task allowlist in GUI OpenCode configuration and role Markdown frontmatter.
- Replaced active agy-first rules in harness bootstrap, router, LRF and ETL with one eligible context-selected worker; agy requires separate explicit request.
- Added nonnormative dated model-order fallback in evaluation/worker-context-priority.json; no comparable official task-consumption figures were confirmed. Context capacity is not consumption.
- Preserved task depth, effect gates, one-writer-per-file and three correction cycles; six is a parallel cap, not a total-call cap.

## Verification
- harness-lint rules level 2, agents SUBAGENTS level 1 and etl passed; rule preflight PASS. Focused git diff --check passed. Independent verifier confirmed all root orchestrator task targets registered, but rejected unattended CODE because worker edit permission remains ask. This is retained intentionally to avoid broad filesystem write permission; the GUI may prompt before an approved edit.
- GUI-only runtime was not restarted, so actual Orchestrator-to-Researcher nested dispatch and runtime edit prompts remain unverified. Existing unrelated dirty files were not changed intentionally. No VCS writing or deployment.
- Prompt-defect assessment: 0% for final scoped instruction; agent/tool limitation: GUI cannot be restarted and its runtime cannot be proven by static checks. Attribution is approximate.

## Remaining
Fully quit/restart OpenCode GUI, then test Orchestrator→Researcher and Orchestrator→Worker and inspect actual child session IDs. Keep worker edit approval unless a narrowly enforceable file-scoped permission is available. Date: 2026-09-26. Tags: opencode, subagents, routing, context-consumption.
