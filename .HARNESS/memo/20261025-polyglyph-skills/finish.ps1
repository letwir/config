$ErrorActionPreference = 'Stop'
$root = 'C:/Users/letwir/.harness/memo/20261025-polyglyph-skills'
$taskId = '20261025-polyglyph-skills'
$utf8 = [Text.UTF8Encoding]::new($false)
$OutputEncoding = $utf8
[Console]::OutputEncoding = $utf8
function Save-Json([string]$p, $o) { [IO.File]::WriteAllText($p, ($o | ConvertTo-Json -Depth 100), $utf8) }
$diary = "$root/final-diary.md"
$text = [IO.File]::ReadAllText($diary).Replace('### 2026-10-25 19:15:00', ('### ' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss')))
[IO.File]::WriteAllText($diary, $text, $utf8)
$runs = @(
    @{id='278e9ee7-b52d-4b46-a629-e903abbb7db9'; role='worker'; result='FAIL'; reason='Broad inventory failed shell quoting; no target edits.'},
    @{id='ede63fc0-23f9-4805-9503-b2d10fd28b7c'; role='worker'; result='FAIL'; reason='Partial AGY conversion; no validation; unsupported language.'},
    @{id='33a319b0-fc75-457b-b4da-f8bff76b5c1c'; role='worker'; result='FAIL'; reason='AGY conversion retained unsupported language; shell validation failed.'},
    @{id='4f306ade-8cbb-4e41-9ef9-15d12ae81d6e'; role='worker'; result='REJECTED_IN_REVIEW'; reason='Parent byte checks passed but independent review rejected ordinary-English representation.'},
    @{id='d2309fec-ce7c-4297-be5a-4d9198c276e4'; role='worker'; result='FAIL'; reason='Manifest hydration failed shell interpolation; no target edits.'},
    @{id='2e0c57e2-d838-4a9e-9b9b-e441b7013515'; role='worker'; result='CORRECTION_REQUIRED'; reason='26 label edits completed but edit tool normalized immutable final CRLF; later repaired.'},
    @{id='69802ce1-34c9-4eff-9612-2e2cfdf63574'; role='worker'; result='CORRECTION_REQUIRED'; reason='39 skill structure and codec passed; parent identified missing local-reference details.'},
    @{id='38ee08ab-33fb-432d-a983-32ad5e6f621c'; role='worker'; result='PASS'; reason='Completed 39-file reuse/sync with nine approved reference unions; parent and independent checks passed.'},
    @{id='825475e4-4833-4fae-aca6-1f50017dac70'; role='verifier'; result='REJECT'; reason='Identified incomplete AGY representation; profile lacked shell for independent command execution.'},
    @{id='55ec2d17-981d-48ae-acc0-cd51394f23c9'; role='worker'; result='PASS'; reason='Genuine Japanese/English/logical AGY conversion; parent opaque-byte tests and independent review passed.'},
    @{id='635e851e-f01a-40ea-9538-8fc7cef8bfd3'; role='verifier'; result='PASS_WITH_LIMITS'; reason='Independent validator and two negative fixtures exited0; sampled semantic review PASS. Inferred native attestation missing structured report is recorded separately.'}
)
$evaluations = @()
foreach ($run in $runs) {
    $path = "A:/TMP/pi-subagents-user-letwir/async-subagent-runs/$($run.id)/status.json"
    $status = [IO.File]::ReadAllText($path) | ConvertFrom-Json
    $step = $status.steps[0]
    $provider = $step.resolvedModel.provider
    $model = $step.resolvedModel.model
    if (-not $provider -or -not $model) { throw "Missing correlated runtime model: $($run.id)" }
    if ($provider -ne 'openai-codex' -or $model -ne 'gpt-6-luna') { throw "Unexpected correlated model: $($run.id)" }
    $evaluations += [pscustomobject]@{
        task_id="$taskId/$($run.id)"; difficulty='medium'; role=$run.role
        requested_provider='openai-codex'; effective_provider=$provider
        requested_model='gpt-6-luna'; effective_model=$model
        requested_reasoning='medium'; effective_reasoning=$step.thinking
        acceptance_result=$run.result; outcome=$run.result
        quality=$(if ($run.result -like 'PASS*') {'high-with-stated-limits'} else {'requires-correction'})
        reason=$run.reason; runtime="$($step.durationMs)ms"; process_exit=$step.exitCode
        native_acceptance=$step.acceptance.status
        native_acceptance_note='Inferred native attestation status is not rewritten or represented as a passed explicit gate.'
        routing_change='none'; date=(Get-Date -Format yyyy-MM-dd); confidence=0.95
    }
}
Save-Json "$root/evaluation.json" ([ordered]@{task_id=$taskId;comparison_key='polyglyph-skills-local-migration';axes=@{implementation=1;structural_verification=1;independent_checks=1;operational_efficiency=0};attribution=@{prompt=0;agent=100};ratio_scope='qualitative attribution for observed task friction, not calibrated model quality';records=$evaluations;database_eval_receipt='NOT_SUBMITTED: task-memory authorization covers report ingest, not unrelated eval writes'})
$sharedPath = 'C:/Users/letwir/.harness/evaluation/subagents_model_eval.json'
$sharedBytes = [IO.File]::ReadAllBytes($sharedPath)
$sharedHash = (Get-FileHash -LiteralPath $sharedPath).Hash
$shared = [Text.Encoding]::UTF8.GetString($sharedBytes) | ConvertFrom-Json
$backup = "$root/shared-evaluation.before.json"
if (-not (Test-Path -LiteralPath $backup)) { [IO.File]::WriteAllBytes($backup, $sharedBytes) }
$existing = @($shared.evaluations.task_id)
$new = @($evaluations | Where-Object { $existing -notcontains $_.task_id })
if ($new.Count -gt 0) {
    $shared.evaluations = @($shared.evaluations) + $new
    $shared.updated_at = Get-Date -Format yyyy-MM-dd
    if ((Get-FileHash -LiteralPath $sharedPath).Hash -ne $sharedHash) { throw 'Shared evaluation changed concurrently; refusing overwrite' }
    Save-Json $sharedPath $shared
}
$check = [IO.File]::ReadAllText($sharedPath) | ConvertFrom-Json
foreach ($e in $evaluations) { if (@($check.evaluations | Where-Object task_id -eq $e.task_id).Count -ne 1) { throw "Missing/duplicate evaluation: $($e.task_id)" } }
"SHARED_EVALUATION: $($evaluations.Count) correlated records verified"
$ledgerPath = "$root/memory-receipts.json"
$ledger = @(if (Test-Path -LiteralPath $ledgerPath) { [IO.File]::ReadAllText($ledgerPath) | ConvertFrom-Json })
$alreadySubmitted = @(
    [IO.Path]::GetFullPath("$root/worker-report.md"),
    [IO.Path]::GetFullPath("$root/retry-agy/worker-report.md")
)
$reports = @(Get-ChildItem -LiteralPath $root -Recurse -File -Filter '*.md' | Where-Object {
    $_.Name -match '^worker-(report|pg-final).*\.md$|^independent-review.*\.md$|^walkthrough\.md$|^final-diary\.md$'
})
$bin = 'C:/Users/letwir/.harness/skills/llm-memory/llm-mem.exe'
foreach ($file in $reports) {
    $full = [IO.Path]::GetFullPath($file.FullName)
    if ($alreadySubmitted -contains $full -or @($ledger | Where-Object path -eq $full).Count -gt 0) { continue }
    $rel = [IO.Path]::GetRelativePath($root, $full)
    $cat = if ($file.Name -eq 'walkthrough.md') {'walkthrough'} else {'diary'}
    $record = [pscustomobject]@{path=$full;category=$cat;sha256=(Get-FileHash $full).Hash;status='ATTEMPT_STARTED';exit_code=$null;response_ids=@();date=(Get-Date -Format o)}
    $ledger = @($ledger) + @($record)
    Save-Json $ledgerPath @($ledger)
    $raw = @(& $bin ingest -file $full -title "$taskId / $rel" -cat $cat 2>&1)
    $code = $LASTEXITCODE
    $record.exit_code = $code
    $record.status = if ($code -eq 0) {'SUCCESS'} else {'FAILED'}
    $record.response_ids = @([regex]::Matches(($raw -join "`n"), '[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}') | ForEach-Object Value | Select-Object -Unique)
    Save-Json $ledgerPath @($ledger)
    "INGEST $rel : exit=$code"
}
$analysis = @(& $bin analyze -file "$root/final-diary.md" -suggest -json 2>&1)
$analysisCode = $LASTEXITCODE
if ($analysisCode -eq 0) {
    $analysisText = $analysis -join "`n"
    $null = $analysisText | ConvertFrom-Json
    [IO.File]::WriteAllText("$root/defect-analysis.json", $analysisText, $utf8)
}
"DEFECT_ANALYZE_EXIT=$analysisCode"
$failed = @($ledger | Where-Object status -ne SUCCESS)
"MEMORY_REPORTS_SUCCESS=$(@($ledger | Where-Object status -eq SUCCESS).Count); FAILED=$($failed.Count); PRIOR_TWO_WORKER_REPORTS_ALREADY_INGESTED=2"
if ($failed.Count -gt 0 -or $analysisCode -ne 0) { exit 1 }
