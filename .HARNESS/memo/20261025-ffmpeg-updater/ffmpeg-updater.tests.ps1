[CmdletBinding()]
param()
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$root = Join-Path $PSScriptRoot ('fixtures/' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $root)
$passed = 0
function Assert([bool] $Condition, [string] $Message) { if (-not $Condition) { throw "ASSERTION FAILED: $Message" } }
function New-FixtureArchive([string] $Path, [string] $FfmpegVersion, [string] $FfprobeVersion) {
    $zip = [IO.Compression.ZipFile]::Open($Path, [IO.Compression.ZipArchiveMode]::Create)
    try {
        foreach ($tool in @(@('ffmpeg', $FfmpegVersion), @('ffprobe', $FfprobeVersion))) {
            $entry = $zip.CreateEntry("build/bin/$($tool[0]).exe")
            $stream = $entry.Open(); $bytes = [Text.Encoding]::UTF8.GetBytes("fake-$($tool[0])`n$($tool[0]) version $($tool[1])`n")
            $stream.Write($bytes, 0, $bytes.Length); $stream.Dispose()
        }
    } finally { $zip.Dispose() }
}
try {
    $cmdExe = Join-Path $env:SystemRoot 'System32/cmd.exe'
    $batch = 'C:\PATH\update-ffmpeg.cmd'
    & $cmdExe /d /c ('"' + $batch + '" -LibraryOnly')
    $libraryExit = $LASTEXITCODE
    Assert ($libraryExit -eq 0) "Real batch LibraryOnly exits 0 (actual: $libraryExit)."
    $batchInstall = Join-Path $root 'batch-install'; [void](New-Item -ItemType Directory -Path $batchInstall)
    & $cmdExe /d /c ('"' + $batch + '" -InstallDirectory "' + $batchInstall + '" -ApiUri "https://invalid.example/"')
    $invalidSourceExit = $LASTEXITCODE
    Assert ($invalidSourceExit -eq 1) "Real batch rejects invalid fixed-source ApiUri with exit 1 (actual: $invalidSourceExit)."
    Assert (-not (Get-ChildItem -LiteralPath $batchInstall -Force | Where-Object { $_.Name -in @('ffmpeg.exe', 'ffprobe.exe') })) 'Invalid-source batch smoke writes no binaries to disposable install directory.'
    $passed += 3
    . 'C:\PATH\update-ffmpeg.ps1' -LibraryOnly
    $archiveName = 'ffmpeg-master-latest-win64-gpl.zip'
    $checksum = Join-Path $root 'checksums.sha256'
    [IO.File]::WriteAllText($checksum, (('a' * 64) + '  other.zip' + "`n" + ('b' * 64) + ' *' + $archiveName + "`n"), [Text.Encoding]::UTF8)
    Assert ((Get-ExpectedChecksum $checksum) -eq ('b' * 64)) 'Expected exact archive checksum entry.'
    $passed++
    [IO.File]::WriteAllText($checksum, (('b' * 64) + '  ' + $archiveName + "`n" + ('c' * 64) + '  ' + $archiveName))
    $checksumError = ''
    try { Get-ExpectedChecksum $checksum | Out-Null } catch { $checksumError = $_.Exception.Message }
    Assert ($checksumError -match '^Checksum file must contain exactly one unambiguous entry') 'Duplicate checksum entries fail for the intended ambiguity reason.'
    $passed++

    Assert ((Convert-VersionProbeOutput "ffmpeg version 7.1`n" 'ffmpeg') -match '^ffmpeg version 7.1') 'ffmpeg probe accepts its own version banner.'
    Assert ((Convert-VersionProbeOutput "ffprobe version 7.1`n" 'ffprobe') -match '^ffprobe version 7.1') 'ffprobe probe accepts its own version banner.'
    $probeError = ''
    try { Convert-VersionProbeOutput "ffmpeg version 7.1`n" 'ffprobe' | Out-Null } catch { $probeError = $_.Exception.Message }
    Assert ($probeError -match '^Version probe output did not identify ffprobe\.') 'Probe rejects a banner for the other executable with the expected identification error.'
    $ffmpegId = [regex]::Match((Convert-VersionProbeOutput 'ffmpeg version 7.1' 'ffmpeg'), '(?i)version\s+(\S+)').Groups[1].Value
    $ffprobeId = [regex]::Match((Convert-VersionProbeOutput 'ffprobe version 7.2' 'ffprobe'), '(?i)version\s+(\S+)').Groups[1].Value
    Assert ($ffmpegId -cne $ffprobeId) 'Different tool version identifiers are distinguishable and rejected by pair validation.'
    $passed++

    $install = Join-Path $root 'install'; $stage = Join-Path $root 'stage'; $backup = Join-Path $root 'backup'
    foreach ($dir in @($install, $stage, $backup)) { [void](New-Item -ItemType Directory -Path $dir) }
    Assert-ExecutableTargetsAreFiles $install
    [IO.File]::WriteAllText((Join-Path $install 'ffmpeg.exe'), 'guard-file')
    Assert-ExecutableTargetsAreFiles $install
    Remove-Item -LiteralPath (Join-Path $install 'ffmpeg.exe') -Force
    Assert-ExecutableTargetsAreFiles $install
    $passed++
    [IO.File]::WriteAllText((Join-Path $install 'ffmpeg.exe'), 'old-ffmpeg')
    [IO.File]::WriteAllText((Join-Path $install 'ffprobe.exe'), 'old-ffprobe')
    [IO.File]::WriteAllText((Join-Path $install 'keep.txt'), 'unrelated')
    Copy-Item (Join-Path $install 'ffmpeg.exe') (Join-Path $backup 'ffmpeg.exe')
    Copy-Item (Join-Path $install 'ffprobe.exe') (Join-Path $backup 'ffprobe.exe')
    [IO.File]::WriteAllText((Join-Path $stage 'ffmpeg.exe'), 'new-ffmpeg')
    [IO.File]::WriteAllText((Join-Path $stage 'ffprobe.exe'), 'new-ffprobe')
    Install-Pair $install $stage $backup $true $true
    Assert ([IO.File]::ReadAllText((Join-Path $install 'ffmpeg.exe')) -eq 'new-ffmpeg') 'Successful Install-Pair replaces ffmpeg.'
    Assert ([IO.File]::ReadAllText((Join-Path $install 'ffprobe.exe')) -eq 'new-ffprobe') 'Successful Install-Pair replaces ffprobe.'
    [IO.File]::WriteAllText((Join-Path $install 'ffmpeg.exe'), 'old-ffmpeg')
    [IO.File]::WriteAllText((Join-Path $install 'ffprobe.exe'), 'old-ffprobe')
    Remove-Item -LiteralPath (Join-Path $stage 'ffprobe.exe') -Force
    $failureMessage = ''
    try { Install-Pair $install $stage $backup $true $true } catch { $failureMessage = $_.Exception.Message }
    Assert ($failureMessage -match '^Pair update failed; original executable pair restored:') 'Missing staged ffprobe must fail during replacement and report successful rollback.'
    Assert ([IO.File]::ReadAllText((Join-Path $install 'ffmpeg.exe')) -eq 'old-ffmpeg') 'First executable restored after partial pair replacement.'
    Assert ([IO.File]::ReadAllText((Join-Path $install 'ffprobe.exe')) -eq 'old-ffprobe') 'Second executable unchanged after partial pair replacement.'
    Assert ([IO.File]::ReadAllText((Join-Path $install 'keep.txt')) -eq 'unrelated') 'Unrelated file remains untouched.'
    $passed++

    Remove-Item -LiteralPath (Join-Path $install 'ffmpeg.exe') -Force
    [void](New-Item -ItemType Directory -Path (Join-Path $install 'ffmpeg.exe'))
    $directoryError = ''
    try { Install-Pair $install $stage $backup $false $true } catch { $directoryError = $_.Exception.Message }
    Assert ($directoryError -match '^Refusing non-file executable target: ffmpeg\.exe$') 'Non-file executable destination is rejected by the pre-replacement guard.'
    Assert (Test-Path -LiteralPath (Join-Path $install 'ffmpeg.exe') -PathType Container) 'Existing executable-named directory remains intact.'
    $passed++
    Remove-Item -LiteralPath (Join-Path $install 'ffmpeg.exe') -Recurse -Force
    [IO.File]::WriteAllText((Join-Path $install 'ffmpeg.exe'), 'old-ffmpeg')

    Remove-Item -LiteralPath (Join-Path $stage 'ffmpeg.exe') -Force
    $zipPath = Join-Path $root 'fixture.zip'
    Add-Type -AssemblyName System.IO.Compression
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $zip = [IO.Compression.ZipFile]::Open($zipPath, [IO.Compression.ZipArchiveMode]::Create)
    try {
        foreach ($leaf in @('ffmpeg.exe', 'ffprobe.exe')) {
            $entry = $zip.CreateEntry("build/bin/$leaf")
            $stream = $entry.Open(); $bytes = [Text.Encoding]::UTF8.GetBytes("fake-$leaf")
            $stream.Write($bytes, 0, $bytes.Length); $stream.Dispose()
        }
    } finally { $zip.Dispose() }
    $wrongHashError = ''
    try { Assert-ArchiveHash $zipPath ('0' * 64) '' | Out-Null } catch { $wrongHashError = $_.Exception.Message }
    Assert ($wrongHashError -match '^Archive SHA-256 mismatch') 'Wrong archive checksum rejected before installation.'
    Assert ([IO.File]::ReadAllText((Join-Path $install 'ffmpeg.exe')) -eq 'old-ffmpeg') 'Checksum rejection preserves existing pair.'
    $passed++
    Get-ZipExecutables $zipPath $stage
    Assert (Test-Path (Join-Path $stage 'ffmpeg.exe')) 'Valid archive extracted expected binary.'
    $passed++
    $badZipPath = Join-Path $root 'bad.zip'
    $zip = [IO.Compression.ZipFile]::Open($badZipPath, [IO.Compression.ZipArchiveMode]::Create)
    try {
        $entry = $zip.CreateEntry('../ffmpeg.exe'); $stream = $entry.Open(); $stream.Dispose()
    } finally { $zip.Dispose() }
    $badZipError = ''
    try { Get-ZipExecutables $badZipPath $stage } catch { $badZipError = $_.Exception.Message }
    Assert ($badZipError -match '^Archive must contain exactly one canonical executable pair\.') 'Unexpected ZIP entry fails for the intended archive-pair reason.'
    $passed++
    $incompletePath = Join-Path $root 'incomplete.zip'
    $zip = [IO.Compression.ZipFile]::Open($incompletePath, [IO.Compression.ZipArchiveMode]::Create)
    try { $entry = $zip.CreateEntry('build/bin/ffmpeg.exe'); $stream = $entry.Open(); $stream.WriteByte(1); $stream.Dispose() }
    finally { $zip.Dispose() }
    $incompleteError = ''
    try { Get-ZipExecutables $incompletePath $stage } catch { $incompleteError = $_.Exception.Message }
    Assert ($incompleteError -match '^Archive must contain exactly one canonical executable pair\.') 'Incomplete ZIP pair fails for the intended missing-pair reason.'
    $passed++

    $lockPath = Join-Path $root 'lock'
    $held = [IO.File]::Open($lockPath, [IO.FileMode]::OpenOrCreate, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None)
    try {
        $lockError = $null
        try { $contender = [IO.File]::Open($lockPath, [IO.FileMode]::OpenOrCreate, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None); $contender.Dispose() } catch { $lockError = $_.Exception }
        Assert ($lockError.InnerException -is [IO.IOException] -and $lockError.Message -match 'being used by another process') 'Second exclusive update lock fails specifically because the held lock denies sharing.'
    } finally { $held.Dispose() }
    $passed++
    # Mock only release/download/version boundaries; checksum, ZIP extraction, install and metadata are real.
    $flowRoot = Join-Path $root 'integration-success'
    $flowInstall = Join-Path $flowRoot 'install'; [void](New-Item -ItemType Directory -Path $flowInstall -Force)
    $flowZip = Join-Path $flowRoot $archiveName
    $InstallDirectory = $flowInstall; $script:FixtureZip = $flowZip; $script:FixtureChecksum = Join-Path $flowRoot 'checksums.sha256'
    New-FixtureArchive $flowZip '7.1' '7.1'
    $flowHash = (Get-FileHash -LiteralPath $flowZip -Algorithm SHA256).Hash.ToLowerInvariant(); $script:FixtureHash = $flowHash
    $flowChecksum = $script:FixtureChecksum
    [IO.File]::WriteAllText($flowChecksum, "$flowHash  *$archiveName`n")
    function Get-ReleaseAssets { [pscustomobject]@{ Tag='fixture-tag'; Archive=[pscustomobject]@{ name=$script:ArchiveName; browser_download_url='fixture://zip'; digest="sha256:$script:FixtureHash" }; Checksums=[pscustomobject]@{ name=$script:ChecksumName; browser_download_url='fixture://checksum'; digest=$null } } }
    function Save-VerifiedAsset([string] $Uri, [string] $Destination) { Copy-Item -LiteralPath $(if ($Uri -eq 'fixture://zip') { $script:FixtureZip } else { $script:FixtureChecksum }) -Destination $Destination }
    function Invoke-VersionProbe([string] $Executable, [int] $TimeoutSeconds) { $tool = [IO.Path]::GetFileNameWithoutExtension($Executable); [IO.File]::ReadAllText($Executable) -split "`n" | Where-Object { $_ -match "^$tool version " } | Select-Object -First 1 }
    Invoke-Update | Out-Null
    Assert ([IO.File]::ReadAllText((Join-Path $flowInstall 'ffmpeg.exe')) -match 'ffmpeg version 7\.1') 'Offline full-flow integration installs extracted ffmpeg fixture.'
    Assert ([IO.File]::ReadAllText((Join-Path $flowInstall 'ffprobe.exe')) -match 'ffprobe version 7\.1') 'Offline full-flow integration installs extracted ffprobe fixture.'
    $metadata = Get-Content -Raw -LiteralPath (Join-Path $flowInstall '.ffmpeg-update/installed.json') | ConvertFrom-Json
    Assert ($metadata.tag -eq 'fixture-tag' -and $metadata.sha256 -eq $flowHash -and $metadata.ffmpegVersion -eq 'ffmpeg version 7.1') 'Full-flow integration writes matching installation metadata.'
    $passed++

    $mismatchInstall = Join-Path (Join-Path $root 'integration-mismatch') 'install'; [void](New-Item -ItemType Directory -Path $mismatchInstall -Force)
    [IO.File]::WriteAllText((Join-Path $mismatchInstall 'ffmpeg.exe'), 'original-ffmpeg')
    [IO.File]::WriteAllText((Join-Path $mismatchInstall 'ffprobe.exe'), 'original-ffprobe')
    $mismatchZip = Join-Path $root 'mismatch.zip'; New-FixtureArchive $mismatchZip '7.1' '7.2'
    $mismatchHash = (Get-FileHash -LiteralPath $mismatchZip -Algorithm SHA256).Hash.ToLowerInvariant()
    $mismatchChecksum = Join-Path $root 'mismatch-checksums.sha256'; [IO.File]::WriteAllText($mismatchChecksum, "$mismatchHash  *$archiveName`n")
    $InstallDirectory = $mismatchInstall; $script:FixtureZip = $mismatchZip; $script:FixtureChecksum = $mismatchChecksum; $script:FixtureHash = $mismatchHash
    function Get-ReleaseAssets { [pscustomobject]@{ Tag='fixture-mismatch'; Archive=[pscustomobject]@{ name=$script:ArchiveName; browser_download_url='fixture://mismatch-zip'; digest="sha256:$script:FixtureHash" }; Checksums=[pscustomobject]@{ name=$script:ChecksumName; browser_download_url='fixture://mismatch-checksum'; digest=$null } } }
    function Save-VerifiedAsset([string] $Uri, [string] $Destination) { Copy-Item -LiteralPath $(if ($Uri -eq 'fixture://mismatch-zip') { $script:FixtureZip } else { $script:FixtureChecksum }) -Destination $Destination }
    $mismatchError = ''
    try { Invoke-Update | Out-Null } catch { $mismatchError = $_.Exception.Message }
    Assert ($mismatchError -match '^Staged ffmpeg and ffprobe versions do not match\.') "Mismatched full-flow pair fails at pair-version validation (actual: $mismatchError)."
    Assert ([IO.File]::ReadAllText((Join-Path $mismatchInstall 'ffmpeg.exe')) -eq 'original-ffmpeg') 'Mismatched pair leaves original ffmpeg content unchanged.'
    Assert ([IO.File]::ReadAllText((Join-Path $mismatchInstall 'ffprobe.exe')) -eq 'original-ffprobe') 'Mismatched pair leaves original ffprobe content unchanged.'
    $passed++
    Write-Output "PASS: $passed deterministic updater fixture checks."
} finally {
    if (Test-Path -LiteralPath $root) { Remove-Item -LiteralPath $root -Recurse -Force }
}
