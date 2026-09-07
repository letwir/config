---
name: hatch-pet
description: Create, repair, validate, visually QA, and package Codex-compatible v2 animated pets (8x11 spritesheet, 9 animation rows, 16 look directions).
---
[SPR/PIDGIN::ρ→max] 🐣hatch 🎨rowgen 🧩assemble 🔍verify 📦package
🐣 生成 hatch ⊢ Concept ∨ Brand ⇒ Codex v2 アニメーションペット（8x11アトラス、全9行＋16方向ルック）構築
🎨 行生成 rowgen ⊢ StateSpec ⇒ `$imagegen` によるスプライト各行の生成
🧩 結合 assemble ⊢ ApprovedRows ⇒ 承認済行（rows 0..10）の決定論的8x11アトラス結合
🔍 検証 verify ⊢ VisualQA ⇒ 16方向ルック＆動作アニメーション品質・寸法検査
📦 出力 package ⊢ Atlas8x11 ⇒ `spriteVersionNumber: 2` パッケージング出力

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=hatch-pet
R|use.pipeline|task:hatch-pet|MUST|LW_SCOPE|stages=hatch,rowgen,assemble,verify,package; image-generator=$imagegen
R|use.atlas|task:hatch-pet|MUST|RO_LOCAL|geometry=8x11; rows=0..10; look-directions=16-clockwise; spriteVersionNumber=2
R|guard.cred|task:hatch-pet|MUST_NOT|CRED|expose=tokens,api-keys,private-env-values
