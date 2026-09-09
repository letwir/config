# H-06 Diary

### 2026-09-09

**Hypothesis:** A bounded manifest runner could close the gap between distribution and actual invocation without overclaiming client discovery.

**Tried:** Separated executable, catalog, and representative status; resolved final paths; recorded hashes and versions; bounded and redacted output; killed timed-out process trees; ran safe core operations and fresh client processes.

**Rejected:** Treating help output as representative operation success; treating Codex help as skill-catalog evidence; classifying a timeout as a missing output token; retaining exploratory receipts.

**Uncertainty:** Codex skill discovery remains unverified until a deterministic catalog API is available. `llm-memory` and `agy` representative operations remain untested because their useful operations cross database/network or external orchestration boundaries.

**Attribution:** PromptDefect not found. AgentDefect found and corrected: timeout was initially classified after required-text checking, and transient client failures were not retried.

**Search correction:** The existing sync health check proved file visibility only, so the new runner kept runtime states distinct.

**Emotion and thoughts:** The most useful result was preserving uncertainty explicitly instead of forcing every installed client into Healthy or Failed.
