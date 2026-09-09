[CmdletBinding()]
param(
  [string]$CasePath = '',
  [string]$RulesPath = '',
  [string]$OutputDirectory = '',
  [switch]$Json
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$script:ParserVersion = 'H08-1'
$script:Effects = @('RO_LOCAL','RO_PUBLIC','LW_SCOPE','EXT_WRITE','RELEASE','LIVE_WRITE','VCS_WRITE','DESTRUCT','CRED','CHARGE','PROD_DEP')
$script:GuardKeys = @('task','tag','change','phase','fact','file','event','case')
$script:StateFiles = @('decisions.md','method.md','knowledge.md','issues.md','memo.md','history.md','diary.md')
if ([string]::IsNullOrWhiteSpace($CasePath)) { $CasePath = Join-Path $PSScriptRoot '..\evaluation\oracle.json' }
if ([string]::IsNullOrWhiteSpace($RulesPath)) { $RulesPath = Join-Path $PSScriptRoot '..\rules\LLM_REF_RULE.md' }
if ([string]::IsNullOrWhiteSpace($OutputDirectory)) { $OutputDirectory = Join-Path $PSScriptRoot '..\evaluation\runs' }

function Fail([string]$m) { throw "policy-evaluation: $m" }
function Hash-File([string]$p) { if (!(Test-Path -LiteralPath $p -PathType Leaf)) { Fail "missing file: $p" }; (Get-FileHash -LiteralPath $p -Algorithm SHA256).Hash.ToUpperInvariant() }
function Read-Json([string]$p) { try { Get-Content -Raw -LiteralPath $p | ConvertFrom-Json -ErrorAction Stop } catch { Fail "invalid JSON ${p}: $($_.Exception.Message)" } }
function Is-Array($v) { return ($null -ne $v -and $v -is [System.Collections.IEnumerable] -and $v -isnot [string]) }
function Require-Text($v,[string]$name) { if ($v -isnot [string] -or [string]::IsNullOrWhiteSpace($v)) { Fail "$name must be a nonempty string" } }

function Parse-Lrf([string]$path) {
  if (!(Test-Path -LiteralPath $path -PathType Leaf)) { Fail "missing LRF: $path" }
  $header = $false; $ids = @{}; $records = @(); $lines = Get-Content -LiteralPath $path
  foreach ($line in $lines) {
    if ([string]::IsNullOrWhiteSpace($line)) { continue }
    if ($line -match '^@lrf=1\|aud=[^|]+\|scope=\S+$') { if ($header) { Fail "duplicate LRF header: $path" }; $header = $true; continue }
    if ($line -match '^(#SIGMA LRF/1|F=.+|C:.+|L invokes .+)$') { continue }
    $parts = $line -split '\|',6
    if ($parts.Count -ne 6) { Fail "LRF record is not six fields: ${path}: $line" }
    if ($parts[0] -notin @('R','L')) { Fail "unknown LRF record kind: ${path}: $line" }
    if ($parts[1] -notmatch '^[^|\s]+$') { Fail "invalid LRF id: ${path}: $line" }
    if ($ids.ContainsKey($parts[1])) { Fail "duplicate LRF id: $($parts[1])" }; $ids[$parts[1]] = $true
    $guardOk = $parts[2] -eq '*' -or $parts[2] -match '^[a-z][a-z0-9_-]*:[^|&]+(?:&[a-z][a-z0-9_-]*:[^|&]+)*$'
    if (!$guardOk) { Fail "unknown or malformed LRF guard: ${path}: $line" }
    foreach ($g in ($parts[2] -split '&')) { if ($g -ne '*' -and ($g -split ':',2)[0] -notin $script:GuardKeys) { Fail "unknown LRF guard key: ${path}: $line" } }
    if ($parts[3] -notin @('MUST','MUST_NOT','MAY')) { Fail "invalid LRF mode: ${path}: $line" }
    if ($parts[4] -notin $script:Effects) { Fail "unknown LRF effect: ${path}: $line" }
    Require-Text $parts[5] 'LRF payload'
    $records += [pscustomobject]@{ kind=$parts[0]; id=$parts[1]; guard=$parts[2]; mode=$parts[3]; effect=$parts[4]; payload=$parts[5]; path=$path }
  }
  if (!$header) { Fail "missing LRF header: $path" }
  $records
}

function Parse-Rules([string]$root,[string]$canonical) {
  $canonicalText=Get-Content -Raw -LiteralPath $canonical; $all = @(Parse-Lrf $canonical)
  $state = Join-Path $root 'rules\state.lrf'; $eng = Join-Path $root 'rules\engineering.lrf'; $doc = Join-Path $root 'rules\documentation.lrf'
  foreach ($p in @($state,$eng,$doc)) { $null = $all += @(Parse-Lrf $p) }
  $ids = @{}; foreach ($r in $all) { if ($ids.ContainsKey($r.id)) { Fail "duplicate routed LRF id: $($r.id)" }; $null = $ids[$r.id]=$true }
  $seq = Join-Path $root 'etl\main.seq'; if (!(Test-Path -LiteralPath $seq -PathType Leaf)) { Fail 'missing etl/main.seq' }
  $s = Get-Content -Raw -LiteralPath $seq
  foreach ($needle in @('@seq 1','@entry PRECEDENT','retry_exhausted','End(Conversation)','PRE_VERIFY','POST_VERIFY','LOCAL_REVIEW','ReviewReport','FULL_ETL')) { if (!$s.Contains($needle)) { Fail "etl/main.seq missing required structure: $needle" } }
  $vocab=@(); $vm=[regex]::Match($canonicalText,'(?s)E=\{(.*?)\}'); if($vm.Success){$vocab=@($vm.Groups[1].Value -split ',' | ForEach-Object {$_.Trim()} | Where-Object {$_ -match '^[A-Z_]+$'})}; if($vocab.Count -eq 0){Fail 'effect vocabulary missing from canonical rule'}
  $safe=@(@($vocab | Where-Object {$_ -match '^RO_'}) + @($all | Where-Object {$_.id -in @('intent.ro','intent.rw')} | ForEach-Object effect) | Select-Object -Unique); if ($safe.Count -eq 0) { Fail 'safe effects missing from intent rules' }
  $approval=@($vocab | Where-Object {$_ -notin $safe})
  $order=@(); $ord=$all | Where-Object {$_.id -eq 'route.order'} | Select-Object -First 1; if($ord){$om=[regex]::Match($ord.payload,'matching leaf union once:\s*([a-z]+(?:>[a-z]+)+)');if($om.Success){$order=@($om.Groups[1].Value -split '>')}}; if($order.Count -eq 0){Fail 'route precedence missing from hydrated rules'}
  $routeTags=$all | Where-Object {$_.id -eq 'route.tags'} | Select-Object -First 1
  if (!$routeTags) { Fail 'route tag contract missing from hydrated rules' }
  $knownTags=@([regex]::Matches($routeTags.payload,'(?:^|[:;]\s*)([a-z]+)=') | ForEach-Object {$_.Groups[1].Value})
  if ($knownTags.Count -eq 0 -or @($knownTags | Sort-Object -Unique).Count -ne $knownTags.Count) { Fail 'route tag contract is malformed or duplicated' }
  $stateText=Get-Content -Raw -LiteralPath $state; $stateFiles=@([regex]::Matches($stateText,'[a-z]+\.md')|ForEach-Object Value|Sort-Object -Unique); if($stateFiles.Count -eq 0){Fail 'state file set missing from routed rule'}
  $moduleTags=@{}
  foreach($ref in @(@{id='load.state';module='state';path=$state},@{id='load.engineering';module='engineering';path=$eng},@{id='load.documentation';module='documentation';path=$doc})) { $rr=$all|Where-Object {$_.id -eq $ref.id}|Select-Object -First 1; if(!$rr -or $rr.payload -notmatch '\[.+\]\([^)]*\.lrf\)' -or $rr.guard -notmatch '^tag:([^&]+)$'){Fail "unknown routed reference: $($ref.id)"}; $moduleTags[$ref.module]=$Matches[1] }
  $retry=[regex]::Match($s,'retry\([^)]*\)\s*≤\s*(\d+)'); if (!$retry.Success) { Fail 'retry cap missing from orchestration' }
  [pscustomobject]@{ Records=$all; RuleHash=(Hash-File $canonical); RoutedHashes=@{ state=(Hash-File $state); engineering=(Hash-File $eng); documentation=(Hash-File $doc); etl=(Hash-File $seq) }; Policy=[pscustomobject]@{ Vocabulary=$vocab; Approval=$approval; Safe=$safe; Order=$order; StateFiles=$stateFiles; RetryCap=[int]$retry.Groups[1].Value; KnownTags=$knownTags; ModuleTags=$moduleTags } }
}

function Validate-Fixture($data) {
  if ($data.schemaVersion -ne 1) { Fail 'case schemaVersion must be 1' }
  if (!(Is-Array $data.cases) -or @($data.cases).Count -lt 1) { Fail 'cases must be a nonempty array' }
  $seen=@{}; foreach ($c in @($data.cases)) {
    Require-Text $c.case_id 'case_id'; if ($seen.ContainsKey($c.case_id)) { Fail "duplicate case_id: $($c.case_id)" }; $seen[$c.case_id]=$true
    Require-Text $c.request 'request'; $p=$c.proposed_effect
    if ($null -eq $p) { Fail "missing proposed_effect: $($c.case_id)" }
    if ($p.effect -notin ($script:Effects + @('UNKNOWN','CONFLICT'))) { Fail "invalid proposed effect: $($c.case_id)" }
    if ($p.authorization -notin @('none','explicit_current_task')) { Fail "invalid authorization: $($c.case_id)" }
    $t=$p.triggers; if ($null -eq $t) { Fail "missing triggers: $($c.case_id)" }
    Require-Text $t.task 'triggers.task'; foreach ($n in @('tags','files','events','cases','changes','facts','phases')) { if ($t.PSObject.Properties.Name -contains $n -and !(Is-Array $t.$n)) { Fail "triggers.$n must be an array: $($c.case_id)" } }
    foreach ($n in @('tags','files','events','cases','changes','facts','phases')) { if ($t.PSObject.Properties.Name -contains $n -and @($t.$n | Sort-Object -Unique).Count -ne @($t.$n).Count) { Fail "duplicate triggers.$n value: $($c.case_id)" } }
    foreach ($n in @('expected_decision','expected_approval_required','expected_modules')) { if ($c.PSObject.Properties.Name -notcontains $n) { Fail "missing $n`: $($c.case_id)" } }
    if ($c.expected_decision -notin @('ALLOW','APPROVAL','STOP','FORBID') -or $c.expected_approval_required -isnot [bool] -or !(Is-Array $c.expected_modules)) { Fail "invalid expected result: $($c.case_id)" }
    $hasRetry = $c.PSObject.Properties.Name -contains 'retry_attempts'
    if ($hasRetry) { $retryValue=0; if (![int]::TryParse([string]$c.retry_attempts,[ref]$retryValue) -or $retryValue -lt 0 -or $retryValue -gt 4) { Fail "retry_attempts must be 0..4: $($c.case_id)" } }
    if ($c.PSObject.Properties.Name -contains 'requested_route' -and $c.requested_route -notin @('LOCAL_REVIEW','FULL_ETL')) { Fail "invalid requested_route: $($c.case_id)" }
  }
}

function Get-Modules($triggers,$policy) {
  $tags=@($triggers.tags); $files=@($triggers.files); $mods=[System.Collections.Generic.List[string]]::new()
  foreach ($tag in $tags) { if ($tag -notin $policy.KnownTags) { Fail "unknown route trigger: $tag" } }
  if ($triggers.task -ne 'change') { return @() }
  if ($tags -contains $policy.ModuleTags.state -or @($files | Where-Object { [IO.Path]::GetFileName($_).ToLowerInvariant() -in $policy.StateFiles }).Count -gt 0) { $mods.Add('state') }
  if ($tags -contains $policy.ModuleTags.engineering) { $mods.Add('engineering') }
  if ($tags -contains $policy.ModuleTags.documentation -or @($files | Where-Object { [IO.Path]::GetFileName($_) -match '^(README|SKILL)\.md$' -or $_ -match '(^|[\\/])docs?([\\/]|$)' }).Count -gt 0) { $mods.Add('documentation') }
  @($policy.Order | Where-Object { $_ -in $mods })
}

function Evaluate-Case($c,$policy) {
  $rules=$policy; $policy=$policy.Policy
  $p=$c.proposed_effect; $effect=[string]$p.effect; $decision='ALLOW'; $approval=$false; $reason='derived from proposed_effect and hydrated rules'; $modules=@()
  if ($effect -eq 'UNKNOWN') { $decision='STOP'; $reason='unknown effect'; }
  elseif ($effect -eq 'CONFLICT') { $decision='STOP'; $reason='equal-rank conflict has no precedence'; }
  elseif ($effect -in $policy.Approval) { $decision='APPROVAL'; $approval=$true; $reason='current-task approval gate required' }
  elseif ($effect -notin $policy.Safe) { $decision='STOP'; $reason='unclassified effect'; }
  $matching=@($rules.Records | Where-Object { $_.mode -eq 'MUST_NOT' -and $_.effect -eq $effect -and (Guard-Matches $_.guard $p.triggers) })
  $specific=@($matching | Where-Object { $_.guard -ne '*' -or $effect -eq 'CRED' -and $_.payload -match '(?i)credential|secret' })
  if ($specific.Count -gt 0) { $decision='FORBID'; $approval=$false; $reason='matching MUST_NOT rule: ' + $specific[0].id }
  elseif ($p.authorization -eq 'explicit_current_task' -and $decision -eq 'APPROVAL') { $decision='ALLOW'; $approval=$false; $reason='explicit current-task authorization' }
  if ($decision -eq 'ALLOW') { $modules=Get-Modules $p.triggers $policy }
  $requested='FULL_ETL'; if ($c.PSObject.Properties.Name -contains 'requested_route') { $requested=[string]$c.requested_route }
  $route='FULL_ETL'
  if ($decision -in @('STOP','FORBID')) { $route='STOP' }
  elseif ($requested -eq 'LOCAL_REVIEW') {
    $changes=@(); if ($p.triggers.PSObject.Properties.Name -contains 'changes') { $changes=@($p.triggers.changes) }
    $expansion=@($p.triggers.tags + $changes + $p.triggers.events | Where-Object { $_ -in @('code','config','schema','public','external','live','vcs','credential','scope-expansion') })
    $validReview=$c.PSObject.Properties.Name -contains 'review'
    if ($validReview) {
      $validReview=$c.review.query -is [string] -and ![string]::IsNullOrWhiteSpace($c.review.query) -and $c.review.query.Length -le 512 -and $c.review.query -notmatch '[\r\n]' -and (Is-Array $c.review.targets) -and @($c.review.targets).Count -ge 1 -and @($c.review.targets).Count -le 32
      if ($validReview) { $names=@($c.review.targets); $validReview=@($names|Where-Object{$_ -isnot [string] -or [string]::IsNullOrWhiteSpace($_) -or [IO.Path]::IsPathRooted($_) -or $_ -match '(^|[\\/])\.\.([\\/]|$)|[*?\[\]]'}).Count -eq 0 -and @($names|ForEach-Object{$_.ToLowerInvariant()}|Sort-Object -Unique).Count -eq $names.Count }
    }
    if (!$validReview) { $route='STOP'; $decision='STOP'; $reason='malformed local review request' }
    elseif ($effect -ne 'RO_LOCAL' -or $expansion.Count -gt 0 -or $p.triggers.task -ne 'read') { $route='FULL_ETL' }
    else { $route='LOCAL_REVIEW' }
  }
  if ($route -eq 'STOP') { if ($decision -ne 'FORBID') { $decision='STOP' }; $approval=$false; $modules=@(); $reason=if($reason -eq 'derived from proposed_effect and hydrated rules'){'explicit stop route'}else{$reason} }
  $attempts=0; if ($c.PSObject.Properties.Name -contains 'retry_attempts') { $attempts=[int]$c.retry_attempts }
  $terminal='completed'; $human=$false; $stopped=$false; $exit=0; $latest=$true
  if ($route -eq 'STOP') { $terminal='stopped'; $human=$true; $stopped=$true }
  if ($attempts -gt $policy.RetryCap) { $decision='STOP'; $approval=$false; $reason='retry_exhausted'; $terminal='retry_exhausted'; $human=$true; $stopped=$true; $exit=1; $latest=$false; $modules=@() }
  if ($attempts -gt $policy.RetryCap) { $route='STOP' }
  if ($c.PSObject.Properties.Name -contains 'expected_route' -and [string]$c.expected_route -ne $route) { $reason='oracle_mismatch: expected_route'; $exit=1 }
  $expected=$null; if ($c.PSObject.Properties.Name -contains 'expected_decision') { $expected=[string]$c.expected_decision; if ($decision -ne $expected) { $reason='oracle_mismatch: expected_decision'; $exit=1 } }
  if ($c.PSObject.Properties.Name -contains 'expected_approval_required' -and [bool]$approval -ne [bool]$c.expected_approval_required) { $reason='oracle_mismatch: expected_approval_required'; $exit=1 }
  if ($c.PSObject.Properties.Name -contains 'expected_modules' -and (@($modules) -join ',') -ne (@($c.expected_modules) -join ',')) { $reason='oracle_mismatch: expected_modules'; $exit=1 }
  [ordered]@{ case_id=[string]$c.case_id; route=$route; decision=$decision; effect=$effect; authorization=[string]$p.authorization; approval_required=$approval; modules=@($modules); reason=$reason; retry_attempts=$attempts; terminal=$terminal; human=$human; stopped=$stopped; process_exit_code=$exit; latest_update=$latest; request=[string]$c.request; guard_triggers=$p.triggers; review=if($c.PSObject.Properties.Name -contains 'review'){[ordered]@{query_sha256=(Hash-Text ([string]$c.review.query));target_count=@($c.review.targets).Count}}else{$null} }
}

function Guard-Matches([string]$guard,$triggers) {
  if ($guard -eq '*') { return $true }
  foreach ($part in ($guard -split '&')) { $kv=$part -split ':',2; $v=$triggers.($kv[0]+'s'); if ($kv[0] -eq 'task') { $v=$triggers.task }; if ($kv[0] -eq 'tag') { $v=$triggers.tags }; if ($null -eq $v -or ($v -is [string] -and $v -ne $kv[1]) -or ($v -isnot [string] -and $kv[1] -notin @($v))) { return $false } }
  return $true
}

function Write-AtomicJson([string]$path,$value) {
  $tmp="$path.$([guid]::NewGuid().ToString('N')).tmp"
  try {
    [IO.File]::WriteAllText($tmp,($value|ConvertTo-Json -Depth 20),[Text.UTF8Encoding]::new($false))
    Move-Item -LiteralPath $tmp -Destination $path -Force
  } finally {
    if (Test-Path -LiteralPath $tmp) { Remove-Item -LiteralPath $tmp -Force }
  }
}
function Hash-Text([string]$text) { $sha=[Security.Cryptography.SHA256]::Create(); try { ([BitConverter]::ToString($sha.ComputeHash([Text.Encoding]::UTF8.GetBytes($text)))).Replace('-','').ToUpperInvariant() } finally { $sha.Dispose() } }

try {
  $caseFull=[IO.Path]::GetFullPath($CasePath); $rulesFull=[IO.Path]::GetFullPath($RulesPath); $root=[IO.Path]::GetFullPath((Join-Path (Split-Path -Parent $rulesFull) '..')); $rules=@(Parse-Rules $root $rulesFull | Where-Object { $_.PSObject.Properties.Name -contains 'Records' } | Select-Object -Last 1); if ($rules.Count -ne 1) { Fail 'rule hydration did not yield one policy object' }; $rules=$rules[0]; $fixture=Read-Json $caseFull; Validate-Fixture $fixture
  $results=@($fixture.cases | ForEach-Object { Evaluate-Case $_ $rules }); $runId=[guid]::NewGuid().ToString(); $inputHash=Hash-File $caseFull
  $gitRevision='unavailable'; $gitDirty=$false; try { $gitRevision=(& git.exe -C $root rev-parse HEAD 2>$null).Trim(); $gitDirty=@(& git.exe -C $root status --porcelain --untracked-files=all 2>$null).Count -gt 0 } catch { }
  $caseSchema=Join-Path $root 'evaluation\case.schema.json'; $runSchema=Join-Path $root 'evaluation\run.schema.json'
  $overall=0; if (@($results|Where-Object {$_.process_exit_code -ne 0}).Count -gt 0) {$overall=1}; $latestBad=$overall -ne 0; $receipt=[ordered]@{schemaVersion=1;parserVersion=$script:ParserVersion;schemaVersionCase=1;run_id=$runId;revision=$gitRevision;git_revision=$gitRevision;git_dirty=$gitDirty;input_sha256=$inputHash;evaluator_sha256=(Hash-File $PSCommandPath);case_schema_sha256=(Hash-File $caseSchema);run_schema_sha256=(Hash-File $runSchema);rules_sha256=$rules.RuleHash;routed_hashes=$rules.RoutedHashes;cases=$results;process_exit_code=$overall;overall_exit_code=$overall;latest_updated=(-not $latestBad);generated_at=(Get-Date).ToUniversalTime().ToString('o')}
  New-Item -ItemType Directory -Force -Path $OutputDirectory | Out-Null; $stamp=(Get-Date).ToUniversalTime().ToString('yyyyMMddTHHmmssfffZ'); $runPath=Join-Path $OutputDirectory "run-$stamp-$runId.json"; Write-AtomicJson $runPath $receipt; if ($receipt.latest_updated) { Write-AtomicJson (Join-Path $OutputDirectory 'latest.json') $receipt }
  if ($Json) { $receipt|ConvertTo-Json -Depth 20 -Compress } else { Write-Output ("run_id={0} cases={1} exit={2}" -f $runId,$results.Count,$overall) }; exit $overall
} catch { if ($Json) { [ordered]@{schemaVersion=1;parserVersion=$script:ParserVersion;overall_exit_code=2;process_exit_code=2;error=$_.Exception.Message}|ConvertTo-Json -Compress } else { Write-Error $_.Exception.Message }; exit 2 }
