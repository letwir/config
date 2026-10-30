#!/usr/bin/env bash
set -Eeuo pipefail
usage(){ printf 'Usage: bash install-pi.sh [--agent-dir DIR] [--dry-run] [--backup] [--help]\n'; }
agent_dir="${HOME:?HOME is not set}/.pi/agent"; dry=0; backup=0
while (($#)); do case "$1" in
 --agent-dir) (($#>=2)) || { usage >&2; exit 2; }; agent_dir=$2; shift 2;;
 --dry-run) dry=1; shift;; --backup) backup=1; shift;; --help|-h) usage; exit 0;; *) echo "Unknown option: $1" >&2; usage >&2; exit 2;; esac; done
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"; shared="$(cd -- "$script_dir/.." && pwd -P)"
agent_dir="$(realpath -m -s -- "$agent_dir")" || { echo 'Cannot normalize destination.' >&2; exit 1; }
paths=(settings.json extensions/subagent agents/pi agents/worker.md agents/WORKER.css agents/SUBAGENTS.lrf agents/AGENT_ROUTER.md AGENTS.md POLYGLYPH.md HARDWARE.md rules skills prompts persona etl scripts evaluation memo)
sources=("$script_dir/agent/settings.json" "$script_dir/agent/extensions/subagent" "$script_dir/agent/agents/pi" "$shared/agents/worker.md" "$shared/agents/WORKER.css" "$shared/agents/SUBAGENTS.lrf" "$shared/agents/AGENT_ROUTER.md" "$shared/AGENTS.md" "$shared/POLYGLYPH.md" "$shared/HARDWARE.md" "$shared/rules" "$shared/skills" "$shared/crabcode/prompts" "$shared/persona" "$shared/etl" "$shared/scripts" "$shared/evaluation" "$shared/memo")
contains(){ [[ "$1" == "$2" || "$1" == "$2/"* ]]; }
for srcroot in "$script_dir" "$shared"; do if contains "$agent_dir" "$srcroot" || contains "$srcroot" "$agent_dir"; then echo "Unsafe source/destination overlap: $agent_dir <> $srcroot" >&2; exit 1; fi; done
# Reject symlinked ancestors and non-directory ancestors; they can redirect writes outside the requested destination.
p="$agent_dir"; while [[ "$p" != / && "$p" != . ]]; do if [[ -L "$p" ]]; then echo "Unsafe symlinked destination ancestor: $p" >&2; exit 1; fi; [[ ! -e "$p" || -d "$p" ]] || { echo "Destination ancestor is not a directory: $p" >&2; exit 1; }; p="$(dirname -- "$p")"; done
conflicts=(); for i in "${!paths[@]}"; do src="${sources[$i]}"; dst="$agent_dir/${paths[$i]}"; [[ -e "$src" ]] || { echo "Missing source: $src" >&2; exit 1; }
 # Existing parent components must be directories, not files or links.
 p="$(dirname -- "$dst")"; while [[ "$p" != "$agent_dir" && "$p" != / ]]; do if [[ -L "$p" || ( -e "$p" && ! -d "$p" ) ]]; then echo "Unsafe destination parent: $p" >&2; exit 1; fi; p="$(dirname -- "$p")"; done
 if [[ -L "$dst" ]] && [[ "$(readlink -- "$dst" || true)" == "$src" ]]; then continue; fi
 if [[ -e "$dst" || -L "$dst" ]]; then conflicts+=("$dst"); fi
done
if ((${#conflicts[@]})) && (( ! backup )); then printf 'Conflict (no changes made): %s\n' "${conflicts[@]}" >&2; exit 1; fi
if ((dry)); then for i in "${!paths[@]}"; do printf '%s -> %s\n' "${sources[$i]}" "$agent_dir/${paths[$i]}"; done; exit 0; fi
backed=(); created=(); fail(){ echo "ERROR: $*" >&2; exit 1; }
rollback(){ local rc=$? i d b dst rollback_failed=0; ((rc==0)) && return; trap - EXIT
 for ((i=${#created[@]}-1;i>=0;i--)); do d="${created[$i]}"; if [[ -e "$d" || -L "$d" ]]; then if ! rm -f -- "$d"; then echo "ROLLBACK FAILED removing link: $d" >&2; rollback_failed=1; fi; fi; done
 for ((i=${#backed[@]}-1;i>=0;i--)); do b="${backed[$i]%%|*}"; dst="${backed[$i]#*|}"; if [[ -e "$dst" || -L "$dst" ]]; then echo "ROLLBACK PRESERVED backup (destination occupied): $b" >&2; rollback_failed=1; elif ! mv -- "$b" "$dst"; then echo "ROLLBACK FAILED restoring $b to $dst" >&2; rollback_failed=1; fi; done
 if ((rollback_failed)); then echo "Install failed (exit $rc); rollback had errors—inspect listed paths." >&2; else echo "Install failed (exit $rc); rollback completed." >&2; fi; exit "$rc"; }
trap rollback EXIT
for i in "${!paths[@]}"; do src="${sources[$i]}"; dst="$agent_dir/${paths[$i]}"; [[ -L "$dst" && "$(readlink -- "$dst" || true)" == "$src" ]] && { echo "Already linked: $dst"; continue; }
 mkdir -p -- "$(dirname -- "$dst")" || fail "cannot create parent for $dst"
 if [[ -e "$dst" || -L "$dst" ]]; then ((backup)) || fail "conflict: $dst"; b="$dst.backup.$(date +%Y%m%d%H%M%S).$$"; n=0; while [[ -e "$b" || -L "$b" ]]; do ((n+=1)); b="$dst.backup.$(date +%Y%m%d%H%M%S).$$.$n"; done; mv -- "$dst" "$b" || fail "cannot back up $dst"; backed+=("$b|$dst"); fi
 if ! MSYS=winsymlinks:nativestrict ln -s -- "$src" "$dst" || [[ ! -L "$dst" ]]; then fail "symbolic link creation/verification failed for $dst; use PowerShell if Git Bash lacks native symlink support"; fi
 created+=("$dst"); echo "Linked: $dst -> $src"
done
trap - EXIT
