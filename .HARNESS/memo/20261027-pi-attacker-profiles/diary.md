# diary — 20261027-pi-attacker-profiles

- Timestamp: 2026-10-27 JST
- Task: Repair missing Pi attacker-medium registration without lowering review depth.
- Request evidence: User quoted Pi runtime unavailability and asked to fix it.
- Action: Investigated native discovery and configured agentDirs; delegated bounded adapter implementation to worker-medium; main checked files and deterministic regression/preflight; corrected invalid acceptanceRole after actual launch rejected it; medium smoke completed with verified Sol/high.
- Result: Registration defect repaired and medium launch verified. Independent final verifier requested separately; its actual outcome is recorded in walkthrough, not assumed here.
- Friction: Parent first supplied a directory to preflight -Root instead of bootstrap file; reran correctly. Initial handoff mixed generated-profile restrictions with worker capability, causing a supervisor question; parent clarified. First generated acceptanceRole was invalid and lightweight regression missed enum validation. Main LSP caught PowerShell automatic-variable collision. Worker automatic acceptance rejected unavailable Git attestation rather than indicating a product-code failure.
- Attribution: AgentDefect, not user ambiguity: erroneous preflight argument, ambiguous worker handoff phrasing, initial schema mismatch, incomplete first regression. The user clearly identified the blocker and requested a repair.
- Ratios: PromptDefect 0%, AgentDefect 100% among observed friction causes; qualitative attribution, not measured population percentages.
- Impact: Additional correction and smoke cycle; no external/VCS writes or PADDOCK changes.
- Feedback/internal correction: Validate native parser enums before claiming registration; distinguish generated-agent tool constraints from implementation-worker tools; report orchestration wrappers honestly instead of literal single-tool guarantees.
- Rewritten request: Repair native Pi attacker-medium registration while preserving Sol/high, read-only constraints, and no automatic lower-depth fallback; verify actual startup.
