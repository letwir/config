$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
$profiles = @(
    @{ Path = 'agents/pi/attacker-low.md'; Name = 'attacker-low'; Model = 'openai-codex/gpt-6-luna'; Thinking = 'low' },
    @{ Path = 'agents/pi/attacker-medium.md'; Name = 'attacker-medium'; Model = 'openai-codex/gpt-6.1-sol'; Thinking = 'high' }
)
foreach ($agentProfile in $profiles) {
    $text = Get-Content -Raw (Join-Path $root $agentProfile.Path)
    foreach ($expected in @("name: $($agentProfile.Name)", 'advertise: true', "model: $($agentProfile.Model)", "thinking: $($agentProfile.Thinking)", 'tools: read', 'extensions:', 'acceptanceRole: read-only', 'maxSubagentDepth: 0', 'agents/blackhat.md', 'agents/BLACKHAT.css', 'agents/SUBAGENTS.lrf', 'POLYGLYPH.md', 'rules/BOOTSTRAP.lrf')) {
        if (-not $text.Contains($expected)) { throw "$($agentProfile.Path): missing '$expected'" }
    }
    if ($text -notmatch '(?s)\A---\r?\n(?<frontmatter>.*?)\r?\n---\r?\n') { throw "$($agentProfile.Path): invalid frontmatter delimiters" }
    $frontmatter = $Matches.frontmatter
    $keys = @([regex]::Matches($frontmatter, '(?m)^([A-Za-z][A-Za-z0-9]*):').Groups | Where-Object { $_.Name -eq '1' } | ForEach-Object { $_.Value })
    if ($keys.Count -ne (@($keys | Select-Object -Unique)).Count) { throw "$($agentProfile.Path): duplicate frontmatter keys" }
    if ($frontmatter -notmatch '(?m)^acceptanceRole: read-only\r?$') { throw "$($agentProfile.Path): invalid acceptanceRole" }
    if ($frontmatter -notmatch '(?m)^maxSubagentDepth: (0|[1-9][0-9]*)\r?$') { throw "$($agentProfile.Path): maxSubagentDepth must be a non-negative integer" }
    if ($text -notmatch '(?m)^tools: read\r?$') { throw "$($agentProfile.Path): tools must be read-only" }
    if ($text -notmatch '(?m)^extensions:\s*\r?$') { throw "$($agentProfile.Path): extensions must be empty" }
    if ($text -notmatch 'no children') { throw "$($agentProfile.Path): missing no-children contract" }
}
$manual = Get-Content -Raw (Join-Path $root 'rules/MANUAL.lrf')
$subagents = Get-Content -Raw (Join-Path $root 'agents/SUBAGENTS.lrf')
$router = Get-Content -Raw (Join-Path $root 'agents/AGENT_ROUTER.md')
$blackhat = Get-Content -Raw (Join-Path $root 'agents/blackhat.md')
foreach ($source in @(@{Name='MANUAL';Text=$manual}, @{Name='SUBAGENTS';Text=$subagents}, @{Name='AGENT_ROUTER';Text=$router}, @{Name='blackhat';Text=$blackhat})) {
    foreach ($expected in @('openai-codex/gpt-6-luna', 'openai-codex/gpt-6.1-sol', 'host', 'never fallback')) {
        if (-not $source.Text.Contains($expected)) { throw "$($source.Name): missing routing invariant '$expected'" }
    }
}
if ($router -notmatch 'hard=attacker-hard\(openai/gpt-6-astra,provider-default\)') { throw 'AGENT_ROUTER: Astra route changed or missing' }
if ($blackhat -notmatch 'attacker-hard=openai/gpt-6-astra/provider-default') { throw 'blackhat: Astra route changed or missing' }
Write-Output 'PASS: Pi attacker profiles, read-only constraints, contract references, and host-selected routing consistency.'
