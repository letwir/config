---
name: google-search
description: 実時間Web情報検索。w3m経由でDuckDuckGo lite(主)を走査し、結果タイトル・出典URL一覧を抽出（lite障害時はhtml形式へフォールバック）。APIキー不要・JavaScript不要。事実確認、最新ニュース、Webリサーチ、URL調査に利用。
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max] 🔍search
🔍 Query→Web current-info lookup; trigger=事実確認∪最新ニュース∪Webリサーチ∪URL調査. Execute `bash $bin <query...> [-n <count>] [--json]`. w3m transport; DuckDuckGo lite=primary, html=fallback iff lite unavailable/fails. Return matching result titles+source URLs; --json requests JSON output. No hits⇒report no results; command/transport failure⇒report error, ¬fabricate results.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=google-search
R|bin.resolve|task:google-search|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/google-search/search.sh,./search.sh); missing=>STOP
R|use.search|task:google-search|MUST|RO_PUBLIC|cmds=search; effect=read-only-retrieval; transport=w3m-dump-source; engine=duckduckgo-lite-primary,duckduckgo-html-fallback; api-key=not-required; javascript=not-required
R|guard.cred|task:google-search|MUST_NOT|CRED|env=none-required; deny=credentials,raw-keys; allow=queries
R|guard.truth|task:google-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=duckduckgo-lite-or-html
