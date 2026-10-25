# 経済・金融・企業財務・マクロ統計 API 利用ガイド & キー申請方法 (USAGE.md)

本スキル (`economics-api`) で対応している 7 種のデータソース API の利用条件、無料アカウント登録・API キー発行手順、および環境変数の設定方法ですわ。

---

## 1. データソース一覧 & API キー申請・登録手順

### ① e-Stat API (日本・政府統計ポータル)
- **費用**: 無料
- **API キー（appId）の取得手順**:
  1. [e-Stat 開発者ポータル](https://www.e-stat.go.jp/api/) にアクセスし、アカウントを新規登録します。
  2. ログイン後、マイページの「API機能（アプリケーションID発行）」を開きます。
  3. 名称・利用目的を入力して「発行」をクリックし、`appId`（32桁などの文字列）を取得します。
- **環境変数設定**:
  - `ESTAT_API_KEY`: 取得した `appId` をセットしてください。

### ② EDINET API (v2) (日本・有価証券報告書・決算情報)
- **費用**: 無料
- **API キー（Subscription-Key）の取得手順**:
  1. [EDINET API ポータルサイト](https://api.edinet-fsa.go.jp/api/auth/index.aspx?mode=1) にアクセスします。
  2. メールアドレスおよび携帯電話番号（SMS認証）を入力してアカウントを作成します。
  3. マイページにログインし、「Subscription-Key（APIキー）」を発行します。
  * ※2024年4月1日以降、Version 2 API では Subscription-Key の指定が必須化されています。
- **環境変数設定**:
  - `EDINET_SUBSCRIPTION_KEY`: 発行された `Subscription-Key` をセットしてください。

### ③ FRED API (St. Louis Fed / 米国・世界マクロ経済指標)
- **費用**: 無料
- **API キーの取得手順**:
  1. [FRED API 申請ページ](https://fred.stlouisfed.org/docs/api/api_key.html) にアクセスし、FRED アカウントを作成/ログインします。
  2. 「Request API Key」をクリックし、使用目的を入力して 32 桁の API Key を取得します。
- **環境変数設定**:
  - `FRED_API_KEY`: 取得した API キーをセットしてください。

### ④ Alpha Vantage API (グローバル株価・為替・財務データ)
- **費用**: 無料枠あり (5 req/min, 100 req/day)
- **API キーの取得手順**:
  1. [Alpha Vantage API Key Claim](https://www.alphavantage.co/support/#api-key) にアクセスします。
  2. アカウントタイプ（Free等）、名前、メールアドレスを入力して「GET FREE API KEY」をクリックします。
  3. 画面上に表示された API Key を保存します。
- **環境変数設定**:
  - `ALPHAVANTAGE_API_KEY`: 取得した API キーをセットしてください。

### ⑤ SEC EDGAR API (米国・企業財務開示情報 / 10-K, 10-Q)
- **費用**: 完全無料
- **事前登録・API キー**: **不要！**
- **注意点**: SEC の規定により HTTP ヘッダーに `User-Agent: CompanyName ContactEmail` 形式の連絡先表記が必須となっています（バイナリ内で自動付与されますが、`--user-agent` オプションで上書き可能です）。

### ⑥ World Bank API (世界銀行 / 世界各国のマクロ経済指標)
- **費用**: 完全無料
- **事前登録・API キー**: **不要！**
- **注意点**: どなたでも即座にクエリを実行して各国 GDP、インフレ率、人口などを取得できます。

### ⑦ Yahoo Finance API (株価・指標 / パブリックエンドポイント)
- **費用**: 完全無料
- **事前登録・API キー**: **不要！**
- **注意点**: ブラウザ互換ヘッダーを自動付与してデータを取り出します。

---

## 2. 環境変数の設定方法 (PowerShell)

キー取得後、`$PROFILE` や環境変数に設定いただくことで、コマンドライン引数での渡しいらずで自動的にキーが読み込まれますわ。

```powershell
$env:ESTAT_API_KEY = "your_estat_app_id"
$env:EDINET_SUBSCRIPTION_KEY = "your_edinet_subscription_key"
$env:FRED_API_KEY = "your_fred_api_key"
$env:ALPHAVANTAGE_API_KEY = "your_alphavantage_api_key"
```

---

## 3. CLI の使い方と使用例

### 統合ラッパー CLI (`economics.exe`) を使用する場合
```powershell
# 日本の有価証券報告書一覧を取得 (EDINET)
.\economics.exe edinet --date 2026-08-14

# 米国Apple社(AAPL)の財務ファクトを取得 (SEC EDGAR)
.\economics.exe sec --ticker AAPL --facts

# 米国政策金利データを取得 (FRED)
.\economics.exe fred --series FEDFUNDS

# 日本の名目GDP時系列を取得 (World Bank)
.\economics.exe worldbank --country JPN --indicator NY.GDP.MKTP.CD

# Apple株価の1ヶ月日足チャートデータを取得 (Yahoo Finance)
.\economics.exe yfinance --symbol AAPL --range 1mo

# トヨタ自動車(7203.T)の株価を取得 (Alpha Vantage)
.\economics.exe alphavantage --symbol 7203.T
```

### 個別 CLI を直接使用する場合
各バイナリ（`estat.exe`, `edinet.exe`, `sec_edgar.exe`, `fred.exe`, `worldbank.exe`, `yfinance.exe`, `alphavantage.exe`）を直接起動することも可能ですわ。
```powershell
.\sec_edgar.exe --ticker AAPL --facts
.\worldbank.exe --country USA --indicator FP.CPI.TOTL.ZG
```
