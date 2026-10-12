#!/usr/bin/env bash
# Bash port of check-rule-distribution.ps1; Python is used only to parse/encode JSON.
set -Eeuo pipefail
fail() { printf 'ERROR: %s\n' "$*" >&2; exit 2; }
MANIFEST="${1:-$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)/rules/client-distribution.json}"
USER_HOME=${HOME:?HOME is not set}
JSON=0
if [[ ${1:-} == --json ]]; then MANIFEST="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)/rules/client-distribution.json"; shift; fi
while (($#)); do case $1 in --json) JSON=1;; --home) (($# >= 2)) || fail 'missing value for --home'; USER_HOME=$2; shift;; --manifest) (($# >= 2)) || fail 'missing value for --manifest'; MANIFEST=$2; shift;; *) [[ -n $MANIFEST && $1 != --* ]] && MANIFEST=$1 || fail "unknown argument: $1";; esac; shift; done
[[ -f $MANIFEST ]] || fail "Manifest missing: $MANIFEST"
MANIFEST=$(realpath -m -- "$MANIFEST")
ROOT=$(dirname -- "$(dirname -- "$MANIFEST")")
export MANIFEST ROOT USER_HOME JSON
# Python emits TSV records only; all path and content checks are performed in Bash.
mapfile -t ROWS < <(python3 - "$MANIFEST" <<'PY'
import json,sys
m=json.load(open(sys.argv[1],encoding='utf-8'))
if m.get('schemaVersion') != 1 or not isinstance(m.get('entries'),list): sys.exit(2)
print('CANON\t'+m['canonical']+'\t'+m['canonicalSha256'])
for e in m['entries']:
 print('ENTRY\t'+e['name']+'\t'+str(e['required']).lower()+'\t'+e['entry'])
 for x in e['chain']:
  print('ITEM\t'+e['name']+'\t'+x['path']+'\t'+x.get('kind','')+'\t'+x['sha256']+'\t'+','.join(x.get('ids',[]))+'\t'+','.join(x.get('allowedDifferenceIds',[]))+'\t'+x.get('contains',''))
PY
) || fail 'Invalid manifest'
resolve() { local p=$1; [[ $p != *'..'* ]] || fail 'Invalid manifest path'; case $p in '{HarnessRoot}/'*) printf '%s/%s' "$ROOT" "${p#\{HarnessRoot\}/}";; '{UserHome}/'*) printf '%s/%s' "$USER_HOME" "${p#\{UserHome\}/}";; *) fail "Unknown manifest path token: $p";; esac; }
contained() { local root cand; root=$(realpath -m -- "$1"); cand=$(realpath -m -- "$2"); [[ $cand == "$root" || $cand == "$root/"* ]]; }
ids() { awk -F'|' '/^@lrf=1\|/{next} /^[RL]\|/{if(NF>=6) print $2}' "$1" | sort -u; }
canonical= canonhash=; results=(); names=() requireds=(); allowed_only=0
for row in "${ROWS[@]}"; do IFS=$'\t' read -r type a b c d e f g h <<<"$row"; case $type in CANON) canonical=$(resolve "$a"); canonhash=$b;; ENTRY) names+=("$a"); requireds+=("$b");; esac; done
[[ -f $canonical ]] || fail 'Canonical LRF is missing'
[[ $(sha256sum "$canonical" | cut -d' ' -f1) == "${canonhash,,}" ]] || fail 'Canonical LRF hash differs from manifest'
for ((i=0;i<${#names[@]};i++)); do
  name=${names[i]}; required=${requireds[i]}; missing=0; mismatch=0; unexpected=0; allowed=0; details=; entrypath=; chainpaths=(); expectedids=; allowedids=; iddetails=; allowed_only=0
 for row in "${ROWS[@]}"; do IFS=$'\t' read -r type a b c d e f g h <<<"$row"; [[ $type == ITEM && $a == "$name" ]] || continue
 p=$(resolve "$b"); chainpaths+=("$b"); [[ -n $entrypath || $b != "${ROWS[*]}" ]] || :
 [[ $p != *$'\n'* ]] || fail 'Invalid path'
 # Reject symlinks/reparse paths, including links in intermediate components.
 q=$p; while [[ $q != / && $q != . ]]; do [[ ! -L $q ]] || fail "Reparse-point path: $b"; q=$(dirname -- "$q"); done
 root=$ROOT; [[ $b == '{UserHome}/'* ]] && root=$USER_HOME
 contained "$root" "$p" || fail "Manifest path escapes roots: $b"
 [[ -f $p ]] || { missing=1; continue; }
  hash=$(sha256sum "$p" | cut -d' ' -f1); [[ ${hash,,} == ${d,,} ]] || { unexpected=1; details+="hash:$b;"; }
 if [[ -n $h ]]; then text=${h//\{UserHomeForward\}/$USER_HOME}; text=${text//\{HarnessRootForward\}/$ROOT}; text=${text//\{UserHome\}/$USER_HOME}; text=${text//\{HarnessRoot\}/$ROOT}; grep -Fq -- "$text" "$p" || mismatch=1; fi
  if [[ $c == lrf ]]; then
   actual=$(ids "$p"); expected=$(printf '%s\n' "${e//,/ }" | tr ' ' '\n' | sed '/^$/d' | sort -u); allowedids=$(printf '%s\n' "${f//,/ }" | tr ' ' '\n' | sed '/^$/d' | sort -u)
   mapfile -t actualset <<<"$actual"; mapfile -t expectedset <<<"$expected"; mapfile -t allowset <<<"$allowedids"
   item_missing=(); item_extra=(); item_diff=();
   for id in "${expectedset[@]}"; do [[ -n $id ]] || continue; if ! printf '%s\n' "$actual" | grep -Fqx -- "$id"; then item_missing+=("$id"); item_diff+=("$id"); fi; done
   for id in "${actualset[@]}"; do [[ -n $id ]] || continue; if ! printf '%s\n' "$expected" | grep -Fqx -- "$id"; then item_extra+=("$id"); item_diff+=("$id"); fi; done
   if ((${#item_diff[@]})); then
    item_allowed=1
    for id in "${item_diff[@]}"; do if ! printf '%s\n' "$allowedids" | grep -Fqx -- "$id"; then item_allowed=0; fi; done
    iddetails+="MissingIds:${item_missing[*]};ExtraIds:${item_extra[*]};"
    if ((item_allowed)); then allowed_only=1; else unexpected=1; fi
   fi
  fi
 done
  if ((missing)); then class=Missing; detail=chain-path-missing; elif ((mismatch)); then class=EntryMismatch; detail=reference-mismatch; elif ((unexpected)); then class=UnexpectedDifference; detail=${details}${iddetails:-id-difference}; elif ((allowed_only)); then class=AllowedDifference; detail=${iddetails}; else class=Exact; detail=all-chain-hashes-and-ids-match; fi
 results+=("$class|$name|$required|$detail")
done
python3 - "$MANIFEST" "$JSON" "${results[@]}" <<'PY'
import datetime,json,sys
manifest,asjson,*rows=sys.argv[1:]
r=[]; failed=False
for row in rows:
 c,n,req,d=row.split('|',3); required=req=='true'; r.append(dict(Name=n,Classification=c,Required=required,Detail=d)); failed |= required and c in ('Missing','EntryMismatch','UnexpectedDifference')
out=dict(SchemaVersion=1,Manifest=manifest,Passed=not failed,Results=r,GeneratedAt=datetime.datetime.now().astimezone().isoformat())
if asjson=='1': print(json.dumps(out,separators=(',',':')))
else:
 for x in r: print('[{}] {}: {}'.format(x['Classification'],x['Name'],x['Detail']))
 print('Rule distribution doctor '+('passed.' if not failed else 'failed.'))
sys.exit(1 if failed else 0)
PY
