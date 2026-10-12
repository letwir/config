#!/usr/bin/env bash
# search.sh - Web search via w3m (DuckDuckGo lite primary, html fallback)
# Usage: search.sh <query words...> [-n <count>] [--json]

set -u

NUM=5
JSON_MODE=false
QUERY_ARGS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    -n)
      if [[ -z "${2:-}" || "$2" =~ ^- ]]; then
        echo "Error: -n requires a numeric argument" >&2
        exit 1
      fi
      NUM="$2"
      shift 2
      ;;
    -n*)
      NUM="${1#-n}"
      shift
      ;;
    --json)
      JSON_MODE=true
      shift
      ;;
    -h|--help)
      echo "Usage: search.sh <query words...> [-n <count>] [--json]"
      echo ""
      echo "Arguments:"
      echo "  <query words...>  Search query keywords"
      echo "  -n <count>        Max results to return (default: 5)"
      echo "  --json            Output results in JSON format"
      echo "  -h, --help        Show this help message"
      exit 0
      ;;
    *)
      QUERY_ARGS+=("$1")
      shift
      ;;
  esac
done

QUERY="${QUERY_ARGS[*]}"

if [[ -z "$QUERY" ]]; then
  echo "Error: Query is required." >&2
  echo "Usage: search.sh <query words...> [-n <count>] [--json]" >&2
  exit 1
fi

if ! [[ "$NUM" =~ ^[0-9]+$ ]] || [[ "$NUM" -le 0 ]]; then
  echo "Error: -n must be a positive integer" >&2
  exit 1
fi

if ! command -v w3m >/dev/null 2>&1; then
  echo "Error: w3m is required but not found." >&2
  exit 1
fi

# URL encode (application/x-www-form-urlencoded)
urlencode() {
  local s="$1"
  if command -v jq >/dev/null 2>&1; then
    jq -rn --arg v "$s" '$v|@uri' | sed 's/%20/+/g'
    return
  fi
  local out="" hex i byte safe
  hex=$(printf '%s' "$s" | od -An -tx1 | tr -d ' \n')
  for (( i=0; i<${#hex}; i+=2 )); do
    byte="${hex:i:2}"
    safe=0
    case "$byte" in
      2d|2e|5f|7e) safe=1 ;;
      *)
        if [[ "$byte" =~ ^[0-9a-f]+$ ]]; then
          local ch
          ch=$(printf "\\x$byte")
          if [[ "$ch" =~ [A-Za-z0-9] ]]; then
            safe=1
          fi
        fi
        ;;
    esac
    if [[ $safe -eq 1 ]]; then
      out+=$(printf "\\x$byte")
    else
      out+="%${byte^^}"
    fi
  done
  printf '%s' "$out"
}

# Fetch page source via w3m (identity encoding so w3m never emits undecoded br/gzip)
w3m_fetch() {
  w3m -dump_source -no-cookie -o accept_encoding=identity -o document_charset=UTF-8 "$1" 2>/dev/null
}

# URL decode percent-encoded strings
urldecode() {
  local url="${1//+/ }"
  printf '%b' "${url//%/\\x}"
}

# Output formatted JSON
output_json() {
  local engine="$1"
  local query="$2"
  shift 2

  if command -v jq >/dev/null 2>&1; then
    local json_items="[]"
    while [[ $# -ge 2 ]]; do
      local t="$1"
      local u="$2"
      shift 2
      json_items=$(echo "$json_items" | jq --arg t "$t" --arg u "$u" '. + [{title: $t, url: $u}]')
    done
    jq -n \
      --arg eng "$engine" \
      --arg q "$query" \
      --argjson items "$json_items" \
      '{
        engine: $eng,
        query: $q,
        transport: "w3m",
        sources: $items,
        results: $items
      }'
  elif command -v python3 >/dev/null 2>&1; then
    python3 -c '
import sys, json
engine = sys.argv[1]
query = sys.argv[2]
raw = sys.argv[3:]
items = []
for i in range(0, len(raw), 2):
    if i + 1 < len(raw):
        items.append({"title": raw[i], "url": raw[i+1]})
out = {
    "engine": engine,
    "query": query,
    "transport": "w3m",
    "sources": items,
    "results": items
}
print(json.dumps(out, ensure_ascii=False, indent=2))
' "$engine" "$query" "$@"
  fi
}

# Output formatted text
output_text() {
  local engine_title="$1"
  shift

  echo "=== $engine_title ==="
  echo ""
  echo "【検索結果】"
  local idx=1
  while [[ $# -ge 2 ]]; do
    local t="$1"
    local u="$2"
    shift 2
    echo "  $idx. $t"
    echo "     $u"
    ((idx++))
  done
}

# Decode common HTML entities
decode_entities() {
  sed -e 's/&quot;/"/g' -e "s/&#x27;/'/g" -e "s/&#39;/'/g" -e 's/&lt;/</g' -e 's/&gt;/>/g' -e 's/&amp;/\&/g'
}

# Primary: DuckDuckGo lite endpoint (no JS required)
search_ddg_lite() {
  local encoded
  encoded=$(urlencode "$QUERY") || return 1
  local src
  src=$(w3m_fetch "https://lite.duckduckgo.com/lite/?q=$encoded")
  [[ -z "$src" ]] && return 1

  local parsed
  parsed=$(echo "$src" | awk '
    BEGIN { RS="</a>" }
    {
      if (match($0, /<a[^>]*class=.?result-link[^>]*>/)) {
        head=substr($0, RSTART, RLENGTH)
        tail=substr($0, RSTART+RLENGTH)
        sub(/<\/?a[^>]*>.*/, "", tail)
        gsub(/\r?\n[ \t]*/, " ", tail)
        url=""
        if (match(head, /uddg=[^&"'"'"'>]+/)) { url=substr(head, RSTART+5, RLENGTH-5) }
        if (url != "") { print tail "\t" url }
      }
    }
  ' | decode_entities)

  [[ -z "$parsed" ]] && return 1

  local items=()
  local count=0
  while IFS=$'\t' read -r raw_title raw_url; do
    [[ -z "$raw_url" ]] && continue
    local clean_url
    clean_url=$(urldecode "$raw_url")
    items+=("$raw_title" "$clean_url")
    ((count++))
    if [[ $count -ge $NUM ]]; then
      break
    fi
  done <<< "$parsed"

  if [[ ${#items[@]} -eq 0 ]]; then
    return 1
  fi

  if [[ "$JSON_MODE" == "true" ]]; then
    output_json "duckduckgo-lite" "$QUERY" "${items[@]}"
  else
    output_text "DuckDuckGo Lite (w3m)" "${items[@]}"
  fi
  return 0
}

# Fallback: DuckDuckGo html endpoint
search_ddg_html() {
  local encoded
  encoded=$(urlencode "$QUERY") || return 1
  local src
  src=$(w3m_fetch "https://html.duckduckgo.com/html/?q=$encoded")
  [[ -z "$src" ]] && return 1

  local parsed
  parsed=$(echo "$src" | awk '
    BEGIN { in_title=0; href=""; title="" }
    /class="result__title"/ { in_title=1; next }
    in_title && /class="result__a"/ {
      if (match($0, /href="([^"]+)"/, m)) { href = m[1] }
      gsub(/.*<a[^>]*>/, "")
      gsub(/<\/a>.*/, "")
      gsub(/<[^>]+>/, "")
      title = $0
      if (href !~ /duckduckgo\.com\/y\.js/) {
        if (match(href, /uddg=([^&]+)/, u)) { href = u[1] }
        print title "\t" href
      }
      in_title=0
    }
  ' | decode_entities)

  [[ -z "$parsed" ]] && return 1

  local items=()
  local count=0
  while IFS=$'\t' read -r raw_title raw_url; do
    [[ -z "$raw_url" ]] && continue
    local clean_url
    clean_url=$(urldecode "$raw_url")
    items+=("$raw_title" "$clean_url")
    ((count++))
    if [[ $count -ge $NUM ]]; then
      break
    fi
  done <<< "$parsed"

  if [[ ${#items[@]} -eq 0 ]]; then
    return 1
  fi

  if [[ "$JSON_MODE" == "true" ]]; then
    output_json "duckduckgo-html" "$QUERY" "${items[@]}"
  else
    output_text "DuckDuckGo HTML (w3m)" "${items[@]}"
  fi
  return 0
}

# Main execution: w3m via DuckDuckGo lite, fall back to html endpoint
if ! search_ddg_lite; then
  if ! search_ddg_html; then
    echo "Error: Search failed for query: $QUERY" >&2
    exit 1
  fi
fi
