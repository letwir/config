# Knowledge Base

## Context
`.codex\AGENTS.md` had independent contents while `.harness\AGENTS.md` is the intended single source of policy.

## Findings
On 2026-09-25, `.codex\AGENTS.md` was replaced with a Windows symbolic link to `C:\Users\letwir\.harness\AGENTS.md`. Its former bytes were preserved at `C:\Users\letwir\.codex\AGENTS.md.pre-harness-symlink-20260925.bak` (SHA-256 `92265534926B97959E322F71A803324CE11CC5AA4F38D79E9DC768BE4058A342`). The link target resolved correctly after creation.

## Morphism
For future policy updates, edit `.harness\AGENTS.md`; the `.codex\AGENTS.md` link resolves to that source.
