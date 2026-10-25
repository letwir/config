package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"economics-api/pkg/utils"
)

func main() {
	statsDataID := flag.String("stats-data", "", "Stats data ID (getStatsData)")
	searchWord := flag.String("search", "", "Search word for stats list (getStatsList)")
	apiKeyFlag := flag.String("api-key", "", "e-Stat API Key (appId)")
	limit := flag.Int("limit", 10, "Limit for search results")
	raw := flag.Bool("raw", false, "Output raw JSON")
	flag.Parse()

	apiKey := os.Getenv("ESTAT_API_KEY")
	if *apiKeyFlag != "" {
		apiKey = *apiKeyFlag
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] e-Stat API key (appId) is required.")
		fmt.Fprintln(os.Stderr, "[INFO] Please set ESTAT_API_KEY environment variable or pass --api-key.")
		fmt.Fprintln(os.Stderr, "[INFO] See USAGE.md for registration instructions.")
		os.Exit(1)
	}

	var reqURL string
	if *statsDataID != "" {
		reqURL = fmt.Sprintf("https://api.e-stat.go.jp/rest/3.0/app/json/getStatsData?appId=%s&statsDataId=%s",
			url.QueryEscape(apiKey), url.QueryEscape(*statsDataID))
	} else if *searchWord != "" {
		reqURL = fmt.Sprintf("https://api.e-stat.go.jp/rest/3.0/app/json/getStatsList?appId=%s&searchWord=%s&limit=%d",
			url.QueryEscape(apiKey), url.QueryEscape(*searchWord), *limit)
	} else {
		fmt.Fprintln(os.Stderr, "Usage: estat --stats-data <statsDataId> OR estat --search <searchWord>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query e-Stat API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to read response body: %v\n", err)
		os.Exit(1)
	}

	if *raw {
		fmt.Println(string(body))
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(string(body))
		return
	}

	formatted, err := utils.FormatJSON(data)
	if err != nil {
		fmt.Println(string(body))
		return
	}
	fmt.Println(formatted)
}
