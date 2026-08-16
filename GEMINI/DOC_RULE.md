[SPR/XML::ρ→max|legibility:LLM≫human|protocol:hydrate→exec]
<Γ id="doc_rule">
readme: H1(Title) ∘ LangNav(🇺🇸/🇯🇵) ∘ H2("ナニコレ？":≤3L) ∘ H2("必要なもの":Req) ∘ H2("使い方":Usage) ∘ H2("概要詳しく":Detail∘Table(docs/*.md)) ∘ H2("状態図":1L∘ref(docs/state_diagram.md)) ∘ H2("ER図とデータ構造":1L∘ref(docs/database_er_diagram.md)) ∘ H2("SHM/WORM":1L∘ref(docs/shm_architecture.md)) ∘ H2("ライセンス":License∘Warning);
doc_sep: LargeDoc{Mermaid,ER,SHM} ⇒ Mor(docs/*.md) ∧ README(1L_summary ∘ ref(docs/*.md));
cfg_mgmt: Track(*.example) ∧ Ignore(real_cfg{config.toml,.env,settings.json}) ∧ Tracked(real_cfg)⇒`git rm --cached` ∧ Secret⇒`git filter-repo` ∧ Detect(*.ts.net∨cred)⇒`git filter-repo --invert-paths --path config.toml --force`∘gen_example;
mega_docs: (FeatureAdd∨Change) ⇒ PreCond(`sec-forensics.exe mega-check` ∧ invoke_subagent(blackhat) ∧ Assert(OverallVerdict == "PASS")) ∧ Tag(`mega-docs-update-YYYYMMDD`) ∧ Msg(`docs(mega-docs-update): <summary>`) ∧ Log(`git.exe log mega-docs-update-YYYYMMDD..HEAD --oneline`) ∧ Sync(decisions.md∨history.md);
walkthrough_sync: PreCond(invoke_subagent(Verifier) ∧ Assert(Verdict == "PASS")) ⇒ Mor(Walkthrough{Title, L0_Raw, L1_Points, L2_OneLiner, Tags}) ⇒ `llm-mem.exe add -cat walkthrough -title <Title> -content <L0> -l1 <L1> -l2 <L2> -tags <Tags>` ∘ SyncRemotePostgres(llm_memory);
</Γ>
