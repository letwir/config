# Knowledge Base: Absolute-Path Cleanup Completion

## Context Findings
- Replaced user-specific home paths in maintained agent prompts, LRF guidance, default-settings docs, harness sync docs, and the roadmap target reference.
- Chrome discovery now builds standard installation paths from `ProgramFiles` and `ProgramFiles(x86)` environment variables.
- LLM-memory's machine-specific repository fallback now uses the configurable `LLM_MEMORY_REPO` root.
- `rules/default.rules` contains exact-command match rules; changing those literal paths can change rule matching, so it was preserved. Historical evaluation/memo records were also preserved.

## Morphism
- Input: generalize absolute filesystem paths without changing expected behavior.
- Operation: replace maintained personal/fixed path literals with home/environment-derived forms, then validate focused formats.
- Result: implementation completed in the approved files; no personal-home or fixed Program Files literals remain in those targets.

## Remaining Unknowns
- Whether `LLM_MEMORY_REPO` is set in every runtime that uses the fallback.
- Exact-command and historical records retain their original machine-specific path evidence.
