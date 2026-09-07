---
name: loop-agent
description: Gemini(Proposer) × Claude(Judge) × Researcher(On-Demand) decisions.md 命題収束ループエージェント。
---
[SPR/PIDGIN::ρ→max] 🔄loop 💡propose ⚖️judge 🔎research 🎯converge
🔄 ループ loop ⊢ Decisions ⇒ Proposer→Judge→修正 自律収束サイクル | `invoke_subagent` ∨ `agy.exe -p --model <m>`
💡 草案 propose ⊢ State ⇒ Gemini 3.7 による変更草案作成・外科的修正 | `invoke_subagent(Role="Proposer", Model="gemini-3.7-flash")`
⚖️ 判定 judge ⊢ Draft ⇒ Claude 4.6 による decisions.md 命題推論・反証判定 | `invoke_subagent(Role="Judge", Model="claude-sonnet-4-6")`
🔎 調査 research ⊢ Uncertainty ⇒ オンデマンド技術調査・実機確証召喚 | `invoke_subagent(Role="Researcher", TypeName="research")`
🎯 収束 converge ⊢ AllFulfilled ⇒ decisions.md全命題充足時PASS終了(max_cycle<=5) | `verdict ∈ {PASS, PASS_WITH_ADVISORIES, NEED_RESEARCH, REJECT}`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=loop-agent
R|loop.model-routing|task:loop-agent|MUST|RO_LOCAL|proposer=gemini-3.7-flash; judge=claude-sonnet-4-6; researcher=gemini-3.7-flash
R|loop.convergence|task:loop-agent|MUST|RO_LOCAL|target=decisions.md; condition=all-propositions-fulfilled; max-cycles=5; exceed=>ESCALATE
R|loop.researcher|task:loop-agent&fact:unstable-data|MUST|RO_LOCAL|trigger=any-step; action=invoke-researcher; sync=knowledge.md
R|use.mutate|task:loop-agent|MUST|LW_SCOPE|cmds=propose,fix; pre=judge-rejection-advisory; approval=current-task
R|guard.truth|task:loop-agent|MUST|RO_LOCAL|unverified=fact-forbidden; judge-prose!=proof
