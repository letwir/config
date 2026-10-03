# walkthrough | 20260923-worker-backend-routing-design

## Design update
- Treat Worker as a task role with two explicit execution routes: `worker-agy` and `worker-llama`.
- `AGENT_ROUTER.md` selects the route by task shape, difficulty, tool/edit needs, and evaluated performance.
- AGY route: bounded workspace edit with current model discovery; one retry only for launch failure; stop on auth, policy, approval, or post-launch backend failures; inspect partial diff before any fallback.
- Llama route: local llama.cpp-compatible code generation; default output is pure code to stdout. Main agent applies output only within the agreed file allowlist, then reviews diff and verifies. Pipe-to-file requires an explicit bounded route.
- Do not silently treat llama as an automatic fallback for AGY failures. If both are eligible, select before execution; parallel use only across disjoint files with explicit independent workstreams.
- Keep Codex custom-agent TOMLs distinct from external CLI backend contracts. Only choose a TOML where an actual Codex subagent runtime is used; AGY and llama2coder may instead be direct executor adapters with their own Markdown protocols.
- Existing root-only nested dispatch finding remains: main agent owns backend selection and invocation.
- No implementation or tests performed.
