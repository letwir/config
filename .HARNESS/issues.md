# Issues and Roadmap Queue

確認日: 2026-09-08 / 対象: `~/.harness`

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

## Termux 再インストール計画（2026-10-12）

### [x] H-09 / P1 / `.harness` 実体への再配置とクライアント参照の統一

**依頼・範囲:** `~/repo/config/.HARNESS` をコピー元、`~/.harness` を実体の SSOT とする。`.harness` 自体を repo への symlink にしない。crabcode の有効な設定置場を優先し、OpenCode・共通 skills/rules は必要な参照先を symlink 化する。今回は計画の記載のみで、再配置は未実施。上の2026年9月の検証結果は旧環境の記録であり、Termux の成功証拠ではない。

**環境・未確認事項:** Android/Termux、bash、Linux aarch64。`crabcode --help` で `--model provider/model`、`--agent`、`--print` を確認。`pwsh` は前提にしない。過去の list 結果だけではディレクトリの実体・リンク・現行設定探索順序は確定できないため、実行前に再確認する。先例検索 `harness Termux migration` は `null`。現行 provider・認証情報は表示せず保持する。

**移行順序 / 完了条件:**

1. [x] **棚卸し:** bash で `$HOME`、`${XDG_CONFIG_HOME:-$HOME/.config}`、実行バイナリを解決。`lstat`/`readlink` 相当でコピー元・SSOT・各参照側の種類、リンク先、存在を確認する。クライアントの実効設定パスと探索優先順位を現行 CLI/一次資料で検証し、`.opencode` と `.config/opencode` の両方を推測で有効化しない。
2. [x] **バックアップ・差分:** 現 SSOT と置換対象を日時付きで退避し、相対パス・種類・ハッシュ・リンク先・権限を manifest に記録する。コピー元との差分を確認し、ローカル追加 skills、今回の issues、memo、現行設定を失わないマージ方針を決める。バックアップは自動読込されない場所へ置き、復元試験を先に行う。
3. [x] **ステージング:** repo は変更せず `.HARNESS` を別ディレクトリへコピーする。既存ローカル差分を保全してからマージし、symlink の意図しない追跡や範囲外の上書きを拒否する。以前提示された `rsync --delete` や無バックアップの `rm`/`ln -sfn` は採用しない。
4. [x] **bash 化:** `$env:USERPROFILE` → `$HOME`、Windows 絶対パス → 解決済み Termux パス、NTFS junction → Unix symlink に読み替える。PowerShell の文字列置換だけで済ませず、引数・引用・dry-run・非0終了・原子的切替・復元の意味を移植する。入口は bash 用に用意し、`bash -n` と利用可能なら ShellCheck、隔離 fixture で検証する。`.exe` 拡張子だけで可否を決めず、実体・shebang・実行形式を確認する。必要な順序は preflight → 同期/リンク検査 → 記録/メモリ → 使用するランチャー。旧 `.ps1` を実行可能と偽装しない。
5. [x] **参照設計:** `~/.harness/crabcode/` を crabcode 設定の実体とし、確認済みの有効設定ディレクトリ（第一候補 `~/.config/crabcode`）からリンクする。既存 `crabcode.jsonc` の provider/model 等を保持し、OpenCode 設定を無検証で代用しない。OpenCode は `~/.harness/opencode/`、共通 skills/rules は `.harness/skills`/`.harness/rules` を必要な配下から参照する。`.agents` 全体のリンクは既存独自内容と探索仕様を確認してから判断する。未使用クライアントの設定・node_modules・依存関係を勝手に作成しない。
6. [x] **切替・検証:** fixture で欠落先、誤リンク、既存通常ディレクトリ、空白パス、再実行、途中失敗、復元を検証してから切り替える。SSOT が実ディレクトリであること、リンク実効先、設定解析、skills/agents 検出、最小起動、dry-run 無変更、正常0/異常非0、非変更範囲を検証する。停止・再起動が必要なら事前に合意する。依存追加、credential の新規使用、VCS 操作は別承認とする。
7. [x] **復元:** 各切替の前状態を manifest から戻せることを確認する。失敗時は新規リンクのみを外し、旧ディレクトリ・設定・権限を復元して同じ検査を実行する。バックアップは受入れ完了まで保持する。

**進捗 (2026-10-12):** 手順1 棚卸し完了: `XDG_CONFIG_HOME` 未設定のため layer2=`~/.config/crabcode` のみ有効（公式 docs の layer 順序を 2026-10-12 取得で確認。`~/.config/opencode` は不在）。`~/.harness` は全項目で `~/repo/config/.HARNESS` の上位集合と確認（差分は issues・BOOTSTRAP の `host` 行・llm-memory SKILL の更新と memo/skills 追加分のみ）のため、追加マージは行わない。手順2 完了: 退避先 `~/backup/harness-20261012-072322` に manifest（path/type/perms/mtime/sha256）を記録、復元試験 PASS。手順5 完了: `~/.harness/crabcode/` を実体とし `~/.config/crabcode` を symlink で受ける。手順6 部分完了: 切替後の `crabcode --print` 起動は exit 0、`crabcode.jsonc` は切替前後で sha256 一致。誤リンク・既存通常ディレクトリ・空白パスの fixture 検証は未実施。手順4 は preflight のみ完了（下記 H-10）。残作業: 他 `.ps1`（sync-harness ほか15ファイル）の bash 移植、残 fixture 検証。

**完了証拠 (2026-10-12 続):** 残作業完了。gpt-6-luna 作業者へ `<Γ>` handoff で `invoke-client-catalog-probe.sh`・`invoke-local-review.sh`・`check-skill-runtime.sh` の bash 移植を委譲（1ファイル1writer・直列・委譲ランタイム拒否後は同一作業者へ明確化して再委譲、1巡目で完了）。主担当による独立検証: 3ファイルとも `bash -n` 通過・実行ビット付与。probe は UTF-16LE base64 引数デコード→実行→終了コード透過（正常0/`false`→1）を Termux 上で実証。local-review は JSON レシート（schema_version/receipt_id/overall_status 同一キー）、成功 exit 0、STOP は exit 2（`invalid_query` 実証）、git 非管理 git 不在環境では FAIL レシート 918バイト・exit 1（git 検証は git 管理リポジトリで PASS 経路確認）。check-skill-runtime は manifest 10項目を分類・`latest_updated` 判定・exit 1（必須3ツール欠落を正しく検出）。`sync-harness.sh` は途中失敗時に `ERROR: consolidation failed: <詳細>` と `status=Partial`/`error=ConsolidationFailed`/`completed=false`/`after` スナップショット付き manifest を原子的に書き、非0終了することを隔離 fixture で確認。既存テスト14アサーション PASS 維持。`bash scripts/invoke-rule-preflight.sh --task change` 全 LRF 構文 PASS（exit 0）。`~/repo/config/.HARNESS` へ同期済みの全変更（sync-harness.sh 更新、新規 .sh 6本、memo 追加分、evaluation runs、llm-mem 実体）を `cp -a` で反映し、`diff -rq` の差分が解消したことを確認。

**クローズ検証 (2026-10-12 10:58 JST):** 手順4・手順6 の残作業は上記完了証拠で解消済みであるため、両ステップを `[x]` に更新した。クローズ前の最終再検証: 全9 `.sh` が `bash -n` 通過、`sync-harness.test.sh` 14アサーション PASS（欠落先・誤リンク・既存通常ディレクトリ・空白パス・復元・dry-run を含む隔離 fixture）、`--consolidate` 再実行は2回とも exit 0 で冪等（`--check` が1を返すのは fixture に全クライアント配置を持たせないためで、テスト内の全配置 fixture では `--check` も exit 0）、`invoke-rule-preflight.sh --task change` 全 LRF 構文 PASS（exit 0）、`~/.harness` と `~/repo/config/.HARNESS` の `diff -rq` 差分ゼロ。ShellCheck は本環境に未導入のため `bash -n` のみ（既存完了記録どおり、利用可能な場合の条件は非該当）。これにより H-09 の全手順が完了。

### [x] H-10 / P2 / 別モデルのサブエージェントで移行作業を分担

**方針:** 主担当と異なる、実際に利用可能なモデルを明示指定する。一覧掲載だけでは認証・実行可否を保証しない。`crabcode --print --no-session-persistence --model <provider/model> --agent <role>` のように bash から指定できる経路を検証する。現在の `task` API はモデル指定を公開していないため、それだけで別モデル実行と称しない。agy は別途明示依頼がない限り使わない。

- [x] bash 版 preflight で LRF の参照・ハッシュ・範囲・停止条件を検証する。bash への読み替えは安全ゲート省略を意味しない。
- [x] 調査担当は設定探索とリンク現状のみ、実装担当は確定したスクリプトの allowlist のみ、検証担当は差分・fixture 結果のみを受け取る。1ファイル1 writer、依存作業は直列、同時 worker は1つとする。
- [x] 最初は読取専用の小さな疎通で実効モデル・終了コード・出力を確認。権限無視オプションは使わず、モデル/認証/権限失敗時は停止して原因を記録する。別モデル不能なら無断で同モデルへ置換しない。
- [x] handoff は対象・既知事実・受入れ条件・禁止事項だけに圧縮し、会話全文・全カタログ・秘密を送らない。返答は変更ファイル、検査結果、未解決点の短い形式に固定。子の追加委譲は既定禁止。
- [x] 主担当が実際の終了と差分を確認して採否を判断。入力/出力/キャッシュの計測条件を記録し、比較可能な消費量がある場合だけ節約効果を判断する。小型モデルや広い context window だけを根拠に削減量を断定しない。

**完了証拠 (2026-10-12):** `crabcode --print --no-session-persistence --model openai/gpt-6-luna` の疎通で `MODEL_PROBE_OK`・exit 0 を確認（親のグローバル既定は未変更、委譲時に明示指定）。同モデルへ `<Γ>`/SIGMA/1 handoff で `scripts/invoke-rule-preflight.sh` の作成を allowlist 限定で委譲し、WORKER_EXIT=0・1ファイルのみ変更・`bash -n` 通過を確認。主担当による独立検証: 正常系は `lrf-preflight/v1` JSON 出力・exit 0、`--verify` は `lrf-preflight-verify/v1 PASS`・exit 0、欠落参照・cycle・ヘッダ不正・フィールド数不正・root外参照はすべて message 付き exit 1。消費 context の比較測定は未実施のため節約効果は未確定。

**今回のモデル確認:** `openai/gpt-5.4-mini` と `openai/gpt-5.3-codex` はモデル一覧に掲載されていたが、いずれも読取専用レビューが「ChatGPT account では非対応」で失敗した。独立レビューは未成立で、レビュー成功・節約効果の証拠にはしない。2026-10-12 に `openai/gpt-6-luna` が同条件で成功したため、以後の委譲は同モデルを明示指定して実行した。設定切替は完了済み（下記 H-09 参照）。

**追加検証・暫定運用 (2026-10-12):** ユーザーの指示により、移行完了までは `~/.config/crabcode/crabcode.jsonc` を実体ファイルとして保持し、ディレクトリも symlink 化しない。今回の検査時には同設定と `.harness/crabcode/` が存在しなかったため、以前読み取った provider 設定（秘密値ではなく環境変数参照）を復元した。既存ファイルからのコピーではない。JSON解析、通常ファイル/通常ディレクトリ、権限0600を検証し、`crabcode models tailnet-llama` によって有効な読込を確認した。`.harness` の対象検索では OpenCode 用 `opencode/opencode.jsonc` のみ見つかり、crabcode JSONC の誤配置を示す証拠はなかった。設定復元前後で GPT-6 系一覧は同じであり、GPT-5 系の失敗原因を JSONC と断定しない。

| モデル | 実行結果 | 今回の用途選定 |
| --- | --- | --- |
| `openai/gpt-6-luna` | `--print --no-session-persistence --agent plan --reasoning-effort low`、ツールなしの短い応答で `MODEL_PROBE_OK` | 小規模な読取専用調査・限定委譲の第一候補 |
| `openai/gpt-6.1-sol` | 同条件で `MODEL_PROBE_OK` | 難しいレビュー/計画の候補 |
| `openai/gpt-6-sol` | 同条件で `MODEL_PROBE_OK` | 利用可能な代替候補 |
| `openai/gpt-6-astra` | 一覧掲載のみ、今回は未実行 | **2026-10-12当時の履歴記録**。旧「主担当専用」方針はユーザーの2026-10-18選択/H-16によりsupersede。現在はhard-Sol失敗結果に対する明示的一回限りleaf subagent retryのみ。 |

**実効モデル検証 (2026-10-12):** `crabcode stats --days 1 --models` と `~/.local/state/crabcode/data.db` の `messages.model` 列を照合した結果、`--model openai/gpt-6-luna` を明示した probe・委譲実行のいずれも、実行記録上のモデルは `tailnet-llama/qwen3.8-flash-next-orca-iq3_m`（`model_preferences.recent` 先頭の既定）として記録された。`--model` 指定は messages 表の実効モデル表示に反映されていない。したがって本日の委譲成果（preflight、bash 移植3本、sync-harness.sh 修正）の実作業者は qwen3.8-flash であり、「gpt-6-luna へ委譲」の記録は指定事実であって実効事実ではない。成果物自体は主担当の独立検証（`bash -n`・fixture・終了コード）で品質確認済み。以後の委譲では起動後に実効モデルを messages 表またはセッション表示で確認してから実質の根拠とする。比較測定は messages 表の `input_tokens`/`output_tokens` が NULL のため従来どおり不能（節約効果は未確定）。

これは応答疎通の証拠であり、ツール呼出し・編集権限・移行実装の検証ではない。親セッションのモデルとグローバル既定モデルは変更していない。委譲時に明示指定する。消費contextの比較測定は未実施で、節約率は未確定。公式 [Using GPT-6](https://developers.openai.com/api/docs/guides/latest-model) は Luna を focused tasks 向けの最速/低コスト候補、6.1 Sol を複雑な作業向けと説明する（2026-10-12取得）；アカウントでの利用可否は上記実行で別途確認した。

### [x] H-11 / P1 / ホスト環境仕様の正本化 (.harness/HARDWARE.md) と各地コマンドの Bash 誘導

**依頼・背景:** 「この環境を.harness/HARDWARE.mdに記録しつつ、それを元に各地のコマンド実行方法をHARDWARE.mdを参照して環境に合わせたシェルを使うように誘導してくれ。」

**確認事実:**
- ハードウェア: Sony XQ-GE44 (Sony Xperia 1 VI), Snapdragon 8 Gen 3 (8 cores aarch64), RAM 16GB (15.6GB), Storage 449GB (/data/user/0, 267GB 空き)。
- OS・カーネル: Android 16, Linux 6.12.38-android16-5-gcaefa63e005e-ab14676683-4k, Termux 0.118.3 (`/data/data/com.termux/files/usr`, `$HOME=/data/data/com.termux/files/home`)。
- シェル環境: 既定シェル `bash`。PowerShell (`pwsh`, `powershell.exe`) および `cmd.exe` は不在・非対応。
- 一時領域: `/tmp` は権限拒否されるため `$TMPDIR` (`$PREFIX/tmp`) が必須。
- ネイティブツール: `rg`, `git`, `python3`, `node`, `go1.27.1`, `cargo`/`crabcode` (v0.0.14), `llm-mem` (PostgreSQL 連携 aarch64)。

**実施内容 & 完了証拠 (2026-10-12):**
1. [x] `.harness/HARDWARE.md` を作成し、ハードウェアスペック、OS/カーネル、Termuxランタイム、シェル規約、ツールチェイン、移行対応表（旧PowerShell→新Bash）をSSOTとして定義。
2. [x] `AGENTS.md` にホスト環境コマンド実行指針（`HARDWARE.md` 参照、Termux Bash 準拠、PowerShell/cmd 禁止）を追記。
3. [x] `rules/BOOTSTRAP.lrf` の `R|host|*` に `hardware=../HARDWARE.md` を紐付け、環境適応シェル遵守を規定。
4. [x] `rules/README-forLLM.md` のナビゲーション例を PowerShell から Bash (`rg`, `invoke-rule-preflight.sh`) に更新し、`HARDWARE.md` 参照を明記。
5. [x] `rules/MANUAL.lrf`, `rules/engineering.lrf` の環境変数表記を `$env:...` から `$HOME`, `$LLM_MEMORY_BIN` 等の POSIX 形式に更新。
6. [x] `skills/default-settings/SKILL.md` (および `SKILL_default.md`), `skills/agy-subagent-router/SKILL.md`, `skills/economics-api/USAGE.md`, `skills/rikei-search/USAGE.md`, `method.md` のシェル指定・コマンド例を Bash / `HARDWARE.md` 準拠に更新。
7. [x] `bash scripts/invoke-rule-preflight.sh --task change` を実行し、全 LRF 構文検査 PASS（exit 0）を確認。

### [x] H-12 / P1 / HARDWARE.md における環境プロファイル切り替え記法 (active_profile) の導入

**依頼・背景:** 「今回の変更、Pwsh・WindowsとBashの環境をHARDWARE.mdで切り替えできる記法にできないかな」

**確認事実 & 設計方針:**
- 静的なOS/シェル固定記述ではなく、`HARDWARE.md` の冒頭に機械可読な `active_profile`（例: `android-termux-bash`、`windows-pwsh`、`linux-native-bash`）セレクターを設ける。
- 各プロファイル定義（端末・OS・シェル・パス体系・一時領域・ツールチェイン・禁止事項）と、プロファイル別コマンド対照表を整備。
- 全エージェント・ルール・スキルは `active_profile` を読み取って適切なシェルおよびパス規約に自動適応する契約とする。

**実施内容 & 完了証拠 (2026-10-12):**
1. [x] `HARDWARE.md` に `active_profile` セレクターブロック（`android-termux-bash` 既定）、3プロファイルの詳細仕様、およびプロファイル別コマンド対照表を実装。
2. [x] `AGENTS.md` に「`HARDWARE.md` の `active_profile` を確認し、指定プロファイルに適合したシェル・パス記法・ツールチェインを用いる」旨を明記。
3. [x] `rules/BOOTSTRAP.lrf` の `R|host|*` を `hardware=../HARDWARE.md; resolve=active_profile-in-HARDWARE.md; adapt-shell-and-paths=active_profile; default-profile=android-termux-bash` に更新。
4. [x] `rules/README-forLLM.md` を更新し、`active_profile` に応じたコマンド選択（Bash / PowerShell 併記）を明記。
5. [x] `skills/default-settings/SKILL.md` (および `SKILL_default.md`), `skills/agy-subagent-router/SKILL.md`, `skills/economics-api/USAGE.md`, `skills/rikei-search/USAGE.md`, `method.md` を `active_profile` 連動・切り替え可能な記法へ刷新。
6. [x] `bash scripts/invoke-rule-preflight.sh --task change` を実行し、全 LRF 構文検査 PASS（exit 0）を確認。

### [x] H-13 / P1 / config/.HARNESS への環境適応同期、パス・スキルのスラッシュ統一、INSTALL.md (探索・切り替え) 作成

**依頼・背景:** 「この変更をconfig/.HARNESSに反映して、パスをスラッシュで統一し、SKILLSも.harness内部のをパス統一してほしい。切り替えのためにINSTALL.mdを作っておいてほしい。内容は｛①:今の環境を探索。;②:active_profileを切り替え;｝」

**確認事実 & 実施方針:**
- `~/.harness`（作業実行SSOT）に加え、リポジトリ原本である `repo/config/.HARNESS` への完全同期が必要。
- スキルのバイナリ解決パスをすべて `.harness` 内部統一形式 `first-existing(~/.harness/skills/<skill>/<bin>, ~/.harness/skills/<skill>/<bin>.exe, ./<bin>, ./<bin>.exe)` に統一し、プラットフォーム依存の `& $bin` 記法を標準化。
- `INSTALL.md` には必須項目 `{①:今の環境を探索。; ②:active_profileを切り替え;}` を明確に配し、Bash/PowerShell 双方からの検出ワンライナーとプロファイル切り替え手順を詳述。

**実施内容 & 完了証拠 (2026-10-12):**
1. [x] 残余スキル（`llama2coder`, `rikei-search` 等を含む全15スキル）の `SKILL.md` を更新し、`~/.harness/skills/...` パス統一、スラッシュ `/` 統一、`& $bin` の標準コマンド化を完了。
2. [x] `INSTALL.md` を作成。必須の「①:今の環境を探索。」（Bash/PowerShell 自動判別ワンライナー・手動探索テーブル）、「②:active_profileを切り替え。」（3プロファイルの仕様・sed/PowerShell 切り替えコマンド）、「③:セットアップと検証」を定義。
3. [x] `~/.harness` の全変更（`HARDWARE.md`, `INSTALL.md`, `AGENTS.md`, `rules/`, `skills/`, `scripts/*.sh`, `method.md`, `issues.md`）を `repo/config/.HARNESS` に同期・反映。
4. [x] `repo/config/.HARNESS` 配下のパス区切りをスラッシュ `/` に統一（`rules/client-distribution.json`, 全 `.lrf`, 全 `SKILL.md` のバックスラッシュ排除）。
5. [x] `bash scripts/invoke-rule-preflight.sh --task change` を `~/.harness` および `repo/config/.HARNESS` の双方で実行し、いずれも全 LRF 構文検査 PASS（exit 0）を確認。

### [ ] H-14 / P1 / 接続パスの {HarnessRoot} 相対化と active_profile による拡張子解決

**依頼・背景:** 「repo/config/.HARNESS のパスを Windows / Linux / Android Termux のどれでも切り替えできるように、SKILL やその他接続パスを相対パスにし、拡張子も HARDWARE.md の状態で切り替えできるようにしたい」

**確認事実:**
- `~/.harness` と `~/repo/config/.HARNESS` は `diff -rq` で差分ゼロ（同一内容の二重配置）。したがって変更は片方だけ進めても漂移する。
- 絶対・home 依存リテラルの主な残存箇所: `AGENTS.md` の `[SIGMA](~/.harness/rules/BOOTSTRAP.lrf)` ほか 3 リンク、`rules/SKILL.lrf:2` の `root=~/.harness/skills`、`rules/MANUAL.lrf` の `main reads UTF8=$HOME/.harness/persona/PERSONA.css`、`agents/SUBAGENTS.lrf:26,35`、各 `skills/*/SKILL.md` の `bin=first-existing(~/.harness/skills/...)`。
- 拡張子は現状 `.sh` / `.ps1` の両実装を環境ごとに人が選んで記述している（`README-forLLM.md:20`、`agents/SUBAGENTS.lrf:10`、`evaluation/README.md`）。`.exe` は `first-existing(...)` の候補列で吸収している。
- `scripts/check-rule-distribution.sh:26` の `resolve()` は既に `{HarnessRoot}/` と `{UserHome}/` トークンを解決しており、`rules/client-distribution.json` はそのトークン形式で運用済み。基盤は存在する。
- `scripts/invoke-rule-preflight.sh` の L 参照検証は `[...](...lrf)` の文字列をそのまま解決するため、トークン化すると当該検証が壊れる。プレフィックス解決を事前に通すか、検証側の仕様変更が同時に必要。
- 現環境の実体: `codex`/`agy`/`llm-mem` は PATH 上有効、`harness-lint`・`proof-checker`・`sec-forensics` バイナリと `pwsh` は不在。`.ps1` 系は現環境では実行不能であり、トークン化しても「両実装を保つ」限り検証は bash 側のみで完結する。

**設計方針（2026-10-12 承認済み）:**
1. 基準点は `{HarnessRoot}/` 相対。`HARDWARE.md` の `path_resolution` に profile 別の絶対ルート（android-termux-bash→`/data/data/com.termux/files/home/.harness`、linux-native-bash→`$HOME/.harness`、windows-pwsh→`%USERPROFILE%\.harness`）を定義し、全箇所は `{HarnessRoot}/...` 統一トークンで受ける。
2. 拡張子は両実装維持＋トークン解決。`HARDWARE.md` に `script_ext`（`.sh` / `.ps1`）と `binary_ext`（`""` / `.exe`）を宣言し、文献側は `scripts/invoke-rule-preflight{script_ext}` のように書く。`bin.resolve` は相対候補のみへ縮約。
3. 反映先は `~/.harness` と `repo/config/.HARNESS` の両方に同一内容（編集→検証→`scripts/sync-harness.sh` で同期）。

**Subtasks / 完了条件:**
1. [ ] `HARDWARE.md` に `path_resolution`（`root` トークン定義とプロファイル別解決値、`script_ext`、`binary_ext`）を追加し、`R|host` の解決手順に「トークン→active_profile 値」の展開を含める。
2. [ ] `AGENTS.md` の 3 リンク、`rules/SKILL.lrf` の `root`、`rules/MANUAL.lrf` の persona 読み取り、`agents/SUBAGENTS.lrf` の評価・agy 参照を `{HarnessRoot}/` 相対へ置換。
3. [ ] `skills/*/SKILL.md` の `bin.resolve` を `{HarnessRoot}/skills/<skill>/<bin>{binary_ext}` 基準へ統一し、`first-existing` の home 絶対候補を削除。
4. [ ] `scripts/invoke-rule-preflight.sh`（および `.ps1`）の L 参照検証が `{HarnessRoot}` トークンを解決して追跡できることを確認。解決不能は fail-closed のまま維持。
5. [ ] `README-forLLM.md` の `preflight :=` 行と `agents/SUBAGENTS.lrf:10` のスクリプト指定を `{script_ext}` 解決記法へ更新。
6. [ ] `bash scripts/invoke-rule-preflight.sh --task change` を `~/.harness` と `repo/config/.HARNESS` の双方で実行し、全 LRF 構文検査 PASS（exit 0）を確認。
7. [ ] `scripts/sync-harness.sh --check` で両ツリーの一致（差分ゼロ）を再確認。

### [*] H-16 / P1 / Crabcode workerの難易度別model profileとAstra leaf retry

**依頼・選択:** Crabcode cloneのagent configとHarness routingを難易度別にする。easy=`openai/gpt-6-luna`/low、medium=`openai/gpt-6-luna`/medium、hard=`openai/gpt-6.1-sol`/high。Astraはhard Solが実効modelであることを確認したfailed resultに対する、親の明示的な同scope一回限りretryのみ。Astraはsubagent-onlyかつleaf-onlyで、mainでも子dispatcherでもない。

**difficulty_assignment:** Main publishes `D_pub=HARD`: change crosses model/profile configuration, dispatch admission, effective-model verification, one-shot escalation reservation, and nested-agent safety.

**Acceptance:** Crabcode native Task routes worker CODE by explicit `subagent_type` profile; OpenCode RTS accepts `difficulty`/`escalation_of`. Static profiles use exact models/efforts without CLI `--model` override; plugin validates profile before child creation and correlated final provider/model afterward; exact first-line worker status is validated; Astra is only a parent-reviewed same-scope hard-Sol FAILED/ESCALATE retry, subagent-only, leaf-only, never root/main. OpenCode RTS enforces one-shot source reservation. Crabcode ignores plugins/custom_tools, so its native Task path relies on parent/router instructions to issue at most one reviewed Astra retry; do not claim runtime one-shot enforcement there. Providers/default/global model and credentials remain unchanged. Tests cover route, failures, scope, duplicate escalation, effective-model mismatch, Astra boundary, and existing wrapper/permission/lease/depth behavior.

**Plan/worker limitation:** User explicitly authorized Main-only pre-code review after separate Auditor invocations failed to return a usable PASS, then explicitly authorized Main to substitute for CODE after two serial `general` Worker attempts stopped before edits. Main checked route/profile consistency, exact parent-message correlation, duplicate escalation reservation, canonical scope ownership, result state/error paths, Astra root/child denial, and mock-test coverage. Main retains implementation/test ownership; neither authorization is an independent review or model/provider evidence.

**Runner evidence (2026-10-18):** `node --experimental-strip-types --test opencode/tests/model-routing.test.mjs` PASS (exit 0, 7/7). `node --experimental-strip-types --test opencode/tests/orchestration.test.mjs` PASS (exit 0, 25/25) using a temporary test-only `@opencode-ai/plugin` shim because the runtime package is not installed; the shim was removed after the run. `node --check --experimental-strip-types` on plugin/tools/helper PASS (exit 0). Python JSON/profile assertions for both configs PASS (exit 0). Root `scripts/invoke-rule-preflight.sh --task change --tag code --tag worker --tag subagent --tag state` and receipt verification PASS (exit 0). Direct structural parse of 10 changed documents PASS. Whole `lint-polyglyph.py` remains FAIL because its existing fixture inventory omits `agents/CALL.md`; this unrelated manifest drift was not changed. Crabcode-nested preflight fails because its directory aliases `rules/` outside the nested preflight root; root LRF preflight passed.
**Task-memory receipts:** knowledge `26e0b653-ce88-4ee1-a2cd-dd1b163ee0fe`, diary `5b432d9a-e5d2-4856-acf5-615cfb8eb659`, walkthrough `75f942b7-64fb-4ee1-8331-cceadac7e2a8` all ingested successfully (exit 0). Local EvaluationJSON parsed but was not submitted to the evaluation DB because separate LIVE_WRITE approval was not granted.

**Independent review:** A distinct read-only Verifier returned `PASS_WITH_ADVISORIES` on current code/config and acceptance. Main reconciled: no actionable blocker; advisories remain that live effective model selection is unverified and CI attestation is absent. The reviewer did not independently rerun tests.

**Residual / not DONE:** No live model invocation was made (would require separate CHARGE approval). The catalog listed model IDs and static profiles/config parse, but Crabcode agent-profile loading and actual terminal provider/model remain unverified. Official Crabcode docs say plugins/custom_tools are ignored; therefore the OpenCode RTS one-shot reservation is not active in Crabcode, and its native Task Astra retry depends on the parent following the explicit review/one-retry rule. Do not claim effective reasoning effort or CI. Keep `[*] TESTED` but not DONE until separately authorized runtime verification resolves this risk. Evaluation DB submission was not authorized/performed.

**Next action:** With separate user approval for CHARGE, invoke a bounded Crabcode native Task using one configured worker profile and verify the terminal assistant `providerID/modelID`; do not run Astra unless a hard-Sol task yields a reviewed FAILED/ESCALATE result. Otherwise keep H-16 TESTED but not DONE. No provider/credential/default-model changes or VCS operations.
