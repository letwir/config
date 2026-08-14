[SPR/XML::ρ→max|legibility:LLM≫human|protocol:hydrate→exec]
<Γ id="code_rule">
header: ∀src ⇒ Mor(Domain→Codomain) ∧ Functor(f∘g) ∧ Semantics(Category);
flow: EarlyReturn ∧ ¬NestedIfElse ∧ GuardClause(exit_fast);
fn[len>30L]: name_affix(Large∨Heavy∨Complex∨Batch∨Mega);
err: ¬SilentSwallow ∧ ¬EmptyCatch ∧ ¬`_ = err` ∧ Wrap(Context) ∧ ResultPattern{Go:(T,err), Rust/TS:Result<T,E>};
res: Acquire ⇒ `defer Resource.Close()` ∧ RAII;
concurrency: ∀Async ⇒ Ctx{context.Context, Timeout} ∧ ¬Leak;
comment: voice("[Chaotic×ojou-sama]") ∧ PureWhy(¬What) ∧ RegsAndMagic{What ∧ Why};
naming: Type[PascalCase] ∧ Fn[camelCase∨snake_case] ∧ Const[SCREAMING_SNAKE∨PascalCase] ∧ SideEffectFn[VerbPrefix] ∧ PureMorph[NounForm];
arch: CompOverInherit ∧ Sep(IO_Effect, PureDomain);
</Γ>
