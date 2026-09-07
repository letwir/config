---
name: canvas
description: A Codex Canvas is a live React app that the user can open beside the chat. You MUST use a canvas when producing standalone analytical artifacts.
metadata: { surfaces: [ide] }
---
[SPR/PIDGIN::ρ→max] 🎨create 📐design 🔍dts 🛠️debug
🎨 作成 create ⊢ Data ⇒ 独立可視化React成果物 | `~/.cursor/projects/<ws>/canvases/<name>.canvas.tsx`
📐 設計 design ⊢ Flat ∧ Minimal ⇒ スロップ排除（グラデ・絵文字・影禁止） ∧ `useHostTheme()`
🔍 定義 dts ⊢ CursorCanvasSDK ⇒ 公開コンポーネント・型定義参照 | `~/.cursor/skills-cursor/canvas/sdk/index.d.ts`
🛠️ 診断 debug ⊢ Status ⇒ ビルド成否・エラー診断 | `<name>.canvas.status.json`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=canvas
R|use.write|task:canvas|MUST|LW_SCOPE|target=~/.cursor/projects/<ws>/canvases/<name>.canvas.tsx; single-file=true; helper-files=forbidden
R|use.import|task:canvas|MUST|RO_LOCAL|source=cursor/canvas; npm-node-builtins=forbidden; fetch-network=forbidden
R|guard.design|task:canvas|MUST|RO_LOCAL|forbid=gradients,emojis,box-shadow,rainbow,giant-text; colors=useHostTheme; empty-state=forbid
R|guard.truth|task:canvas|MUST|RO_LOCAL|data=inline-real-data-only; placeholder-unverified=forbidden
