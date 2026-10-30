---
name: kavita-fetch-search
description: Search an authorized Kavita EPUB library by series, author, and full text through its local cache.
---

# Kavita EPUB search

Use this skill only to search books available in the configured Kavita library.

## Prerequisite

Configure the CLI externally using environment variables (do not print or inspect their values):

- `KAVITA_URL` overrides `server_url` and must be an absolute HTTP(S) server URL. Reverse-proxy prefix paths are supported; URL userinfo, query strings (including an empty `?`), and fragments (including an empty `#`) are rejected without echoing the URL.
- `KAVITA_OPDS` overrides `opds_url`. Its final nonempty path component is the API key, decoded from percent encoding. The path must end in `/api/opds/<API key>` (an optional reverse-proxy prefix and trailing slashes are supported). Query strings, fragments and URL userinfo are rejected. An explicitly set invalid or empty value fails closed, without falling back or echoing the value.
- `KAVITA_IMG` is not consumed by this CLI; it does not provide image functionality.

Legacy `server_url` and `opds_url` must also be absolute HTTP(S) URLs without userinfo, query strings (even an empty `?`), or fragments (even an empty `#`); reverse-proxy prefix paths remain supported. Validation errors identify the field, never the supplied URL.

When `KAVITA_OPDS` is absent, legacy `config.toml` values `api_key` and `opds_url` remain the fallback; an independent `KAVITA_URL` override still applies. A complete valid environment configuration needs no file. Otherwise, the CLI looks next to the executable, then in the current directory, for a file supplying the remaining `server_url`, `api_key`, and `opds_url` fields. Keep legacy configuration private; never include credentials or key-bearing URLs in prompts or output.

## Commands

```powershell
# Find a series or author. Add --sync only when a fresh library index is needed.
.\kavita-fetch-search.exe search [--sync] <query>

# Search text within a known series. Missing EPUBs are fetched into the local cache.
.\kavita-fetch-search.exe query <series_id> <query>

# List EPUB chapter/file information; credential-bearing download URLs are omitted.
.\kavita-fetch-search.exe dl <series_id>

# Refresh the local series metadata cache.
.\kavita-fetch-search.exe sync
```

Search results provide a series ID, title, author and summary when available. Use the series ID with `query` for full-text retrieval. OPDS/download URLs containing credentials are intentionally omitted (including from `dl`); output is not a curl-ready download command. HTTP failures retain operation/status context but omit request URLs, response bodies and underlying errors that may disclose secrets.
