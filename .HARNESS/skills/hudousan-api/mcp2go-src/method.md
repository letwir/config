# method.md

## ## [Target]
不動産情報ライブラリAPIのクライアント機能設計。

## アプローチと設計概要
### 1. Python PoC (`reinfolib_client.py`)
- Python 3.x の標準ライブラリ（`urllib.request`, `gzip`, `json`）を利用して実装。
- 実装するコマンドライン引数の仕様:
  - `--price`: 不動産価格情報取得API (`xit001`) を実行するフラグ。
  - `--list`: 都道府県内市区町村一覧取得API (`xit002`) を実行するフラグ。
  - `--area`: 都道府県コード (例: 13)。
  - `--year`: 取引年 (例: 2025)。
  - `--quarter`: 四半期 (例: 2)。
  - `--city`: 市区町村コード (例: 13102)。
- 実行時に `KOKKOU_HUDOUSAN_API` 環境変数からAPIキーを取得し、ヘッダー `Ocp-Apim-Subscription-Key` に付与します。
- レスポンスがgzipエンコードされている場合（`Content-Encoding: gzip`）は、自動でデコードしてJSON標準出力を行います。

### 2. Goへのリファクタリング (`main.go`)
- Goの標準パッケージ `net/http` を使用し、リクエストを送信。
- レスポンスヘッダーの `Content-Encoding` に `gzip` が含まれている場合、`compress/gzip` を用いてデコードを行います。
- APIキーの取得:
  - Windows環境下での実行をサポートするため、環境変数 `KOKKOU_HUDOUSAN_API` だけでなく、Windowsレジストリからの読み込みもフォールバックとして組み込みますわ。
  - レジストリキー: `HKCU\Environment` の `KOKKOU_HUDOUSAN_API`。
- コマンドラインフラグの処理にはGoの `flag` パッケージを使用します。

### 3. SKILLパッケージング
- SKILL名は `reinfolib-client`。
- バイナリ名: `reinfolib.exe`。
- `C:\Users\letwir\.agents\skills\reinfolib-client` 配下にSKILLディレクトリを作成し、`SKILL.md` とバイナリを格納いたします。
