# H-05 Diary

### 2026-09-09

**Hypothesis:** Existing declarations could be made reproducible with a deterministic consumer and receipt.

**Tried:** Checked `harness-lint`, separated input/result schemas, parsed live LRF and `main.seq`, added isolated failure fixtures, and ran the live oracle.

**Rejected:** Treating oracle expectations as authority; inferring effects from request prose; applying every `MUST_NOT` record as a blanket effect prohibition; retaining exploratory receipts.

**Uncertainty:** The evaluator validates policy mechanics, not semantic natural-language effect classification.

**Attribution:** PromptDefect not found. AgentDefect found and corrected: the first implementation ignored explicit authorization and did not compare oracle expectations.

**Search correction:** Existing-runner search showed no consumer, so a focused runner was justified.

**Emotion and thoughts:** The independent oracle boundary mattered more than adding case volume.