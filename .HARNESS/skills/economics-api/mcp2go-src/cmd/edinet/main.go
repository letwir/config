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
	dateStr := flag.String("date", "", "Target date for document list (YYYY-MM-DD)")
	docID := flag.String("doc-id", "", "Document ID to download")
	docType := flag.Int("type", 1, "Document type to download (1:XBRL/PDF, 2:PDF, 5:CSV)")
	keyFlag := flag.String("key", "", "EDINET Subscription-Key")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	subKey := os.Getenv("EDINET_SUBSCRIPTION_KEY")
	if *keyFlag != "" {
		subKey = *keyFlag
	}

	if subKey == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] EDINET Subscription-Key is required.")
		fmt.Fprintln(os.Stderr, "[INFO] Please set EDINET_SUBSCRIPTION_KEY environment variable or pass --key.")
		fmt.Fprintln(os.Stderr, "[INFO] See USAGE.md for registration instructions.")
		os.Exit(1)
	}

	headers := map[string]string{
		"Subscription-Key": subKey,
	}

	var targetURL string
	if *docID != "" {
		targetURL = fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents/%s?type=%d",
			url.PathEscape(*docID), *docType)
	} else if *dateStr != "" {
		targetURL = fmt.Sprintf("https://api.edinet-fsa.go.jp/api/v2/documents.json?date=%s&type=2",
			url.QueryEscape(*dateStr))
	} else {
		fmt.Fprintln(os.Stderr, "Usage: edinet --date YYYY-MM-DD OR edinet --doc-id <docID>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	body, err := utils.GetWithHeaders(targetURL, headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query EDINET API: %v\n", err)
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
