---
name: gemini-grounding-search
description: Use GEMINI_GROUNDING_API_KEY to perform Google search via Gemini's Grounding API, returning structured search results with source attribution. Works from any project.
---
[SPR/PIDGIN::ρ→max] 🌐search 📰qiita
🌐 検索 search ⊢ Query ⇒ Google Search Grounding実時間Web検索(ソース引用付) | `& $bin --search "<query>"`
📰 記事 qiita ⊢ ItemID∨Word ⇒ Qiita記事メタデータ取得(自動検知連携 ∨ 直接検索) | `& $qiita item <id> [--json]` ∨ `--search "<query>"`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=gemini-grounding-search
R|bin.resolve|task:gemini-grounding-search|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/search.exe,$env:USERPROFILE/.harness/skills/gemini-grounding-search/search.exe); qiita=first-existing($PSScriptRoot/qiita-search.exe,$env:USERPROFILE/.harness/skills/gemini-grounding-search/qiita-search.exe); missing=>STOP
R|use.inspect|task:gemini-grounding-search|MUST|RO_PUBLIC|cmds=search,qiita; effect=read-only-retrieval; model=gemini-2.5-flash
R|guard.cred|task:gemini-grounding-search|MUST_NOT|CRED|env=GEMINI_GROUNDING_API_KEY,QIITA_ACCESS_TOKEN; deny=credentials,raw-keys; allow=queries,article-ids
R|guard.truth|task:gemini-grounding-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=google-search-grounding
