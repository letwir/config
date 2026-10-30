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

Purpose: current∧cited guidance(OpenAI product, API, model, Codex). Primary reference読込数∈{0,1}.

**First substantive action:** Search the user's exact requested official OpenAI documentation topic and any explicitly named model using a concise, topic-specific query of 2-6 essential terms. When an already-available direct official documentation search and page-retrieval capability is present, use it first: search, then fetch or open the matching official page before general web search. Otherwise, immediately use official-domain web search, then actually open or fetch the relevant official page. Complete this source order before reading a reference, inspecting local or repository files, running a Codex manual or model resolver, drafting a plan, or answering from memory. Use the actual fetched page, not a search snippet or an unopened link. If one official search or page does not establish the answer, search another appropriate official domain and actually open or fetch the result. Preserve the exact requested model; never substitute a newer model.

**Only exception:** An explicitly requested, genuinely broad, cross-topic Codex setup, orientation, or system-map synthesis may use the manual first when shell execution and an allowed temporary cache are available. A specific Codex feature, setting, command, error, model, or requested citation remains docs-first. Mixed Chat/Work/Codex comparisons are official documentation questions, not manual-first Codex requests.

Generic software task⇒直接そのtaskに回答. OpenAI implementation/debugging/SDK/API/prompting/agent/eval requests≠generic.

Straightforward factual∨citation-only request⇒source orderを守る∧route reference読込禁止. Includes{straightforward API facts, ChatGPT Work∨mixed Chat/Work/Codex comparisons, model tiers, aliases, Pro mode, reasoning settings, factual migration baselines, narrow Codex facts}. ChatGPT Work⇒`learn.chatgpt.com`優先.

## Choose one primary route

Routing=first match; reference読込 only if requested task needs its specialized workflow:

- **Explicitly requested local documentation integration:** user explicitly requests that local integration⇒only then read [integration guidance](references/mcp-diagnostics.md).
- **Model migration, upgrades, or model-specific prompting:** actual migration planning/implementation/dynamic target resolution/prompt changes⇒read [model-migration.md](references/model-migration.md). Explicitly requested target=保持.
- **Model selection and comparisons:** nuanced current/latest/default/cost/latency/quality/modality tradeoffs need more guidance⇒only then read [model-selection.md](references/model-selection.md). Selection alone⇒migration resolver実行禁止.
- **Product, API, ChatGPT Work, and mixed Chat/Work/Codex documentation:** fetched official pages leave source selection/API schemas/requested implementation unresolved⇒only then read [official-docs.md](references/official-docs.md). This route≠manual-first.
- **Explicitly broad Codex setup, orientation, or cross-topic synthesis:** eligible Codex manual∨deeper Codex procedures needed⇒read [codex-self-knowledge.md](references/codex-self-knowledge.md).

Primary reference≤1. 全routes/bundled model guides/helper scriptsの一括読込禁止. Supporting reference読込∨helper実行 only if chosen workflow demonstrably needs it.

## Source and execution boundaries

- Search, open, fetch, and cite only `developers.openai.com`, `platform.openai.com`, and `learn.chatgpt.com`. Cite the page that supports the claim. State uncertainty when official sources do not establish pricing, availability, account access, limits, or behavior.
- Preserve an explicitly requested model for selection, migration, and prompting. Resolve an unspecified latest or current migration target only after searching and fetching current official guidance.
- Use `references/latest-model.md` only as a disclosed fallback after current official model guidance does not answer the question. Read `references/upgrading-to-gpt-6-astra.md` only for an actual, requested GPT-6 migration; read `references/prompting-guide.md` only for requested prompting work.
- Before building, running, editing, debugging, or testing an API-backed app or tool, use `openai-platform-api-key` first when available. Documentation, conceptual examples, model selection, and read-only guidance do not require an API key.
- Say "OpenAI Docs" or "official OpenAI documentation" in user-facing answers. Keep exact official citations and examples concise.
