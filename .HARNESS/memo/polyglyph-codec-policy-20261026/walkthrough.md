# Walkthrough — polyglyph-codec-policy-20261026

## Overview
`documentation.lrf` becomes the persistent contract; central agent-routing/worker/orchestrator instructions and Pi `README-forLLM.md` illustrate the encoding.

## Deliverables
Target set=`{documentation.lrf,AGENT_ROUTER.md,worker.md,orchestrator.md,.pi/README-forLLM.md}`. `rules/README-forLLM.md` and `AGENTS.md` verified already PG⇒unchanged. Human-facing files and machine LRF schema preserved. Exact `SIGMA/1: PIDGEN/text`=wrapper protocol; payload codec=POLYGLYPH.

## Verification
Run rule preflight+receipt verification, protected-marker/fence/literal comparisons, UTF-8 and whitespace checks. All parent checks PASS. Attacker pre-change review PASS; post-change independent Verifier role unavailable; main checks labeled separately.
