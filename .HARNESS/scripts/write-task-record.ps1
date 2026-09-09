[CmdletBinding()]
param(
    [Parameter(Mandatory)][ValidateNotNullOrEmpty()][string]$TaskId,
    [string]$RootPath = (Split-Path -Parent $PSScriptRoot),
    [ValidateSet('SUCCESS','N_A','FAILED')][string]$IngestStatus = 'N_A',
    [ValidateSet('SUCCESS','N_A','FAILED')][string]$EvaluationStatus = 'N_A',
    [string]$IngestReasonCode = 'pending_approval',
    [string]$EvaluationReasonCode = 'pending_approval',
    [string]$IngestReceiptId,
    [string]$EvaluationReceiptId,
    [ValidateSet('SUCCESS','N_A','FAILED')][string]$MirrorStatus = 'N_A',
    [string]$MirrorReasonCode = 'not_requested',
    [string]$MirrorReceiptId,
    [ValidateRange(1,60000)][int]$LockTimeoutMs = 5000
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$script:SchemaVersion = 1
$script:Reasons = @('pending_approval','not_requested','boundary_unavailable','completed','process_error','missing_executable','invalid_task_id','root_missing','reparse_point','path_escape','artifact_missing','artifact_drift','evaluation_task_id_mismatch','evaluation_invalid','lock_timeout','replace_failed','io_error','unknown_failure')
function New-Receipt([string]$status, [string]$reason, [string]$receiptId) {
    $allowed = switch ($status) { 'SUCCESS' { @('completed') } 'N_A' { @('pending_approval','not_requested','boundary_unavailable') } 'FAILED' { @('process_error','missing_executable') } }
    if ($reason -notin $allowed) { throw 'io_error' }
    if ($status -eq 'SUCCESS' -and ($receiptId -notmatch '^[A-Za-z0-9._-]{1,64}$')) { throw 'io_error' }
    if ($status -ne 'SUCCESS' -and -not [string]::IsNullOrEmpty($receiptId)) { throw 'io_error' }
    $value = [ordered]@{ status = $status; reason_code = $reason }
    if ($status -eq 'SUCCESS') { $value.receipt_id = $receiptId }
    $value
}
function Get-SafeFullPath([string]$path) { if ([string]::IsNullOrWhiteSpace($path)) { throw 'path_escape' }; [IO.Path]::GetFullPath($path) }
function Assert-NoReparse([string]$path) { $item = Get-Item -LiteralPath $path -Force -ErrorAction Stop; if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { throw 'reparse_point' } }
function Assert-Contained([string]$root, [string]$candidate) {
    $r = (Get-SafeFullPath $root).TrimEnd('\') + '\'; $c = Get-SafeFullPath $candidate
    if (-not $c.StartsWith($r, [StringComparison]::OrdinalIgnoreCase)) { throw 'path_escape' }
    $cursor = $r.TrimEnd('\'); foreach ($part in (($c.Substring($r.Length)) -split '\\')) { if ([string]::IsNullOrEmpty($part)) { continue }; $cursor = Join-Path $cursor $part; if (Test-Path -LiteralPath $cursor) { Assert-NoReparse $cursor } }
}
function Read-StableArtifact([string]$path, [string]$relativePath) {
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw 'artifact_missing' }; Assert-NoReparse $path
    $a = [IO.File]::ReadAllBytes($path); $b = [IO.File]::ReadAllBytes($path); if (-not [Linq.Enumerable]::SequenceEqual($a, $b)) { throw 'artifact_drift' }
    $hash = [Security.Cryptography.SHA256]::Create(); try { $digest = (($hash.ComputeHash($a) | ForEach-Object ToString x2) -join '') } finally { $hash.Dispose() }
    [ordered]@{ path = $relativePath; sha256 = $digest; bytes = $a.Length; bytes_value = $a }
}
function Write-AtomicBytes([string]$target, [byte[]]$bytes) {
    $dir = Split-Path -Parent $target; $tmp = Join-Path $dir ('.record.' + [guid]::NewGuid().ToString('N') + '.tmp'); $backup = Join-Path $dir ('.record.' + [guid]::NewGuid().ToString('N') + '.bak')
    try {
        $fs = [IO.File]::Open($tmp, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None); try { $fs.Write($bytes, 0, $bytes.Length); $fs.Flush($true) } finally { $fs.Dispose() }
        if (Test-Path -LiteralPath $target -PathType Leaf) { Assert-NoReparse $target; if ($env:H07_INJECT_REPLACE_FAILURE -eq '1') { throw 'replace_failed' }; try { [IO.File]::Replace($tmp, $target, $backup) } catch { throw 'replace_failed' } } else { [IO.File]::Move($tmp, $target) }
    } finally { if (Test-Path -LiteralPath $tmp) { Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue }; if (Test-Path -LiteralPath $backup) { Remove-Item -LiteralPath $backup -Force -ErrorAction SilentlyContinue } }
}
function Get-FailureReason([string]$message) { if ($message -in $script:Reasons) { return $message }; if ($message -match 'evaluation') { return 'evaluation_invalid' }; if ($message -match 'lock') { return 'lock_timeout' }; return 'unknown_failure' }
$lock = $null
try {
    if ($TaskId -notmatch '^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$' -or $TaskId -match '(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(?:\.|$)') { throw 'invalid_task_id' }
    $root = Get-SafeFullPath $RootPath; if (-not (Test-Path -LiteralPath $root -PathType Container)) { throw 'root_missing' }; Assert-NoReparse $root
    $memo = Join-Path $root 'memo'; if (-not (Test-Path -LiteralPath $memo -PathType Container)) { New-Item -ItemType Directory -Path $memo | Out-Null }; Assert-NoReparse $memo
    $taskMatches = @(Get-ChildItem -LiteralPath $memo -Directory -Force | Where-Object { $_.Name.Equals($TaskId, [StringComparison]::OrdinalIgnoreCase) })
    if ($taskMatches.Count -ne 1) { throw 'artifact_missing' }
    $taskItem = $taskMatches[0]; if ($taskItem.Name -cne $TaskId) { throw 'invalid_task_id' }
    $taskDir = $taskItem.FullName; Assert-NoReparse $taskDir
    Assert-Contained $root $taskDir; $recordPath = Join-Path $taskDir 'record.json'; $lockPath = Join-Path $taskDir '.record.lock'
    $deadline = [Diagnostics.Stopwatch]::GetTimestamp() + ([Diagnostics.Stopwatch]::Frequency * $LockTimeoutMs / 1000)
    while ($true) { try { $lock = [IO.File]::Open($lockPath, [IO.FileMode]::OpenOrCreate, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None); break } catch { if ([Diagnostics.Stopwatch]::GetTimestamp() -ge $deadline) { throw 'lock_timeout' }; Start-Sleep -Milliseconds 25 } }
    $names = [ordered]@{ knowledge = 'knowledge.md'; diary = 'diary.md'; walkthrough = 'walkthrough.md'; evaluation = 'evaluation.json' }; $artifacts = [ordered]@{}
    foreach ($name in $names.Keys) { $relative = 'memo/' + $TaskId + '/' + $names[$name]; $path = Join-Path $taskDir $names[$name]; Assert-Contained $root $path; $artifacts[$name] = Read-StableArtifact $path $relative }
    try { $evaluation = ([Text.Encoding]::UTF8.GetString([byte[]]$artifacts.evaluation.bytes_value) | ConvertFrom-Json -ErrorAction Stop) } catch { throw 'evaluation_invalid' }
    if ([string]$evaluation.task_id -cne $TaskId) { throw 'evaluation_task_id_mismatch' }
    foreach ($name in $names.Keys) { $again = Read-StableArtifact (Join-Path $taskDir $names[$name]) $artifacts[$name].path; if ($again.sha256 -cne $artifacts[$name].sha256 -or $again.bytes -ne $artifacts[$name].bytes) { throw 'artifact_drift' } }
    $record = [ordered]@{ schema_version = $script:SchemaVersion; task_id = $TaskId; generated_at = (Get-Date).ToUniversalTime().ToString('o'); artifact_root = 'memo/' + $TaskId; artifacts = [ordered]@{}; llm_memory_ingest = (New-Receipt $IngestStatus $IngestReasonCode $IngestReceiptId); llm_memory_eval = (New-Receipt $EvaluationStatus $EvaluationReasonCode $EvaluationReceiptId); external_diary_mirror = (New-Receipt $MirrorStatus $MirrorReasonCode $MirrorReceiptId); memory_report_complete = ($IngestStatus -eq 'SUCCESS' -and $EvaluationStatus -eq 'SUCCESS') }
    foreach ($name in $names.Keys) { $record.artifacts[$name] = [ordered]@{ path = $artifacts[$name].path; sha256 = $artifacts[$name].sha256; bytes = $artifacts[$name].bytes } }
    Write-AtomicBytes $recordPath ([Text.Encoding]::UTF8.GetBytes(($record | ConvertTo-Json -Depth 10 -Compress))); $lock.Dispose(); $lock = $null; $record | ConvertTo-Json -Depth 10 -Compress; exit 0
} catch { if ($lock) { $lock.Dispose() }; [ordered]@{ schema_version = $script:SchemaVersion; task_id = $TaskId; status = 'FAILED'; reason_code = (Get-FailureReason $_.Exception.Message); memory_report_complete = $false } | ConvertTo-Json -Depth 5 -Compress; exit 1 }
