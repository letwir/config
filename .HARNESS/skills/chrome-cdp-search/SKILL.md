---
name: chrome-cdp-search
description: Search local directories or execute live Web search ETL and page scraping by launching headless Google Chrome via Chrome DevTools Protocol (CDP). Returns structured JSON.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# Live web search

A dependency-free Go CLI lives in `src/`. Build and run it with `go run ./src --query "Chrome DevTools Protocol" --engine bing --count 10`. Supported engines are `bing`, `google`, and `duckduckgo`.

Options: `--query`/`-q` (required), `--engine` (default `bing`), `--count` (1–50, default 10), `--browser` (explicit Chromium/Chrome executable; otherwise searches PATH), and `--timeout` (Go duration, default `30s`). Output is one JSON object on stdout with status, query, engine, count, results (`rank`, `title`, `url`, `snippet`), or an error. Example: `go run ./src -q "Go language" --engine duckduckgo --count 5 --timeout 45s`.

Requires system Chromium/Chrome and outbound access to the selected public search service; Go itself is needed only to build/run via `go run`. A deployed binary does not require Go. Search pages change and may block or present consent/CAPTCHA pages: this tool fails with JSON error rather than treating missing extraction as success. Extraction uses each engine's current DOM selectors inside Chromium, not a stable search API; completeness and ranking are not guaranteed. The timeout bounds browser startup/search; a fresh isolated temporary profile is removed on normal exit, error, interrupt, and timeout. Existing Node/PowerShell functionality remains available below.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=chrome-cdp-search
R|use.web|task:chrome-cdp-search&case:web|MUST|RO_PUBLIC|engines=bing,google,duckduckgo; json-only=true
R|use.local|task:chrome-cdp-search&case:local|MUST|RO_LOCAL|target=path; query=text-or-jsonpath; max-depth=32
R|use.page|task:chrome-cdp-search&case:page|MUST|RO_PUBLIC|target=url; selector=css-selector; json-only=true
R|guard.cleanup|task:chrome-cdp-search|MUST|LW_SCOPE|ephemeral-port=DevToolsActivePort; chrome-process=deterministic-cleanup
R|guard.cred|task:chrome-cdp-search|MUST_NOT|CRED|credential-pages=forbidden; tokens-in-query=forbidden
