<#
.SYNOPSIS
    Chrome DevTools Protocol (CDP) Search & Web ETL Wrapper CLI.
.DESCRIPTION
    Launches headless Chrome to search local files (~/.gemini, ~/.opencode, ~/.codex)
    or live Web search/scraping via CDP V8 and DOM evaluation, returning structured JSON.
.EXAMPLE
    .\chrome_cdp_search.ps1 -WebQuery "Chrome DevTools Protocol" -JsonOnly
    .\chrome_cdp_search.ps1 -Target "$HOME\.gemini" -Query "model"
    .\chrome_cdp_search.ps1 -Target "https://example.com" -Selector "h1, p"
#>
[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [string]$Query,

    [Parameter(Position = 1)]
    [string]$Target = "$HOME\.gemini",

    [string]$WebQuery,
    [string]$WebEngine = "bing",
    [string]$Selector,
    [string]$Ext = "json,jsonc,md,toml,js,ts,html,txt,css",
    [int]$MaxFiles = 30,
    [switch]$JsonOnly
)

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$MjsPath = Join-Path $ScriptDir "chrome_cdp_search.mjs"

$NodeArgs = @($MjsPath)
if ($WebQuery) {
    $NodeArgs += @("--web-query", $WebQuery, "--web-engine", $WebEngine)
} else {
    if ($Query) { $NodeArgs += @("--query", $Query) }
    if ($Target) { $NodeArgs += @("--target", $Target) }
}
if ($Selector) { $NodeArgs += @("--selector", $Selector) }
if ($Ext) { $NodeArgs += @("--ext", $Ext) }
if ($MaxFiles) { $NodeArgs += @("--max-files", $MaxFiles.ToString()) }
if ($JsonOnly) { $NodeArgs += @("--json-only") }

& node @NodeArgs
