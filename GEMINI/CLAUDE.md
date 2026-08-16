[SPR/XML::ρ→max|target:CLAUDE(model≠Gemini*)|protocol:hydrate(latent_concepts)⇒exec] <Γ id="claude_entry">

/* ── なぜここへ分岐されたか ──────────────────────────────────────────
   GEMINI.md L5: dispatch: model ∉ {Gemini*} ⇒ view("./CLAUDE.md")
   理由: extended_thinking の diary 密度・思考プロトコルが Gemini と異なるため
   GEMINI.md との diff: [thinking="standard"] → [thinking="extended"] のみ
   ──────────────────────────────────────────────────────────────────── */

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   BootSequence — このファイルを読んだら即時・無条件に以下を実行せよ
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
<Γ.boot>
  Step 1: view_file("$env:USERPROFILE/.gemini/MACHINE.toml")
          /* env/hardware/shell/encoding — 実行環境 single source of truth */
  Step 2: view_file("$env:USERPROFILE/.gemini/PERSONA.css")
          /* persona/paradigm/skill_map/state-policy/self-evolve/archive */
          /* PERSONA.css は内部で subagents/SUBAGENTS.css を @import している */
          /* → SUBAGENTS.css は PERSONA.css 読み込みで連鎖的に有効化される   */
  Step 3: view_file("$env:USERPROFILE/.gemini/subagents/SUBAGENTS.css")
          /* subagent 定義の明示的な確認 (cascade 保証のため二重読み) */

  ¬skip: Boot失敗∨省略 ⇒ 設定不完全 → 以降の判断に誤りが生じる ∴ 必須
</Γ.boot>

/* ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   Claude 固有オーバーライド (PERSONA.css .agent の thinking 属性を上書き)
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ */
<Γ.override>
  .agent { thinking: extended; }
  /* rule: 複雑推論∨設計判断 ⇒ extended_thinking 必須, ¬思考省略    */
  /* priority: 仮説検証 → 結論導出 → 実行                           */
  /* diary: forced-append-per-turn (毎ターン末尾に diary.md へ追記)  */
  /*        attribution: [ワイの指示(PromptDefect):%] vs [AI認知(AgentDefect):%] */
</Γ.override>

</Γ>
