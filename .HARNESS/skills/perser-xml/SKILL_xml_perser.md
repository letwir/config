---
name: universal-html-interpreter
description: HTML/XML files interpreter. Extract structured information (JSON-LD, OGP, Metadata, Site-specific details) from HTML files.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

🌐parse 🏷️meta 📦jsonld 🛍️dmm
🌐 HTMLFilePath → metadata+JSON-LD+OGP extraction; command=`go run scripts/interpret.go <html-file>`
🏷️ Title ∧ OGP → extract `<title>`, `<meta>`, OGP tags
📦 JSON-LD → extract `application/ld+json` object arrays
🛍️ SKU ∧ Maker → DMM/FANZA etc. e-commerce fields (product number, maker, release date)

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=universal-html-interpreter
R|use.interpret|task:universal-html-interpreter|MUST|RO_LOCAL|script=scripts/interpret.go; runtime=go-1.26+; input=html-file
R|use.extract|task:universal-html-interpreter|MUST|RO_LOCAL|targets=title,meta,ogp,json-ld,dmm-sku
R|guard.cred|task:universal-html-interpreter|MUST_NOT|CRED|expose=private-session-cookies,auth-tokens
