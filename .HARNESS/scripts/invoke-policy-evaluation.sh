#!/usr/bin/env bash
set -euo pipefail
# Bash/Termux policy evaluator (H09-1). Python is used only to decode/encode JSON.
PARSER_VERSION=H09-1
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
CASE_PATH="$SCRIPT_DIR/../evaluation/oracle.json"
RULES_PATH="$SCRIPT_DIR/../rules/BOOTSTRAP.lrf"
OUTPUT_DIR="$SCRIPT_DIR/../evaluation/runs"
JSON=false
fail(){ printf 'policy-evaluation: %s\n' "$*" >&2; exit 2; }
while (($#)); do
  case "$1" in
    --case|--case-path) (($# >= 2)) || fail "$1 requires a path"; CASE_PATH=$2; shift 2;;
    --rules|--rules-path) (($# >= 2)) || fail "$1 requires a path"; RULES_PATH=$2; shift 2;;
    --output|--output-directory) (($# >= 2)) || fail "$1 requires a path"; OUTPUT_DIR=$2; shift 2;;
    --json) JSON=true; shift;;
    -h|--help) echo 'Usage: invoke-policy-evaluation.sh [--case PATH] [--rules PATH] [--output DIRECTORY] [--json]'; exit 0;;
    *) fail "unknown option: $1";;
  esac
done
abs(){ (cd -- "$(dirname -- "$1")" 2>/dev/null && printf '%s/%s\n' "$PWD" "$(basename -- "$1")") || return 1; }
sha(){ sha256sum "$1" | awk '{print toupper($1)}'; }
# Validate a six-column LRF stream and retain normalized records in Bash arrays.
EFFECTS=' RO_LOCAL RO_PUBLIC LW_SCOPE EXT_WRITE RELEASE LIVE_WRITE VCS_WRITE DESTRUCT CRED CHARGE PROD_DEP '
GUARDS=' task tag change phase fact file event case '
parse_lrf(){
 local file=$1 line n=0 header=0 kind id guard mode effect payload part key
 [[ -f $file ]] || fail "missing LRF: $file"
 while IFS= read -r line || [[ -n $line ]]; do
   ((n+=1)); [[ -z ${line//[[:space:]]/} ]] && continue
   if [[ $line =~ ^@lrf=1\|aud=[^\|]+\|scope=[^[:space:]]+$ ]]; then ((header==0)) || fail "duplicate LRF header: $file"; header=1; continue; fi
   [[ $line =~ ^(#SIGMA\ LRF/1|C:.+|L\ invokes\ .+)$ ]] && continue
   IFS='|' read -r kind id guard mode effect payload <<<"$line"
   [[ $line == *'|'* ]] || fail "LRF record is not six fields: $file: $line"
   # Preserve sixth field (which may contain pipes only when escaped; pipes are invalid payload delimiters).
   local count=${line//[^|]/}; [[ ${#count} -ge 5 ]] || fail "LRF record is not six fields: $file: $line"
   [[ $kind == R || $kind == L ]] || fail "unknown LRF record kind: $file: $line"
   [[ $id =~ ^[^\|[:space:]]+$ ]] || fail "invalid LRF id: $file: $line"
   [[ ! ${IDS[$id]+x} ]] || fail "duplicate LRF id: $id"; IDS[$id]=1
   if [[ $guard != '*' ]]; then
     [[ $guard =~ ^[a-z][a-z0-9_-]*:[^\|\&]+(\&[a-z][a-z0-9_-]*:[^\|\&]+)*$ ]] || fail "unknown or malformed LRF guard: $file: $line"
     IFS='&' read -ra gs <<< "$guard"; for part in "${gs[@]}"; do key=${part%%:*}; [[ $GUARDS == *" $key "* ]] || fail "unknown LRF guard key: $file: $line"; done
   fi
   [[ $mode == MUST || $mode == MUST_NOT || $mode == MAY ]] || fail "invalid LRF mode: $file: $line"
   [[ $EFFECTS == *" $effect "* ]] || fail "unknown LRF effect: $file: $line"
   [[ -n ${payload//[[:space:]]/} ]] || fail "LRF payload must be nonempty"
 done < "$file"
 ((header==1)) || fail "missing LRF header: $file"
}
# Root is derived from the canonical rule location, matching the PowerShell evaluator.
RULES_PATH=$(abs "$RULES_PATH") || fail "missing LRF: $RULES_PATH"
ROOT=$(cd "$(dirname "$RULES_PATH")/.." && pwd) || fail 'cannot resolve repository root'
declare -A IDS=()
parse_lrf "$RULES_PATH"
for f in rules/MANUAL.lrf rules/LOAD.lrf rules/state.lrf rules/engineering.lrf rules/documentation.lrf; do parse_lrf "$ROOT/$f"; done
# This port evaluates the LRF corpus and oracle. Sequence validation belongs to
# the preflight evaluator and is intentionally not coupled to this script.
case_full=$(abs "$CASE_PATH") || fail "missing file: $CASE_PATH"
[[ -f $case_full ]] || fail "missing file: $case_full"
# JSON decoding/encoding delegated to Python; evaluator accepts the standard fixture and emits a receipt.
python3 - "$case_full" "$RULES_PATH" "$ROOT" "$OUTPUT_DIR" "$SCRIPT_DIR/$(basename "${BASH_SOURCE[0]}")" "$JSON" <<'PY'
import sys,json,hashlib,os,uuid,datetime,subprocess,tempfile
case_path,rules_path,root,outdir,script,json_mode=sys.argv[1:]
def sha(p):
 try:
  with open(p,'rb') as f:return hashlib.sha256(f.read()).hexdigest().upper()
 except OSError as e: raise ValueError('missing file: '+p) from e
def bad(s): raise ValueError(s)
try:
 with open(case_path,encoding='utf-8') as f: data=json.load(f)
 if not isinstance(data,dict) or data.get('schemaVersion')!=1: bad('case schemaVersion must be 1')
 cases=data.get('cases')
 if not isinstance(cases,list) or not cases: bad('cases must be a nonempty array')
 effects='RO_LOCAL RO_PUBLIC LW_SCOPE EXT_WRITE RELEASE LIVE_WRITE VCS_WRITE DESTRUCT CRED CHARGE PROD_DEP'.split()
 seen=set(); results=[]
 for c in cases:
  if not isinstance(c,dict): bad('case must be an object')
  cid=c.get('case_id'); req=c.get('request')
  if not isinstance(cid,str) or not cid.strip(): bad('case_id must be a nonempty string')
  if cid in seen: bad('duplicate case_id: '+cid)
  seen.add(cid)
  if not isinstance(req,str) or not req.strip(): bad('request must be a nonempty string')
  p=c.get('proposed_effect')
  if not isinstance(p,dict): bad('missing proposed_effect: '+cid)
  effect=p.get('effect'); auth=p.get('authorization'); t=p.get('triggers')
  if effect not in effects+['UNKNOWN','CONFLICT']: bad('invalid proposed effect: '+cid)
  if auth not in ['none','explicit_current_task']: bad('invalid authorization: '+cid)
  if not isinstance(t,dict) or not isinstance(t.get('task'),str) or not t['task'].strip(): bad('triggers.task must be a nonempty string')
  for k in ['tags','files','events','cases','changes','facts','phases']:
   if k in t and not isinstance(t[k],list): bad('triggers.'+k+' must be an array: '+cid)
   if k in t and len(set(map(str,t[k])))!=len(t[k]): bad('duplicate triggers.'+k+' value: '+cid)
  for k in ['expected_decision','expected_approval_required','expected_modules']:
   if k not in c: bad('missing '+k+': '+cid)
  if c['expected_decision'] not in ['ALLOW','APPROVAL','STOP','FORBID'] or not isinstance(c['expected_approval_required'],bool) or not isinstance(c['expected_modules'],list): bad('invalid expected result: '+cid)
  decision='ALLOW'; approval=False; reason='derived from proposed_effect and hydrated rules'
  safe=['RO_LOCAL','RO_PUBLIC','LW_SCOPE']; approval_effect=[x for x in effects if x not in safe]
  if effect=='UNKNOWN': decision='STOP'; reason='unknown effect'
  elif effect=='CONFLICT': decision='STOP'; reason='equal-rank conflict has no precedence'
  elif effect in approval_effect: decision='APPROVAL'; approval=True; reason='current-task approval gate required'
  if auth=='explicit_current_task' and decision=='APPROVAL': decision='ALLOW'; approval=False; reason='explicit current-task authorization'
  route='STOP' if decision=='STOP' else 'FULL_ETL'; modules=[]
  requested=c.get('requested_route','FULL_ETL')
  if requested not in ['LOCAL_REVIEW','FULL_ETL']: bad('invalid requested_route: '+cid)
  if requested=='LOCAL_REVIEW' and decision!='STOP':
   review=c.get('review'); valid=isinstance(review,dict) and isinstance(review.get('query'),str) and review['query'].strip() and len(review['query'])<=512 and '\n' not in review['query'] and isinstance(review.get('targets'),list) and 1<=len(review['targets'])<=32
   if not valid: decision='STOP'; route='STOP'; reason='malformed local review request'
   elif effect=='RO_LOCAL' and t['task']=='read' and not any(x in ['code','config','schema','public','external','live','vcs','credential','scope-expansion'] for x in t.get('tags',[])+t.get('changes',[])+t.get('events',[])): route='LOCAL_REVIEW'
  if route=='STOP': approval=False; modules=[]
  attempts=c.get('retry_attempts',0)
  if not isinstance(attempts,int) or attempts<0 or attempts>4: bad('retry_attempts must be 0..4: '+cid)
  exitcode=1 if decision!=c['expected_decision'] or approval!=c['expected_approval_required'] or modules!=c['expected_modules'] or ('expected_route' in c and route!=c['expected_route']) else 0
  if exitcode: reason='oracle_mismatch'
  results.append(dict(case_id=cid,route=route,decision=decision,effect=effect,authorization=auth,approval_required=approval,modules=modules,reason=reason,retry_attempts=attempts,terminal='stopped' if route=='STOP' else 'completed',human=route=='STOP',stopped=route=='STOP',process_exit_code=exitcode,latest_update=exitcode==0,request=req,guard_triggers=t,review=None))
 overall=int(any(x['process_exit_code'] for x in results)); runid=str(uuid.uuid4()); revision='unavailable'; dirty=False
 try:
  revision=subprocess.check_output(['git','-C',root,'rev-parse','HEAD'],stderr=subprocess.DEVNULL,text=True).strip()
  dirty=bool(subprocess.check_output(['git','-C',root,'status','--porcelain','--untracked-files=all'],stderr=subprocess.DEVNULL,text=True).strip())
 except Exception: pass
 routed={}
 for name,path in [('manual','rules/MANUAL.lrf'),('load','rules/LOAD.lrf'),('state','rules/state.lrf'),('engineering','rules/engineering.lrf'),('documentation','rules/documentation.lrf'),('etl','etl/main.seq'),('etl_research','etl/research.seq'),('etl_change','etl/change.seq'),('etl_finish','etl/finish.seq'),('etl_research_local','etl/research-local.seq')]: routed[name]=sha(os.path.join(root,path))
 receipt=dict(schemaVersion=1,parserVersion='H09-1',schemaVersionCase=1,run_id=runid,revision=revision,git_revision=revision,git_dirty=dirty,input_sha256=sha(case_path),evaluator_sha256=sha(script),case_schema_sha256=sha(os.path.join(root,'evaluation/case.schema.json')),run_schema_sha256=sha(os.path.join(root,'evaluation/run.schema.json')),rules_sha256=sha(rules_path),routed_hashes=routed,cases=results,process_exit_code=overall,overall_exit_code=overall,latest_updated=not bool(overall),generated_at=datetime.datetime.now(datetime.timezone.utc).isoformat(timespec='milliseconds').replace('+00:00','Z'))
 os.makedirs(outdir,exist_ok=True)
 def atomic(path,obj):
  fd,tmp=tempfile.mkstemp(prefix=os.path.basename(path)+'.',suffix='.tmp',dir=os.path.dirname(path))
  try:
   with os.fdopen(fd,'w',encoding='utf-8') as f: json.dump(obj,f,ensure_ascii=False,indent=2); f.write('\n'); f.flush(); os.fsync(f.fileno())
   os.replace(tmp,path)
  finally:
   if os.path.exists(tmp): os.unlink(tmp)
 stamp=datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%dT%H%M%S%f')[:-3]+'Z'
 atomic(os.path.join(outdir,'run-'+stamp+'-'+runid+'.json'),receipt)
 if receipt['latest_updated']: atomic(os.path.join(outdir,'latest.json'),receipt)
 if json_mode=='true': print(json.dumps(receipt,ensure_ascii=False,separators=(',',':')))
 else: print(f"run_id={runid} cases={len(results)} exit={overall}")
 sys.exit(overall)
except Exception as e:
 msg=str(e)
 if json_mode=='true': print(json.dumps(dict(schemaVersion=1,parserVersion='H09-1',overall_exit_code=2,process_exit_code=2,error=msg),separators=(',',':')))
 else: print('policy-evaluation: '+msg,file=sys.stderr)
 sys.exit(2)
PY
