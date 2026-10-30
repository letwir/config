### 2026-10-26 — Pi configuration symlink distribution

#### Hypothesis
A small per-entry symlink installer can publish reusable Pi settings without displacing user-owned extensions or including runtime/secret material.

#### Tried
Inspected Pi settings/resource path behavior and existing repo installation conventions; implemented Bash/PowerShell installers; tested fresh installs, reruns, dry runs, destination conflicts, backups (including dangling links), missing sources, overlaps, symlink parents, and rollback against temporary directories.

#### Rejected
Replacing the complete `~/.pi/agent` tree or using the `.HARNESS` root `install-pi.sh` for a different purpose. Neither preserves independently managed Pi resources safely.

#### Uncertainty
Native Linux and Termux link behavior remains untested. PowerShell LSP was unavailable and other diagnostics timed out.

#### Attribution
Prompt defect 0%: the requested `.HARNESS/.pi` target and two installer shells were clear. Agent defect: initial test failure and link-path omissions were corrected in a follow-up; both hermetic test suites then passed.

#### Completion evidence
Bash syntax and test suites passed; PowerShell real-link test suite passed; JSON parsed. No live user configuration, credentials, package manager, or VCS state was changed.
