### 2026-09-25 01:01:47 +09:00
- task: Align the installed harness-lint SKILL.md HSEQ addition with the existing notation.
- request-evidence: User asked to install under `.harness/skills` and add the HSEQ instructions using the current SKILL.md notation.
- action: Confirmed the CLI binary is already installed; replaced prose bullets with PIDGIN-style dialect, structural check, command, and verdict cues; updated the skill description.
- result: Single-file strict lint of `etl/main.seq` and the scoped diff check both exited 0.
- friction: None. Precedent search returned empty for this exact documentation follow-up.
- attribution: PromptDefect=0%; AgentDefect=0%. The target and requested notation were clear; no failure occurred.
- impact: LLM-facing skill documentation now exposes HSEQ invocation and result semantics in the existing notation.
- feedback: None.
- rewritten-request: Document installed HSEQ lint in the harness-lint skill using its established PIDGIN command and result format.
