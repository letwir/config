@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# V2 Animation Rows

∀ newly hatched pet: atlas=8-column x 11-row; cells=192x208; final=1536x2288 ∧ `spriteVersionNumber: 2`.

| Row | State             | Used columns | Durations                                              |
| --- | ----------------- | -----------: | ------------------------------------------------------ |
| 0   | idle              |          0-5 | 280, 110, 110, 140, 140, 320 ms                        |
| 1   | running-right     |          0-7 | 120 ms each, final 220 ms                              |
| 2   | running-left      |          0-7 | 120 ms each, final 220 ms                              |
| 3   | waving            |          0-3 | 140 ms each, final 280 ms                              |
| 4   | jumping           |          0-4 | 140 ms each, final 280 ms                              |
| 5   | failed            |          0-7 | 140 ms each, final 240 ms                              |
| 6   | waiting           |          0-5 | 150 ms each, final 260 ms                              |
| 7   | running           |          0-5 | 120 ms each, final 220 ms                              |
| 8   | review            |          0-5 | 150 ms each, final 280 ms                              |
| 9   | look directions A |          0-7 | 000, 022.5, 045, 067.5, 090, 112.5, 135, 157.5 degrees |
| 10  | look directions B |          0-7 | 180, 202.5, 225, 247.5, 270, 292.5, 315, 337.5 degrees |

∀ standard animation row: cells after final used column MUST be fully transparent. Look-row cells: all used.

`000` degrees=looking up / 12 o'clock. Neutral/front=pointer deadzone ⇒ normal idle animation fallback.

## Row Purposes

- `idle`: calm low-distraction breathing/blinking loop ∧ reduced-motion first frame.
- `running-right`: rightward locomotion ∧ readable alternating cadence.
- `running-left`: leftward locomotion; mirror ONLY IF identity ∧ prop handedness remain correct; preserve frame order.
- `waving`: greeting ∨ attention gesture; clear start → raised gesture → return.
- `jumping`: anticipation → lift → peak → descent → settle.
- `failed`: readable error ∨ sad ∨ deflated reaction; no noisy detached effects.
- `waiting`: expectant asking pose for approval ∨ help ∨ user input.
- `running`: active task work ∨ processing; not literal foot-running.
- `review`: focused inspection of completed output.
- rows `9-10`: 1 continuous clockwise 16-pose look loop; pet-specific eye ∧ head ∧ body ∧ appendage ∧ prop mechanics.
