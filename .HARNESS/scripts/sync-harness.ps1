<#
.SYNOPSIS
    sync-harness.ps1 - LLM Agent Harness SSOT, Skill Synchronization & Junction Linker
.DESCRIPTION
    Consolidates skills across all agent harnesses (.gemini, .codex, .agents, .claude, .opencode)
    into C:\Users\letwir\.harness\skills as the Single Source of Truth (SSOT).
    Creates transparent NTFS directory junctions pointing each agent's skills/ directory directly
    to .harness\skills, ensuring zero-drift instant synchronization.
#>

[CmdletBinding(SupportsShouldProcess = $true)]
param(
    [Parameter()]
    [switch]$Consolidate,

    [Parameter()]
    [switch]$LinkSkills,

    [Parameter()]
    [switch]$Check,

    [Parameter()]
    [ValidateSet('Junction', 'Copy')]
    [string]$Mode = 'Junction',

    [Parameter()]
    [switch]$Force,

    [Parameter()]
    [string]$UserHome = $env:USERPROFILE,

    [Parameter()]
    [switch]$PassThru,

    [Parameter()]
    [ValidateSet('Fail', 'PreferSource', 'PreferDestination')]
    [string]$ConflictPolicy = 'Fail',

    [Parameter()]
    [string]$RestoreManifest,

    [Parameter(DontShow = $true)]
    [int]$TestFailAfter = -1
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$HarnessRoot = Join-Path $UserHome '.harness'
$HarnessSkills = Join-Path $HarnessRoot 'skills'
$HarnessRules = Join-Path $HarnessRoot 'rules'
$script:OperationFailed = $false

$AgentSkillTargets = @(
    @{ Name = 'gemini';              Path = (Join-Path $UserHome '.gemini\skills') },
    @{ Name = 'gemini-plugin-ltw';    Path = (Join-Path $UserHome '.gemini\config\plugins\ltw-skills-plugin\skills') },
    @{ Name = 'codex';               Path = (Join-Path $UserHome '.codex\skills') },
    @{ Name = 'agents';              Path = (Join-Path $UserHome '.agents\skills') },
    @{ Name = 'claude';              Path = (Join-Path $UserHome '.claude\skills') },
    @{ Name = 'opencode';            Path = (Join-Path $UserHome '.opencode\skills') }
)

function Write-Info([string]$msg) {
    Write-Host "[INFO] $msg" -ForegroundColor Cyan
}

function Write-Success([string]$msg) {
    Write-Host "[OK]   $msg" -ForegroundColor Green
}

function Write-WarnMsg([string]$msg) {
    Write-Host "[WARN] $msg" -ForegroundColor Yellow
}

function Write-ErrMsg([string]$msg) {
    Write-Host "[ERR]  $msg" -ForegroundColor Red
}

function Test-IsJunction([string]$path) {
    if (-not (Test-Path $path)) { return $false }
    $item = Get-Item $path -Force
    return [bool]($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint)
}

function Get-JunctionTarget([string]$path) {
    if (-not (Test-IsJunction $path)) { return $null }
    $item = Get-Item $path -Force
    return $item.Target
}

function ConvertTo-NormalizedPath([string]$path) {
    if ([string]::IsNullOrWhiteSpace($path)) { return $null }
    return [System.IO.Path]::GetFullPath($path).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
}

# -----------------------------------------------------------------------------
# 1. Consolidate skills into .harness/skills (SSOT)
# -----------------------------------------------------------------------------
function Get-TreeSnapshot([string]$root) {
    $root = ConvertTo-NormalizedPath $root
    $items = [System.Collections.Generic.List[object]]::new()
    if (-not (Test-Path -LiteralPath $root -PathType Container)) { return @() }
    $queue = [System.Collections.Generic.Queue[string]]::new()
    $queue.Enqueue($root)
    while ($queue.Count -gt 0) {
        $dir = $queue.Dequeue()
        foreach ($item in @(Get-ChildItem -LiteralPath $dir -Force)) {
            $relative = $item.FullName.Substring($root.Length).TrimStart('\')
            $reparse = [bool]($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint)
            $type = if ($item.PSIsContainer) { if ($reparse) { 'reparse-directory' } else { 'directory' } } else { if ($reparse) { 'reparse-file' } else { 'file' } }
            $hash = if ($item.PSIsContainer -or $reparse) { $null } else { (Get-FileHash -LiteralPath $item.FullName -Algorithm SHA256).Hash }
            $target = if ($reparse) { @($item.Target) -join ';' } else { $null }
            $items.Add([pscustomobject]@{ RelativePath = $relative; Type = $type; Hash = $hash; Target = $target })
            if ($item.PSIsContainer -and -not $reparse) { $queue.Enqueue($item.FullName) }
        }
    }
    return @($items | Sort-Object RelativePath)
}

function Test-ContainedPath([string]$root, [string]$candidate) {
    $r = (ConvertTo-NormalizedPath $root) + '\'
    $c = ConvertTo-NormalizedPath $candidate
    return $c.StartsWith($r, [System.StringComparison]::OrdinalIgnoreCase) -or [string]::Equals($c, $r.TrimEnd('\'), [System.StringComparison]::OrdinalIgnoreCase)
}

function Test-NoReparseAncestors([string]$root, [string]$path) {
    if (-not (Test-ContainedPath $root $path)) { return $false }
    $root = ConvertTo-NormalizedPath $root
    $current = ConvertTo-NormalizedPath $path
    while ($true) {
        if (Test-Path -LiteralPath $current) {
            $item = Get-Item -LiteralPath $current -Force
            if (($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint) -and $current -ne $root) { return $false }
        }
        if ([string]::Equals($current, $root, [System.StringComparison]::OrdinalIgnoreCase)) { break }
        $parent = Split-Path $current -Parent
        if ([string]::Equals($parent, $current, [System.StringComparison]::OrdinalIgnoreCase)) { return $false }
        $current = $parent
    }
    return $true
}

function Get-ManifestPath {
    return (Join-Path $HarnessRoot (Join-Path 'manifests' ('sync-{0}.json' -f (Get-Date -Format 'yyyyMMdd-HHmmss-fff'))))
}

function Write-JsonArtifact([string]$path, [object]$value) {
    if ($PSCmdlet.ShouldProcess($path, 'Write synchronization manifest')) {
        $parent = Split-Path $path -Parent
        if (-not (Test-Path -LiteralPath $parent)) { New-Item -ItemType Directory -Path $parent -Force | Out-Null }
        $value | ConvertTo-Json -Depth 12 | Set-Content -LiteralPath $path -Encoding utf8NoBOM
        return $true
    }
    return $false
}

function Invoke-SkillConsolidation {
    Write-Info "Starting skill consolidation into SSOT: $HarnessSkills (ConflictPolicy: $ConflictPolicy)"
    $sourceDirs = @(
        (Join-Path $UserHome '.agents\skills'),
        (Join-Path $UserHome '.codex\skills'),
        (Join-Path $UserHome '.gemini\skills'),
        (Join-Path $UserHome '.gemini\config\plugins\ltw-skills-plugin\skills')
    )
    $usableSources = @()
    foreach ($src in $sourceDirs) {
        if (-not (Test-Path -LiteralPath $src -PathType Container)) { continue }
        if (Test-IsJunction $src) { Write-Info "Skipping $src because it is already a junction."; continue }
        if (-not (Test-NoReparseAncestors $UserHome $src)) { throw "Unsafe reparse source root: $src" }
        $usableSources += $src
    }

    $destinationExisted = Test-Path -LiteralPath $HarnessSkills -PathType Container
    $destinationSnapshot = if ($destinationExisted) { Get-TreeSnapshot $HarnessSkills } else { @() }
    $destByPath = @{}
    $sourceByPath = @{}
    foreach ($e in $destinationSnapshot) {
        $key = $e.RelativePath.ToLowerInvariant()
        $destByPath[$key] = $e
        $sourceByPath[$key] = $HarnessSkills
    }
    $plans = [System.Collections.Generic.List[object]]::new()
    $conflicts = [System.Collections.Generic.List[object]]::new()
    foreach ($src in $usableSources) {
        foreach ($entry in (Get-TreeSnapshot $src)) {
            $key = $entry.RelativePath.ToLowerInvariant()
            $dest = if ($destByPath.ContainsKey($key)) { $destByPath[$key] } else { $null }
            $isConflict = $null -ne $dest -and ($dest.Type -ne $entry.Type -or ($entry.Type -eq 'file' -and $dest.Hash -ne $entry.Hash))
            $previousSource = if ($sourceByPath.ContainsKey($key)) { $sourceByPath[$key] } else { $null }
            if ($isConflict) {
                $conflicts.Add([pscustomobject]@{ Source = $src; PreviousSource = $previousSource; RelativePath = $entry.RelativePath; ExistingType = $dest.Type; ExistingHash = $dest.Hash; IncomingType = $entry.Type; IncomingHash = $entry.Hash; Policy = $ConflictPolicy; Decision = $ConflictPolicy })
            }
            $plans.Add([pscustomobject]@{ Source = $src; Entry = $entry; Conflict = $isConflict; BinarySeed = $false; Before = $dest; BeforeSource = $previousSource; Decision = if ($isConflict) { $ConflictPolicy } else { 'Apply' } })
            if ($null -eq $dest -or $ConflictPolicy -eq 'PreferSource') {
                $destByPath[$key] = $entry
                $sourceByPath[$key] = $src
            }
        }
    }
    # Keep the historical binary seed, but make selection and conflict handling
    # deterministic and visible in the same plan as skill files.
    $candidateBins = @(
        (Join-Path $UserHome '.gemini\config\plugins\ltw-skills-plugin\skills\llm-memory\llm-mem.exe'),
        (Join-Path $UserHome '.codex\skills\llm-memory\llm-mem.exe')
    )
    $selectedBin = @($candidateBins | Where-Object { Test-Path -LiteralPath $_ -PathType Leaf } | Select-Object -First 1)
    if ($selectedBin.Count -gt 0) {
        $binPath = [string]$selectedBin[0]
        $binEntry = [pscustomobject]@{ RelativePath = 'llm-memory\llm-mem.exe'; Type = 'file'; Hash = (Get-FileHash -LiteralPath $binPath -Algorithm SHA256).Hash; Target = $null }
        $binKey = $binEntry.RelativePath.ToLowerInvariant()
        $binDest = if ($destByPath.ContainsKey($binKey)) { $destByPath[$binKey] } else { $null }
        $binConflict = $null -ne $binDest -and ($binDest.Type -ne 'file' -or $binDest.Hash -ne $binEntry.Hash)
        $binPreviousSource = if ($sourceByPath.ContainsKey($binKey)) { $sourceByPath[$binKey] } else { $null }
        if ($binConflict) {
            $conflicts.Add([pscustomobject]@{ Source = $binPath; PreviousSource = $binPreviousSource; RelativePath = $binEntry.RelativePath; ExistingType = $binDest.Type; ExistingHash = $binDest.Hash; IncomingType = $binEntry.Type; IncomingHash = $binEntry.Hash; Policy = $ConflictPolicy; Decision = $ConflictPolicy })
        }
        if (-not (Test-NoReparseAncestors $UserHome $binPath)) { throw "Unsafe reparse binary source: $binPath" }
        $plans.Add([pscustomobject]@{ Source = $binPath; Entry = $binEntry; Conflict = $binConflict; BinarySeed = $true; Before = $binDest; BeforeSource = $binPreviousSource; Decision = if ($binConflict) { $ConflictPolicy } else { 'Apply' } })
        if ($null -eq $binDest -or $ConflictPolicy -eq 'PreferSource') {
            $destByPath[$binKey] = $binEntry
            $sourceByPath[$binKey] = $binPath
        }
    } else {
        Write-WarnMsg 'No compiled llm-mem.exe found to seed into the SSOT.'
    }
    if ($conflicts.Count -gt 0) {
        foreach ($c in $conflicts) { Write-WarnMsg ("Conflict {0}: {1}" -f $c.Source, $c.RelativePath) }
        if ($ConflictPolicy -eq 'Fail') {
            $script:OperationFailed = $true
            Write-ErrMsg "Consolidation aborted: $($conflicts.Count) conflict(s); no SSOT or artifact writes performed."
            return
        }
    }

    $backupPath = Join-Path $HarnessRoot (Join-Path 'backups' ('sync-{0}' -f (Get-Date -Format 'yyyyMMdd-HHmmss-fff')))
    $manifestPath = Get-ManifestPath
    $before = Get-TreeSnapshot $HarnessSkills
    $runId = [guid]::NewGuid().ToString('N')
    $manifest = [ordered]@{
        Version = 1; RunId = $runId; Status = 'Planned'; Completed = $false; Error = $null; ErrorDetail = $null
        ConflictPolicy = $ConflictPolicy; HarnessSkills = $HarnessSkills; DestinationExisted = $destinationExisted; BackupPath = $backupPath
        Sources = $usableSources; Conflicts = @($conflicts); Before = @($before); After = @(); Operations = @()
        Created = (Get-Date).ToString('o')
    }
    if ($PSCmdlet.ShouldProcess($backupPath, 'Create pre-change SSOT backup')) {
        New-Item -ItemType Directory -Path (Split-Path $backupPath -Parent) -Force | Out-Null
        if ($destinationExisted) {
            Copy-Item -LiteralPath $HarnessSkills -Destination $backupPath -Recurse -Force
        } else {
            New-Item -ItemType Directory -Path $backupPath -Force | Out-Null
        }
    } else { return }
    if (-not (Write-JsonArtifact $manifestPath $manifest)) { return }

    try {
        if (-not $destinationExisted -and -not (Test-Path -LiteralPath $HarnessSkills -PathType Container)) {
            if ($PSCmdlet.ShouldProcess($HarnessSkills, 'Create SSOT directory')) { New-Item -ItemType Directory -Path $HarnessSkills -Force | Out-Null }
        }
        $appliedCount = 0
        foreach ($src in $usableSources) {
            $current = Get-TreeSnapshot $src
            $expected = @($plans | Where-Object { $_.Source -eq $src -and -not $_.BinarySeed } | ForEach-Object { $_.Entry })
            if ((($current | ConvertTo-Json -Depth 8 -Compress) -ne ($expected | Sort-Object RelativePath | ConvertTo-Json -Depth 8 -Compress))) { throw "TOCTOU source change detected: $src" }
        }
        foreach ($plan in $plans) {
            if ($plan.Conflict -and $ConflictPolicy -eq 'PreferDestination') { continue }
            if ($TestFailAfter -ge 0 -and $appliedCount -ge $TestFailAfter) { throw "Injected test failure after $appliedCount applied operation(s)." }
            $entry = $plan.Entry
            $from = if ($plan.BinarySeed) { [string]$plan.Source } else { Join-Path $plan.Source $entry.RelativePath }
            $to = Join-Path $HarnessSkills $entry.RelativePath
            if (-not (Test-NoReparseAncestors $UserHome $from) -or -not (Test-NoReparseAncestors $UserHome $to)) { throw "Unsafe reparse path: $($entry.RelativePath)" }
            if ($plan.BinarySeed) {
                if (-not (Test-Path -LiteralPath $from -PathType Leaf)) { throw "TOCTOU binary source missing: $from" }
                $sourceItem = Get-Item -LiteralPath $from -Force
                $sourceHash = (Get-FileHash -LiteralPath $from -Algorithm SHA256).Hash
                if (($sourceItem.Attributes -band [System.IO.FileAttributes]::ReparsePoint) -or $sourceHash -ne $entry.Hash) { throw "TOCTOU binary source changed: $from" }
            }
            if ($entry.Type -eq 'directory') {
                if ($PSCmdlet.ShouldProcess($to, "Create directory from $from")) { New-Item -ItemType Directory -Path $to -Force | Out-Null }
            } elseif ($entry.Type -eq 'file') {
                if ($PSCmdlet.ShouldProcess($to, "Copy from $from")) {
                    $parent = Split-Path $to -Parent; if (-not (Test-Path -LiteralPath $parent)) { New-Item -ItemType Directory -Path $parent -Force | Out-Null }
                    Copy-Item -LiteralPath $from -Destination $to -Force
                }
            }
            $appliedCount++
        }
        $manifest.After = @(Get-TreeSnapshot $HarnessSkills)
        $afterByPath = @{}
        foreach ($afterEntry in (Get-TreeSnapshot $HarnessSkills)) { $afterByPath[$afterEntry.RelativePath.ToLowerInvariant()] = $afterEntry }
        $manifest.Operations = @($plans | ForEach-Object {
            $afterEntry = if ($afterByPath.ContainsKey($_.Entry.RelativePath.ToLowerInvariant())) { $afterByPath[$_.Entry.RelativePath.ToLowerInvariant()] } else { $null }
            [pscustomobject]@{
                Source = $_.Source; RelativePath = $_.Entry.RelativePath; Conflict = $_.Conflict
                Decision = $_.Decision; Applied = -not ($_.Conflict -and $ConflictPolicy -eq 'PreferDestination')
                BeforeType = if ($null -ne $_.Before) { $_.Before.Type } else { $null }
                BeforeHash = if ($null -ne $_.Before) { $_.Before.Hash } else { $null }
                BeforeSource = $_.BeforeSource
                AfterType = if ($null -ne $afterEntry) { $afterEntry.Type } else { $null }
                AfterHash = if ($null -ne $afterEntry) { $afterEntry.Hash } else { $null }
                AfterSource = if ($null -ne $afterEntry -and $sourceByPath.ContainsKey($_.Entry.RelativePath.ToLowerInvariant())) { $sourceByPath[$_.Entry.RelativePath.ToLowerInvariant()] } else { $null }
            }
        })
        $manifest.Completed = $true
        $manifest.Status = 'Applied'
        Write-JsonArtifact $manifestPath $manifest | Out-Null
        Write-Success "Consolidation complete. Manifest: $manifestPath"
    } catch {
        $script:OperationFailed = $true
        $manifest.Status = 'Partial'
        $manifest.Error = 'ConsolidationFailed'
        $manifest.ErrorDetail = $_.Exception.Message
        $manifest.After = @(Get-TreeSnapshot $HarnessSkills)
        Write-JsonArtifact $manifestPath $manifest | Out-Null
        Write-ErrMsg "Consolidation failed: $($_.Exception.Message)"
    }
}

# -----------------------------------------------------------------------------
# 2. Link Skills (Create Junctions from agent skills/ to .harness/skills)
# -----------------------------------------------------------------------------
function Invoke-LinkSkills {
    Write-Info "Linking agent skills/ directories to SSOT: $HarnessSkills (Mode: $Mode)"

    if (-not (Test-Path $HarnessSkills)) {
        Write-ErrMsg "SSOT directory $HarnessSkills does not exist. Run with -Consolidate first."
        return
    }

    $timestamp = (Get-Date -Format 'yyyyMMdd-HHmmss')

    foreach ($target in $AgentSkillTargets) {
        $targetPath = $target.Path
        $parentDir = Split-Path $targetPath -Parent

        if (-not (Test-Path $parentDir)) {
            Write-Info "Parent directory $parentDir does not exist for $($target.Name); skipping."
            continue
        }

        # Check if already a valid junction
        if (Test-IsJunction $targetPath) {
            $existingTarget = Get-JunctionTarget $targetPath
            if ($existingTarget -eq $HarnessSkills -or $existingTarget -eq "$HarnessSkills\") {
                Write-Success "$($target.Name) ($targetPath) is ALREADY linked to SSOT."
                continue
            } else {
                Write-WarnMsg "$($target.Name) junction points to unexpected target: $existingTarget. Re-linking..."
                if ($PSCmdlet.ShouldProcess($targetPath, "Remove outdated junction")) {
                    (Get-Item $targetPath).Delete()
                }
            }
        } elseif (Test-Path $targetPath) {
            # Target is a standard directory. Backup before replacing.
            $backupPath = "$targetPath.bak-$timestamp"
            Write-WarnMsg "$($target.Name) ($targetPath) is a regular directory. Backing up to $backupPath ..."
            if ($PSCmdlet.ShouldProcess($targetPath, "Move to backup $backupPath")) {
                Move-Item -Path $targetPath -Destination $backupPath -Force
            }
        }

        if ($Mode -eq 'Junction') {
            try {
                if ($PSCmdlet.ShouldProcess($targetPath, "Create Junction -> $HarnessSkills")) {
                    New-Item -ItemType Junction -Path $targetPath -Target $HarnessSkills -Force | Out-Null
                    Write-Success "$($target.Name) junction created successfully -> $HarnessSkills"
                }
            } catch {
                Write-ErrMsg "Failed to create junction for $($target.Name): $($_.Exception.Message)"
                Write-WarnMsg "Falling back to Mode: Copy for $($target.Name)..."
                if ($PSCmdlet.ShouldProcess($targetPath, "Copy from $HarnessSkills (Fallback)")) {
                    Copy-Item -Path $HarnessSkills -Destination $targetPath -Recurse -Force
                    Write-Success "$($target.Name) copied successfully via fallback."
                }
            }
        } else {
            # Mode Copy
            if ($PSCmdlet.ShouldProcess($targetPath, "Copy from $HarnessSkills")) {
                Copy-Item -Path $HarnessSkills -Destination $targetPath -Recurse -Force
                Write-Success "$($target.Name) copied successfully."
            }
        }
    }
}

# -----------------------------------------------------------------------------
# 3. Health Check & Validation
# -----------------------------------------------------------------------------
function Invoke-HealthCheck {
    $results = [System.Collections.Generic.List[object]]::new()

    function Add-HealthResult([string]$check, [string]$target, [bool]$passed, [string]$detail) {
        $results.Add([pscustomobject]@{
            Check  = $check
            Target = $target
            Passed = $passed
            Detail = $detail
        })
    }

    Write-Host "=================================================================" -ForegroundColor Cyan
    Write-Host "  Harness Health Check & Verification" -ForegroundColor Cyan
    Write-Host "=================================================================" -ForegroundColor Cyan

    Write-Info "Checking SSOT: $HarnessSkills"
    if (Test-Path $HarnessSkills) {
        $skills = @(Get-ChildItem $HarnessSkills -Directory)
        Write-Success "SSOT exists. Total skills: $($skills.Count)"
        Add-HealthResult 'ssot-directory' $HarnessSkills $true "skills=$($skills.Count)"
    } else {
        Write-ErrMsg "SSOT $HarnessSkills DOES NOT EXIST!"
        Add-HealthResult 'ssot-directory' $HarnessSkills $false 'missing'
    }

    # Check llm-mem.exe in SSOT
    $ssotBin = Join-Path $HarnessSkills 'llm-memory\llm-mem.exe'
    if (Test-Path $ssotBin) {
        $item = Get-Item $ssotBin
        Write-Success "SSOT llm-mem.exe exists (Size: $($item.Length) bytes, LastWrite: $($item.LastWriteTime))"
        Add-HealthResult 'ssot-llm-memory' $ssotBin $true "size=$($item.Length)"
    } else {
        Write-ErrMsg "SSOT llm-mem.exe is MISSING in $ssotBin"
        Add-HealthResult 'ssot-llm-memory' $ssotBin $false 'missing'
    }

    Write-Host "`nTarget Harness Status:" -ForegroundColor Cyan
    foreach ($target in $AgentSkillTargets) {
        $path = $target.Path
        $name = $target.Name.PadRight(20)

        if (-not (Test-Path $path)) {
            Write-ErrMsg "$name : Path not found ($path)"
            Add-HealthResult 'target-link' $path $false 'missing'
            Add-HealthResult 'target-llm-memory' (Join-Path $path 'llm-memory\llm-mem.exe') $false 'target-missing'
            continue
        }

        $isJunc = Test-IsJunction $path
        if ($isJunc) {
            $juncTarget = Get-JunctionTarget $path
            $normalizedTarget = ConvertTo-NormalizedPath $juncTarget
            $normalizedHarness = ConvertTo-NormalizedPath $HarnessSkills
            $match = [string]::Equals($normalizedTarget, $normalizedHarness, [System.StringComparison]::OrdinalIgnoreCase)
            if ($match) {
                Write-Success "$name : [JUNCTION -> SSOT] OK"
                Add-HealthResult 'target-link' $path $true "junction=$normalizedTarget"
            } else {
                Write-WarnMsg "$name : [JUNCTION -> UNKNOWN ($juncTarget)]"
                Add-HealthResult 'target-link' $path $false "unexpected-junction=$normalizedTarget"
            }
        } else {
            Write-WarnMsg "$name : [REGULAR DIRECTORY] (Not a junction)"
            Add-HealthResult 'target-link' $path $false 'not-a-junction'
        }

        # Check access to llm-mem.exe via this target
        $targetLlmMem = Join-Path $path 'llm-memory\llm-mem.exe'
        if (Test-Path $targetLlmMem) {
            Write-Success "  └─ llm-mem.exe accessible directly via $name"
            Add-HealthResult 'target-llm-memory' $targetLlmMem $true 'accessible'
        } else {
            Write-ErrMsg  "  └─ llm-mem.exe NOT ACCESSIBLE via $name"
            Add-HealthResult 'target-llm-memory' $targetLlmMem $false 'missing-or-inaccessible'
        }
    }
    Write-Host "=================================================================" -ForegroundColor Cyan

    $failureCount = @($results | Where-Object { -not $_.Passed }).Count
    return [pscustomobject]@{
        Results      = $results
        FailureCount = $failureCount
        Healthy      = ($failureCount -eq 0)
    }
}

function Invoke-RestoreManifest([string]$path) {
    if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "Restore manifest not found: $path" }
    $manifest = Get-Content -Raw -LiteralPath $path | ConvertFrom-Json
    if ($manifest.Version -ne 1 -or [string]::IsNullOrWhiteSpace([string]$manifest.BackupPath)) { throw "Unsupported or incomplete restore manifest: $path" }
    $target = ConvertTo-NormalizedPath ([string]$manifest.HarnessSkills)
    $backup = ConvertTo-NormalizedPath ([string]$manifest.BackupPath)
    if (-not [string]::Equals($target, (ConvertTo-NormalizedPath $HarnessSkills), [System.StringComparison]::OrdinalIgnoreCase)) { throw "Manifest target does not match current SSOT: $target" }
    if (-not (Test-ContainedPath $HarnessRoot $backup) -or -not (Test-NoReparseAncestors $HarnessRoot $backup)) { throw "Unsafe backup path in manifest: $backup" }
    if (-not (Test-Path -LiteralPath $backup -PathType Container)) { throw "Manifest backup is missing: $backup" }

    $current = Get-TreeSnapshot $HarnessSkills
    $expected = @($manifest.After | ForEach-Object { [pscustomobject]@{ RelativePath = $_.RelativePath; Type = $_.Type; Hash = $_.Hash; Target = $_.Target } })
    $currentJson = @($current | ConvertTo-Json -Depth 8 -Compress)
    $expectedJson = @($expected | Sort-Object RelativePath | ConvertTo-Json -Depth 8 -Compress)
    if (-not $Force -and $currentJson -ne $expectedJson) { throw 'SSOT changed after the recorded synchronization; use -Force to restore.' }
    if (-not $PSCmdlet.ShouldProcess($HarnessSkills, "Restore from $backup")) { return }
    if (-not (Test-NoReparseAncestors $HarnessRoot $HarnessSkills)) { throw "Unsafe SSOT path: $HarnessSkills" }
    if ($manifest.PSObject.Properties.Name -contains 'DestinationExisted' -and -not [bool]$manifest.DestinationExisted) {
        if (Test-Path -LiteralPath $HarnessSkills) { Remove-Item -LiteralPath $HarnessSkills -Recurse -Force }
        Write-Success "Restored absent SSOT state from manifest: $path"
        return
    }
    $staging = Join-Path $HarnessRoot ('.restore-' + [guid]::NewGuid().ToString('N'))
    try {
        New-Item -ItemType Directory -Path $staging -Force | Out-Null
        foreach ($item in @(Get-ChildItem -LiteralPath $backup -Force)) {
            Copy-Item -LiteralPath $item.FullName -Destination (Join-Path $staging $item.Name) -Recurse -Force
        }
        if (Test-Path -LiteralPath $HarnessSkills) { Remove-Item -LiteralPath $HarnessSkills -Recurse -Force }
        Move-Item -LiteralPath $staging -Destination $HarnessSkills -Force
        Write-Success "Restored SSOT from manifest backup: $path"
    } catch {
        $script:OperationFailed = $true
        if (Test-Path -LiteralPath $staging) { Remove-Item -LiteralPath $staging -Recurse -Force -ErrorAction SilentlyContinue }
        throw
    }
}

# -----------------------------------------------------------------------------
# Main Execution Dispatch
# -----------------------------------------------------------------------------
$healthResult = $null
if (-not [string]::IsNullOrWhiteSpace($RestoreManifest)) {
    try { Invoke-RestoreManifest $RestoreManifest } catch { Write-ErrMsg $_.Exception.Message; $script:OperationFailed = $true }
} elseif (-not $Consolidate -and -not $LinkSkills -and -not $Check) {
    # Default: Run full pipeline (Consolidate -> LinkSkills -> Check)
    Invoke-SkillConsolidation
    if (-not $script:OperationFailed) { Invoke-LinkSkills }
    if (-not $script:OperationFailed) { $healthResult = Invoke-HealthCheck }
} else {
    if ($Consolidate) { Invoke-SkillConsolidation }
    if ($LinkSkills -and -not $script:OperationFailed) { Invoke-LinkSkills }
    if ($Check -and -not $script:OperationFailed) { $healthResult = Invoke-HealthCheck }
}

if ($null -ne $healthResult) {
    if ($PassThru) { $healthResult }
    if (-not $healthResult.Healthy) {
        Write-ErrMsg "Health check failed: $($healthResult.FailureCount) check(s) failed."
        exit 1
    }
}
if ($script:OperationFailed) { exit 1 }
