[CmdletBinding()]
param(
    [string]$ManifestPath = '',
    [string]$UserHome = $env:USERPROFILE,
    [switch]$Json
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
if ([string]::IsNullOrWhiteSpace($ManifestPath)) {
    $ManifestPath = Join-Path (Split-Path -Parent $MyInvocation.MyCommand.Definition) '..\rules\client-distribution.json'
}

function New-Receipt([string]$classification, [string]$name, [string]$detail, [bool]$required) {
    [pscustomobject]@{ Name = $name; Classification = $classification; Required = $required; Detail = $detail }
}

function Resolve-ManifestPath([string]$tokenPath, [string]$manifestRoot, [string]$userRoot) {
    if ([string]::IsNullOrWhiteSpace($tokenPath) -or $tokenPath -match '(^|[\\/])\.\.([\\/]|$)') { throw "Invalid manifest path token: $tokenPath" }
    if ($tokenPath -match '^[A-Za-z]:[\\/]') { throw "Absolute paths are not allowed in manifest: $tokenPath" }
    $relative = $tokenPath -replace '/', '\\'
    $userToken = '{UserHome}\'
    $harnessToken = '{HarnessRoot}\'
    if ($relative.StartsWith($userToken, [StringComparison]::OrdinalIgnoreCase)) {
        return Join-Path $userRoot $relative.Substring($userToken.Length)
    }
    if ($relative.StartsWith($harnessToken, [StringComparison]::OrdinalIgnoreCase)) {
        return Join-Path $manifestRoot $relative.Substring($harnessToken.Length)
    }
    throw "Unknown manifest path token: $tokenPath"
}

function Test-Contained([string]$root, [string]$candidate) {
    $r = ([IO.Path]::GetFullPath($root)).TrimEnd('\') + '\'
    $c = ([IO.Path]::GetFullPath($candidate)).TrimEnd('\')
    return $c.StartsWith($r, [StringComparison]::OrdinalIgnoreCase) -or $c.Equals($r.TrimEnd('\'), [StringComparison]::OrdinalIgnoreCase)
}

function Get-Lrf([string]$path) {
    $ids = [System.Collections.Generic.List[string]]::new()
    $errors = [System.Collections.Generic.List[string]]::new()
    $records = @{}
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { return [pscustomobject]@{ Ids = @(); Records = @{}; Errors = @('missing') } }
    $header = $false
    $effects = @('RO_LOCAL','RO_PUBLIC','LW_SCOPE','EXT_WRITE','RELEASE','LIVE_WRITE','VCS_WRITE','DESTRUCT','CRED','CHARGE','PROD_DEP')
    foreach ($line in (Get-Content -LiteralPath $path)) {
        if ([string]::IsNullOrWhiteSpace($line)) { continue }
        if ($line -match '^#SIGMA LRF/1$' -or $line -match '^F=.+$' -or $line -match '^C:.+$' -or $line -match '^L invokes .+$') { continue }
        if ($line.StartsWith('@lrf=1|')) { if ($header) { $errors.Add('duplicate-header') }; $header = $true; continue }
        $parts = $line -split '\|', 6
        $guardValid = $parts.Count -ge 3 -and ($parts[2] -eq '*' -or $parts[2] -match '^[a-z][a-z0-9_-]*:[^|&]+(?:&[a-z][a-z0-9_-]*:[^|&]+)*$')
        $recordValid = $parts.Count -eq 6 -and $parts[0] -in @('R','L') -and $parts[1] -match '^[^|\s]+$' -and $guardValid -and $parts[3] -in @('MUST','MUST_NOT','MAY') -and $parts[4] -in $effects -and $parts[5] -match '^.+$'
        if (-not $recordValid) { $errors.Add("malformed:$line"); continue }
        if ($ids.Contains($parts[1])) {
            $errors.Add("duplicate-id:$($parts[1])")
        } else {
            $ids.Add($parts[1])
            $records[$parts[1]] = $line
        }
    }
    if (-not $header) { $errors.Add('missing-header') }
    [pscustomobject]@{ Ids = @($ids); Records = $records; Errors = @($errors) }
}

function Compare-Ids($actual, $expected, $allowed) {
    $a = @($actual | Sort-Object -Unique); $e = @($expected | Sort-Object -Unique); $allow = @($allowed | Sort-Object -Unique)
    $missing = @($e | Where-Object { $_ -notin $a }); $extra = @($a | Where-Object { $_ -notin $e })
    $diff = @($missing + $extra | Sort-Object -Unique)
    [pscustomobject]@{ Missing = $missing; Extra = $extra; Difference = $diff; Allowed = (@($diff | Where-Object { $_ -notin $allow }).Count -eq 0) }
}

function Compare-LrfRecords($actualRecords, $canonicalRecords) {
    $allIds = @($actualRecords.Keys + $canonicalRecords.Keys | Sort-Object -Unique)
    return @($allIds | Where-Object {
        -not $actualRecords.ContainsKey($_) -or
        -not $canonicalRecords.ContainsKey($_) -or
        $actualRecords[$_] -cne $canonicalRecords[$_]
    })
}

function Resolve-ExpectedText([string]$value, [string]$userRoot, [string]$manifestRoot) {
    $userForward = ([IO.Path]::GetFullPath($userRoot)).TrimEnd('\').Replace('\', '/')
    $harnessForward = ([IO.Path]::GetFullPath($manifestRoot)).TrimEnd('\').Replace('\', '/')
    return $value.Replace('{UserHomeForward}', $userForward).Replace('{HarnessRootForward}', $harnessForward).Replace('{UserHome}', $userRoot).Replace('{HarnessRoot}', $manifestRoot)
}

function Invoke-Doctor {
    $manifestFull = [IO.Path]::GetFullPath($ManifestPath)
    if (-not (Test-Path -LiteralPath $manifestFull -PathType Leaf)) { throw "Manifest missing: $manifestFull" }
    $manifest = Get-Content -Raw -LiteralPath $manifestFull | ConvertFrom-Json
    if ($manifest.schemaVersion -ne 1) { throw 'Unsupported manifest schemaVersion' }
    if ($null -eq $manifest.entries -or @($manifest.entries).Count -eq 0) { throw 'Manifest entries are required' }
    if ([string]$manifest.canonical -notmatch '^\{(UserHome|HarnessRoot)\}\\.+$' -or [string]$manifest.canonical -match '(^|[\\/])\.\.([\\/]|$)') { throw 'Invalid canonical path token' }
    if ([string]$manifest.canonicalSha256 -notmatch '^[A-Fa-f0-9]{64}$') { throw 'Invalid canonical SHA256' }
    $manifestRoot = Split-Path (Split-Path $manifestFull -Parent) -Parent
    $canonicalPath = Resolve-ManifestPath ([string]$manifest.canonical) $manifestRoot $UserHome
    if (-not (Test-Contained $manifestRoot $canonicalPath) -and -not (Test-Contained $UserHome $canonicalPath)) { throw 'Canonical path escapes roots' }
    if (-not (Test-Path -LiteralPath $canonicalPath -PathType Leaf)) { throw 'Canonical LRF is missing' }
    if ((Get-FileHash -LiteralPath $canonicalPath -Algorithm SHA256).Hash -ne ([string]$manifest.canonicalSha256).ToUpperInvariant()) { throw 'Canonical LRF hash differs from manifest' }
    $canonicalLrf = Get-Lrf $canonicalPath
    if ($canonicalLrf.Errors.Count -gt 0) { throw "Canonical LRF invalid: $($canonicalLrf.Errors -join ';')" }
    $names = @{}
    $receipts = [System.Collections.Generic.List[object]]::new()
    foreach ($entry in @($manifest.entries)) {
        if ([string]::IsNullOrWhiteSpace([string]$entry.name) -or $names.ContainsKey([string]$entry.name)) { throw 'Manifest entry names must be unique and nonempty' }
        $names[[string]$entry.name] = $true
        if ([string]$entry.kind -notin @('HarnessProject','Codex','Gemini','OpenCode','Agents','Claude')) { throw "Unknown entry kind: $($entry.kind)" }
        if ($entry.required -isnot [bool]) { throw "Entry required must be boolean: $($entry.name)" }
        if ($null -eq $entry.chain -or @($entry.chain).Count -eq 0) { throw "Entry chain is required: $($entry.name)" }
        if ([string]$entry.entry -notmatch '^\{(UserHome|HarnessRoot)\}\\.+$' -or [string]$entry.entry -match '(^|[\\/])\.\.([\\/]|$)') { throw "Invalid entry path token: $($entry.name)" }
        $entryPath = Resolve-ManifestPath ([string]$entry.entry) $manifestRoot $UserHome
        $entryRoot = Split-Path $entryPath -Parent
        $chainPaths = @($entry.chain | ForEach-Object { [string]$_.path })
        if ([string]$entry.entry -notin $chainPaths) { $receipts.Add((New-Receipt 'EntryMismatch' $entry.name 'entry-not-in-chain' $entry.required)); continue }
        $missing = $false; $entryMismatch = $false; $unexpected = $false; $allowedOnly = $false; $idErrors = @(); $idDiffs = @(); $details = @()
        foreach ($item in @($entry.chain)) {
            if ([string]$item.path -notmatch '^\{(UserHome|HarnessRoot)\}\\.+$' -or [string]$item.path -match '(^|[\\/])\.\.([\\/]|$)') { throw "Invalid chain path token: $($entry.name)" }
            if ([string]$item.sha256 -notmatch '^[A-Fa-f0-9]{64}$') { throw "Invalid SHA256 in chain: $($entry.name)" }
            $actualPath = Resolve-ManifestPath ([string]$item.path) $manifestRoot $UserHome
            if (-not (Test-Contained $manifestRoot $actualPath) -and -not (Test-Contained $UserHome $actualPath)) { throw "Manifest path escapes roots: $($item.path)" }
            if (-not (Test-Path -LiteralPath $actualPath -PathType Leaf)) { $missing = $true; continue }
            if ([string]$item.kind -notin @('entry','lrf','text')) { throw "Unknown chain kind: $($item.kind)" }
            $hash = (Get-FileHash -LiteralPath $actualPath -Algorithm SHA256).Hash.ToUpperInvariant()
            if ($hash -ne ([string]$item.sha256).ToUpperInvariant()) { $unexpected = $true; $details += "hash:$($item.path)" }
            if ($item.PSObject.Properties.Name -contains 'contains' -and -not [string]::IsNullOrWhiteSpace([string]$item.contains)) {
                $expectedText = Resolve-ExpectedText ([string]$item.contains) $UserHome $manifestRoot
                if (-not (Get-Content -Raw -LiteralPath $actualPath).Contains($expectedText)) { $entryMismatch = $true; $details += "reference:$($item.path)" }
            }
            if ($item.kind -eq 'lrf') {
                $parsed = Get-Lrf $actualPath; $idErrors += @($parsed.Errors)
                $snapshotCmp = Compare-Ids $parsed.Ids @($item.ids) @()
                $recordDiffs = @(Compare-LrfRecords $parsed.Records $canonicalLrf.Records)
                $idDiffs += @($recordDiffs)
                $unallowedRecordDiffs = @($recordDiffs | Where-Object { $_ -notin @($item.allowedDifferenceIds) })
                if ($parsed.Errors.Count -gt 0 -or -not $snapshotCmp.Allowed -or $unallowedRecordDiffs.Count -gt 0) {
                    $unexpected = $true
                } elseif ($recordDiffs.Count -gt 0) {
                    $allowedOnly = $true
                }
            }
        }
        if ($missing) { $receipts.Add((New-Receipt $(if($entry.required){'Missing'}else{'NotConfigured'}) $entry.name 'chain-path-missing' $entry.required)); continue }
        if ($entryMismatch) { $receipts.Add((New-Receipt 'EntryMismatch' $entry.name ($details -join ';') $entry.required)); continue }
        if ($unexpected) { $receipts.Add((New-Receipt 'UnexpectedDifference' $entry.name (($details + $idErrors + $idDiffs) -join ';') $entry.required)); continue }
        if ($allowedOnly) { $receipts.Add((New-Receipt 'AllowedDifference' $entry.name (($idDiffs | Sort-Object -Unique) -join ',') $entry.required)); continue }
        $receipts.Add((New-Receipt 'Exact' $entry.name 'all-chain-hashes-and-ids-match' $entry.required))
    }
    $failed = @($receipts | Where-Object { $_.Classification -in @('Missing','EntryMismatch','UnexpectedDifference') -and $_.Required }).Count -gt 0
    [pscustomobject]@{ SchemaVersion = 1; Manifest = $manifestFull; Passed = -not $failed; Results = @($receipts); GeneratedAt = (Get-Date).ToString('o') }
}

try {
    $receipt = Invoke-Doctor
    if ($Json) { $receipt | ConvertTo-Json -Depth 12 -Compress } else {
        foreach ($result in $receipt.Results) { Write-Host ("[{0}] {1}: {2}" -f $result.Classification, $result.Name, $result.Detail) }
        if ($receipt.Passed) { Write-Host 'Rule distribution doctor passed.' } else { Write-Error 'Rule distribution doctor failed.' }
    }
    if (-not $receipt.Passed) { exit 1 }
} catch {
    $failure = [pscustomobject]@{ SchemaVersion = 1; Passed = $false; Error = $_.Exception.Message; Manifest = $ManifestPath }
    if ($Json) { $failure | ConvertTo-Json -Depth 8 -Compress } else { Write-Error $_.Exception.Message }
    exit 2
}
