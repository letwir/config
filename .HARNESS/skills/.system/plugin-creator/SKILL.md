---
name: plugin-creator
description: Create and scaffold plugin directories for Codex with a required `.codex-plugin/plugin.json`, optional plugin folders/files, valid manifest defaults, and personal-marketplace entries by default. Use when Codex needs to create a new personal plugin, add optional plugin structure, generate or update marketplace entries for plugin ordering and availability metadata, or update an existing local plugin during development with the CLI-driven cachebuster and reinstall flow.
---

@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# Plugin Creator

## Quick Start

🔎 New plugin → scaffold from skill root; default personal path `~/plugins/<plugin-name>`.
1. Name → lower-case hyphen-case; ≤64 chars; generated directory=name in `plugin.json`.

```bash
# Plugin names are normalized to lower-case hyphen-case and must be <= 64 chars.
# The generated folder and plugin.json name are always the same.
# Run from the skill root (the directory containing this `SKILL.md`).
# By default creates in `~/plugins/<plugin-name>`.
python3 scripts/create_basic_plugin.py <plugin-name>
```
2. Manifest: edit `<plugin-path>/.codex-plugin/plugin.json` only for requested metadata; scaffold defaults valid; no `[TODO: ...]`.

```bash
# Personal marketplace entries default to `~/.agents/plugins/marketplace.json`.
python3 scripts/create_basic_plugin.py my-plugin --with-marketplace
```
3. UI ordering requested → personal marketplace default `~/.agents/plugins/marketplace.json`; custom `--marketplace-name` only if default `personal` is already taken/installed.

```bash
python3 scripts/create_basic_plugin.py my-plugin \
  --with-marketplace \
  --marketplace-name team-local
```

```bash
python3 scripts/create_basic_plugin.py my-plugin \
  --path <repo-root>/plugins \
  --marketplace-path <repo-root>/.agents/plugins/marketplace.json \
  --with-marketplace
```
4. Destination(repo/team) ⇒ explicit user request; custom marketplace path ⇒ MUST ensure marketplace installed before advising reinstall; personal default ⇒ implicitly discovered; other paths ⇒ ¬implicitly discovered. Windows path ⇒ equivalent user-profile path.

```bash
python3 scripts/create_basic_plugin.py my-plugin \
  --path <parent-plugin-directory> \
  --marketplace-path <marketplace-json-path> \
  --with-skills --with-hooks --with-scripts --with-assets --with-mcp --with-apps --with-marketplace
```
5. Optional companion folders only as needed; `<parent-plugin-directory>` contains new `<plugin-name>` folder (e.g. `~/plugins`).

```bash
python3 scripts/validate_plugin.py <plugin-path>
```
6. Validate generated plugin before handoff: `python3 scripts/validate_plugin.py <plugin-path>`.

```bash
python3 scripts/read_marketplace_name.py
python3 scripts/update_plugin_cachebuster.py <plugin-path>
```

UpdateFlow(existing local plugin ∧ development): MUST keep scaffold flow ∧ use cachebuster/reinstall references ∧ MUST_NOT hand-edit marketplace files; repo/team marketplace ⇒ pass `--marketplace-path <marketplace-json-path>` to marketplace-name reader; helper default cachebuster ⇒ preferred; override iff explicitly requested. Reference: `references/installing-and-updating.md` = iteration/reinstall flow.

## What this skill creates

- 🧱 Root: `/<parent-plugin-directory>/<plugin-name>/`; default storage `~/plugins/<plugin-name>/`.
- Required manifest: `/<parent-plugin-directory>/<plugin-name>/.codex-plugin/plugin.json`.
- Manifest shape: validated schema accepted by ingestion path.
- Marketplace-backed default: personal file `~/.agents/plugins/marketplace.json`; `--with-marketplace` creates/updates it; absent file → seed personal marketplace root before first entry.
- Name normalization (skill-creator rules): `My Plugin` → `my-plugin`; `My--Plugin` → `my-plugin`; underscores, spaces, punctuation → `-`; lowercase hyphen-delimited result, consecutive hyphens collapsed.
- Optional companion structure: `skills/`, `hooks/`, `scripts/`, `assets/`, `.mcp.json`, `.app.json`.

## Marketplace workflow

- MarketplacePolicy(default): personal iff file=`~/.agents/plugins/marketplace.json`; repo/team iff explicit user request ∧ both `--path` and `--marketplace-path` supplied (opt-in).
- MarketplacePolicy(`--marketplace-name`): MAY seed different new marketplace file iff default `personal` already taken; MUST_NOT rename in place; existing file ⇒ top-level `name` MUST already match.
- MarketplacePolicy(custom explicit path): MUST install via `codex plugin marketplace add` before advising reinstall.
- Plugin name regex: `[A-Za-z0-9_-]+(\.[A-Za-z0-9_-]+)*`; marketplace name regex: `[A-Za-z0-9_-]+`.
- MarketplacePolicy(existing marketplace): MUST validate names via `scripts/read_marketplace_name.py` (no argument = default personal; explicit path supports repo/team); validation failure ⇒ STOP. New marketplace ⇒ scaffold validates.
- Before updating existing plugins, validate both identifiers before any file mutation or install-command construction.
- Generated source path in either location: `./plugins/<plugin-name>`.
- Marketplace root metadata: top-level `name`, optional `interface.displayName`. `plugins[]` order is Codex render order; append unless user explicitly requests reorder. `displayName` belongs in `interface`, never an entry.
- Every generated entry requires `policy.installation`, `policy.authentication`, and `category`. Defaults: `policy.installation: "AVAILABLE"`; `policy.authentication: "ON_INSTALL"`. Change defaults only to another allowed value.
- Allowed `policy.installation`: `NOT_AVAILABLE`, `AVAILABLE`, `INSTALLED_BY_DEFAULT`. Allowed `policy.authentication`: `ON_INSTALL`, `ON_USE`.
- `policy.products` is an override: omit unless user explicitly requests product gating.
- Generated plugin entry shape:

```json
{
  "name": "plugin-name",
  "source": {
    "source": "local",
    "path": "./plugins/plugin-name"
  },
  "policy": {
    "installation": "AVAILABLE",
    "authentication": "ON_INSTALL"
  },
  "category": "Productivity"
}
```

- `--force` only for intentional replacement of an existing marketplace entry with the same plugin name.
- Missing target marketplace file ⇒ seed root with top-level `"name"`, an `"interface"` object containing `"displayName"`, and a `plugins` array, then add the entry.
- New marketplace root shape:

```json
{
  "name": "personal",
  "interface": {
    "displayName": "Personal"
  },
  "plugins": [
    {
      "name": "plugin-name",
      "source": {
        "source": "local",
        "path": "./plugins/plugin-name"
      },
      "policy": {
        "installation": "AVAILABLE",
        "authentication": "ON_INSTALL"
      },
      "category": "Productivity"
    }
  ]
}
```

## Required behavior

- directory.name ⇔ normalized plugin name ⇔ `plugin.json` `"name"` (always equal).
- `.codex-plugin/plugin.json` MUST exist; do not remove required structure.
- manifest `[TODO: ...]` ⇒ MUST_NOT.
- `apps`/`mcpServers` ∉ `plugin.json` unless respective companion file exists/is created.
- Unsupported validation-rejected manifest fields (incl. `hooks`) ⇒ omit.
- Existing plugin path + creating files ⇒ `--force` iff overwrite intentional.
- Existing marketplace `interface.displayName` ⇒ preserve.
- Generated entry MUST include `policy.installation` + `policy.authentication` + `category`, defaults included.
- `policy.products` ⇒ include iff user explicitly requests override.
- marketplace `source.path` relative to selected marketplace root = `./plugins/<plugin-name>`.
- `--marketplace-name` MAY only create new file with non-`personal` name iff `personal` already taken or installed elsewhere.
- Marketplace write needs Codex approval ⇒ obtain approval before proceeding. User-run write preferred ⇒ give exact scaffold command, then continue with validation or subsequent plugin edits (not vague).
- Existing local-plugin development update ⇒ MUST_NOT hand-edit marketplace config/`marketplace.json`; use `references/installing-and-updating.md` + `scripts/update_plugin_cachebuster.py` update flow.
- Default personal marketplace (`~/.agents/plugins/marketplace.json`) ⇒ MUST_NOT tell user `codex plugin marketplace add`; command is only for explicit non-default configuration.
- User supplied non-default `--marketplace-path` ⇒ marketplace installed before reinstall instructions; if not configured, use `codex plugin marketplace add <path-to-marketplace-root>`.
- Workflow created/updated marketplace-backed plugin ⇒ final user response MUST end with short Codex app handoff: exact intro `To view this in the Codex app:`; Markdown links (not raw URLs/code spans), labels `View <normalized plugin name>` and `Share <normalized plugin name>`.
- View URL template: `codex://plugins/<normalized plugin name>?marketplacePath=<absolute marketplace.json path>`.
- Share URL template: same URL + `&mode=share`.
- Substitute actual normalized name + absolute scaffolded `marketplace.json` path; URL-encode path segment + query value as needed.
- Deeplinks MUST_NOT include `pluginName` or `hostId` query keys; Codex derives both after click.
- No marketplace entry created/updated ⇒ MUST_NOT emit either View/Share link.

## References

Reference(canonical exact sample JSON for plugin manifests + marketplace entries):

- Reference: `references/plugin-json-spec.md`
- Reference: `references/installing-and-updating.md` = update/reinstall guidance while iterating on existing local plugin + new-thread pickup behavior after reinstall

## Validation

Validation(after editing `SKILL.md`): MUST run:

```bash
python3 ../skill-creator/scripts/quick_validate.py .
```

Validation(before handing back a generated plugin): MUST run:

```bash
python3 scripts/validate_plugin.py <plugin-path>
```
