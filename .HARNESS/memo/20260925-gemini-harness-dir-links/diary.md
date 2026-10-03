# diary | 20260925-gemini-harness-dir-links

timestamp: 2026-09-25 Asia/Tokyo
task: Compare `.gemini` and `.harness` directories and add links for missing counterparts.
request-evidence: The user asked to link `.gemini` folders to `.harness` when replaceable, then clarified to add links for folders missing from `.gemini`, including `etl`, `agents`, and `rules`.
action: Enumerated top-level directories and inspected link types and destinations. Rechecked all `.harness` directories against `.gemini` and validated the `skills` junction plus the existing `GEMINI.md`/`AGENTS.md` links.
result: All applicable directories were already linked; no additional links or repairs were needed. Gemini-only config/runtime/history/temp directories remain distinct.
friction: Initial response concluded no additions were needed before validating each existing destination; the user clarified the requested operation, then the link audit was completed.
attribution: PromptDefect 0%; AgentDefect 100% for premature conclusion, corrected before any configuration change.
impact: No `.gemini` or `.harness` runtime configuration changed.
feedback: The follow-up naming example folders usefully clarified that existing links should be checked, not assumed absent or valid.
rewritten-request: Compare every top-level `.harness` directory with `.gemini`; add a symlink only when the matching `.gemini` entry is missing, and repair an entry only if its target is invalid.
