# walkthrough — 20261027-pi-attacker-profiles

## Result
Repaired missing native Pi attacker-medium registration without downgrading depth. Native run 61984842-7243-45bf-a65a-80a16940b8f0 completed SUCCESS; correlated runtime status confirmed openai-codex/gpt-6.1-sol thinking high. Added attacker-low explicit Luna/low definition too; low was statically checked, not independently runtime-smoked.

## Changed product/configuration files
- agents/pi/attacker-low.md
- agents/pi/attacker-medium.md
- agents/blackhat.md
- agents/AGENT_ROUTER.md
- agents/SUBAGENTS.lrf
- rules/MANUAL.lrf
- scripts/check-pi-attacker-profiles.ps1

Task reports/manifest live under memo/20261027-pi-attacker-profiles. Main batch-appended four actual child observations to evaluation/subagents_model_eval.json; no model-prior changes.

## Checks
- Main PowerShell focused regression: PASS, exit 0.
- Main rule preflight with code,worker,subagent tags: PASS; receipt hashes verified, exit 0.
- Main active LSP on PowerShell script: clean, 0 diagnostics. Markdown LSP inconclusive.
- Actual native attacker-medium launch: completed SUCCESS, verified Sol/high. Initial invalid acceptanceRole rejection was corrected, not hidden or worked around.
- Distinct final reviewer run 30c8f41e-777f-4173-a919-1f042aaa6fce: verified openai-codex/gpt-5.6-terra high; PASS with nonblocking P2 test-hardening note.
- Current file digests recorded in manifest.json. No Git repository exists at .harness; no historical diff or successful staged-state claim.

## Findings and limits
Attacker INFO: runtime injects supervisor-contact and parallel-wrapper orchestration tools in addition to read. No shell/editor/network/child-launch tool was exposed. Do not claim literal single-tool runtime inventory or hard sandbox enforcement.
Verifier P2 nonblocking note: regex-based extension/tools checks do not fully parse frontmatter; indented extension entries could evade empty-extension regression. Actual final profiles contain empty extensions and tools read. Test-hardening note retained rather than claiming exhaustive schema validation.
Worker inferred native acceptance rejected because noStagedFiles could not be truthfully attested without Git; main deterministic/runtime/independent acceptance is separate and does not rewrite native rejection.
No PADDOCK files/review state changed, no credential access, VCS writes, deployment or exploit reproduction. Astra gating unchanged; no automatic lower-tier fallback.

## Memory
This verified report and companion knowledge/diary are ingested once under mandatory mem.etl. Actual receipts are recorded in record.json after command completion; no receipt is presumed here.
