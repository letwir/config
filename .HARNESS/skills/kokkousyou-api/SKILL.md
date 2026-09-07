---
name: kokkou-api
description: 国土交通省データプラットフォーム (MLIT DPF) API検索クライアント。インフラ・地域データ検索、空間検索、件数取得。
---
[SPR/PIDGIN::ρ→max] 🔍search 🔢count 📍around 🗂️dataset 🗾pref
🔍 検索 search ⊢ Keyword ⇒ インフラ・地域データ検索(GraphQL) | `& $bin "<query>" [--limit <N>] [-j]`
🔢 件数 count ⊢ Keyword ∧ Opts ⇒ ヒット総件数のみ高速出力(トークン節約) | `& $bin "<query>" -c [-p <Pref>] [-d <DatasetID>]`
📍 空間 around ⊢ Lat⊗Lon⊗Dist ⇒ 座標中心半径空間検索 | `& $bin --around <lat,lon,dist> [--limit <N>] [-j]`
🗂️ 絞込 filter ⊢ DatasetID∨Pref ⇒ データセットID/都道府県名指定絞込 | `& $bin "<query>" -p "<Pref>" -d "<DatasetID>" [-j]`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=kokkou-api
R|bin.resolve|task:kokkou-api|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/kokkou.exe,$env:USERPROFILE/.harness/skills/kokkousyou-api/kokkou.exe); missing=>STOP
R|use.inspect|task:kokkou-api|MUST|RO_LOCAL|cmds=search,count,around,filter; effect=read-only-inspection; output=text-or-json
R|guard.cred|task:kokkou-api|MUST_NOT|CRED|env=KOKKOU_API; deny=credentials,raw-keys; allow=task-scoped-queries
R|guard.truth|task:kokkou-api|MUST|RO_LOCAL|unverified=fact-forbidden; source=mlit-dpf
