#!/usr/bin/env bash
# Local-only Bash port of write-task-record.ps1.
set -Eeuo pipefail

fail() { printf 'ERROR: %s\n' "$1" >&2; exit 1; }
usage() { printf 'Usage: %s TASK_ID [options]\n' "${0##*/}" >&2; }
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
ROOT_PATH=$(cd -- "$SCRIPT_DIR/.." && pwd -P)
INGEST_STATUS=N_A EVALUATION_STATUS=N_A MIRROR_STATUS=N_A
INGEST_REASON=pending_approval EVALUATION_REASON=pending_approval MIRROR_REASON=not_requested
INGEST_RECEIPT= EVALUATION_RECEIPT= MIRROR_RECEIPT=
LOCK_TIMEOUT_MS=5000
while (($#)); do
  case $1 in
    --root) (($# >= 2)) || fail 'missing root'; ROOT_PATH=$2; shift 2;;
    --ingest-status) (($# >= 2)) || fail 'missing ingest status'; INGEST_STATUS=$2; shift 2;;
    --evaluation-status) (($# >= 2)) || fail 'missing evaluation status'; EVALUATION_STATUS=$2; shift 2;;
    --mirror-status) (($# >= 2)) || fail 'missing mirror status'; MIRROR_STATUS=$2; shift 2;;
    --ingest-reason) (($# >= 2)) || fail 'missing ingest reason'; INGEST_REASON=$2; shift 2;;
    --evaluation-reason) (($# >= 2)) || fail 'missing evaluation reason'; EVALUATION_REASON=$2; shift 2;;
    --mirror-reason) (($# >= 2)) || fail 'missing mirror reason'; MIRROR_REASON=$2; shift 2;;
    --ingest-receipt) (($# >= 2)) || fail 'missing ingest receipt'; INGEST_RECEIPT=$2; shift 2;;
    --evaluation-receipt) (($# >= 2)) || fail 'missing evaluation receipt'; EVALUATION_RECEIPT=$2; shift 2;;
    --mirror-receipt) (($# >= 2)) || fail 'missing mirror receipt'; MIRROR_RECEIPT=$2; shift 2;;
    --lock-timeout-ms) (($# >= 2)) || fail 'missing lock timeout'; LOCK_TIMEOUT_MS=$2; shift 2;;
    --*) fail "unknown option $1";;
    *) [[ -z ${TASK_ID:-} ]] || fail 'unexpected argument'; TASK_ID=$1; shift;;
  esac
done
[[ -n ${TASK_ID:-} ]] || { usage; fail 'task id required'; }
[[ $TASK_ID =~ ^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$ ]] || fail invalid_task_id
lower_task_id=${TASK_ID,,}
[[ ! $lower_task_id =~ ^(con|prn|aux|nul|com[1-9]|lpt[1-9])([.]|$) ]] || fail invalid_task_id
[[ $LOCK_TIMEOUT_MS =~ ^[0-9]+$ ]] && ((LOCK_TIMEOUT_MS >= 1 && LOCK_TIMEOUT_MS <= 60000)) || fail io_error

check_receipt() {
  local status=$1 reason=$2 receipt=$3
  case $status:$reason in
    SUCCESS:completed) [[ $receipt =~ ^[A-Za-z0-9._-]{1,64}$ ]] || fail io_error;;
    N_A:pending_approval|N_A:not_requested|N_A:boundary_unavailable|FAILED:process_error|FAILED:missing_executable) [[ -z $receipt ]] || fail io_error;;
    *) fail io_error;;
  esac
}
check_receipt "$INGEST_STATUS" "$INGEST_REASON" "$INGEST_RECEIPT"
check_receipt "$EVALUATION_STATUS" "$EVALUATION_REASON" "$EVALUATION_RECEIPT"
check_receipt "$MIRROR_STATUS" "$MIRROR_REASON" "$MIRROR_RECEIPT"

# Reject symlink/reparse components, including ancestors, before canonicalization is trusted.
check_components() {
  local p=$1 rest part current=/
  [[ $p == /* ]] || p="$PWD/$p"
  rest=${p#/}
  while [[ -n $rest ]]; do
    part=${rest%%/*}; rest=${rest#*/}; [[ $part == "$rest" ]] && rest=
    [[ -z $part || $part == . ]] && continue
    if [[ $part == .. ]]; then current=${current%/*}; [[ -n $current ]] || current=/; continue; fi
    current="${current%/}/$part"
    [[ ! -L $current ]] || fail reparse_point
  done
  return 0
}
check_components "$ROOT_PATH"
[[ -d $ROOT_PATH ]] || fail root_missing
ROOT_PATH=$(cd -- "$ROOT_PATH" 2>/dev/null && pwd -P) || fail root_missing
MEMO="$ROOT_PATH/memo"
if [[ ! -e $MEMO ]]; then mkdir -- "$MEMO" 2>/dev/null || fail io_error; fi
[[ -d $MEMO ]] || fail path_escape
check_components "$MEMO"
shopt -s nullglob
matches=("$MEMO"/*)
TASK_DIR=
for p in "${matches[@]}"; do
  base=${p##*/}; lower_base=${base,,}; lower_task=${TASK_ID,,}
  [[ $lower_base == "$lower_task" ]] && { [[ $base == "$TASK_ID" ]] || fail invalid_task_id; TASK_DIR=$p; }
done
[[ -n $TASK_DIR && -d $TASK_DIR ]] || fail artifact_missing
check_components "$TASK_DIR"
TASK_DIR=$(cd -- "$TASK_DIR" && pwd -P)
case "$TASK_DIR/" in "$ROOT_PATH"/*/) ;; *) fail path_escape;; esac

LOCK="$TASK_DIR/.record.lock"
check_components "$LOCK"
exec {lock_fd}>"$LOCK" || fail io_error
if ! flock -w "$(awk -v m="$LOCK_TIMEOUT_MS" 'BEGIN{printf "%.3f",m/1000}')" "$lock_fd"; then fail lock_timeout; fi
TMP_FILES=()
cleanup() { for f in "${TMP_FILES[@]}"; do [[ -e $f ]] && rm -f -- "$f"; done; return 0; }
trap cleanup EXIT
names=(knowledge diary walkthrough evaluation)
files=(knowledge.md diary.md walkthrough.md evaluation.json)
hashes=() sizes=()
for i in "${!files[@]}"; do
  path="$TASK_DIR/${files[$i]}"
  [[ -f $path ]] || fail artifact_missing
  check_components "$path"
  a=$(mktemp "$TASK_DIR/.record.read.XXXXXX") || fail io_error
  b=$(mktemp "$TASK_DIR/.record.read.XXXXXX") || fail io_error
  TMP_FILES+=("$a" "$b")
  cat -- "$path" >"$a" || fail io_error
  cat -- "$path" >"$b" || fail io_error
  cmp -s -- "$a" "$b" || fail artifact_drift
  hashes+=("$(sha256sum "$a" | cut -d' ' -f1)")
  sizes+=("$(wc -c <"$a" | tr -d '[:space:]')")
done
python3 - "$TASK_DIR/evaluation.json" "$TASK_ID" <<'PY' || fail evaluation_invalid
import json, sys
try:
    with open(sys.argv[1], 'r', encoding='utf-8') as f:
        obj = json.load(f)
    if not isinstance(obj, dict) or obj.get('task_id') != sys.argv[2]:
        raise ValueError()
except Exception:
    raise SystemExit(1)
PY
# Re-read every artifact after metadata validation to detect changes across the snapshot window.
for i in "${!files[@]}"; do
  path="$TASK_DIR/${files[$i]}"
  c=$(mktemp "$TASK_DIR/.record.read.XXXXXX") || fail io_error; TMP_FILES+=("$c")
  cat -- "$path" >"$c" || fail io_error
  [[ $(sha256sum "$c" | cut -d' ' -f1) == "${hashes[$i]}" && $(wc -c <"$c" | tr -d '[:space:]') == "${sizes[$i]}" ]] || fail artifact_drift
done
STAMP=$(date -u +'%Y-%m-%dT%H:%M:%S.%NZ')
TMP_RECORD=$(mktemp "$TASK_DIR/.record.XXXXXX") || fail io_error
TMP_FILES+=("$TMP_RECORD")
python3 - "$TASK_ID" "$STAMP" "$INGEST_STATUS" "$INGEST_REASON" "$INGEST_RECEIPT" "$EVALUATION_STATUS" "$EVALUATION_REASON" "$EVALUATION_RECEIPT" "$MIRROR_STATUS" "$MIRROR_REASON" "$MIRROR_RECEIPT" "${hashes[*]}" "${sizes[*]}" >"$TMP_RECORD" <<'PY'
import json, sys
(task, stamp, is_, ir, iid, es, er, eid, ms, mr, mid, hs, bs) = sys.argv[1:]
hashes, sizes = hs.split(), list(map(int, bs.split()))
names = ['knowledge', 'diary', 'walkthrough', 'evaluation']
files = ['knowledge.md', 'diary.md', 'walkthrough.md', 'evaluation.json']
receipts = []
for status, reason, receipt in [(is_,ir,iid),(es,er,eid),(ms,mr,mid)]:
    item = {'status': status, 'reason_code': reason}
    if status == 'SUCCESS': item['receipt_id'] = receipt
    receipts.append(item)
record = {'schema_version': 1, 'task_id': task, 'generated_at': stamp,
          'artifact_root': 'memo/' + task,
          'artifacts': {n: {'path': 'memo/' + task + '/' + f, 'sha256': h, 'bytes': b}
                        for n,f,h,b in zip(names,files,hashes,sizes)},
          'llm_memory_ingest': receipts[0], 'llm_memory_eval': receipts[1],
          'external_diary_mirror': receipts[2],
          'memory_report_complete': is_ == 'SUCCESS' and es == 'SUCCESS'}
print(json.dumps(record, ensure_ascii=False, separators=(',', ':')))
PY
[[ -s $TMP_RECORD ]] || fail io_error
RECORD="$TASK_DIR/record.json"
if [[ -e $RECORD || -L $RECORD ]]; then [[ -f $RECORD && ! -L $RECORD ]] || fail reparse_point; check_components "$RECORD"; fi
mv -f -- "$TMP_RECORD" "$RECORD" || fail replace_failed
printf '%s\n' "$(cat -- "$RECORD")"
