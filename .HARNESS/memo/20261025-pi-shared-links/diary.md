# diary: 20261025-pi-shared-links
timestamp: 2026-10-25T15:30:01.0343875+09:00
task: Link missing .harness directories into .pi/agent.
request-evidence: User explicitly requested equivalent links for other folders.
action: Added symbolic links for crabcode, etl, evaluation, memo, opencode, persona, pet-runs, pets, scripts and HARDWARE.md.
result: All 10 new links verified as SymbolicLink and reachable. Required bootstrap, ETL, persona, evaluation and hardware references resolve.
friction: Existing agents is a real local directory; preserved without merging. Existing rules and skills retained.
attribution: No defensible prompt defect; PromptDefect 0%. No independent reviewer used for this filesystem-only operation.
impact: Shared reference discovery restored; no software implementation or credential/config edits.
feedback: No prompt clarification required.
rewritten-request: N_A
residual-risk: Existing agents content not synchronized. agy.exe not tested.
public-research: N_A local filesystem repair.
