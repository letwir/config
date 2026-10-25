# Privacy/portability notice: personal paths and credentials are intentionally omitted.
# Source files are selected relative to this installer; target paths are resolved from the active PowerShell profile.
# Provide local-only values through environment variables; never commit secrets.
#requires -Version 5.1

[CmdletBinding(SupportsShouldProcess = $true, ConfirmImpact = 'Medium')]
param(
    [string]$DestinationDirectory
)

$ErrorActionPreference = 'Stop'
$profileFiles = @('profile.ps1', 'starship_init.ps1', 'ojosudo.ps1')

if ([string]::IsNullOrWhiteSpace($DestinationDirectory)) {
    $profilePath = $PROFILE.CurrentUserAllHosts
    if ([string]::IsNullOrWhiteSpace($profilePath)) {
        throw 'The current PowerShell host did not provide a CurrentUserAllHosts profile path.'
    }
    $DestinationDirectory = Split-Path -Parent $profilePath
}

$missing = @($profileFiles | Where-Object { -not (Test-Path -LiteralPath (Join-Path $PSScriptRoot $_) -PathType Leaf) })
if ($missing.Count) {
    throw ('Required bundle files are missing: ' + ($missing -join ', '))
}

if (-not (Test-Path -LiteralPath $DestinationDirectory -PathType Container)) {
    if ($PSCmdlet.ShouldProcess('PowerShell profile directory', 'Create destination directory')) {
        New-Item -ItemType Directory -Path $DestinationDirectory -Force | Out-Null
    }
}

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss-fff'
foreach ($name in $profileFiles) {
    $source = Join-Path $PSScriptRoot $name
    $target = Join-Path $DestinationDirectory $name

    if (Test-Path -LiteralPath $target -PathType Leaf) {
        $sourceHash = (Get-FileHash -LiteralPath $source -Algorithm SHA256).Hash
        $targetHash = (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash
        if ($sourceHash -eq $targetHash) {
            Write-Output "Already current: $name"
            continue
        }

        $backup = "$target.bak.$stamp"
        if (Test-Path -LiteralPath $backup) {
            throw "Backup already exists for $name; refusing to overwrite it."
        }
        if ($PSCmdlet.ShouldProcess($name, 'Back up the existing profile file')) {
            Copy-Item -LiteralPath $target -Destination $backup -ErrorAction Stop
        }
    }

    if ($PSCmdlet.ShouldProcess($name, 'Install profile file')) {
        Copy-Item -LiteralPath $source -Destination $target -Force -ErrorAction Stop
        Write-Output "Installed: $name"
    } else {
        Write-Output "Would install: $name"
    }
}
