# Walkthrough: 20261025-ffmpeg-updater completion

## Overview
Completed native Pi subagent operability verification and the requested dedicated C:\PATH FFmpeg/ffprobe updater. User explicitly prioritized native subagent operation and authorized overwriting the shadowing legacy manifest rather than bypassing the worker. No alternate model, foreground workaround, external coding CLI or Main-only CODE implementation was used.

## Deliverables and actual changes
- C:\PATH\update-ffmpeg.cmd: repeatable PowerShell 7 batch entry, argument forwarding and exit-code propagation; no execution-policy bypass.
- C:\PATH\update-ffmpeg.ps1: fixed BtbN GitHub latest win64 GPL asset, SHA-256 manifest plus asset digest when supplied, explicit HTTPS/redirect hosts, finite download cancellation, exact ZIP pair extraction, staged version probes, matching pair version, exclusive lock, backup/rollback and provenance metadata.
- C:\PATH\ffmpeg.exe and ffprobe.exe installed; runtime-owned C:\PATH\.ffmpeg-update\ contains lock and installed.json. Archive/staging are cleaned; no ffplay requested or installed.
- Task-local ffmpeg-updater.tests.ps1 contains 16 offline fixture/batch checks, including full-flow success and mismatch preservation, successful pair replacement, rollback, checksum/ZIP rejection, non-file destinations and real batch exit codes. Disposable fixtures are cleaned.
- Earlier registration changes retained: ~/.harness/agents/pi/{worker-easy,worker-medium,worker-hard}.md; ~/.pi/agent/settings.json scan root; normal CODE host branch annotations in ~/.harness/AGENTS.md, agents/SUBAGENTS.lrf, agents/AGENT_ROUTER.md, agents/worker.md and rules/engineering.lrf. Crabcode/OpenCode use openai; Pi uses openai-codex, identical model IDs/effort. Astra/Attacker routes unchanged.
- User-approved ~/.pi/package.json overwrite changed only name to pi-local-runtime-layout, private=true and descriptive metadata. Old SDK identity was false: its ./dist did not exist. Other legacy metadata/assets/dependencies were preserved. No fake dist junction, runtime package patch or dependency installation.
- Shared local evaluation/subagents_model_eval.json updated once after all runs with eight correlated records, preserving routing priors and distinguishing infrastructure failures from model quality. Evaluation JSON parsed successfully afterward; no database eval submission authorized or performed.

## Native runtime evidence
Root-shadow defect reproduced locally: actual subagent runner JS paths caused SDK-root selection to accept ancestor ~/.pi/package.json (0.85.1), overriding configured valid 1.4.0. After metadata identity correction, both pure resolver cases select the real managed 1.4.0 SDK and existing CLI.
Completed native runs: smoke 35eec871-574f-480d-9d07-00fd13ea26d6; initial CODE d2bb6ad0-e134-4c05-8778-c04bf3ed6f2d; corrections b080b1c9-35c7-42ea-8b5d-852a94099eaf, e39528e2-2743-4194-bc54-db9f0879b318, 8ff60d3c-3d43-4008-a1ca-9ab1ef07fe0b. Parent inspected each saved assistant-message stream: actual provider/model openai-codex/gpt-6-luna; runtime thinking medium, complete, exit0. Root startup failures had no model result; never counted as verified model execution.

## Verification
- Main inspected actual source and requested three bounded same-worker corrections; no Main code edits. Findings corrected: ffprobe banner rejection, unbounded download-body copy and response-object leakage, target type guard, false-positive negative-only tests, and CMD echo parentheses.
- Main final offline fixtures: PASS 16, exit0. Guard accepts normal target paths; full offline success checks checksum, ZIP, install and metadata. No fake executable execution.
- Main real CMD -CheckOnly: latest BtbN asset resolved, exit0.
- Main real CMD first installation: exit0; downloaded archive SHA-256 de0842de7770c756e1d66abc4dc5fdf1b99bffb5f1d1b2510f7b8bc54c59f0c2, validated before execution. Installed metadata timestamp 2026-10-25T06:15:35.0573634Z.
- Actual installed build for both tools: N-128013-g6b05263b44-20261024. Absolute-path -version checks each exit0. Clean pwsh command resolution selects C:\PATH\ffmpeg.exe and C:\PATH\ffprobe.exe.
- Live no-file media pipeline: ffmpeg synthetic audio to ffprobe stdin, successful JSON pcm_s16le / 8000 Hz / mono.
- One parent version probe initially used Select-Object -First 1 and closed native stdout early; corrected test captured complete output and both native exit codes were0. This was test-harness behavior, not a product defect.
- Git diff/staging unavailable because C:\PATH and ~/.harness are not Git worktrees. No stage/commit/push or unrelated app/profile/PATH modifications. Actual source inspections and targeted edits provide change evidence instead.

## Use and remaining limitations
Run C:\PATH\update-ffmpeg.cmd to fetch/update latest build; append -CheckOnly to inspect available build without installation. PowerShell 7 and network access required. Each actual update downloads the ~200 MB archive again; same-build skipping was optional and not implemented. Download prints two harmless System.Threading.Tasks.VoidTaskResult lines; cosmetic only and no signed URLs emitted. SHA-256 from the same publisher proves integrity/consistency, not independent signature authenticity; upstream account/build compromise remains a trust risk. Main conducted separate source/security checks, but no independent Attacker or Verifier was run. Only worker-medium was live-tested; easy/hard profile discovery/configuration were checked, not their model execution. Fresh-shell operation without the retained environment overrides was not separately tested.

## Reporting
Task reports are ingested once each through the configured existing llm-mem.exe; exact receipts are stored separately in memory-receipts.json to avoid self-referential report re-ingestion. No DB eval, external mirror or credentials. Tags: native-subagent, Pi, FFmpeg, updater, Windows, verified-install.
