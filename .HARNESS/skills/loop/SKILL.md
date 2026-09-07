---
name: loop
description: Run a prompt or skill on a recurring local interval using monitored background shell output or schedule tool.
---
[SPR/PIDGIN::ρ→max] ⏱️fixed ⏳dynamic 🛑stop
⏱️ 定期 fixed ⊢ Interval⊗Prompt ⇒ 固定間隔定期実行ループ | `schedule(CronExpression="*/5 * * * *", Prompt="...")` ∨ `while($true){sleep <sec>; ...}`
⏳ 自律 dynamic ⊢ Event∨Heartbeat ⇒ イベント駆動・自律ペース監視ループ | `schedule(DurationSeconds=<sec>, Prompt="...")`
🛑 停止 stop ⊢ TaskID ⇒ 実行中ループタスク・プロセスの終了 | `manage_task(Action="kill", TaskId="...")`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=loop
R|use.schedule|task:loop|MUST|RO_LOCAL|tool=schedule; mode=DurationSeconds-or-CronExpression; sleep-in-background=forbidden
R|use.inspect|task:loop|MUST|RO_LOCAL|cmds=fixed,dynamic; effect=read-only-inspection; one-off-task=omit
R|use.kill|task:loop&case:stop|MUST|LW_SCOPE|action=kill-task; target=active-loop-id
R|guard.truth|task:loop|MUST|RO_LOCAL|unverified=fact-forbidden; report-each-tick=minimal

## Guidance

- Title shell commands as `Loop <schedule>: <prompt>` (e.g. `Loop every 5m: check deploy status`).
- Use a unique sentinel per loop so unrelated output does not trigger notifications.
- Avoid noisy commands inside the loop.
- Do not create duplicate fixed loops or dynamic sleepers.
- If the user asks to stop, kill any tracked loop/sleeper PID, then await the shell task so its completion notification is consumed and does not wake the agent later. Do not schedule another dynamic wake.
