# H-07 Diary

### 2026-09-09

**Hypothesis:** The existing task-local artifact convention could become the source of truth by adding one bounded index rather than moving or duplicating records.

**Tried:** Compared local diary and llm-memory contracts, defined receipt states, added strict path and identity checks, serialized writers, and tested both PowerShell hosts.

**Rejected:** Keeping `.gemini/diary.md` as the only diary; letting the local writer perform external effects; free-form error text; writing the index before external outcomes were final; silently backfilling H-01 through H-06.

**Uncertainty:** The external diary mirror remains `N_A` unless a task explicitly authorizes it. External receipt shapes remain owned by their respective operations.

**Attribution:** PromptDefect not found. AgentDefect found and corrected: the first writer had incomplete tests, open reason strings, and incorrect FINISH receipt ordering.

**Search correction:** Public guidance did not define this local convention; the implementation therefore follows verified repository contracts and existing task directories.

**Emotion and thoughts:** Separating local durability from external success removed the temptation to hide partial failure behind one completion flag.