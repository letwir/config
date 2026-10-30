$ErrorActionPreference='Stop'
$script=Join-Path (Split-Path -Parent $PSScriptRoot) 'install-pi.ps1'
$tmp=Join-Path ([IO.Path]::GetTempPath()) ('pi-install-test-'+[guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $tmp | Out-Null
function Assert($condition,[string]$message) { if (-not $condition) { throw "FAIL: $message" } }
function Entry([string]$path) { try { Get-Item -LiteralPath $path -Force -ErrorAction Stop } catch { $null } }
try {
 $agent=Join-Path $tmp 'agent path with spaces'
 & $script -AgentDir $agent -DryRun | Out-Null
 Assert (-not (Test-Path $agent)) 'dry-run wrote files'
 & $script -AgentDir $agent
 $link=Entry (Join-Path $agent 'settings.json')
 Assert ($link -and (($link.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0)) 'fresh install did not make a real symbolic link'
 & $script -AgentDir $agent
 Assert ((Entry (Join-Path $agent 'settings.json')).Target -eq (Join-Path (Split-Path -Parent $script) 'agent/settings.json')) 'rerun changed correct link'
 # Entire preflight must reject regular conflict, wrong link, and ancestor regular-file conflicts without partial writes.
 $regular=Join-Path $tmp 'regular'; New-Item -ItemType Directory $regular | Out-Null; Set-Content (Join-Path $regular 'settings.json') 'original'
 try { & $script -AgentDir $regular; throw 'regular conflict unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } }
 Assert ((Get-Content (Join-Path $regular 'settings.json') -Raw).Trim() -eq 'original' -and -not (Test-Path (Join-Path $regular 'agents/pi'))) 'regular conflict mutated destination'
 $parent=Join-Path $tmp 'parent-conflict'; New-Item -ItemType Directory $parent | Out-Null; Set-Content (Join-Path $parent 'agents') 'keep'
 try { & $script -AgentDir $parent; throw 'parent conflict unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } }
 Assert (-not (Test-Path (Join-Path $parent 'settings.json'))) 'parent conflict partially installed'
 $wrong=Join-Path $tmp 'wrong'; New-Item -ItemType Directory $wrong | Out-Null
 New-Item -ItemType SymbolicLink -Path (Join-Path $wrong 'settings.json') -Target (Join-Path $tmp 'not-source') | Out-Null
 try { & $script -AgentDir $wrong; throw 'wrong link conflict unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } }
 & $script -AgentDir $wrong -Backup
 Assert ((Entry (Join-Path $wrong 'settings.json')).Attributes -band [IO.FileAttributes]::ReparsePoint) 'wrong-link backup did not install link'
 # Separate backups for regular files and broken symlinks.
 foreach ($kind in @('backup-regular','backup-dangling')) {
  $d=Join-Path $tmp $kind; New-Item -ItemType Directory $d | Out-Null; $p=Join-Path $d 'settings.json'
  if ($kind -eq 'backup-regular') { Set-Content $p 'preserve' } else { New-Item -ItemType SymbolicLink -Path $p -Target (Join-Path $tmp 'absent-target') | Out-Null }
  & $script -AgentDir $d -Backup
  Assert ((Get-ChildItem -LiteralPath $d -Force | Where-Object Name -Like 'settings.json.backup.*').Count -eq 1) "$kind original not preserved"
 }
 # Source/destination overlap and unsafe symlinked parent are rejected.
 $dist=Split-Path -Parent $script; $hroot=Split-Path -Parent $dist
 foreach ($candidate in @($dist,$hroot,(Join-Path $dist 'child'))) { try { & $script -AgentDir $candidate; throw "overlap unexpectedly succeeded: $candidate" } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } } }
 $sourceLink=Join-Path $tmp 'source-link'; New-Item -ItemType SymbolicLink -Path $sourceLink -Target $dist | Out-Null
 try { & $script -AgentDir (Join-Path $sourceLink 'child'); throw 'symlinked parent unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } }
 # Missing source fixture: reproduce every link source, then remove just one required file.
 $fake=Join-Path $tmp 'fake/.HARNESS/.pi'; New-Item -ItemType Directory -Force (Join-Path $fake 'agent/extensions/subagent') | Out-Null
 New-Item -ItemType Directory -Force (Join-Path $fake 'agent/agents/pi') | Out-Null
 Copy-Item (Join-Path (Split-Path -Parent $script) 'agent/settings.json') (Join-Path $fake 'agent/settings.json')
 Copy-Item (Join-Path (Split-Path -Parent $script) 'agent/extensions/subagent/config.json') (Join-Path $fake 'agent/extensions/subagent/config.json')
 Copy-Item (Join-Path (Split-Path -Parent $script) 'agent/agents/pi/*.md') (Join-Path $fake 'agent/agents/pi')
 $fakeRoot=Split-Path -Parent $fake
 foreach ($f in @('AGENTS.md','POLYGLYPH.md','HARDWARE.md')) { Set-Content (Join-Path $fakeRoot $f) 'stub' }
 New-Item -ItemType Directory -Force (Join-Path $fakeRoot 'agents') | Out-Null
 foreach ($f in @('worker.md','WORKER.css','SUBAGENTS.lrf','AGENT_ROUTER.md')) { Set-Content (Join-Path (Join-Path $fakeRoot 'agents') $f) 'stub' }
 foreach ($d in @('rules','skills','crabcode/prompts','persona','etl','scripts','evaluation','memo')) { New-Item -ItemType Directory -Force (Join-Path $fakeRoot $d) | Out-Null }
 Set-Content (Join-Path $fakeRoot 'rules/BOOTSTRAP.lrf') 'stub'
 Copy-Item $script (Join-Path $fake 'install-pi.ps1')
 Remove-Item (Join-Path $fakeRoot 'rules') -Recurse -Force
 $missingAgent=Join-Path $tmp 'missing-agent'
 try { & (Join-Path $fake 'install-pi.ps1') -AgentDir $missingAgent; throw 'missing source unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw } }
 Assert (-not (Test-Path $missingAgent)) 'missing source mutated destination'
 # Inject New-Item symlink failure; the original conflicting file must be restored and reported.
 $fail=Join-Path $tmp 'link-failure'; New-Item -ItemType Directory $fail | Out-Null; Set-Content (Join-Path $fail 'settings.json') 'restore-me'
 function New-Item { param([string]$ItemType,[string]$Path,[string]$Target,[switch]$Force); if ($ItemType -eq 'SymbolicLink') { throw 'injected link failure' }; Microsoft.PowerShell.Management\New-Item @PSBoundParameters }
 try { & $script -AgentDir $fail -Backup; throw 'injected link failure unexpectedly succeeded' } catch { if ($_.Exception.Message -like '*unexpectedly succeeded*') { throw }; Assert ($_.Exception.Message -like '*rollback completed*') 'failure did not report successful rollback' }
 Remove-Item Function:\New-Item -ErrorAction SilentlyContinue
 Assert ((Get-Content (Join-Path $fail 'settings.json') -Raw).Trim() -eq 'restore-me') 'link failure rollback lost original'
 'PASS PowerShell installer scenarios'
} finally { Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue }
