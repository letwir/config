---
name: agy-subagent-router
description: Route a bounded external coding task through agy.exe with current model discovery, bounded permissions, reasoning effort, and compact structured output.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# AGY subagent router

## Trigger and authority

Use this procedure only when the current user explicitly asks to route a bounded CODE task through AGY and `AGENT_ROUTER.md` selects `worker-agy`. Invoking this skill alone is not task permission; the current explicit AGY delegation request is required. AGY is opt-in, not a default or automatic fallback. This skill describes procedure; it grants no authority, permission, or effects. Preserve the user's exact task boundary. Do not recurse, delegate children, perform external writes, or make VCS changes.

## Preflight

1. Read [`HARDWARE.md`](../../HARDWARE.md) relative to HarnessRoot (never relative to the process CWD) and select the shell/tool syntax from its current `active_profile`.
2. Resolve `agy` / `agy.exe` according to that profile and the canonical `agy` resolution rule in [`rules/agy.lrf`](../../rules/agy.lrf). Stop if no executable resolves; do not substitute another executable or stale path.
3. Read the applicable [`rules/agy.lrf`](../../rules/agy.lrf), [`agents/AGENT_ROUTER.md`](../../agents/AGENT_ROUTER.md), [`agents/worker-agy.md`](../../agents/worker-agy.md), and [`agents/SUBAGENTS.lrf`](../../agents/SUBAGENTS.lrf) requirements. Discover inventory by running the resolved executable's `models` subcommand before choosing a model. If it exits nonzero, produces invalid/unparseable inventory, or produces no model identifiers, STOP before delegation. Select only an exact identifier in the current output; do not invent or reuse a stale ID.
4. Confirm the request is bounded and the mode/effects are authorized. Use `plan` for design or uncertainty. Use `accept-edits` only for specifically requested bounded workspace edits. Inspection is read-only: use a read-only prompt, not edit mode.
5. Use the sandbox when compatible with the task, and add only the exact workspace directory required. Never use `--dangerously-skip-permissions`.

## Handoff

Construct the exact wrapper and field order below. Replace placeholders with task-specific verified facts only; do not infer permission, evidence, scope, or acceptance. Keep the target, acceptance, and scope concrete and bounded. The final payload line is exactly `SIGMA/1`, immediately before `</Γ>`.

```text
<Γ>
SIGMA/1: PIDGEN/text
Role: worker-agy; difficulty=<easy|medium|hard>
Target: <one bounded CODE target>
Acceptance: <observable acceptance and checks>
Scope: <exact allowlist and effect constraints>
Known facts: <verified facts only, or N_A>
SIGMA/1
</Γ>
```

## Invocation examples

These are profile-specific templates, not authorization to dispatch. Before use, inspect the active profile and resolve the executable. The `models` command must succeed and include at least one parseable `MODEL_ID<TAB>DISPLAY_NAME` row with nonempty fields; otherwise stop before delegation. Display the actual model inventory, inspect it, and manually choose an exact ID from a valid row. Substitute that chosen ID in the separate delegation command; do not infer or automatically select an ID. Also substitute only the approved workspace and constructed handoff. Installed CLI effort values are `low`, `medium`, `high`, and `max` (not `med`). Choose effort appropriate to the task and current routing rules.

Windows PowerShell (`windows-pwsh`):

```powershell
# Resolve in canonical candidate order; stop if none exists. $agy is the resolved executable.
$agy = $null
foreach ($candidate in @('agy', 'agy.exe', "$HOME/bin/agy", "$env:USERPROFILE/bin/agy.exe")) {
    if ($candidate -in @('agy', 'agy.exe')) {
        $command = Get-Command $candidate -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($command) { $agy = $command.Source }
    } elseif (Test-Path -LiteralPath $candidate -PathType Leaf) {
        $agy = (Resolve-Path -LiteralPath $candidate).Path
    }
    if ($agy) { break }
}
if (-not $agy) { throw 'AGY executable unresolved; stopping.' }
$models = & $agy models
if ($LASTEXITCODE -ne 0 -or -not $models -or -not ($models | Where-Object { $_ -match '\S' })) { throw 'Model inventory failed, is invalid, or is empty; stopping before delegation.' }
$validModels = @($models | Where-Object { $_ -match '^[^\t\r\n]+\t[^\t\r\n]+$' })
if ($validModels.Count -eq 0) { throw 'Model inventory contains no parseable MODEL_ID<TAB>DISPLAY_NAME rows; stopping before delegation.' }
Write-Output 'Current model inventory (inspect and choose one exact ID before continuing):'
$models | ForEach-Object { Write-Output $_ }
# After inspecting the displayed inventory, replace the placeholder below with the chosen exact ID.
$p = @'
<Γ>
SIGMA/1: PIDGEN/text
Role: worker-agy; difficulty=<easy|medium|hard>
Target: <one bounded CODE target>
Acceptance: <observable acceptance and checks>
Scope: <exact allowlist and effect constraints>
Known facts: <verified facts only, or N_A>
SIGMA/1
</Γ>
'@
& $agy --print $p --model '<current-model-id>' --effort medium --mode plan --sandbox --add-dir '<exact-workspace>' --output-format json
```

POSIX shell (`android-termux-bash` or `linux-native-bash`):

```sh
# Resolve first existing candidate in canonical order; stop if none exists.
agy_bin=
for candidate in agy agy.exe "$HOME/bin/agy"; do
  if [ "$candidate" = agy ] || [ "$candidate" = agy.exe ]; then
    resolved=$(command -v "$candidate" 2>/dev/null) || resolved=
  elif [ -f "$candidate" ]; then
    resolved=$candidate
  else
    resolved=
  fi
  if [ -n "$resolved" ]; then agy_bin=$resolved; break; fi
done
if [ -z "$agy_bin" ]; then echo 'AGY executable unresolved; stopping.' >&2; exit 1; fi
models_output=$("$agy_bin" models)
models_status=$?
if [ "$models_status" -ne 0 ]; then
  echo 'Model inventory command failed; stopping before delegation.' >&2
  exit 1
fi
valid_models=$(printf '%s\n' "$models_output" | awk -F '\t' 'NF == 2 && $1 != "" && $2 != "" { print; found = 1 } END { if (!found) exit 1 }') || {
  echo 'Model inventory contains no parseable MODEL_ID<TAB>DISPLAY_NAME rows; stopping before delegation.' >&2
  exit 1
}
printf '%s\n' 'Current model inventory (inspect and choose one exact ID before continuing):'
printf '%s\n' "$valid_models"
# After inspecting the displayed inventory, replace the placeholder below with the chosen exact ID.
p=$(cat <<'EOF'
<Γ>
SIGMA/1: PIDGEN/text
Role: worker-agy; difficulty=<easy|medium|hard>
Target: <one bounded CODE target>
Acceptance: <observable acceptance and checks>
Scope: <exact allowlist and effect constraints>
Known facts: <verified facts only, or N_A>
SIGMA/1
</Γ>
EOF
)
"$agy_bin" --print "$p" --model '<current-model-id>' --effort medium --mode plan --sandbox --add-dir '<exact-workspace>' --output-format json
```

For explicitly requested workspace edits, change `--mode plan` to `--mode accept-edits`; do not otherwise broaden scope. Use `--mode plan` for planning. Read-only inspection must not use an edit mode. Apply only options supported by the installed CLI; these examples use its `--print`, `--model`, `--effort`, `--mode`, `--sandbox`, `--add-dir`, and `--output-format` flags. Do not assume JSON output proves success.

### Authentication

If authentication is required, use an unsandboxed invocation only after explicit approval for this task. Pause for the user to complete browser sign-in or grant access. Never receive, request, display, paste, or transmit credentials, tokens, secrets, or private environment values. Authentication, policy, approval, credential, or permission blocks are terminal: report FAILED and STOP; do not retry or fall back around them.

## Failure and fallback

- Retry once only for a pre-launch startup failure (including executable/OS/CLI startup failure). Do not retry an invocation that has launched.
- For timeout or any other post-launch operational failure, stop or confirm the process has exited, then inspect the real partial workspace diff before the parent considers any fallback. Preserve unrelated/existing edits.
- After an eligible operational failure, the parent/main agent alone may select one currently available GPT-family worker for the same bounded task and scope, using AGENT_ROUTER's difficulty mapping. This is serial; never run AGY and fallback concurrently. Never select Astra as an AGY-failure fallback. No fallback for authentication, policy, approval, credential, or permission blocks.
- Malformed result or generation/backend failure is operational failure: stop/confirm exit and inspect partial diff before parent decision. Do not claim success from JSON structure or a status field alone.

## STOP and success

STOP and report the blocker on missing executable/model, unresolved profile or permissions, invalid handoff, prohibited effect, terminal block, or unconfirmed process state. Do not dispatch on uncertainty.

Success requires confirmed process completion, parent/main inspection of the actual diff, proportionate deterministic checks with observed exit results, and reconciliation against every acceptance condition. The worker returns compact evidence: observed changed files, checks and exit results, failures, and residual risk. Report the effective model only when verified from actual invocation evidence. JSON is structured output, not proof of success.

The child may return evaluation fields as evidence, but must not write the shared evaluation artifact. After wait-all, the parent/main alone updates `~/.harness/evaluation/subagents_model_eval.json` as required by [`agents/SUBAGENTS.lrf`](../../agents/SUBAGENTS.lrf).

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=agy-subagent-router
R|bin.resolve|task:agy|MUST|RO_LOCAL|hardware=../../HARDWARE.md; resolve=active_profile-in-HARDWARE.md; bin=first-existing(agy,agy.exe,$HOME/bin/agy,$env:USERPROFILE/bin/agy.exe); missing=>STOP
R|use.inspect|task:agy|MUST|RO_LOCAL|cmds=models; effect=read-only-inspection
R|use.delegate|task:agy|MUST|LIVE_WRITE|cmds=delegate,exec; pre=explicit-user-request-for-external-delegation; approval=current-task
R|guard.perm|task:agy|MUST_NOT|LW_SCOPE|flag=--dangerously-skip-permissions; external-authority-expansion=forbidden
R|guard.cred|task:agy|MUST_NOT|CRED|deny=credentials,tokens,raw-secret-logs; allow=task-prompts,model-ids
R|guard.truth|task:agy|MUST|RO_LOCAL|unverified=fact-forbidden; main-agent-diff-inspection=mandatory