#!/usr/bin/env bash
set -Eeuo pipefail
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
MANIFEST_PATH="$SCRIPT_DIR/../rules/skill-runtime.json"; OUTPUT_DIRECTORY="$SCRIPT_DIR/../evaluation/runs"; JSON=0; NO_RECEIPT=0
while (($#)); do case $1 in
 --ManifestPath|-ManifestPath) (($#>=2))||exit 2; MANIFEST_PATH=$2; shift;;
 --OutputDirectory|-OutputDirectory) (($#>=2))||exit 2; OUTPUT_DIRECTORY=$2; shift;;
 --Json|-Json) JSON=1;; --NoReceipt|-NoReceipt) NO_RECEIPT=1;; *) printf 'unknown argument: %s\n' "$1" >&2; exit 2;; esac; shift; done
MANIFEST_PATH=$(realpath -m -- "$MANIFEST_PATH"); root=$(cd -- "$(dirname -- "$MANIFEST_PATH")/.." && pwd -P)
[[ $JSON == 1 ]] && export SKILL_RUNTIME_JSON=1
work=$(mktemp -d "${TMPDIR:?TMPDIR is not set}/skill-runtime.XXXXXX"); trap 'rm -rf -- "$work"' EXIT
set +e
python3 - "$MANIFEST_PATH" "$root" "$HOME" "$SCRIPT_DIR" "$work" "$OUTPUT_DIRECTORY" "$NO_RECEIPT" <<'PY'
import datetime,hashlib,json,os,re,shutil,signal,subprocess,sys,uuid
manifest,root,home,scripts,work,outdir,no_receipt=sys.argv[1:]
try:
 with open(manifest,encoding='utf-8-sig') as f:m=json.load(f)
 if m.get('schemaVersion')!=1:raise ValueError('unsupported manifest schemaVersion')
except Exception as e:
 print(json.dumps({'schemaVersion':1,'overall_exit_code':2,'error':str(e)},separators=(',',':')));raise SystemExit(2)
limit=int(m.get('maxOutputChars',4096)); ps=shutil.which('pwsh') or ''
def expand(s):
 if not s:return s
 return str(s).replace('{HarnessRoot}',root).replace('{UserHome}',home).replace('{PowerShell}',ps).replace('\\','/')
def final(p):return os.path.realpath(p)
def contained(p,roots):
 c=final(p);return any(c==final(r) or c.startswith(final(r)+os.sep) for r in roots)
def norm(raw):
 s=re.sub(r'(?i)(api[_-]?key|token|password|secret)\s*[=:]\s*[^\s,;]+',r'\1=<REDACTED>',raw)
 s=re.sub(r'(?i)bearer\s+[A-Za-z0-9._~-]+','Bearer <REDACTED>',s).replace('\r\n','\n').replace('\r','\n').strip(); truncated=len(s)>limit;s=s[:limit]
 return {'text':s,'truncated':truncated,'sha256':hashlib.sha256(s.encode()).hexdigest().upper()}
def run(path,args,timeout,cwd):
 try:
  p=subprocess.Popen([path,*map(str,args)],cwd=cwd,stdout=subprocess.PIPE,stderr=subprocess.PIPE,env={**os.environ,'CI':'true','NO_COLOR':'1'},start_new_session=True)
  try:o,e=p.communicate(timeout=timeout/1000); timed=False; code=p.returncode
  except subprocess.TimeoutExpired:
   os.killpg(p.pid,signal.SIGKILL);o,e=p.communicate();timed=True;code=-1
  return code,timed,o.decode(errors='replace')+ ('\n'+e.decode(errors='replace') if e else '')
 except Exception as e:return 127,False,str(e)
def hashfile(p):
 try:
  h=hashlib.sha256()
  with open(p,'rb') as f:
   for b in iter(lambda:f.read(131072),b''):h.update(b)
  return h.hexdigest().upper()
 except OSError:return None
def check(spec,client=False):
 name=str(spec.get('name',''));required=bool(spec.get('required'));path=expand(spec.get('target') if client else spec.get('path')); roots=[expand(x) for x in m.get('approvedRoots',[])]
 result={'name':name,'required':required,'path':path,'finalPath':None,'classification':'Unknown','executableHealthy':'Missing','catalogObserved':'NotObserved','representativeOperation':'UnTested','commandIdentity':{'executable':path,'args':[]},'fileSha256':None,'version':None,'output':None,'error':None}
 wrapper=expand(spec.get('wrapper')) if client else None
 if client:
  result['commandIdentity']={'wrapper':wrapper,'target':path,'args':spec.get('args',[]),'wrapperSha256':None,'targetSha256':None}
  if not os.path.isfile(path):result['classification']='Missing' if required else 'NotConfigured';return result
  if not wrapper or not os.path.isfile(wrapper):result.update(classification='WrapperMismatch',error='wrapper-missing');return result
 if not contained(path,roots):result.update(classification='PathRejected',error='path-outside-approved-roots');return result
 if not os.path.isfile(path):result['classification']='Missing' if required else 'NotConfigured';return result
 result.update(executableHealthy='Healthy',finalPath=final(path),fileSha256=hashfile(path));
 if spec.get('sha256') and result['fileSha256']!=str(spec['sha256']).upper():result.update(classification='HashMismatch',error='file-hash-mismatch');return result
 if client:
  if not contained(wrapper,roots):result.update(classification='WrapperMismatch',error='wrapper-outside-approved-roots');return result
  result['commandIdentity'].update(wrapperSha256=hashfile(wrapper),targetSha256=result['fileSha256'])
  code,timed,raw=run(wrapper,[path,*spec.get('args',[])],int(spec.get('timeoutMs',5000)),root); attempts=1
  for _ in range(int(spec.get('retryCount',0))):
   if not timed and code==0:break
   code,timed,raw=run(wrapper,[path,*spec.get('args',[])],int(spec.get('timeoutMs',5000)),root);attempts+=1
  n=norm(raw);result['output']={'exitCode':code,'timedOut':timed,'attempts':attempts,**n,'checkedBeforeNormalization':True}
  if timed:result.update(classification='Timeout',catalogObserved='Timeout',error='catalog-timeout')
  elif code:result.update(classification='CommandNonZero',catalogObserved='Unverified',error=f'catalog-exit-{code}')
  elif spec.get('requiredText') and spec['requiredText'] not in raw:result.update(classification='ContractMismatch',error='required-output-missing')
  else:result.update(catalogObserved='Unverified' if spec.get('mode')=='catalog-api-or-unverified' else 'Observed',classification='Unverified' if spec.get('mode')=='catalog-api-or-unverified' else 'Healthy')
  return result
 args=spec.get('helpArgs',[]);result['commandIdentity']['args']=args;code,timed,raw=run(path,args,int(spec.get('timeoutMs',m.get('defaultTimeoutMs',5000))),root);n=norm(raw);result['output']={'exitCode':code,'timedOut':timed,**n,'checkedBeforeNormalization':True}
 if timed:result.update(classification='Timeout',catalogObserved='Timeout',error='help-timeout');return result
 result['catalogObserved']='Observed'; expected=int(spec.get('requiredExitCode',0))
 if code!=expected:result.update(classification='CommandNonZero',error=f'help-exit-{code}')
 elif spec.get('requiredText') and spec['requiredText'] not in raw:result.update(classification='ContractMismatch',error='required-output-missing')
 else:result['classification']='Healthy'
 return result
tools=[check(x) for x in m.get('tools',[])];clients=[check(x,True) for x in m.get('clients',[])+m.get('optionalClients',[])]
bad={'Missing','PathRejected','Timeout','CommandNonZero','ContractMismatch','WrapperMismatch','HashMismatch','RepresentativeFailed'}
failed=any(x['required'] and x['classification'] in bad for x in tools+clients)
r={'schemaVersion':1,'run_id':str(uuid.uuid4()),'generated_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'overall_exit_code':int(failed),'tools':tools,'clients':clients,'latest_updated':not failed}
if no_receipt=='0':
 os.makedirs(outdir,exist_ok=True)
 def atomic(p):
  t=p+'.tmp.'+uuid.uuid4().hex
  with open(t,'w',encoding='utf-8') as f:json.dump(r,f,ensure_ascii=False,indent=2);f.write('\n');f.flush();os.fsync(f.fileno())
  os.replace(t,p)
 atomic(os.path.join(outdir,'skill-runtime-'+r['run_id']+'.json'))
 if not failed:atomic(os.path.join(outdir,'skill-runtime-latest.json'))
print(json.dumps(r,separators=(',',':')) if os.environ.get('SKILL_RUNTIME_JSON')=='1' else json.dumps(r,ensure_ascii=False,indent=2))
raise SystemExit(r['overall_exit_code'])
PY
rc=$?
exit "$rc"
