---
name: loop-agent
description: Gemini(Proposer) × Claude(Judge) × Researcher(On-Demand) decisions.md 命題収束ループエージェント。
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max] 🔄⟨loop⟩ 💡⟨propose⟩ ⚖️⟨judge⟩ 🔎⟨research⟩ 🎯⟨converge⟩
🔄 Trigger: loop-agent task ∧ decisions.md proposition work ⇒ repeat {Proposer draft/change proposal → Judge reason over propositions + attempt refutation → rejection/advisory ⇒ Proposer surgical revision → Judge}; eligible invocation path `invoke_subagent` ∨ `agy.exe -p --model <m>`; no other trigger/authority inferred.
💡 Proposer: state/context ⇒ bounded change draft or surgical fix | `invoke_subagent(Role="Proposer", Model="gemini-3.7-flash")`
⚖️ Judge: draft ⇒ evaluate decisions.md propositions, fulfillment, reasoning and counterevidence; verdict is not proof | `invoke_subagent(Role="Judge", Model="claude-sonnet-4-6")`
🔎 Research on demand: any step with unstable/uncertain fact ⇒ invoke researcher, obtain technical/real-world evidence; synchronize relevant knowledge to knowledge.md; unverified claim ≠ fact | `invoke_subagent(Role="Researcher", TypeName="research")`
🎯 Termination: all decisions.md propositions fulfilled ⇒ PASS and stop; cycle limit max_cycle<=5; limit exceeded ⇒ ESCALATE, no silent continuation. Output verdict ∈ {PASS, PASS_WITH_ADVISORIES, NEED_RESEARCH, REJECT}; unresolved evidence ⇒ NEED_RESEARCH; rejected proposal ⇒ REJECT or revise only within allowed approval boundary.
Failure boundaries: missing/failed proposer, judge, or research evidence; unresolved uncertainty; malformed/contradictory decision; or unavailable invocation ⇒ do not claim fulfillment/PASS, preserve uncertainty and surface failure/escalation. propose/fix mutations require judge rejection/advisory and current-task approval; no execution or approval inferred.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=loop-agent
R|loop.model-routing|task:loop-agent|MUST|RO_LOCAL|proposer=gemini-3.7-flash; judge=claude-sonnet-4-6; researcher=gemini-3.7-flash
R|loop.convergence|task:loop-agent|MUST|RO_LOCAL|target=decisions.md; condition=all-propositions-fulfilled; max-cycles=5; exceed=>ESCALATE
R|loop.researcher|task:loop-agent&fact:unstable-data|MUST|RO_LOCAL|trigger=any-step; action=invoke-researcher; sync=knowledge.md
R|use.mutate|task:loop-agent|MUST|LW_SCOPE|cmds=propose,fix; pre=judge-rejection-advisory; approval=current-task
R|guard.truth|task:loop-agent|MUST|RO_LOCAL|unverified=fact-forbidden; judge-prose!=proof
