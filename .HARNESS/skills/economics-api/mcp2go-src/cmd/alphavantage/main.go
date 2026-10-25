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
	symbol := flag.String("symbol", "", "Stock/Crypto symbol (e.g. IBM, AAPL, BTC)")
	function := flag.String("function", "GLOBAL_QUOTE", "Function (GLOBAL_QUOTE, TIME_SERIES_DAILY, FX_DAILY, CURRENCY_EXCHANGE_RATE)")
	fromCurr := flag.String("from", "", "From currency code (for FX_DAILY or CURRENCY_EXCHANGE_RATE)")
	toCurr := flag.String("to", "", "To currency code (for FX_DAILY or CURRENCY_EXCHANGE_RATE)")
	keyFlag := flag.String("key", "", "Alpha Vantage API Key")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	if *keyFlag != "" {
		apiKey = *keyFlag
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] Alpha Vantage API key is required.")
		fmt.Fprintln(os.Stderr, "[INFO] Please set ALPHAVANTAGE_API_KEY environment variable or pass --key.")
		fmt.Fprintln(os.Stderr, "[INFO] See USAGE.md for registration instructions.")
		os.Exit(1)
	}

	var targetURL string
	if *fromCurr != "" && *toCurr != "" {
		targetURL = fmt.Sprintf("https://www.alphavantage.co/query?function=%s&from_currency=%s&to_currency=%s&apikey=%s",
			url.QueryEscape(*function), url.QueryEscape(*fromCurr), url.QueryEscape(*toCurr), url.QueryEscape(apiKey))
	} else if *symbol != "" {
		targetURL = fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s",
			url.QueryEscape(*function), url.QueryEscape(*symbol), url.QueryEscape(apiKey))
	} else {
		fmt.Fprintln(os.Stderr, "Usage: alphavantage --symbol IBM OR alphavantage --function CURRENCY_EXCHANGE_RATE --from USD --to JPY")
		flag.PrintDefaults()
		os.Exit(1)
	}

	body, err := utils.GetWithHeaders(targetURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query Alpha Vantage API: %v\n", err)
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
