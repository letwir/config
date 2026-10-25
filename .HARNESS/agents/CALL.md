# サブエージェント呼び出しガイド

このガイドは `.harness/opencode` の `orchestration_dispatch` 実装・README と `agents/AGENT_ROUTER.md` / 各役割定義からまとめた、OpenCode RTS 用の呼び出し手順です。Markdown に書かれたテンプレートだけでは呼び出しは実行されません。OpenCode を所定のランチャーで起動し、ツールが実際に利用可能なセッションからディスパッチしてください。

## 呼び出し手順

1. `AGENT_ROUTER.md` でタスクに適合する役割を選び、`agents/<role>.md` と `<role>.toml` の契約を確認する。Verifier は変更後の独立検証専用で、実装を依頼しない。
2. 対象をリポジトリ相対パスで限定する。`scope_files` には役割が読む必要のある既存ファイルだけを指定する。Verifier では `write_files: []`。
3. 人の依頼・計画が確定してから、以下の形式でタスク専用の handoff を作る。5項目の名前・順序、外側の囲み、末尾の `SIGMA/1` 行を厳守する。未確認情報は事実として書かず `N_A` とする。
4. `orchestration_dispatch` を使って `role`, `handoff`, `scope_files`, `write_files` を渡す。返る task ID は予約 ID であって完了証拠ではない。
5. `orchestration_status` で状態を確認し、完了後に `orchestration_result` で結果を取得する。要求した全タスクの最終結果を待ってから報告する。親が実 diff と必要な独立チェックを確認する。

```text
<Γ>
SIGMA/1: PIDGEN/text
Role: <AGENT_ROUTERで選んだ役割名>
Target: <単一の限定タスク>
Acceptance: <観測可能な完了条件・チェック>
Scope: <作業境界と明示した制約>
Known facts: <検証済み事実、またはN_A>
SIGMA/1
</Γ>
```

## Verifier 呼び出し例

変更後にのみ使います。`<...>` は実際の値に置換してください。Verifier の sandbox は read-only なので変更を許可しません。

```text
role: verifier
scope_files: ["<変更差分と仕様・関連テストのリポジトリ相対パス>"]
write_files: []

<Γ>
SIGMA/1: PIDGEN/text
Role: verifier
Target: 指定スコープ内の変更がAcceptanceを満たすか独立検証する
Acceptance: 実差分と仕様を照合し、読み取り専用で実行可能な決定的チェックの証拠・終了コードを提示する。未実行や不足はINCONCLUSIVEとし、推測でPASSにしない
Scope: 指定されたファイルだけを読む。変更・書き込み・外部操作は行わない
Known facts: <親が確認した変更対象・チェック手段。未確認ならN_A>
SIGMA/1
</Γ>
```

Verifier 実行のために必要なコマンドがスコープ・read-only 制約内で実行できない場合は、無理に実行せず不足条件を報告します。PASS はエージェントの文章だけでは確定せず、親が実行証拠を確認します。

## 対応範囲と制限

- RTS プラグインの現行 role allowlist は `explorer`, `researcher`, `proposer`, `auditor`, `verifier`, `refactorer`, `critic`, `blackhat`, `compressor`, `worker`。呼び出し先エージェントの実在・選択可否は実行時の `app.agents` 解決にも依存します。
- `FileReader` と `Orchestrator` は Router に定義されていますが、現行プラグインの子 role allowlist には含まれません。したがって、この dispatcher 経由で呼べるとは扱いません。
- 書き込みを伴う委譲は `worker` のみ。`worker` 以外の `write_files` は拒否されます。Verifier は read-only であるため `write_files` を空にします。
- ネストは `researcher` / `proposer` の read-only 子からの限定された孫委譲だけが許可されます。Verifier や Worker は葉として扱います。
- 親タスクの権限・効果ルールが適用されます。Handoff は許可を追加しません。認証情報・秘密値を含めません。
- `orchestration_status` は完了を待たず、`orchestration_result` の pending 状態は成功を意味しません。ランタイム未接続、役割未解決、枠不足などは失敗・停止として記録し、呼び出し成功と主張しません。

## 根拠ファイル

- `agents/AGENT_ROUTER.md` — 役割と選択条件
- `agents/SUBAGENTS.lrf` — 権限・ハンドオフ・制約
- `agents/verifier.md`, `agents/verifier.toml` — Verifier 契約と read-only 設定
- `opencode/README.md` — RTS 呼び出し説明
- `opencode/plugins/orchestration.ts` — role allowlist、hand-off 検査、dispatch/status/result 実装
