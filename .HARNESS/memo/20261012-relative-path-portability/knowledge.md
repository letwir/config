# Knowledge Base: 接続パスの {HarnessRoot} 相対化と拡張子解決

## Context Findings
- `~/.harness`（実行 SSOT）と `~/repo/config/.HARNESS`（リポジトリ原本）は内容が同一の二重配置。片側だけ変更すると漂移する。
- home 絶対リテラルは `AGENTS.md` のリンク 3 箇所、`rules/SKILL.lrf:2` の `root=~/.harness/skills`、`rules/MANUAL.lrf` の persona 読み取り、`agents/SUBAGENTS.lrf:26,35`、各 `skills/*/SKILL.md` の `bin=first-existing(~/.harness/skills/...)` に存在。
- トークン解決基盤は既存: `scripts/check-rule-distribution.sh:26` の `resolve()` が `{HarnessRoot}/` と `{UserHome}/` を解決し、`rules/client-distribution.json` はその形式で運用されている。
- 制約: `scripts/invoke-rule-preflight.sh` の L 参照検証は `[...](...lrf)` の文字列をそのまま解決する。トークンをそのまま書くと追跡が壊れるため、プレフィックス展開を事前に通すか検証側の仕様変更が同時に必要。
- 現環境の実体: `codex` / `agy` / `llm-mem` は PATH 上有効。`harness-lint`・`proof-checker`・`sec-forensics` のバイナリと `pwsh` は不在。`.ps1` 系は現環境では実行不能であり、両実装を維持する場合の検証は bash 側のみで完結する。

## Morphism
- Input: 3 環境（android-termux-bash / linux-native-bash / windows-pwsh）で同一に動くパス指定。
- Operation: 全接続パスを `{HarnessRoot}/` 相対へ統一し、`HARDWARE.md` の `path_resolution` に profile 別の絶対ルート・`script_ext`（`.sh`/`.ps1`）・`binary_ext`（`""`/`.exe`）を宣言して 1 点解決させる。
- Result: 環境ごとの書き換え対象は `HARDWARE.md` の `active_profile` に集約され、文献・スキル・ルール側は環境非依存のトークン表記になる。

## Uncertainty
- 実装は未着手（issues.md H-14 に 7 項目の完了条件を移譲）。
- preflight の L 参照検証を「トークン展開を事前に通す」方式で受け入れるか、検証側を拡張するかの実装選択は H-14 の CODE 段階で確定する。
