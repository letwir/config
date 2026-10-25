# web-search

APIキー不要のGoogle検索とURL本文取得を行う補助CLIです。

```powershell
go run . --search "Chrome DevTools Protocol" --json
go run . --url "https://example.com" --json
```

検索結果のURLを `--url` に渡して深掘りできます。取得は20秒、本文は最大4 MiBに制限します。
