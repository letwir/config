# Knowledge: FileReader and multilingual codec plan status

- task_id: `20260926-file-reader-codec-plan-block`; timestamp: `2026-09-26T16:53:42+09:00`; target_environment: Windows 11, PowerShell 7, `C:\Users\letwir\.harness`.
- Evidence scope: active LRF, agent routing/prompt files, HSEQ local-research sequence, task-scoped LLM-memory search; verification status: repository comparison completed, implementation blocked before CODE.
- Remaining uncertainty: priority between minimum total tokens and the compression×uniqueness word score is not specified.
- Live repository reads confirm existing FileReader entries and files are present in `agents/AGENT_ROUTER.md`, `agents/README.md`, `agents/SUBAGENTS.lrf`, `agents/filereader.md`, `agents/filereader.toml`, `agents/FILEREADER.css`, and `etl/research-local.seq`; these were preserved and not edited in this task.
- Existing codec copies name `Langᵂ₅` without an explicit set and use a frequency-weighted lexical objective; the user now specifies `{EN,ES,DE,JA,FR}` and prioritizes compression plus uniqueness.
- A remaining specification conflict blocks implementation: the inherited codec minimizes total `tok(y)`, while the new lexeme criterion could make `compression×uniqueness` primary. The priority or tie-break relation is not yet user-confirmed.
- At stop, the relevant FileReader paths were already modified or untracked in Git; no clean baseline was available. This task itself made no rule/code edits.
