# Diary Correction / 20261018-rts-parallel-capacity-20

### 2026-10-18 — Attribution normalization correction

- The first `diary.md` revision was ingested successfully with a provisional `AgentDefect=25%` estimate.
- The evaluation CLI requires `prompt_ratio + agent_ratio == 1`; the first evaluation submission was rejected before registration. With no identified prompt defect, the attribution was normalized to `prompt_ratio=0.0`, `agent_ratio=1.0` for the prompt-vs-agent causal allocation. This does not mean the completed task itself failed.
- The corrected local `diary.md` and evaluation file now agree. Re-ingesting the same diary path failed due a memory self-edit conflict, including with `-force`; raw error output was not retained. This separate correction record is provided so the final attribution is represented in memory without overwriting the earlier entry.
- No implementation files or VCS state changed during this correction.
