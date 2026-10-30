#!/usr/bin/env bash
set -Eeuo pipefail
source_dist="$(cd -- "$(dirname -- "$0")/.." && pwd -P)"; script="$source_dist/install-pi.sh"
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT
fail(){ echo "FAIL: $*" >&2; exit 1; }
ln(){ MSYS=winsymlinks:nativestrict command ln "$@"; }
agent="$tmp/agent dir"
bash "$script" --agent-dir "$agent" --dry-run >/dev/null; [[ ! -e "$agent" ]] || fail 'dry-run wrote files'
bash "$script" --agent-dir "$agent" >/dev/null
[[ -L "$agent/settings.json" ]] || fail 'fresh install missing symlink'
bash "$script" --agent-dir "$agent" >/dev/null
# All conflicts are preflighted before writes; regular file, parent file, and wrong symlink.
mkdir -p "$tmp/conflict"; printf user > "$tmp/conflict/settings.json"
if bash "$script" --agent-dir "$tmp/conflict" >/dev/null 2>&1; then fail 'regular conflict unexpectedly succeeded'; fi
[[ "$(<"$tmp/conflict/settings.json")" == user && ! -e "$tmp/conflict/agents/pi" ]] || fail 'conflict changed destination partially'
mkdir -p "$tmp/parent-conflict"; printf user > "$tmp/parent-conflict/agents"
if bash "$script" --agent-dir "$tmp/parent-conflict" >/dev/null 2>&1; then fail 'parent conflict unexpectedly succeeded'; fi
[[ -f "$tmp/parent-conflict/agents" && ! -e "$tmp/parent-conflict/settings.json" ]] || fail 'parent conflict mutated destination'
mkdir -p "$tmp/wrong"; ln -s "$tmp/not-the-source" "$tmp/wrong/settings.json"
if bash "$script" --agent-dir "$tmp/wrong" >/dev/null 2>&1; then fail 'wrong symlink conflict unexpectedly succeeded'; fi
bash "$script" --agent-dir "$tmp/wrong" --backup >/dev/null
[[ -L "$tmp/wrong/settings.json" ]] || fail 'wrong-link backup install failed'
# Back up regular file and dangling link in separate independent installs.
mkdir -p "$tmp/regular"; printf preserve > "$tmp/regular/settings.json"
bash "$script" --agent-dir "$tmp/regular" --backup >/dev/null
compgen -G "$tmp/regular/settings.json.backup.*" >/dev/null || fail 'regular backup missing'
[[ "$(<"$(compgen -G "$tmp/regular/settings.json.backup.*" | head -1)")" == preserve ]] || fail 'regular original not preserved'
mkdir -p "$tmp/dangling"; ln -s "$tmp/missing-target" "$tmp/dangling/settings.json"
if bash "$script" --agent-dir "$tmp/dangling" >/dev/null 2>&1; then fail 'dangling conflict unexpectedly succeeded'; fi
bash "$script" --agent-dir "$tmp/dangling" --backup >/dev/null
compgen -G "$tmp/dangling/settings.json.backup.*" >/dev/null || fail 'dangling link backup missing'
# Destination/source equality, ancestor/descendant overlap, and symlinked ancestors refuse.
for target in "$source_dist" "$(dirname "$source_dist")" "$source_dist/child"; do if bash "$script" --agent-dir "$target" >/dev/null 2>&1; then fail "overlap accepted: $target"; fi; done
ln -s "$source_dist" "$tmp/source-link"
if bash "$script" --agent-dir "$tmp/source-link/child" >/dev/null 2>&1; then fail 'symlinked destination parent accepted'; fi
# Missing-source fixture includes all shared entries except one specifically removed.
fixture="$tmp/fake/.HARNESS/.pi"; mkdir -p "$fixture/agent/extensions/subagent" "$fixture/agent/agents/pi" "$fixture/../agents" "$fixture/../crabcode/prompts" "$fixture/../rules" "$fixture/../skills" "$fixture/../persona" "$fixture/../etl" "$fixture/../scripts" "$fixture/../evaluation" "$fixture/../memo"
cp "$source_dist/install-pi.sh" "$fixture/"; cp "$source_dist/agent/settings.json" "$fixture/agent/"
cp "$source_dist/agent/extensions/subagent/config.json" "$fixture/agent/extensions/subagent/"
cp "$source_dist/agent/agents/pi/"*.md "$fixture/agent/agents/pi/"
for f in AGENTS.md POLYGLYPH.md HARDWARE.md; do printf x > "$fixture/../$f"; done
for f in worker.md WORKER.css SUBAGENTS.lrf AGENT_ROUTER.md; do printf x > "$fixture/../agents/$f"; done
for f in rules skills persona etl scripts evaluation memo; do [[ -e "$fixture/../$f" ]] || mkdir -p "$fixture/../$f"; done
printf x > "$fixture/../rules/BOOTSTRAP.lrf"
rm "$fixture/../rules/BOOTSTRAP.lrf"; rmdir "$fixture/../rules"
if bash "$fixture/install-pi.sh" --agent-dir "$tmp/missing-agent" >/dev/null 2>&1; then fail 'missing source unexpectedly succeeded'; fi
[[ ! -e "$tmp/missing-agent" ]] || fail 'missing source mutated destination'
# Inject link failure after backing up the first original and verify rollback restores it.
mkdir -p "$tmp/failbin" "$tmp/linkfail"; printf original > "$tmp/linkfail/settings.json"
printf '#!/bin/sh\nexit 91\n' > "$tmp/failbin/ln"; chmod +x "$tmp/failbin/ln"
if PATH="$tmp/failbin:$PATH" bash "$script" --agent-dir "$tmp/linkfail" --backup >/dev/null 2>&1; then fail 'injected link failure unexpectedly succeeded'; fi
[[ -f "$tmp/linkfail/settings.json" && "$(<"$tmp/linkfail/settings.json")" == original ]] || fail 'link-failure rollback lost original'
compgen -G "$tmp/linkfail/settings.json.backup.*" >/dev/null && fail 'rollback left unnecessary backup'
echo 'PASS bash installer scenarios'
