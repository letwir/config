# Knowledge Base: 20261025-native-worker-and-ffmpeg

<api id="pi-native-subagent-sdk-shadow">
## Context
Pi 1.4.0 managed SDK, pi-subagents 0.82.0; legacy ~/.pi/package.json claimed SDK identity 0.85.1 while ~/.pi/dist was absent. Native detached workers failed before model execution.

## Findings
The resolver's direct-JS-entry ownership check walks ancestors past the pi-subagents nested package to the legacy manifest. That branch runs before the explicit SDK-root override. Pure local reproductions with both actual runner entry paths returned the legacy root even when both environment roots were valid. User-approved name/private/description metadata correction removed the false SDK identity; native worker-medium then read, wrote, tested, corrected and completed using actual openai-codex/gpt-6-luna at medium. Local source inspections and actual session messages, not self-reported model names, provide evidence. Other environments/host versions are untested.

## Morphism
Runner entry + package ancestry + root override -> resolve SDK root -> import actual SDK -> native model session -> parent-verified result. A missing-dist error before model execution is infrastructure evidence, not model-quality evidence. Removing an ancestor false identity must preserve other configuration and needs approval if it affects a legacy installation.
</api>

<api id="ffmpeg-batch-update-validation">
## Context
User requested C:\PATH dedicated latest FFmpeg/ffprobe with repeatable batch updates. FFmpeg official downloads link to BtbN; BtbN documents floating latest win64 GPL static builds. Sources inspected 2026-10-25: https://ffmpeg.org/download.html and https://raw.githubusercontent.com/BtbN/FFmpeg-Builds/master/README.md.

## Findings
PowerShell parser success and negative-only fixtures do not prove valid inputs work. Real successful full-flow fixtures exposed guard regressions; a real cmd.exe invocation exposed an unescaped parenthesis inside IF. ResponseHeadersRead needs an explicit streaming cancellation token. Suppress returned HTTP objects to avoid signed-URL output. Check ffmpeg and ffprobe banners separately, require matching version IDs, hash before executable probes, and never indiscriminately extract ZIP paths. Full parent verification passed16 offline checks, real CheckOnly/install, both versions, command resolution and a synthetic audio pipeline.

## Morphism
Approved fixed HTTPS release metadata -> exact asset/checksum -> archive integrity -> exact pair staging -> both version probes -> lock-protected backup/install/metadata -> rollback on failure -> parent live behavior checks. Same-origin hashes do not authenticate a compromised upstream publisher. Tests and review required; no independent review claimed here.
</api>
