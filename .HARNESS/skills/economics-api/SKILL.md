---
name: economics-api
description: Fetch Japanese & global economic, financial, corporate filing, macro statistics, and stock market data via unified CLI wrapper (economics.exe).
---
[SPR/PIDGIN::ρ→max] 🇯🇵estat 🏢edinet 🇺🇸sec 📈fred 🌐worldbank 💹yfinance 💱alphavantage
🇯🇵 統計 estat ⊢ StatsID∨Word ⇒ e-Stat 日本政府統計データ取得 | `& $bin estat --stats-data <ID>` ∨ `--search <WORD>`
🏢 開示 edinet ⊢ Date∨DocID ⇒ EDINET 日本有価証券報告書・提出書類取得 | `& $bin edinet --date <YYYY-MM-DD>` ∨ `--doc-id <ID>`
🇺🇸 米開示 sec ⊢ Ticker ⇒ SEC EDGAR 米国企業提出書類・財務ファクト取得 | `& $bin sec --ticker <SYMBOL> [--facts|--submissions]`
📈 マクロ fred ⊢ SeriesID ⇒ FRED 米国・世界マクロ経済時系列データ | `& $bin fred --series <ID>`
🌐 世界指標 worldbank ⊢ Country⊗Indicator ⇒ 世界銀行グローバル開発指標取得 | `& $bin worldbank --country <CODE> --indicator <CODE>`
💹 市場 yfinance ⊢ Ticker ⇒ Yahoo Finance 株価チャート・指標要約 | `& $bin yfinance --symbol <TICKER> [--summary]`
💱 為替株価 alphavantage ⊢ Ticker∨FX ⇒ AlphaVantage 株価・為替レート | `& $bin alphavantage --symbol <TICKER>` ∨ `--from USD --to JPY`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=economics-api
R|bin.resolve|task:economics-api|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/economics.exe,$env:USERPROFILE/.harness/skills/economics-api/economics.exe); missing=>STOP
R|use.inspect|task:economics-api|MUST|RO_LOCAL|cmds=estat,edinet,sec,fred,worldbank,yfinance,alphavantage; effect=read-only-inspection; output=text-or-json
R|guard.cred|task:economics-api|MUST_NOT|CRED|env=ESTAT_API_KEY,EDINET_SUBSCRIPTION_KEY,FRED_API_KEY,ALPHAVANTAGE_API_KEY; deny=credentials,raw-keys; allow=symbols,tickers,series-id
R|guard.truth|task:economics-api|MUST|RO_LOCAL|unverified=fact-forbidden; source=official-economic-apis
