# Knowledge Base: Absolute-Path Cleanup Request

## Context Findings
- The request targets absolute filesystem paths in the current harness repository.
- Maintained agent prompts, LRF guidance, skills, scripts, configuration, and historical evaluation artifacts contain distinct classes of path references.
- Some entries in `rules/default.rules` encode exact command strings; changing path fragments can affect rule matching.

## Morphism
- Input: request to convert absolute paths to relative paths or `$env:USERPROFILE`-based forms.
- Operation: inventory occurrences, classify operational versus explanatory/history data, then apply a bounded transformation.
- Result: implementation paused before edits because the applicable change sequence and handoff contract conflict on whether a Japanese retransformation and user confirmation are required.

## Uncertainty
- User direction is needed to resolve the workflow conflict before CODE can begin.
