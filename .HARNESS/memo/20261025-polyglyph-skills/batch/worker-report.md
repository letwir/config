RTS_STATUS: FAILED

## Outcome
No target skill files were changed. The required inventory/snapshot manifest could not be safely read through the available command interface, so the exact 39-file allowlist and immutable baseline hashes were not established. I therefore made no edits and did not produce an unverifiable manifest.

## Hydration and blockers
- Read the requested POLYGLYPH codec, BOOTSTRAP.lrf, worker.md, WORKER.css, SUBAGENTS.lrf, and README-forLLM.md.
- Routed documentation lookup requires loading matching MANUAL/LOAD records; rg queries returned worker routing records, but the attempt to inspect the manifest using pwsh failed because the command was parsed/interpreted by the shell adapter as an invalid command (`/usr/bin/bash.PSObject.Properties.Name`). A prior rg command also failed at shell quote parsing.
- Tool output displayed the Unicode codec and policy material as mojibake, preventing byte-faithful reuse/copying of the canonical codec from observed output.
- As a result, neither candidate-body inspections nor comparisons, preservation merges, label repairs, synchronization, or the new root codec copy were safely completed.

## Evidence
- Commands attempted: bounded `rg.exe` lookups for MANUAL/LOAD (returned worker records); manifest inspection via pwsh (failed); additional combined rg query (failed at shell quote parsing).
- Changed files: only this required report.
- Tests/checks: not run; parent-owned deterministic validation remains pending.
- No staged files: no VCS commands were run.

## Residual risk
All requested migration work remains outstanding. No claims are made regarding the 39 targets, local-only capability preservation, UTF-8/no-BOM, frontmatter/fence/LRF preservation, or root POLYGLYPH.md parity. A functioning pwsh invocation/file reading route and task-scoped byte-preserving review are needed before safe completion.
