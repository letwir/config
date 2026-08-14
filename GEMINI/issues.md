# Issues & Roadmap Queue

## [x] COMPLETED: 時間発展型ナレッジグラフ（Temporal Knowledge Graph）と記憶の自動縮約
- **ID**: `ISSUE-001`
- **概要**: `knowledge.md`（1000行上限）および `diary.md` の線形検索（`rg.exe`）の限界を突破するため、時間経過による事実更新・陳腐化・因果関係を追跡できるTemporal Graph/Vector層を構築し、行数上限到達時に自動で知識を蒸留・縮約（Compaction / Consolidation）するパイプラインを実装する。
- **状態**: `[x] COMPLETED` (llm-memory / JITMIND + Graphiti 融合、Bi-Temporal L0〜L3多段縮約＆トリプル抽出実装・実機検証済)

## [-] IN_PROGRESS: diary.md の因果帰属データからの「自己進化（Self-Evolving）ループ」構築
- **ID**: `ISSUE-002`
- **概要**: `diary.md` に蓄積された `Attribution: PromptDefect vs AgentDefect` のデータをバッチ集計・分析し、`PERSONA.css` や `MACHINE.toml`、各種ルール定義への自動リファクタリング・最適化提案を行う自己改善ループを設計する。
- **サブタスク**:
  - `[-] IN_PROGRESS` **独立した検証・批判エージェント（Verifier / Critic-in-the-loop）の分離**: 生成側とテスト/検証側を独立したサブエージェント（Orchestrator-Worker-Verifier）として分離し、自己確証バイアスを排除したテスト自動実行ゲートを構築する。
- **状態**: `[-] IN_PROGRESS`

## [ ] TODO: 機械的・数学的プルーフ（Proof Checking）のCI直結
- **ID**: `ISSUE-003`
- **概要**: `CODE_RULE.md` にて規定した圏論的射（Morphism: Domain → Codomain）のアノテーションを自然言語コメントに留めず、静的型解析ツール（Go `staticcheck` / `golangci-lint`, Rust `cargo-kani` 等）や形式的証明器と連携させ、機械的・数学的健全性を検証するゲートを構築する。
- **状態**: `[ ] TODO`

## [ ] TODO: MCP / ツールスキーマの動的プルーニング（Schema Token Pruning）
- **ID**: `ISSUE-004`
- **概要**: 多数のプラグイン・MCPツールを初期プロンプトに常駐させず、ユーザー指示やタスクのドメインに応じて必要なツールスキーマのみを動的に選択・インジェクションし、初期コンテキストのトークン消費削減とAttentionの純度向上を実現する。
- **状態**: `[ ] TODO`
