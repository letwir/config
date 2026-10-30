# Installation

## Build on Linux

Install Go 1.22 or later and system Chromium/Chrome. From this directory:

```sh
go test ./src
go build -o chrome-search ./src
./chrome-search --query "Go language" --engine bing --count 10
```

Cross-compile a portable Linux amd64 executable from another platform:

```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o chrome-search-linux-amd64 ./src
```

The resulting executable has no Go runtime dependency. The destination machine still needs Chromium/Chrome available on `PATH` (or pass `--browser /path/to/chromium`) and public network access. Chromium is intentionally an external runtime requirement, not bundled or downloaded by this package.

## Runtime requirements and limits

The CLI launches Chromium headlessly with an isolated temporary profile and ephemeral DevTools port. It only performs live public-web searches; it does not operate offline. Search providers may change markup, restrict automated requests, or show consent/CAPTCHA pages. Such cases produce a JSON error/no-results failure; there is no API fallback. Results come from provider-specific DOM selectors and are not guaranteed complete or stable. See `SKILL.md` for usage and the legacy Node/PowerShell capabilities.
