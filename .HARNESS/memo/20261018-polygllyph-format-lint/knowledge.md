# Knowledge / 20261018-polygllyph-format-lint

- topic: POLYGLYPH migration and static lint attempt
- verified-date: 2026-10-18
- scope: repository `.HARNESS` LLM instruction/contract Markdown
- facts:
  1. The fixed inventory drafted for this task contains 59 tracked files at baseline Git revision `6a8889674bd01d867e9374430dbab46ebd3115b8`: 1 root AGENTS file, 16 agent Markdown files excluding README, 40 SKILL.md files excluding node_modules, 1 LLM navigation doc, and 1 OpenCode agent prompt.
  2. The user directly requested the format migration and new lint, supplied `.HARNESS/POLYGLYPH.md`, and explicitly permitted subagent work. The router limits simultaneous workers to six despite the user's allowance of twenty.
  3. A partial worker change added `MAY` to the modality set and inserted §8 decoder/embedding text, but §8 currently appears before §7; it is not fully validated.
  4. Two untracked fixture files exist with misspelled `polyglglyph` paths. Their target list contains 59 Git-root-prefixed paths; preservation fence hashes were independently found to mismatch for 8 baseline files. They are not a valid verified baseline.
  5. No linter, migration of AGENTS/README-forLLM or agent/skill prompts, semantic review, or tests were completed.
  6. The worktree also contained pre-existing changes to the llm-memory binary and its untracked backup, plus reports from an earlier interrupted request; these were preserved.
- source: local Git status and read-only baseline signature comparison, 2026-10-18
