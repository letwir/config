$syncScript = Join-Path $PSScriptRoot 'sync-harness.ps1'
$pwshPath = (Get-Process -Id $PID).Path
$targetPaths = @(
    '.gemini\skills',
    '.gemini\config\plugins\ltw-skills-plugin\skills',
    '.codex\skills',
    '.agents\skills',
    '.claude\skills',
    '.opencode\skills'
)

function New-HealthFixture([string]$case) {
    $root = Join-Path ([System.IO.Path]::GetTempPath()) ("sync-harness-h01-{0}" -f [guid]::NewGuid().ToString('N'))
    $skills = Join-Path $root '.harness\skills'
    New-Item -ItemType Directory -Path (Join-Path $skills 'llm-memory') -Force | Out-Null

    if ($case -ne 'missing-binary') {
        New-Item -ItemType File -Path (Join-Path $skills 'llm-memory\llm-mem.exe') -Force | Out-Null
    }

    foreach ($relativePath in $targetPaths) {
        if ($case -eq 'missing-target' -and $relativePath -eq '.opencode\skills') { continue }

        $targetPath = Join-Path $root $relativePath
        New-Item -ItemType Directory -Path (Split-Path $targetPath -Parent) -Force | Out-Null

        if ($case -eq 'regular-directory' -and $relativePath -eq '.claude\skills') {
            New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
            continue
        }

        $junctionTarget = $skills
        if ($case -eq 'wrong-junction' -and $relativePath -eq '.agents\skills') {
            $junctionTarget = Join-Path $root 'wrong-skills'
            New-Item -ItemType Directory -Path $junctionTarget -Force | Out-Null
        }
        New-Item -ItemType Junction -Path $targetPath -Target $junctionTarget -ErrorAction Stop | Out-Null
    }

    return $root
}

function Invoke-HealthFixture([string]$root) {
    & $pwshPath -NoProfile -File $syncScript -Check -UserHome $root *> $null
    return [int]$LASTEXITCODE
}

function Remove-HealthFixture([string]$root) {
    $tempRoot = [System.IO.Path]::GetFullPath([System.IO.Path]::GetTempPath()).TrimEnd('\')
    $resolvedRoot = [System.IO.Path]::GetFullPath($root)
    if (-not $resolvedRoot.StartsWith("$tempRoot\sync-harness-h01-", [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to remove unexpected fixture path: $resolvedRoot"
    }
    Remove-Item -LiteralPath $resolvedRoot -Recurse -Force
}

function New-WhatIfFixture {
    $root = Join-Path ([System.IO.Path]::GetTempPath()) ("sync-harness-h02-{0}" -f [guid]::NewGuid().ToString('N'))
    New-Item -ItemType Directory -Path (Join-Path $root '.harness\skills') -Force | Out-Null

    $sourceSkill = Join-Path $root '.gemini\skills\source-skill'
    New-Item -ItemType Directory -Path $sourceSkill -Force | Out-Null
    [System.IO.File]::WriteAllText((Join-Path $sourceSkill 'SKILL.md'), "fixture`n", [System.Text.UTF8Encoding]::new($false))

    foreach ($relativePath in @(
        '.gemini\config\plugins\ltw-skills-plugin',
        '.codex',
        '.agents',
        '.claude',
        '.opencode'
    )) {
        New-Item -ItemType Directory -Path (Join-Path $root $relativePath) -Force | Out-Null
    }

    $wrongSkills = Join-Path $root 'wrong-skills'
    New-Item -ItemType Directory -Path $wrongSkills -Force | Out-Null
    New-Item -ItemType Junction -Path (Join-Path $root '.agents\skills') -Target $wrongSkills -ErrorAction Stop | Out-Null
    return $root
}

function Get-FixtureSnapshot([string]$root) {
    $rootPath = [System.IO.Path]::GetFullPath($root).TrimEnd('\')
    $pending = [System.Collections.Generic.Queue[string]]::new()
    $pending.Enqueue($rootPath)
    $entries = [System.Collections.Generic.List[object]]::new()

    while ($pending.Count -gt 0) {
        $directory = $pending.Dequeue()
        foreach ($item in @(Get-ChildItem -LiteralPath $directory -Force)) {
            $isReparsePoint = [bool]($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint)
            $relativePath = $item.FullName.Substring($rootPath.Length).TrimStart('\')
            $hash = if (-not $item.PSIsContainer) { (Get-FileHash -LiteralPath $item.FullName -Algorithm SHA256).Hash } else { $null }
            $target = if ($isReparsePoint) { @($item.Target) -join ';' } else { $null }
            $entries.Add([pscustomobject]@{
                Path   = $relativePath
                Type   = if ($item.PSIsContainer) { 'directory' } else { 'file' }
                Reparse = $isReparsePoint
                Target = $target
                Length = if ($item.PSIsContainer) { $null } else { $item.Length }
                Hash   = $hash
            })
            if ($item.PSIsContainer -and -not $isReparsePoint) { $pending.Enqueue($item.FullName) }
        }
    }

    return @($entries | Sort-Object Path) | ConvertTo-Json -Depth 4 -Compress
}

function Remove-WhatIfFixture([string]$root) {
    $tempRoot = [System.IO.Path]::GetFullPath([System.IO.Path]::GetTempPath()).TrimEnd('\')
    $resolvedRoot = [System.IO.Path]::GetFullPath($root)
    if (-not $resolvedRoot.StartsWith("$tempRoot\sync-harness-h02-", [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to remove unexpected fixture path: $resolvedRoot"
    }
    Remove-Item -LiteralPath $resolvedRoot -Recurse -Force
}

Describe 'sync-harness health-check exit contract' {
    It 'returns the documented result object for a healthy fixture' {
        $fixtureRoot = New-HealthFixture 'healthy'
        try {
            $result = . $syncScript -Check -UserHome $fixtureRoot -PassThru
            $result.Healthy | Should Be $true
            $result.FailureCount | Should Be 0
            @($result.Results).Count | Should Be 14
        } finally {
            Remove-HealthFixture $fixtureRoot
        }
    }

    foreach ($case in @('healthy', 'missing-target', 'wrong-junction', 'regular-directory', 'missing-binary')) {
        It "returns the expected exit code for $case" {
            $fixtureRoot = New-HealthFixture $case
            try {
                $actual = Invoke-HealthFixture $fixtureRoot
                $expected = if ($case -eq 'healthy') { 0 } else { 1 }
                $actual | Should Be $expected
            } finally {
                Remove-HealthFixture $fixtureRoot
            }
        }
    }
}

Describe 'sync-harness WhatIf contract' {
    foreach ($mode in @('Junction', 'Copy')) {
        It "does not change the fixture in $mode mode" {
            $fixtureRoot = New-WhatIfFixture
            try {
                $before = Get-FixtureSnapshot $fixtureRoot
                & $pwshPath -NoProfile -File $syncScript -Consolidate -LinkSkills -Mode $mode -UserHome $fixtureRoot -WhatIf *> $null
                $LASTEXITCODE | Should Be 0
                $after = Get-FixtureSnapshot $fixtureRoot
                $after | Should Be $before
            } finally {
                Remove-WhatIfFixture $fixtureRoot
            }
        }
    }

    It 'does not create a missing SSOT during consolidation preview' {
        $fixtureRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("sync-harness-h02-{0}" -f [guid]::NewGuid().ToString('N'))
        New-Item -ItemType Directory -Path $fixtureRoot -Force | Out-Null
        try {
            $before = Get-FixtureSnapshot $fixtureRoot
            & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $fixtureRoot -WhatIf *> $null
            $LASTEXITCODE | Should Be 0
            $after = Get-FixtureSnapshot $fixtureRoot
            $after | Should Be $before
        } finally {
            Remove-WhatIfFixture $fixtureRoot
        }
    }
}

function New-H03Fixture {
    $root = Join-Path ([System.IO.Path]::GetTempPath()) ("sync-harness-h03-{0}" -f [guid]::NewGuid().ToString('N'))
    New-Item -ItemType Directory -Path (Join-Path $root '.harness\skills') -Force | Out-Null
    New-Item -ItemType Directory -Path (Join-Path $root '.agents\skills\shared') -Force | Out-Null
    New-Item -ItemType Directory -Path (Join-Path $root '.codex\skills\shared') -Force | Out-Null
    return $root
}

function Remove-H03Fixture([string]$root) {
    $tempRoot = [System.IO.Path]::GetFullPath([System.IO.Path]::GetTempPath()).TrimEnd('\')
    $resolvedRoot = [System.IO.Path]::GetFullPath($root)
    if (-not $resolvedRoot.StartsWith("$tempRoot\sync-harness-h03-", [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to remove unexpected fixture path: $resolvedRoot"
    }
    Remove-Item -LiteralPath $resolvedRoot -Recurse -Force
}

Describe 'sync-harness H-03 consolidation and restore' {
    It 'does not continue the default pipeline after a Fail-policy conflict' {
        $root = New-H03Fixture
        try {
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\shared\a.txt') -Value 'one'
            Set-Content -LiteralPath (Join-Path $root '.codex\skills\shared\a.txt') -Value 'two'
            $before = Get-FixtureSnapshot $root
            $null = & $pwshPath -NoProfile -File $syncScript -UserHome $root -ConflictPolicy Fail 2>&1
            $LASTEXITCODE | Should Be 1
            (Get-FixtureSnapshot $root) | Should Be $before
            (Test-Path -LiteralPath (Join-Path $root '.gemini\skills')) | Should Be $false
            (Test-Path -LiteralPath (Join-Path $root '.claude\skills')) | Should Be $false
            (Test-Path -LiteralPath (Join-Path $root '.harness\manifests')) | Should Be $false
        } finally { Remove-H03Fixture $root }
    }

    It 'fails all writes on a conflict and records every conflict' {
        $root = New-H03Fixture
        try {
            Set-Content -LiteralPath (Join-Path $root '.harness\skills\shared.txt') -Value 'old'
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\shared\a.txt') -Value 'one'
            Set-Content -LiteralPath (Join-Path $root '.codex\skills\shared\a.txt') -Value 'two'
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\other.txt') -Value 'one'
            Set-Content -LiteralPath (Join-Path $root '.codex\skills\other.txt') -Value 'two'
            $before = Get-FixtureSnapshot $root
            $output = & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $root -ConflictPolicy Fail 2>&1
            $LASTEXITCODE | Should Be 1
            @($output | Where-Object { $_ -match 'Conflict A:' }).Count | Should Be 2
            (Get-FixtureSnapshot $root) | Should Be $before
            (Test-Path (Join-Path $root '.harness\manifests')) | Should Be $false
            (Test-Path (Join-Path $root '.harness\backups')) | Should Be $false
        } finally { Remove-H03Fixture $root }
    }

    It 'uses PreferSource, writes a manifest and restores after drift with Force' {
        $root = New-H03Fixture
        try {
            Set-Content -LiteralPath (Join-Path $root '.harness\skills\shared.txt') -Value 'destination'
            Set-Content -LiteralPath (Join-Path $root '.harness\skills\.hidden-seed') -Value 'hidden'
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\shared\a.txt') -Value 'source'
            Set-Content -LiteralPath (Join-Path $root '.codex\skills\shared\a.txt') -Value 'higher-priority-source'
            $null = & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $root -ConflictPolicy PreferSource 2>&1
            $LASTEXITCODE | Should Be 0
            (Get-Content -Raw (Join-Path $root '.harness\skills\shared\a.txt')).Trim() | Should Be 'higher-priority-source'
            $manifest = Get-ChildItem -LiteralPath (Join-Path $root '.harness\manifests') -Filter '*.json' | Select-Object -First 1
            $manifest | Should Not Be $null
            $manifestData = Get-Content -Raw $manifest.FullName | ConvertFrom-Json
            $manifestData.RunId | Should Not BeNullOrEmpty
            $manifestData.Status | Should Be 'Applied'
            $manifestData.Completed | Should Be $true
            $manifestData.Error | Should Be $null
            @($manifestData.Operations | Where-Object { $_.RelativePath -eq 'shared\a.txt' }).Count | Should Be 2
            ($manifestData.Operations | Where-Object { $_.RelativePath -eq 'shared\a.txt' } | Select-Object -First 1).Decision | Should Be 'Apply'
            [string]::IsNullOrEmpty([string](($manifestData.Operations | Where-Object { $_.RelativePath -eq 'shared\a.txt' } | Select-Object -First 1).BeforeSource)) | Should Be $true
            $sourceConflict = @($manifestData.Conflicts | Where-Object { $_.Source -eq (Join-Path $root '.codex\skills') -and $_.RelativePath -eq 'shared\a.txt' }) | Select-Object -First 1
            $sourceConflict.PreviousSource | Should Be (Join-Path $root '.agents\skills')
            $sourceConflict.Decision | Should Be 'PreferSource'
            (Test-Path (Join-Path $root '.harness\backups')) | Should Be $true
            Set-Content -LiteralPath (Join-Path $root '.harness\skills\drift.txt') -Value 'drift'
            & $pwshPath -NoProfile -File $syncScript -RestoreManifest $manifest.FullName -UserHome $root *> $null
            $LASTEXITCODE | Should Be 1
            & $pwshPath -NoProfile -File $syncScript -RestoreManifest $manifest.FullName -UserHome $root -Force *> $null
            $LASTEXITCODE | Should Be 0
            (Test-Path (Join-Path $root '.harness\skills\shared\a.txt')) | Should Be $false
            (Get-Content -Raw (Join-Path $root '.harness\skills\shared.txt')).Trim() | Should Be 'destination'
            (Get-Content -Raw (Join-Path $root '.harness\skills\.hidden-seed')).Trim() | Should Be 'hidden'
        } finally { Remove-H03Fixture $root }
    }

    It 'keeps destination content under PreferDestination' {
        $root = New-H03Fixture
        try {
            New-Item -ItemType Directory -Path (Join-Path $root '.harness\skills\shared') -Force | Out-Null
            Set-Content -LiteralPath (Join-Path $root '.harness\skills\shared\a.txt') -Value 'destination'
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\shared\a.txt') -Value 'source'
            $null = & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $root -ConflictPolicy PreferDestination 2>&1
            $LASTEXITCODE | Should Be 0
            (Get-Content -Raw (Join-Path $root '.harness\skills\shared\a.txt')).Trim() | Should Be 'destination'
        } finally { Remove-H03Fixture $root }
    }

    It 'records an incomplete manifest on a deterministic partial failure' {
        $root = New-H03Fixture
        try {
            Set-Content -LiteralPath (Join-Path $root '.agents\skills\shared\a.txt') -Value 'source'
            & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $root -TestFailAfter 1 *> $null
            $LASTEXITCODE | Should Be 1
            $manifest = Get-ChildItem -LiteralPath (Join-Path $root '.harness\manifests') -Filter '*.json' | Select-Object -First 1
            $manifest | Should Not Be $null
            $manifestData = Get-Content -Raw $manifest.FullName | ConvertFrom-Json
            $manifestData.Completed | Should Be $false
            $manifestData.Status | Should Be 'Partial'
            $manifestData.Error | Should Be 'ConsolidationFailed'
        } finally { Remove-H03Fixture $root }
    }

    It 'seeds a binary from the injected UserHome with its exact source hash' {
        $root = New-H03Fixture
        try {
            $source = Join-Path $root '.codex\skills\llm-memory\llm-mem.exe'
            New-Item -ItemType Directory -Path (Split-Path $source -Parent) -Force | Out-Null
            [System.IO.File]::WriteAllBytes($source, [byte[]](1, 2, 3, 4, 5))
            $expectedHash = (Get-FileHash -LiteralPath $source -Algorithm SHA256).Hash
            $null = & $pwshPath -NoProfile -File $syncScript -Consolidate -UserHome $root -ConflictPolicy PreferSource 2>&1
            $LASTEXITCODE | Should Be 0
            $destination = Join-Path $root '.harness\skills\llm-memory\llm-mem.exe'
            (Test-Path -LiteralPath $destination -PathType Leaf) | Should Be $true
            (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash | Should Be $expectedHash
        } finally { Remove-H03Fixture $root }
    }
}
