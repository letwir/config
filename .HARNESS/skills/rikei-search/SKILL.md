---
name: rikei-search
description: Unified CLI & skill suite for natural science and engineering APIs (Materials Project, AFLOW, PubChem, IEEE Xplore, NASA ADS, COD).
---
[SPR/PIDGIN::ρ→max] 🔬search 💎mp ⚛️aflow 🧪pubchem 📜ieee 🔭ads 🧊cod
🔬 統合 search ⊢ Query ⇒ 自然科学・工学全領域横断検索 | `& $bin search "<query>"`
💎 物質 mp ⊢ Formula ⇒ Materials Project バンドギャップ・結晶構造・生成エネルギー | `& $bin mp <formula>`
⚛️ 熱力学 aflow ⊢ Species ⇒ AFLOW/AFLUX 空間群・エンタルピー・熱力学特性 | `& $bin aflow <species>`
🧪 化学 pubchem ⊢ Name ⇒ PubChem 分子量・IUPAC・SMILES・物性スペクトル | `& $bin pubchem <name>`
📜 論文 ieee ⊢ Query ⇒ IEEE Xplore 工学・物理・情報科学論文・DOI | `& $bin ieee "<query>"`
🔭 天文 ads ⊢ Query ⇒ NASA ADS 天体物理学・物理学文献・Bibcode | `& $bin ads "<query>"`
🧊 結晶 cod ⊢ Formula ⇒ COD 結晶構造オープンDB・格子定数・CIF取得 | `& $bin cod <formula>`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=rikei-search
R|bin.resolve|task:rikei-search|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/rikei.exe,$env:USERPROFILE/.harness/skills/rikei-search/rikei.exe); missing=>STOP
R|use.inspect|task:rikei-search|MUST|RO_LOCAL|cmds=search,mp,aflow,pubchem,ieee,ads,cod; effect=read-only-inspection; output=text-or-json
R|fallback.keyless|task:rikei-search&fact:key-missing|MUST|RO_LOCAL|fallback=AFLOW,PubChem,COD; keyless-first=true
R|guard.cred|task:rikei-search|MUST_NOT|CRED|env=MP_API_KEY,ADS_API_TOKEN,IEEE_API_KEY,CHEMSPIDER_API_KEY; deny=credentials,raw-keys; allow=chemical-formula,query,doi
R|guard.truth|task:rikei-search|MUST|RO_LOCAL|unverified=fact-forbidden; source=natural-science-apis
