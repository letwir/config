### [golang.org/x/sys/windows/registry]
- **参照URL**: [https://pkg.go.dev/golang.org/x/sys/windows/registry](https://pkg.go.dev/golang.org/x/sys/windows/registry)
- **知見**:
  - Windows環境におけるレジストリ操作に特化したGoパッケージ。
  - 環境変数 `KOKKOU_API` をレジストリから読み出すために `registry.OpenKey` と `GetStringValue` を利用可能。
  - ユーザー環境変数キー: `registry.CURRENT_USER` の `"Environment"`
  - システム環境変数キー: `registry.LOCAL_MACHINE` の `"SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment"`
  - 利用には `go get golang.org/x/sys/windows/registry` で依存モジュールの追加が必要。

### [Kavita API]
- **参照URL**: [https://www.kavitareader.com/docs/api/](https://www.kavitareader.com/docs/api/)
  - OpenAPI Spec: [https://raw.githubusercontent.com/Kareadita/Kavita/develop/openapi.json](https://raw.githubusercontent.com/Kareadita/Kavita/develop/openapi.json)
- **知見**:
  - **認証フロー**:
    1. APIキー（Auth Key）を作成（Kavita設定画面の「3rd Party Clients」等から）。
    2. `POST /api/Plugin/authenticate?apiKey={apiKey}&pluginName={pluginName}` を叩いて `UserDto` を取得。
    3. `UserDto` 内の `token`（JWTトークン）を抽出し、以降のリクエストヘッダーに `Authorization: Bearer {token}` を付与。
  - **シリーズ検索**:
    - `GET /api/Search/search?queryString={query}` を用いてシリーズを検索。結果の `series` リストから対象の `id`（シリーズID）を取得。
  - **チャプター取得**:
    - `GET /api/Search/chapters-by-series?seriesId={seriesId}` を用いて、シリーズに紐づくチャプター一覧（`ChapterDto`）を取得。
    - .
  - **EPUBのダウンロード**:
    - `GET /api/Download/chapter?chapterId={chapterId}` を用いてチャプターファイルをダウンロード。
    - 返却されるファイルはZIPアーカイブ（またはEPUB単体）。EPUB自体がZIPフォーマットであるため、Go言語の `archive/zip` を用いて解凍・パース可能。
