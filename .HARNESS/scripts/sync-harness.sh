#!/usr/bin/env bash
# Bash/Linux port of sync-harness.ps1. Paths are always quoted; manifest JSON is
# generated with Python 3 (the system utility), while synchronization is Bash.
set -Eeuo pipefail

usage() { cat <<'EOF'
Usage: sync-harness.sh [--consolidate] [--link-skills] [--check]
  [--mode symlink|copy] [--conflict-policy Fail|PreferSource|PreferDestination]
  [--restore-manifest PATH] [--dry-run] [--pass-thru] [--home PATH]
EOF
}
fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
CONSOLIDATE=0 LINK=0 CHECK=0 DRY=0 PASS=0 MODE=symlink POLICY=Fail RESTORE= HOME_PATH=${HOME:?HOME is not set}
while (($#)); do
  case $1 in
    --consolidate) CONSOLIDATE=1;; --link-skills) LINK=1;; --check) CHECK=1;;
    --dry-run) DRY=1;; --pass-thru) PASS=1;;
    --mode|--conflict-policy|--restore-manifest|--home)
      (($# >= 2)) || fail "missing value for $1"
      case $1 in --mode) MODE=$2;; --conflict-policy) POLICY=$2;; --restore-manifest) RESTORE=$2;; --home) HOME_PATH=$2;; esac
      shift;;
    -h|--help) usage; exit 0;; *) fail "unknown argument: $1";;
  esac
  shift
done
[[ $MODE == symlink || $MODE == copy ]] || fail "mode must be symlink or copy"
[[ $POLICY == Fail || $POLICY == PreferSource || $POLICY == PreferDestination ]] || fail "invalid conflict policy"
[[ -n $HOME_PATH ]] || fail '--home cannot be empty'
HOME_PATH=$(cd -- "$HOME_PATH" 2>/dev/null && pwd -P) || fail "home does not exist: $HOME_PATH"
HARNESS=$HOME_PATH/.harness; SKILLS=$HARNESS/skills
TARGETS=("$HOME_PATH/.gemini/skills" "$HOME_PATH/.gemini/config/plugins/ltw-skills-plugin/skills" "$HOME_PATH/.codex/skills" "$HOME_PATH/.agents/skills" "$HOME_PATH/.claude/skills" "$HOME_PATH/.opencode/skills")

# Atomic JSON replacement: temporary file is created beside the destination.
write_manifest() {
  local path=$1; shift
  mkdir -p -- "$(dirname -- "$path")" || fail "cannot create manifest directory"
  local tmp; tmp=$(mktemp "${path}.tmp.XXXXXX") || fail 'cannot create manifest temporary file'
  if ! python3 - "$tmp" "$@" <<'PY'
import json, os, sys
out, payload = sys.argv[1], sys.argv[2]
with open(out, 'w', encoding='utf-8') as f:
    f.write(payload); f.write('\n'); f.flush(); os.fsync(f.fileno())
os.replace(out, os.path.join(os.path.dirname(out), os.path.basename(os.path.splitext(out)[0])))
PY
  then rm -f -- "$tmp"; fail 'manifest write failed'; fi
}
# Use Python only for safe JSON encoding/atomic rename, not for sync behavior.
manifest_file() { printf '%s/manifests/sync-%s.json' "$HARNESS" "$RUN_ID"; }
RUN_ID=$(date -u +%Y%m%dT%H%M%S)-$$-$RANDOM
json_escape() { python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'; }
snapshot() {
  local root=$1
  [[ -d $root ]] || { printf '[]'; return; }
  python3 - "$root" <<'PY'
import hashlib,json,os,stat,sys
root=os.path.abspath(sys.argv[1]); result=[]
for base, dirs, files in os.walk(root, followlinks=False):
    dirs.sort(); files.sort()
    for name in dirs+files:
        p=os.path.join(base,name); rel=os.path.relpath(p,root); st=os.lstat(p)
        if stat.S_ISLNK(st.st_mode): typ='symlink'; target=os.readlink(p); digest=None
        elif stat.S_ISDIR(st.st_mode): typ='directory'; target=None; digest=None
        elif stat.S_ISREG(st.st_mode):
            typ='file'; target=None; h=hashlib.sha256()
            with open(p,'rb') as f:
                for b in iter(lambda:f.read(131072),b''): h.update(b)
            digest=h.hexdigest()
        else: typ='other'; target=None; digest=None
        result.append(dict(relative_path=rel,type=typ,hash=digest,target=target))
print(json.dumps(sorted(result,key=lambda x:x['relative_path']),separators=(',',':')))
PY
}

if [[ -n $RESTORE ]]; then
  [[ -f $RESTORE ]] || fail "restore manifest not found: $RESTORE"
  [[ -d $SKILLS ]] || fail "SSOT directory does not exist: $SKILLS"
  backup=$(python3 - "$RESTORE" "$HARNESS" <<'PY'
import json, os, sys
manifest, harness = sys.argv[1:]
try:
    with open(manifest, encoding='utf-8') as f: m=json.load(f)
    backup=os.path.abspath(m['backup_path'])
    root=os.path.abspath(os.path.join(harness, 'backups')) + os.sep
    if m.get('version') != 1 or not backup.startswith(root) or not os.path.isdir(backup): raise ValueError()
    if not isinstance(m.get('after'), list): raise ValueError()
    print(backup)
except Exception:
    sys.exit(1)
PY
  ) || fail 'unsupported, incomplete, or unsafe restore manifest'
  current=$(snapshot "$SKILLS")
  expected=$(python3 - "$RESTORE" <<'PY'
import json,sys
with open(sys.argv[1],encoding='utf-8') as f: print(json.dumps(json.load(f)['after'],sort_keys=True,separators=(',',':')))
PY
  ) || fail 'cannot read restore snapshot'
  normalized=$(python3 -c 'import json,sys; print(json.dumps(json.loads(sys.stdin.read()),sort_keys=True,separators=(",",":")))' <<<"$current")
  [[ $normalized == "$expected" ]] || fail 'SSOT changed since manifest; refusing restore'
  if ((DRY)); then printf '[DRY-RUN] would restore %s from %s\n' "$SKILLS" "$backup"
  else
    staging=$(mktemp -d "$HARNESS/.restore.XXXXXX") || fail 'cannot create restore staging directory'
    cp -a -- "$backup/." "$staging/" || { rm -rf -- "$staging"; fail 'cannot stage restore backup'; }
    old="$HARNESS/.skills-old-$RUN_ID"
    mv -- "$SKILLS" "$old" && mv -- "$staging" "$SKILLS" || { [[ -d $old ]] && mv -- "$old" "$SKILLS"; rm -rf -- "$staging"; fail 'restore switch failed'; }
    rm -rf -- "$old"
    printf '[OK] restored SSOT from %s\n' "$RESTORE"
  fi
  exit 0
fi

status=0
if ((CONSOLIDATE)); then
  sources=("$HOME_PATH/.agents/skills" "$HOME_PATH/.codex/skills" "$HOME_PATH/.gemini/skills" "$HOME_PATH/.gemini/config/plugins/ltw-skills-plugin/skills")
  existing=(); for s in "${sources[@]}"; do [[ -d $s && ! -L $s ]] && existing+=("$s"); done
  before=$(snapshot "$SKILLS")
  # Stage desired files and conflict records using a temporary tree and Python planner.
  tmp=$(mktemp -d "${TMPDIR:-/tmp}/sync-harness.XXXXXX") || fail 'cannot create staging directory'
  trap 'rm -rf -- "$tmp"' EXIT
  if python3 - "$SKILLS" "$tmp/plan.json" "$POLICY" "${existing[@]}" <<'PY'
import hashlib,json,os,stat,sys
root,out,policy,*sources=sys.argv[1:]
def tree(path):
 d={}
 if not os.path.isdir(path): return d
 for base,ds,fs in os.walk(path,followlinks=False):
  ds.sort(); fs.sort()
  for n in ds+fs:
   p=os.path.join(base,n); rel=os.path.relpath(p,path); st=os.lstat(p)
   if stat.S_ISLNK(st.st_mode): v=('symlink',os.readlink(p))
   elif stat.S_ISDIR(st.st_mode): v=('directory',None)
   elif stat.S_ISREG(st.st_mode):
    h=hashlib.sha256(open(p,'rb').read()).hexdigest(); v=('file',h)
   else: v=('other',None)
   d[rel]=v
 return d
chosen=tree(root); owner={k:root for k in chosen}; operations=[]; conflicts=[]; data={}
for src in sources:
 for rel, val in tree(src).items():
  prev=chosen.get(rel); conflict=prev is not None and prev!=val
  decision=policy if conflict else 'Apply'
  if conflict: conflicts.append(dict(source=src,relative_path=rel,existing_type=prev[0],incoming_type=val[0],decision=decision))
  operations.append(dict(source=src,relative_path=rel,conflict=conflict,decision=decision,applied=not(conflict and policy=='PreferDestination')))
  if prev is None or not conflict or policy=='PreferSource':
   chosen[rel]=val; owner[rel]=src
   if val[0]=='file': data[rel]=os.path.join(src,rel)
   elif val[0]=='symlink': data[rel]=('link',val[1])
   else: data[rel]=None
if conflicts and policy=='Fail':
 json.dump(dict(conflicts=conflicts,operations=operations),open(out,'w')); sys.exit(3)
json.dump(dict(conflicts=conflicts,operations=operations,files=data,owners=owner),open(out,'w'))
PY
  then rc=0; else rc=$?; fi
  if ((rc != 0)); then
    rm -rf -- "$tmp"; trap - EXIT
    ((rc==3)) && fail 'consolidation conflicts; no changes made'
    fail 'consolidation planning failed'
  fi
  if ((DRY)); then printf '[DRY-RUN] consolidation planned; no filesystem changes\n'; rm -rf "$tmp"; trap - EXIT
  else
    plan=$tmp/plan.json; backup=$HARNESS/backups/sync-$RUN_ID
    mkdir -p -- "$HARNESS/manifests" "$HARNESS/backups" || fail 'cannot prepare manifest/backup directories'
    [[ ! -e $SKILLS ]] || cp -a -- "$SKILLS" "$backup" || fail 'backup failed'
    [[ -e $SKILLS ]] || mkdir -p -- "$backup"
    mkdir -p -- "$SKILLS"
    apply_error=$tmp/apply-error.txt
    if ! python3 - "$plan" "$SKILLS" 2>"$apply_error" <<'PY'
import json,os,shutil,sys
try:
 p=json.load(open(sys.argv[1])); dest=sys.argv[2]
 for rel,src in p['files'].items():
  out=os.path.join(dest,rel); os.makedirs(os.path.dirname(out),exist_ok=True)
  if isinstance(src,list) and src[0]=='link':
   if os.path.lexists(out): os.unlink(out)
   os.symlink(src[1],out)
  elif src:
   shutil.copy2(src,out)
  else: os.makedirs(out,exist_ok=True)
except Exception as exc:
 print(str(exc),file=sys.stderr)
 sys.exit(1)
PY
    then
      detail=$(cat -- "$apply_error")
      [[ -n $detail ]] || detail='unknown error'
      after=$(snapshot "$SKILLS")
      planjson=$(cat "$plan")
      payload=$(python3 - "$RUN_ID" "$POLICY" "$backup" "$before" "$after" "$planjson" "$detail" <<'PY'
import json,sys
run,policy,backup,before,after,plan,detail=sys.argv[1:]; p=json.loads(plan)
print(json.dumps(dict(version=1,run_id=run,status='Partial',conflict_policy=policy,backup_path=backup,sources=[x['source'] for x in p['operations']],conflicts=p['conflicts'],before=json.loads(before),after=json.loads(after),operations=p['operations'],error='ConsolidationFailed',error_detail=detail,completed=False),separators=(',',':')))
PY
      )
      mf=$(manifest_file); t=$(mktemp "${mf}.XXXXXX") || fail 'cannot create manifest temp'
      printf '%s\n' "$payload" >"$t" && mv -f -- "$t" "$mf" || fail 'cannot atomically write manifest'
      printf 'ERROR: consolidation failed: %s\n' "$detail" >&2
      exit 1
    fi
    after=$(snapshot "$SKILLS")
    planjson=$(cat "$plan")
    payload=$(python3 - "$RUN_ID" "$POLICY" "$backup" "$before" "$after" "$planjson" <<'PY'
import json,sys
run,policy,backup,before,after,plan=sys.argv[1:]; p=json.loads(plan)
print(json.dumps(dict(version=1,run_id=run,status='Applied',conflict_policy=policy,backup_path=backup,sources=[x['source'] for x in p['operations']],conflicts=p['conflicts'],before=json.loads(before),after=json.loads(after),operations=p['operations']),separators=(',',':')))
PY
)
    mf=$(manifest_file); t=$(mktemp "${mf}.XXXXXX") || fail 'cannot create manifest temp'; printf '%s\n' "$payload" >"$t" && mv -f -- "$t" "$mf" || fail 'cannot atomically write manifest'
    printf '[OK] consolidated; manifest: %s\n' "$mf"
    rm -rf -- "$tmp"; trap - EXIT
  fi
fi

if ((LINK)); then
  [[ -d $SKILLS ]] || fail "SSOT does not exist: $SKILLS"
  for target in "${TARGETS[@]}"; do
    parent=$(dirname -- "$target"); [[ -d $parent ]] || continue
    if [[ -L $target ]] && [[ $(readlink -f -- "$target" 2>/dev/null || true) == $(readlink -f -- "$SKILLS") ]]; then continue; fi
    if ((DRY)); then printf '[DRY-RUN] would link %s\n' "$target"; continue; fi
    backup="${target}.bak-${RUN_ID}"
    if [[ -e $target || -L $target ]]; then mv -- "$target" "$backup" || fail "cannot back up $target"; fi
    if [[ $MODE == symlink ]]; then ln -s -- "$SKILLS" "$target" || fail "cannot create symlink $target"
    else cp -a -- "$SKILLS" "$target" || fail "cannot copy skills to $target"; fi
  done
fi
if ((CHECK)); then
  failures=0
  if [[ -d $SKILLS ]]; then printf '[OK] SSOT exists: %s\n' "$SKILLS"; else printf '[ERR] SSOT missing: %s\n' "$SKILLS" >&2; failures=$((failures+1)); fi
  for target in "${TARGETS[@]}"; do
    [[ -e $(dirname -- "$target") ]] || continue
    if [[ $MODE == symlink ]]; then
      if [[ -L $target && $(readlink -f -- "$target" 2>/dev/null || true) == $(readlink -f -- "$SKILLS" 2>/dev/null || true) ]]; then printf '[OK] target: %s\n' "$target"; else printf '[ERR] invalid target: %s\n' "$target" >&2; failures=$((failures+1)); fi
    elif [[ -d $target ]]; then printf '[OK] copy target: %s\n' "$target"; else printf '[ERR] missing target: %s\n' "$target" >&2; failures=$((failures+1)); fi
  done
  ((failures==0)) || status=1
fi
if ((PASS)); then printf '{"success":%s,"dry_run":%s}\n' "$([[ $status == 0 ]] && echo true || echo false)" "$([[ $DRY == 1 ]] && echo true || echo false)"; fi
exit "$status"
