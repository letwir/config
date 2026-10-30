---
name: shell
description: Runs the rest of a /shell request as a literal shell command. Use only when the user explicitly invokes /shell.
disable-model-invocation: true
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🐚exec 📋status
🔎 /shell明示呼出 ∧ args≠∅ ⊢ LiteralCmd=user_args; ∴ `run_command(CommandLine="<user-command>")` を即時実行。引数改変・解釈・前置説明なし。
📋 結果 ⊢ ExitCode ∧ Stdout ∧ Stderr を簡潔に報告。

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=shell
R|use.exec|task:shell|MUST|RO_LOCAL|target=verbatim-user-args; rewrite=forbidden; explain-first=forbidden
R|use.ask|task:shell&case:no-args|MUST|RO_LOCAL|prompt="Which command to run?"
R|guard.truth|task:shell|MUST|RO_LOCAL|output=exit-status+stdout+stderr; fabrications=forbidden
