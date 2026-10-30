---
name: universal-html-interpreter
description: HTML/XML files interpreter. Extract structured information (JSON-LD, OGP, Metadata, Site-specific details) from HTML files.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max|target:HTML-metadata-interpreter|legibility:LLM≫human] 🌐parse 🏷️meta 📦jsonld 🛍️dmm
🔎 HTML/XML-file request ⊢ HTML file; `go run scripts/interpret.go <html-file>` → interpret supported HTML → title∪meta∪OGP∪`application/ld+json` object arrays∪[DMM/FANZA: SKU,item#,maker,release-date where present]. XML extent=unspecified. ⛔ unsupported claims∪fabricated results; ✅ return only extracted information.
Extraction reference: output keys={title,description,site_name,ogp}; all application/ld+json blocks → array; DMM/FANZA fields={maker,SKU/item#,release-date} where present. Canonical runtime and script invocation remain controlling.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=universal-html-interpreter
R|use.interpret|task:universal-html-interpreter|MUST|RO_LOCAL|script=scripts/interpret.go; runtime=go-1.26+; input=html-file
R|use.extract|task:universal-html-interpreter|MUST|RO_LOCAL|targets=title,meta,ogp,json-ld,dmm-sku
R|guard.cred|task:universal-html-interpreter|MUST_NOT|CRED|expose=private-session-cookies,auth-tokens
