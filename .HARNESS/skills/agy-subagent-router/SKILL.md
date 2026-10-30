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

🔎 適用条件=current userがbounded CODEのAGY委託を明示要求 ∧ `AGENT_ROUTER.md`選択=`worker-agy`。Skill呼出≠task permission；current-task explicit AGY delegation request MUST。AGY=opt-in；既定選択・automatic fallbackは禁止。本文=手順記述のみ；authority/permission/effectsを付与しない。user task boundaryを厳密保持；recursion・child delegation・external writes・VCS changesは禁止。

Obj(bounded TaskContext∧verified Facts) → Mor(profile解決→model inventory→authorized handoff→bounded invocation→parent検証) → Obj(completion Evidence∪Failure∪ResidualRisk)；表現変換≠権限付与。

## Preflight

1. [`HARDWARE.md`](../../HARDWARE.md)をHarnessRoot基準で読込（process CWD基準は禁止）；current `active_profile`→shell/tool syntax選択。
2. `agy` / `agy.exe`をactive profile∧[`rules/agy.lrf`](../../rules/agy.lrf)のcanonical解決規則でresolve。実行体なし⇒STOP；別executable・stale pathでの代用は禁止。
3. 適用契約=[`rules/agy.lrf`](../../rules/agy.lrf)∧[`agents/AGENT_ROUTER.md`](../../agents/AGENT_ROUTER.md)∧[`agents/worker-agy.md`](../../agents/worker-agy.md)∧[`agents/SUBAGENTS.lrf`](../../agents/SUBAGENTS.lrf)を読込。model選択前：resolved executable→`models`→current inventory。exit非0∨invalid/unparseable inventory∨model identifierなし⇒委託前STOP。選択IDはcurrent output内の完全一致のみ；捏造・stale ID再利用は禁止。
4. requestのbounded性∧mode/effectsの承認を確認。design∨uncertainty⇒`plan`；明示要求されたbounded workspace editのみ⇒`accept-edits`。inspection=read-only prompt；edit modeは禁止。
5. taskと両立⇒sandbox使用；追加directoryは必要なexact workspaceのみ。`--dangerously-skip-permissions`は禁止。

## Handoff

Canonical wrapper∧field順序を厳密保持；placeholder→task固有のverified factsのみ。permission/evidence/scope/acceptanceの推定は禁止。target∧acceptance∧scope=具体的かつbounded。冒頭2行=`[SIGMA/1]` then `<∑>`、末尾=`</∑>`；追加header/footer・独立した末尾`SIGMA/1`行は禁止。field値のみPOLYGLYPH；codec≠authority。

Canonical handoff（非fenced、dispatch用）:

[SIGMA/1]
<∑>
Role: worker-agy; difficulty=<easy|medium|hard>
Target: <one bounded CODE target>
Acceptance: <observable acceptance and checks>
Scope: <exact allowlist and effect constraints>
Known facts: <verified facts only, or N_A>
</∑>

下記fenceはimmutableなhistorical legacy例；そのままdispatchは禁止。

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

以下fences=immutableなhistorical profile別template；legacy promptはそのままdispatch禁止。使用時は`$p`/`p`を上記canonical handoffに置換；dispatch authorizationではない。使用前：active profile確認→executable解決。`models`成功 ∧ 両field非空のparseable `MODEL_ID<TAB>DISPLAY_NAME` row≥1を要求；不成立⇒委託前STOP。actual inventory表示→内容確認→valid rowからexact IDを手動選択→別のdelegation commandへ代入；IDの推定・自動選択は禁止。代入対象=approved workspace∧構築済handoffのみ。installed effort∈{`low`,`medium`,`high`,`max`}；`med`は禁止；task∧current routing rulesで選択。

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

明示要求されたworkspace editのみ：`--mode plan`→`--mode accept-edits`；その他のscope拡大は禁止。planning⇒`--mode plan`；read-only inspection⇒edit mode禁止。使用options=installed CLIが対応するもののみ；例のflags=`--print`∪`--model`∪`--effort`∪`--mode`∪`--sandbox`∪`--add-dir`∪`--output-format`。JSON出力≠success証明。

### Authentication

認証が必要⇒当taskのexplicit approval後のみunsandboxed invocation可；user自身のbrowser sign-in/access grant完了までpause。credentials/tokens/secrets/private environment valuesの受領・要求・表示・貼付・送信は禁止。authentication/policy/approval/credential/permission block=terminal：FAILED報告∧STOP；retry/fallbackによる回避は禁止。

## Failure and fallback

- retry≤1、pre-launch startup failureのみ（executable/OS/CLI startup失敗を含む）。起動済invocationの再試行は禁止。
- timeout∨その他post-launch operational failure⇒process停止またはexit確認→actual partial workspace diff検査→parentがfallback検討；unrelated/existing editsは保持。
- eligible operational failure後、parent/mainのみがAGENT_ROUTER difficulty mappingに従いcurrently available GPT-family workerを1つ選択可；同一bounded task∧scope、serialのみ。AGYとfallbackの同時実行は禁止；AGY失敗fallbackにAstraは禁止。authentication/policy/approval/credential/permission blockにはfallback禁止。
- malformed result∨generation/backend failure=operational failure：停止/exit確認→partial diff検査→parent判断。JSON構造・status field単独でsuccessと主張しない。

## STOP and success

executable/model不在∨profile/permissions未解決∨invalid handoff∨prohibited effect∨terminal block∨process state未確認⇒STOP∧blocker報告。uncertaintyあり⇒dispatch禁止。

process completion確認 ∧ parent/mainがactual diff検査 ∧ proportionate deterministic checksのexit resultsを実測 ∧ 全acceptance条件との照合⇒success。worker返却=compact evidence{observed changed files,checks∧exit results,failures,residual risk}。effective model報告はactual invocation evidenceで検証済の場合のみ。JSON=構造化出力、success証明ではない。

childはevaluation fieldsをevidenceとして返却可；shared evaluation artifactへの書込は禁止。wait-all後、parent/mainのみが[`agents/SUBAGENTS.lrf`](../../agents/SUBAGENTS.lrf)に従い`~/.harness/evaluation/subagents_model_eval.json`を更新。

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=agy-subagent-router
R|bin.resolve|task:agy|MUST|RO_LOCAL|hardware=../../HARDWARE.md; resolve=active_profile-in-HARDWARE.md; bin=first-existing(agy,agy.exe,$HOME/bin/agy,$env:USERPROFILE/bin/agy.exe); missing=>STOP
R|use.inspect|task:agy|MUST|RO_LOCAL|cmds=models; effect=read-only-inspection
R|use.delegate|task:agy|MUST|LIVE_WRITE|cmds=delegate,exec; pre=explicit-user-request-for-external-delegation; approval=current-task
R|guard.perm|task:agy|MUST_NOT|LW_SCOPE|flag=--dangerously-skip-permissions; external-authority-expansion=forbidden
R|guard.cred|task:agy|MUST_NOT|CRED|deny=credentials,tokens,raw-secret-logs; allow=task-prompts,model-ids
R|guard.truth|task:agy|MUST|RO_LOCAL|unverified=fact-forbidden; main-agent-diff-inspection=mandatory