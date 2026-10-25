# Build script for economics-api

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$pluginDir = "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\economics-api"

if (-not (Test-Path $pluginDir)) {
    New-Item -ItemType Directory -Path $pluginDir -Force | Out-Null
}

Set-Location $scriptDir

Write-Host "[1/8] Building economics.exe (Wrapper)..."
go build -o "$pluginDir\economics.exe" .\cmd\economics\main.go

Write-Host "[2/8] Building estat.exe..."
go build -o "$pluginDir\estat.exe" .\cmd\estat\main.go

Write-Host "[3/8] Building edinet.exe..."
go build -o "$pluginDir\edinet.exe" .\cmd\edinet\main.go

Write-Host "[4/8] Building sec_edgar.exe..."
go build -o "$pluginDir\sec_edgar.exe" .\cmd\sec_edgar\main.go

Write-Host "[5/8] Building fred.exe..."
go build -o "$pluginDir\fred.exe" .\cmd\fred\main.go

Write-Host "[6/8] Building worldbank.exe..."
go build -o "$pluginDir\worldbank.exe" .\cmd\worldbank\main.go

Write-Host "[7/8] Building yfinance.exe..."
go build -o "$pluginDir\yfinance.exe" .\cmd\yfinance\main.go

Write-Host "[8/8] Building alphavantage.exe..."
go build -o "$pluginDir\alphavantage.exe" .\cmd\alphavantage\main.go

Write-Host "[Copying Docs] Copying SKILL.md and USAGE.md..."
Copy-Item "$scriptDir\SKILL.md" -Destination "$pluginDir\SKILL.md" -Force
Copy-Item "$scriptDir\USAGE.md" -Destination "$pluginDir\USAGE.md" -Force

Write-Host "[SUCCESS] All binaries and docs built & deployed to $pluginDir successfully!"
