# Implementation Plan: 20261025-ffmpeg-updater

## Purpose and authorized scope
User approved installing dedicated FFmpeg/ffprobe in C:\PATH and adding a repeatable batch updater for latest builds. C:\PATH already appears on PATH. Do not use or alter the unrelated app-bundled copies, PowerShell profile, global PATH, or any other tools. Main owns local research/plan/review; one worker-medium owns this tightly coupled updater workstream. No VCS writes, credentials, external publishing, scheduled tasks, or elevated installation.

## Facts and sources (2026-10-25)
- No existing update-ffmpeg.cmd, update-ffmpeg.ps1, ffmpeg.exe or ffprobe.exe was found at C:\PATH.
- PowerShell 7.6.7 is installed. Windows x64 active profile.
- ffmpeg.org/download.html links to BtbN/FFmpeg-Builds as a Windows build provider.
- BtbN README documents floating latest, win64 x86_64, GPL and static builds. Source: https://raw.githubusercontent.com/BtbN/FFmpeg-Builds/master/README.md.
- GitHub releases/latest API returned tag latest, published 2026-10-24T13:21:42Z, assets checksums.sha256 and ffmpeg-master-latest-win64-gpl.zip (~200 MB). Re-query at execution; do not hardcode a dated version or a checksum from this note.
- Same-origin SHA-256 protects integrity/consistency, not independent publisher authenticity. Trust remains FFmpeg-listed BtbN + GitHub HTTPS; do not market this as signature verification.
- Earlier native smoke runs failed before CODE due compiled-host package path discovery. Current shell now has valid PI_PACKAGE_DIR and PI_SUBAGENTS_PI_CODING_AGENT_PACKAGE_ROOT overrides. Native execution is to be retried under same provider/profile, never another CLI/model.

## Components and exclusive ownership
Worker may create ONLY:
1. C:\PATH\update-ffmpeg.cmd: batch entry point invoking the adjacent PowerShell helper using pwsh -NoLogo -NoProfile -File; preserve arguments and exit code; no execution-policy bypass; clear prerequisite error if pwsh unavailable.
2. C:\PATH\update-ffmpeg.ps1: readable PowerShell 7 updater, default install directory=PSScriptRoot. Dedicated -InstallDirectory parameter is acceptable for disposable tests. Provide -CheckOnly for metadata/status without installation. No startup/profile hook.
3. C:\Users\letwir\.harness\memo\20261025-ffmpeg-updater\ffmpeg-updater.tests.ps1: self-contained deterministic local fixture tests; no network or external test dependency.
Tests may create disposable fixtures ONLY under C:\Users\letwir\.harness\memo\20261025-ffmpeg-updater\fixtures\ (unique owned children). Runtime installation may use ONLY target ffmpeg.exe, ffprobe.exe and an owned C:\PATH\.ffmpeg-update\ subtree for lock/staging/metadata/backups. Worker does not download or install live binaries; main performs first live install after actual-source review and deterministic checks.

## Required behavior
- Fetch fixed HTTPS GitHub BtbN/FFmpeg-Builds release metadata/asset URLs, exact static master-latest-win64-gpl.zip plus checksums.sha256. Validate exact expected asset names/host/repo; reject missing/duplicate checksum entries. Decode checksum response bytes correctly. Cross-check GitHub SHA256 asset digest when supplied. No arbitrary URLs or shell-evaluated remote text. Network timeouts bounded; failure exits nonzero before touching installed pair. Floating latest can race: fail safely on mismatch, explain retry, never relax verification.
- Verify downloaded archive SHA-256 before any execution. Extract only exact archive bin/ffmpeg.exe and bin/ffprobe.exe entries via ZipArchive streams; reject missing/duplicate/ambiguous entries. Never extract attacker-selected ZIP paths to arbitrary locations. Preserve unrelated files.
- Stage both executables in owned same-volume storage; run -version for BOTH with bounded timeout and exit-code checks before replacement. Do not kill existing user processes.
- One filesystem update lock; no simultaneous installs. Dispose streams, processes, zip archives and lock reliably. Refuse unsafe reparse/symlink target paths.
- Back up old pair before replacing. On partial replacement failure restore original presence/content of both files; if rollback itself fails, retain backups and clearly report their exact recovery location. Never report success or delete the only good backup on failure. First install failure leaves no half-installed pair where rollback succeeds.
- Record installed provenance/build/hash metadata safely only after successful pair replacement. If metadata participates in failure, handle it consistently with rollback. Avoid unnecessary downloads when remote identity and installed local hashes demonstrably match; no-update skip is optional if it would overcomplicate the implementation.
- Print concise result, installed version and update command; script failures and batch exit codes must be meaningful. No blanket silent catch or suppression of critical failures.

## Acceptance and checks
- Both scripts parse; batch launches helper and propagates failure.
- Fixture tests cover wrong checksum (existing pair untouched), incomplete/ambiguous ZIP (existing pair untouched), partial replacement rollback (original pair restored), lock conflict (no install), and unrelated-file preservation. Prefer testable narrow functions and mocked downloads/version probes; fixtures never execute fake downloaded binaries.
- Main inspects actual new files, reruns deterministic checks, then runs one live update and absolute-path -version for both; confirm command resolution from clean pwsh with unchanged PATH.
- No independent Attacker or Verifier run claimed: exact approved Pi CODE profiles are available, but independently routed matching review profiles are not registered. Main performs separately labeled pre-design security review and post-source review; this independence limitation is explicit.

## Main PRE_VERIFY / security design assessment
PASS for bounded implementation after user effect approval. Single seam, no artificial parallel split. Trust boundaries reviewed: mutable release metadata/checksum pair, ZIP paths, execution of verified third-party binaries, existing local files and rollback, concurrent updaters. Defensive requirements above address checksum races, path traversal, partial installs and symlink hazards. Remaining risk is upstream account/build compromise; same-origin hashes are not independent authenticity proof. No exploit reproduction or security scanning requested or authorized. Main is not an independent auditor.

## Completion evidence
Require actual worker terminal result and correlated provider/model verification (openai-codex/gpt-6-luna,medium), actual file inspection/test exit codes, first-install result or explicit blocker, and bounded task-memory receipt. No model fallback on operational failure.
