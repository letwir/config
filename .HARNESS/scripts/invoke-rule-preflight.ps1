[CmdletBinding(DefaultParameterSetName = 'Resolve')]
param(
    [Parameter(ParameterSetName = 'Resolve')]
    [string]$Root,
    [Parameter(Mandatory, ParameterSetName = 'Resolve')]
    [ValidatePattern('^[a-z][a-z0-9.-]*$')]
    [string]$Task,
    [Parameter(ParameterSetName = 'Resolve')]
    [ValidateScript({ @($_) | ForEach-Object { if ($_ -notmatch '^[a-z][a-z0-9.-]*$') { throw "tag invalid: $_" } }; $true })]
    [string[]]$Tag = @(),
    [Parameter(ParameterSetName = 'Resolve')]
    [string]$OutFile,
    [Parameter(ParameterSetName = 'Resolve')]
    [switch]$Detailed,
    [Parameter(Mandatory, ParameterSetName = 'Verify')]
    [string]$VerifyReceipt,
    [int]$MaxDepth = 16,
    [int]$MaxModules = 32,
    [int64]$MaxBytes = 1048576
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$harnessRoot = [IO.Path]::GetFullPath((Split-Path $PSScriptRoot -Parent)).TrimEnd('\')
if (-not $Root) { $Root = Join-Path $harnessRoot 'rules\BOOTSTRAP.lrf' }
$allowedModalities = @('MUST', 'MUST_NOT', 'MAY')
$allowedEffects = @('RO_LOCAL', 'RO_PUBLIC', 'LW_SCOPE', 'EXT_WRITE', 'RELEASE', 'LIVE_WRITE', 'VCS_WRITE', 'DESTRUCT', 'CRED', 'CHARGE', 'PROD_DEP')
$allowedGuards = @('task', 'tag', 'change', 'phase', 'fact', 'file', 'event', 'case')
$lint = Join-Path $harnessRoot 'skills\harness-lint\harness-lint.exe'

function Fail([string]$message) { throw "Rule preflight failed: $message" }

function Assert-SafeLeaf([string]$full) {
    $full = [IO.Path]::GetFullPath($full)
    if (-not ($full.Equals($harnessRoot, [StringComparison]::OrdinalIgnoreCase) -or $full.StartsWith($harnessRoot + '\', [StringComparison]::OrdinalIgnoreCase))) {
        Fail "path escapes harness root: $full"
    }
    if (-not (Test-Path -LiteralPath $full -PathType Leaf)) { Fail "missing reference: $full" }
    $leaf = Get-Item -LiteralPath $full -Force
    if (($leaf.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { Fail "reparse-point leaf is forbidden: $full" }
    $cursor = Split-Path $full -Parent
    while ($cursor.Length -ge $harnessRoot.Length) {
        $item = Get-Item -LiteralPath $cursor -Force
        if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { Fail "reparse-point path is forbidden: $cursor" }
        if ($cursor.Equals($harnessRoot, [StringComparison]::OrdinalIgnoreCase)) { break }
        $cursor = Split-Path $cursor -Parent
    }
    return $full
}

function Resolve-SafePath([string]$candidate, [string]$baseDirectory) {
    if ([IO.Path]::IsPathRooted($candidate)) { Fail "rooted reference is forbidden: $candidate" }
    $full = [IO.Path]::GetFullPath((Join-Path $baseDirectory $candidate))
    $null = Assert-SafeLeaf $full
    return $full
}

function Test-Guard([string]$guard) {
    if ($guard -eq '*') { return $true }
    $matched = $true
    foreach ($atom in $guard -split '&') {
        if ($atom -notmatch '^([a-z]+):([A-Za-z0-9._-]+)$') { Fail "invalid guard: $guard" }
        $key = $Matches[1]; $value = $Matches[2]
        if ($allowedGuards -notcontains $key) { Fail "unknown guard key: $key" }
        if ($key -eq 'task') { $matched = $matched -and ($value -eq $Task) }
        elseif ($key -eq 'tag') { $matched = $matched -and ($Tag -contains $value) }
        else { $matched = $false }
    }
    return $matched
}

function Read-Lrf([string]$path) {
    $raw = [IO.File]::ReadAllText($path, [Text.UTF8Encoding]::new($false, $true))
    $bytes = [Text.Encoding]::UTF8.GetBytes($raw)
    $script:totalBytes += $bytes.Length
    if ($script:totalBytes -gt $MaxBytes) { Fail "byte limit exceeded ($MaxBytes)" }
    $lines = $raw -split "`r?`n"
    $headerIndex = -1
    for ($h = 0; $h -lt $lines.Count; $h++) {
        if ($lines[$h] -match '^@lrf=1\|aud=[^|]+\|scope=[^|]+$') { $headerIndex = $h; break }
    }
    if ($headerIndex -lt 0) { Fail "invalid header: $path" }
    $records = @()
    for ($i = $headerIndex + 1; $i -lt $lines.Count; $i++) {
        $line = $lines[$i].Trim()
        if (-not $line -or $line.StartsWith('#')) { continue }
        $parts = $line -split '\|', 6
        if ($parts.Count -ne 6) { Fail "invalid field count at ${path}:$($i + 1)" }
        $kind, $id, $guard, $modality, $effect, $payload = $parts
        if ($kind -notin @('R', 'L')) { Fail "unknown record kind at ${path}:$($i + 1)" }
        if ($id -notmatch '^[a-z][a-z0-9.-]*$') { Fail "invalid id at ${path}:$($i + 1)" }
        if ($script:ids.ContainsKey($id)) { Fail "duplicate id '$id'" }
        $script:ids[$id] = $true
        if ($allowedModalities -notcontains $modality) { Fail "invalid modality at ${path}:$($i + 1)" }
        if ($allowedEffects -notcontains $effect) { Fail "invalid effect at ${path}:$($i + 1)" }
        if (-not $payload) { Fail "empty payload at ${path}:$($i + 1)" }
        $isMatch = Test-Guard $guard
        $records += [pscustomobject]@{ kind=$kind; id=$id; guard=$guard; payload=$payload; line=$i + 1; matched=$isMatch }
    }
    $sha = [Security.Cryptography.SHA256]::Create()
    try { $hash = ([BitConverter]::ToString($sha.ComputeHash($bytes))).Replace('-', '') }
    finally { $sha.Dispose() }
    return [pscustomobject]@{ records=@($records); sha256=$hash }
}

function Visit-Lrf([string]$path, [int]$depth) {
    if ($depth -gt $MaxDepth) { Fail "depth limit exceeded ($MaxDepth)" }
    $key = [IO.Path]::GetFullPath($path).ToUpperInvariant()
    if ($script:stack.ContainsKey($key)) { Fail "reference cycle: $path" }
    if ($script:visited.ContainsKey($key)) { return }
    if ($script:visited.Count -ge $MaxModules) { Fail "module limit exceeded ($MaxModules)" }
    $script:stack[$key] = $true
    $parsed = Read-Lrf $path
    if (Test-Path -LiteralPath $lint -PathType Leaf) {
        & $lint -path $path -level 0 -strict *> $null
        if ($LASTEXITCODE -ne 0) { Fail "harness-lint rejected: $path" }
    }
    $script:files += [pscustomobject]@{ path=$path; sha256=$parsed.sha256 }
    $script:visited[$key] = $true
    foreach ($record in $parsed.records) {
        if (-not $record.matched) { continue }
        $script:selected += [pscustomobject]@{ path=$path; line=$record.line; id=$record.id; kind=$record.kind }
        if ($record.kind -eq 'L') {
            if ($record.payload -notmatch '\[[^\]]+\]\(([^)]+\.lrf)\)') { Fail "invalid L reference at ${path}:$($record.line)" }
            $target = Resolve-SafePath $Matches[1] (Split-Path $path -Parent)
            Visit-Lrf $target ($depth + 1)
        }
    }
    $script:stack.Remove($key)
}

if ($PSCmdlet.ParameterSetName -eq 'Verify') {
    $receipt = Get-Content -Raw -LiteralPath $VerifyReceipt | ConvertFrom-Json
    if ($receipt.schema -ne 'lrf-preflight/v1') { Fail 'unknown receipt schema' }
    if ($receipt.parser -ne 'invoke-rule-preflight.ps1/v1') { Fail 'unknown receipt parser' }
    if ([string]$receipt.task -notmatch '^[a-z][a-z0-9.-]*$') { Fail 'invalid receipt task' }
    if (@($receipt.files).Count -lt 1 -or @($receipt.selected).Count -lt 1) { Fail 'receipt graph is empty' }
    $receiptRoot = [IO.Path]::GetFullPath([string]$receipt.root)
    $receiptPaths = @{}
    foreach ($file in $receipt.files) {
        $path = Assert-SafeLeaf ([string]$file.path)
        if ((Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash -ne [string]$file.sha256) { Fail "stale receipt: $path" }
        $receiptPaths[$path.ToUpperInvariant()] = $true
    }
    if (-not $receiptPaths.ContainsKey($receiptRoot.ToUpperInvariant())) { Fail 'receipt root is not in graph' }
    foreach ($selection in $receipt.selected) {
        $selectedPath = [IO.Path]::GetFullPath([string]$selection.path)
        if (-not $receiptPaths.ContainsKey($selectedPath.ToUpperInvariant())) { Fail 'selected record path is not in graph' }
        if ([string]$selection.id -notmatch '^[a-z][a-z0-9.-]*$' -or [string]$selection.kind -notin @('R', 'L') -or [int]$selection.line -lt 1) { Fail 'invalid selected record' }
    }
    [pscustomobject]@{ schema='lrf-preflight-verify/v1'; status='PASS'; file_count=@($receipt.files).Count } | ConvertTo-Json -Compress
    exit 0
}

if ($MaxDepth -lt 1 -or $MaxModules -lt 1 -or $MaxBytes -lt 1) { Fail 'limits must be positive' }
if (-not $Root) { $Root = Join-Path $harnessRoot 'rules\BOOTSTRAP.lrf' }
$Root = [IO.Path]::GetFullPath($Root)
$Root = Assert-SafeLeaf $Root
$script:ids = @{}; $script:visited = @{}; $script:stack = @{}; $script:files = @(); $script:selected = @(); $script:totalBytes = 0
Visit-Lrf $Root 0
$receipt = [pscustomobject]@{
    schema='lrf-preflight/v1'; status='PASS'; parser='invoke-rule-preflight.ps1/v1'; task=$Task; tags=@($Tag)
    root=$Root; files=@($script:files); selected=@($script:selected); total_bytes=$script:totalBytes
}
if ($OutFile) {
    $parent = Split-Path ([IO.Path]::GetFullPath($OutFile)) -Parent
    if (-not (Test-Path -LiteralPath $parent -PathType Container)) { Fail "receipt parent missing: $parent" }
    [IO.File]::WriteAllText([IO.Path]::GetFullPath($OutFile), ($receipt | ConvertTo-Json -Depth 6), [Text.UTF8Encoding]::new($false))
}
if ($Detailed) { $receipt | ConvertTo-Json -Depth 6 -Compress }
else {
    [pscustomobject]@{
        schema='lrf-preflight-summary/v1'; status='PASS'; task=$Task; tags=@($Tag)
        file_count=@($script:files).Count; selected_count=@($script:selected).Count; total_bytes=$script:totalBytes
        receipt_path=if ($OutFile) { [IO.Path]::GetFullPath($OutFile) } else { $null }
    } | ConvertTo-Json -Compress
}
