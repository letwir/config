---
name: hudousan-api
description: 国土交通省「不動産情報ライブラリ」API (XIT001/XIT002) CLI クライアント。不動産取引価格・成約価格および市区町村一覧を取得。
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🔎🏠API trigger=real-estate transaction/contract price or municipality lookup; ⊢ choose exactly one method:
💴 price ⇒ XIT001; required `-year <YYYY>` + exactly one location selector from `-area <Code>`, `-city <CityCode>`, `-station <StationCode>`; optional `-quarter <1..4>`; CLI: `$bin -price -year <YYYY> -quarter <1..4> -area <Code> [-city <CityCode>]` (CLI example is illustrative; selectors are alternatives, not cumulative requirements).
🏙️ list ⇒ XIT002; required `-area <PrefCode>`; CLI: `$bin -list -area <PrefCode>`.
→ successful response: API response body to stdout (gzip decoded when applicable); request/validation/HTTP/non-200/decompression/output failures → error on stderr and nonzero exit. Credentials: read `KOKKOU_HUDOUSAN_API`; never expose or print key. Network: invocation sends authenticated HTTP GET to selected API endpoint; this is an external network effect, not local inspection. Do not invoke API absent explicit authorization.
Legacy configuration-loader reference: KOKKOU_HUDOUSAN_API may be sourced from process environment or Windows user-registry configuration. Registry fallback must be verified in the selected installed binary before reliance. Never query or print key values during migration; canonical API validation/gates unchanged.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=hudousan-api
R|bin.resolve|task:hudousan-api|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/hudousan-api/hudousan,~/.harness/skills/hudousan-api/hudousan.exe,./hudousan,./hudousan.exe); missing=>STOP
R|use.inspect|task:hudousan-api|MUST|RO_LOCAL|cmds=price,list; effect=read-only-inspection; stdout=text-or-json
R|guard.cred|task:hudousan-api|MUST_NOT|CRED|env=KOKKOU_HUDOUSAN_API; deny=credentials,raw-keys; allow=task-scoped-queries
R|guard.truth|task:hudousan-api|MUST|RO_LOCAL|unverified=fact-forbidden; source=mlit-real-estate-library
