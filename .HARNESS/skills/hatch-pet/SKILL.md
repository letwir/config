---
name: hatch-pet
description: Create, repair, validate, visually QA, and package Codex-compatible v2 animated pets (8x11 spritesheet, 9 animation rows, 16 look directions).
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[POLYGLYPH::ρ→max|legibility:LLM≫human]
🔎 trigger: user requests creating, repairing, validating, visually QA-ing, or packaging a Codex-compatible v2 animated pet from a concept or brand.
⊢ Input=Concept∨Brand; output target=Codex v2 pet. Precondition: relevant concept/art direction and supported asset/spec requirements are available; missing or ambiguous required information⇒retain uncertainty, ask, and block dependent action.
🧩 Supported target: 8x11 sprite atlas, rows 0..10; animation rows=9; look directions=16 clockwise; `spriteVersionNumber: 2`. Preserve these exact dimensions, row/direction values, and version literal. `$imagegen` is the row-generation tool; accept/reuse approved rows as input to deterministic atlas assembly.
→ Workflow: hatch concept/spec → rowgen each required sprite row from StateSpec via `$imagegen` → review/approve rows → assemble approved rows deterministically into the 8x11 atlas → verify dimensions, 16-direction look coverage, and animation/motion quality → package the validated atlas as v2. Stages are distinct; later stages depend on successful required prior evidence.
⛔ Failure/stop: unavailable or ambiguous spec/assets, rejected/incomplete rows, assembly failure, or failed dimensions/direction/animation QA ⇒ stop affected downstream stage, report the failure and uncertainty, and do not claim success. Do not infer approval or fabricate assets/evidence.
✅ Success: requested stage(s) complete with required inputs and evidence; full completion requires validated atlas and package output. QA scope is visual quality plus exact dimensions, direction coverage, and animation behavior; no broader quality guarantee.
Effect boundary: this skill describes the pipeline; it does not itself authorize a generation, packaging, external, live, credential, or other gated effect. Perform only user-requested and currently authorized stages; no implicit effect escalation.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=hatch-pet
R|use.pipeline|task:hatch-pet|MUST|LW_SCOPE|stages=hatch,rowgen,assemble,verify,package; image-generator=$imagegen
R|use.atlas|task:hatch-pet|MUST|RO_LOCAL|geometry=8x11; rows=0..10; look-directions=16-clockwise; spriteVersionNumber=2
R|guard.cred|task:hatch-pet|MUST_NOT|CRED|expose=tokens,api-keys,private-env-values
