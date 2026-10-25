---
name: kavita-fetch-search
description: Search an authorized Kavita EPUB library by series, author, and full text through its local cache.
---

# Kavita EPUB search

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
