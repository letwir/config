# walkthrough | 20260925-gemini-harness-dir-links

## Plan
Enumerate `.harness` top-level directories, check same-name entries in `.gemini`, and validate link destinations. Create links only for missing entries; preserve Gemini-only state.

## Evidence and result
- Checked `.harness` directories: `agents`, `etl`, `evaluation`, `memo`, `persona`, `pet-runs`, `pets`, `rules`, `scripts`, and `skills` (excluding `.git`). Every name already exists in `.gemini`.
- The first nine are symbolic links to `..\.harness\<name>`; `skills` is a junction resolving to `C:\Users\letwir\.harness\skills`.
- `.gemini\GEMINI.md` and `.gemini\AGENTS.md` already target `.harness\AGENTS.md`.
- `.gemini` contains distinct runtime/configuration directories without `.harness` counterparts. Keep them in place.

## Completion
No symlinks were added or modified because none were missing or invalid. No tests were needed for this read-only link audit.
