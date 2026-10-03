### 2026-09-25

Hypothesis: Pointing `.codex\AGENTS.md` at `.harness\AGENTS.md` will make the harness file the shared policy source.
Tried: Inspected both files, saved the previous Codex file as a dated backup, created the symbolic link, and rechecked the link target and backup hash.
Rejected: Overwriting the existing contents without a recoverable copy.
Uncertainty: Fresh Codex session loading behavior was not tested.
Attribution: PromptDefect=0%; AgentDefect=100% for an initial memory-search call using a stale executable path; the configured executable path was found and used for final ingestion.
Search: Existing local precedent did not return a hit; workspace memory notes and current filesystem state were inspected.
Correction: Resolve `llm-mem.exe` from `LLM_MEMORY_BIN` rather than assuming the older plugin path.
Emotion: N/A.
Thoughts: The filesystem link is verified; instruction-chain reload needs a fresh session for runtime confirmation.
