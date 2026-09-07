---
name: sdk
description: Guide users building apps, scripts, CI pipelines, or automations on top of the Codex SDK (TypeScript or Python).
---
[SPR/PIDGIN::ρ→max] ⚡prompt 🔄session 🧵resume 📦install
⚡ 実行 prompt ⊢ OneShot ⇒ 単発エージェント実行 | TS: `await Agent.prompt("...")` / Py: `Agent.prompt("...")`
🔄 対話 session ⊢ Stateful ⇒ 複数ターン対話実行 | TS: `agent = await Agent.create(); await agent.send("...")`
🧵 再開 resume ⊢ ConversationId ⇒ 既存セッション再開 | `Agent.resume(conversationId)`
📦 導入 install ⊢ Package ⇒ SDK導入 | TS: `npm i @cursor/sdk` / Py: `pip install cursor-sdk`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=sdk
R|use.client|task:sdk|MUST|RO_LOCAL|variants=ts(@cursor/sdk),py(cursor-sdk); patterns=prompt,session,resume
R|use.runtime|task:sdk|MUST|RO_LOCAL|runtimes=local(cwd),cloud(codex-vm); explicit-design-first=true
R|guard.cred|task:sdk|MUST_NOT|CRED|expose=api-keys,bearer-tokens,auth-headers
