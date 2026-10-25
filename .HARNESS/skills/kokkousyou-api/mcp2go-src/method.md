# 詳細設計とアプローチ

## [Target] MLIT Client Go 移植設計 (v2: LLM最適化仕様)

### 1. CLI引数の設計
Go言語の標準 `flag` パッケージを用いて、ショートカットとロングネームの重複バインドによるパニックを避けつつ、以下の引数を実装する。
- `flag.BoolVar(&countOnly, "count", false, "...")` / `flag.BoolVar(&cFlag, "c", false, "...")`
- `flag.StringVar(&around, "around", "", "円形空間検索 (形式: lat,lon,distance_meters)")`
- `flag.StringVar(&dataset, "dataset", "", "...")` / `flag.StringVar(&dFlag, "d", "", "...")`
- `flag.StringVar(&prefecture, "prefecture", "", "...")` / `flag.StringVar(&pFlag, "p", "", "...")`
- `flag.BoolVar(&jsonOut, "json", false, "...")` / `flag.BoolVar(&jFlag, "j", false, "...")`
- `flag.IntVar(&limit, "limit", 10, "取得する最大件数 (デフォルト: 10)")`

`flag.Parse()` 実行後、以下のように変数を統合する。
- `isCount := countOnly || cFlag`
- `targetDataset := dataset` (空なら `dFlag`)
- `targetPrefecture := prefecture` (空なら `pFlag`)
- `isJson := jsonOut || jFlag`
- 検索キーワード (`term`) は、位置引数の最初 (`flag.Arg(0)`) から取得する。位置引数が存在しない場合も、空間検索オプション（`--around`）が指定されていれば検索可能とする。

### 2. 環境変数および Windows レジストリ探索アプローチ
- `os.Getenv("KOKKOU_API")` を最優先。
- 存在しない場合、`registry.CURRENT_USER` -> `Environment` キーから `KOKKOU_API` を探索。
- それでもない場合、`registry.LOCAL_MACHINE` -> `SYSTEM\CurrentControlSet\Control\Session Manager\Environment` から探索。

### 3. API通信 (GraphQL) アプローチ
- エンドポイント: `https://data-platform.mlit.go.jp/api_docs/reference/queries/search.html` に基づく `https://data-platform.mlit.go.jp/api/v1/`
- **位置空間検索の動的構築**:
  `--around 35.6895,139.6917,1000` が指定された場合、カンマ区切りでパースし、以下を構築。
  `locationFilter: { geoDistance: { lat: 35.6895, lon: 139.6917, distance: 1000.0 } }`
- **属性フィルタの動的構築**:
  - `prefecture` と `dataset` の両方が指定された場合:
    `attributeFilter: { AND: [ { attributeName: "DPF:dataset_id", is: "dataset_id" }, { attributeName: "DPF:prefecture_name", is: "prefecture_name" } ] }`
  - 片方のみの場合:
    `attributeFilter: { attributeName: "DPF:dataset_id", is: "dataset_id" }`
- **レスポンス構造**:
  件数取得 (`-c` / `--count`) モードの時は `totalNumber` のみを取得し、結果を標準出力に数値のみ出力して終了する。

### 4. ビルドと配置設計
- `go build -o KOKKOU.exe mlit_client.go`
- `C:\Users\letwir\notes\.agents\skills\mlit-client\` または `C:\Users\letwir\.agents\skills\mlit-client\` に配置する。
- スキル定義 `SKILL.md` には、地理空間検索や件数取得などの全引数の詳細なスキーマを定義し、エージェントが自律的に使いこなせるようにする。
