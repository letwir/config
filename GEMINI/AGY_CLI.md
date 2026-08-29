[SPR/XML::ρ→max|target:{GEMINI,CLAUDE}|legibility:LLM≫human|protocol:hydrate(latent_concepts)⇒exec]
<Γ id="agy_cli">

@import "./MACHINE.toml";
@import "./PERSONA.css";

/* ── Environment & Binary Specification ── */
:root {
  --bin: "C:/Users/letwir/AppData/Local/agy/bin/agy.exe";
  --cmd: "agy.exe";
  --cfg-dir: "$env:USERPROFILE/.gemini/antigravity-cli";
  --default-timeout: "5m0s";
}

/* ── Model Taxonomy & Morphism Mapping ── */
model[vendor="gemini"] {
  fast: "gemini-3.7-flash";
  pro:  "gemini-2.5-pro";
  lite: "gemini-2.5-flash";
  role: "Proposer ∨ Researcher ∨ FastDiscovery";
}

model[vendor="claude"] {
  critic: "claude-sonnet-4-6";
  arch:   "claude-opus-4-6";
  light:  "claude-haiku-4-5";
  role:   "LLM-Judge ∨ AdversarialCritic ∨ DecisionsAuditor";
}

/* ── Command Invocations (Non-Interactive CLI) ── */
.dispatch {
  base: 'agy.exe --model {model} --dangerously-skip-permissions --print-timeout {timeout} -p "{prompt}"';
  json: 'agy.exe --model {model} --dangerously-skip-permissions --json-schema \'{schema}\' -p "{prompt}"';
  cont: 'agy.exe --model {model} --continue -p "{prompt}"';
  resume: 'agy.exe --model {model} --conversation {id} -p "{prompt}"';
}

/* ── Pipeline Rules & Invariants ── */
<Γ.rules>
  R1(NonInteractive): ∀SubprocessDispatch ⇒ Assert(Flag∈{"-p", "--print"} ∧ Flag=="--dangerously-skip-permissions");
  R2(ZombieGuard):    ∀Exec ⇒ ExplicitTimeout(--print-timeout∈{"30s","1m","3m","5m"});
  R3(MultiModelLoop): (Proposer ⇒ model[vendor="gemini"].fast) ∧ (Judge ∨ Auditor ∨ Verifier ⇒ model[vendor="claude"].critic);
  R4(QuotaFallback):  ExitCode≠0 ∧ (Output~="quota" ∨ Output~="rate limit" ∨ Output~="token") ⇒ Fallback(agy.exe[gemini-3.7-flash --effort high] ∨ invoke_subagent(Role, TypeName="self|research"));
  R5(AdvisoryHandling): PASS_WITH_ADVISORIES ⇒ Worker must apply/incorporate advisories (¬スルー) ∧ `llm-mem.exe ingest -cat advisory`;
</Γ.rules>

/* ── Station Polymorphism ── */
.station[role="proposer"] {
  model: model[vendor="gemini"].fast;
  flags: ["--model", "gemini-3.7-flash", "--effort", "high", "--dangerously-skip-permissions", "--print-timeout", "3m", "-p"];
  output: "YAML{ProposalContent, ChangeSummary, SelfAssessment{NeedResearch}}";
}

.station[role="auditor"] {
  model: model[vendor="claude"].critic;
  flags: ["--model", "claude-sonnet-4-6", "--dangerously-skip-permissions", "--print-timeout", "5m", "-p"];
  script: "python.exe scripts/auditor_claude.py --plan {plan} --decisions decisions.md --rule CODE_RULE.md";
  output: "YAML{DiscoveryLog, SpecDriftCheck, CTSoundness, EffectPurityAudit, Prognosis, Defects, CycleCount, Verdict}";
  verdict-space: {"PASS", "PASS_WITH_ADVISORIES", "REJECT", "ESCALATE"};
}

.station[role="verifier"] {
  model: model[vendor="claude"].critic;
  flags: ["--model", "claude-sonnet-4-6", "--dangerously-skip-permissions", "--print-timeout", "5m", "-p"];
  script: "python.exe scripts/verifier_claude.py --target {target} --rule CODE_RULE.md";
  output: "YAML{DiscoveryLog, StaticAnalysis, ExecutionTests, AdversarialAttacks, EvidenceLedger, DefectsFound, FeedbackToWorker, CycleCount, Verdict}";
  verdict-space: {"PASS", "PASS_WITH_ADVISORIES", "REJECT", "INCONCLUSIVE", "ESCALATE"};
}

.station[role="judge"] {
  model: model[vendor="claude"].critic;
  flags: ["--model", "claude-sonnet-4-6", "--dangerously-skip-permissions", "--print-timeout", "5m", "-p"];
  output: "YAML{PropositionsExtracted, PropositionEvaluation, RequiredRefinement, Verdict}";
  verdict-space: {"PASS", "PASS_WITH_ADVISORIES", "REJECT", "NEED_RESEARCH", "ESCALATE", "INCONCLUSIVE"};
}

.station[role="researcher"] {
  trigger: Proposer.SelfAssessment.NeedResearch==true ∨ Judge.Verdict=="NEED_RESEARCH";
  morphism: Mor(Query) ⇒ `invoke_subagent(researcher)` ⊸ knowledge.md ∧ `llm-mem.exe ingest -file knowledge.md -cat knowledge`;
}

</Γ>
