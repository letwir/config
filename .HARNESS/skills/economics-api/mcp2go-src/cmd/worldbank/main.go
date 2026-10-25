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
	country := flag.String("country", "JPN", "Country code (2 or 3 letters, e.g. JPN, USA, CHN, WLD)")
	indicator := flag.String("indicator", "NY.GDP.MKTP.CD", "Indicator ID (e.g. NY.GDP.MKTP.CD for GDP, FP.CPI.TOTL.ZG for inflation, SP.POP.TOTL for population)")
	dateRange := flag.String("date", "", "Date range (e.g. 2010:2025)")
	raw := flag.Bool("raw", false, "Output raw response")
	flag.Parse()

	targetURL := fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/%s?format=json",
		url.PathEscape(*country), url.PathEscape(*indicator))

	if *dateRange != "" {
		targetURL += "&date=" + url.QueryEscape(*dateRange)
	}

	body, err := utils.GetWithHeaders(targetURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to query World Bank API: %v\n", err)
		os.Exit(1)
	}

	if *raw {
		fmt.Println(string(body))
		return
	}

	var data interface{}
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
