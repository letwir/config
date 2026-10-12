# diary / 20261012-relative-path-portability

- timestamp: 2026-10-12T11:20:00+09:00
- task: 接続パスの {HarnessRoot} 相対化と active_profile による拡張子解決の依頼受け付け
- request-evidence: ユーザー指示「repo/config/.HARNESSのパスをWindows,Linux,Android termuxのどれでも切り替えできるようにSKILLやその他接続パスを相対パスかつ拡張子をHARDWARE.mdの状態で切り替えできるようにしたい」→ 後続指示「Issues.mdを作成して引き継ぎ可能にしてね」
- action:
  1. `~/.harness` と `~/repo/config/.HARNESS` が `diff -rq` で差分ゼロの二重配置であることを確認。
  2. home 絶対リテラルの残存箇所を棚卸し（`AGENTS.md` の 3 リンク、`rules/SKILL.lrf:2` の `root`、`rules/MANUAL.lrf` の persona 読み取り、`agents/SUBAGENTS.lrf:26,35`、各 `skills/*/SKILL.md` の `bin.resolve`）。
  3. `scripts/check-rule-distribution.sh:26` の `resolve()` が `{HarnessRoot}` / `{UserHome}` トークンを既に解決していることを確認（基盤は既存）。
  4. `scripts/invoke-rule-preflight.sh` の L 参照検証が `[...](...lrf)` をそのまま解決するため、トークン化には事前展開または検証側仕様変更が必要と判定。
  5. 3 点（基準点／拡張子解決方式／反映先）を質問し、{HarnessRoot} 相対・両実装維持＋トークン解決・両ツリー同一内容を採用方針として確定。
  6. `issues.md` に H-14 として上記の設計方針と 7 項目の完了条件を追記し、`cp -p` で `repo/config/.HARNESS/issues.md` へ同一内容を反映。
- result: H-14 を issues.md に登録（未完了 `[ ]`、Subtasks 7 項目）。両ツリーは `diff -rq` 一致（exit 0）を再確認。`bash scripts/invoke-rule-preflight.sh --task change` は exit 0。CODE 段階は未着手。
- friction: なし。
- attribution: PromptDefect 0%; AgentDefect 0%
- impact: 環境移行の障壁である home 絶対パスと拡張子の選択が、`HARDWARE.md` の `active_profile` 単一変更へ集約される見通しが確定した。実装は H-14 へ引き継ぎ済み。
- feedback: なし。
- rewritten-request: 現依頼のまま
