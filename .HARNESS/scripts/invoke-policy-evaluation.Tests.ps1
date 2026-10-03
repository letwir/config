$script:Runner = Join-Path $PSScriptRoot 'invoke-policy-evaluation.ps1'
$script:NL = [Environment]::NewLine

function Write-Utf8([string]$p,[string]$v) { [IO.File]::WriteAllText($p,$v,[Text.UTF8Encoding]::new($false)) }

# Minimal valid LRF header for fixture files
function Get-LrfHead { '@lrf=1|aud=test-fixture|scope=test' + $script:NL }

function New-EvalFixture([string]$mode='valid') {
  $tmpBase = [IO.Path]::GetTempPath()
  $leafName = 'h05-' + [guid]::NewGuid().ToString('N')
  $root = Join-Path $tmpBase $leafName
  $eval  = Join-Path $root 'evaluation'
  $rules = Join-Path $root 'rules'
  $etl   = Join-Path $root 'etl'
  New-Item -ItemType Directory -Force -Path $eval,$rules,$etl | Out-Null

  Copy-Item -LiteralPath (Join-Path $PSScriptRoot '..\evaluation\case.schema.json') -Destination $eval
  Copy-Item -LiteralPath (Join-Path $PSScriptRoot '..\evaluation\run.schema.json')  -Destination $eval

  $head = Get-LrfHead
  $NL   = $script:NL

  # --- BOOTSTRAP.lrf: vocab R-record, routing contracts, bootstrap L-records ---
  $bootstrap  = $head
  $bootstrap += 'R|vocab|*|MUST|RO_LOCAL|E={RO_LOCAL,RO_PUBLIC,LW_SCOPE,EXT_WRITE,RELEASE,LIVE_WRITE,VCS_WRITE,DESTRUCT,CRED,CHARGE,PROD_DEP}' + $NL
  $bootstrap += 'R|route.tags|*|MUST|RO_LOCAL|after ALLOW collect all: code=software; doc=README; state=state files; agy=external; subagent=invocation' + $NL
  $bootstrap += 'R|route.order|*|MUST|RO_LOCAL|matching leaf union once: state>engineering>documentation' + $NL
  $bootstrap += 'L|bootstrap.manual|*|MUST|RO_LOCAL|[manual](MANUAL.lrf)' + $NL
  $bootstrap += 'L|bootstrap.load|*|MUST|RO_LOCAL|[load](LOAD.lrf)' + $NL

  # --- MANUAL.lrf ---
  $manual  = $head
  $manual += 'R|manual.one|*|MUST|RO_LOCAL|manual content' + $NL

  # --- LOAD.lrf: the module L-records live here ---
  $load  = $head
  $load += 'L|load.state|tag:state|MUST|RO_LOCAL|[state](state.lrf)' + $NL
  $load += 'L|load.engineering|tag:code|MUST|RO_LOCAL|[engineering](engineering.lrf)' + $NL
  $load += 'L|load.documentation|tag:doc|MUST|RO_LOCAL|[documentation](documentation.lrf)' + $NL

  # --- state.lrf ---
  $stateLrf  = $head
  $stateLrf += 'R|state.one|*|MUST|RO_LOCAL|state' + $NL
  $stateLrf += 'R|state.files|*|MUST|RO_LOCAL|decisions.md method.md knowledge.md issues.md memo.md history.md diary.md' + $NL

  # --- engineering.lrf ---
  $engLrf  = $head
  $engLrf += 'R|engineering.one|*|MUST|RO_LOCAL|engineering' + $NL

  # --- documentation.lrf ---
  $docLrf  = $head
  $docLrf += 'R|documentation.one|*|MUST|RO_LOCAL|documentation' + $NL

  # Unicode ≤ for retry cap
  $leq = [char]0x2264
  $arr = [char]0x21D2

  # --- ETL sequence files ---
  $mainSeq = "@seq 2$NL@entry PRECEDENT${NL}retry(gate) $leq 3${NL}retry_exhausted${NL}LOCAL_REVIEW${NL}research.seq${NL}change.seq${NL}finish.seq${NL}research-local.seq"
  $resSeq  = "@seq 1${NL}@entry RESEARCH${NL}PLAN${NL}FINISH"
  $chgSeq  = "@seq 1${NL}@entry CHANGE${NL}PRE_VERIFY${NL}CODE @agy${NL}POST_VERIFY${NL}PASS $arr FINISH"
  $finSeq  = "@seq 1${NL}@entry FINISH${NL}mem.etl${NL}llm-mem.ingest${NL}End(Conversation)"
  $rlSeq   = "@seq 1${NL}@entry RESEARCH_LOCAL${NL}LOCAL_REVIEW${NL}ReviewReport${NL}FINISH"

  # Apply mode mutations before writing
  switch ($mode) {
    'malformed'          { $bootstrap = $head + 'R|one|*|MUST|RO_LOCAL' + $NL }
    'unknown'            { $bootstrap = $head + 'R|one|*|MUST|NOPE|one' + $NL }
    'duplicate'          { $bootstrap += $bootstrap }
    'header'             { $bootstrap = 'R|vocab|*|MUST|RO_LOCAL|E={RO_LOCAL,RO_PUBLIC,LW_SCOPE,EXT_WRITE,RELEASE,LIVE_WRITE,VCS_WRITE,DESTRUCT,CRED,CHARGE,PROD_DEP}' + $NL }
    'malformed-effects'  { $bootstrap = $bootstrap.Replace('E={RO_LOCAL,RO_PUBLIC,LW_SCOPE,EXT_WRITE,RELEASE,LIVE_WRITE,VCS_WRITE,DESTRUCT,CRED,CHARGE,PROD_DEP}','E={}') }
    'missing-load-ref'   { $bootstrap = $bootstrap.Replace("L|bootstrap.load|*|MUST|RO_LOCAL|[load](LOAD.lrf)$NL",'') }
    'wrong-entry-main'   { $mainSeq = $mainSeq.Replace('@entry PRECEDENT','@entry WRONG') }
    'wrong-entry-leaf'   { $chgSeq = $chgSeq.Replace('@entry CHANGE','@entry WRONG') }
    'wrong-entry-suffix' { $chgSeq = $chgSeq.Replace('@entry CHANGE','@entry CHANGE_EXTRA') }
    'seq-version-suffix' { $chgSeq = $chgSeq.Replace('@seq 1','@seq 10') }
    'wrong-load-target'  { $load = $load.Replace('L|load.state|tag:state|MUST|RO_LOCAL|[state](state.lrf)','L|load.state|tag:state|MUST|RO_LOCAL|[state](wrong.lrf)') }
    'missing-leaf'       { } # handled by removing file after write
  }

  Write-Utf8 (Join-Path $rules 'BOOTSTRAP.lrf') $bootstrap
  Write-Utf8 (Join-Path $rules 'MANUAL.lrf')    $manual
  Write-Utf8 (Join-Path $rules 'LOAD.lrf')      $load
  Write-Utf8 (Join-Path $rules 'state.lrf')     $stateLrf
  Write-Utf8 (Join-Path $rules 'engineering.lrf') $engLrf
  Write-Utf8 (Join-Path $rules 'documentation.lrf') $docLrf

  Write-Utf8 (Join-Path $etl 'main.seq')          $mainSeq
  Write-Utf8 (Join-Path $etl 'research.seq')       $resSeq
  Write-Utf8 (Join-Path $etl 'change.seq')         $chgSeq
  Write-Utf8 (Join-Path $etl 'finish.seq')         $finSeq
  Write-Utf8 (Join-Path $etl 'research-local.seq') $rlSeq

  if ($mode -eq 'missing-leaf') { Remove-Item -LiteralPath (Join-Path $etl 'change.seq') -Force }

  $trigBase = [ordered]@{task='read';tags=@();files=@();events=@();cases=@()}
  $p   = [ordered]@{effect='RO_LOCAL';authorization='none';triggers=$trigBase}
  $case = [ordered]@{case_id='one';request='read';proposed_effect=$p;expected_decision='ALLOW';expected_approval_required=$false;expected_modules=@()}
  if ($mode -eq 'schema') { $case.Remove('proposed_effect') }
  if ($mode -eq 'retry')  { $case.retry_attempts=4; $case.expected_decision='STOP' }

  $obj = [ordered]@{schemaVersion=1;cases=@($case)}
  $casePath = Join-Path $eval 'case.json'
  Write-Utf8 $casePath ($obj | ConvertTo-Json -Depth 10)
  [pscustomobject]@{Root=$root;Case=$casePath;Rules=(Join-Path $rules 'BOOTSTRAP.lrf');Out=(Join-Path $eval 'runs')}
}

function Clean-EvalFixture($f) {
  $full   = [IO.Path]::GetFullPath($f.Root)
  $parent = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd([IO.Path]::DirectorySeparatorChar,[IO.Path]::AltDirectorySeparatorChar)
  $leaf   = Split-Path $full -Leaf
  $actualParent = [IO.Path]::GetFullPath((Split-Path $full -Parent)).TrimEnd([IO.Path]::DirectorySeparatorChar,[IO.Path]::AltDirectorySeparatorChar)
  if ($actualParent -ne $parent) { throw "Clean-EvalFixture: path escapes temp: $full" }
  if ($leaf -notmatch '^h05-[a-fA-F0-9]{32}$') { throw "Clean-EvalFixture: leaf name invalid: $leaf" }
  if (Test-Path -LiteralPath $full) { Remove-Item -LiteralPath $full -Recurse -Force }
}

function Run-Eval($f) {
  $o = & powershell.exe -NoProfile -ExecutionPolicy Bypass -File $script:Runner -CasePath $f.Case -RulesPath $f.Rules -OutputDirectory $f.Out -Json 2>&1
  [pscustomobject]@{Exit=$LASTEXITCODE;Text=($o -join '')}
}

Describe 'H-05 policy evaluator' {

  It 'derives decisions and writes unique atomic receipts' {
    $f = New-EvalFixture
    try {
      $r = Run-Eval $f
      if ($r.Exit -ne 0) { Write-Host $r.Text }
      $r.Exit | Should Be 0
      (Get-ChildItem $f.Out -Filter 'run-*.json').Count | Should Be 1
      (Get-Content (Join-Path $f.Out 'latest.json') -Raw | ConvertFrom-Json).cases[0].decision | Should Be 'ALLOW'
    } finally { Clean-EvalFixture $f }
  }

  It 'fails closed for malformed unknown duplicate and header-invalid rules' {
    foreach ($m in 'malformed','unknown','duplicate','header') {
      $f = New-EvalFixture $m
      try { (Run-Eval $f).Exit | Should Be 2 } finally { Clean-EvalFixture $f }
    }
  }

  It 'fails closed for split layout violations: missing leaf, load ref, malformed effects, wrong entry, wrong load target' {
    foreach ($m in 'missing-leaf','missing-load-ref','malformed-effects','wrong-entry-main','wrong-entry-leaf','wrong-load-target','wrong-entry-suffix','seq-version-suffix') {
      $f = New-EvalFixture $m
      try { (Run-Eval $f).Exit | Should Be 2 } finally { Clean-EvalFixture $f }
    }
  }

  It 'fails closed for schema input and rejects nonzero retry exhaustion' {
    foreach ($m in 'schema','retry') {
      $f = New-EvalFixture $m
      try {
        $r = Run-Eval $f
        if ($m -eq 'schema') {
          $r.Exit | Should Be 2
        } else {
          $r.Exit | Should Be 1
          ($r.Text | ConvertFrom-Json).cases[0].terminal | Should Be 'retry_exhausted'
          Test-Path (Join-Path $f.Out 'latest.json') | Should Be $false
        }
      } finally { Clean-EvalFixture $f }
    }
  }

  It 'rejects unknown and duplicate route triggers and records expectation mismatches' {
    foreach ($m in 'unknown-trigger','duplicate-trigger','mismatch') {
      $f = New-EvalFixture
      try {
        $o = Get-Content -Raw $f.Case | ConvertFrom-Json
        if ($m -eq 'unknown-trigger')    { $o.cases[0].proposed_effect.triggers.tags = @('unknown-route') }
        elseif ($m -eq 'duplicate-trigger') { $o.cases[0].proposed_effect.triggers.tags = @('code','code') }
        else                             { $o.cases[0].expected_decision = 'STOP' }
        Write-Utf8 $f.Case ($o | ConvertTo-Json -Depth 10)
        $r = Run-Eval $f
        if ($m -eq 'mismatch') { $r.Exit | Should Be 1; ($r.Text | ConvertFrom-Json).cases[0].reason | Should Match 'oracle_mismatch' }
        else                   { $r.Exit | Should Be 2 }
      } finally { Clean-EvalFixture $f }
    }
  }

  It 'fails closed when the receipt destination cannot be created' {
    $f = New-EvalFixture
    try {
      Write-Utf8 $f.Out 'occupied'
      (Run-Eval $f).Exit | Should Be 2
      @(Get-ChildItem (Split-Path $f.Out) -Filter '*.tmp').Count | Should Be 0
    } finally { Clean-EvalFixture $f }
  }

  It 'routes a bounded local review without claiming collector execution' {
    $f = New-EvalFixture
    try {
      $o = Get-Content -Raw $f.Case | ConvertFrom-Json
      $o.cases[0] | Add-Member -NotePropertyName requested_route -NotePropertyValue 'LOCAL_REVIEW'
      $o.cases[0] | Add-Member -NotePropertyName expected_route  -NotePropertyValue 'LOCAL_REVIEW'
      $o.cases[0] | Add-Member -NotePropertyName review -NotePropertyValue ([pscustomobject]@{query='x';targets=@('note.md')})
      Write-Utf8 $f.Case ($o | ConvertTo-Json -Depth 10)
      $r = Run-Eval $f
      $r.Exit | Should Be 0
      ($r.Text | ConvertFrom-Json).cases[0].route | Should Be 'LOCAL_REVIEW'
      ($r.Text | ConvertFrom-Json).cases[0].review.target_count | Should Be 1
    } finally { Clean-EvalFixture $f }
  }

  It 'fails closed for malformed local review structure' {
    $f = New-EvalFixture
    try {
      $o = Get-Content -Raw $f.Case | ConvertFrom-Json
      $o.cases[0] | Add-Member -NotePropertyName requested_route -NotePropertyValue 'LOCAL_REVIEW'
      $o.cases[0] | Add-Member -NotePropertyName expected_route  -NotePropertyValue 'STOP'
      $o.cases[0].expected_decision = 'STOP'
      Write-Utf8 $f.Case ($o | ConvertTo-Json -Depth 10)
      $r = Run-Eval $f
      $r.Exit | Should Be 0
      ($r.Text | ConvertFrom-Json).cases[0].route | Should Be 'STOP'
    } finally { Clean-EvalFixture $f }
  }

  It 'Guard-Matches handles missing optional event array without error' {
    # A case whose triggers object has NO events property; a task-guard rule must still match correctly
    $f = New-EvalFixture
    try {
      $docPath = Join-Path $f.Root 'rules\documentation.lrf'
      $docContent = (Get-Content -Encoding UTF8 -Raw -LiteralPath $docPath) + 'R|fixture.task-deny|task:read|MUST_NOT|RO_LOCAL|guard proof' + $script:NL
      Write-Utf8 $docPath $docContent

      $o = Get-Content -Raw $f.Case | ConvertFrom-Json
      # Remove events from triggers to exercise Guard-Matches safe property access
      $t = $o.cases[0].proposed_effect.triggers
      $props = $t.PSObject.Properties | Where-Object { $_.Name -ne 'events' }
      $newT = [ordered]@{}
      foreach ($pr in $props) { $newT[$pr.Name] = $pr.Value }
      $o.cases[0].proposed_effect.triggers = $newT
      $o.cases[0].expected_decision = 'FORBID'
      Write-Utf8 $f.Case ($o | ConvertTo-Json -Depth 10)
      $r = Run-Eval $f
      $r.Exit | Should Be 0
      ($r.Text | ConvertFrom-Json).cases[0].decision | Should Be 'FORBID'
    } finally { Clean-EvalFixture $f }
  }

}
