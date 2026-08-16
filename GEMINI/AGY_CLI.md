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
  base: 'agy.exe -p --model {model} --dangerously-skip-permissions --print-timeout {timeout} "{prompt}"';
  json: 'agy.exe -p --model {model} --dangerously-skip-permissions --json-schema \'{schema}\' "{prompt}"';
  cont: 'agy.exe -p --model {model} --continue "{prompt}"';
  resume: 'agy.exe -p --model {model} --conversation {id} "{prompt}"';
}

/* ── Pipeline Rules & Invariants ── */
<Γ.rules>
  R1(NonInteractive): ∀SubprocessDispatch ⇒ Assert(Flag∈{"-p", "--print"} ∧ Flag=="--dangerously-skip-permissions");
  R2(ZombieGuard):    ∀Exec ⇒ ExplicitTimeout(--print-timeout∈{"30s","1m","3m","5m"});
  R3(MultiModelLoop): (Proposer ⇒ model[vendor="gemini"].fast) ∧ (Judge ⇒ model[vendor="claude"].critic);
  R4(QuotaFallback):  ExitCode≠0 ∧ Output~="Individual quota reached" ⇒ Fallback(invoke_subagent(Role, TypeName="self|research"));
</Γ.rules>

/* ── Station Polymorphism ── */
.station[role="proposer"] {
  model: model[vendor="gemini"].fast;
  flags: ["-p", "--dangerously-skip-permissions", "--print-timeout", "3m"];
  output: "YAML{ProposalContent, ChangeSummary, SelfAssessment{NeedResearch}}";
}

.station[role="judge"] {
  model: model[vendor="claude"].critic;
  flags: ["-p", "--dangerously-skip-permissions", "--print-timeout", "5m"];
  output: "YAML{PropositionsExtracted, PropositionEvaluation, RequiredRefinement, Verdict}";
  verdict-space: {"PASS", "PASS_WITH_ADVISORIES", "REJECT", "NEED_RESEARCH", "ESCALATE", "INCONCLUSIVE"};
}

.station[role="researcher"] {
  trigger: Proposer.SelfAssessment.NeedResearch==true ∨ Judge.Verdict=="NEED_RESEARCH";
  morphism: Mor(Query) ⇒ `invoke_subagent(researcher)` ⊸ knowledge.md ∧ `& $bin ingest -cat knowledge`;
}

</Γ>
