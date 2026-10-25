# llm-memory (LLM Multi-Client Bi-Temporal Memory CLI)

[English](README.md) | [日本語](README.md)

## ナニコレ？
Windows 4台やDebianからTailscale経由で同時に読み書き・自己編集できる、PostgreSQL 18バックエンドのエージェント記憶基盤CLIです。  
JITMINDの自己編集ライフサイクル（ADD/UPDATE/NOOP）とGraphitiの多段縮約（L0〜L3）・ナレッジグラフ抽出を統合し、愚痴・因果分析から設計方針までを自動で構造化蓄積します。

## 必要なもの
- **PostgreSQL 18+**: `llm_memory` データベースおよび `llm` ユーザー（リモートDebianサーバー等）
- **Go 1.22+**: ビルド用（単一静的バイナリ `llm-mem.exe`）
- **環境変数**:
  - `LLM_MEMORY_DB_URL`: `host:port:dbname:user:password` または `postgres://user:pass@host:port/dbname`
  - `GEMINI_GROUNDING_API_KEY` または `GEMINI_API_KEY`（任意・推奨: 高精度LLM抽出に利用）

## 使い方

### 1. 統合ナレッジ取り込み (`ingest`) ★推奨
ファイル（`diary.md`, `decisions.md`, `Walkthrough.md` 等）やテキストを投入すると、**自己編集判定 ＋ 多段縮約 ＋ 知識グラフ抽出** を一括実行します。

```powershell
$bin = "$env:USERPROFILE\.gemini\config\plugins\ltw-skills-plugin\skills\llm-memory\llm-mem.exe"

# ファイルから自動カテゴリ推定取り込み
& $bin ingest -file "path/to/diary.md"
& $bin ingest -file "path/to/decisions.md"
& $bin ingest -file "path/to/Walkthrough.md"

# 直接テキストを取り込み
& $bin ingest -title "新設計方針" -text "本文..." -cat "knowledge"
```

### 評価イベント (`eval` / `evals`)

評価は通常の `ingest` と分離し、LLM抽出・グラフ化・supersedeを経由せずに保存します。
`eval` は同じ `comparison_key` の直近評価との差分をGo側で計算して新規Insertし、`evals` は履歴参照だけを行います。評価の再利用・モデル推薦は現段階では無効です。

```powershell
& $bin eval -file "evaluation.json"
& $bin evals -key "verifier/code-review/v1" -limit 20 -json
```

入力JSONの最小形式:

```json
{
  "comparison_key": "verifier/code-review/v1",
  "role": "verifier",
  "persona": "adversarial",
  "axes": {
    "counterexample_quality": 0.78,
    "scope_discipline": 0.92
  }
}
```

### 2. 記憶の横断検索 (`search`)
```powershell
# キーワード検索 (L1要点表示)
& $bin search -q "JITMIND" -level 1

# カテゴリ別検索 (diaryの愚痴・反省・因果比率を1行表示)
& $bin search -cat "diary" -level 2

# タグ検索 & JSON構造化出力
& $bin search -tag "postgres" -json
```

### 3. 未縮約レコードのバッチ多段要約 (`compact`)
```powershell
& $bin compact -limit 20
```

### 4. 知識グラフの確認 (`graph`)
```powershell
& $bin graph list
```

### 5. 疎通確認 & 端末一覧 (`status` / `clients`)
```powershell
& $bin status
& $bin clients
```

## 概要詳しく

各状態ファイルに応じた詳細な機能・仕様は以下のドキュメントをご参照ください。

| ドキュメント | 概要 | リンク |
|---|---|---|
| **状態遷移図** | JITMIND自己編集ライフサイクル（ADD / UPDATE / NOOP / DEPRECATED） | [state_diagram.md](docs/state_diagram.md) |
| **ER図とデータ構造** | PostgreSQL 18 のテーブル定義（clients, memories, nodes, edges） | [database_er_diagram.md](docs/database_er_diagram.md) |
| **アーキテクチャ & 縮約** | Multi-client Tailscale 接続と L0〜L3 多段縮約パイプライン | [shm_architecture.md](docs/shm_architecture.md) |

## 状態図
二重時間軸（Valid Time / Transaction Time）における記憶の自己編集・バージョン昇格のライフサイクルを定義しています。  
詳細は [docs/state_diagram.md](docs/state_diagram.md) をご覧ください。

## ER図とデータ構造
端末管理 (`clients`)、二重時間軸記憶 (`memories`)、および知識グラフ (`knowledge_nodes`, `knowledge_edges`) のリレーショナル構造です。  
詳細は [docs/database_er_diagram.md](docs/database_er_diagram.md) をご覧ください。

## SHM/WORM
多端末からの同時参照に耐えうるTrigramインデックス全文検索と、追記型・二重時間軸による改ざん不能な不変履歴（WORM構造）を両立しています。  
詳細は [docs/shm_architecture.md](docs/shm_architecture.md) をご覧ください。

## ライセンス
MIT License / Proprietary to LTW Ecosystem.  
**Warning**: 本ツールの本番環境接続情報（Tailscale IP・DB認証情報）は `.env` や環境変数で管理し、Gitリポジトリへコミットしないでください。
