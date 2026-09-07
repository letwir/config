---
name: universal-html-interpreter
description: HTML/XML files interpreter. Extract structured information (JSON-LD, OGP, Metadata, Site-specific details) from HTML files.
---
[SPR/PIDGIN::ρ→max] 🌐parse 🏷️meta 📦jsonld 🛍️dmm
🌐 解析 parse ⊢ HTMLFilePath ⇒ メタデータ・JSON-LD・OGP一括抽出 | `go run scripts/interpret.go <html-file>`
🏷️ メタ meta ⊢ Title ∧ OGP ⇒ `<title>`, `<meta>`, OGPタグ抽出
📦 構造化 jsonld ⊢ JSON-LD ⇒ `application/ld+json` オブジェクト配列抽出
🛍️ 特化 dmm ⊢ SKU ∧ Maker ⇒ DMM/FANZA等ECサイト特定項目（品番、メーカー、配信日）抽出

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=universal-html-interpreter
R|use.interpret|task:universal-html-interpreter|MUST|RO_LOCAL|script=scripts/interpret.go; runtime=go-1.26+; input=html-file
R|use.extract|task:universal-html-interpreter|MUST|RO_LOCAL|targets=title,meta,ogp,json-ld,dmm-sku
R|guard.cred|task:universal-html-interpreter|MUST_NOT|CRED|expose=private-session-cookies,auth-tokens
