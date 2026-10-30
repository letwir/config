# diary / 20261025-crabcode-config-precedence

- timestamp: 2026-10-25 13:11 local
- task: Investigate whether OpenCode configuration is taking precedence over Crabcode configuration.
- request-evidence: User asked in Japanese to investigate the apparent priority of OpenCode settings over Crabcode settings.
- action: Read the current Crabcode/OpenCode configs, checked relevant environment overrides and config-directory presence, inspected Crabcode 0.0.14 source and official config documentation, and searched prior memory. No config was edited and no live model request was made.
- result: `XDG_CONFIG_HOME`, `CRABCODE_CONFIG`, `OPENCODE_CONFIG` unset; `~/.config/opencode/opencode.jsonc` exists; `~/.config/crabcode` is absent; intended file exists only at `~/.harness/crabcode/crabcode.jsonc`, outside Crabcode 0.0.14's standard candidates. This explains why shared OpenCode config can apply. Correctly placed Crabcode-global values would normally override conflicting OpenCode-global values; project-level OpenCode can override global Crabcode.
- friction: I initially used incorrect shell quoting for a PowerShell timestamp/version probe and first invoked a fallback memory binary before checking `LLM_MEMORY_BIN`. Both were corrected; the configured `LLM_MEMORY_BIN` path does not exist, and the available local fallback binary successfully returned search results. No secret-bearing output was requested or recorded.
- impact: Read-only investigation only. App configuration/provider/model state and VCS were untouched.
- prompt feedback: PromptDefect 0%; the request was sufficiently clear to investigate.
- agent feedback: AgentDefect 100% of observed execution friction; correction is to validate shell quoting and the configured binary path before fallback.
- feedback: No user clarification needed; current runtime CWD/effective config remains unverified.
- rewritten-request: Unchanged — investigate why Crabcode appears to use OpenCode settings first.
