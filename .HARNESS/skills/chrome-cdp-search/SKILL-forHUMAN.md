# Chrome CDP Search (Go CLI)

Chromium/Chrome を使って Bing、Google、DuckDuckGo の公開検索結果を JSON で取得する、Go 標準ライブラリのみの CLI です。

## 使い方

```sh
go run ./src --query "検索語" --engine bing --count 10
```

- `--query` / `-q`: 検索語（必須）
- `--engine`: `bing` / `google` / `duckduckgo`（既定値: `bing`）
- `--count`: 最大結果数 1〜50（既定値: 10）
- `--browser`: Chromium/Chrome 実行ファイル。省略時は PATH から検索
- `--timeout`: 全体の制限時間（例: `45s`、既定値: `30s`）

出力は stdout の JSON 1 件です。正常時は `status`, `query`, `searchEngine`, `resultsCount`, `results` を含み、各結果は `rank`, `title`, `url`, `snippet` を持ちます。エラー時は `status: error` と `error` が出力され、終了コードは 0 以外です。

## 導入

Go 1.22 以降で `go build -o chrome-search ./src` を実行します。生成バイナリの実行に Go は不要です。対象 Linux に Chromium/Chrome を別途インストールしてください。詳細は `INSTALL.md` を参照。

## 注意事項

この CLI はライブ検索専用です。インターネット接続と Chromium/Chrome が必須で、完全オフライン検索には対応しません。検索エンジンの DOM 構造変更、同意画面、CAPTCHA、アクセス制限により抽出に失敗することがあります。結果の完全性・順位は保証されません。ブラウザーは一時プロファイルで起動し、終了時に削除されます。
