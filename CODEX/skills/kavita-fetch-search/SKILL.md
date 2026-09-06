---
name: kavita-fetch-search
description: Search an authorized Kavita EPUB library by series, author, and full text through its local cache.
---

# Kavita EPUB search

[SPR/XML::ρ→max|protocol:authorize⇒resolve⇒search/fetch⇒evidence|cues:🔎⊕⊢⊕→⊕⛔⊕✅]
Obj(AuthorizedLibrary ∘ Query) → Mor(SeriesSearch ∨ EPUBTextSearch ∨ CacheSync) → Obj(SeriesEvidence ∨ TextEvidence ∨ CacheReceipt)

<Γ.morphisms>
series_search: Mor(Query) ⇒ `.\kavita-fetch-search.exe search <query>` ⊸ SeriesID ∘ Metadata;
text_search: Mor(SeriesID ∘ Query) ⇒ `.\kavita-fetch-search.exe query <series_id> <query>` ⊸ EPUBTextHits;
cache_sync: Mor(AuthorizedLibrary) ⇒ `.\kavita-fetch-search.exe sync` ⊸ LocalMetadataCache;
</Γ.morphisms>

@lrf=1|aud=GPT-5.6|scope=kavita-fetch-search-contract
R|trigger|task:kavita-fetch-search|MUST|RO_LOCAL|library=authorized-configured-Kavita-only; arbitrary-remote-library=false
R|precondition|task:kavita-fetch-search|MUST|RO_LOCAL|required=private-config+resolved-executable; config-values=presence-check-only
R|search|task:kavita-fetch-search|MUST|RO_LOCAL|return=series-id,title,author,summary?,OPDS-link-or-text-hits; empty=explicit
R|cache|task:kavita-fetch-search&case:sync-or-fetch|MUST|LW_SCOPE|write=skill-local-cache-only; remote-library-mutation=false
R|secrets|task:kavita-fetch-search|MUST_NOT|CRED|prompt-or-output-or-log=api-key,config-content,credential-bearing-url
R|failure|task:kavita-fetch-search|MUST|RO_LOCAL|on=config,connection,authorization,download,parse-failure; return=FAILED+non-secret-cause; fabricate-results=false
R|success|task:kavita-fetch-search|MUST|RO_LOCAL|required=command-success+source-series-id+bounded-results; query-needs-series-id-from-search

Use this skill only to search books available in the configured Kavita library.

## Prerequisite

Place a local `config.toml` next to the executable or in the current directory. It must define non-empty `server_url`, `api_key`, and `opds_url` values. Keep this file private; do not place credentials in prompts, output, or this skill.

## Commands

```powershell
# Find a series or author. Add --sync only when a fresh library index is needed.
.\kavita-fetch-search.exe search [--sync] <query>

# Search text within a known series. Missing EPUBs are fetched into the local cache.
.\kavita-fetch-search.exe query <series_id> <query>

# Show EPUB download URLs for a series.
.\kavita-fetch-search.exe dl <series_id>

# Refresh the local series metadata cache.
.\kavita-fetch-search.exe sync
```

Search results provide a series ID, title, author, summary when available, and an OPDS link. Use the series ID with `query` for full-text retrieval.
