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
	symbol := flag.String("symbol", "", "Stock/Index symbol (e.g. AAPL, TSLA, ^N225, 7203.T)")
	timeRange := flag.String("range", "1mo", "Chart range (1d, 5d, 1mo, 3mo, 6mo, 1y, 5y, max)")
	interval := flag.String("interval", "1d", "Chart interval (1m, 5m, 15m, 1d, 1wk, 1mo)")
	summary := flag.Bool("summary", false, "Fetch financial quote summary instead of chart")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	if *symbol == "" {
		fmt.Fprintln(os.Stderr, "Usage: yfinance --symbol AAPL [--range 1mo --interval 1d] OR yfinance --symbol 7203.T --summary")
		flag.PrintDefaults()
		os.Exit(1)
	}

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	var targetURL string
	if *summary {
		targetURL = fmt.Sprintf("https://query2.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=financialData,summaryDetail,defaultKeyStatistics",
			url.PathEscape(*symbol))
	} else {
		targetURL = fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=%s&interval=%s",
			url.PathEscape(*symbol), url.QueryEscape(*timeRange), url.QueryEscape(*interval))
	}

	body, err := utils.GetWithHeaders(targetURL, headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query Yahoo Finance API: %v\n", err)
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
