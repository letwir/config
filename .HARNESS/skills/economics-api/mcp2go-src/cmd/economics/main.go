package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func printUsage() {
	fmt.Println(`=======================================================
 Economics & Finance Unified CLI (economics)
=======================================================
Usage:
  economics <subcommand> [flags...]

Available Subcommands:
  estat        Japan Official Statistics API (e-Stat)
  edinet       Japan Financial Disclosure API (EDINET v2)
  sec          US SEC EDGAR Company Financials API
  fred         Federal Reserve Economic Data API (FRED)
  worldbank    World Bank Global Indicators API
  yfinance     Yahoo Finance Realtime & Historical Market Data
  alphavantage Alpha Vantage Global Market & Forex Data

Global Options:
  --help, -h   Show help message

Example Usage:
  economics estat --stats-data 0003444458
  economics edinet --date 2026-08-14
  economics sec --ticker AAPL --facts
  economics fred --series FEDFUNDS
  economics worldbank --country JPN --indicator NY.GDP.MKTP.CD
  economics yfinance --symbol AAPL --range 1mo
  economics alphavantage --symbol 7203.T
=======================================================`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	subCmd := os.Args[1]
	if subCmd == "-h" || subCmd == "--help" || subCmd == "help" {
		printUsage()
		os.Exit(0)
	}

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	var targetBin string
	switch subCmd {
	case "estat":
		targetBin = "estat" + ext
	case "edinet":
		targetBin = "edinet" + ext
	case "sec", "sec_edgar":
		targetBin = "sec_edgar" + ext
	case "fred":
		targetBin = "fred" + ext
	case "worldbank", "wb":
		targetBin = "worldbank" + ext
	case "yfinance", "yf":
		targetBin = "yfinance" + ext
	case "alphavantage", "av":
		targetBin = "alphavantage" + ext
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] Unknown subcommand: '%s'\n", subCmd)
		printUsage()
		os.Exit(1)
	}

	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		execDir = "."
	}

	binPath := filepath.Join(execDir, targetBin)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		// Fallback to searching PATH if not found in same directory
		binPath = targetBin
	}

	subArgs := os.Args[2:]
	cmd := exec.Command(binPath, subArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to execute %s: %v\n", targetBin, err)
		os.Exit(1)
	}
}
