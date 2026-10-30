# walkthrough: 20261025-pi-worker-registration

Request: Make implementation workers usable by Pi from ~/.harness/agents, before continuing FFmpeg updater work.

Verified facts:
- Pi 1.4.0 README and pi-subagents agent/config/model docs inspected.
- Pi settings currently have no subagents block. Supported subagents.agentScanDirs can recursively scan an extra agent root.
- ~/.harness/agents/worker.md exists as a shared contract with OpenCode-style frontmatter; static worker-easy/medium/hard Pi definitions do not exist in the inspected root.
- Live `pi --list-models` shows openai-codex/gpt-6-luna and openai-codex/gpt-6.1-sol. No exact openai/gpt-6-luna or openai/gpt-6.1-sol match was returned.
- Governing worker routing explicitly specifies provider openai, with Luna low/medium and Sol high. It requires exact provider/model agreement and fail-closed handling of unavailable routes. Do not silently translate to openai-codex.

Proposed config-only approach, pending provider decision: Store three Pi-compatible profiles in ~/.harness/agents/pi, reference the shared worker.md contract, grant only the needed built-in file/check tools without nested delegation, and append that narrow directory to subagents.agentScanDirs. Avoid scanning all Harness prose as agent definitions. Preserve all existing settings and source contracts. Verify discovery and actual model on a bounded smoke task once the mapping is approved.

Status: BLOCKED awaiting current-task permission to use the available openai-codex provider for the same model ids and difficulty/thinking mapping. No configuration or product edits, model calls, credential use, external writes beyond standing report-memory ingestion, or FFmpeg changes. No worker run or independent review claimed.
