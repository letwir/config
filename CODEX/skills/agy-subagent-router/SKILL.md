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

## Construct the external-agent handoff

Build the handoff in Codex before invoking `agy.exe`. Resolve the nearest applicable project `AGENTS.md`, compile its routed rules for the current task, and select exactly one matching role from `C:\Users\letwir\.codex\agents`. Extract only the role and project contracts that affect the bounded task. The external model must not discover, choose, reinterpret, or expand its own authority.

[SPR/XML::ρ→max|protocol:select⇒project⇒handoff⇒invoke|cues:🔎⊕⊢⊕→⊕⛔⊕✅]
Obj(TaskRequest ∘ ProjectRules ∘ AgentRole) → Mor(ResolveModel ∘ BoundHandoff ∘ InvokeAgy) → Obj(StructuredResult ∘ Evidence)
<ExternalAgentHandoff>
role: one selected agent role;
target: concrete requested outcome;
acceptance: observable completion checks;
scope: workspace root and exact allowed files;
known_facts: verified task facts only;
project_contracts: task-matching conventions selected from the nearest applicable project rules;
constraints: applicable prohibited effects, ownership, timeout, and non-interference requirements;
output: SUCCESS | FAILED{cause,evidence} | ADVICE{action} | ESCALATE{reason};
</ExternalAgentHandoff>

@lrf=1|aud=GPT-5.6|scope=agy-external-agent-handoff
R|handoff.select|task:agy|MUST|RO_LOCAL|selector=Codex; sources=nearest-applicable-AGENTS.md+compiled-task-rules+exactly-one-selected-agent-role; external-self-selection=false
R|handoff.project|task:agy|MUST|RO_LOCAL|include=task-matching-project-conventions; omit=unrelated-rules,full-conversation,other-roles,raw-evaluation-history
R|handoff.fields|task:agy|MUST|RO_LOCAL|required=role,target,acceptance,scope,allowed-files,known-facts,project-contracts,constraints,output-schema
R|handoff.evidence|task:agy|MUST|RO_LOCAL|known-facts=verified-only; inference=labelled; unknown=explicit; completion-claim=acceptance-evidence-required
R|handoff.secrets|task:agy|MUST_NOT|CRED|transmit=credentials,tokens,private-environment-values,raw-secret-logs
R|handoff.enforcement|task:agy|MUST|RO_LOCAL|prompt-is-not-enforcement; caller-enforces=allowed-files,process-mode,timeout,cancellation,prohibited-effects
R|handoff.effects|task:agy|MUST|RO_LOCAL|external-agent-cannot-expand-current-task-authorization; unclassified-or-conflicting-effect=STOP
R|handoff.verify|task:agy|MUST|RO_LOCAL|owner=main-Codex; inspect=real-output+real-diff; run=proportionate-checks; external-self-report=insufficient

## Invoke a bounded task

Use non-interactive print mode with an exact model, effort, and safe mode:

```powershell
agy.exe --print --model <resolved-id> --effort <low|medium|high> --mode plan --output-format json --print-timeout 5m --prompt <bounded-task>
```

Use `--mode accept-edits` only for an explicitly requested, bounded workspace edit with exclusive file ownership. Never use `--dangerously-skip-permissions`. Serialize the constructed `ExternalAgentHandoff` into the prompt without adding omitted context. Do not pass secrets, raw logs, or unrelated context.

## Return and evaluate

Normalize the result to `SUCCESS`, `FAILED`, `ADVICE`, or `ESCALATE`. The main Codex agent must inspect any real diff and run its own checks. At completion, record requested and effective provider, model, reasoning, outcome, quality, acceptance result, runtime when available, and rerouting advice in `C:\Users\letwir\.codex\prompt-ref\eval\subagents_model_eval.json`.
