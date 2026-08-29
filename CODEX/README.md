# Codex user configuration map

このディレクトリは、Codexのユーザー設定と実行状態を分けて管理します。

## 手動管理する領域

- `AGENTS.md`, `AGENTS.override.md`: ルートの指示とOverride
- `agents/`: サブエージェント定義 (`*.toml`)
- `persona/`: メインエージェントの性格CSS
- `rules/`: LRFと安全・実行ルール
- `evaluation/`: エージェント評価の入力・スキーマ・根拠
- `skills/`: 手動追加したCodex skill

## 管理対象外の領域

- `sqlite/`, `sessions/`, `logs`, `cache/`: Codexの実行状態
- `plugins/`: プラグインの管理・キャッシュ
- `auth.json`, `cap_sid`, `.sandbox-secrets/`: 認証・サンドボックス状態
- `attachments/`, `generated_images/`, `visualizations/`: 作業生成物

管理対象外の領域は、ファイル名やサイズだけで削除・移動しないこと。稼働中のCodexが参照している可能性があります。

## 参照の正本

SIGMAの正本は `rules/LLM_REF_RULE.md` です。サブエージェント評価の現行JSONは `evaluation/subagents_model_eval.json`、旧スナップショットは `evaluation/archive/` に置きます。
