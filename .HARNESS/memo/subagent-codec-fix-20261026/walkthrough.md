# Walkthrough — subagent-codec-fix-20261026

## Overview
Updated exactly two subagent role bodies found by the previous audit to contain the retired tokenizer optimizer.

## Change
`compressor.md`: PG objective replaces fixed languages/token counts/uniqueness/NFKC tie-break. `filereader.md`: same; language lookup and output protocol remain exact. No authority/effect changes.

## Verification
Preflight PASS; obsolete signature scan PASS; frontmatter/bootstrap/schema/handoff invariants PASS; UTF-8 and diff checks PASS.
