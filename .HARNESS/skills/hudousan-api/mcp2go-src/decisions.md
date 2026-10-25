# decisions.md

## 開発の命題 (Goal)
旦那様のために「不動産情報ライブラリ」のAPIを操作するクライアントを実装し、最終的にGo言語のバイナリを含むSKILLとして導入いたしますわ！

## プロジェクト方針
1. **フェーズ1: 調査と記録 (Research)**
   - APIマニュアルから仕様を特定し、`knowledge.md` に蓄積いたします。（済）
2. **フェーズ2: Pythonでの概念実証 (Python PoC)**
   - `reinfolib_client.py` として動作するPoCを作成し、APIキー認証やgzipデコード処理の妥当性を検証いたします。
3. **フェーズ3: Go言語への移植 (Go Refactoring)**
   - Python PoCのロジックをGo言語でリライトし、単一の実行可能バイナリとしてビルド可能な構成にします。
4. **フェーズ4: テストとSKILL化 (Verification & Skill Packaging)**
   - 実装をテストし、`~/.agent/skills/` に配置可能な `SKILL.md` および実行可能バイナリをパッケージングいたしますわ。

## APIキーおよび環境変数の管理
- 不動産情報ライブラリAPIのキーは環境変数 `KOKKOU_HUDOUSAN_API` から読み込むことといたします。
- Windows環境下での実行を考慮し、レジストリからの読み込みやフォールバック手段も設計に含めますわ。
