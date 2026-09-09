# Issues and Roadmap Queue

確認日: 2026-09-08 / 対象: `C:\Users\letwir\.harness`

サブエージェントを使わず、rg・lsd・dustによる限定検索と同期の読み取り専用チェックからまとめた改善候補。優先度は P1=誤動作・誤判定防止、P2=再現性・運用改善。各項目は未実装の TODO であり、今回の作業は列挙のみ。

## 調査範囲と確認済み事項

- rules、etl、scripts、evaluation の入口と llm-memory の実行契約を調査。全スキルの網羅監査ではない。
- `scripts/sync-harness.ps1 -Check` は終了コード0。gemini、gemini-plugin-ltw、codex、agents、claude、opencode の6共有先すべてで SSOT junction と llm-mem.exe の可視性を確認。
- `skills/harness-lint/harness-lint.exe` は存在する。rg の通常ファイル一覧では `.gitignore` の `*.exe` により表示されないため、欠落とは扱わない。
- `llm-mem search -q 'harness validation workflow' -level 2 -json` は `null`。過去事例は今回の現状証拠に代用しない。
- ローカル構成の不足を対象とするため、公開Web/EPUBは N_A。GUIでの読込、異常系の実行、全スキルの動作、外部の未調査runnerは未検証。
- 既存の変更 `agents/verifier.md` と未追跡の `skills/harness-lint/` は保持。

## Queue

### [x] H-01 / P1 / 同期チェックに機械判定できる失敗結果が必要

**確認事実:** [sync-harness.ps1](scripts/sync-harness.ps1) 59–60行の Write-ErrMsg は Write-Host のみ。216–270行の HealthCheck は欠落・不一致を表示するが、失敗集計や非0終了を返す処理がない。

**不足と影響:** 呼出し側が終了コードだけを見た場合、診断上の異常を成功と誤認し得る。今回の正常系が成功したこととは別の問題。

**Subtasks / 完了条件:** 結果オブジェクトと失敗数を返す。対象欠落・誤リンク・バイナリ欠落を隔離fixtureで再現し、異常時は非0、正常時は0になることを確認する。

**完了証拠 (2026-09-08):** `Invoke-HealthCheck` は `Results`、`FailureCount`、`Healthy` を返し、`-PassThru` で取得可能。Pester 3.4隔離fixture 6件が成功。正常系はexit 0、missing-target / wrong-junction / regular-directory / missing-binaryはexit 1。実環境の `-Check` もexit 0。独立Verifier PASS。

### [x] H-02 / P1 / WhatIfの無変更保証を検証するテストが必要

**確認事実:** 同スクリプト11行で SupportsShouldProcess を宣言しているが、120–121行の llm-memory ディレクトリ作成は ShouldProcess の外側にある。

**不足と影響:** ディレクトリが存在しない条件では、dry-runでも作成されるコード経路がある。今回は実環境で欠落状態を作っていない。

**Subtasks / 完了条件:** 書込み経路をすべて ShouldProcess で保護し、隔離環境の実行前後のファイル一覧・ハッシュで `-WhatIf` が無変更であることを確認する。

**完了証拠 (2026-09-08):** 全filesystem mutationを`ShouldProcess`配下に置き、missing SSOTのpreviewも正常終了する。reparse pointを辿らないpath/type/hash/targetスナップショットによりJunction・Copy両モードとmissing SSOTを検証。Pester全体9/9、独立VerifierのH-02単独実行3/3、変更なしでPASS。

### [x] H-03 / P1 / スキル統合の競合検出と復元可能なmanifestが必要

**確認事実:** 同スクリプト103–106行は各配布元のスキルを `Copy-Item -Recurse -Force` で統合する。111–132行の「latest」バイナリ選択は、日時・バージョン比較ではなく最初に存在する候補を採る。

**不足と影響:** 同名異内容の扱いがコピー順に依存し、採用理由や元のSSOT内容を追いにくい。リンク先のbackup処理はあるが、統合によるSSOT上書きの履歴とは別である。

**Subtasks / 完了条件:** コピー前に相対パスとhashの差分を提示し、競合時の採用方針を明示する。実行manifestとSSOTバックアップを用意し、同名競合の検出・復元を隔離fixtureで確認する。

**完了証拠 (2026-09-08):** `Fail`（既定）/ `PreferSource` / `PreferDestination`、相対path・type・SHA256 preflight、source優先順と`PreviousSource`を実装。Fail/WhatIfはSSOT・manifest・backup・後続linkを変更しない。version 1 manifest、事前backup、drift guard付き`-RestoreManifest`、`-Force`、partial status、binary seed hashを隔離検証。Pester 15/15、独立Verifier PASS。

### [x] H-04 / P1 / ルール配布元と各クライアントの差分検出が必要

**確認事実:** 現在の `.harness/rules/LLM_REF_RULE.md` と `.codex/rules/LLM_REF_RULE.md` は一致していない。前者には `load.index`、`workflow.default` があり、persona参照先も異なる。同期スクリプトは skills を対象とし、`HarnessRules` は宣言されるだけである。

**不足と影響:** 意図したクライアント差分と移行漏れを機械的に区別できない。差分があるだけで権限の矛盾や実行不良とは断定しない。

**Subtasks / 完了条件:** 正本・配布先・許容差分のmanifestを定める。読み取り専用doctorで差分を分類し、各クライアントが実際に読む入口まで確認する。

**完了証拠 (2026-09-09):** `rules/client-distribution.json`に正本hash、必須4 clientのentry chain、client別許容LRF record ID、任意2 clientを定義。read-only doctorはhash、参照文字列、canonical record payload、ID集合、malformed/duplicate、manifest path/schemaを検証。live結果はHarness/Gemini Exact、Codex/OpenCode AllowedDifference、Agents/Claude NotConfigured、exit 0。Pester 8/8をWindows PowerShellとpwshで確認、独立Verifier PASS。

### [x] H-05 / P2 / 評価と工程の再実行手順・証跡が必要

**確認事実:** [main.seq](etl/main.seq) は工程順とretryを定義し、[oracle.lrf](evaluation/oracle.lrf) と [result.schema.json](evaluation/result.schema.json) は判定ケースと出力形を持つ。[evaluation/README.md](evaluation/README.md) はファイル説明のみ。調査したscriptsには `.seq` 実行やoracle/schema照合の参照が見つからない。

**不足と影響:** 宣言があることと、その通り実行・判定されたことを区別する証跡が弱い。ハーネス外のrunnerの有無は未確認。

**Subtasks / 完了条件:** 新しいrunnerを増やす前に既存harness-lint等の担当範囲を確認する。単一の再実行コマンド、対象revision、ケース結果、終了コードを記録する。未承認write・未知effect・複合route・retry上限の結果を再現する。

**完了証拠 (2026-09-09):** `invoke-policy-evaluation.ps1` は構造化 oracle を入力にし、実 `LLM_REF_RULE.md`・routed LRF・`main.seq` から判定、route順、retry上限を導出する。自然言語requestからeffectを再分類せず、期待値は比較専用。単一コマンドでGit revision/dirty、入力・runner・schema・rule hash、21 case、終了コードをUTC+UUID receiptへ原子的に記録。未承認write、未知/競合、CRED禁止、明示承認、複合route、guard不一致、retry 0/1/3/4を再現し、4回目は`retry_exhausted`/Human/停止/exit 1/latest非更新。Pester 5/5、schema適合、独立Verifier PASS。

### [x] H-06 / P2 / スキルの存在確認から実呼出しまでの疎通検査が必要

**確認事実:** HealthCheck はスキル数、junction、llm-mem.exe の存在を確認する。個々のスキルの実行ファイル解決、help契約、クライアントでの読込成功は確認していない。

**不足と影響:** 「配布されている」と「利用できる」の間に検証の空白がある。現状すべてのスキルが壊れているという意味ではない。

**Subtasks / 完了条件:** 対象を重要スキルに絞り、バイナリpath・version/hash・安全なhelp呼出しを一覧化する。クライアントごとに新規セッションでの発見と代表操作の証跡を残す。

**完了証拠 (2026-09-09):** `rules/skill-runtime.json` と `check-skill-runtime.ps1` により exact/final path、SHA256/version、help契約、代表操作、client catalogを分離記録。core 5 toolはHealthy、harness-lint/proof-checker/sec-forensics代表操作はPassed、Gemini catalogはObserved、決定論的catalog APIのないCodexはUnverified、任意3 clientはNotConfigured。timeoutはprocess treeを終了し、reparse越境、hash/token/exit/wrapper不一致をfail closed。Pester 7/7、最終receipt schema適合、独立Verifier PASS。

### [x] H-07 / P2 / 共有ルート移行後の記録先統一が必要

**確認事実:** [DIARY.lrf](rules/DIARY.lrf) 2行は `.gemini/diary.md` への追記を要求する。[state.lrf](rules/state.lrf) は作業場所の状態ファイル、llm-memoryスキルは `memo/<task_id>/` の評価入力前提を持つ。

**不足と影響:** 記録が複数ルートに分かれ、workspace外書込みや重複記録の扱いが不明瞭になりやすい。本件では許可されたworkspace内にタスク記録を置く。

**Subtasks / 完了条件:** task_id単位の正本を決め、diary・eval・ingestの対応を一つの入口から追えるようにする。ルート外書込み権限がない場合の手順も定義し、不要な記録複製を減らす。

**完了証拠 (2026-09-09):** `memo/<task_id>/` を正本とし、local-only `write-task-record.ps1` が固定4成果物のstable hashと最終 ingest/eval/external-diary receiptを`record.json`へ原子的に集約する。task_idのcase/reserved名、reparse、receipt整合、lock timeout、replace失敗をfail closedとし、外部mirrorは別承認・既定N_A。pwshとWindows PowerShellで各Pester 8/8、H-07 live recordはschema適合・memory_report_complete=true・mirror=N_A、独立Verifier PASS。H-01〜H-06は未変更。

### [x] H-08 / P2 / 小規模なローカル調査に適した工程分岐が必要

**確認事実:** [research.lrf](rules/research.lrf) は通常Web-firstとResearcher委譲を要求し、ローカルincidentにのみ軽量経路を定義する。[LLM_REF_RULE.md](rules/LLM_REF_RULE.md) は工程ゲートを規定する。今回のようなローカル構成棚卸し専用の分岐は見当たらない。

**不足と影響:** 対象を絞った文書作業にも不要な探索工程が掛かりやすい。今回は「サブエージェント無し」という明示指示を優先し、主エージェントで検証した。

**Subtasks / 完了条件:** local-reviewの条件、必須証拠、停止条件を明示する。小規模レビューは対象検索→根拠→文書差分チェックで完了でき、外部作用やコード変更に広がると通常ゲートへ戻るケースを確認する。

**完了証拠 (2026-09-10):** `LOCAL_REVIEW` を PRECEDENT 後の独立終端として追加し、固定文字列 `rg`、対象限定 `git status --short`、`git diff --check` のみを実行する collector と hash-only receipt schema を実装した。編集・public・external・live・VCS・混在は `FULL_ETL`、unknown・conflict・不正 target は `STOP` に分類する。collector は pwsh/Windows PowerShell で各 6/6、policy evaluator は pwsh/Windows PowerShell で各 7/7、30-case oracle は local/full/stop と retry-4 を期待どおり分類し、独立 Verifier PASS。軽量経路は FINISH、memo、diary、llm-memory、VCS write を呼ばない。

## 推奨順序

H-01とH-02 → H-03 → H-04 → H-05/H-06 → H-07/H-08。新規機能を増やすより、既存の同期・評価・記録が検証可能であることを先に整える。

## 今回の検証と残余

根拠は実ファイルと正常系チェック。異常系実行・修正・配布は行っていない。検索で見つからなかったものは対象範囲限定の未確認として記載した。PromptDefectを示す根拠はない。調査初期に必要以上のルール/スキル全文を読んだ点はAgentDefectとして完了記録に分離する。

Tags: harness, local-review, sync, verification, roadmap
