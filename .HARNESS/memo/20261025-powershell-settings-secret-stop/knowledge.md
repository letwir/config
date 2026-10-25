# Knowledge Base — PowerShell settings export blocked — 20261025-powershell-settings-secret-stop

- task_id: `20261025-powershell-settings-secret-stop`
- date: 2026-10-25
- target environment: Windows PowerShell configuration on Windows; destination is the config repository's `powershell/` directory.
- evidence/source scope: filenames and syntax-only scans of the current PowerShell profile, Starship initializer, and `ojosudo.ps1`; no source values were emitted.
- verification status: both Harness `change` preflights reported PASS; PowerShell parser reported zero syntax errors in the three scanned files.
- remaining uncertainty: the credential-like marker in `profile.ps1` was not inspected or copied; the personal-path marker in `starship_init.ps1` was not inspected or copied. Harness secret-found policy requires stopping for user review.

## Findings

- The requested destination directory did not exist at inspection time.
- A value-free scan found a credential-related marker at `profile.ps1:13` and a personal-path marker at `starship_init.ps1:13`. Their contents were not printed, retained, or copied into this report.
- No PowerShell settings were copied into the requested destination; no installer or installation guide was created.
- The approved Harness profile changes are the only config-repository edit made before the stop: home and project Hardware profiles were set to `windows-pwsh`.
- No commit or push was performed.

## Constraints

- Do not reproduce, disclose, or ingest source credential values or personal paths.
- Resume only after the user has reviewed the source marker and explicitly directs a safe exclusion/redaction route.
