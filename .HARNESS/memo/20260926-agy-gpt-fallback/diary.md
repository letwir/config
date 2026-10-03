# Diary: agy GPT fallback routing change

- Request: change `FAILED: agy *` handling to GPT fallback and reselect the GPT implementation model by difficulty instead of fixing it to Luna. The user authorized any GPT-family implementation backend for this task.
- Exploration: compiled BOOTSTRAP and routed LRF, inspected AGENT_ROUTER, agy worker contracts, SUBAGENTS, engineering workflow, HSEQ, and the agy skill. Precedent search returned no matching record.
- Implementation: aligned bootstrap, manual, router, agy, engineering, subagent, worker, skill, and CODE sequence contracts. Added a dedicated `worker-gpt-fallback` protocol.
- Verification: contract assertions passed; LRF rules, SUBAGENTS, and HSEQ strict checks passed; `git diff --check` passed. The skill Level 3 Go symmetry check was inapplicable because the documentation-only skill has no Go source.
- Scope: no VCS write, publish, deploy, credential operation, or external message was performed.

