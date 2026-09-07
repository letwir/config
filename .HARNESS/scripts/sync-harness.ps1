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
    [switch]$Force
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$UserHome = $env:USERPROFILE
$HarnessRoot = Join-Path $UserHome '.harness'
$HarnessSkills = Join-Path $HarnessRoot 'skills'
$HarnessRules = Join-Path $HarnessRoot 'rules'

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

# -----------------------------------------------------------------------------
# 1. Consolidate Skills into .harness/skills (SSOT)
# -----------------------------------------------------------------------------
function Invoke-SkillConsolidation {
    Write-Info "Starting skill consolidation into SSOT: $HarnessSkills"

    if (-not (Test-Path $HarnessSkills)) {
        if ($PSCmdlet.ShouldProcess($HarnessSkills, "Create directory")) {
            New-Item -ItemType Directory -Path $HarnessSkills -Force | Out-Null
        }
    }

    # Consolidation priority order: agents (base) -> codex -> gemini-skills -> ltw-plugin (highest priority)
    $sourceDirs = @(
        (Join-Path $UserHome '.agents\skills'),
        (Join-Path $UserHome '.codex\skills'),
        (Join-Path $UserHome '.gemini\skills'),
        (Join-Path $UserHome '.gemini\config\plugins\ltw-skills-plugin\skills')
    )

    foreach ($src in $sourceDirs) {
        if (-not (Test-Path $src)) { continue }
        if (Test-IsJunction $src) {
            Write-Info "Skipping $src because it is already a junction."
            continue
        }

        Write-Info "Merging skills from $src ..."
        Get-ChildItem -Path $src -Directory | ForEach-Object {
            $destDir = Join-Path $HarnessSkills $_.Name
            if ($PSCmdlet.ShouldProcess($destDir, "Copy from $($_.FullName)")) {
                Copy-Item -Path $_.FullName -Destination $HarnessSkills -Recurse -Force
            }
        }
    }

    # Seed latest llm-mem.exe
    $candidateBins = @(
        'a:\Users\letwir\repo\llm-memory\llm-mem.exe',
        (Join-Path $UserHome '.gemini\config\plugins\ltw-skills-plugin\skills\llm-memory\llm-mem.exe'),
        (Join-Path $UserHome '.codex\skills\llm-memory\llm-mem.exe')
    )

    $seeded = $false
    $harnessLlmMemDir = Join-Path $HarnessSkills 'llm-memory'
    if (-not (Test-Path $harnessLlmMemDir)) {
        New-Item -ItemType Directory -Path $harnessLlmMemDir -Force | Out-Null
    }

    $targetBin = Join-Path $harnessLlmMemDir 'llm-mem.exe'
    foreach ($cand in $candidateBins) {
        if (Test-Path $cand) {
            Write-Info "Seeding llm-mem.exe from $cand to $targetBin"
            if ($PSCmdlet.ShouldProcess($targetBin, "Seed from $cand")) {
                Copy-Item -Path $cand -Destination $targetBin -Force
            }
            $seeded = $true
            break
        }
    }

    if (-not $seeded) {
        Write-WarnMsg "No compiled llm-mem.exe found to seed into $targetBin. Please build it first."
    } else {
        Write-Success "Seeded llm-mem.exe into SSOT."
    }

    $skillCount = (Get-ChildItem $HarnessSkills -Directory).Count
    Write-Success "Consolidation complete. Total skills in SSOT ($HarnessSkills): $skillCount"
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
    Write-Host "=================================================================" -ForegroundColor Cyan
    Write-Host "  Harness Health Check & Verification" -ForegroundColor Cyan
    Write-Host "=================================================================" -ForegroundColor Cyan

    Write-Info "Checking SSOT: $HarnessSkills"
    if (Test-Path $HarnessSkills) {
        $skills = Get-ChildItem $HarnessSkills -Directory
        Write-Success "SSOT exists. Total skills: $($skills.Count)"
    } else {
        Write-ErrMsg "SSOT $HarnessSkills DOES NOT EXIST!"
    }

    # Check llm-mem.exe in SSOT
    $ssotBin = Join-Path $HarnessSkills 'llm-memory\llm-mem.exe'
    if (Test-Path $ssotBin) {
        $item = Get-Item $ssotBin
        Write-Success "SSOT llm-mem.exe exists (Size: $($item.Length) bytes, LastWrite: $($item.LastWriteTime))"
    } else {
        Write-ErrMsg "SSOT llm-mem.exe is MISSING in $ssotBin"
    }

    Write-Host "`nTarget Harness Status:" -ForegroundColor Cyan
    foreach ($target in $AgentSkillTargets) {
        $path = $target.Path
        $name = $target.Name.PadRight(20)

        if (-not (Test-Path $path)) {
            Write-ErrMsg "$name : Path not found ($path)"
            continue
        }

        $isJunc = Test-IsJunction $path
        if ($isJunc) {
            $juncTarget = Get-JunctionTarget $path
            $match = ($juncTarget -eq $HarnessSkills -or $juncTarget -eq "$HarnessSkills\")
            if ($match) {
                Write-Success "$name : [JUNCTION -> SSOT] OK"
            } else {
                Write-WarnMsg "$name : [JUNCTION -> UNKNOWN ($juncTarget)]"
            }
        } else {
            Write-WarnMsg "$name : [REGULAR DIRECTORY] (Not a junction)"
        }

        # Check access to llm-mem.exe via this target
        $targetLlmMem = Join-Path $path 'llm-memory\llm-mem.exe'
        if (Test-Path $targetLlmMem) {
            Write-Success "  └─ llm-mem.exe accessible directly via $name"
        } else {
            Write-ErrMsg  "  └─ llm-mem.exe NOT ACCESSIBLE via $name"
        }
    }
    Write-Host "=================================================================" -ForegroundColor Cyan
}

# -----------------------------------------------------------------------------
# Main Execution Dispatch
# -----------------------------------------------------------------------------
if (-not $Consolidate -and -not $LinkSkills -and -not $Check) {
    # Default: Run full pipeline (Consolidate -> LinkSkills -> Check)
    Invoke-SkillConsolidation
    Invoke-LinkSkills
    Invoke-HealthCheck
} else {
    if ($Consolidate) { Invoke-SkillConsolidation }
    if ($LinkSkills)  { Invoke-LinkSkills }
    if ($Check)       { Invoke-HealthCheck }
}
