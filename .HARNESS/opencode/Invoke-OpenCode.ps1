[CmdletBinding()]
param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$OpenCodeArguments
)

$ErrorActionPreference = 'Stop'
$expectedVersion = '1.18.25'
$harnessRoot = Split-Path -Parent $PSCommandPath
$configPath = Join-Path $harnessRoot 'opencode.jsonc'
$command = Get-Command opencode -ErrorAction SilentlyContinue

if ($null -eq $command) {
    throw 'OpenCode CLI was not found on PATH. Install or expose OpenCode 1.18.25, then retry.'
}

$actualVersion = (@(& $command.Source --version 2>&1)[0]).ToString().Trim()
$normalizedVersion = $actualVersion -replace '^v', ''
if ($normalizedVersion -ne $expectedVersion) {
    throw "OpenCode version mismatch: expected $expectedVersion, found $actualVersion. Re-convert against that version before use."
}

$previousConfig = $env:OPENCODE_CONFIG
$previousConfigDir = $env:OPENCODE_CONFIG_DIR
try {
    $env:OPENCODE_CONFIG = $configPath
    $env:OPENCODE_CONFIG_DIR = $harnessRoot
    & $command.Source @OpenCodeArguments
}
finally {
    $env:OPENCODE_CONFIG = $previousConfig
    $env:OPENCODE_CONFIG_DIR = $previousConfigDir
}
