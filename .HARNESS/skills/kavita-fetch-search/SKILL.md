---
name: kavita-fetch-search
description: Search an authorized Kavita EPUB library by series, author, and full text through its local cache.
---
[SPR/PIDGIN::ρ→max] 📚series 📖query 📥dl 🔄sync
📚 シリーズ search ⊢ Query ⇒ シリーズ・著者名検索(OPDS/メタデータ) | `& $bin search [--sync] <query>`
📖 全文 query ⊢ SeriesID⊗Query ⇒ シリーズ内EPUB全文テキスト検索(自動キャッシュ) | `& $bin query <series_id> <query>`
📥 ダウンロード dl ⊢ SeriesID ⇒ EPUBダウンロードURL取得 | `& $bin dl <series_id>`
🔄 同期 sync ⊢ ∅ ⇒ ローカルメタデータキャッシュ更新 | `& $bin sync`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=kavita-fetch-search
R|bin.resolve|task:kavita-fetch-search|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/kavita-fetch-search.exe,$env:USERPROFILE/.harness/skills/kavita-fetch-search/kavita-fetch-search.exe); missing=>STOP
R|use.inspect|task:kavita-fetch-search|MUST|RO_LOCAL|cmds=search,query,dl,sync; effect=read-only-inspection; library=authorized-kavita-only
R|use.cache|task:kavita-fetch-search&case:sync-or-fetch|MUST|LW_SCOPE|write=skill-local-cache-only; remote-mutation=false
R|guard.cred|task:kavita-fetch-search|MUST_NOT|CRED|file=config.toml; deny=credentials,api_key,server_url; allow=series-id,metadata
R|guard.truth|task:kavita-fetch-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=kavita-epub-library
