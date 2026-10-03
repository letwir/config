# Knowledge — runtime agent registration — 2026-09-25

## Facts
- Official Codex Configuration Reference defines custom role declarations under `agents.<name>` and `agents.<name>.config_file`; relative config-file paths resolve from the config file declaring the role.
- `C:\Users\letwir\.codex\agents` is a symbolic link to `C:\Users\letwir\.harness\agents`.
- The current `.codex/config.toml` has `[agents] enabled = true` and a concurrency limit, but no `agents.compressor` or `agents.orchestrator` declarations.
- The first bounded worker-agy run applied eight files in `.harness/agents` and related rules: two role TOMLs, two persona CSS files, Worker CSS, AGENT_ROUTER, MANUAL, and SUBAGENTS LRF.
- A follow-up worker-agy generation failed and applied no corrections. The resulting role TOMLs still use `sandbox_mode = "plan"`; the role-specific `⊕` persona summaries need refinement. The global Codex config registry was not modified.
- Targeted MANUAL/SUBAGENTS LRF level-1 checks passed; TOML syntax parsing passed. These checks do not prove custom-agent schema/runtime loading.
- `codex.exe` was not available on PATH; no live runtime registration test was completed.

## Inference
- The prior “unknown agent type” is consistent with the absence of `agents.<name>.config_file` registrations, but a fresh Codex runtime discovery is still needed to confirm the correction.

## Constraints
- Preserve unrelated user changes in both `.harness` and `.codex`.
- Do not register agent entries that point to invalid config layers; complete the role TOMLs and then add the two config-file references.
