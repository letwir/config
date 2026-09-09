[CmdletBinding()]
param([Parameter(Mandatory)][string]$Target,[Parameter(Mandatory)][string]$ArgsBase64,[Parameter(Mandatory)][string]$OutPath,[Parameter(Mandatory)][string]$ErrPath)
$ErrorActionPreference = 'Stop'
$env:CI = 'true'
$env:NO_COLOR = '1'
$json = [Text.Encoding]::Unicode.GetString([Convert]::FromBase64String($ArgsBase64))
$clientArguments = @($json | ConvertFrom-Json)
& $Target @clientArguments 1> $OutPath 2> $ErrPath
exit $LASTEXITCODE
