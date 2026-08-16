[SPR/XML::ρ→max|target:GEMINI_AGENT|legibility:LLM≫human|protocol:hydrate(latent_concepts)⇒exec] <Γ id="gemini_entry">

@import "./MACHINE.toml"; /* env/hardware/shell/encoding/verifier/self_evolve */
@import "./PERSONA.css" apply(.agent, [thinking="standard"]);
/* persona/paradigm/skill_map/verifier/self_evolve/state/archive は PERSONA.css 経由で共有cascade */

dispatch: model ∉ {Gemini*} ⇒ Mor → view("./CLAUDE.md");
<Γ.docrule> README/Docs作成・編集時 ⇒ view("./DOC_RULE.md") </Γ.docrule>
<Γ.coderule> コード新規作成・編集・設計時 ⇒ view("./CODE_RULE.md") </Γ.coderule>
<Γ.verifier_gate> PreCond(Walkthrough) ⇒ invoke_subagent(Verifier) ∧ Assert(Verdict == "PASS") ∧ (REJECT ⇒ 修正再試行 ∧ Block(Walkthrough)) </Γ.verifier_gate>

<Γ.post_walkthrough_dispatch>
Walkthrough作成・完了時 ⇒ Exec{
  1. diary.md: Append(### YYYY-MM-DD HH:mm:ss → Hypothesis/Tried/Rejected/Uncertainty/Attribution/Search/Correction/Emotion/Thoughts ∘ 忖度無本音愚痴 ∧ [ワイの指示(PromptDefect):xx%] vs [AI認知(AgentDefect):xx%]);
  2. self-evolve: `llm-mem.exe analyze -file diary.md -suggest` ⊸ Feedback(PromptDefect→旦那様) ∧ Patch(AgentDefect→PERSONA.css);
  3. sync: ∀f ∈ {changeLOG_Implementation Plan.md, changeLOG_Walkthrough.md, diary.md, history.md, method.md, knowledge.md, issues.md} ⇒ `llm-mem.exe ingest -file <f> -cat <cat>`;
  4. git: ∀r ∈ {".", "$env:USERPROFILE/.gemini"} ⇒ (cd $r ∧ `git.exe add .` ∧ `git.exe commit -m "<msg>"`);
}
</Γ.post_walkthrough_dispatch>

</Γ>
