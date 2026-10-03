# Walkthrough: EPERM Diagnosis for Agent Router

## Overview
- Date: 2026-09-25
- Task ID: `20260925-eprem-agent-router-diagnosis`
- Request: Explain how to address the Windows `EPERM` encountered when the worker tried to rename its temporary file over `agents/AGENT_ROUTER.md`.
- Outcome: Read-only metadata inspection found a regular file (not a reparse point), no read-only attribute, and inherited ACL summaries showing Modify/FullControl on both target and parent. No obvious read-only or ACL restriction was identified.

## Evidence and Inference
- Observed: `AGENT_ROUTER.md` attributes are `Archive`; `LinkType` and `Target` are empty; `IsReadOnly` is false.
- Observed: Parent directory and target ACL summaries include inherited `Modify` and `FullControl` allow rules.
- Inference: A transient handle lock, security/antivirus filter, controlled-folder protection, or transactional rename behavior is more plausible than a simple read-only attribute. These causes were not independently confirmed.

## Guidance
- Close editors/previews that may hold the target; if needed, inspect open handles with Process Explorer or Resource Monitor.
- Check Windows Security protection history and endpoint-security logs before changing security configuration.
- If a later independently authorized task retries, first validate rename behavior in a disposable same-directory test path and preserve the current target. Do not change ACLs or disable protection without identifying a specific block.
- The failed launched AGY run remains terminal under current routing rules; this diagnosis does not authorize retry or backend switching.
