#!/usr/bin/env bash
set -Eeuo pipefail
TARGET= ARGS_BASE64= OUT_PATH= ERR_PATH=
while (($#)); do case $1 in
  --Target|-Target) (($#>=2)) || exit 2; TARGET=$2; shift;;
  --ArgsBase64|-ArgsBase64) (($#>=2)) || exit 2; ARGS_BASE64=$2; shift;;
  --OutPath|-OutPath) (($#>=2)) || exit 2; OUT_PATH=$2; shift;;
  --ErrPath|-ErrPath) (($#>=2)) || exit 2; ERR_PATH=$2; shift;;
  *) printf 'unknown argument: %s\n' "$1" >&2; exit 2;; esac; shift; done
[[ -n $TARGET && -n $ARGS_BASE64 && -n $OUT_PATH && -n $ERR_PATH ]] || { echo 'Target, ArgsBase64, OutPath, and ErrPath are required' >&2; exit 2; }
export CI=true NO_COLOR=1
python3 - "$ARGS_BASE64" <<'PY' | python3 -c 'import json,sys; a=json.load(sys.stdin); assert isinstance(a,list) and all(isinstance(x,str) for x in a); print(json.dumps(a,separators=(",",":")))' > "${TMPDIR:?TMPDIR is not set}/client-args.$$.json"
import base64,sys
raw=base64.b64decode(sys.argv[1],validate=True)
if len(raw)%2: raise SystemExit('invalid UTF-16 payload')
print(raw.decode('utf-16-le'))
PY
args_file="$TMPDIR/client-args.$$.json"
trap 'rm -f -- "$args_file"' EXIT
python3 - "$args_file" <<'PY' > "${TMPDIR:?TMPDIR is not set}/client-argv.$$.nul"
import json,sys
for a in json.load(open(sys.argv[1],encoding='utf-8')): sys.stdout.buffer.write(a.encode()+b'\0')
PY
argv_file="$TMPDIR/client-argv.$$.nul"; trap 'rm -f -- "$args_file" "$argv_file"' EXIT
mapfile -d '' -t args < "$argv_file"
set +e
"$TARGET" "${args[@]}" > "$OUT_PATH" 2> "$ERR_PATH"
rc=$?
set -e
exit "$rc"
