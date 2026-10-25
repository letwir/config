# Install and build staged skills

This directory documents source trees copied from sibling repository `../mcp2go/.harness/skill/src/` into `.HARNESS/skills/<skill-name>/mcp2go-src/`. Each Go module remains an independent build boundary; do not combine module roots or assume a workspace-wide `go test ./...` is valid. Python-backed skills retain their original scripts and Python dependencies.

## Repeatable Go build check (Windows PowerShell)

From the config repository root, run this PowerShell command. It discovers only `go.mod` files underneath the imported `mcp2go-src` subtrees, builds each module's packages with isolated temporary cache/output locations, and removes only the temporary directory it creates. No executable is written into a skill source tree.

```powershell
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path '.HARNESS/skills').Path
$temp = Join-Path ([System.IO.Path]::GetTempPath()) ('mcp2go-skill-build-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $temp | Out-Null
try {
    $env:GOCACHE = Join-Path $temp 'gocache'
    $env:GOTMPDIR = Join-Path $temp 'gotmp'
    New-Item -ItemType Directory -Path $env:GOTMPDIR | Out-Null
    $modules = Get-ChildItem -Path $root -Filter go.mod -Recurse | Where-Object {
        $_.FullName -match '[\\/]mcp2go-src[\\/]' -and
        # Exclude only Kavita's obsolete outer module; its nested src module is included.
        $_.FullName -notmatch '[\\/]kavita-fetch-search[\\/]mcp2go-src[\\/]go\.mod$'
    } | Sort-Object FullName
    foreach ($module in $modules) {
        $moduleDir = $module.DirectoryName
        Write-Host "Testing/building $($module.FullName)"
        Push-Location $moduleDir
        try {
            $packages = @(go list ./...)
            if ($LASTEXITCODE -ne 0) { throw "go list failed in $moduleDir (exit $LASTEXITCODE)" }
            foreach ($package in $packages) {
                $output = Join-Path $temp ([guid]::NewGuid().ToString('N'))
                go build -o $output $package
                if ($LASTEXITCODE -ne 0) { throw "go build failed for $package in $moduleDir (exit $LASTEXITCODE)" }
            }
            go test -vet=off ./...
            if ($LASTEXITCODE -ne 0) { throw "go test failed in $moduleDir (exit $LASTEXITCODE)" }
        }
        finally { Pop-Location }
    }
} finally {
    Remove-Item -LiteralPath $temp -Recurse -Force
}
```

The per-package `go build -o <temporary-path>` loop compiles without placing executables in source trees, and `go test -vet=off ./...` runs tests. Vet is disabled because default vet reports the existing `gemini-grounding-search` redundant-newline warning in `main.go:122`; that warning remains unresolved. Keep `GOCACHE` and `GOTMPDIR` outside this tree. Some modules require newer Go releases; check each module's `go` directive before treating a toolchain mismatch as a source failure. Kavita's supported build root is `.HARNESS/skills/kavita-fetch-search/mcp2go-src/src`; its outer `.HARNESS/skills/kavita-fetch-search/mcp2go-src/go.mod` is an obsolete CLI duplicate excluded by the command. The nested module is included and built/tested normally.

Python modules do not have a common repository-wide build command. Follow each skill's own instructions and use a disposable virtual environment when validation is required.
