# Knowledge Base — Portable PowerShell profile — 20261025-powershell-config-sanitization

- task_id: `20261025-powershell-config-sanitization`
- timestamp: `2026-10-25T04:19:35Z` (completion work began at this time; tool clock recorded fractional seconds)
- target environment: Windows PowerShell 5.1 and PowerShell 7; repository bundle under `powershell/`.
- evidence/source scope: top-level local `profile.ps1`, `starship_init.ps1`, and `ojosudo.ps1`; repository bundle files; Harness change preflight and isolated installer/parser tests.
- verification status: complete for sanitized export and installer behavior; no actual profile startup or WSL/third-party command execution.
- remaining uncertainty: runtime integration with local Starship/Carapace/WSL installations was not exercised. Local backup snapshots and the standalone MySQL helper were intentionally excluded and remain outside this export.

## Facts

- The three canonical profile files and repository copies are byte-identical. All have a privacy/portability note at the beginning and UTF-8 BOM for Windows PowerShell 5.1 compatibility.
- Personal path literals were removed from the current profile and sample comments. Bundled files resolve relative to `$PSScriptRoot`; optional tools are resolved through `PATH` or `JAVA_HOME`; the WSL home is resolved at runtime.
- Hardcoded executable locations were replaced by optional command discovery. `AOZORAEPUB_JAR` and `POWERSHELL_WSL_DISTRO` are local environment settings.
- `powershell/README.md` documents prerequisites, installation, backups, local environment settings, and syntax validation. `powershell/Install.ps1` copies the three profile files and backs up differing destination files.
- The local MySQL query script and installed module tree were not copied; they are not loaded by the shell profile and are database-specific.

## Verification evidence

- Project Harness `change` preflight: PASS.
- Windows PowerShell 5.1 and PowerShell 7 parser checks: 0 syntax errors for all four `.ps1` files in the bundle.
- `Write-OjosamaMessage` isolated behavior check: 5 output lines on both hosts.
- Installer test in a temporary directory: installed bytes match, old file backup is restorable, second install is idempotent, and `-WhatIf` performs no writes.
- Final bundle audit: 0 credential signatures, 0 personal-path/email markers, 0 hardcoded absolute-path literals in PowerShell files, 0 trailing-whitespace lines; source/copy equality 3/3; top-note/BOM coverage 3/3.
- `git diff --check`: PASS after `.gitattributes` normalized PowerShell script line endings (`*.ps1 text eol=crlf`). No commit or push is represented in this record.

## Constraints

- Never commit local credentials or machine-specific values; use environment variables outside versioned files.
- Do not export backups or database-connection utilities without a separate scope and redaction review.
