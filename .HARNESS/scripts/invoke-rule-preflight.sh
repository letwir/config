#!/usr/bin/env bash
set -euo pipefail

usage() { echo "usage: $0 --task TASK [--tag TAG ...] [--root FILE] [--out FILE] [--detailed] [--verify RECEIPT]" >&2; exit 2; }
fail() { echo "Rule preflight failed: $*" >&2; exit 1; }
harness_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
root="$harness_root/rules/BOOTSTRAP.lrf" task= out= detailed=0 verify= max_depth=16 max_modules=32 max_bytes=1048576
declare -a tags=()
while (($#)); do
  case "$1" in
    --root) (($# >= 2)) || usage; root=$2; shift 2;;
    --task) (($# >= 2)) || usage; task=$2; shift 2;;
    --tag) (($# >= 2)) || usage; tags+=("$2"); shift 2;;
    --out) (($# >= 2)) || usage; out=$2; shift 2;;
    --detailed) detailed=1; shift;;
    --verify) (($# >= 2)) || usage; verify=$2; shift 2;;
    --max-depth) (($# >= 2)) || usage; max_depth=$2; shift 2;;
    --max-modules) (($# >= 2)) || usage; max_modules=$2; shift 2;;
    --max-bytes) (($# >= 2)) || usage; max_bytes=$2; shift 2;;
    *) usage;;
  esac
done
valid_name='^[a-z][a-z0-9.-]*$'
[[ $task =~ $valid_name || -n $verify ]] || usage
for tag in "${tags[@]}"; do [[ $tag =~ $valid_name ]] || fail "tag invalid: $tag"; done
for n in "$max_depth" "$max_modules" "$max_bytes"; do [[ $n =~ ^[1-9][0-9]*$ ]] || fail 'limits must be positive'; done

python3 - "$harness_root" "$root" "$task" "$out" "$detailed" "$verify" "$max_depth" "$max_modules" "$max_bytes" "${tags[@]}" <<'PY'
import hashlib,json,os,re,sys
from pathlib import Path
hroot=Path(sys.argv[1]).resolve(); root_arg,task,out,detailed,verify=sys.argv[2:7]
max_depth,max_modules,max_bytes=map(int,sys.argv[7:10]); tags=sys.argv[10:]
def fail(s): raise ValueError(s)
def safe(path):
 p=Path(path).resolve(strict=True)
 try: p.relative_to(hroot)
 except ValueError: fail(f'path escapes harness root: {p}')
 if not p.is_file(): fail(f'missing reference: {p}')
 cur=p
 while cur != hroot:
  if cur.is_symlink(): fail(f'reparse-point path is forbidden: {cur}')
  cur=cur.parent
 return p
def output(obj): print(json.dumps(obj,separators=(',',':'),ensure_ascii=False))
try:
 if verify:
  rec=json.loads(Path(verify).read_text(encoding='utf-8'))
  if rec.get('schema')!='lrf-preflight/v1' or rec.get('parser')!='invoke-rule-preflight.ps1/v1': fail('unknown receipt schema or parser')
  if not re.fullmatch(r'[a-z][a-z0-9.-]*',str(rec.get('task',''))): fail('invalid receipt task')
  files=rec.get('files',[]); selected=rec.get('selected',[])
  if not files or not selected: fail('receipt graph is empty')
  paths=set()
  for f in files:
   p=safe(f['path']); paths.add(str(p))
   if hashlib.sha256(p.read_bytes()).hexdigest().upper()!=f['sha256'].upper(): fail(f'stale receipt: {p}')
  rr=str(Path(rec['root']).resolve())
  if rr not in paths: fail('receipt root is not in graph')
  for s in selected:
   if str(Path(s['path']).resolve()) not in paths or not re.fullmatch(r'[a-z][a-z0-9.-]*',s['id']) or s['kind'] not in ('R','L') or int(s['line'])<1: fail('invalid selected record')
  output({'schema':'lrf-preflight-verify/v1','status':'PASS','file_count':len(files)}); sys.exit(0)
 root=safe(root_arg); allowed_mod={'MUST','MUST_NOT','MAY'}
 effects={'RO_LOCAL','RO_PUBLIC','LW_SCOPE','EXT_WRITE','RELEASE','LIVE_WRITE','VCS_WRITE','DESTRUCT','CRED','CHARGE','PROD_DEP'}
 guards={'task','tag','change','phase','fact','file','event','case'}
 ids=set(); visited=set(); stack=set(); files=[]; selected=[]; total=0
 lint=hroot/'skills/harness-lint/harness-lint'
 def visit(path,depth):
  nonlocal_dummy=None
  global total
  key=str(path)
  if depth>max_depth: fail(f'depth limit exceeded ({max_depth})')
  if key in stack: fail(f'reference cycle: {path}')
  if key in visited: return
  if len(visited)>=max_modules: fail(f'module limit exceeded ({max_modules})')
  stack.add(key)
  raw=path.read_bytes();
  try: text=raw.decode('utf-8')
  except UnicodeDecodeError: fail(f'invalid UTF-8: {path}')
  total+=len(raw)
  if total>max_bytes: fail(f'byte limit exceeded ({max_bytes})')
  lines=text.splitlines(); hi=next((i for i,x in enumerate(lines) if re.fullmatch(r'@lrf=1\|aud=[^|]+\|scope=[^|]+',x)),None)
  if hi is None: fail(f'invalid header: {path}')
  records=[]
  for i,line in enumerate(lines[hi+1:],hi+2):
   line=line.strip()
   if not line or line.startswith('#'): continue
   parts=line.split('|',5)
   if len(parts)!=6: fail(f'invalid field count at {path}:{i}')
   kind,rid,guard,mod,effect,payload=parts
   if kind not in ('R','L'): fail(f'unknown record kind at {path}:{i}')
   if not re.fullmatch(r'[a-z][a-z0-9.-]*',rid): fail(f'invalid id at {path}:{i}')
   if rid in ids: fail(f"duplicate id '{rid}'")
   ids.add(rid)
   if mod not in allowed_mod: fail(f'invalid modality at {path}:{i}')
   if effect not in effects: fail(f'invalid effect at {path}:{i}')
   if not payload: fail(f'empty payload at {path}:{i}')
   matched=True
   if guard!='*':
    for atom in guard.split('&'):
     m=re.fullmatch(r'([a-z]+):([A-Za-z0-9._-]+)',atom)
     if not m: fail(f'invalid guard: {guard}')
     k,v=m.groups()
     if k not in guards: fail(f'unknown guard key: {k}')
     matched &= (v==task if k=='task' else v in tags if k=='tag' else False)
   records.append((kind,rid,payload,i,matched))
  if lint.is_file():
   import subprocess
   if subprocess.run([str(lint),'-path',str(path),'-level','0','-strict'],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL).returncode: fail(f'harness-lint rejected: {path}')
  files.append({'path':str(path),'sha256':hashlib.sha256(raw).hexdigest().upper()}); visited.add(key)
  for kind,rid,payload,line,matched in records:
   if not matched: continue
   selected.append({'path':str(path),'line':line,'id':rid,'kind':kind})
   if kind=='L':
    m=re.search(r'\[[^\]]+\]\(([^)]+\.lrf)\)',payload)
    if not m: fail(f'invalid L reference at {path}:{line}')
    candidate=m.group(1)
    if os.path.isabs(candidate): fail(f'rooted reference is forbidden: {candidate}')
    visit(safe(path.parent/candidate),depth+1)
  stack.remove(key)
 visit(root,0)
 receipt={'schema':'lrf-preflight/v1','status':'PASS','parser':'invoke-rule-preflight.ps1/v1','task':task,'tags':tags,'root':str(root),'files':files,'selected':selected,'total_bytes':total}
 if out:
  op=Path(out).resolve();
  if not op.parent.is_dir(): fail(f'receipt parent missing: {op.parent}')
  op.write_text(json.dumps(receipt,indent=2,ensure_ascii=False)+'\n',encoding='utf-8')
 if detailed: output(receipt)
 else: output({'schema':'lrf-preflight-summary/v1','status':'PASS','task':task,'tags':tags,'file_count':len(files),'selected_count':len(selected),'total_bytes':total,'receipt_path':str(Path(out).resolve()) if out else None})
except Exception as e:
 print(f'Rule preflight failed: {e}',file=sys.stderr); sys.exit(1)
PY
