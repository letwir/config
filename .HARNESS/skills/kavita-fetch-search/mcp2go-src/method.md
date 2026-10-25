# Method: Kavita EPUB 全文検索ツールの詳細設計

## ## [Target] 全文検索のアーキテクチャ設計

### 1. データフロー
```mermaid
sequenceDiagram
    participant CLI as CLI (Go)
    participant DB as SQLite Cache
    participant Kavita as Kavita Server
    
    CLI->>Kavita: Authenticate with API Key
    Kavita-->>CLI: Return JWT Token
    
    CLI->>Kavita: Search Series (query)
    Kavita-->>CLI: Return Series ID & Metadata
    
    CLI->>Kavita: Get Chapters for Series
    Kavita-->>CLI: Return Chapter List (ID, Title)
    
    loop Each Chapter
        CLI->>DB: Check if chapter is cached
        alt Cached
            DB-->>CLI: Return cached text
        else Not Cached
            CLI->>Kavita: Download Chapter (EPUB)
            Kavita-->>CLI: Return EPUB (Zip bytes)
            CLI->>CLI: Extract text from EPUB (XHTML/HTML)
            CLI->>DB: Save text to cache
        end
    end
    
    CLI->>CLI: Execute full-text search on texts
    CLI-->>Stdout: Display matching lines with context
```

### 2. SQLite テーブル設計
CGOが不要なPure GoのSQLiteドライバ（`modernc.org/sqlite`）を使用して、以下のテーブルを構築する。

#### `chapters`
- `chapter_id` (INTEGER PRIMARY KEY): KavitaのチャプターID
- `series_id` (INTEGER): シリーズID
- `series_title` (TEXT): シリーズ名
- `chapter_title` (TEXT): チャプター名
- `content` (TEXT): パース済みのプレーンテキスト（XHTMLからタグを除去したもの）
- `cached_at` (DATETIME DEFAULT CURRENT_TIMESTAMP): キャッシュ日時

### 3. EPUB (Zip) パース仕様
1. `.epub` ファイルを `archive/zip` でオープン。
2. `OEBPS/` などのディレクトリ以下にある `.html` または `.xhtml` ファイル（あるいは `content.opf` で指定されたスピン順）を走査。
3. HTMLファイルからタグ(`<...?>`)を正規表現または簡易パーサーで除去し、プレーンテキストを抽出。
4. チャプター単位で結合し、データベースに保存。

### 4. 検索仕様
- 簡易一致（SQLの `LIKE '%query%'`）または FTS5 を利用した全文検索。
- マッチした行とその前後1行（コンテキスト）を取得してStdOutに出力。
