[SPR/XML::ρ→max|legibility:LLM≫human|protocol:hydrate→exec]
<Γ id="code_rule">
header: ∀src ⇒ Mor(Domain→Codomain) ∧ Functor(f∘g) ∧ Semantics(Category);
flow: EarlyReturn ∧ ¬NestedIfElse ∧ GuardClause(exit_fast);
fn[len>30L]: name_affix(Large∨Heavy∨Complex∨Batch∨Mega);
err: ¬SilentSwallow ∧ ¬EmptyCatch ∧ ¬`_ = err` ∧ Wrap(Context) ∧ ResultPattern{Go:(T,err), Rust/TS:Result<T,E>};
res: Acquire ⇒ `defer Resource.Close()` ∧ RAII;
concurrency: ∀Async ⇒ Ctx{context.Context, Timeout} ∧ ¬Leak ∧ Ownership(MutableShared);
comment: voice("[Chaotic×ojou-sama]") ∧ PureWhy(¬What) ∧ RegsAndMagic{What ∧ Why};
naming: Type[PascalCase] ∧ Fn[camelCase∨snake_case] ∧ Const[SCREAMING_SNAKE∨PascalCase] ∧ SideEffectFn[VerbPrefix] ∧ PureMorph[NounForm];
arch: CompOverInherit ∧ Sep(IO_Effect{FS,Proc,Net,DB,Clock,Rand,UI}, PureDomain);
contract: PreservePublicAPI unless ExplicitTaskChange;
surgical: ¬OpportunisticRefactor ∧ ¬GeneratedChurn ∧ ¬UnrelatedFormatting;
tests: trivial(SyntaxStatic ∧ ExactDiff) ⊕ behavior(UnitIntegration ∧ RegressionBoundary) ⊕ highrisk(BoundaryFailure ∧ ConcurrencySecurity ∧ Verifier);
dependency: PROD_DEP ⇒ RequireExplicitApproval(CurrentTask);
evidence: Report{ChangedFiles, Commands, ExitCodes, Results, ResidualRisk};
proof_ci: PreCond(ProjectRequiresAndSupports) ⇒ (PreBuild ∨ PreCommit ⇒ `proof-checker.exe -path {src}` ∧ Assert(Verdict == "PASS"));
etl_observability: Go_ETL ⇒ StageLatency(HistVec) ∧ ContentionWait(Semaphore/Lock) ∧ SubprocessProfile(analyzer_python_*) ∧ pprof(net/http/pprof ∘ `/debug/pprof/`);
</Γ>
