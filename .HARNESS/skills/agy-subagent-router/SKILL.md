---
name: agy-subagent-router
description: Route a bounded external coding task through agy.exe with current model discovery, bounded permissions, reasoning effort, and compact structured output.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

§ROUTE
Model selection ⇒ first discover current inventory; use returned identifiers only; absent/stale identifier ⇒ ¬select.
External handoff/implementation ⇒ bounded Task⊗Model; explicit current-task approval required before delegate/exec (LRF governs).
Plan ⇒ plan-mode; requested scoped edits ⇒ accept-edits; inspections ⇒ read-only. Preserve permissions; ¬permission-skip.
```bash
# HARDWARE.md active_profile: android-termux-bash (windows-pwsh 時は $p = @"..." を使用)
p=$(cat <<'EOF'
<Γ>
SIGMA/1: PIDGEN/text
Role: <role>
Target: <target>
Acceptance: <acc>
Scope: <scope>
Known facts: <facts>
</Γ>
EOF
)
agy --print="$p" --model <id> --effort <low|med|high> --mode plan --output-format json
```
⚡ Exec ⇒ requested Task⊗Files within exclusive bounded scope; preserve unrelated edits.
📊 Result evaluation ⇒ verify structured RunResult; record SUCCESS∨FAILED∨ADVICE∨ESCALATE; ¬claim unverified success.
Failure ⇒ STOP/confirm process exit; inspect partial diff before permitted fallback; auth/policy/approval/credential/permission block ⇒ FAILED+report+STOP; ¬concurrent fallback.
Credentials/secrets/private values ⇒ ¬request∧¬print∧¬paste∧¬transmit.
VCS ⇒ ¬stage∧¬commit∧¬push∧¬rewrite-history∧¬branch-alter.
```bash
agy --print="$p" --model <id> --mode accept-edits --output-format json
```
Eval(RunResult) ⇒ verify∧record(subagents_model_eval.json); outcome∈{SUCCESS,FAILED,ADVICE,ESCALATE}.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=agy-subagent-router
R|bin.resolve|task:agy|MUST|RO_LOCAL|hardware=../../HARDWARE.md; resolve=active_profile-in-HARDWARE.md; bin=first-existing(agy,agy.exe,$HOME/bin/agy,$env:USERPROFILE/bin/agy.exe); missing=>STOP
R|use.inspect|task:agy|MUST|RO_LOCAL|cmds=models; effect=read-only-inspection
R|use.delegate|task:agy|MUST|LIVE_WRITE|cmds=delegate,exec; pre=explicit-user-request-for-external-delegation; approval=current-task
R|guard.perm|task:agy|MUST_NOT|LW_SCOPE|flag=--dangerously-skip-permissions; external-authority-expansion=forbidden
R|guard.cred|task:agy|MUST_NOT|CRED|deny=credentials,tokens,raw-secret-logs; allow=task-prompts,model-ids
R|guard.truth|task:agy|MUST|RO_LOCAL|unverified=fact-forbidden; main-agent-diff-inspection=mandatory
