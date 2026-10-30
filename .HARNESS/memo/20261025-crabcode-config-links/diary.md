# diary / 20261025-crabcode-config-links

- timestamp: 2026-10-25 13:23 local
- task: Place multiple symlinks under `.config/crabcode` to align Crabcode's global config view with `.harness`.
- request-evidence: User explicitly asked to put many symbolic links in `.config/crabcode` pointing into `.harness`.
- action: Confirmed `.config/crabcode` was absent and `XDG_CONFIG_HOME` unset; inspected H-09 guidance and Crabcode config-relative prompt references; created a normal config directory and 17 symlinks to the Harness SSOT. Ran the Windows change preflight and checked link listings/target paths.
- result: All 17 links are present and five representative critical target paths resolve. The canonical config JSONC is linked, not copied; `.config/agents`, OpenCode configuration, and the source `.harness` files were not edited.
- friction: An initial `cmd.exe /c` quoting attempt for `mklink` failed before creating a link. I switched to PowerShell `New-Item -ItemType SymbolicLink`; all intended links were then created. This was agent-side tool syntax friction.
- impact: Crabcode's standard global config path now sees the Harness policy/config tree; this should address the previously found path mismatch. The running application was not restarted or invoked, and no provider request was made.
- attribution: PromptDefect 0%; AgentDefect 100% for the corrected command-quoting friction.
- feedback: No clarification needed. Internal correction: use PowerShell 7 `New-Item` for Windows links when `cmd` quoting is ambiguous.
- rewritten-request: Unchanged — add multiple symlinks in `.config/crabcode` to align it with `.harness`.
