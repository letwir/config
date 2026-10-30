[CmdletBinding()]
param([string]$AgentDir = (Join-Path $HOME '.pi/agent'), [switch]$DryRun, [switch]$Backup, [switch]$Help)
if ($Help) { 'Usage: .\install-pi.ps1 [-AgentDir DIR] [-DryRun] [-Backup] [-Help]'; exit 0 }
$ErrorActionPreference = 'Stop'
function Get-Entry([string]$Path) { try { Get-Item -LiteralPath $Path -Force -ErrorAction Stop } catch { if ($_.Exception -is [Management.Automation.ItemNotFoundException] -or $_.CategoryInfo.Category -eq [Management.Automation.ErrorCategory]::ObjectNotFound) { return $null }; throw } }
function Is-Link($Item) { return $null -ne $Item -and (($Item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) }
$dist = $PSScriptRoot; $root = Split-Path -Parent $dist
$items = @(
 @('agent/settings.json','settings.json'), @('agent/extensions/subagent','extensions/subagent'), @('agent/agents/pi','agents/pi'),
 @('../agents/worker.md','agents/worker.md'), @('../agents/WORKER.css','agents/WORKER.css'), @('../agents/SUBAGENTS.lrf','agents/SUBAGENTS.lrf'), @('../agents/AGENT_ROUTER.md','agents/AGENT_ROUTER.md'),
 @('../AGENTS.md','AGENTS.md'), @('../POLYGLYPH.md','POLYGLYPH.md'), @('../HARDWARE.md','HARDWARE.md'), @('../rules','rules'), @('../skills','skills'), @('../crabcode/prompts','prompts'), @('../persona','persona'), @('../etl','etl'), @('../scripts','scripts'), @('../evaluation','evaluation'), @('../memo','memo')
)
$dest = [IO.Path]::GetFullPath($AgentDir).TrimEnd('\','/')
$srcRoots = @([IO.Path]::GetFullPath($dist).TrimEnd('\','/'), [IO.Path]::GetFullPath($root).TrimEnd('\','/'))
foreach ($srcRoot in $srcRoots) { if ($dest.Equals($srcRoot,[StringComparison]::OrdinalIgnoreCase) -or $dest.StartsWith($srcRoot+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase) -or $srcRoot.StartsWith($dest+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)) { throw "Unsafe source/destination overlap: $dest <> $srcRoot" } }
$parent = $dest; while ($parent -and $parent -ne [IO.Path]::GetPathRoot($parent)) { $ent=Get-Entry $parent; if ($ent -and (Is-Link $ent)) { throw "Unsafe symlinked destination ancestor: $parent" }; if ($ent -and -not $ent.PSIsContainer) { throw "Destination ancestor is not a directory: $parent" }; $parent=Split-Path -Parent $parent }
$conflicts = [Collections.Generic.List[string]]::new()
foreach ($item in $items) {
 $src=[IO.Path]::GetFullPath((Join-Path $dist $item[0])); $dst=Join-Path $dest $item[1]
 if (-not (Test-Path -LiteralPath $src)) { throw "Missing source: $src" }
 $p=Split-Path -Parent $dst; while ($p -and $p -ne $dest) { $pe=Get-Entry $p; if (($pe -and (Is-Link $pe)) -or ($pe -and -not $pe.PSIsContainer)) { throw "Unsafe destination parent: $p" }; $p=Split-Path -Parent $p }
 $ent=Get-Entry $dst
 if ($ent -and (Is-Link $ent) -and ([IO.Path]::GetFullPath([string]$ent.Target).TrimEnd('\','/') -eq $src.TrimEnd('\','/'))) { continue }
 if ($ent) { $conflicts.Add($dst) }
}
if ($conflicts.Count -and -not $Backup) { throw "Conflicts found; no changes made: $($conflicts -join ', ')" }
if ($DryRun) { foreach ($item in $items) { "$(Join-Path $dist $item[0]) -> $(Join-Path $dest $item[1])" }; exit 0 }
$created=[Collections.Generic.List[string]]::new(); $backed=[Collections.Generic.List[object]]::new()
try {
 foreach ($item in $items) {
  $src=[IO.Path]::GetFullPath((Join-Path $dist $item[0])); $dst=Join-Path $dest $item[1]; $ent=Get-Entry $dst
  if ($ent -and (Is-Link $ent) -and ([IO.Path]::GetFullPath([string]$ent.Target).TrimEnd('\','/') -eq $src.TrimEnd('\','/'))) { "Already linked: $dst"; continue }
  New-Item -ItemType Directory -Force -Path (Split-Path -Parent $dst) | Out-Null
  if ($ent) {
   if (-not $Backup) { throw "Conflict appeared after preflight: $dst" }
   $b="$dst.backup.$([DateTime]::UtcNow.ToString('yyyyMMddHHmmssfff')).$PID"; $n=0; while (Get-Entry $b) { $n++; $b="$dst.backup.$([DateTime]::UtcNow.ToString('yyyyMMddHHmmssfff')).$PID.$n" }
   Move-Item -LiteralPath $dst -Destination $b; $backed.Add(@($b,$dst))
  }
  New-Item -ItemType SymbolicLink -Path $dst -Target $src | Out-Null
  $created.Add($dst); $made=Get-Entry $dst
  if (-not (Is-Link $made) -or [IO.Path]::GetFullPath([string]$made.Target).TrimEnd('\','/') -ne $src.TrimEnd('\','/')) { throw "Symbolic link verification failed: $dst" }
  "Linked: $dst -> $src"
 }
} catch {
 $original=$_.Exception.Message; $rbErrors=[Collections.Generic.List[string]]::new()
 for ($i=$created.Count-1; $i -ge 0; $i--) { try { Remove-Item -LiteralPath $created[$i] -Force -ErrorAction Stop } catch { $rbErrors.Add("remove $($created[$i]): $($_.Exception.Message)") } }
 for ($i=$backed.Count-1; $i -ge 0; $i--) { $b=$backed[$i][0]; $dst=$backed[$i][1]; if (Get-Entry $dst) { $rbErrors.Add("preserved backup $b; destination occupied: $dst"); continue }; try { Move-Item -LiteralPath $b -Destination $dst -ErrorAction Stop } catch { $rbErrors.Add("restore $b to $dst`: $($_.Exception.Message)") } }
 if ($rbErrors.Count) { throw "Install failed: $original; ROLLBACK ERRORS: $($rbErrors -join '; ')" }
 throw "Install failed; rollback completed: $original"
}
