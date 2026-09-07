---
name: hudousan-api
description: 国土交通省「不動産情報ライブラリ」API (XIT001/XIT002) CLI クライアント。不動産取引価格・成約価格および市区町村一覧を取得。
---
[SPR/PIDGIN::ρ→max] 💴price 🏙️list
💴 価格 price ⊢ Area⊗Year⊗Quarter ⇒ 取引・成約価格取得(XIT001) | `& $bin -price -year <YYYY> -quarter <1..4> -area <Code> [-city <CityCode>]`
🏙️ 一覧 list ⊢ Area ⇒ 都道府県内市区町村一覧取得(XIT002) | `& $bin -list -area <PrefCode>`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=hudousan-api
R|bin.resolve|task:hudousan-api|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/hudousan.exe,$env:USERPROFILE/.harness/skills/hudousan-api/hudousan.exe); missing=>STOP
R|use.inspect|task:hudousan-api|MUST|RO_LOCAL|cmds=price,list; effect=read-only-inspection; stdout=text-or-json
R|guard.cred|task:hudousan-api|MUST_NOT|CRED|env=KOKKOU_HUDOUSAN_API; deny=credentials,raw-keys; allow=task-scoped-queries
R|guard.truth|task:hudousan-api|MUST|RO_LOCAL|unverified=fact-forbidden; source=mlit-real-estate-library
