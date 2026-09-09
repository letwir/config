$script:Runner=Join-Path $PSScriptRoot 'invoke-local-review.ps1'
$script:Root=Join-Path $PSScriptRoot '..'
function Invoke-Review([string]$Q='local-review',[string]$T='rules/research.lrf,etl/main.seq',[string]$O=''){$a=@('-NoProfile','-ExecutionPolicy','Bypass','-File',$script:Runner,'-RootPath',$script:Root,'-Query',$Q,'-Target',$T,'-Json');if($O){$a+=@('-OutputPath',$O)};$out=& powershell.exe @a 2>&1;[pscustomobject]@{Exit=$LASTEXITCODE;Text=($out-join'')}}
Describe 'H-08 bounded local review collector' {
 It 'runs the fixed collectors and emits only bounded hashes' {
  $r=Invoke-Review;$r.Exit|Should Be 0;$o=$r.Text|ConvertFrom-Json;$o.schema_version|Should Be 1;$o.overall_status|Should Be 'PASS';$o.query.mode|Should Be 'fixed_string';@($o.targets).Count|Should Be 2
  @($o.commands.PSObject.Properties.Name)|Should Be @('rg','git_status','git_diff_check');$o.commands.rg.executable|Should Be 'rg.exe';$o.commands.git_status.executable|Should Be 'git.exe status --short';$o.commands.git_diff_check.executable|Should Be 'git.exe diff --check'
  $r.Text|Should Not Match '"text"|args|local-review\s+route';foreach($c in @($o.commands.rg,$o.commands.git_status,$o.commands.git_diff_check)){$c.bounded|Should Be $true;$c.output_sha256|Should Match '^[A-F0-9]{64}$'}
 }
 It 'accepts a fixed-string no-match as completed evidence' {$r=Invoke-Review 'definitely-absent-h08-fixed-string';$r.Exit|Should Be 0;($r.Text|ConvertFrom-Json).commands.rg.exit_code|Should Be 1}
 It 'fails closed for malformed finite targets' {
  $bad=@((Join-Path $script:Root 'rules\research.lrf'),'..\rules\research.lrf','rules\*.lrf','rules/research.lrf,rules/research.lrf','missing-h08.txt','rules')
  foreach($x in $bad){$r=Invoke-Review 'local-review' $x;$r.Exit|Should Be 2;($r.Text|ConvertFrom-Json).overall_status|Should Be 'STOP'}
 }
 It 'rejects empty multiline and overlong queries and too many targets' {
  foreach($q in @(' ',("a`nb"),('x'*513))){(Invoke-Review $q).Exit|Should Be 2}
  $many=(1..33|ForEach-Object{'rules/research.lrf'})-join ',';(Invoke-Review 'local-review' $many).Exit|Should Be 2
 }
 It 'writes a schema-valid receipt atomically only at a contained requested path' {
  $dir=Join-Path $script:Root ('evaluation\runs\h08-test-'+[guid]::NewGuid().ToString('N'));$path=Join-Path $dir 'receipt.json'
  try{$r=Invoke-Review 'local-review' 'rules/research.lrf,etl/main.seq' $path;$r.Exit|Should Be 0;Test-Path -LiteralPath $path|Should Be $true;$raw=Get-Content -Raw $path;$saved=$raw|ConvertFrom-Json;if(Get-Command Test-Json -ErrorAction SilentlyContinue){$raw|Test-Json -SchemaFile (Join-Path $script:Root 'evaluation\local-review.schema.json')|Should Be $true}else{$saved.receipt_id|Should Match '^[0-9a-f-]{36}$';$saved.commands.rg.output_sha256|Should Match '^[A-F0-9]{64}$'};@(Get-ChildItem $dir -Filter '*.tmp').Count|Should Be 0}finally{Remove-Item -LiteralPath $dir -Recurse -Force -ErrorAction SilentlyContinue}
  (Invoke-Review 'local-review' 'rules/research.lrf' (Join-Path ([IO.Path]::GetTempPath()) 'escaped-h08.json')).Exit|Should Be 2
 }
 It 'does not create FINISH or memory artifacts' {
  $before=@(Get-ChildItem -LiteralPath (Join-Path $script:Root 'memo') -ErrorAction SilentlyContinue).Count;$null=Invoke-Review;$after=@(Get-ChildItem -LiteralPath (Join-Path $script:Root 'memo') -ErrorAction SilentlyContinue).Count;$after|Should Be $before
 }
}
