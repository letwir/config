# Knowledge / 20261018-polygllyph-application-resume

- topic: Resumed POLYGLYPH format application and structural lint
- verified-date: 2026-10-18
- scope: `.HARNESS` LLM instruction/contract Markdown
- facts:
  1. The exact target manifest contains 59 harness-relative paths: root AGENTS, 16 agent prompts, 40 SKILL.md files excluding node_modules, the LLM navigation document, and one OpenCode agent prompt.
  2. The fixed source baseline is Git revision `6a8889674bd01d867e9374430dbab46ebd3115b8`. Main independently compared all frozen frontmatter/fence/LRF-block signatures and verified 59/59.
  3. The linter models the real embedded `#SIGMA LRF/1` block in `skills/loop/SKILL.md`; its block is followed by a blank separator and further PG prose. The parser preserves the block and resumes PG/fence parsing after it. EOF trailing blanks are retained.
  4. All 59 current targets carry valid metadata and parse structurally; only `skills/rikei-search/SKILL.md` had unchanged body content because its existing body was already a compact `⊢/⇒` science-query mapping.
  5. The LRF modality set includes MUST, MUST_NOT, and MAY; symbols/emoji are never sole safety or authority carriers. PG does not replace SIGMA/LRF authority.
  6. Representative semantic review passed for worker handoff, read-only Google search, LIVE_WRITE llm-memory, plugin creator with code/config exceptions, and OpenCode orchestrator. This review is not a formal proof or cross-model guarantee.
  7. The immediately preceding RTS task already raised admission to 20 descendants per Runtime; this resumed PG task preserved those capacity changes.
- source: repository contents and test outputs verified 2026-10-18
