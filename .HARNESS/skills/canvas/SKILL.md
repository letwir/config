---
name: canvas
description: A Codex Canvas is a live React app that the user can open beside the chat. You MUST use a canvas when producing standalone analytical artifacts.
metadata: { surfaces: [ide] }
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🎨create 📐design 🔍dts 🛠️debug
🎨 作成 create 🔎standalone analytical artifact ⊢ user-provided Data ⇒ independent visualization React artifact at `~/.cursor/projects/<ws>/canvases/<name>.canvas.tsx`; required input=Data.
📐 設計 design ⊢ Flat ∧ Minimal ⇒ avoid gradients, emojis, shadows ∧ use `useHostTheme()`.
🔍 定義 dts ⊢ CursorCanvasSDK ⇒ consult public component/type definitions at `~/.cursor/skills-cursor/canvas/sdk/index.d.ts`.
🛠️ 診断 debug ⊢ `<name>.canvas.status.json` ⇒ build success/failure and error diagnosis; report status evidence, not assumed runtime capabilities.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=canvas
R|use.write|task:canvas|MUST|LW_SCOPE|target=~/.cursor/projects/<ws>/canvases/<name>.canvas.tsx; single-file=true; helper-files=forbidden
R|use.import|task:canvas|MUST|RO_LOCAL|source=cursor/canvas; npm-node-builtins=forbidden; fetch-network=forbidden
R|guard.design|task:canvas|MUST|RO_LOCAL|forbid=gradients,emojis,box-shadow,rainbow,giant-text; colors=useHostTheme; empty-state=forbid
R|guard.truth|task:canvas|MUST|RO_LOCAL|data=inline-real-data-only; placeholder-unverified=forbidden
