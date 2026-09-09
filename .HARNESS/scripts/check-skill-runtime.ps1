[CmdletBinding()]
param(
  [string]$ManifestPath = (Join-Path $PSScriptRoot '..\rules\skill-runtime.json'),
  [string]$OutputDirectory = (Join-Path $PSScriptRoot '..\evaluation\runs'),
  [switch]$Json,
  [switch]$NoReceipt
)
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$script:MaxOutput = 4096

function Expand-Token([string]$p,[string]$root) {
  if ([string]::IsNullOrWhiteSpace($p)) { return $p }
  $ps = (Get-Command pwsh.exe -ErrorAction SilentlyContinue).Source
  if (!$ps) { $ps = (Get-Command powershell.exe -ErrorAction SilentlyContinue).Source }
  return $p.Replace('{HarnessRoot}', $root).Replace('{UserHome}', $env:USERPROFILE).Replace('{PowerShell}', $ps)
}
function Resolve-FinalPath([string]$path) {
  $full=[IO.Path]::GetFullPath($path).TrimEnd('\'); if (!(Test-Path -LiteralPath $full)) { return $full }
  $volume=[IO.Path]::GetPathRoot($full).TrimEnd('\'); $current=$volume
  foreach ($part in $full.Substring($volume.Length).TrimStart('\').Split('\')) {
    $current=Join-Path $current $part; $item=Get-Item -Force -LiteralPath $current
    if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { $target=$item.ResolveLinkTarget($true); if ($null -eq $target) { throw "unresolved reparse point: $current" }; $current=$target.FullName }
  }
  [IO.Path]::GetFullPath($current).TrimEnd('\')
}
function Contained([string]$candidate,[object[]]$roots) {
  try {$c = Resolve-FinalPath $candidate} catch {return $false}
  foreach ($root in $roots) { try {$r=Resolve-FinalPath $root} catch {continue}; if ($c.Equals($r,[StringComparison]::OrdinalIgnoreCase) -or $c.StartsWith($r+'\',[StringComparison]::OrdinalIgnoreCase)) { return $true } }
  return $false
}
function Redact([string]$s) {
  if ($null -eq $s) { return '' }
  $s = [regex]::Replace($s,'(?i)(api[_-]?key|token|password|secret)\s*[=:]\s*[^\s,;]+','$1=<REDACTED>')
  return [regex]::Replace($s,'(?i)bearer\s+[A-Za-z0-9._~-]+','Bearer <REDACTED>')
}
function Normalize([string]$s,[int]$limit) {
  $r = Redact $s
  $r = [regex]::Replace($r,'\r\n?','\n').Trim()
  $truncated = $r.Length -gt $limit
  if ($truncated) { $r = $r.Substring(0,$limit) }
  [pscustomobject]@{ text=$r; truncated=$truncated; sha256=([Convert]::ToHexString(([Security.Cryptography.SHA256]::Create()).ComputeHash([Text.Encoding]::UTF8.GetBytes($r)))) }
}
function Invoke-Bounded([string]$exe,[string[]]$argumentList,[int]$timeout,[string]$root) {
  $out=[IO.Path]::GetTempFileName(); $err=[IO.Path]::GetTempFileName(); $p=$null; $timed=$false
  try {
    $si=[Diagnostics.ProcessStartInfo]::new(); $si.FileName=$exe; $si.UseShellExecute=$false; $si.RedirectStandardOutput=$true; $si.RedirectStandardError=$true; $si.CreateNoWindow=$true; $si.WorkingDirectory=$root; $si.Environment['CI']='true'; $si.Environment['NO_COLOR']='1'
    foreach($a in $argumentList){$null=$si.ArgumentList.Add([string]$a)}
    $p=[Diagnostics.Process]::new(); $p.StartInfo=$si; [void]$p.Start()
    $outTask=$p.StandardOutput.ReadToEndAsync(); $errTask=$p.StandardError.ReadToEndAsync()
    if (!$p.WaitForExit($timeout)) { $timed=$true; try { & taskkill.exe /PID $p.Id /T /F 2>$null | Out-Null } catch {}; try{$p.Kill($true)}catch{}; $p.WaitForExit() }
    $rawOut=$outTask.Result; $rawErr=$errTask.Result; $exit=if($timed){-1}else{$p.ExitCode}
    [pscustomobject]@{ exitCode=$exit; timedOut=$timed; raw=($rawOut+$(if($rawErr){"`n$rawErr"}else{''})); pid=$p.Id }
  } finally { if($p){$p.Dispose()}; Remove-Item -LiteralPath $out,$err -Force -ErrorAction SilentlyContinue }
}
function Hash-File([string]$p) { if(!(Test-Path -LiteralPath $p -PathType Leaf)){return $null}; (Get-FileHash -LiteralPath $p -Algorithm SHA256).Hash.ToUpperInvariant() }
function Invoke-ClientBounded([string]$wrapper,[string]$target,[object[]]$clientSpecArgs,[int]$timeout,[string]$root) {
  $out=[IO.Path]::GetTempFileName();$err=[IO.Path]::GetTempFileName()
  try {$probe=Join-Path $PSScriptRoot 'invoke-client-catalog-probe.ps1';$payload=@($clientSpecArgs)|ConvertTo-Json -Compress;$encoded=[Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes($payload));$probeArgs=@('-NoProfile','-File',$probe,'-Target',$target,'-ArgsBase64',$encoded,'-OutPath',$out,'-ErrPath',$err);$p=Invoke-Bounded $wrapper $probeArgs $timeout $root;$raw=(Get-Content -Raw -LiteralPath $out -ErrorAction SilentlyContinue)+"`n"+(Get-Content -Raw -LiteralPath $err -ErrorAction SilentlyContinue);[pscustomobject]@{exitCode=$p.exitCode;timedOut=$p.timedOut;raw=$raw;pid=$p.pid}}
  finally {Remove-Item -LiteralPath $out,$err -Force -ErrorAction SilentlyContinue}
}
function Run-Command([string]$name,[string]$path,[object]$spec,[string]$root,[object[]]$approved,[bool]$required) {
  $result=[ordered]@{name=$name; required=$required; path=$path; finalPath=$null; classification='Unknown'; executableHealthy='Missing'; catalogObserved='NotObserved'; representativeOperation='UnTested'; commandIdentity=[ordered]@{executable=$path;args=@()}; fileSha256=$null; version=$null; output=$null; error=$null}
  if(!(Contained $path $approved)){ $result.error='path-outside-approved-roots'; $result.classification='PathRejected'; return [pscustomobject]$result }
  if(!(Test-Path -LiteralPath $path -PathType Leaf)){ $result.classification=if($required){'Missing'}else{'NotConfigured'}; return [pscustomobject]$result }
  $result.executableHealthy='Healthy'; $result.finalPath=Resolve-FinalPath $path; $result.fileSha256=Hash-File $path; $result.version=(Get-Item -LiteralPath $path).VersionInfo.FileVersion
  if($spec.PSObject.Properties.Name -contains 'sha256' -and $result.fileSha256 -ne ([string]$spec.sha256).ToUpperInvariant()){ $result.classification='HashMismatch';$result.error='file-hash-mismatch';return [pscustomobject]$result }
  $help=@($spec.helpArgs); $result.commandIdentity.args=$help
  $h=Invoke-Bounded $path $help ([int]$spec.timeoutMs) $root
  $norm=Normalize $h.raw $script:MaxOutput; $result.output=[ordered]@{exitCode=$h.exitCode;timedOut=$h.timedOut;text=$norm.text;truncated=$norm.truncated;sha256=$norm.sha256;checkedBeforeNormalization=$true}
  if($h.timedOut){$result.catalogObserved='Timeout';$result.classification='Timeout';$result.error='help-timeout';return [pscustomobject]$result}
  $expectedHelpExit=0;if($spec.PSObject.Properties.Name -contains 'requiredExitCode'){$expectedHelpExit=[int]$spec.requiredExitCode};$result.catalogObserved='Observed'; if($h.exitCode -ne $expectedHelpExit){$result.classification='CommandNonZero';$result.error="help-exit-$($h.exitCode)"} else {$result.classification='Healthy'}
  if($result.classification -eq 'Healthy' -and $spec.PSObject.Properties.Name -contains 'requiredText' -and -not $h.raw.Contains([string]$spec.requiredText)){ $result.classification='ContractMismatch';$result.error='required-output-missing';return [pscustomobject]$result }
  if($null -ne $spec.representative){ $ra=@($spec.representative.args|ForEach-Object{Expand-Token ([string]$_) $root}); $result.commandIdentity.representativeArgs=$ra; $rr=Invoke-Bounded $path $ra ([int]$spec.timeoutMs) $root; $rn=Normalize $rr.raw $script:MaxOutput; $tokenOk=$true;if($spec.representative.PSObject.Properties.Name -contains 'requiredText'){$tokenOk=$rr.raw.Contains([string]$spec.representative.requiredText)}; if($rr.timedOut){$result.representativeOperation='Timeout';$result.classification='RepresentativeFailed'}elseif($rr.exitCode -eq [int]$spec.representative.requiredExitCode -and $tokenOk){$result.representativeOperation='Passed'}else{$result.representativeOperation='Failed';$result.classification='RepresentativeFailed';$result.error='representative-contract-failed'}; $result.output.representative=[ordered]@{exitCode=$rr.exitCode;timedOut=$rr.timedOut;text=$rn.text;truncated=$rn.truncated;sha256=$rn.sha256;checkedBeforeNormalization=$true} }
  [pscustomobject]$result
}
function Run-Client([object]$spec,[string]$root,[object[]]$approved) {
  $name=[string]$spec.name; $target=Expand-Token ([string]$spec.target) $root; $wrapper=Expand-Token ([string]$spec.wrapper) $root
  $base=[ordered]@{name=$name;required=[bool]$spec.required;path=$target;finalPath=$null;classification='Unknown';executableHealthy='Missing';catalogObserved='NotObserved';representativeOperation='UnTested';commandIdentity=[ordered]@{wrapper=$wrapper;target=$target;args=@($spec.args);wrapperSha256=$null;targetSha256=$null};fileSha256=$null;version=$null;output=$null;error=$null}
  if(!(Test-Path -LiteralPath $target -PathType Leaf)){ $base.classification=if($spec.required){'Missing'}else{'NotConfigured'}; return [pscustomobject]$base }
  if([string]::IsNullOrWhiteSpace($wrapper) -or !(Test-Path -LiteralPath $wrapper -PathType Leaf)){ $base.classification='WrapperMismatch';$base.error='wrapper-missing';return [pscustomobject]$base }
  if(!(Contained $target $approved)){ $base.classification='PathRejected';$base.error='target-outside-approved-roots';return [pscustomobject]$base }
  if(!(Contained $wrapper $approved)){ $base.classification='WrapperMismatch';$base.error='wrapper-outside-approved-roots';return [pscustomobject]$base }
  $base.executableHealthy='Healthy';$base.finalPath=Resolve-FinalPath $target;$base.fileSha256=Hash-File $target;$base.version=(Get-Item -LiteralPath $target).VersionInfo.FileVersion;$base.commandIdentity.wrapperSha256=Hash-File $wrapper;$base.commandIdentity.targetSha256=$base.fileSha256
  $r=Invoke-ClientBounded $wrapper $target @($spec.args) ([int]$spec.timeoutMs) $root;$attempts=1
  $retryCount=if($spec.PSObject.Properties.Name -contains 'retryCount'){[int]$spec.retryCount}else{0}
  while($attempts -le $retryCount -and ($r.timedOut -or $r.exitCode -ne 0)){$r=Invoke-ClientBounded $wrapper $target @($spec.args) ([int]$spec.timeoutMs) $root;$attempts++}
  $n=Normalize $r.raw $script:MaxOutput;$base.output=[ordered]@{exitCode=$r.exitCode;timedOut=$r.timedOut;attempts=$attempts;text=$n.text;truncated=$n.truncated;sha256=$n.sha256;checkedBeforeNormalization=$true}
  if($r.timedOut){$base.classification='Timeout';$base.catalogObserved='Timeout';$base.error='catalog-timeout';return [pscustomobject]$base}
  if($r.exitCode -ne 0){$base.classification='CommandNonZero';$base.catalogObserved='Unverified';$base.error="catalog-exit-$($r.exitCode)";return [pscustomobject]$base}
  if($spec.PSObject.Properties.Name -contains 'requiredText' -and -not $r.raw.Contains([string]$spec.requiredText)){ $base.classification='ContractMismatch';$base.error='required-output-missing';return [pscustomobject]$base }
  $base.catalogObserved=if($spec.mode -eq 'catalog-api-or-unverified'){'Unverified'}else{'Observed'};$base.classification=if($base.catalogObserved -eq 'Observed'){'Healthy'}else{'Unverified'}
  [pscustomobject]$base
}
function Write-Atomic([string]$path,$value){$tmp="$path.$([guid]::NewGuid().ToString('N')).tmp";try{[IO.File]::WriteAllText($tmp,($value|ConvertTo-Json -Depth 20),[Text.UTF8Encoding]::new($false));Move-Item -LiteralPath $tmp -Destination $path -Force}finally{Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue}}
try {
  $manifestFull=[IO.Path]::GetFullPath($ManifestPath);$root=[IO.Path]::GetFullPath((Join-Path (Split-Path -Parent $manifestFull) '..'));$m=Get-Content -Raw -LiteralPath $manifestFull|ConvertFrom-Json;if($m.schemaVersion -ne 1){throw 'unsupported manifest schemaVersion'};$script:MaxOutput=[int]$m.maxOutputChars
  $approved=@($m.approvedRoots|ForEach-Object{Expand-Token ([string]$_) $root});$clients=@($m.clients|ForEach-Object{Run-Client $_ $root $approved});if($m.PSObject.Properties.Name -contains 'optionalClients'){$clients+=@($m.optionalClients|ForEach-Object{Run-Client $_ $root $approved})};$tools=@($m.tools|ForEach-Object{Run-Command ([string]$_.name) (Expand-Token ([string]$_.path) $root) $_ $root $approved ([bool]$_.required)})
  $bad=@($tools+$clients|Where-Object{$_.required -and $_.classification -in @('Missing','PathRejected','Timeout','CommandNonZero','ContractMismatch','WrapperMismatch','HashMismatch','RepresentativeFailed')}).Count;$overall=if($bad){1}else{0};$receipt=[ordered]@{schemaVersion=1;run_id=[guid]::NewGuid().ToString();generated_at=(Get-Date).ToUniversalTime().ToString('o');overall_exit_code=$overall;tools=$tools;clients=$clients;latest_updated=(-not $bad)}
  if(!$NoReceipt){New-Item -ItemType Directory -Force -Path $OutputDirectory|Out-Null;$run=Join-Path $OutputDirectory "skill-runtime-$($receipt.run_id).json";Write-Atomic $run $receipt;if(!$bad){Write-Atomic (Join-Path $OutputDirectory 'skill-runtime-latest.json') $receipt}}
  if($Json){$receipt|ConvertTo-Json -Depth 20 -Compress}else{$receipt|ConvertTo-Json -Depth 20};exit $overall
}catch{if($Json){[ordered]@{schemaVersion=1;overall_exit_code=2;error=$_.Exception.Message}|ConvertTo-Json -Compress}else{Write-Error $_.Exception.Message};exit 2}
