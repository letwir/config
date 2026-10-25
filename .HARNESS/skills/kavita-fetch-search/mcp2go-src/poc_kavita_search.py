import os
import sys
import re
import urllib.request
import urllib.parse
import json
import zipfile
import io
from html.parser import HTMLParser

# HTMLタグを剥がすための簡易パーサー
class MLStripper(HTMLParser):
    def __init__(self):
        super().__init__()
        self.reset()
        self.strict = False
        self.convert_charrefs = True
        self.text = io.StringIO()
    def handle_data(self, d):
        self.text.write(d)
    def get_data(self):
        return self.text.getvalue()

def strip_tags(html):
    s = MLStripper()
    s.feed(html)
    return s.get_data()

def load_config():
    config = {}
    config_path = "config.toml"
    if not os.path.exists(config_path):
        print(f"Error: {config_path} not found.")
        sys.exit(1)
    with open(config_path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            if "=" in line:
                k, v = line.split("=", 1)
                config[k.strip()] = v.strip().strip('"').strip("'")
    return config

def make_request(url, token=None, method="GET", data=None):
    req = urllib.request.Request(url, method=method)
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    if data:
        req.add_header("Content-Type", "application/json")
        jsondata = json.dumps(data).encode("utf-8")
        req.data = jsondata
    
    try:
        with urllib.request.urlopen(req) as res:
            return res.read()
    except Exception as e:
        print(f"HTTP Request failed to {url}: {e}")
        if hasattr(e, "read"):
            print("Response body:", e.read().decode("utf-8", errors="ignore"))
        return None

def authenticate(server_url, api_key):
    auth_url = f"{server_url.rstrip('/')}/api/Plugin/authenticate?apiKey={api_key}&pluginName=AntigravitySearchPoC"
    # POST request (empty body)
    res_bytes = make_request(auth_url, method="POST")
    if not res_bytes:
        print("Authentication failed.")
        sys.exit(1)
    
    user_info = json.loads(res_bytes.decode("utf-8"))
    token = user_info.get("token")
    if not token:
        print("Failed to retrieve token from auth response.")
        sys.exit(1)
    return token

def search_series(server_url, token, query):
    search_url = f"{server_url.rstrip('/')}/api/Search/search?queryString={urllib.parse.quote(query)}"
    res_bytes = make_request(search_url, token=token)
    if not res_bytes:
        print("Search request failed.")
        return []
    
    search_result = json.loads(res_bytes.decode("utf-8"))
    return search_result.get("series", [])

def get_chapters(server_url, token, series_id):
    chapters_url = f"{server_url.rstrip('/')}/api/Search/chapters-by-series?seriesId={series_id}"
    res_bytes = make_request(chapters_url, token=token)
    if not res_bytes:
        print("Failed to get chapters.")
        return []
    return json.loads(res_bytes.decode("utf-8"))

def download_chapter(server_url, token, chapter_id):
    dl_url = f"{server_url.rstrip('/')}/api/Download/chapter?chapterId={chapter_id}"
    res_bytes = make_request(dl_url, token=token)
    return res_bytes

def parse_epub_and_search(epub_bytes, search_query):
    matches = []
    try:
        # ZIPとしてパース
        with zipfile.ZipFile(io.BytesIO(epub_bytes)) as z:
            for file_info in z.infolist():
                filename = file_info.filename.lower()
                # HTML / XHTML のテキストを探す
                if filename.endswith(".html") or filename.endswith(".xhtml") or filename.endswith(".htm"):
                    with z.open(file_info) as f:
                        content = f.read().decode("utf-8", errors="ignore")
                        # HTMLタグを剥がす
                        plain_text = strip_tags(content)
                        # 行ごとに分割して検索
                        lines = plain_text.splitlines()
                        for i, line in enumerate(lines):
                            if search_query.lower() in line.lower():
                                # 前後のコンテキストを取得
                                prev_line = lines[i-1].strip() if i > 0 else ""
                                next_line = lines[i+1].strip() if i < len(lines)-1 else ""
                                matches.append({
                                    "file": file_info.filename,
                                    "line_num": i + 1,
                                    "match_line": line.strip(),
                                    "prev_line": prev_line,
                                    "next_line": next_line
                                })
    except Exception as e:
        print(f"Error parsing EPUB: {e}")
    return matches

def main():
    # Windowsのコンソールでのエンコーディングエラー対策
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
        
    if len(sys.argv) < 3:
        print("Usage: python poc_kavita_search.py <series_name> <search_query>")
        sys.exit(1)
        
    series_name = sys.argv[1]
    search_query = sys.argv[2]
    
    print("Loading config...")
    config = load_config()

    server_url = config.get("server_url")
    api_key = config.get("api_key")
    
    if not server_url or not api_key:
        print("Error: server_url or api_key missing in config.toml.")
        sys.exit(1)
        
    print(f"Connecting to Kavita at {server_url}...")
    token = authenticate(server_url, api_key)
    print("Authenticated successfully!")
    
    print(f"Searching for series matching '{series_name}'...")
    series_list = search_series(server_url, token, series_name)
    if not series_list:
        print("No series found matching the query.")
        sys.exit(1)
    
    print(f"Found {len(series_list)} series:")
    for idx, series in enumerate(series_list):
        print(f"  [{idx}] ID: {series['seriesId']} - Name: {series['name']}")
    
    # 最初のシリーズを使用
    target_series = series_list[0]
    print(f"\nUsing series: {target_series['name']} (ID: {target_series['seriesId']})")

    
    print("Fetching chapters...")
    chapters = get_chapters(server_url, token, target_series['seriesId'])
    print(f"Found {len(chapters)} chapters.")

    
    for chapter in chapters:
        chapter_id = chapter['id']
        chapter_title = chapter.get('title', f"Chapter {chapter.get('number', '?')}")
        print(f"\nDownloading chapter '{chapter_title}' (ID: {chapter_id})...")
        epub_bytes = download_chapter(server_url, token, chapter_id)
        if not epub_bytes:
            print("Failed to download chapter.")
            continue
        
        print("Parsing and searching content...")
        matches = parse_epub_and_search(epub_bytes, search_query)
        if matches:
            print(f"Found {len(matches)} matches in '{chapter_title}':")
            for m in matches:
                print(f"  --- File: {m['file']} (Line {m['line_num']}) ---")
                if m['prev_line']:
                    print(f"    - {m['prev_line']}")
                print(f"    * \033[93m{m['match_line']}\033[0m")
                if m['next_line']:
                    print(f"    - {m['next_line']}")
        else:
            print("No matches found in this chapter.")

if __name__ == "__main__":
    main()
