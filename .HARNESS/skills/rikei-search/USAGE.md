# 📖 自然科学・工学 (Rikei-Search) APIキー取得・申請マニュアル (USAGE.md)

本スキル (`rikei-search`) で使用する各種学術・材料・物理・工学データベースAPIの **無料アカウント作成、API Key / Bearer Token の取得手順、および環境変数設定方法** のマニュアルですわ。

---

## 🔑 1. Materials Project API (材料科学・結晶構造・物性)

- **費用**: 無料 (学術・個人・非営利)
- **環境変数名**: `MP_API_KEY` (または `MATERIALSPROJECT_API_KEY`)

### 【取得手順】
1. ブラウザで [Materials Project 公式サイト](https://next-gen.materialsproject.org/) にアクセスいたします。
2. 右上の **"Login"** をクリックし、Google / GitHub アカウントまたは Email で無料登録・ログインいたします。
3. ログイン後、画面右上の **"Dashboard"** (または [https://next-gen.materialsproject.org/api](https://next-gen.materialsproject.org/api)) を開きます。
4. **"API Key"** セクションに表示されている英数字のキー文字列をコピーいたします。
5. PowerShell で以下の環境変数を設定（永続化）します：
   ```powershell
   [System.Environment]::SetEnvironmentVariable('MP_API_KEY', 'あなたのAPIキー', 'User')
   ```

---

## 🔭 2. NASA ADS API (Astrophysics Data System / 物理学・天文学論文)

- **費用**: 完全無料
- **環境変数名**: `ADS_API_TOKEN` (または `NASA_ADS_TOKEN`)

### 【取得手順】
1. [NASA ADS UI (Harvard)](https://ui.adsabs.harvard.edu/) にアクセスいたします。
2. 右上の **"Sign Up"** から無料アカウントを作成し、ログインいたします。
3. ログイン後、右上アイコンの **"Account Settings"** -> **"API Token"** を開きます。
4. **"Generate new token"** をクリックして表示される Bearer Token をコピーいたします。
5. PowerShell で以下のように設定いたします：
   ```powershell
   [System.Environment]::SetEnvironmentVariable('ADS_API_TOKEN', 'あなたのADSトークン', 'User')
   ```

---

## ⚡ 3. IEEE Xplore API (電子・情報・物理・工学論文)

- **費用**: 無料（アカウント申請・審査制）
- **環境変数名**: `IEEE_API_KEY` (または `IEEE_XPLORE_API_KEY`)

### 【取得手順】
1. [IEEE Developer Portal](https://developer.ieee.org/) にアクセスいたします。
2. **"Register"** をクリックし、IEEEアカウントを作成してログインいたします。
3. **"App Management"** -> **"Create New App"** を選択し、アプリ名（例: `RikeiSearchClient`）と利用目的（学術調査・個人研究等）を入力いたします。
4. 申請承認後、発行された **Consumer Key (API Key)** をコピーいたします。
5. PowerShell で以下のように設定いたします：
   ```powershell
   [System.Environment]::SetEnvironmentVariable('IEEE_API_KEY', 'あなたのIEEEキー', 'User')
   ```

---

## 🧪 4. ChemSpider API (RSC / 有機・無機化学定数・物性)

- **費用**: 無料 (1,000 req/月)
- **環境変数名**: `CHEMSPIDER_API_KEY`

### 【取得手順】
1. [RSC Developer Portal](https://www.rsc.org/developers/) にアクセスいたします。
2. アカウント登録およびログインを行い、Developer Key 発行を申請いたします。
3. 発行された API Key を `CHEMSPIDER_API_KEY` に設定いたします。

---

## ✨ 5. API Key 不要で即時使える完全オープンAPI

以下のAPIは **API Key の設定を行わなくても即座に完全無料** でご利用いただけますわ！

- **AFLOW / AFLUX API**: 材料熱力学・空間群・結晶物性 DB
- **PubChem PUG-REST API**: 分子構造・物理化学定数・NMR/IR等スペクトル参照 DB (SDBS代替)
- **Crystallography Open Database (COD)**: CIF 結晶データ DB

---

## 💡 トラブルシューティング

API Key が設定されていない状態でコマンドを実行した場合、`rikei-search` ツールはエラー終了せず、キー不要なオープンAPI（PubChem, AFLOW, COD等）へ自動フォールバックして検索結果を提示し、キー取得方法の案内を出力する安全設計になっておりますわ！
