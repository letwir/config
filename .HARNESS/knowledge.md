# Knowledge Base

2026-09-08 / scope: .harness / source: scripts/sync-harness.ps1 -Check / supersedes: none

6つのスキル共有先はSSOT junctionとして正常であり、llm-mem.exeへアクセスできた。これは各クライアントでのスキル読込成功を証明しない。

2026-09-09 / scope: H-06 runtime probe / source: `evaluation/runs/skill-runtime-e064e7bd-fa97-4ae0-901a-4e5fa44782e1.json` / supersedes: H-06未検証部分

配布存在と実行可能性は別の状態である。core toolの実行、代表操作、client catalogを別フィールドで記録した結果、Geminiでcatalogを観測できた一方、CodexはCLIに決定論的なskill catalog APIがないためUnverifiedとして保持するのが正確だった。

2026-09-09 / scope: H-07 task records / source: `memo/H-07/record.json` / supersedes: split diary/eval/ingest record uncertainty

`memo/<task_id>/` をworkspace正本にし、外部作用の最終receiptだけをlocal `record.json`へ集約すれば、外部書込みがN_AまたはFAILEDでもlocal証跡を保持できる。共有memory完了はingestとevalの両方がSUCCESSの場合だけ成立し、任意のexternal diary mirrorは完了条件に含めない。
