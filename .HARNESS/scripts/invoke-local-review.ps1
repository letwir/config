[CmdletBinding()]
param([Parameter(Mandatory=$true)][string]$Query,[Parameter(Mandatory=$true)][string[]]$Target,[string]$RootPath='',[string]$OutputPath='',[switch]$Json)
Set-StrictMode -Version Latest; $ErrorActionPreference='Stop'
$script:MaxTargets=32; $script:MaxQueryLength=512; $script:MaxCapturedCharacters=4096
if([string]::IsNullOrWhiteSpace($RootPath)){$RootPath=Join-Path $PSScriptRoot '..'}
function Fail([string]$Code){throw "local-review:$Code"}
function Hash-Bytes([byte[]]$Value){$h=[Security.Cryptography.SHA256]::Create();try{([BitConverter]::ToString($h.ComputeHash($Value))).Replace('-','')}finally{$h.Dispose()}}
function Hash-Text([string]$Value){Hash-Bytes ([Text.Encoding]::UTF8.GetBytes($Value))}
function Is-Reparse([string]$Path){((Get-Item -LiteralPath $Path -Force).Attributes -band [IO.FileAttributes]::ReparsePoint)-ne 0}
function Get-Root([string]$Path){$full=[IO.Path]::GetFullPath($Path).TrimEnd('\');if(!(Test-Path -LiteralPath $full -PathType Container)){Fail 'invalid_root'};if(Is-Reparse $full){Fail 'reparse_root'};$full}
function Get-Targets([string]$Root,[string[]]$Names){
 $flat=@($Names|ForEach-Object{$_ -split ','}|ForEach-Object{$_.Trim()}|Where-Object{$_});if($flat.Count-lt 1-or $flat.Count-gt $script:MaxTargets){Fail 'invalid_target_count'}
 $seen=@{};$result=@();foreach($name in $flat){
  if([IO.Path]::IsPathRooted($name)-or $name-match'(^|[\\/])\.\.([\\/]|$)'-or $name-match'[*?\[\]]'){Fail 'invalid_target_path'}
  $full=[IO.Path]::GetFullPath((Join-Path $Root $name));if(!$full.StartsWith($Root+'\',[StringComparison]::OrdinalIgnoreCase)){Fail 'target_escape'}
  if(!(Test-Path -LiteralPath $full -PathType Leaf)){Fail 'target_missing_or_not_file'}
  $cursor=Split-Path -Parent $full;while($cursor.Length-ge $Root.Length){if(Is-Reparse $cursor){Fail 'reparse_target'};if($cursor-eq $Root){break};$cursor=Split-Path -Parent $cursor};if(Is-Reparse $full){Fail 'reparse_target'}
  $relative=$full.Substring($Root.Length+1).Replace('\','/');$key=$relative.ToLowerInvariant();if($seen.ContainsKey($key)){Fail 'duplicate_target'};$seen[$key]=$true
  $result+=[pscustomobject]@{path=$relative;sha256=(Get-FileHash -LiteralPath $full -Algorithm SHA256).Hash}
 };$result
}
function Quote-Arg([string]$Value){'"'+$Value.Replace('"','\"')+'"'}
function Invoke-Fixed([string]$Executable,[string[]]$CommandArgs,[string]$Root){
 $si=New-Object Diagnostics.ProcessStartInfo;$si.FileName=$Executable;$si.WorkingDirectory=$Root;$si.UseShellExecute=$false;$si.RedirectStandardOutput=$true;$si.RedirectStandardError=$true;$si.CreateNoWindow=$true;$si.Arguments=(($CommandArgs|ForEach-Object{Quote-Arg ([string]$_)})-join' ')
 $p=New-Object Diagnostics.Process;$p.StartInfo=$si;try{[void]$p.Start();$stdout=$p.StandardOutput.ReadToEndAsync();$stderr=$p.StandardError.ReadToEndAsync();$timedOut=-not $p.WaitForExit(30000);if($timedOut){try{$p.Kill()}catch{};$p.WaitForExit()};$text=[regex]::Replace(($stdout.Result+"`n"+$stderr.Result),'\r\n?','\n');if($text.Length-gt $script:MaxCapturedCharacters){$text=$text.Substring(0,$script:MaxCapturedCharacters)};[pscustomobject]@{exit_code=if($timedOut){-1}else{$p.ExitCode};output_sha256=(Hash-Text $text);bounded=$true}}finally{$p.Dispose()}
}
function Command-Record([string]$Name,$Result){[ordered]@{executable=$Name;exit_code=[int]$Result.exit_code;output_sha256=$Result.output_sha256;bounded=$true}}
function Write-Atomic([string]$Root,[string]$Path,$Value){
 $full=[IO.Path]::GetFullPath($Path);if(!$full.StartsWith($Root+'\',[StringComparison]::OrdinalIgnoreCase)){Fail 'output_escape'};$dir=Split-Path -Parent $full;if(!(Test-Path -LiteralPath $dir -PathType Container)){New-Item -ItemType Directory -Path $dir -Force|Out-Null};$cursor=$dir;while($cursor.Length-ge $Root.Length){if(Is-Reparse $cursor){Fail 'reparse_output'};if($cursor-eq $Root){break};$cursor=Split-Path -Parent $cursor};$tmp="$full.$([guid]::NewGuid().ToString('N')).tmp";try{[IO.File]::WriteAllText($tmp,($Value|ConvertTo-Json -Depth 12),[Text.UTF8Encoding]::new($false));Move-Item -LiteralPath $tmp -Destination $full -Force}finally{Remove-Item -LiteralPath $tmp -Force -ErrorAction SilentlyContinue}
}
try{
 if([string]::IsNullOrWhiteSpace($Query)-or $Query.Length-gt $script:MaxQueryLength-or $Query-match'[\r\n]'){Fail 'invalid_query'}
 $root=Get-Root $RootPath;$targets=@(Get-Targets $root $Target);$paths=@($targets.path);$rg=(Get-Command rg.exe -ErrorAction Stop).Source;$git=(Get-Command git.exe -ErrorAction Stop).Source
 $rr=Invoke-Fixed $rg (@('--fixed-strings','--line-number','--no-heading','--color','never','--max-count','20',$Query,'--')+$paths) $root;$sr=Invoke-Fixed $git (@('status','--short','--')+$paths) $root;$dr=Invoke-Fixed $git (@('diff','--check','--')+$paths) $root
 $ok=($rr.exit_code-in @(0,1))-and $sr.exit_code-eq 0-and $dr.exit_code-eq 0
 $receipt=[ordered]@{schema_version=1;receipt_id=[guid]::NewGuid().ToString();repo_root=$root;query=[ordered]@{mode='fixed_string';sha256=(Hash-Text $Query);length=$Query.Length};targets=@($targets);commands=[ordered]@{rg=(Command-Record 'rg.exe' $rr);git_status=(Command-Record 'git.exe status --short' $sr);git_diff_check=(Command-Record 'git.exe diff --check' $dr)};overall_status=if($ok){'PASS'}else{'FAIL'};generated_at=(Get-Date).ToUniversalTime().ToString('o')}
 if($OutputPath){Write-Atomic $root $OutputPath $receipt};if($Json-or !$OutputPath){$receipt|ConvertTo-Json -Depth 12 -Compress};if($ok){exit 0}else{exit 1}
}catch{if($Json-or !$OutputPath){[ordered]@{schema_version=1;overall_status='STOP';error_code=($_.Exception.Message-replace'^local-review:','')}|ConvertTo-Json -Compress};exit 2}
