package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// getAPIKey は環境変数またはWindowsレジストリからAPIキーを取得します
func getAPIKey() (string, error) {
	key := os.Getenv("KOKKOU_HUDOUSAN_API")
	if key != "" {
		return key, nil
	}

	// Windowsレジストリからの取得を試みる
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		val, _, err := k.GetStringValue("KOKKOU_HUDOUSAN_API")
		if err == nil && val != "" {
			return val, nil
		}
	}

	return "", fmt.Errorf("KOKKOU_HUDOUSAN_API 環境変数またはレジストリキーが設定されていません")
}

func main() {
	priceFlag := flag.Bool("price", false, "不動産価格（取引価格・成約価格）情報取得API (XIT001)")
	listFlag := flag.Bool("list", false, "都道府県内市区町村一覧取得API (XIT002)")

	areaFlag := flag.String("area", "", "都道府県コード (例: 13)")
	yearFlag := flag.String("year", "", "取引年 (例: 2025, priceの時に必須)")
	quarterFlag := flag.String("quarter", "", "四半期 (1-4)")
	cityFlag := flag.String("city", "", "市区町村コード (例: 13102)")
	stationFlag := flag.String("station", "", "駅コード")

	flag.Parse()

	// 1. APIキーの取得
	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 2. フラグのバリデーション
	if (*priceFlag && *listFlag) || (!*priceFlag && !*listFlag) {
		fmt.Fprintln(os.Stderr, "Error: -price または -list のいずれか一方のみを指定してください。")
		os.Exit(1)
	}

	var baseURL string
	params := url.Values{}

	if *priceFlag {
		baseURL = "https://www.reinfolib.mlit.go.jp/ex-api/external/XIT001"
		if *yearFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: -price 実行には -year が必須です。")
			os.Exit(1)
		}
		if *areaFlag == "" && *cityFlag == "" && *stationFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: -price 実行には -area, -city, -station のいずれか一つが必須です。")
			os.Exit(1)
		}

		params.Set("year", *yearFlag)
		if *quarterFlag != "" {
			params.Set("quarter", *quarterFlag)
		}
		if *areaFlag != "" {
			params.Set("area", *areaFlag)
		}
		if *cityFlag != "" {
			params.Set("city", *cityFlag)
		}
		if *stationFlag != "" {
			params.Set("station", *stationFlag)
		}
	} else {
		baseURL = "https://www.reinfolib.mlit.go.jp/ex-api/external/XIT002"
		if *areaFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: -list 実行には -area が必須です。")
			os.Exit(1)
		}
		params.Set("area", *areaFlag)
	}

	requestURL := baseURL + "?" + params.Encode()

	// 3. HTTPリクエストの作成と送信
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTP Error: %d %s\n", resp.StatusCode, resp.Status)
		// エラーレスポンスボディも読みだして出力する
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Response: %s\n", string(body))
		os.Exit(1)
	}

	// 4. gzipのデコードと出力
	var reader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating gzip reader: %v\n", err)
			os.Exit(1)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// 標準出力への書き出し
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying response output: %v\n", err)
		os.Exit(1)
	}
}
