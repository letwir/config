# Diary — PowerShell settings export blocked — 20261025-powershell-settings-secret-stop

- task_id: `20261025-powershell-settings-secret-stop`
- timestamp: 2026-10-25 (system date; exact clock time not captured)
- target environment: Windows PowerShell; destination config repository `powershell/`.
- evidence/source scope: Harness bootstrap/routed rules, change preflight summaries, and value-free PowerShell syntax/structure scan.
- verification status: preflight PASS in both Harness roots; parser syntax errors = 0 for scanned profile files; task blocked before migration.
- remaining uncertainty: exact meaning/content of the credential-like and personal-path markers is intentionally unknown.

## Request evidence

The user requested a sanitized export of their PowerShell setup into `~/repo/repo/config/powershell`, validation, installation instructions, an installer script, and eventual Git push.

## Action and result

Applied the user-approved `windows-pwsh` profile setting to the home and project Harness files. Both change preflights passed. A scanner that emitted only filenames, line numbers, syntax status, and marker counts detected a credential-related marker in `profile.ps1` and a personal-path marker in `starship_init.ps1`. It emitted no source values. Per the routed `secret-found` stop rule, no settings were copied, no deployment script/docs were created, and no VCS operation was performed. The temporary scanner was removed.

## Friction and attribution

- PromptDefect: 0%. The destination and desired security treatment were stated clearly; the blocking source marker was not inferable from the request.
- AgentDefect: present. The main persona was applied after the initial response/action instead of before it, and the first two PowerShell preflight invocations used incorrect quoting; both process issues were corrected before the source scan result was acted upon. They did not expose source values.
- Impact: delivery is incomplete; user review is required before any further source inspection or export.

## Feedback and rewritten request

No prompt rewrite is needed. To resume safely: review and remove or manually redact the credential-like marker in the local profile, then ask to resume while explicitly excluding any remaining credential- or personal-data-bearing lines from the export.
