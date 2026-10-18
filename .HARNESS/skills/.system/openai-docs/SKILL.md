---
name: "openai-docs"
description: "Use for Codex models/pricing, scheduled tasks, skills, settings, setup, troubleshooting, customization, automations, and self-knowledge—including 'you,' 'your,' 'this app,' or 'this coding agent' when they refer to Codex—and for OpenAI APIs/products and ChatGPT Work. Also use for model choice/migration, prompting, SDKs, Responses, Realtime, agents, evals, and Chat/Work/Codex comparisons. Do not use for generic app/software tasks that merely mention Codex."
metadata:
  short-description: "Codex models/pricing, scheduled tasks, skills, settings, setup, troubleshooting, and self-knowledge; OpenAI APIs and ChatGPT Work. 'You'/'this app' means Codex only."
---

@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# OpenAI Docs

🎯 Output:=current+cited(OpenAI product∪API∪model∪Codex guidance); primary_reference∈{0,1}.

🔎 First substantive action:=query(exact_user_topic∪explicit_model, concise 2–6 essential terms) against official OpenAI docs.
⊢ If direct official search+page-retrieval capability already exists: search→fetch/open matching official page before general web search.
⊢ Else: immediately official-domain web-search→actually fetch/open relevant page.
→ Complete source-order before references, local/repo inspection, Codex manual/model resolver, planning, or memory-based answer. Evidence=actually fetched page, ¬snippet/ unopened link. If insufficient, search another suitable official domain and fetch/open it. Preserve requested model exactly; ¬newer-model substitution.

⚠️ Sole exception: explicitly requested + genuinely broad + cross-topic Codex setup/orientation/system-map synthesis may use manual-first iff shell execution∧allowed temporary cache available. Specific feature/setting/command/error/model/citation⇒docs-first. Mixed Chat/Work/Codex comparison⇒official-docs question, ¬manual-first.

Generic software task⇒answer software task directly. OpenAI implementation/debugging/SDK/API/prompting/agent/eval request¬generic.

Straight factual∨citation-only: follow source order; route-reference=none. Includes simple API facts, ChatGPT Work or mixed comparisons, model tiers/aliases/Pro mode/reasoning settings/factual migration baselines/narrow Codex facts. ChatGPT Work source priority=`learn.chatgpt.com`.

## Choose one primary route

Route:=first match; open its reference iff specialized workflow needed:
- Explicit local-doc integration requested → [integration guidance](references/mcp-diagnostics.md) only iff explicitly requested.
- Migration/upgrades/model-specific prompting → [model-migration.md](references/model-migration.md) for actual planning/implementation/dynamic target/prompt changes; retain explicit target.
- Selection/comparison → [model-selection.md](references/model-selection.md) iff nuanced current/latest/default/cost/latency/quality/modality tradeoffs needed; selection alone⇒¬migration resolver.
- Product/API/ChatGPT Work/mixed comparisons → [official-docs.md](references/official-docs.md) iff fetched official pages leave source selection/API schema/requested implementation unresolved; ¬manual-first.
- Explicit broad Codex setup/orientation/cross-topic synthesis → [codex-self-knowledge.md](references/codex-self-knowledge.md) iff eligible manual/deeper Codex procedures needed.

Primary references≤1. ¬open every route/bundled model guide/helper. Supporting reference or helper only if chosen workflow demonstrably requires it.

## Source and execution boundaries

- 🔎/open/fetch/cite domains only: `developers.openai.com`, `platform.openai.com`, `learn.chatgpt.com`. Cite claim-supporting page. If official sources do not establish pricing/availability/account access/limits/behavior⇒state uncertainty.
- Preserve explicit model in selection/migration/prompting. Unspecified latest/current migration target⇒resolve only after search+fetch current official guidance.
- `references/latest-model.md`⇒disclosed fallback only after current official model guidance insufficient. `references/upgrading-to-gpt-6-astra.md`⇒actual requested GPT-6 migration only; `references/prompting-guide.md`⇒requested prompting only.
- Before build/run/edit/debug/test of API-backed app/tool, use `openai-platform-api-key` first iff available. Documentation/conceptual examples/model selection/read-only guidance⇒API key not required.
- User-facing answers name "OpenAI Docs" or "official OpenAI documentation". Exact official citations/examples concise.
