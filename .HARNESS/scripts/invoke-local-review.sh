#!/usr/bin/env bash
set -Eeuo pipefail
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
QUERY= TARGETS= ROOT_PATH="$SCRIPT_DIR/.." OUTPUT_PATH= JSON=0
fail(){ printf '%s\n' "${1#local-review:}" >&2; exit 2; }
while (($#)); do case $1 in
 --Query|-Query) (($#>=2))||fail local-review:missing_query; QUERY=$2; shift;;
 --Target|-Target) (($#>=2))||fail local-review:missing_target; TARGETS+="${TARGETS:+,}$2"; shift;;
 --RootPath|-RootPath) (($#>=2))||fail local-review:missing_root; ROOT_PATH=$2; shift;;
 --OutputPath|-OutputPath) (($#>=2))||fail local-review:missing_output; OUTPUT_PATH=$2; shift;;
 --Json|-Json) JSON=1;; *) fail "local-review:unknown_argument_$1";; esac; shift; done
json_stop(){ local code=$1; if ((JSON)) || [[ -z $OUTPUT_PATH ]]; then python3 - "$code" <<'PY'
import json,sys
print(json.dumps({'schema_version':1,'overall_status':'STOP','error_code':sys.argv[1]},separators=(',',':')))
PY
fi; exit 2; }
[[ -n ${QUERY//[[:space:]]/} && ${#QUERY} -le 512 && $QUERY != *$'\n'* && $QUERY != *$'\r'* ]] || json_stop invalid_query
[[ -n $ROOT_PATH && -d $ROOT_PATH && ! -L $ROOT_PATH ]] || json_stop invalid_root
ROOT_PATH=$(cd -- "$ROOT_PATH" && pwd -P) || json_stop invalid_root
IFS=, read -r -a raw_targets <<< "$TARGETS"
((${#raw_targets[@]}>=1 && ${#raw_targets[@]}<=32)) || json_stop invalid_target_count
paths=(); seen='|'
for name in "${raw_targets[@]}"; do name=${name##+([[:space:]])}; name=${name%%+([[:space:]])}; [[ -n $name ]] || json_stop invalid_target_count
 [[ $name != /* && ! $name =~ (^|/)\.\.(/|$) && ! $name =~ [*?\[\]] ]] || json_stop invalid_target_path
 full="$ROOT_PATH/$name"; [[ -f $full && ! -L $full ]] || json_stop target_missing_or_not_file
 IFS=/ read -r -a components <<< "$name"; cur=$ROOT_PATH; for c in "${components[@]}"; do cur="$cur/$c"; [[ ! -L $cur ]] || json_stop reparse_target; done
 rel=${name#./}; [[ $seen != *"|$rel|"* ]] || json_stop duplicate_target; seen+="$rel|"; paths+=("$rel")
done
if [[ -n $OUTPUT_PATH ]]; then [[ $OUTPUT_PATH == "$ROOT_PATH"/* ]] || json_stop output_escape; fi
hash_text(){ printf %s "$1" | sha256sum | awk '{print toupper($1)}'; }
fixed=$(mktemp "${TMPDIR:?TMPDIR is not set}/local-review.XXXXXX"); status=$(mktemp "${TMPDIR:?TMPDIR is not set}/local-review.XXXXXX"); diffcheck=$(mktemp "${TMPDIR:?TMPDIR is not set}/local-review.XXXXXX"); trap 'rm -f -- "$fixed" "$status" "$diffcheck"' EXIT
set +e
rg --fixed-strings --line-number --no-heading --color never --max-count 20 "$QUERY" -- "${paths[@]}" >"$fixed" 2>&1; rr=$?
git -C "$ROOT_PATH" status --short -- "${paths[@]}" >"$status" 2>&1; sr=$?
git -C "$ROOT_PATH" diff --check -- "${paths[@]}" >"$diffcheck" 2>&1; dr=$?
set -e
((rr==0 || rr==1)) && ((sr==0)) && ((dr==0)); ok=$?
receipt_tmp=$(mktemp "${TMPDIR:?TMPDIR is not set}/local-review-json.XXXXXX")
python3 - "$ROOT_PATH" "$QUERY" "$rr" "$sr" "$dr" "$ok" "$(hash_text "$(cat "$fixed")")" "$(hash_text "$(cat "$status")")" "$(hash_text "$(cat "$diffcheck")")" "${paths[@]}" <<'PY' > "$receipt_tmp"
import datetime,hashlib,json,os,sys,uuid
root,q,rr,sr,dr,ok,rh,sh,dh,*paths=sys.argv[1:]
def digest(p):
 h=hashlib.sha256()
 with open(os.path.join(root,p),'rb') as f:
  for b in iter(lambda:f.read(131072),b''):h.update(b)
 return h.hexdigest().upper()
receipt={'schema_version':1,'receipt_id':str(uuid.uuid4()),'repo_root':root,'query':{'mode':'fixed_string','sha256':hashlib.sha256(q.encode()).hexdigest().upper(),'length':len(q)},'targets':[{'path':p,'sha256':digest(p)} for p in paths],'commands':{'rg':{'executable':'rg','exit_code':int(rr),'output_sha256':rh,'bounded':True},'git_status':{'executable':'git status --short','exit_code':int(sr),'output_sha256':sh,'bounded':True},'git_diff_check':{'executable':'git diff --check','exit_code':int(dr),'output_sha256':dh,'bounded':True}},'overall_status':'PASS' if ok=='0' else 'FAIL','generated_at':datetime.datetime.now(datetime.timezone.utc).isoformat()}
print(json.dumps(receipt,separators=(',',':')))
PY
rc=0; ((ok==0)) || rc=1
if [[ -n $OUTPUT_PATH ]]; then mkdir -p -- "$(dirname -- "$OUTPUT_PATH")" || json_stop output_error; tmp=$(mktemp "${TMPDIR:?TMPDIR is not set}/local-review-receipt.XXXXXX") || json_stop output_error; cp -- "$receipt_tmp" "$tmp"; mv -f -- "$tmp" "$OUTPUT_PATH" || { rm -f -- "$tmp"; json_stop output_error; }; ((JSON)) && cat "$OUTPUT_PATH"; elif ((JSON)); then cat "$receipt_tmp"; fi
rm -f -- "$receipt_tmp"
exit "$rc"
