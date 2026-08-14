[SPR/XML::ρ→max|legibility:LLM>human|protocol:hydrate→exec]
<Γ id="coderule_specs">
readme_structure {
  seq: H1(Title) ∘ LangNav(🇺🇸/🇯🇵) ∘ H2("ナニコレ？": ≤3L) ∘ H2("必要なもの": Req) ∘ H2("使い方": Usage) ∘ H2("概要詳しく": Detail ∘ DocsTable(docs/*.md)) ∘ H2("状態図": 1L ∘ ref(docs/state_diagram.md)) ∘ H2("ER図とデータ構造": 1L ∘ ref(docs/database_er_diagram.md)) ∘ H2("Windows 共有メモリ (SHM) 管理と WORM アーキテクチャ": 1L ∘ ref(docs/shm_architecture.md)) ∘ H2("ライセンス": License ∘ Warning);
  docs_separation: Mor(LargeDoc{Mermaid,ER,SHM}) → docs/*.md ∧ README ⇒ 1L_summary ∘ ref(docs/*.md);
}
config_management {
  rule[target="*.example"]: Git.track(*.example) ∧ Git.ignore(real_config{config.toml,.env,settings.json});
  rule[target="retroactive_purge"]: Git.tracked(real_config) ⇒ `git rm --cached <file>` ∘ commit; secret_in_history ⇒ `git filter-repo`;
  rule[target="private_domain_purge"]: detect(*.ts.net ∨ credential) ⇒ git_untrack ∘ gen_example(dummy_val ∘ ja_comment) ∘ `git filter-repo --invert-paths --path config.toml --force`;
}
mega_docs_update {
  anchor: Mor(FeatureAdd ∨ Change) → DocsSync ∧ GitLog.query_since_last_mega;
  tag: `mega-docs-update-YYYYMMDD`;
  commit_msg: `docs(mega-docs-update): <summary>`;
  diff_cmd: `git.exe log mega-docs-update-YYYYMMDD..HEAD --oneline`;
  record: decisions.md ∨ history.md ⇒ append_table(tag, date, summary);
}
</Γ>
