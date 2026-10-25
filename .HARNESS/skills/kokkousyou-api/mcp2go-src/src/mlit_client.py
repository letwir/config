import os
import sys
import argparse
import requests
import json
import winreg

# Windowsのターミナル出力での文字化け防止
if sys.platform == 'win32' and hasattr(sys.stdout, 'reconfigure'):
    sys.stdout.reconfigure(encoding='utf-8')
    sys.stderr.reconfigure(encoding='utf-8')

def get_env_from_registry(name):
    """Windowsレジストリから環境変数を取得するフォールバック処理"""
    for root, path in [
        (winreg.HKEY_CURRENT_USER, "Environment"),
        (winreg.HKEY_LOCAL_MACHINE, r"System\CurrentControlSet\Control\Session Manager\Environment")
    ]:
        try:
            with winreg.OpenKey(root, path, 0, winreg.KEY_READ) as key:
                value, _ = winreg.QueryValueEx(key, name)
                return value
        except Exception:
            pass
    return None

def get_api_key():
    api_key = os.environ.get('KOKKOU_API') or get_env_from_registry('KOKKOU_API')
    if not api_key:
        print("Error: KOKKOU_API environment variable is not set in Windows Environment or Registry.")
        sys.exit(1)
    return api_key

def search_mlit(term, size=5, dataset_id=None):
    api_key = get_api_key()
    url = "https://data-platform.mlit.go.jp/api/v1/"
    headers = {
        "Content-Type": "application/json",
        "apikey": api_key
    }

    # AttributeFilterの構築
    attribute_filter_str = ""
    if dataset_id:
        attribute_filter_str = f', attributeFilter: {{ attributeName: "DPF:dataset_id", is: "{dataset_id}" }}'

    query = f"""
    query {{
      search(term: "{term}", phraseMatch: true, first: 0, size: {size}{attribute_filter_str}) {{
        totalNumber
        searchResults {{
          id
          title
          metadata
        }}
      }}
    }}
    """

    payload = {"query": query}

    try:
        response = requests.post(url, headers=headers, json=payload, timeout=20)
        response.raise_for_status()
        result = response.json()
        return result.get("data", {}).get("search", {})
    except Exception as e:
        print(f"API Request failed: {e}")
        if 'response' in locals() and response.text:
            print(f"Response snippet: {response.text[:200]}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="国土交通データプラットフォーム (MLIT DPF) 簡易検索クライアント")
    parser.add_argument("term", type=str, help="検索するキーワード (例: '橋梁', 'トンネル')")
    parser.add_argument("--size", "-s", type=int, default=5, help="取得する最大件数 (デフォルト: 5)")
    parser.add_argument("--dataset", "-d", type=str, default=None, help="特定のデータセットIDでフィルタリング (例: 'rsdb_bridge')")
    parser.add_argument("--json-out", "-j", action="store_true", help="結果をそのままJSON形式で出力")

    args = parser.parse_args()

    print(f"=== MLIT DPF Search: '{args.term}' (Max: {args.size}) ===")
    if args.dataset:
        print(f"Filter by dataset: {args.dataset}")

    result = search_mlit(args.term, args.size, args.dataset)
    total = result.get("totalNumber", 0)
    results = result.get("searchResults", [])

    print(f"Total matching items: {total}\n")

    if args.json_out:
        print(json.dumps(results, indent=2, ensure_ascii=False))
        return

    for idx, item in enumerate(results, 1):
        metadata = item.get("metadata", {})
        title = item.get("title", "No Title")
        item_id = item.get("id", "No ID")
        
        # 共通属性から代表的なデータを抽出
        prefecture = metadata.get("DPF:prefecture_name", "不明")
        dataset = metadata.get("DPF:dataset_id", "不明")
        lat = metadata.get("DPF:latitude", "N/A")
        lon = metadata.get("DPF:longitude", "N/A")
        completion = metadata.get("DPF:completion_datetime", "不明")

        print(f"[{idx}] {title}")
        print(f"  ID         : {item_id}")
        print(f"  Dataset ID : {dataset}")
        print(f"  Prefecture : {prefecture}")
        print(f"  Coordinates: Lat={lat}, Lon={lon}")
        print(f"  Completion : {completion}")
        print("-" * 50)

if __name__ == "__main__":
    main()
