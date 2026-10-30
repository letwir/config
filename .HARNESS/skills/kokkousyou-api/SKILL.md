---
name: kokkou-api
description: 国土交通省データプラットフォーム (MLIT DPF) API検索クライアント。インフラ・地域データ検索、空間検索、件数取得。
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🔍search 🔢count 📍around 🗂️dataset 🗾pref
🔍 search ⊢ query:string ⇒ GraphQL を介し、インフラ・地域データをキーワード検索; optional --limit <N> = 最大結果数; -j = JSON出力 | `$bin "<query>" [--limit <N>] [-j]`
🔢 count ⊢ query:string ∧ Opts ⇒ 検索一致の総件数だけを高速・省トークン出力; -p <Pref> = 都道府県限定; -d <DatasetID> = データセット限定 | `$bin "<query>" -c [-p <Pref>] [-d <DatasetID>]`
📍 around ⊢ lat⊗lon⊗dist ⇒ 指定座標(lat,lon)を中心とする半径(dist)内の空間検索; optional --limit <N> = 最大結果数; -j = JSON出力 | `$bin --around <lat,lon,dist> [--limit <N>] [-j]`
🗂️ filter ⊢ query:string ∧ (Pref∨DatasetID) ⇒ 都道府県名またはデータセットIDで検索を絞込; -p = 都道府県名, -d = データセットID; 両方指定時は複合条件 | `$bin "<query>" -p "<Pref>" -d "<DatasetID>" [-j]`

意味解析: 検索/絞込結果は検索に合致するデータ。count は結果一覧ではなく総件数のみ。`--limit <N>` は結果数を制限し、総件数を意味しない。
失敗・不確定: API/通信/入力/認証エラーは成功結果として扱わず、エラーとして提示。件数または結果が未取得/不確定なら推測・捏造せず、その旨を明示。
安全範囲: API操作は読み取り検索のみ。範囲は task-scoped queries; 更新・削除その他の書込み、範囲外アクセス、認証情報の表示/出力/共有は禁止。`KOKKOU_API` と raw keys は参照・開示しない。
Option reference: -c/--count; -d/--dataset; -p/--prefecture; -j/--json; --limit defaults to 10; --around radius is meters. For LLM inspection of result data, request -j/--json. Documentation only; no API call or runtime verification.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=kokkou-api
R|bin.resolve|task:kokkou-api|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/kokkousyou-api/kokkou,~/.harness/skills/kokkousyou-api/kokkou.exe,./kokkou,./kokkou.exe); missing=>STOP
R|use.inspect|task:kokkou-api|MUST|RO_LOCAL|cmds=search,count,around,filter; effect=read-only-inspection; output=text-or-json
R|guard.cred|task:kokkou-api|MUST_NOT|CRED|env=KOKKOU_API; deny=credentials,raw-keys; allow=task-scoped-queries
R|guard.truth|task:kokkou-api|MUST|RO_LOCAL|unverified=fact-forbidden; source=mlit-dpf
