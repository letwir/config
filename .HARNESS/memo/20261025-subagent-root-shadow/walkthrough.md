# Walkthrough: 20261025-subagent-root-shadow

## Overview
User clarified that native subagent operability is itself the purpose of the work, not merely a means to FFmpeg installation. Main acknowledged misprioritization. FFmpeg remains paused; no alternate CLI, model or Main implementation fallback.

## Findings and evidence
- Failed native run 435b05f8-e773-411d-87c6-ef34c7888d7a is terminal/failed with observed process proof, zero child results, error ENOENT lstat C:\Users\letwir\.pi\dist; resume unavailable because no child session persisted.
- C:\Users\letwir\.pi\package.json declares @earendil-works/pi-coding-agent version 0.85.1 and main/root import ./dist/index.js. Its dist directory does not exist.
- Valid managed SDK under ~/.pi/agent/install/releases/1.4.0/node_modules/@earendil-works/pi-coding-agent declares version 1.4.0 and has dist.
- Pi settings load pi-subagents 0.82.0 from ~/.pi/agent/npm/node_modules/pi-subagents/bundle/index.js. The relevant bundle and source logic agree.
- Root cause locally reproduced without any worker/model/network call: import src/runs/shared/pi-spawn.js resolveHostPiPackageRoot and pass argv1 as each of the real subagent-runner-bootstrap.js and subagent-runner.js paths, plus the correct SDK environment override. Both return C:\Users\letwir\.pi, not the configured 1.4.0 root.
- selectPiPackage accepts a JS entry before the explicit override and findPiPackageRootFromEntry walks past the nested pi-subagents package to the ancestor .pi/package.json. This mistakes the runner's ancestor legacy manifest for ownership by the Pi SDK. The missing ./dist export then matches the observed error.
- Get-Command pi -All in the tool shell resolves managed agent/bin/pi.ps1, pi.cmd and pi launcher; no need to ask another blind restart. Correct environment variables alone cannot fix the reproduced ancestor-selection precedence.

## Proposed bounded repair and constraints
Preferred config-only diagnostic repair: retain the exact old ~/.pi/package.json bytes under a clearly named backup and reversibly remove that legacy manifest from SDK discovery, then rerun the same pure resolver test before any native child retry. This can affect the legacy standalone 0.85.1 Pi installation, so obtain explicit owner approval for that exact target and effect first. Do not delete legacy assets, create a fake dist junction, edit package runtime source, move dependencies, change models or bypass the subagent protocol. Alternative source-code repair would require an explicit Main-only bootstrap exception because the required CODE worker cannot currently start.

## Completion and attribution
No product/config/runtime edits in this diagnostic pass. Native model execution remains unverified. The local root-selection defect is reproduced, but no repair or successful child run claimed. Agent fault: previous environment-only advice and repeated FFmpeg/Main fallback proposals missed the operability acceptance criterion; user instructions were clear. Standing memory report ingestion only; no evaluation write, external mirror, package update or credential changes. Tags: pi, subagents, SDK-root-shadow, local-repro, blocked-pending-config-decision.
