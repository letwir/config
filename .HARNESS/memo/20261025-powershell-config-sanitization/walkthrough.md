# Walkthrough — Portable PowerShell profile — 20261025-powershell-config-sanitization

- task_id: `20261025-powershell-config-sanitization`
- timestamp: `2026-10-25T04:19:35Z`
- target environment: Windows PowerShell 5.1 / PowerShell 7.
- evidence/source scope: three top-level canonical PowerShell files and five repository bundle files.
- verification status: static privacy/syntax scans and isolated installer tests passed.
- remaining uncertainty: no live profile startup, WSL navigation, or optional external command was invoked.

## Overview

Created a sanitized portable bundle under `powershell/` and applied the same privacy/path changes to the three canonical files. Every profile copy starts with a privacy notice; user-specific filesystem paths are not embedded.

## 実行計画

1. Recheck Harness rules and inspect current PowerShell configuration without printing sensitive values.
2. Mask personal path examples and replace fixed executable paths with portable references in canonical and copied files.
3. Add an installer and installation guide.
4. Parse and scan the files, test installation/backups in an isolated temporary directory, and review Git state.
5. Complete task records, then follow the user's requested VCS step after verification.

## Bundle contents

- `powershell/profile.ps1`
- `powershell/starship_init.ps1`
- `powershell/ojosudo.ps1`
- `powershell/Install.ps1`
- `powershell/README.md`

The installer targets the current host's profile folder by default and preserves any differing existing files as timestamped backups. Optional applications must be installed separately and found through `PATH`/environment variables.

## Verification

- Harness preflight: PASS.
- Windows PowerShell 5.1 / PowerShell 7 parse: 0 errors across four `.ps1` files.
- Isolated `ojosudo` message function output: 5 lines on both hosts.
- Temporary installer tests: file hashes, backup restoration, idempotent rerun, and no-write `-WhatIf` all PASS.
- Static bundle audit: no credential signatures, personal paths/emails, or static absolute file paths; no trailing whitespace.
- Canonical/copy equality and top-of-file note coverage: 3/3 each.
- `git diff --check`: PASS after `.gitattributes` normalized PowerShell script line endings.
- No commit/push has yet been performed.

## Residual risk

Actual profile startup and WSL/tool integrations remain untested to avoid invoking user shell actions. The local standalone database query script and historical backups remain outside the export. An earlier broad search emitted an absolute path in tool output; it was not written into repository files or these records.
