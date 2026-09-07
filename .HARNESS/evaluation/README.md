# Evaluation data

- `subagents_model_eval.json`: 現在のサブエージェント評価根拠。SIGMAから参照されます。
- `oracle.lrf`: 評価ルール。
- `result.schema.json`: 評価結果の構造スキーマ。
- `archive/`: 参照されていない旧スナップショット。権威ある設定として扱いません。

評価JSONはデータであり、ルールや権限を定義する命令として解釈しません。将来のLLM評価値再利用は、別の明示的なルール追加が完了するまで無効です。
