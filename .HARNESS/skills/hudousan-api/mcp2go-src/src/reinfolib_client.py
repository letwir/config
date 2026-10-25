import argparse
import gzip
import json
import os
import sys
import urllib.request
import urllib.parse

def main():
    # Windows環境で標準出力が文字化けするのを防ぐため、UTF-8に再設定しますの
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
    if hasattr(sys.stderr, 'reconfigure'):
        sys.stderr.reconfigure(encoding='utf-8')

    parser = argparse.ArgumentParser(description="不動産情報ライブラリ API PoC クライアント")
    parser.add_argument("--price", action="store_true", help="不動産価格（取引価格・成約価格）情報取得API (XIT001)")
    parser.add_argument("--list", action="store_true", help="都道府県内市区町村一覧取得API (XIT002)")
    
    parser.add_argument("--area", type=str, help="都道府県コード (例: 13)")
    parser.add_argument("--year", type=str, help="取引年 (例: 2025, XIT001必須)")
    parser.add_argument("--quarter", type=str, help="四半期 (1-4)")
    parser.add_argument("--city", type=str, help="市区町村コード (例: 13102)")
    
    args = parser.parse_args()
    
    # APIキーの取得
    api_key = os.environ.get("KOKKOU_HUDOUSAN_API")
    if not api_key:
        print("Error: 環境変数 KOKKOU_HUDOUSAN_API が設定されていません。", file=sys.stderr)
        sys.exit(1)
        
    if not args.price and not args.list:
        print("Error: --price または --list のいずれかを指定してください。", file=sys.stderr)
        sys.exit(1)
        
    if args.price and args.list:
        print("Error: --price と --list を同時に指定することはできません。", file=sys.stderr)
        sys.exit(1)
        
    params = {}
    if args.price:
        url = "https://www.reinfolib.mlit.go.jp/ex-api/external/XIT001"
        if not args.year:
            print("Error: --price 実行には --year が必須です。", file=sys.stderr)
            sys.exit(1)
        if not args.area and not args.city:
            print("Error: --price 実行には --area または --city のいずれかが必要です。", file=sys.stderr)
            sys.exit(1)
            
        params["year"] = args.year
        if args.quarter:
            params["quarter"] = args.quarter
        if args.area:
            params["area"] = args.area
        if args.city:
            params["city"] = args.city
            
    else:  # args.list
        url = "https://www.reinfolib.mlit.go.jp/ex-api/external/XIT002"
        if not args.area:
            print("Error: --list 実行には --area が必須です。", file=sys.stderr)
            sys.exit(1)
        params["area"] = args.area
        
    query_string = urllib.parse.urlencode(params)
    request_url = f"{url}?{query_string}"
    
    req = urllib.request.Request(request_url)
    req.add_header("Ocp-Apim-Subscription-Key", api_key)
    # 明示的にgzipを受け入れるヘッダーを追加
    req.add_header("Accept-Encoding", "gzip")
    
    try:
        with urllib.request.urlopen(req) as response:
            data = response.read()
            
            # gzip デコード
            content_encoding = response.headers.get("Content-Encoding", "").lower()
            if "gzip" in content_encoding or data.startswith(b'\x1f\x8b'):
                data = gzip.decompress(data)
                
            json_data = json.loads(data.decode("utf-8"))
            print(json.dumps(json_data, ensure_ascii=False, indent=2))
            
    except urllib.error.HTTPError as e:
        print(f"HTTP Error: {e.code} {e.reason}", file=sys.stderr)
        try:
            error_body = e.read()
            if error_body:
                # エラーレスポンスもgzip圧縮されている可能性がある
                if error_body.startswith(b'\x1f\x8b'):
                    error_body = gzip.decompress(error_body)
                print(error_body.decode("utf-8"), file=sys.stderr)
        except Exception:
            pass
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
