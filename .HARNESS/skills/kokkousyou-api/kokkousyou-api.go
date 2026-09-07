package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// SearchResponse は API レスポンスをマッピングするための構造体です
type SearchResponse struct {
	Data struct {
		Search struct {
			TotalNumber   int `json:"totalNumber"`
			SearchResults []struct {
				ID       string                 `json:"id"`
				Title    string                 `json:"title"`
				Metadata map[string]interface{} `json:"metadata"`
			} `json:"searchResults"`
		} `json:"search"`
	} `json:"data"`
}

// getEnvFromRegistry は環境変数が設定されていない場合に Windows レジストリから取得を試みます
func getEnvFromRegistry(name string) string {
	// ユーザー環境変数 (HKCU)
	k, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err == nil {
		val, _, err := k.GetStringValue(name)
		k.Close()
		if err == nil && val != "" {
			return val
		}
	}
	// システム環境変数 (HKLM)
	k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	if err == nil {
		val, _, err := k.GetStringValue(name)
		k.Close()
		if err == nil && val != "" {
			return val
		}
	}
	return ""
}

// getAPIKey は環境変数またはレジストリから KOKKOU_API キーを取得します
func getAPIKey() string {
	apiKey := os.Getenv("KOKKOU_API")
	if apiKey == "" {
		apiKey = getEnvFromRegistry("KOKKOU_API")
	}
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: KOKKOU_API environment variable is not set in Windows Environment or Registry.")
		os.Exit(1)
	}
	return apiKey
}

// parseAround は "lat,lon,distance" 形式の文字列をパースします
func parseAround(aroundStr string) (float64, float64, float64, error) {
	parts := strings.Split(aroundStr, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid around format: expected lat,lon,distance")
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid latitude: %v", err)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid longitude: %v", err)
	}
	dist, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid distance: %v", err)
	}
	return lat, lon, dist, nil
}

func main() {
	var (
		isCount          bool
		around           string
		targetDataset    string
		targetPrefecture string
		isJson           bool
		limit            int = 10
		term             string
	)

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-c" || arg == "--count" {
			isCount = true
		} else if arg == "-j" || arg == "--json" {
			isJson = true
		} else if (arg == "-d" || arg == "--dataset") && i+1 < len(args) {
			targetDataset = args[i+1]
			i++
		} else if (arg == "-p" || arg == "--prefecture") && i+1 < len(args) {
			targetPrefecture = args[i+1]
			i++
		} else if arg == "--around" && i+1 < len(args) {
			around = args[i+1]
			i++
		} else if arg == "--limit" && i+1 < len(args) {
			val, err := strconv.Atoi(args[i+1])
			if err == nil {
				limit = val
			}
			i++
		} else if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		} else {
			if term == "" {
				term = arg
			}
		}
	}

	// APIキーの取得
	apiKey := getAPIKey()

	// GraphQL クエリパラメータの動的構築
	var queryParams []string
	queryParams = append(queryParams, "phraseMatch: true")
	queryParams = append(queryParams, "first: 0")
	queryParams = append(queryParams, fmt.Sprintf("size: %d", limit))

	if term != "" {
		// ダブルクォート等をエスケープ
		escapedTerm := strings.ReplaceAll(term, `"`, `\"`)
		queryParams = append(queryParams, fmt.Sprintf(`term: "%s"`, escapedTerm))
	}

	// 空間検索フィルタの構築
	if around != "" {
		lat, lon, dist, err := parseAround(around)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --around: %v\n", err)
			os.Exit(1)
		}
		queryParams = append(queryParams, fmt.Sprintf(`locationFilter: { geoDistance: { lat: %f, lon: %f, distance: %f } }`, lat, lon, dist))
	}

	// メタデータフィルタの構築 (dataset, prefecture)
	if targetDataset != "" && targetPrefecture != "" {
		queryParams = append(queryParams, fmt.Sprintf(`attributeFilter: { AND: [ { attributeName: "DPF:dataset_id", is: "%s" }, { attributeName: "DPF:prefecture_name", is: "%s" } ] }`, targetDataset, targetPrefecture))
	} else if targetDataset != "" {
		queryParams = append(queryParams, fmt.Sprintf(`attributeFilter: { attributeName: "DPF:dataset_id", is: "%s" }`, targetDataset))
	} else if targetPrefecture != "" {
		queryParams = append(queryParams, fmt.Sprintf(`attributeFilter: { attributeName: "DPF:prefecture_name", is: "%s" }`, targetPrefecture))
	}

	// GraphQLクエリ組み立て
	query := fmt.Sprintf(`
	query {
		search(%s) {
			totalNumber
			searchResults {
				id
				title
				metadata
			}
		}
	}
	`, strings.Join(queryParams, ", "))

	// ペイロード構築
	payloadMap := map[string]string{"query": query}
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Payload Marshal error: %v\n", err)
		os.Exit(1)
	}

	// HTTPリクエスト送信
	url := "https://data-platform.mlit.go.jp/api/v1/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create HTTP request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API HTTP Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "API returned non-200 status: %d. Response: %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON response: %v\n", err)
		os.Exit(1)
	}

	searchResult := searchResp.Data.Search

	// --- 出力フェーズ ---

	// 件数のみ取得モードの場合
	if isCount {
		fmt.Println(searchResult.TotalNumber)
		return
	}

	// JSON出力モードの場合
	if isJson {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false) // 日本語の文字化け（エスケープ）防止
		enc.SetIndent("", "  ")
		if err := enc.Encode(searchResult.SearchResults); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format JSON output: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(buf.String())
		return
	}

	// テキスト出力モードの場合
	fmt.Printf("=== MLIT DPF Search: '%s' (Limit: %d) ===\n", term, limit)
	if targetDataset != "" {
		fmt.Printf("Filter by dataset   : %s\n", targetDataset)
	}
	if targetPrefecture != "" {
		fmt.Printf("Filter by prefecture: %s\n", targetPrefecture)
	}
	if around != "" {
		fmt.Printf("Filter by around coordinates: %s\n", around)
	}
	fmt.Printf("Total matching items: %d\n\n", searchResult.TotalNumber)

	for idx, item := range searchResult.SearchResults {
		prefectureName := "不明"
		if pref, ok := item.Metadata["DPF:prefecture_name"]; ok && pref != nil {
			prefectureName = fmt.Sprintf("%v", pref)
		}
		datasetID := "不明"
		if ds, ok := item.Metadata["DPF:dataset_id"]; ok && ds != nil {
			datasetID = fmt.Sprintf("%v", ds)
		}
		latitude := "N/A"
		if lat, ok := item.Metadata["DPF:latitude"]; ok && lat != nil {
			latitude = fmt.Sprintf("%v", lat)
		}
		longitude := "N/A"
		if lon, ok := item.Metadata["DPF:longitude"]; ok && lon != nil {
			longitude = fmt.Sprintf("%v", lon)
		}
		completion := "不明"
		if comp, ok := item.Metadata["DPF:completion_datetime"]; ok && comp != nil {
			completion = fmt.Sprintf("%v", comp)
		}

		fmt.Printf("[%d] %s\n", idx+1, item.Title)
		fmt.Printf("  ID         : %s\n", item.ID)
		fmt.Printf("  Dataset ID : %s\n", datasetID)
		fmt.Printf("  Prefecture : %s\n", prefectureName)
		fmt.Printf("  Coordinates: Lat=%s, Lon=%s\n", latitude, longitude)
		fmt.Printf("  Completion : %s\n", completion)
		fmt.Println("--------------------------------------------------")
	}
}
