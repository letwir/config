package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"economics-api/pkg/utils"
)

func main() {
	seriesID := flag.String("series", "", "FRED series ID (e.g. FEDFUNDS, CPIAUCSL, GDP, DEXJPUS)")
	searchWord := flag.String("search", "", "Search FRED series by keyword")
	keyFlag := flag.String("key", "", "FRED API Key")
	limit := flag.Int("limit", 10, "Limit for search results")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	apiKey := os.Getenv("FRED_API_KEY")
	if *keyFlag != "" {
		apiKey = *keyFlag
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] FRED API key is required.")
		fmt.Fprintln(os.Stderr, "[INFO] Please set FRED_API_KEY environment variable or pass --key.")
		fmt.Fprintln(os.Stderr, "[INFO] See USAGE.md for registration instructions.")
		os.Exit(1)
	}

	var targetURL string
	if *seriesID != "" {
		targetURL = fmt.Sprintf("https://api.stlouisfed.org/fred/series/observations?series_id=%s&api_key=%s&file_type=json",
			url.QueryEscape(*seriesID), url.QueryEscape(apiKey))
	} else if *searchWord != "" {
		targetURL = fmt.Sprintf("https://api.stlouisfed.org/fred/series/search?search_text=%s&api_key=%s&file_type=json&limit=%d",
			url.QueryEscape(*searchWord), url.QueryEscape(apiKey), *limit)
	} else {
		fmt.Fprintln(os.Stderr, "Usage: fred --series FEDFUNDS OR fred --search GDP")
		flag.PrintDefaults()
		os.Exit(1)
	}

	body, err := utils.GetWithHeaders(targetURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query FRED API: %v\n", err)
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
