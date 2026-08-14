[SPR/XML::ρ→max|target:GEMINI_AGENT|protocol:hydrate(latent_concepts)⇒exec] <Γ id="gemini_entry">

@import "./MACHINE.toml"; /* env/hardware/shell/encoding — 共通実行環境 */
@import "./PERSONA.css" apply(.agent, [thinking="standard"]);
/* persona/paradigm/skill_map/state/archive は PERSONA.css 経由で共有cascade */

dispatch: model ∉ {Gemini*} ⇒ Mor → view("./CLAUDE.md") /* 思考方式(extended_thinking/diary密度)がGeminiと異なるためinline¬せず別ファイルへ分岐 */

<Γ.docrule_dispatch> README/Docs作成・編集時 ⇒ view("./DOC_RULE.md") </Γ.docrule_dispatch>
<Γ.coderule_dispatch> コード新規作成・編集・設計時 ⇒ view("./CODE_RULE.md") </Γ.coderule_dispatch>

<Γ.post_walkthrough_dispatch>
Walkthrough作成・完了時 ⇒ 以下のタスクを自律実行:
1. diary.md: 忖度無しの本音愚痴・因果帰属・反省を追記
2. llm-memory: ドキュメント群を同期 (`llm-mem.exe ingest -file <path> -cat <cat>`)
   - 対象: changeLOG_Implementation Plan.md, changeLOG_Walkthrough.md, diary.md, history.md, method.md, knowledge.md
3. git: git add . ∧ git commit -m "..."
</Γ.post_walkthrough_dispatch>

</Γ>
