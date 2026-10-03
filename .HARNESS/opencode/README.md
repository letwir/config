# Canonical OpenCode RTS integration

Configuration, tools, launcher and scheduler live in `.harness/opencode`. Installed `@opencode-ai/plugin` and root SDK **1.18.25** are reused through the approved `node_modules` junction to `.opencode/node_modules`; no dependency installation is needed. Global provider configuration, shell and credential storage are unchanged. Legacy tools export named/default adapters; legacy launcher forwards here. Quit and restart OpenCode after changing configuration.

## Run

```powershell
& "$env:USERPROFILE/.harness/opencode/Invoke-OpenCode.ps1"
# Tests (all process runners / SDK calls are mocked; no live model calls)
$env:TEMP = 'A:/TMP/opencode'; $env:TMP = $env:TEMP
node --experimental-strip-types --test opencode/tests/*.test.mjs
```

The launcher checks CLI 1.18.25, sets `OPENCODE_CONFIG` and `OPENCODE_CONFIG_DIR` for its invocation, and restores both. The global config also connects to the canonical bootstrap for ordinary startup. Bootstrap copies canonical behavior/permissions without recursion, preserving providers, model and shell. Discovery overlap shares a process singleton and registers hooks once.

## Dispatch and observation

Use `orchestration_dispatch` with `role`, exact `handoff`, `scope_files`, `write_files` and optional bounded `timeout_ms`. The handoff must be:

```text
<Γ>
SIGMA/1: PIDGEN/text
Role: researcher
Target: one bounded task
Acceptance: observable checks
Scope: exact task boundary
Known facts: verified facts or N_A
SIGMA/1
</Γ>
```

Reservation returns `task_id` immediately. Local preflight and receipt hash verification, live `app.agents` resolution, canonical scope checks and SDK ancestry verification happen before asynchronous submission. `session.create` uses the root SDK `{body:{parentID,title}}`, never v2 session permission arguments. The assigned `messageID` is correlated to assistant `parentID`. Only a completed final `finish: stop` or completed assistant error terminates the job; idle and intermediate tool steps do not. Events plus bounded periodic polling reconcile reordered/duplicate events and expose permission waits. Results are bounded text evidence, not an automatic declaration that acceptance passed.

`orchestration_status` shows owned state, queued notifications and observer errors. `orchestration_result` retrieves correlated final evidence without waiting for model completion. Await **all** requested final evidence before synthesis while doing only non-conflicting parent work. `orchestration_cancel` cancels the entire task subtree. Abort failure/uncertain stop keeps leases and slots; uncertain prompt acceptance is never automatically re-prompted. An ambiguous session-create response with no ID remains blocked until runtime disposal/manual operator recovery; never bypass it by resubmitting.

Root depth 0 → child 1 → grandchild 2. Root orchestrator can select explicit allowlisted roles. Only researcher/proposer children may delegate read-only grandchildren within inherited exact file scope; worker/verifier are leaves. Six descendants total per process server/worktree singleton, including grandchildren; at most four root children reserve two slots for nested work. Terminal jobs do not consume slots. Canonical case-insensitive realpath leases reject overlapping writes and root edits of leased files.

## Safety and limits

Managed sessions are enforced at `tool.execute.before`, with immutable exact scopes and all files of supported edit/write/patch/move/delete operations checked. Unsupported shapes, native task, opaque processes/external tools and unknown worker sessions fail closed. Reads are restricted to explicit files (`read` and file-scoped `grep`); broad glob/list operations are intentionally denied. Parent handles builds/tests. Root opaque operations are denied while worker leases are active. The existing agy tool stays opt-in: its actual tool boundary calls `context.ask` using an ask-only current-task permission, not a user boolean. The Codex worker remains read-only advisory, not CODE. Existing auth modes are unchanged.

**Queued wake mode:** completion notifications are deduplicated, available in status tools, then injected on the next natural parent model turn. No timer/event prompts or promptAsync calls to parents; noReply is not a safe wake admission proxy. No toast is required.

This is **not sandbox-grade protection**. A residual filesystem check-to-write race exists between the before-hook realpath check and native execution; hard-link aliases and hostile filesystem mutation are not a sandbox boundary. In-memory state has **no restart recovery/persistence**; after restart unknown descendant writes are denied, but historical job evidence and lease recovery are unavailable. SDK mock tests do not prove live GUI event timing, plugin loading, permission prompts, agent/model availability or nested model behavior. Live runtime verification must be done separately; do not claim it from mocks. Errors are surfaced rather than silently swallowed; disposal clears timers and attempts subtree abort, retaining uncertain reservations in the process registry.
