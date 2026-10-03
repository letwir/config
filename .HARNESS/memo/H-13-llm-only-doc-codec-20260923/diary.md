### 2026-09-23 — LLM-only DOC_RULE codec

- Hypothesis: the token-minimal semantic codec should apply only to explicitly LLM-only documents; otherwise it can damage human-facing clarity.
- Tried: added an audience-gated LRF rule and connected the SKILL contract to it.
- Rejected: treating H-readability=empty as permission to weaken authority, safety, exact schemas, or protocol requirements.
- Verification: Level 1 LRF lint and diff whitespace checks passed. Level 2 surfaced an existing diagnostic at line 11.
- Uncertainty: no repository-wide behavior test is applicable to this documentation-only change; this task did not ask for broader verification.
- Attribution: prompt defect none; implementation defect none identified. Existing Level 2 diagnostic is pre-existing.
- Completion: local rule updated; unrelated dirty files preserved; no publication requested.
