#!/usr/bin/env bash
set -Eeuo pipefail
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
SCRIPT=$SCRIPT_DIR/sync-harness.sh
[[ -x $SCRIPT ]] || { echo 'sync-harness.sh must be executable' >&2; exit 1; }
root=$(mktemp -d "${TMPDIR:-/tmp}/sync harness tests.XXXXXX")
cleanup() { rm -rf -- "$root"; }
trap cleanup EXIT
ok=0
run_expect() { local want=$1; shift; local actual=0; "$@" >/dev/null 2>&1 || actual=$?; if [[ $actual != "$want" ]]; then echo "FAIL expected exit $want got $actual: $*" >&2; exit 1; fi; ok=$((ok+1)); }
fixture() { local name=$1; mkdir -p "$root/$name/.harness/skills" "$root/$name/.gemini" "$root/$name/.codex" "$root/$name/.agents" "$root/$name/.claude" "$root/$name/.opencode"; }
# Consolidate source files through roots containing spaces and verify manifest fields.
fixture 'space home'
mkdir -p "$root/space home/.agents/skills/example"
printf 'fixture\n' > "$root/space home/.agents/skills/example/SKILL.md"
run_expect 0 "$SCRIPT" --consolidate --home "$root/space home" --pass-thru
mf=$(printf '%s\n' "$root/space home/.harness/manifests/"*.json)
python3 - "$mf" <<'PY'
import json,sys
m=json.load(open(sys.argv[1])); required=('version','run_id','status','conflict_policy','backup_path','sources','conflicts','before','after','operations')
assert all(k in m for k in required),m
assert m['version']==1 and m['status']=='Applied' and m['run_id']
PY
[[ $(cat "$root/space home/.harness/skills/example/SKILL.md") == fixture ]] || exit 1
ok=$((ok+1))
# Restore a completed run to its pre-run snapshot.
fixture 'restore home'
mkdir -p "$root/restore home/.harness/skills" "$root/restore home/.agents/skills"
printf old > "$root/restore home/.harness/skills/old"
printf new > "$root/restore home/.agents/skills/new"
run_expect 0 "$SCRIPT" --consolidate --home "$root/restore home"
restore_manifest=$(printf '%s\n' "$root/restore home/.harness/manifests/"*.json)
run_expect 0 "$SCRIPT" --restore-manifest "$restore_manifest" --home "$root/restore home"
[[ -f "$root/restore home/.harness/skills/old" && ! -e "$root/restore home/.harness/skills/new" ]] || { echo 'restore did not recover pre-run snapshot' >&2; exit 1; }
ok=$((ok+1))
# Preview is strictly non-mutating, including creation of SSOT.
fixture 'dry run'
mkdir -p "$root/dry run/.agents/skills/x"; echo data > "$root/dry run/.agents/skills/x/a"
run_expect 0 "$SCRIPT" --consolidate --dry-run --home "$root/dry run"
[[ ! -e "$root/dry run/.harness/skills/x" && ! -e "$root/dry run/.harness/manifests" ]] || { echo 'dry-run changed files' >&2; exit 1; }
ok=$((ok+1))
# Link: missing target, wrong symlink target, ordinary directory, plus paths with spaces.
fixture 'links home'
printf source > "$root/links home/.harness/skills/file"
mkdir -p "$root/links home/.agents" "$root/links home/.codex" "$root/links home/.gemini/config/plugins/ltw-skills-plugin" "$root/links home/.claude" "$root/links home/.opencode"
mkdir -p "$root/links home/wrong"; ln -s "$root/links home/wrong" "$root/links home/.agents/skills"
mkdir -p "$root/links home/.codex/skills/old"
run_expect 0 "$SCRIPT" --link-skills --home "$root/links home" --mode symlink
[[ -L "$root/links home/.agents/skills" && $(readlink -f "$root/links home/.agents/skills") == $(readlink -f "$root/links home/.harness/skills") ]] || exit 1
backup=$(compgen -G "$root/links home/.codex/skills.bak-*" || true)
[[ -n $backup && -d $backup ]] || { echo 'ordinary directory backup missing' >&2; exit 1; }
run_expect 0 "$SCRIPT" --check --home "$root/links home"
ok=$((ok+1))
# Check failure exit for a missing SSOT and for a wrong or ordinary target.
fixture 'missing ssot'; run_expect 1 "$SCRIPT" --check --home "$root/missing ssot"
fixture 'bad links'; mkdir -p "$root/bad links/.harness/skills" "$root/bad links/.agents" "$root/bad links/.codex" "$root/bad links/.gemini" "$root/bad links/.claude" "$root/bad links/.opencode"
mkdir -p "$root/bad links/wrong" "$root/bad links/.agents/skills"; ln -s "$root/bad links/wrong" "$root/bad links/.codex/skills"
run_expect 1 "$SCRIPT" --check --home "$root/bad links"
ok=$((ok+1))
# Fail-policy conflicts are an error and preserve destination data.
fixture 'conflict'; echo destination > "$root/conflict/.harness/skills/a"; mkdir -p "$root/conflict/.agents/skills"; echo source > "$root/conflict/.agents/skills/a"
run_expect 1 "$SCRIPT" --consolidate --home "$root/conflict"
[[ $(cat "$root/conflict/.harness/skills/a") == destination ]] || exit 1
printf 'PASS (%d assertions)\n' "$ok"
