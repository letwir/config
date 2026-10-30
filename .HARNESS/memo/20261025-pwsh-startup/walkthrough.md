# walkthrough: 20261025-pwsh-startup

- Request: Speed up the local PowerShell profile while preserving existing features.
- Target: A:\Users\letwir\Documents\PowerShell\profile.ps1.
- PowerShell: 7.6.7. Noninteractive process timing, three runs per mode, same parent environment.
- NoProfile: 1159, 1116, 947 ms; median 1116 ms.
- Profile: 4607, 3816, 3844 ms; median 3844 ms. Observed median difference 2728 ms; not an interactive first-prompt measurement.
- Separate fresh-process lookup probe: carapace.exe 586.4 ms; absent 7z/jtrim/sakura/ffmpeg/ffprobe/nkf lookups 192.3/226.1/212.0/204.1/209.0/194.6 ms. Probes are not exclusive attribution of full startup time.
- Fact: Profile eagerly resolves twelve tool aliases with Get-Command and initializes cached Carapace. Starship sibling contains comments only.
- Inference: Repeated missing-executable searches and initial command discovery are optimization candidates. Compare exact application-only resolution with module autoload disabled in a narrow scope, or remove redundant aliases where ordinary executable lookup already preserves behavior; preserve JAVA_HOME fallback and existing aliases. No implementation approved or applied.
- Product files changed: none. Existing startup cache may be touched by normal profile execution during checks.
- Checks: All six measured child processes exited 0; no raw startup output emitted. Interactive completion, prompt rendering and VS Code integration remain untested.
- Workflow: Main performed local research only. Current standing rules require a fixed-profile CODE worker or explicit Main-only implementation exception. Ask user for direct-edit exception before implementation in this session. No worker execution or independent review claimed.
- Memory precedent search: command exited 0 and returned parseable JSON; results were not reviewed and no precedent claim is made.
- Public/EPUB sources: not applicable to local diagnosis. VCS writes: none.
- Evaluation submission: blocked; standing task-memory approval covers report ingest, not a separately authorized eval live write.
- Friction attribution: agent overexplored workflow configuration; no defensible user prompt defect identified. No performance improvement claimed.
