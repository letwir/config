param (
    [string]$Output = ""
)

$urlFile = "url.txt"
$destDir = "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\llama2coder"

if (-not (Test-Path $urlFile)) {
    Write-Error "url.txt not found. Please copy url.txt.temp to url.txt and configure it."
    exit 1
}

$lines = Get-Content $urlFile
$buildCount = 0

foreach ($line in $lines) {
    $trimmed = $line.Trim()
    if ($trimmed -eq "" -or $trimmed.StartsWith("#")) {
        continue
    }

    # スペースで分割
    $parts = -split $trimmed
    if ($parts.Length -lt 2) {
        Write-Warning "Invalid line format (expected '<name> <url>'): $trimmed"
        continue
    }

    $name = $parts[0]
    $url = $parts[1]

    Write-Output "Preparing build for $name ($url)..."

    # 出力ファイル名決定
    $winOut = "../../../dist/coder-$name.exe"
    $linuxOut = "../../../dist/coder-$name"

    # Windows (amd64) ビルド
    Write-Output "  Building Windows binary: coder-$name.exe"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go.exe build -ldflags "-X main.embeddedURL=$url" -o $winOut main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Output "  Succeeded: $winOut"
        if (Test-Path $destDir) {
            Copy-Item -Path $winOut -Destination "$destDir\coder-$name.exe" -Force
        }
    } else {
        Write-Error "  Failed to build Windows binary for $name"
    }

    # Linux (amd64) ビルド
    Write-Output "  Building Linux (Debian) binary: coder-$name"
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    go.exe build -ldflags "-X main.embeddedURL=$url" -o $linuxOut main.go
    if ($LASTEXITCODE -eq 0) {
        Write-Output "  Succeeded: $linuxOut"
        if (Test-Path $destDir) {
            Copy-Item -Path $linuxOut -Destination "$destDir\coder-$name" -Force
        }
    } else {
        Write-Error "  Failed to build Linux binary for $name"
    }

    $buildCount++
}

# 環境変数を元に戻す
$env:GOOS = $null
$env:GOARCH = $null

Write-Output "Completed. Processed $buildCount targets."
