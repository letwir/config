---
name: chrome-cdp-search
description: Search local directories (~/.gemini, ~/.opencode, ~/.codex, etc.) or execute live Web search ETL and page scraping by launching headless Google Chrome and evaluating queries inside Chrome's V8 engine and DOM renderer via Chrome DevTools Protocol (CDP). Returns structured JSON with line numbers, JSON key-paths, DOM selector matches, and Web search result rankings. Zero npm dependencies (Node 22 or newer required).
---

# Chrome CDP Searcher & Web ETL Engine

[SPR/XML::ρ→max|protocol:resolve⇒launch⇒query⇒collect⇒cleanup|cues:🔎⊕⊢⊕→⊕⛔⊕✅]
Obj(LocalTarget ∨ PublicQuery ∨ PublicURL) → Mor(ChromeCDP ∘ DOM/V8Query) → Obj(StructuredJSON ∘ SourceEvidence)

<Γ.morphisms>
local_search: Mor(Path ∘ Query ∘ Limits) ⇒ `node scripts/chrome_cdp_search.mjs --target <path> --query <query> --json-only` ⊸ LineMatches ∨ JSONPathMatches;
web_search: Mor(Query ∘ Engine) ⇒ `node scripts/chrome_cdp_search.mjs --web-query <query> --web-engine <engine> --json-only` ⊸ RankedPublicResults;
page_extract: Mor(PublicURL ∘ Selector) ⇒ `node scripts/chrome_cdp_search.mjs --target <url> --selector <selector> --json-only` ⊸ DOMMatches;
</Γ.morphisms>

@lrf=1|aud=GPT-5.6|scope=chrome-cdp-search-contract
R|trigger|task:chrome-cdp-search|MUST|RO_LOCAL|input=local-target-or-public-query-or-public-url; unrelated-browser-automation=false
R|precondition|task:chrome-cdp-search|MUST|RO_LOCAL|required=Node>=22+Chrome-CDP; target-and-limits=resolved-before-launch
R|local|task:chrome-cdp-search&case:local|MUST|RO_LOCAL|read=selected-paths-only; bounds=extensions+max-files+query
R|public|task:chrome-cdp-search&case:web|MUST|RO_PUBLIC|read=public-search-or-public-page; source-url+query+engine=reported
R|cleanup|task:chrome-cdp-search|MUST|LW_SCOPE|ephemeral-profile-and-Chrome-process=deterministic-cleanup; persistent-target-content=unchanged
R|secrets|task:chrome-cdp-search|MUST_NOT|CRED|read-or-return=credentials,tokens,private-environment-values,credential-bearing-pages
R|failure|task:chrome-cdp-search|MUST|RO_LOCAL|on=launch,CDP,timeout,parse,cleanup-failure; return=FAILED+operation+target; silent-fallback=false
R|success|task:chrome-cdp-search|MUST|RO_LOCAL|required=valid-JSON+mode+query-or-selector+bounded-results; empty-results=success-empty-not-fabricated

A zero-dependency search and web discovery skill powered by **Google Chrome DevTools Protocol (CDP)** and headless Chrome. Evaluates queries, structural JSON paths, live web search engines, and DOM selectors inside a real browser instance, returning structured JSON results to the CUI/CLI.

## Features

- **Local & Web ETL in One Tool**: Search local code/markdown/JSON or execute live search engine queries across Bing/Google/DuckDuckGo and scrape arbitrary URLs.
- **Page Load Event Synchronization**: Uses `Page.loadEventFired` CDP events to guarantee complete DOM hydration before evaluation.
- **Chrome V8 & DOM Evaluation**: Uses Chrome's real JavaScript engine to render files (`file:///...`) and execute search queries.
- **Deep JSON Key-Path & Value Traversal**: Recursively traverses nested JSON/JSONC documents with depth-bounding (max depth $\le 32$) and Set-based cyclic reference protection.
- **Multi-Engine Web Parsing & Fallbacks**: Organic search result parser for Bing, Google, and DuckDuckGo with semantic fallback (`li.b_algo`, `div.g`, `<article>`, `h2 a`).
- **DOM CSS Selector Search**: Query local HTML or remote web pages with standard CSS selectors (`document.querySelectorAll`).
- **Concurrency & Race Guard**: Ephemeral port binding (`--remote-debugging-port=0`) with `DevToolsActivePort` detection and exponential retry backoff.
- **Zero External Dependencies**: Pure ESM script utilizing native Node.js 22+ `WebSocket`.
- **Automatic Resource Teardown**: Registers `exit`, `SIGINT`, `SIGTERM`, and `uncaughtException` handlers ensuring zero orphaned Chrome instances.

## Usage

### 1. Web Search ETL Mode
```powershell
# Live Web search via Bing (default)
node scripts/chrome_cdp_search.mjs --web-query "Chrome DevTools Protocol" --json-only

# Live Web search via Google
node scripts/chrome_cdp_search.mjs --web-query "PostgreSQL JSONB indexing" --web-engine google --json-only

# Direct Web URL Scraping with CSS Selector
node scripts/chrome_cdp_search.mjs --target "https://example.com" --selector "h1, p" --json-only
```

### 2. Local File / Directory Search Mode
```powershell
# Local directory search
node scripts/chrome_cdp_search.mjs --target "$HOME\.gemini" --query "model" --json-only
node scripts/chrome_cdp_search.mjs --target "$HOME\.opencode" --query "agent" --json-only

# PowerShell Wrapper
.\scripts\chrome_cdp_search.ps1 -WebQuery "Chrome DevTools Protocol" -JsonOnly
.\scripts\chrome_cdp_search.ps1 -Target "$HOME\.gemini" -Query "model" -JsonOnly
```

### Options

| Flag | Short | Default | Description |
| :--- | :--- | :--- | :--- |
| `-w, --web-query` | `-w` | `""` | Search keywords to query live Web search engines |
| `--web-engine` | | `bing` | Web search engine to use (`bing`, `google`, or `duckduckgo`) |
| `-t, --target` | `-t` | `$HOME\.gemini` | Target file, directory, or remote `http(s)://` URL |
| `-q, --query` | `-q` | `""` | Search text for local lines or JSON key-paths |
| `-s, --selector` | `-s` | `null` | DOM CSS Selector to extract matching HTML elements |
| `--ext` | | `json,jsonc,md,toml,js,ts,html,txt,css` | Comma-separated file extensions for directory scan |
| `-m, --max-files`| `-m` | `30` | Maximum number of candidate files to scan in directory |
| `--json-only` | | `false` | Silence stderr debug logs and output clean JSON only |

---

## Output JSON Schemas

### Mode 1: Web Search Mode Output Schema
```json
{
  "status": "success",
  "engine": "Google Chrome DevTools Protocol (CDP)",
  "mode": "web",
  "timestamp": "2026-08-30T12:00:00.000Z",
  "query": "Chrome DevTools Protocol",
  "searchEngine": "bing",
  "resultsCount": 10,
  "results": [
    {
      "rank": 1,
      "title": "Google Chrome ウェブブラウザ",
      "url": "https://www.bing.com/ck/...",
      "snippet": "Google の最先端技術を搭載し、さらにシンプル、安全、高速になった Chrome をご活用ください"
    }
  ]
}
```

### Mode 2: Local File Mode Output Schema
```json
{
  "status": "success",
  "engine": "Google Chrome DevTools Protocol (CDP)",
  "mode": "local",
  "timestamp": "2026-08-30T12:00:00.000Z",
  "query": "model",
  "selector": null,
  "targetPath": "C:\\Users\\letwir\\.gemini",
  "totalScannedFiles": 5,
  "matchedFilesCount": 1,
  "results": [
    {
      "file": "C:\\Users\\letwir\\.gemini\\AGY_CLI.md",
      "url": "file:///C:/Users/letwir/.gemini/AGY_CLI.md",
      "query": "model",
      "selector": null,
      "totalLines": 83,
      "lineMatchCount": 16,
      "lineMatches": [
        {
          "lineNumber": 15,
          "content": "/* ── Model Taxonomy & Morphism Mapping ── */"
        }
      ],
      "jsonMatchCount": 0,
      "jsonMatches": null,
      "domMatchCount": 0,
      "domMatches": null
    }
  ]
}
```
