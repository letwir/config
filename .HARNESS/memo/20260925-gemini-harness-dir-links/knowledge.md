# knowledge | 20260925-gemini-harness-dir-links

## Verified
- Every top-level directory in `.harness` except `.git` already has a same-name entry in `.gemini`.
- `agents`, `etl`, `evaluation`, `memo`, `persona`, `pet-runs`, `pets`, `rules`, and `scripts` are symbolic links to the corresponding `.harness` directories. `skills` is a junction to `.harness\skills`.
- `.gemini\GEMINI.md` and `.gemini\AGENTS.md` already link to `.harness\AGENTS.md`.
- `.gemini` also has Gemini/runtime-specific directories such as `config`, `antigravity`, `antigravity-cli`, `history`, and `tmp` that have no `.harness` counterpart.

## Conclusion
No additional `.harness` directory symlinks were needed. No existing links were changed. Do not replace the `.gemini` root with `.harness`; retain Gemini-specific state and configuration.
