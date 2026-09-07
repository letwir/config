package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

func main() {
	var query string
	var output string
	var format string
	var dbURI string

	flag.StringVar(&query, "query", "", "SQL query to execute")
	flag.StringVar(&query, "q", "", "SQL query to execute (shorthand)")
	flag.StringVar(&output, "output", "", "Output file path (default stdout)")
	flag.StringVar(&output, "o", "", "Output file path (shorthand)")
	flag.StringVar(&format, "format", "text", "Output format: m3u, json, csv, text")
	flag.StringVar(&format, "f", "text", "Output format (shorthand)")
	flag.StringVar(&dbURI, "db", "", "PostgreSQL Database Connection URI")
	flag.Parse()

	// If no query is provided, show usage
	if query == "" {
		fmt.Fprintln(os.Stderr, "Usage: fb2k-sql [options]")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Resolve database connection URI
	if dbURI == "" {
		dbURI = os.Getenv("DATABASE_URL")
	}
	if dbURI == "" {
		// Attempt to read from config.toml in current or parent directories
		dbURI = readDBURIFromConfig("config.toml")
	}
	if dbURI == "" {
		dbURI = readDBURIFromConfig("../config.toml")
	}
	if dbURI == "" {
		dbURI = readDBURIFromConfig("../../config.toml")
	}
	if dbURI == "" {
		dbURI = readDBURIFromConfig("../../../config.toml")
	}
	if dbURI == "" {
		fmt.Fprintln(os.Stderr, "Error: PostgreSQL Database URI not specified.")
		fmt.Fprintln(os.Stderr, "Please provide it via -db argument, DATABASE_URL environment variable, or config.toml.")
		os.Exit(1)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	colNames := make([]string, len(fields))
	for i, f := range fields {
		colNames[i] = f.Name
	}

	// Prepare output writer
	var out io.Writer = os.Stdout
	if output != "" {
		file, err := os.Create(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		out = file
	}

	// Fetch all rows
	var rawData []map[string]any
	var rawRows [][]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read row values: %v\n", err)
			os.Exit(1)
		}
		rowMap := make(map[string]any)
		for i, colName := range colNames {
			rowMap[colName] = values[i]
		}
		rawData = append(rawData, rowMap)
		rawRows = append(rawRows, values)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Rows iteration error: %v\n", err)
		os.Exit(1)
	}

	// Format and write output
	switch strings.ToLower(format) {
	case "m3u":
		writeM3U(out, colNames, rawData)
	case "json":
		writeJSON(out, rawData)
	case "csv":
		writeCSV(out, colNames, rawRows)
	case "text":
		writeText(out, colNames, rawRows)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s. Defaulting to text.\n", format)
		writeText(out, colNames, rawRows)
	}
}

// readDBURIFromConfig parses config.toml to find target database URL
func readDBURIFromConfig(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	// Simple regex parsing: url = "..."
	re := regexp.MustCompile(`(?m)^\s*url\s*=\s*["']([^"']+)["']`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// writeM3U extracts filepath (and metadata if present) and writes as M3U playlist
func writeM3U(w io.Writer, cols []string, data []map[string]any) {
	fmt.Fprintln(w, "#EXTM3U")

	// Find columns for filepath, artist, title
	filepathCol := ""
	artistCol := ""
	titleCol := ""
	for _, c := range cols {
		lower := strings.ToLower(c)
		if lower == "filepath" || lower == "path" || lower == "uri" {
			filepathCol = c
		} else if lower == "artist" || lower == "album_artist" {
			artistCol = c
		} else if lower == "title" {
			titleCol = c
		}
	}

	// If filepathCol not found, default to first column
	if filepathCol == "" && len(cols) > 0 {
		filepathCol = cols[0]
	}

	for _, row := range data {
		pathVal, exists := row[filepathCol]
		if !exists || pathVal == nil {
			continue
		}
		pathStr := fmt.Sprintf("%v", pathVal)

		artistStr := ""
		if artistCol != "" {
			if v, ok := row[artistCol]; ok && v != nil {
				artistStr = fmt.Sprintf("%v", v)
			}
		}
		titleStr := ""
		if titleCol != "" {
			if v, ok := row[titleCol]; ok && v != nil {
				titleStr = fmt.Sprintf("%v", v)
			}
		}

		if artistStr != "" || titleStr != "" {
			fmt.Fprintf(w, "#EXTINF:-1,%s - %s\n", artistStr, titleStr)
		} else {
			fmt.Fprintln(w, "#EXTINF:-1,")
		}
		fmt.Fprintln(w, pathStr)
	}
}

// writeJSON writes query results in JSON format
func writeJSON(w io.Writer, data []map[string]any) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "JSON encode error: %v\n", err)
	}
}

// writeCSV writes query results in CSV format
func writeCSV(w io.Writer, cols []string, rows [][]any) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	if err := writer.Write(cols); err != nil {
		fmt.Fprintf(os.Stderr, "CSV write header error: %v\n", err)
		return
	}

	// Write rows
	for _, r := range rows {
		strRow := make([]string, len(r))
		for i, val := range r {
			if val == nil {
				strRow[i] = ""
			} else {
				strRow[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(strRow); err != nil {
			fmt.Fprintf(os.Stderr, "CSV write row error: %v\n", err)
			return
		}
	}
}

// writeText writes query results in a clean, tab-separated format
func writeText(w io.Writer, cols []string, rows [][]any) {
	// Write header
	fmt.Fprintln(w, strings.Join(cols, "\t"))

	// Write rows
	for _, r := range rows {
		strRow := make([]string, len(r))
		for i, val := range r {
			if val == nil {
				strRow[i] = "NULL"
			} else {
				strRow[i] = fmt.Sprintf("%v", val)
			}
		}
		fmt.Fprintln(w, strings.Join(strRow, "\t"))
	}
}
