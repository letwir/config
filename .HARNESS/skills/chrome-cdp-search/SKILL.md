---
name: chrome-cdp-search
description: Search local directories or execute live Web search ETL and page scraping by launching headless Google Chrome via Chrome DevTools Protocol (CDP). Returns structured JSON.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max] 🌐web 🔍local 📑page
🔎 🌐web: Query∧Engine ⊢ engines{bing,google,duckduckgo} ⇒ headless Chrome/CDP live Web search ETL → structured JSON | `node scripts/chrome_cdp_search.mjs --web-query "<query>" --web-engine bing --json-only`
🔎 🔍local: Path∧Query(text∨JSONPath) ⊢ local file/JSON search via Chrome V8/DOM; max-depth=32 → structured JSON | `node scripts/chrome_cdp_search.mjs --target "<path>" --query "<query>" --json-only`
🔎 📑page: URL∧CSS Selector ⊢ public page DOM element extraction → structured JSON | `node scripts/chrome_cdp_search.mjs --target "<url>" --selector "<sel>" --json-only`
⛔ credential-pages forbidden ∧ tokens-in-query forbidden. Unknown/failure ⇒ stop; no invented result. ✅ Success ⇔ requested operation returns structured JSON; temporary DevTools port + Chrome process ⇒ deterministic cleanup.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=chrome-cdp-search
R|use.web|task:chrome-cdp-search&case:web|MUST|RO_PUBLIC|engines=bing,google,duckduckgo; json-only=true
R|use.local|task:chrome-cdp-search&case:local|MUST|RO_LOCAL|target=path; query=text-or-jsonpath; max-depth=32
R|use.page|task:chrome-cdp-search&case:page|MUST|RO_PUBLIC|target=url; selector=css-selector; json-only=true
R|guard.cleanup|task:chrome-cdp-search|MUST|LW_SCOPE|ephemeral-port=DevToolsActivePort; chrome-process=deterministic-cleanup
R|guard.cred|task:chrome-cdp-search|MUST_NOT|CRED|credential-pages=forbidden; tokens-in-query=forbidden
