---
name: agy-subagent-router
description: Route a bounded external coding task through agy.exe with current model discovery, bounded permissions, reasoning effort, and compact structured output. Use only after an explicit current-task user request for agy coding delegation, not for ordinary Codex subagents, research, or conversation.
---

# agy Subagent Router

Use this Skill as middleware only after an explicit current-task user request selects `agy.exe` for bounded coding delegation. Agent autonomy, provider availability, or a Codex failure is not sufficient authorization. Do not use it for research, ordinary conversation, external messaging, or final verification. Keep evaluation and routing evidence in the Codex-side evaluation artifact; do not make agy configuration the source of truth.

## Resolve before invoking

1. Run `agy.exe models` immediately before model selection.
2. If discovery fails or the intended model is absent, return `FAILED` with a non-secret reason. Do not guess a model identifier or silently substitute a provider.
3. Select the returned model for the bounded coding role and difficulty according to `agy.lrf`, the active subagent LRF, and the Codex evaluation artifact.
4. Map supported agy reasoning to `--effort low|medium|high`; record any requested depth that agy cannot express as unresolved rather than pretending it was applied.

## Invoke a bounded task

Use non-interactive print mode with an exact model, effort, and safe mode:

```powershell
agy.exe --print --model <resolved-id> --effort <low|medium|high> --mode plan --output-format json --print-timeout 5m --prompt <bounded-task>
```

Use `--mode accept-edits` only for an explicitly requested, bounded workspace edit with exclusive file ownership. Never use `--dangerously-skip-permissions`. Pass target, acceptance, verified facts, constraints, allowed files, and the required compact output schema. Do not pass secrets, raw logs, or unrelated context.

## Return and evaluate

Normalize the result to `SUCCESS`, `FAILED`, `ADVICE`, or `ESCALATE`. The main Codex agent must inspect any real diff and run its own checks. At completion, record requested and effective provider, model, reasoning, outcome, quality, acceptance result, runtime when available, and rerouting advice in `C:\Users\letwir\.codex\prompt-ref\eval\subagents_model_eval.json`.
