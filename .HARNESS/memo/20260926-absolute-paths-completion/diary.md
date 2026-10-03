# Diary: Absolute-Path Cleanup Completion

### 2026-09-26
- **Hypothesis:** Maintained prompt/config references can use home/environment-derived paths while preserving their role.
- **Tried:** Updated the approved active files; parsed the changed PowerShell, JavaScript, and TOML; ran focused LRF lint and path searches.
- **Rejected:** Rewriting exact command patterns in `rules/default.rules`, because literal replacement would alter command matching; rewriting historical records, because that would modify recorded evidence.
- **Uncertainty:** The optional `LLM_MEMORY_REPO` environment variable may be unset in some environments; existing executable and `LLM_MEMORY_BIN` candidates remain.
- **Attribution:** Initial stop was caused by conflicting reusable instructions; the user clarified the condition for unsettled versus settled handoffs, after which implementation resumed.
- **PromptDefect vs AgentDefect:** PromptDefect 0%; AgentDefect 0% for the completed change. Earlier equal-rank workflow conflict was resolved by the user.
