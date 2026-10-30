# walkthrough / 20261025-crabcode-config-precedence

## 実行計画
1. Locate the Crabcode and OpenCode configs and check current config-related environment overrides.
2. Compare installed Crabcode 0.0.14 resolver source with official Crabcode/OpenCode docs.
3. Separate path-resolution effects from unsupported plugin/custom-tool features; do not alter settings.
4. Record findings, verification limits, and a bounded recommendation.

## 実施・結果
- Confirmed the current global OpenCode config at `~/.config/opencode/opencode.jsonc` and Crabcode JSONC at `~/.harness/crabcode/crabcode.jsonc`.
- Confirmed `XDG_CONFIG_HOME`, `CRABCODE_CONFIG`, and `OPENCODE_CONFIG` are unset, and no `~/.config/crabcode` directory is present.
- Crabcode 0.0.14 source searches global OpenCode config before global Crabcode config, then project OpenCode before project Crabcode. Later conflicting values win; object values are deep-merged.
- The `.harness/crabcode/crabcode.jsonc` path is not in that resolver's candidates. Thus the current Crabcode process can use OpenCode-global values without seeing the intended Crabcode-specific values.
- Crabcode docs explicitly say `plugin` and custom-tool config are ignored; distinguish that expected unsupported behavior from key precedence.
- Recommendation only (not applied): restore/register the Crabcode config at the supported global path `~/.config/crabcode/crabcode.jsonc` (or a verified link from that path to `~/.harness/crabcode`), then restart Crabcode and verify the resolved provider/model without printing credentials.

## Verification
- Crabcode binary: `crabcode 0.0.14`.
- Local source: `.cargo/registry/src/index.crates.io-1949cf8c6b5b557f/crabcode-0.0.14/src/config/configuration.rs`, resolver candidate paths and merge loop.
- Official docs: https://crabcode.rs/config/ ; https://crabcode.rs/config/opencode-compatibility/ ; https://opencode.ai/docs/config/ (accessed 2026-10-25).
- No runtime config dump/model invocation; actual Crabcode process CWD/project root remains unknown.
- Changed app files: none. Task-scoped memory records are the only writes.
