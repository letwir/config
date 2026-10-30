---
name: loop
description: Run a prompt or skill on a recurring local interval using monitored background shell output or schedule tool.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] ⏱️fixed ⏳dynamic 🛑stop
🔎 User requests recurring execution of a prompt/skill ⊢ choose exactly one cadence: fixed interval (Interval⊗Prompt) or dynamic/event-driven (Event∨Heartbeat; autonomous-paced monitoring). Fixed example: `schedule(CronExpression="*/5 * * * *", Prompt="...")` ∨ `while($true){sleep <sec>; ...}`. Dynamic example: `schedule(DurationSeconds=<sec>, Prompt="...")`.
⏱️ fixed ⇒ rerun at the specified interval; ⏳ dynamic ⇒ schedule for the specified duration and monitor at the event/heartbeat-driven autonomous pace. Use `schedule`, not a background sleep loop. Inspect only the corresponding `fixed` or `dynamic` monitored output; report each tick minimally and never state an unverified fact.
⛔ Before scheduling: use a unique sentinel per loop so unrelated output cannot trigger notifications; avoid noisy commands; do not create a duplicate fixed loop or dynamic sleeper. Name shell commands `Loop <schedule>: <prompt>` (e.g. `Loop every 5m: check deploy status`).
🛑 User requests stop ⊢ identify the active loop task ID and end that task with `manage_task(Action="kill", TaskId="...")`; if a tracked loop/sleeper PID is active, kill it. Then await the shell task so its completion notification is consumed and cannot wake the agent later. Do not schedule another dynamic wake.
✅ Success = requested loop is scheduled and its corresponding output is monitored; stop = active loop/sleeper is terminated and its shell task completion is awaited. Scheduling does not authorize unrelated effects.

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
