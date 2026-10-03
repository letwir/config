### 2026-09-25 00:50:40 +09:00
- task: Add static linting for `.seq` files and identify their DSL dialect on line 1.
- request-evidence: User asked for `.seq` lint and a dialect name recognizable by an LLM from the first line; user pointed to `repo/repo/mcp2go` as the Go source location.
- action: Read the active SIGMA workflow and existing CLI; added HSEQ (`harness-seq/1`) validation and CLI dispatch, migrated the five local workflow sequences, documented the format, built the installed binary, and linted the real sequence directory.
- result: Build and both targeted lint invocations exited 0; the sequence directory produced no diagnostics. The first line now identifies the dialect. No runtime execution was claimed.
- friction: The initial source path was outside writable roots; user moved `mcp2go` into the permitted scope. The second AGY call returned API 503 after its file changes were present. No backend switch was made. The source tree was already untracked and other workspace files were already dirty.
- attribution: PromptDefect=0; AgentDefect=tool availability failure (AGY API capacity 503); implementation decisions followed the current request and local LRF contracts.
- impact: `harness-lint` now provides structural checks for `.seq`; the current five workflow files pass strict lint.
- feedback: None required; the API capacity failure was transient infrastructure friction and did not require a prompt change.
- rewritten-request: Add HSEQ lint to the existing CLI, make `@dialect harness-seq/1` the first line of each workflow file, preserve existing sequence versions and unrelated edits, and verify with a local build and strict lint.
