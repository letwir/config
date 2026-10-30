# Task report: 20261025-anthropic-personnel-web-search

## Purpose
Check whether this pi-agent environment can perform web searches and try the topic “Anthropic personnel moves.”

## Verified observations (accessed 2026-10-25)
- The session exposes the `chrome-cdp-search`, `google-search`, and `gemini-grounding-search` skills.
- `chrome-cdp-search` was executed with Chrome/CDP against Bing and DuckDuckGo and returned structured JSON with `status: success`; Google returned zero results for one targeted query. This verifies the Chrome path is operational for at least Bing and DuckDuckGo.
- Search-engine output for Anthropic personnel queries was noisy/generic and did not provide a reliable, dateable report of specific hires or departures. No personnel change is asserted here.
- The `google-search` w3m method is documented and its script exists; execution via the current PowerShell route returned no result, so that route was not verified operational.
- The Gemini grounding skill is documented, but this task did not verify the executable or API availability.
- No web/browser/search entries were found in `.pi/agent/settings.json`; the methods were discovered as skills exposed in the current session, rather than a verified built-in pi setting.
- Prior-memory lookup for Anthropic personnel/leadership returned `null` (no usable precedent surfaced).

## Constraints / unknowns
- Need a better targeted search and primary/credible source verification before summarizing personnel moves.
- Do not infer staff changes from these search results.
