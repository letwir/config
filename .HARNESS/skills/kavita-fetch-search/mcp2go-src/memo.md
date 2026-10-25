# Memo: 開発ログ

### 2026-06-09 10:56:08
- プロジェクト開始。Kavita Reader の API を利用した EPUB 全文検索ツールの作成タスク。
- Kavita API に関する調査を実施：
  - OpenAPI 仕様を `https://raw.githubusercontent.com/Kareadita/Kavita/develop/openapi.json` から取得しローカルに保存。
  - 認証は `/api/Plugin/authenticate` に API キーを投げて JWT トークンを取得し、以降のリクエストで `Authorization: Bearer <token>` を使用することが判明。
  - シリーズ検索は `/api/Search/search`、チャプター一覧は `/api/Search/chapters-by-series`、ダウンロードは `/api/Download/chapter` を使用する。
  - EPUBファイルはGoの `archive/zip` で解凍してプレーンテキストを抽出し、SQLiteでキャッシュする方針を決定。
  - `decisions.md`、`method.md`、`issues.md`、`memo.md`、`knowledge.md` の階層的ドキュメントを初期作成・更新した。

### 2026-06-09 11:08:00
- Python PoCの動作検証を実施：
  - 実機のKavita接続情報を `config.toml` に設定。
  - 実機からのシリーズ検索・EPUBダウンロードおよびプレーンテキスト抽出を検証。コンソールのエンコーディング調整を追加。
- Go言語での本実装完了：
  - 認証、検索、ダウンロード、ZIPパース、Pure Go SQLite（`modernc.org/sqlite`）によるキャッシュ処理を構築。
  - ユニットテスト (`go test`) を無事パス。
  - `kavita_api.exe` をビルドし、実環境への検索動作がキャッシュの構築と高速同期を含めて期待通り動作することを確認。
- SKILL化の完了：
  - `kavita_search_epub/` ディレクトリを作成。
  - ソースコード、バイナリ、接続用 config.toml を同ディレクトリ内に完全自己同梱。
  - エージェント向けの `SKILL.md` を作成し、ツールの機能と使用方法を詳細にドキュメント化。

### 2026-06-09 11:50:00
- ユーザーの要望に基づき「自律修復型ETL＋サブコマンド設計（search, query, sync）」へ大幅な機能拡張を決定・承認獲得。
- `kavita-fetch-search/` フォルダを作成し、ソースを `src/` サブフォルダへ配置するディレクトリレイアウトへ移行。
- Go APIクライアントを拡張し、全シリーズのフェッチ（all-v2）および詳細メタデータ（metadata）からの著者名マッピングを実装。
- ローカルDBに `series` メタデータキャッシュテーブルを新規導入。
- 各サブコマンド (`sync`, `search`, `query`) を実装。
- メタデータがローカルに存在しない、または検索ヒットしなかった場合に自動でETL同期（自己修復）を回す仕組みを組み込み。
- 接続エラー発生時には即時 `Failed: <error>` で異常終了するハンドリングを徹底。
- `go test` による自動テストをパスし、実機への `sync`, `search`, `query` の動作検証も完璧に完了。
- 成果物フォルダの直下に説明書 `SKILL.md` を作成し、全タスクの完了を記録。
