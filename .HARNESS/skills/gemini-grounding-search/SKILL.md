---
name: gemini-grounding-search
description: Use GEMINI_GROUNDING_API_KEY to perform Google search via Gemini's Grounding API, returning structured search results with source attribution. Works from any project.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

🔎 search🌐: request current web information/source-backed research → Google Search Grounding.
📰 Qiita: request article metadata by ItemID ∨ search words → Qiita lookup.
⊢ Search(Query) → `$bin --search "<query>"` → web search results with source attribution. Qiita(ItemID ∨ Word) → `$qiita item <id> [--json]` ∨ `--search "<query>"` → article metadata (automatic-detection integration ∨ direct search).
∴ Output requested retrieved information with source attribution; do not present unverified claims as facts. ⛔ missing command/setup or retrieval failure → report failure; no fabricated result or citation. ✅ success iff requested search results/metadata returned with attribution as applicable.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=gemini-grounding-search
R|bin.resolve|task:gemini-grounding-search|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/gemini-grounding-search/search,~/.harness/skills/gemini-grounding-search/search.exe,./search,./search.exe); qiita=first-existing(~/.harness/skills/gemini-grounding-search/qiita-search,~/.harness/skills/gemini-grounding-search/qiita-search.exe,./qiita-search,./qiita-search.exe); missing=>STOP
R|use.inspect|task:gemini-grounding-search|MUST|RO_PUBLIC|cmds=search,qiita; effect=read-only-retrieval; model=gemini-2.5-flash
R|guard.cred|task:gemini-grounding-search|MUST_NOT|CRED|env=GEMINI_GROUNDING_API_KEY,QIITA_ACCESS_TOKEN; deny=credentials,raw-keys; allow=queries,article-ids
R|guard.truth|task:gemini-grounding-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=google-search-grounding
