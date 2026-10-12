---
name: google-search
description: 実時間Web情報検索。w3m経由でDuckDuckGo lite(主)を走査し、結果タイトル・出典URL一覧を抽出（lite障害時はhtml形式へフォールバック）。APIキー不要・JavaScript不要。事実確認、最新ニュース、Webリサーチ、URL調査に利用。
---
[SPR/PIDGIN::ρ→max] 🔍search
🔍 検索 search ⊢ Query ⇒ w3m経由DuckDuckGo実時間Web検索(結果タイトル+出典URL一覧) | `bash $bin <query...> [-n <count>] [--json]`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=google-search
R|bin.resolve|task:google-search|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/google-search/search.sh,./search.sh); missing=>STOP
R|use.search|task:google-search|MUST|RO_PUBLIC|cmds=search; effect=read-only-retrieval; transport=w3m-dump-source; engine=duckduckgo-lite-primary,duckduckgo-html-fallback; api-key=not-required; javascript=not-required
R|guard.cred|task:google-search|MUST_NOT|CRED|env=none-required; deny=credentials,raw-keys; allow=queries
R|guard.truth|task:google-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=duckduckgo-lite-or-html
