$script:Doctor = Join-Path $PSScriptRoot 'check-rule-distribution.ps1'

function Set-Utf8NoBom([string]$path, [string]$value) {
    [IO.File]::WriteAllText($path, $value, [Text.UTF8Encoding]::new($false))
}

function New-Fixture([string]$mode = 'exact') {
    $root = Join-Path ([IO.Path]::GetTempPath()) ('rule-doctor-' + [guid]::NewGuid().ToString('N'))
    $harness = Join-Path $root 'harness'; $userRoot = Join-Path $root 'home'
    New-Item -ItemType Directory -Force -Path (Join-Path $harness 'rules'), (Join-Path $userRoot '.codex\rules') | Out-Null
    $canonical = '@lrf=1|aud=GPT-5.6|scope=x' + [Environment]::NewLine + 'R|alpha|*|MUST|RO_LOCAL|x' + [Environment]::NewLine + 'R|beta|*|MUST|RO_LOCAL|x' + [Environment]::NewLine
    $clientLrf = $canonical
    if ($mode -eq 'allowed') { $clientLrf = $clientLrf -replace '(?m)^R\|beta\|[^\r\n]*\r?\n', '' }
    if ($mode -eq 'unexpected') { $clientLrf = $clientLrf + 'R|gamma|*|MUST|RO_LOCAL|x' + [Environment]::NewLine }
    if ($mode -eq 'broken') { $clientLrf = $clientLrf + 'R|alpha|*|MUST|RO_LOCAL|x' + [Environment]::NewLine }
    $canonicalPath = Join-Path $harness 'rules\LLM_REF_RULE.md'
    $lrfPath = Join-Path $userRoot '.codex\rules\LLM_REF_RULE.md'; $entryPath = Join-Path $userRoot '.codex\AGENTS.md'
    Set-Utf8NoBom $canonicalPath $canonical
    Set-Utf8NoBom $lrfPath $clientLrf
    $entryReference = ([IO.Path]::GetFullPath($lrfPath)).Replace('\', '/')
    Set-Utf8NoBom $entryPath "READ $entryReference"
    $sha = (Get-FileHash -LiteralPath $lrfPath -Algorithm SHA256).Hash
    $entrySha = (Get-FileHash -LiteralPath $entryPath -Algorithm SHA256).Hash
    $actualIds = @('alpha','beta'); if ($mode -eq 'allowed') { $actualIds = @('alpha') }; if ($mode -eq 'unexpected') { $actualIds = @('alpha','beta','gamma') }
    $allowed = @(); if ($mode -eq 'allowed') { $allowed = @('beta') }
    $manifest = [ordered]@{ schemaVersion = 1; canonical = '{HarnessRoot}\rules\LLM_REF_RULE.md'; canonicalSha256=(Get-FileHash -LiteralPath $canonicalPath -Algorithm SHA256).Hash; entries = @(
        [ordered]@{ name='codex'; kind='Codex'; required=$true; entry='{UserHome}\.codex\AGENTS.md'; chain=@(
            [ordered]@{ path='{UserHome}\.codex\AGENTS.md'; kind='entry'; sha256=$entrySha; contains='{UserHomeForward}/.codex/rules/LLM_REF_RULE.md'; ids=@(); allowedDifferenceIds=@() },
            [ordered]@{ path='{UserHome}\.codex\rules\LLM_REF_RULE.md'; kind='lrf'; sha256=$sha; ids=$actualIds; allowedDifferenceIds=$allowed }
        ) }
    ) }
    $manifestPath = Join-Path $harness 'rules\client-distribution.json'
    Set-Utf8NoBom $manifestPath ($manifest | ConvertTo-Json -Depth 10)
    [pscustomobject]@{ Root=$root; Harness=$harness; UserRoot=$userRoot; Manifest=$manifestPath }
}

function Invoke-Fixture($fixture, [switch]$Json) {
    $args = @('-NoProfile','-ExecutionPolicy','Bypass','-File',$script:Doctor,'-ManifestPath',$fixture.Manifest,'-UserHome',$fixture.UserRoot)
    if ($Json) { $args += '-Json' }
    $output = & powershell.exe @args 2>&1
    [pscustomobject]@{ ExitCode=$LASTEXITCODE; Output=@($output) }
}

Describe 'check-rule-distribution doctor' {
    It 'accepts an exact chain and leaves the fixture unchanged' {
        $f=New-Fixture; try { $before=(Get-ChildItem -LiteralPath $f.Root -Recurse -Force | ForEach-Object FullName); $r=Invoke-Fixture $f -Json; $after=(Get-ChildItem -LiteralPath $f.Root -Recurse -Force | ForEach-Object FullName); $r.ExitCode | Should Be 0; ($before -join '|') | Should Be ($after -join '|'); (($r.Output -join '') | ConvertFrom-Json).Results[0].Classification | Should Be 'Exact' } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'classifies an allowed LRF ID difference' {
        $f=New-Fixture allowed; try { $r=Invoke-Fixture $f -Json; $r.ExitCode | Should Be 0; (($r.Output -join '') | ConvertFrom-Json).Results[0].Classification | Should Be 'AllowedDifference' } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'rejects an unexpected LRF difference' {
        $f=New-Fixture unexpected; try { $r=Invoke-Fixture $f -Json; $r.ExitCode | Should Not Be 0; (($r.Output -join '') | ConvertFrom-Json).Results[0].Classification | Should Be 'UnexpectedDifference' } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'rejects a payload drift even when the LRF ID set is unchanged' {
        $f=New-Fixture; try {
            $clientLrf=Join-Path $f.UserRoot '.codex\rules\LLM_REF_RULE.md'
            Set-Utf8NoBom $clientLrf ((Get-Content -Raw $clientLrf).Replace('R|beta|*|MUST|RO_LOCAL|x','R|beta|*|MUST|RO_LOCAL|changed'))
            $m=Get-Content -Raw $f.Manifest|ConvertFrom-Json
            $m.entries[0].chain[1].sha256=(Get-FileHash -LiteralPath $clientLrf -Algorithm SHA256).Hash
            Set-Utf8NoBom $f.Manifest ($m|ConvertTo-Json -Depth 10)
            $r=Invoke-Fixture $f -Json
            $r.ExitCode | Should Not Be 0
            (($r.Output -join '')|ConvertFrom-Json).Results[0].Classification | Should Be 'UnexpectedDifference'
        } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'rejects a broken entry reference as EntryMismatch' {
        $f=New-Fixture; try {
            $entry=Join-Path $f.UserRoot '.codex\AGENTS.md'
            Set-Utf8NoBom $entry 'READ a different rule'
            $m=Get-Content -Raw $f.Manifest|ConvertFrom-Json
            $m.entries[0].chain[0].sha256=(Get-FileHash -LiteralPath $entry -Algorithm SHA256).Hash
            Set-Utf8NoBom $f.Manifest ($m|ConvertTo-Json -Depth 10)
            $r=Invoke-Fixture $f -Json
            $r.ExitCode | Should Not Be 0
            (($r.Output -join '')|ConvertFrom-Json).Results[0].Classification | Should Be 'EntryMismatch'
        } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'rejects malformed or duplicate LRF IDs' {
        $f=New-Fixture broken; try { $r=Invoke-Fixture $f -Json; $r.ExitCode | Should Not Be 0; (($r.Output -join '') | ConvertFrom-Json).Results[0].Classification | Should Be 'UnexpectedDifference' } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'treats an absent optional entry as NotConfigured and success' {
        $f=New-Fixture; try { $m=Get-Content -Raw $f.Manifest | ConvertFrom-Json; $m.entries[0].required=$false; $m.entries[0].name='claude'; $m.entries[0].kind='Claude'; $m.entries[0].entry='{UserHome}\.claude\CLAUDE.md'; $m.entries[0].chain[0].path=$m.entries[0].entry; Set-Utf8NoBom $f.Manifest ($m | ConvertTo-Json -Depth 10); $r=Invoke-Fixture $f -Json; $r.ExitCode | Should Be 0; (($r.Output -join '') | ConvertFrom-Json).Results[0].Classification | Should Be 'NotConfigured' } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
    It 'fails closed on an invalid manifest' {
        $f=New-Fixture; try { $m=Get-Content -Raw $f.Manifest | ConvertFrom-Json; $m.schemaVersion=2; Set-Utf8NoBom $f.Manifest ($m | ConvertTo-Json -Depth 10); $r=Invoke-Fixture $f -Json; $r.ExitCode | Should Be 2; (($r.Output -join '') | ConvertFrom-Json).Passed | Should Be $false } finally { Remove-Item -LiteralPath $f.Root -Recurse -Force }
    }
}
