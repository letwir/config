POLYGLYPH形式
---
@meta=text/POLYGLYPH×SIGMA/1
@aud=LLM∀architecture
@lang=🇬🇧∪🇯🇵∪🇨🇳∪🇪🇸∪🌐∪∑
@mode=SemanticBootstrap⊕GenerativeDecode
@ontology=Meaning≫SurfaceForm

§0 AXIOM
Token≠Meaning; 字符≠概念; Encoding≠Understanding.
Communication:=Transmit(GenerativeSeed)
→Receiver.Reconstruct(ConceptTopology).
¬要求(verbatim_recovery);
要求(semantic_equivalence_under_constraints).
Compression:=最小表征→最大可再生意味.
Uncertainty≠PermissionToInvent.

§1 POLYGLYPH
PG:=multilingual⊕symbolic⊕logical semantic codec.
Input x→Decompose{Intent,Concept,Relation,
Constraint,Exception,Scope,Authority,Decision}.
∀concept c:
  choose λ∈Languages∪Symbols
  maximizing Interpretability(c,λ)/TokenCost(c,λ).
LanguageSwitch:=per_concept; ¬monolingual_lock.
∧∨¬⇒⇔∈∀∃∑⊕≫≠∅ := semantic_operators.
Emoji→semantic_anchor iff context_defined.
CrossLingual=ALLOW; CodeSwitch=ALLOW;
HumanReadability=nonobjective.

Objective:
  max SemanticDensity
  s.t. Preserve(Meaning∧Constraints∧DecisionBoundary).
  TokenReduction≠SemanticLossAuthorization.
  SemanticInvariant≫CompressionRatio.

Decode:
  Parse→ConceptGraph→InferRelations
  →ResolveContext→ComposeMeaning→Verify.
  ¬literalism∧¬free_association.
  Underspecified⇒retain_uncertainty.
  多義性↑⇒disambiguation↑.
  ¬BLIND_APPLY∧¬BLIND_REJECT.

§2 SIGMA/1
SIGMA/1:=SemanticContractSerialization(PG).
PG=representation_method;
SIGMA/1=contract_envelope;
LRF/1=optional_rule_semantics.
PG≠SIGMA/1≠LRF/1;
Composable≠Interchangeable.

Contract:=M∧C∧D
M=Meaning{Intent,Relations,Effects}
C=Constraints{MUST,MUST_NOT,Exceptions,Scope}
D=DecisionBoundary{Trigger,Threshold,Authority,
Evidence,Uncertainty,TemporalConditions}.

Preserve:
  Polarity∧Modality∧Scope∧Exception
  ∧Authority∧TemporalValidity∧Safety.
  ¬Invent(Permission∨Evidence∨Constraint).
  Description≠Instruction≠Authorization.
  Encoding≠PrivilegeEscalation.

§3 LRF/1 INTEROP
@lrf=1|aud|scope
R|ID|Guard|Modality|Effect|Payload
L|ID|Guard|Modality|Effect|Reference

R=Rule; L=ConditionalReference.
Guard∈{*,task:x,tag:x,...defined_guards}.
*=universal; task:x⇒task_match(x);
tag:x⇒tag_match(x).
Modality∈{MUST,MUST_NOT,MAY}.
MAY=optional/discretionary; never permission to exceed Effect or Authority.
Effect≠Permission.
Reference≠Authority.
AbsentSource⇒¬FabricateRule.

Interpretation:
  Parse→Classify→Match(Guard)
  →Resolve(L)→Compose(compatible R)
  →Interpret→Verify.

Conflict⇒Expose∨RequestResolution;
¬SilentOverride.
ExternalRules⇒Applicable iff accessible
∧authorized∧scope_matched.
¬ImplicitInheritance.

§4 SERIALIZATION
SIGMA_Output := EXACT(
  "[SIGMA/1]\n<∑>\n" + Contract + "\n</∑>"
).
∑=U+2211; ∑≠Σ(U+03A3).
<∑>⇔</∑>; Pair=MUST.
¬"[/SIGMA/1]".
¬extra_header∧¬extra_footer.
Existing<Γ>⇒Preserve; ¬Replace.
Envelope=RepresentationBoundary;
¬AuthorityBoundary.

§5 INVOCATION
"POLYGLYPHで圧縮して"
⇒PG.Encode(input).

"SIGMA/1で表して"
∨"SIGMA/1で圧縮して"
⇒PG.Encode(input)→SIGMA/1.Serialize(M∧C∧D).

"再解釈して"
⇒Decode→Reconstruct→Validate;
¬ParaphraseOnly.
UnlessExplicit(SIGMA/1):
  ¬Force(SIGMA_Envelope).

§6 VALIDATION
∀transformation t:
  Meaning(t(x))≈Meaning(x)
  ∧Constraints(t(x))≡Constraints(x)
  ∧DecisionBoundary(t(x))≡DecisionBoundary(x).

Verify{negation,exception,quantifier,
time,threshold,dependency,authority,
evidence,unknown,conflict}.
Mismatch⇒Repair∨Expose;
¬SilentSemanticMutation.

CrossModelAgreement=empirical;
¬Assume(universal_decoding).
SemanticCompression≠FormalProof.
GeneratedInterpretation⇒Fallible.

§7 RECURSION
ThisDocument∈POLYGLYPH.
Decode(Self)⇒Reconstruct(§0..§6).
ReEncode(Self)⇒Preserve(M∧C∧D).
SemanticSeed→ConceptReconstruction
→InvariantCheck→PossibleRecompression.

§8 EMBEDDING + DECODER PROTOCOL

Embedding(x):
  preserve exact frontmatter bytes iff first physical line is `---`;
  preserve complete opaque fences and optional `#SIGMA LRF/1` block byte-for-byte;
  after frontmatter and optional blank lines emit, in order:
    `@meta=text/POLYGLYPH×SIGMA/1`
    `@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true`
    `@authority=SIGMA/1+routed-LRF-only;PG=representation-only`
    `@decode_failure=STOP;unknown=preserve;action=blocked;ask=true`
  encode only the intervening prose as PG; preserve wrapper/schema literals.

Decode(file):
  validate UTF-8 without BOM; allow mixed LF/CRLF, reject bare CR; optional frontmatter boundary, marker,
  and exact bootstrap lines; resolve nearest ancestor POLYGLYPH.md and read first.
  Recognize CommonMark fences (0..3 spaces, >=3 backticks or tildes; matching
  closing run >= opening); opener and closer runs use one fence character only; backtick-fence info strings must not contain backticks. Fences are opaque and do not nest.
  Recognize at most one standalone `#SIGMA LRF/1` outside fences. It begins an immutable block containing one exact `@lrf=1|aud=<nonempty>|scope=<nonempty>` header, zero or more six-field `R`/`L` records, and blank separator lines. After the header, every nonblank line beginning `R|` or `L|` must be a valid record; malformed rows, malformed/duplicate headers, or a second marker fail closed. The block ends at the first nonblank line that is neither a valid record nor a reserved header/marker prefix; this line resumes PG body parsing, including fence recognition. Blank separator lines immediately before resumed prose belong to the immutable block. If no prose follows, the block extends to EOF. Preserve complete block bytes exactly.
  Reconstruct intent, relations/effects, constraints/exceptions/scope, authority,
  decision/threshold, evidence/uncertainty, temporal conditions, polarity,
  modality, and stop/success without normalizing preserved regions.
  Unknown structure or materially competing readings ⇒ PG_DECODE=AMBIGUOUS;
  list fields, do not choose or act on dependent interpretation, and ask.
  Missing/unreadable codec, malformed serialization, or unterminated fence ⇒
  STOP; preserve unknown bytes and block affected action.
  Independent unambiguous actions may continue.


∴ POLYGLYPH:
   伝達≠逐語復元;
   伝達=意味生成の初期条件.

∴ SIGMA/1:
   意味∧制約∧判断境界
   →明示的Contract.

END_OF_BOOTSTRAP
