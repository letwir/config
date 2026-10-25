# PowerShell Profile Configuration

> **Privacy:** User-specific paths are replaced with `$PSScriptRoot`, `$HOME`, runtime command discovery, or relative paths. No local credentials are stored here. Keep machine-specific values in environment variables and out of Git.

A portable Windows PowerShell profile bundle. The profile loads Starship setup and the `ojosudo` helper from paths relative to this directory; optional external commands are discovered at runtime instead of using machine-specific executable paths.

## Included files

- `profile.ps1` — profile functions, aliases, PSReadLine settings, and optional Carapace initialization.
- `starship_init.ps1` — Starship prompt and PSReadLine integration.
- `ojosudo.ps1` — optional elevated-command helper.
- `Install.ps1` — backup-aware installer for the current user's PowerShell profile directory.

The separate local `Scripts/Invoke-MySqlQuery.ps1` and installed `Modules/` tree are intentionally not included: they are not loaded by this profile and are outside the portable shell-configuration bundle.

## Requirements

- Windows PowerShell 5.1 or PowerShell 7; PowerShell 7 (`pwsh.exe`) is recommended.
- Starship and Carapace are optional. Their initializers are used only when the corresponding executable is available on `PATH`.
- Aliases for 7-Zip, FLAC, Android platform tools, Java, ffmpeg, Sakura, NKF, and related utilities are optional; they are created only when their executables are discoverable on `PATH`. Java can also be found under `JAVA_HOME`.
- `cdwsl` requires WSL. It uses `POWERSHELL_WSL_DISTRO` when set and otherwise uses the `Ubuntu` distribution; the Linux home directory is resolved at runtime.
- `aozoraepub` requires Java and `AOZORAEPUB_JAR`. A relative `AOZORAEPUB_JAR` value is resolved relative to this bundle; no jar or machine-specific path is included.

## Install

From this directory, preview the changes first:

```powershell
pwsh.exe -NoProfile -File .\Install.ps1 -WhatIf
```

Then install:

```powershell
pwsh.exe -NoProfile -File .\Install.ps1
```

The installer targets the directory containing the current PowerShell host's `$PROFILE.CurrentUserAllHosts`. Existing files with the same names are backed up beside the originals with a `.bak.<timestamp>` suffix before replacement. No administrator privileges are required.

To test the installer without changing the active profile, choose a temporary destination:

```powershell
$testProfile = Join-Path $env:TEMP 'pwsh-profile-test'
pwsh.exe -NoProfile -File .\Install.ps1 -DestinationDirectory $testProfile
```

After installation, open a new PowerShell session. The profile's optional commands remain inactive when their dependencies are not installed or discoverable.

## Local settings and validation

Keep private values out of the repository. Set local values in the process/user environment as needed; for example, set `JAVA_HOME`, `POWERSHELL_WSL_DISTRO`, or `AOZORAEPUB_JAR` outside these tracked files.

Validate syntax without executing profile startup actions:

```powershell
$files = 'profile.ps1', 'starship_init.ps1', 'ojosudo.ps1'
foreach ($file in $files) {
    $tokens = $null
    $errors = $null
    [System.Management.Automation.Language.Parser]::ParseFile((Join-Path $PWD $file), [ref]$tokens, [ref]$errors) | Out-Null
    if ($errors.Count) { throw "Syntax errors in $file" }
}
```

To restore a previous profile, close PowerShell, remove the installed file, and rename the matching `.bak.<timestamp>` file back to its original name.
