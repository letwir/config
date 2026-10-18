---
name: kavita-fetch-search
description: Search an authorized Kavita EPUB library by series, author, and full text through its local cache.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

🔎 Trigger: authorized Kavita EPUB library search/lookup only. ⊢ Local bin resolves per `bin.resolve`; missing ⇒ STOP.
📚 `search [--sync] <query>` → OPDS/metadata series + author lookup. `--sync` exact; metadata cache refresh may be skill-local write only.
📖 `query <series_id> <query>` → full-text search within that series' EPUBs; automatic cache. Series-scoped; no broader library scope.
📥 `dl <series_id>` → EPUB download URL retrieval. 🔄 `sync` → local metadata-cache refresh.
⊢ Inputs/flags/field names stay exact: `query`, `series_id`, `--sync`; `search`, `dl`, `sync`. Cache writes, when needed for sync/fetch, are skill-local only; remote mutation=false.
✅ Return command-derived series/author metadata, search results or download URL as applicable; unverified facts forbidden. No-hit presentation and command-specific error behavior are unspecified here: do not infer results, error recovery, or success; preserve and report this ambiguity.
⛔ Authorized Kavita library only; inspection read-only except permitted skill-local cache. Never access `config.toml` credentials, `api_key`, or `server_url`; series-id and metadata only. Do not query any library unless specifically invoked under this contract. Missing binary ⇒ STOP; unknown/error ⇒ stop affected operation, preserve uncertainty, ask rather than invent or broaden.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=kavita-fetch-search
R|bin.resolve|task:kavita-fetch-search|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/kavita-fetch-search/kavita-fetch-search,~/.harness/skills/kavita-fetch-search/kavita-fetch-search.exe,./kavita-fetch-search,./kavita-fetch-search.exe); missing=>STOP
R|use.inspect|task:kavita-fetch-search|MUST|RO_LOCAL|cmds=search,query,dl,sync; effect=read-only-inspection; library=authorized-kavita-only
R|use.cache|task:kavita-fetch-search&case:sync-or-fetch|MUST|LW_SCOPE|write=skill-local-cache-only; remote-mutation=false
R|guard.cred|task:kavita-fetch-search|MUST_NOT|CRED|file=config.toml; deny=credentials,api_key,server_url; allow=series-id,metadata
R|guard.truth|task:kavita-fetch-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=kavita-epub-library
