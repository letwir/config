---
name: chrome-cdp-search
description: Search local directories or execute live Web search ETL and page scraping by launching headless Google Chrome via Chrome DevTools Protocol (CDP). Returns structured JSON.
---
[SPR/PIDGIN::ρ→max] 🌐web 🔍local 📑page
🌐 検索 web ⊢ Query ∧ Engine ⇒ Chrome CDP経由ライブWeb検索ETL | `node scripts/chrome_cdp_search.mjs --web-query "<query>" --web-engine bing --json-only`
🔍 探索 local ⊢ Path ∧ Query ⇒ Chrome V8/DOM走査によるローカルファイル/JSON探索 | `node scripts/chrome_cdp_search.mjs --target "<path>" --query "<query>" --json-only`
📑 抽出 page ⊢ URL ∧ Selector ⇒ DOMセレクタによるWebページ要素抽出 | `node scripts/chrome_cdp_search.mjs --target "<url>" --selector "<sel>" --json-only`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=chrome-cdp-search
R|use.web|task:chrome-cdp-search&case:web|MUST|RO_PUBLIC|engines=bing,google,duckduckgo; json-only=true
R|use.local|task:chrome-cdp-search&case:local|MUST|RO_LOCAL|target=path; query=text-or-jsonpath; max-depth=32
R|use.page|task:chrome-cdp-search&case:page|MUST|RO_PUBLIC|target=url; selector=css-selector; json-only=true
R|guard.cleanup|task:chrome-cdp-search|MUST|LW_SCOPE|ephemeral-port=DevToolsActivePort; chrome-process=deterministic-cleanup
R|guard.cred|task:chrome-cdp-search|MUST_NOT|CRED|credential-pages=forbidden; tokens-in-query=forbidden
