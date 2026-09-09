$script:Runner = Join-Path $PSScriptRoot 'invoke-policy-evaluation.ps1'
function Write-Utf8([string]$p,[string]$v) { [IO.File]::WriteAllText($p,$v,[Text.UTF8Encoding]::new($false)) }
function New-EvalFixture([string]$mode='valid') {
  $root=Join-Path ([IO.Path]::GetTempPath()) ('h05-'+[guid]::NewGuid().ToString('N')); $eval=Join-Path $root 'evaluation'; $rules=Join-Path $root 'rules'; $etl=Join-Path $root 'etl'; New-Item -ItemType Directory -Force -Path $eval,$rules,$etl|Out-Null
  Copy-Item -LiteralPath (Join-Path $PSScriptRoot '..\evaluation\case.schema.json') -Destination $eval
  Copy-Item -LiteralPath (Join-Path $PSScriptRoot '..\evaluation\run.schema.json') -Destination $eval
  $head='@lrf=1|aud=GPT-5.6|scope=x'+[Environment]::NewLine+'F=e in E={RO_LOCAL,RO_PUBLIC,LW_SCOPE,EXT_WRITE,RELEASE,LIVE_WRITE,VCS_WRITE,DESTRUCT,CRED,CHARGE,PROD_DEP}'+[Environment]::NewLine; $rec='R|one|*|MUST|RO_LOCAL|one'+[Environment]::NewLine+'R|route.tags|*|MUST|RO_LOCAL|after ALLOW collect all: code=software; doc=README; state=state files; agy=external; subagent=invocation'+[Environment]::NewLine+'R|route.order|*|MUST|RO_LOCAL|matching leaf union once: state>engineering>documentation'+[Environment]::NewLine+'L|load.state|tag:state|MUST|RO_LOCAL|[state](state.lrf)'+[Environment]::NewLine+'L|load.engineering|tag:code|MUST|RO_LOCAL|[engineering](engineering.lrf)'+[Environment]::NewLine+'L|load.documentation|tag:doc|MUST|RO_LOCAL|[documentation](documentation.lrf)'+[Environment]::NewLine
  if($mode -eq 'malformed'){$rec='R|one|*|MUST|RO_LOCAL'+[Environment]::NewLine}; if($mode -eq 'unknown'){$rec='R|one|*|MUST|NOPE|one'+[Environment]::NewLine}; if($mode -eq 'duplicate'){$rec=$rec+$rec}; if($mode -eq 'header'){$head=''}
  Write-Utf8 (Join-Path $rules 'LLM_REF_RULE.md') ($head+$rec); foreach($n in 'state','engineering','documentation'){$nr='R|'+$n+'.one|*|MUST|RO_LOCAL|'+$n+[Environment]::NewLine;if($n -eq 'state'){$nr += 'R|state.files|*|MUST|RO_LOCAL|decisions.md method.md knowledge.md issues.md memo.md history.md diary.md'+[Environment]::NewLine};Write-Utf8 (Join-Path $rules ($n+'.lrf')) ($head+$nr)}
  Write-Utf8 (Join-Path $etl 'main.seq') ("@seq 1`n@entry PRECEDENT`nforall gate: retry(gate, ADV or DENY) " + [char]0x2264 + " 3`nretry_exhausted`nEnd(Conversation)`nPRE_VERIFY`nPOST_VERIFY`nLOCAL_REVIEW`nReviewReport`nFULL_ETL")
  $p=[ordered]@{effect='RO_LOCAL';authorization='none';triggers=[ordered]@{task='read';tags=@();files=@();events=@();cases=@()}}; $case=[ordered]@{case_id='one';request='read';proposed_effect=$p;expected_decision='ALLOW';expected_approval_required=$false;expected_modules=@()}; if($mode -eq 'schema'){$case.Remove('proposed_effect')}; if($mode -eq 'retry'){$case.retry_attempts=4;$case.expected_decision='STOP'}
  $obj=[ordered]@{schemaVersion=1;cases=@($case)}; $casePath=Join-Path $eval 'case.json'; Write-Utf8 $casePath ($obj|ConvertTo-Json -Depth 10); [pscustomobject]@{Root=$root;Case=$casePath;Rules=(Join-Path $rules 'LLM_REF_RULE.md');Out=(Join-Path $eval 'runs')}
}
function Run-Eval($f){$o=& powershell.exe -NoProfile -ExecutionPolicy Bypass -File $script:Runner -CasePath $f.Case -RulesPath $f.Rules -OutputDirectory $f.Out -Json 2>&1; [pscustomobject]@{Exit=$LASTEXITCODE;Text=($o -join '')}}
Describe 'H-05 policy evaluator' {
  It 'derives decisions and writes unique atomic receipts' { $f=New-EvalFixture;try{$r=Run-Eval $f;if($r.Exit -ne 0){Write-Host $r.Text};$r.Exit|Should Be 0;(Get-ChildItem $f.Out -Filter 'run-*.json').Count|Should Be 1;(Get-Content (Join-Path $f.Out 'latest.json') -Raw|ConvertFrom-Json).cases[0].decision|Should Be 'ALLOW'}finally{Remove-Item $f.Root -Recurse -Force} }
  It 'fails closed for malformed unknown duplicate and header-invalid rules' { foreach($m in 'malformed','unknown','duplicate','header'){$f=New-EvalFixture $m;try{(Run-Eval $f).Exit|Should Be 2}finally{Remove-Item $f.Root -Recurse -Force}} }
  It 'fails closed for schema input and rejects nonzero retry exhaustion' { foreach($m in 'schema','retry'){$f=New-EvalFixture $m;try{$r=Run-Eval $f;if($m -eq 'schema'){$r.Exit|Should Be 2}else{$r.Exit|Should Be 1;($r.Text|ConvertFrom-Json).cases[0].terminal|Should Be 'retry_exhausted';Test-Path (Join-Path $f.Out 'latest.json')|Should Be $false}}finally{Remove-Item $f.Root -Recurse -Force}} }
  It 'rejects unknown and duplicate route triggers and records expectation mismatches' {
    foreach($m in 'unknown-trigger','duplicate-trigger','mismatch') {
      $f=New-EvalFixture
      try {
        $o=Get-Content -Raw $f.Case|ConvertFrom-Json
        if($m -eq 'unknown-trigger'){$o.cases[0].proposed_effect.triggers.tags=@('unknown-route')}
        elseif($m -eq 'duplicate-trigger'){$o.cases[0].proposed_effect.triggers.tags=@('code','code')}
        else{$o.cases[0].expected_decision='STOP'}
        Write-Utf8 $f.Case ($o|ConvertTo-Json -Depth 10)
        $r=Run-Eval $f
        if($m -eq 'mismatch'){$r.Exit|Should Be 1;($r.Text|ConvertFrom-Json).cases[0].reason|Should Match 'oracle_mismatch'}else{$r.Exit|Should Be 2}
      } finally { Remove-Item $f.Root -Recurse -Force }
    }
  }
  It 'fails closed when the receipt destination cannot be created' {
    $f=New-EvalFixture
    try { Write-Utf8 $f.Out 'occupied'; (Run-Eval $f).Exit|Should Be 2; @(Get-ChildItem (Split-Path $f.Out) -Filter '*.tmp').Count|Should Be 0 }
    finally { Remove-Item $f.Root -Recurse -Force }
  }
  It 'routes a bounded local review without claiming collector execution' {
    $f=New-EvalFixture
    try {
      $o=Get-Content -Raw $f.Case|ConvertFrom-Json
      $o.cases[0]|Add-Member -NotePropertyName requested_route -NotePropertyValue 'LOCAL_REVIEW';$o.cases[0]|Add-Member -NotePropertyName expected_route -NotePropertyValue 'LOCAL_REVIEW';$o.cases[0]|Add-Member -NotePropertyName review -NotePropertyValue ([pscustomobject]@{query='x';targets=@('note.md')})
      Write-Utf8 $f.Case ($o|ConvertTo-Json -Depth 10);$r=Run-Eval $f;$r.Exit|Should Be 0;($r.Text|ConvertFrom-Json).cases[0].route|Should Be 'LOCAL_REVIEW';($r.Text|ConvertFrom-Json).cases[0].review.target_count|Should Be 1
    } finally { Remove-Item $f.Root -Recurse -Force }
  }
  It 'fails closed for malformed local review structure' {
    $f=New-EvalFixture
    try { $o=Get-Content -Raw $f.Case|ConvertFrom-Json;$o.cases[0]|Add-Member -NotePropertyName requested_route -NotePropertyValue 'LOCAL_REVIEW';$o.cases[0]|Add-Member -NotePropertyName expected_route -NotePropertyValue 'STOP';$o.cases[0].expected_decision='STOP';Write-Utf8 $f.Case ($o|ConvertTo-Json -Depth 10);$r=Run-Eval $f;$r.Exit|Should Be 0;($r.Text|ConvertFrom-Json).cases[0].route|Should Be 'STOP' }
    finally { Remove-Item $f.Root -Recurse -Force }
  }
}
