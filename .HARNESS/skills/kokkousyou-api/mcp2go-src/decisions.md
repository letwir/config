# 開発の命題とプロジェクト方針

## 1. 命題 (Core Proposition)
- 国土交通省データプラットフォーム (MLIT DPF) を検索するための Python クライアント `mlit_client.py` を Go 言語に移植する。
- 移植した Go コードを Windows 向けバイナリ `KOKKOU.exe` としてビルドする。
- ビルドした `KOKKOU.exe` をエージェントスキルとして `C:\Users\letwir\.agents\skills\mlit-client\` ディレクトリに `SKILL.md` と共に配置し、Antigravity などの AI エージェントがいつでも国土交通省のデータを瞬時に検索できるようにする。

## 2. プロジェクト方針 (Project Guidelines)
- **言語仕様と依存性**: Go言語で記述する。Windows API を叩くために最小限の準標準ライブラリ（`golang.org/x/sys`）を除き、極力標準ライブラリ（`net/http` 等）のみで実装し、ポータビリティとビルド速度を高める。
- **エラー処理と堅牢性**: Windows ターミナルでの日本語文字化けを防ぐため、UTF-8出力を明示的に担保する設計とする。
- **レジストリフォールバック**: 環境変数 `KOKKOU_API` が設定されていない場合、Windows Registry (HKCU / HKLM) から自動で取得を試みる。
- **UI/UX (引数仕様の決定)**:
  - 基本検索: `KOKKOU.exe <term>`
  - 件数のみ取得: `-c` / `--count` (標準出力に数値のみ出力)
  - 地理空間検索: `--around <lat>,<lon>,<distance_meters>`
  - 属性フィルタ: `--dataset <id> / -d <id>`, `--prefecture <name> / -p <name>`
  - 出力フォーマット: `-j` / `--json`
  - 件数制限: `--limit <int>` (デフォルト: 10件)
- **スキル設計**: `SKILL.md` は Antigravity が解釈可能なスキルフォーマットに従い、呼び出し方法や引数を詳細に記載する。
