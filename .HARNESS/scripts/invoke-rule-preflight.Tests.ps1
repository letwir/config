$sut = Join-Path $PSScriptRoot 'invoke-rule-preflight.ps1'

Describe 'rule reference preflight' {
    It 'resolves the live subagent rule graph into a compact receipt' {
        $json = & $sut -Task subagent -Tag subagent
        $LASTEXITCODE | Should Be 0
        $receipt = $json | ConvertFrom-Json
        $receipt.status | Should Be 'PASS'
        $receipt.file_count | Should BeGreaterThan 2
    }

    It 'rejects invalid task and tag syntax before traversal' {
        & pwsh.exe -NoProfile -File $sut -Task 'bad task' -Tag subagent *> $null
        $LASTEXITCODE | Should Not Be 0
        & pwsh.exe -NoProfile -File $sut -Task subagent -Tag 'bad tag' *> $null
        $LASTEXITCODE | Should Not Be 0
    }

    It 'rejects a missing canonical root' {
        & pwsh.exe -NoProfile -File $sut -Task subagent -Root (Join-Path (Split-Path $PSScriptRoot -Parent) 'rules\missing.lrf') *> $null
        $LASTEXITCODE | Should Not Be 0
    }

    It 'detects a stale receipt' {
        $receiptPath = Join-Path $TestDrive 'receipt.json'
        & $sut -Task subagent -Tag subagent -OutFile $receiptPath *> $null
        $receipt = Get-Content -Raw $receiptPath | ConvertFrom-Json
        $receipt.files[0].sha256 = '0' * 64
        $receipt | ConvertTo-Json -Depth 6 | Set-Content -LiteralPath $receiptPath -Encoding utf8
        & pwsh.exe -NoProfile -File $sut -VerifyReceipt $receiptPath *> $null
        $LASTEXITCODE | Should Not Be 0
    }

    It 'rejects an empty forged receipt' {
        $receiptPath = Join-Path $TestDrive 'empty-receipt.json'
        '{"schema":"lrf-preflight/v1","parser":"invoke-rule-preflight.ps1/v1","task":"subagent","root":"C:\\Users\\letwir\\.harness\\rules\\BOOTSTRAP.lrf","files":[],"selected":[]}' | Set-Content -LiteralPath $receiptPath -Encoding utf8
        & pwsh.exe -NoProfile -File $sut -VerifyReceipt $receiptPath *> $null
        $LASTEXITCODE | Should Not Be 0
    }

    It 'enforces graph resource limits' {
        & pwsh.exe -NoProfile -File $sut -Task subagent -Tag subagent -MaxModules 1 *> $null
        $LASTEXITCODE | Should Not Be 0
        & pwsh.exe -NoProfile -File $sut -Task subagent -Tag subagent -MaxBytes 10 *> $null
        $LASTEXITCODE | Should Not Be 0
    }
}
