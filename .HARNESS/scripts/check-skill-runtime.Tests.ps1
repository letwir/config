$here = Split-Path -Parent $MyInvocation.MyCommand.Definition
$runner = Join-Path $here 'check-skill-runtime.ps1'

Describe 'H-06 skill runtime runner' {
  BeforeAll {
    $root = Join-Path $TestDrive 'fixture'; New-Item -ItemType Directory -Force $root | Out-Null
    $pwsh = (Get-Command pwsh.exe).Source; $manifest = Join-Path $root 'manifest.json'
    $base = [ordered]@{schemaVersion=1;maxOutputChars=32;approvedRoots=@('{PowerShell}','{HarnessRoot}');tools=@();clients=@()}
    $base.tools += [ordered]@{name='ok';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','Write-Output "catalog token=abc and a deliberately long suffix"');requiredText='token=abc';timeoutMs=10000;representative=$null}
    $base.tools += [ordered]@{name='slow';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','Start-Sleep -Seconds 3');timeoutMs=100;representative=$null}
    $base.clients += [ordered]@{name='OptionalMissing';required=$false;wrapper='{PowerShell}';target=(Join-Path $root 'missing.ps1');args=@();mode='optional';timeoutMs=100}
    $base.clients += [ordered]@{name='WrapperMismatch';required=$true;wrapper=(Join-Path $root 'missing-wrapper.exe');target=$pwsh;args=@();mode='catalog-only';timeoutMs=100}
    $base | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $manifest -Encoding utf8
    $script:manifest=$manifest; $script:root=$root
  }

  It 'separates executable, catalog, and representative status' {
    $out=& pwsh.exe -NoProfile -File $runner -ManifestPath $manifest -NoReceipt -Json 2>$null;$LASTEXITCODE|Should Be 1;$j=($out-join '')|ConvertFrom-Json;$ok=$j.tools|Where-Object name -eq ok
    $ok.executableHealthy|Should Be 'Healthy';$ok.catalogObserved|Should Be 'Observed';$ok.representativeOperation|Should Be 'UnTested'
  }
  It 'classifies timeout, wrapper mismatch, and optional absence' {
    $j=((& pwsh.exe -NoProfile -File $runner -ManifestPath $manifest -NoReceipt -Json 2>$null)-join '')|ConvertFrom-Json
    ($j.tools|Where-Object name -eq slow).classification|Should Be 'Timeout';($j.clients|Where-Object name -eq WrapperMismatch).classification|Should Be 'WrapperMismatch';($j.clients|Where-Object name -eq OptionalMissing).classification|Should Be 'NotConfigured'
  }
  It 'checks raw output before redaction and bounds stored output' {
    $j=((& pwsh.exe -NoProfile -File $runner -ManifestPath $manifest -NoReceipt -Json 2>$null)-join '')|ConvertFrom-Json;$o=$j.tools|Where-Object name -eq ok
    ($o.output.text.Length -le 32)|Should Be $true;$o.output.text|Should Match '<REDACTED>';$o.output.text|Should Not Match 'abc';$o.output.sha256|Should Match '^[A-F0-9]{64}$';$o.output.checkedBeforeNormalization|Should Be $true
  }
  It 'distinguishes nonzero, token mismatch, representative failure, and hash mismatch' {
    $m=[ordered]@{schemaVersion=1;maxOutputChars=256;approvedRoots=@('{PowerShell}','{HarnessRoot}');clients=@();tools=@(
      [ordered]@{name='nonzero';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','exit 7');timeoutMs=5000;representative=$null},
      [ordered]@{name='token';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','Write-Output other');requiredText='wanted';timeoutMs=5000;representative=$null},
      [ordered]@{name='representative';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','Write-Output help');timeoutMs=5000;representative=[ordered]@{args=@('-NoProfile','-Command','Write-Output wrong');requiredExitCode=0;requiredText='wanted'}},
      [ordered]@{name='hash';required=$true;path='{PowerShell}';sha256=('0'*64);helpArgs=@('--help');timeoutMs=5000;representative=$null})}
    $p=Join-Path $root 'failures.json';$m|ConvertTo-Json -Depth 10|Set-Content $p -Encoding utf8;$j=((& pwsh.exe -NoProfile -File $runner -ManifestPath $p -NoReceipt -Json 2>$null)-join '')|ConvertFrom-Json
    ($j.tools|Where-Object name -eq nonzero).classification|Should Be 'CommandNonZero';($j.tools|Where-Object name -eq token).classification|Should Be 'ContractMismatch';($j.tools|Where-Object name -eq representative).classification|Should Be 'RepresentativeFailed';($j.tools|Where-Object name -eq hash).classification|Should Be 'HashMismatch'
  }
  It 'kills a timed out process tree' {
    $pidFile=Join-Path $root 'child.pid';$childFile=Join-Path $root 'tree.ps1'
    @'
$child=Start-Process pwsh.exe -ArgumentList '-NoProfile','-Command','Start-Sleep -Seconds 30' -PassThru
Set-Content -LiteralPath $args[0] -Value $child.Id
Start-Sleep -Seconds 30
'@ | Set-Content -LiteralPath $childFile -Encoding utf8
    $m=[ordered]@{schemaVersion=1;maxOutputChars=64;approvedRoots=@('{PowerShell}','{HarnessRoot}');clients=@();tools=@([ordered]@{name='tree';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-File',$childFile,$pidFile);timeoutMs=5000;representative=$null})};$p=Join-Path $root 'tree.json';$m|ConvertTo-Json -Depth 10|Set-Content $p -Encoding utf8
    $null=& pwsh.exe -NoProfile -File $runner -ManifestPath $p -NoReceipt -Json 2>$null;(Test-Path $pidFile)|Should Be $true;$childPid=[int](Get-Content $pidFile);Start-Sleep -Milliseconds 300;(Get-Process -Id $childPid -ErrorAction SilentlyContinue)|Should BeNullOrEmpty
  }
  It 'rejects a reparse target that escapes an approved root' {
    $link=Join-Path $root 'escape';New-Item -ItemType Junction -Path $link -Target $here|Out-Null
    $m=[ordered]@{schemaVersion=1;maxOutputChars=64;approvedRoots=@('{PowerShell}','{HarnessRoot}');tools=@();clients=@([ordered]@{name='escape';required=$true;wrapper='{PowerShell}';target=(Join-Path $link 'check-skill-runtime.Tests.ps1');args=@();mode='catalog-only';timeoutMs=1000})};$p=Join-Path $root 'escape.json';$m|ConvertTo-Json -Depth 10|Set-Content $p -Encoding utf8;$j=((& pwsh.exe -NoProfile -File $runner -ManifestPath $p -NoReceipt -Json 2>$null)-join '')|ConvertFrom-Json;$j.clients[0].classification|Should Be 'PathRejected'
  }
  It 'writes an identical latest receipt only for success' {
    $m=[ordered]@{schemaVersion=1;maxOutputChars=64;approvedRoots=@('{PowerShell}');clients=@();tools=@([ordered]@{name='ok';required=$true;path='{PowerShell}';helpArgs=@('-NoProfile','-Command','Write-Output ok');requiredText='ok';timeoutMs=5000;representative=$null})};$p=Join-Path $root 'success.json';$m|ConvertTo-Json -Depth 10|Set-Content $p -Encoding utf8;$receipts=Join-Path $root 'receipts';$null=& pwsh.exe -NoProfile -File $runner -ManifestPath $p -OutputDirectory $receipts -Json 2>$null;$LASTEXITCODE|Should Be 0
    $latest=Get-Content -Raw (Join-Path $receipts 'skill-runtime-latest.json');$run=Get-ChildItem $receipts -Filter 'skill-runtime-*.json'|Where-Object Name -ne 'skill-runtime-latest.json'|Select-Object -First 1;$latest|Should Be (Get-Content -Raw $run.FullName)
  }
}
