# Semantic Review / 20261018-polygllyph-application-resume

## Method and limit

Compared representative current PG bodies against their source versions at baseline `6a8889674bd01d867e9374430dbab46ebd3115b8`; checked Intent, Relations/effects, polarity/negation, modality, exceptions, scope, temporal/threshold, evidence/uncertainty, stop/success, and DecisionBoundary. Frontmatter, code fences, and LRF blocks are frozen by SHA-256 fixture (59/59). This is a bounded human review, not formal proof or a universal-decoding guarantee.

## 1. Agent handoff — `agents/worker.md`

PASS: Bounded CODE role remains one eligible worker; exact `<Γ>` wrapper and ordered fields `Role,Target,Acceptance,Scope,Known facts` remain in the fenced protocol; final `SIGMA/1` position is preserved. Worker role/permission does not grant authority; credential/external/publish/deploy/VCS/production-write prohibitions remain. Agy opt-in, llama single-file/stdout, one-writer, verified-facts, failure reporting and root-owned verification remain explicit. Stop/success is tied to actual diff/check evidence.

## 2. Read-only retrieval — `skills/google-search/SKILL.md`

PASS: Search trigger and query input preserved; DuckDuckGo lite remains primary, HTML fallback only on lite failure. Output is result titles and source URLs; no-hit/error is reported and fabricated claims are forbidden. Existing LRF still defines `RO_PUBLIC`, credential denial and sourced-fact boundary. No external write or API-key authority introduced.

## 3. Live mutation — `skills/llm-memory/SKILL.md`

PASS: Inspection and mutation command sets are distinguished; mutations require current-task approval and verified evidence. `eval` prerequisites/receipt behavior, survey threshold and submit approval, credential guards, and uncertainty are preserved. `survey.classify` remains explicitly unclassified rather than inferred read-only or mutating. Original LRF block bytes remain exact.

## 4. Code/config and exceptions — `skills/.system/plugin-creator/SKILL.md`

PASS: Name normalization/length, default personal destination, repo/team explicit opt-in, marketplace-name exception/no in-place rename, policy enums/defaults, override gating, `--force` conditions, approval-before-write, update/cachebuster flow, validation commands and app handoff conditions are represented as structured relations (`⇔`, `iff`, `⇒`, `MUST`, `MUST_NOT`, named fields). Exact 11 code/JSON fences, commands, YAML and literals match frozen baseline. Critic and final Verifier found no specific semantic loss; remaining English technical terms are allowed by per-concept multilingual encoding.

## 5. Runtime orchestration — `opencode/agents/orchestrator.md`

PASS: Dispatch/status/result/cancel tools, exact role/depth constraints, inherited read-only grandchildren, explicit scopes, wait-and-verify, parent-owned tests/builds, queued natural-turn notifications, authority separation, lease retention on uncertain transport/cancellation, and honest runtime limitation reporting remain encoded. Exact frontmatter/runtime literals preserved.

## Global boundaries

`POLYGLYPH.md` explicitly distinguishes representation from authority and decoding from authorization. SIGMA/1 plus routed LRF remain sole policy authority. LRF `MAY` remains supported, and symbols/emoji are never sole safety carriers. Ambiguity blocks dependent actions and requires uncertainty/reporting. No detected permission/effect expansion in the reviewed cases.
