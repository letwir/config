---
name: skill-installer
description: Install Codex skills into $CODEX_HOME/skills from a curated list or a GitHub repo path. Use when a user asks to list installable skills, install a curated skill, or install a skill from another repo (including private repos).
metadata:
  short-description: Install curated skills from openai/skills or other repos
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/XML::ρ→max|target:skill-installer|legibility:LLM≫human]
# Skill Installer

🎯 skill-management: default source=https://github.com/openai/skills/tree/main/skills/.curated; user may specify another location. Experimental source=https://github.com/openai/skills/tree/main/skills/.experimental; same installation procedure.

⊢ choose helper action by request:
- 🔎 list when user asks availability ∨ invokes skill without specifying action. Default `.curated`; experimental request⇒pass `--path skills/.experimental`.
- 📦 curated install when user supplies skill name.
- 📦 other-repository install when user supplies GitHub repo/path, private repositories included.

⊢ install only through helper scripts; every script uses network. In sandbox, request escalation before running them.

## Communication

🔎 For listing, adapt the following output to request context. Experimental request⇒list `.experimental` and label source:
"""
Skills from {repo}:
1. skill-1
2. skill-2 (already installed)
3. ...
Which ones would you like installed?
"""

✅ After successful installation, tell user skill is available on their next turn.

## Scripts

- `scripts/list-skills.py` (prints skills list with installed annotations)
- `scripts/list-skills.py --format json`
- Example (experimental list): `scripts/list-skills.py --path skills/.experimental`
- `scripts/install-skill-from-github.py --repo <owner>/<repo> --path <path/to/skill> [<path/to/skill> ...]`
- `scripts/install-skill-from-github.py --url https://github.com/<owner>/<repo>/tree/<ref>/<path>`
- Example (experimental skill): `scripts/install-skill-from-github.py --repo openai/skills --path skills/.experimental/<skill-name>`

## Behavior and Options

- Public GitHub repo default⇒direct download.
- ⛔ Download auth/permission failure⇒git sparse-checkout fallback.
- ⛔ Destination skill directory already exists⇒abort.
- 📍 destination `$CODEX_HOME/skills/<skill-name>`; default `~/.codex/skills`.
- Multiple `--path` values⇒install each skill; name=path basename unless `--name` supplied.
- Options: `--ref <ref>` (default `main`), `--dest <path>`, `--method auto|download|git`.

## Notes

- Curated listing uses GitHub API at `https://github.com/openai/skills/tree/main/skills/.curated`. If unavailable, explain error and exit.
- Private GitHub repo access may use existing git credentials or optional `GITHUB_TOKEN`/`GH_TOKEN` for download. Do not request or expose credentials.
- Git fallback order: HTTPS, then SSH.
- Skills at https://github.com/openai/skills/tree/main/skills/.system are preinstalled; normally explain this rather than installing. If user insists, download and overwrite.
- Installed annotations source: `$CODEX_HOME/skills`.
