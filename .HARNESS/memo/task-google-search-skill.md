# google-search スキル作成タスク

## 概要
Termux(Android aarch64)環境向けに、CLIからWeb検索（ググる）を実現するスキル `google-search` を作成した。

## 変更ファイル
- `~/.harness/skills/google-search/SKILL.md` — 新規作成（YAML frontmatter + SPR/PIDGEN + LRF/1）
- `~/.harness/skills/google-search/search.sh` — 新規作成（実行権限付与済）

## 実装仕様
- 主経路: Gemini API `gemini-2.5-flash` + `google_search` grounding（`x-goog-api-key: $GEMINI_API`）で回答テキスト＋出典URL一覧を取得
- フォールバック: DuckDuckGo HTML（`html.duckduckgo.com/html/`）からタイトル＋URL抽出
- 引数: `search.sh <クエリ語...> [-n 件数] [--json]`
- APIキーは標準出力・標準エラーに一切出力しない

## 検証結果
- `bash -n` 構文チェック: PASS
- Gemini Grounding 主経路（東京タワー 高さ）: PASS（回答＋出典3件取得）
- DuckDuckGo フォールバック（GEMINI_API=""）: PASS（Wikipedia等3件取得）
- `--json` モード: PASS（構造化JSON出力確認）

## 残存リスク
- Google HTMLスクレイピングはbot遮断（429/CAPTCHA）多発のため主経路に不使用
- DuckDuckGo は連続アクセス時に一時的IP制限の可能性あり
