---
name: sdk
description: Guide users building apps, scripts, CI pipelines, or automations on top of the Codex SDK (TypeScript or Python).
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

SDK-task ⊢ user builds app/script/CI/automation atop Codex SDK ⇒ guide within source contract.
Platforms/languages: TypeScript(@cursor/sdk) ∨ Python(cursor-sdk); runtimes: local(cwd) ∨ cloud(codex-vm); design runtime explicitly first.
OneShot prompt ⇒ single agent run: TS `await Agent.prompt("...")` / Py `Agent.prompt("...")`.
Stateful session ⇒ multi-turn: TS `agent = await Agent.create(); await agent.send("...")`.
ConversationId ⇒ resume existing session: `Agent.resume(conversationId)`.
Install: TS `npm i @cursor/sdk`; Py `pip install cursor-sdk`.
Results: prompt=one-shot execution; session=stateful multi-turn interaction; resume=existing conversation continued.
Limits: only described variants/patterns and runtimes are supported by this contract; no unsupported platform/behavior or deploy/auth implementation inferred. Never expose API keys, bearer tokens, or auth headers.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=sdk
R|use.client|task:sdk|MUST|RO_LOCAL|variants=ts(@cursor/sdk),py(cursor-sdk); patterns=prompt,session,resume
R|use.runtime|task:sdk|MUST|RO_LOCAL|runtimes=local(cwd),cloud(codex-vm); explicit-design-first=true
R|guard.cred|task:sdk|MUST_NOT|CRED|expose=api-keys,bearer-tokens,auth-headers
