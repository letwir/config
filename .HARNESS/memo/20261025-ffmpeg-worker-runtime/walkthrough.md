# Walkthrough: 20261025-ffmpeg-worker-runtime

## Overview
User asked to resume FFmpeg updater implementation after approving host-specific worker provider routes. No FFmpeg implementation began: the authorized native worker registration smoke run failed before model execution.

## Deliverables and retained state
Prior local configuration edits are retained: ~/.pi/agent/settings.json adds ~/.harness/agents/pi as a scan root; worker-easy/medium/hard definitions exist there; shared worker.md, AGENT_ROUTER.md, SUBAGENTS.lrf, rules/engineering.lrf and AGENTS.md explain ordinary CODE provider branches Crabcode/OpenCode=openai and Pi=openai-codex. Astra/Attacker routing unchanged. These are configuration/document edits, not FFmpeg implementation.

## Verification
- Pi discovery and get report all three profiles with intended model/thinking; JSON, six-field LRF record, UTF-8/text format and defaultReads-path checks passed.
- Git status/diff unavailable: ~/.harness is not a Git repository; edits were exact targeted replacements plus three new profile files. No commits or branch changes.
- Native smoke run 6748198b-a8b5-46c8-84b1-e5ba66a01c21: FAILED before child model result, `ENOENT: no such file or directory, lstat 'C:\Users\letwir\.pi\dist'`. Status results count 0. process-terminal.json records state observed. No effective model verified and no child result accepted.
- Host package root C:\Users\letwir\.pi\agent\install\releases\1.4.0\node_modules\@earendil-works\pi-coding-agent has expected package name/version and an existing dist directory.
- Supported env override PI_SUBAGENTS_PI_CODING_AGENT_PACKAGE_ROOT is unset in available shell environment. Proposed next step is restart Pi with this exact validated package root set, then retry the same native read-only smoke protocol. This remedy is untested; do not claim success or silently switch execution mode.
- No update-ffmpeg.cmd/.ps1, ffmpeg.exe or ffprobe.exe existed in C:\PATH at start. No FFmpeg target files created or modified; no installation executed.
- Public read-only primary evidence 2026-10-25: GitHub API BtbN/FFmpeg-Builds latest release tag latest, published 2026-10-24T13:21:42Z. Assets include checksums.sha256 and ffmpeg-master-latest-win64-gpl.zip (200487045 bytes). No asset downloaded or executed. Planned updater uses checksums and staged replace with rollback; design/implementation not finalized.

## Limitations and attribution
Native worker path infrastructure blocks implementation under the standing worker-only CODE rule. No fallback, auth/credential changes, live binary execution, independent review, or model upgrade. Parent over-spent effort on policy navigation; no defensible defect in the user's current request. Memory ingest permitted by standing report directive; separate evaluation live write not performed. Tags: pi, subagents, FFmpeg, Windows, blocked-runtime.
