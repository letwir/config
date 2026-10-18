# knowledge / 20261018-crabcode-config-sync-push

- topic: Crabcode and `.harness` settings sync attempt
- facts:
  1. The active Crabcode config directory resolves to `~/.harness/crabcode`; its shared agents, rules, ETL, scripts, persona, and skills are symlinks to common `.harness` content.
  2. The repository's `.HARNESS/crabcode` does not yet mirror those runtime references. The current config contains a credential-like value and a private-looking endpoint; neither value was copied or recorded.
  3. Crabcode documents OpenCode-compatible configuration. OpenCode's config docs support `{env:VARIABLE}` substitutions for configuration values, including provider keys: https://opencode.ai/docs/config/ (checked 2026-10-18). Crabcode compatibility reference: https://crabcode.rs/config/opencode-compatibility/ (checked 2026-10-18).
  4. The required worker dispatch was unavailable in this execution environment; no implementation or Git push occurred.
- verified-date: 2026-10-18
- scope: `repo/config/.HARNESS`
