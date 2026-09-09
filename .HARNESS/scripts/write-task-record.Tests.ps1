$here = Split-Path -Parent $PSCommandPath
$writer = Join-Path $here 'write-task-record.ps1'
$schema = Join-Path (Split-Path -Parent $here) 'evaluation\task-record.schema.json'

Describe 'H-07 task record writer' {
    BeforeEach {
        $root = Join-Path ([IO.Path]::GetTempPath()) ('h07-' + [guid]::NewGuid().ToString('N'))
        New-Item -ItemType Directory -Path $root | Out-Null
        $task = Join-Path $root 'memo\T-01'; New-Item -ItemType Directory -Path $task -Force | Out-Null
        [IO.File]::WriteAllText((Join-Path $task 'knowledge.md'),'k')
        [IO.File]::WriteAllText((Join-Path $task 'diary.md'),'d')
        [IO.File]::WriteAllText((Join-Path $task 'walkthrough.md'),'w')
        [IO.File]::WriteAllText((Join-Path $task 'evaluation.json'),'{"task_id":"T-01","axes":{"acceptance":1.0}}')
        $hostExe = (Get-Process -Id $PID).Path
    }
    AfterEach {
        $temp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath()).TrimEnd('\') + '\'
        $full = [IO.Path]::GetFullPath($root)
        if ($full.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $full) -match '^h07-[0-9a-f]{32}$') {
            Remove-Item -LiteralPath $full -Recurse -Force -ErrorAction SilentlyContinue
        }
    }

    It 'writes a schema-valid complete record from fixed artifacts' {
        $out = & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root -IngestStatus SUCCESS -IngestReasonCode completed -IngestReceiptId ingest-1 -EvaluationStatus SUCCESS -EvaluationReasonCode completed -EvaluationReceiptId eval-1
        $LASTEXITCODE | Should Be 0
        $recordPath=Join-Path $task 'record.json';$r=($out-join '')|ConvertFrom-Json
        $r.memory_report_complete|Should Be $true;$r.artifacts.knowledge.path|Should Be 'memo/T-01/knowledge.md';$r.external_diary_mirror.status|Should Be 'N_A'
        $validator=Join-Path $root 'validate-schema.ps1'
        [IO.File]::WriteAllText($validator,"param([string]`$Data,[string]`$Schema)`nGet-Content -Raw -LiteralPath `$Data|Test-Json -SchemaFile `$Schema -ErrorAction Stop")
        $valid=& pwsh.exe -NoProfile -File $validator -Data $recordPath -Schema $schema
        ($valid-join '')|Should Be 'True'
    }

    It 'preserves numeric evaluation bytes and records their exact hash' {
        $before=(Get-FileHash (Join-Path $task 'evaluation.json') -Algorithm SHA256).Hash.ToLowerInvariant()
        & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root | Out-Null;$LASTEXITCODE|Should Be 0
        $r=Get-Content -Raw (Join-Path $task 'record.json')|ConvertFrom-Json
        $r.artifacts.evaluation.sha256|Should Be $before
        ((Get-Content -Raw (Join-Path $task 'evaluation.json')|ConvertFrom-Json).axes.acceptance)|Should Be 1
    }

    It 'rejects traversal, reserved names, and actual-directory case mismatch' {
        & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId '..\x' -RootPath $root | Out-Null;$LASTEXITCODE|Should Not Be 0
        & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId CON -RootPath $root | Out-Null;$LASTEXITCODE|Should Not Be 0
        $out=& $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId t-01 -RootPath $root
        $LASTEXITCODE|Should Not Be 0;((($out-join '')|ConvertFrom-Json).reason_code)|Should Be 'invalid_task_id'
    }

    It 'rejects a task-directory reparse point' {
        $outside=Join-Path $root 'outside';New-Item -ItemType Directory $outside|Out-Null
        foreach($n in @('knowledge.md','diary.md','walkthrough.md')){[IO.File]::WriteAllText((Join-Path $outside $n),'x')}
        [IO.File]::WriteAllText((Join-Path $outside 'evaluation.json'),'{"task_id":"J-01","axes":{"a":1}}')
        New-Item -ItemType Junction -Path (Join-Path $root 'memo\J-01') -Target $outside|Out-Null
        $out=& $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId J-01 -RootPath $root
        $LASTEXITCODE|Should Not Be 0;((($out-join '')|ConvertFrom-Json).reason_code)|Should Be 'reparse_point'
    }

    It 'rejects incoherent or free-form receipt data without persisting it' {
        $out=& $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root -IngestStatus SUCCESS -IngestReasonCode pending_approval -IngestReceiptId secret-value
        $LASTEXITCODE|Should Not Be 0;((($out-join '')|ConvertFrom-Json).reason_code)|Should Be 'io_error';Test-Path (Join-Path $task 'record.json')|Should Be $false
        $out=& $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root -IngestStatus FAILED -IngestReasonCode 'raw failure details'
        $LASTEXITCODE|Should Not Be 0;((($out-join '')|ConvertFrom-Json).reason_code)|Should Be 'io_error'
    }

    It 'preserves the prior record and removes temporary files on replacement failure' {
        & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root | Out-Null;$before=Get-Content -Raw (Join-Path $task 'record.json')
        $env:H07_INJECT_REPLACE_FAILURE='1'
        try { & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root | Out-Null;$LASTEXITCODE|Should Not Be 0 } finally { Remove-Item Env:H07_INJECT_REPLACE_FAILURE -ErrorAction SilentlyContinue }
        (Get-Content -Raw (Join-Path $task 'record.json'))|Should Be $before
        @(Get-ChildItem $task -Force|Where-Object Name -match '^\.record\..+\.(tmp|bak)$').Count|Should Be 0
    }

    It 'returns a bounded lock timeout without corrupting an existing record' {
        & $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root | Out-Null;$before=Get-Content -Raw (Join-Path $task 'record.json')
        $lock=[IO.File]::Open((Join-Path $task '.record.lock'),[IO.FileMode]::OpenOrCreate,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
        try { $out=& $hostExe -NoProfile -ExecutionPolicy Bypass -File $writer -TaskId T-01 -RootPath $root -LockTimeoutMs 100;$LASTEXITCODE|Should Not Be 0;((($out-join '')|ConvertFrom-Json).reason_code)|Should Be 'lock_timeout' } finally { $lock.Dispose() }
        (Get-Content -Raw (Join-Path $task 'record.json'))|Should Be $before
    }

    It 'serializes concurrent complete snapshots and never mixes receipt states' {
        $out1=Join-Path $root 'p1.out';$out2=Join-Path $root 'p2.out'
        $common=@('-NoProfile','-ExecutionPolicy','Bypass','-File',$writer,'-TaskId','T-01','-RootPath',$root)
        $p1=Start-Process -FilePath $hostExe -ArgumentList ($common+@('-IngestStatus','SUCCESS','-IngestReasonCode','completed','-IngestReceiptId','ingest-1','-EvaluationStatus','SUCCESS','-EvaluationReasonCode','completed','-EvaluationReceiptId','eval-1')) -RedirectStandardOutput $out1 -PassThru -WindowStyle Hidden
        $p2=Start-Process -FilePath $hostExe -ArgumentList ($common+@('-IngestStatus','FAILED','-IngestReasonCode','process_error')) -RedirectStandardOutput $out2 -PassThru -WindowStyle Hidden
        $p1.WaitForExit(15000)|Should Be $true;$p2.WaitForExit(15000)|Should Be $true;$p1.WaitForExit();$p2.WaitForExit()
        ((Get-Content -Raw $out1|ConvertFrom-Json).task_id)|Should Be 'T-01';((Get-Content -Raw $out2|ConvertFrom-Json).task_id)|Should Be 'T-01'
        $r=Get-Content -Raw (Join-Path $task 'record.json')|ConvertFrom-Json
        @('SUCCESS','FAILED') -contains $r.llm_memory_ingest.status|Should Be $true
        if($r.llm_memory_ingest.status -eq 'SUCCESS'){$r.llm_memory_eval.status|Should Be 'SUCCESS';$r.memory_report_complete|Should Be $true}else{$r.memory_report_complete|Should Be $false}
    }
}
