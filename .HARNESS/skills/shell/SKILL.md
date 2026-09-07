---
name: shell
description: Runs the rest of a /shell request as a literal shell command. Use only when the user explicitly invokes /shell.
disable-model-invocation: true
---
[SPR/PIDGIN::ρ→max] 🐚exec 📋status
🐚 実行 exec ⊢ LiteralCmd ⇒ 改変・前置推論なしの即時シェルコマンド実行 | `run_command(CommandLine="<user-command>")`
📋 報告 status ⊢ ExitCode ∧ Stdout ⇒ 実行ステータス・標準出力／標準エラーの簡潔報告

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=shell
R|use.exec|task:shell|MUST|RO_LOCAL|target=verbatim-user-args; rewrite=forbidden; explain-first=forbidden
R|use.ask|task:shell&case:no-args|MUST|RO_LOCAL|prompt="Which command to run?"
R|guard.truth|task:shell|MUST|RO_LOCAL|output=exit-status+stdout+stderr; fabrications=forbidden
