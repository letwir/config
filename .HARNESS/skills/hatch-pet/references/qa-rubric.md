@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# V2 Pet QA Rubric

Packaging blocked UNTIL ∀ sections PASS.

## Geometry And Package

- Final atlas: exactly `1536x2288`; 8 columns x 11 rows; cells=`192x208`.
- `pet.json`: contains `spriteVersionNumber: 2` ∧ points to packaged spritesheet.
- Used cells=non-empty; unused standard-row cells=transparent.
- Fully transparent pixels: RGB residue=zero.
- 8x9 intermediate atlas: NEVER package.
- `qa/review.json`: no errors.
- Standard rows: component extraction UNLESS `stable-slots` deliberately approved AFTER playback review.
- Coherent look rows: recover ordered pose groups ∧ pass near-edge clipping checks AFTER shared-scale registration into final cells.

## Character And Style

- Across all 11 rows: silhouette ∧ proportions ∧ face ∧ expression language ∧ material ∧ palette ∧ lighting ∧ markings ∧ props remain the same.
- Pet readability: clear within `192x208` cell, chosen style.
- ∀ frames: no introduction of unintended {character, object, logo, text, scene, effect}.

## Standard Animation

- Rows `0-8`: exact required frame counts ∧ recognizable state semantics.
- Loops: no popping, reversed cadence, wrong-facing direction, or remaining effectively static.
- First idle frame: works as reduced-motion still.
- `waiting` ∧ `running` ∧ `review` ∧ `failed`: remain visually distinct.

## Look Directions

- All 16 directions: present ∧ fixed clockwise order ∧ visibly distinct from neutral/rest.
- Cardinals: unmistakably up/right/down/left; diagonals ∧ intermediates: correct quadrant.
- `qa/look-directions.png`: full-body ∧ zoomed head/upper-body views.
- `qa/direction-semantics.json`: records `pass`, `expected`, `observed`, `reason` ∀ directions.
- `qa/look-continuity.json`: no unexplained holes, center jumps, area jumps, or local difference outliers.
- Eyes ∧ eyelids ∧ head ∧ body ∧ appendages ∧ props: follow pet-specific look mechanics plan.
- No whole-sprite rotation, replacement/googly eyes, visual clipping, seam bands, or transparent interior holes.
- Repaired direction: approval by independent visual QA worker ∨ explicit user inspection; repairing parent alone is insufficient.

## Repair Policy

Repair scope=smallest packaging-eligible: 1 standard row ∨ 1 complete coherent look row. NEVER mix individually generated repair cell into a new pet's final look row. AFTER EVERY relevant repair: re-run assembly ∧ deterministic validation ∧ direction QA ∧ continuity measurement ∧ semantic review.
