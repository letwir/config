# knowledge / 20261025-crabcode-config-precedence

- topic: Crabcode / OpenCode config precedence on the current Windows profile
- timestamp: 2026-10-25 13:11 local
- target environment: Windows, user profile `C:\Users\letwir`, Crabcode 0.0.14
- evidence/source scope: local configs, Crabcode 0.0.14 installed crate source, Crabcode/OpenCode official configuration docs; no live Crabcode session was started
- verified facts:
  1. `XDG_CONFIG_HOME`, `CRABCODE_CONFIG`, and `OPENCODE_CONFIG` are unset. The default global config root is therefore `~/.config`.
  2. `~/.config/opencode/opencode.jsonc` exists and defines shared OpenCode-compatible provider/shell/plugin settings. `~/.opencode/opencode.jsonc` also exists and only declares a plugin.
  3. `~/.config/crabcode` and `~/.crabcode` do not exist in the current inventory. The Crabcode file present is `~/.harness/crabcode/crabcode.jsonc`.
  4. Crabcode 0.0.14's resolver source (`.cargo/registry/src/.../crabcode-0.0.14/src/config/configuration.rs`, `resolve_sources`, lines 1018-1084) searches OpenCode global, Crabcode global, OpenCode project, then Crabcode project candidate paths. `~/.harness/crabcode/crabcode.jsonc` is not one of those candidates. Source merge code deep-merges files in resolver order.
  5. Official Crabcode docs state later layers override earlier conflicting keys and object values deep-merge: https://crabcode.rs/config/ (accessed 2026-10-25). Therefore a valid Crabcode-global file should override conflicting OpenCode-global keys; an OpenCode-project file can override Crabcode-global keys.
  6. Official Crabcode compatibility docs say `plugin` and custom-tool settings are intentionally ignored: https://crabcode.rs/config/opencode-compatibility/ (accessed 2026-10-25). Missing plugin behavior is not a precedence problem.
  7. Historical `~/.harness/issues.md` notes report a `.config/crabcode` setup in October, but the current `.config` inventory has no `crabcode` directory; current filesystem evidence supersedes that historical state.
- conclusion: observed OpenCode-over-Crabcode behavior is consistent with Crabcode loading `~/.config/opencode/opencode.jsonc` while not finding the intended config at `~/.harness/crabcode/crabcode.jsonc`. The main issue is path registration/link drift, not a general rule that OpenCode global overrides Crabcode global.
- remaining uncertainty: actual Crabcode launch CWD/project root and resolved runtime config were not observed; project-level OpenCode settings could add another override. No runtime provider/model invocation was made.
- verification status: source/docs/config inventory cross-checked; no application config changes.
