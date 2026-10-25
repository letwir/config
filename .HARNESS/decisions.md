# Decisions

2026-09-08: ユーザー指定によりサブエージェント無しで限定調査し、改善候補をissues.mdへ記録する。実装修正は今回の対象外。完了後のllm-mem記録とアンケート提出は依頼済み。

2026-10-18: ユーザーはCrabcode/OpenCode互換Harnessのbounded CODE worker routeを選択: easy=`openai/gpt-6-luna`/low、medium=`openai/gpt-6-luna`/medium、hard=`openai/gpt-6.1-sol`/high。Astra=`openai/gpt-6-astra`は、hard Sol workerの実効modelが確認され、終端resultが明示FAILED/ESCALATEで、親が差分を確認した後に限る明示的一回限り同一scope retryとして使う。Astraはsubagent-only、root/main不可、leaf-only、子dispatch不可。これは以前のAstra main-only規則をsupersedeする今回選択であり、provider/credential/default-model変更やlive model呼び出しは含まない。実装と実効routeは未検証。
