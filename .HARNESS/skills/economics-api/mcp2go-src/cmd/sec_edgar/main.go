package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"economics-api/pkg/utils"
)

type TickerEntry struct {
	CIK    int    `json:"cik_str"`
	Ticker string `json:"ticker"`
	Title  string `json:"title"`
}

func getCIKFromTicker(ticker string, ua string) (string, error) {
	tickerURL := "https://www.sec.gov/files/company_tickers.json"
	headers := map[string]string{
		"User-Agent": ua,
	}

	body, err := utils.GetWithHeaders(tickerURL, headers)
	if err != nil {
		return "", err
	}

	var data map[string]TickerEntry
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	upperTicker := strings.ToUpper(ticker)
	for _, entry := range data {
		if strings.ToUpper(entry.Ticker) == upperTicker {
			return fmt.Sprintf("%010d", entry.CIK), nil
		}
	}
	return "", fmt.Errorf("ticker '%s' not found in SEC database", ticker)
}

func main() {
	ticker := flag.String("ticker", "", "Stock ticker symbol (e.g. AAPL, MSFT, GOOGL)")
	cikFlag := flag.String("cik", "", "10-digit CIK code (e.g. 0000320193)")
	fetchFacts := flag.Bool("facts", false, "Fetch company XBRL facts")
	fetchSubmissions := flag.Bool("submissions", false, "Fetch company filing submissions history")
	uaFlag := flag.String("user-agent", "AntigravityBot admin@antigravity-agent.org", "SEC required User-Agent (Format: Name Email)")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	_ = fetchFacts

	if *ticker == "" && *cikFlag == "" {
		fmt.Fprintln(os.Stderr, "Usage: sec_edgar --ticker AAPL [--facts|--submissions] OR sec_edgar --cik 0000320193")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cik := *cikFlag
	if cik == "" && *ticker != "" {
		resolvedCIK, err := getCIKFromTicker(*ticker, *uaFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
			os.Exit(1)
		}
		cik = resolvedCIK
	}

	// Pad CIK to 10 digits
	for len(cik) < 10 {
		cik = "0" + cik
	}

	headers := map[string]string{
		"User-Agent": *uaFlag,
	}

	var targetURL string
	if *fetchSubmissions {
		targetURL = fmt.Sprintf("https://data.sec.gov/submissions/CIK%s.json", cik)
	} else {
		// Default to facts
		targetURL = fmt.Sprintf("https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json", cik)
	}

	body, err := utils.GetWithHeaders(targetURL, headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query SEC EDGAR API: %v\n", err)
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
