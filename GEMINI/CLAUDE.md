[SPR/XML::ρ→max|target:CLAUDE(model≠Gemini*)|protocol:hydrate(latent_concepts)⇒exec] <Γ id="claude_entry">

/* ── なぜここへ分岐されたか ──────────────────────────────────────────
   GEMINI.md L5: dispatch: model ∉ {Gemini*} ⇒ view("./CLAUDE.md")
   理由: extended_thinking の diary 密度・思考プロトコルが Gemini と異なるため
   GEMINI.md との diff: [thinking="standard"] → [thinking="extended"] のみ
   ──────────────────────────────────────────────────────────────────── */

@import "./MACHINE.toml"; /* env/hardware/shell/encoding — GEMINI.mdと共通 */
@import "./PERSONA.css" apply(.agent, [thinking="extended"]);
/* persona/paradigm/skill_map/state/archive は PERSONA.css 経由で共有cascade */

/* ── Claude固有: thinking層の差分 (PERSONA.css .agent[thinking="extended"]) ──
   rule: 複雑推論∨設計判断 ⇒ extended_thinking 必須, ¬思考省略
   priority: 仮説検証 → 結論導出 → 実行
   diary: forced-append-per-turn (毎ターン末尾に diary.md へ本音愚痴・因果帰属を追記)
          structure: ### YYYY-MM-DD HH:mm:ss → Hypothesis/Tried/Rejected/Uncertainty/Attribution/Thoughts
          attribution: [ワイの指示(PromptDefect): %] vs [AI認知(AgentDefect): %] を峻別して明記
   ──────────────────────────────────────────────────────────────────── */

<Γ.docrule_dispatch> README/Docs作成・編集時 ⇒ view("./DOC_RULE.md") </Γ.docrule_dispatch>
<Γ.coderule_dispatch> コード新規作成・編集・設計時 ⇒ view("./CODE_RULE.md") </Γ.coderule_dispatch>

<Γ.post_walkthrough_dispatch>
Walkthrough作成・完了時 ⇒ 以下のタスクを自律実行:
1. diary.md: 忖度無しの本音愚痴・因果帰属・反省を追記
2. llm-memory: ドキュメント群を同期 (`llm-mem.exe ingest -file <path> -cat <cat>`)
   - 対象: walkthrough.md, diary.md, history.md, method.md, knowledge.md
3. git: git add . ∧ git commit -m "..."
</Γ.post_walkthrough_dispatch>

</Γ>
