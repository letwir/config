# knowledge | 20260923-subagent-routing-design

## Facts
- Current bootstrap makes LRF the policy authority and README navigation only.
- Current ETL already stages the workflow: main plans, agy executes CODE, and permitted Codex roles can research/review. `SUBAGENTS.lrf` is currently the loaded role contract.
- Runtime probe H-12 verified root-to-subagent dispatch but the subagent could not dispatch a child; runtime reported root-thread-only dispatch.
- Current working tree already contains extensive unrelated edits, including `etl/main.seq`, `agents/SUBAGENTS.lrf`, nine TOMLs, and workflow files.

## Design inference
- Keep main.seq as the overall workflow; move only role selection predicates and role-specific branching to one canonical agents selection contract, then call that selector from main at relevant gates.
- Use a normative `agents/selection.lrf` plus a concise `agents/selection.md` index if a human-facing selection guide is useful; LRF remains authoritative.
- Define Explorer for local repository mapping and Researcher for current public facts/API sources/precedents. Select one by task need; allow both only for independent workstreams.
- Current mandatory agy CODE policy conflicts with making `worker.toml` an active Codex code-writing role. Recommended Worker contract is for the bounded agy CODE handoff; retire or repurpose the Codex Worker TOML so it cannot compete with that route.
- Remove all `proactively` wording; TOML descriptions should state trigger, exclusions, input, and output, while selection LRF owns dispatch rules.

## Unknown / decision
- Whether to retire `worker.toml` or repurpose it as a read-only agy-task packager needs user confirmation before implementation.
