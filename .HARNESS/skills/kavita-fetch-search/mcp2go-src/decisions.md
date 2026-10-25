# Decisions: Kavita EPUB 全文検索ツールの開発

## 開発の命題
KavitaのAPIを使用してEPUB文書をフェッチし、Go製バイナリで全文検索を行って必要な知識（テキストマッチ）をStdOutに表示する。

## プロジェクト方針
1. **技術スタック**:
   - PoC検証: Python
   - 本実装: Go (Golang)
   - キャッシュバッファ: SQLite (Goの `modernc.org/sqlite` または標準SQLインターフェース。Pure GoのSQLiteライブラリがWindowsでCGO不要のため望ましい)
2. **動作環境**:
   - OS: Windows 11 (PowerShell 7)
3. **実行フェーズ**:
   - フェーズ1: 概念実証 (PoC) をPythonで実装し、APIの仕様やEPUBのパース挙動を検証。
   - フェーズ2: Go言語へのリファクタリングおよび高速化。
   - フェーズ3: SQLiteによるローカルキャッシュの実装。
   - フェーズ4: 自動テストと最終動作検証。
4. **仕様方針**:
   - KavitaのURLおよびAPIキーは環境変数 (`KAVITA_URL`, `KAVITA_API_KEY`) または設定ファイルから読み込む。
   - EPUBファイルは `archive/zip` を用いて解凍し、中身の `.xhtml` または `.html` からプレーンテキストを抽出して検索する。
