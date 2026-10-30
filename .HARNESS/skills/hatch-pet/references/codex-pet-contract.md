@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# Codex V2 Pet Contract

## Sprite Atlas

- Version: `spriteVersionNumber: 2`.
- Format: PNG or WebP.
- Dimensions: `1536x2288`.
- Grid: 8 columns x 11 rows.
- Cell: `192x208`.
- Background: transparent.
- Rows `0-8`: standard animation states.
- Rows `9-10`: 16 clockwise look directions.
- Unused standard-row cells: fully transparent.

8x9 `1536x1872` atlas=intermediate assembly artifact ONLY; NEVER package as newly hatched pet.

## Look Directions

- Row `9`: `000`, `022.5`, `045`, `067.5`, `090`, `112.5`, `135`, `157.5` degrees.
- Row `10`: `180`, `202.5`, `225`, `247.5`, `270`, `292.5`, `315`, `337.5` degrees.
- `000`=up / 12 o'clock; not neutral/front.
- Neutral/front=no-vector deadzone ⇒ idle fallback.

## Local Custom Pet Package

Package files → directory:

```text
${CODEX_HOME:-$HOME/.codex}/pets/<pet-name>/
├── pet.json
└── spritesheet.webp
```

Required manifest schema:

```json
{
  "id": "pet-name",
  "displayName": "Pet Name",
  "description": "One short sentence.",
  "spriteVersionNumber": 2,
  "spritesheetPath": "spritesheet.webp"
}
```

`spriteVersionNumber: 2` ⇒ app derives 11-row layout ∧ look-direction behavior. Omission ⇒ pet defaults to v1 ∧ 2288-pixel-tall spritesheet rejected.
