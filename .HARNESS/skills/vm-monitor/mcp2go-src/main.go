package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	var (
		rawURL     string
		flagStatus bool
		flagMetrics bool
		searchWord  string
		queryStr    string
		rangeQuery  string
		startTime   string
		endTime     string
		stepTime    string
		flagNodes   bool
		flagSuggest bool
		format      string
	)

	defaultURL := os.Getenv("VM_ADDR")
	if defaultURL == "" {
		defaultURL = "http://docker.tigris-tailor.ts.net:8428"
	}

	flag.StringVar(&rawURL, "url", defaultURL, "VictoriaMetrics Base URL (env: VM_ADDR)")
	flag.StringVar(&rawURL, "u", defaultURL, "VictoriaMetrics Base URL (shorthand)")
	flag.BoolVar(&flagStatus, "status", false, "Check TSDB status and health")
	flag.BoolVar(&flagMetrics, "metrics", false, "List registered metric names")
	flag.StringVar(&searchWord, "search", "", "Search metric names by keyword")
	flag.StringVar(&queryStr, "query", "", "Execute instant MetricsQL/PromQL query")
	flag.StringVar(&queryStr, "q", "", "Execute instant MetricsQL/PromQL query (shorthand)")
	flag.StringVar(&rangeQuery, "query-range", "", "Execute range MetricsQL/PromQL query")
	flag.StringVar(&startTime, "start", "", "Start time for range query (RFC3339 or relative)")
	flag.StringVar(&endTime, "end", "", "End time for range query (RFC3339 or relative)")
	flag.StringVar(&stepTime, "step", "1m", "Step resolution for range query")
	flag.BoolVar(&flagNodes, "nodes", false, "Show cluster node resource statistics")
	flag.BoolVar(&flagNodes, "top", false, "Show cluster node resource statistics (alias for -nodes)")
	flag.BoolVar(&flagSuggest, "suggest", false, "Analyze spare compute resources and suggest deployment targets")
	flag.BoolVar(&flagSuggest, "capacity", false, "Analyze spare compute resources (alias for -suggest)")
	flag.StringVar(&format, "format", "table", "Output format: table, json")
	flag.StringVar(&format, "f", "table", "Output format (shorthand)")
	flag.Parse()

	// If no action flags are specified, default to -suggest
	if !flagStatus && !flagMetrics && searchWord == "" && queryStr == "" && rangeQuery == "" && !flagNodes && !flagSuggest {
		flagSuggest = true
	}

	client, err := NewVMClient(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing VictoriaMetrics client: %v\n", err)
		os.Exit(1)
	}

	// 1. Status check
	if flagStatus {
		handleStatus(client, format)
		return
	}

	// 2. Metrics listing / search
	if flagMetrics || searchWord != "" {
		handleMetrics(client, searchWord, format)
		return
	}

	// 3. Instant query
	if queryStr != "" {
		handleQuery(client, queryStr, format)
		return
	}

	// 4. Range query
	if rangeQuery != "" {
		handleQueryRange(client, rangeQuery, startTime, endTime, stepTime, format)
		return
	}

	// 5. Node resources table
	if flagNodes {
		handleNodes(client, format)
		return
	}

	// 6. Spare resource analysis & deployment suggestion
	if flagSuggest {
		handleSuggest(client, format)
		return
	}
}

func handleStatus(client *VMClient, format string) {
	_, err := client.HealthCheck()
	healthStr := "HEALTHY"
	if err != nil {
		healthStr = fmt.Sprintf("UNHEALTHY (%v)", err)
	}

	status, err := client.GetTSDBStatus()
	if err != nil {
		if format == "json" {
			out := map[string]any{
				"endpoint": client.BaseURL,
				"health":   healthStr,
				"error":    err.Error(),
			}
			json.NewEncoder(os.Stdout).Encode(out)
			return
		}
		fmt.Fprintf(os.Stderr, "TSDB Status Error: %v\n", err)
		fmt.Printf("Endpoint: %s\nHealth:   %s\n", client.BaseURL, healthStr)
		return
	}

	if format == "json" {
		out := map[string]any{
			"endpoint":    client.BaseURL,
			"health":      healthStr,
			"totalSeries": status.Data.TotalSeries,
			"status":      status.Data,
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	fmt.Println("==================================================")
	fmt.Printf(" VictoriaMetrics TSDB Status\n")
	fmt.Println("==================================================")
	fmt.Printf("Endpoint:     %s\n", client.BaseURL)
	fmt.Printf("Health:       %s\n", healthStr)
	fmt.Printf("Total Series: %d\n", status.Data.TotalSeries)
	fmt.Printf("Total Labels: %d\n", status.Data.TotalLabelValuePairs)

	if len(status.Data.SeriesCountByMetricName) > 0 {
		fmt.Println("\nTop Metric Names by Series Count:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Metric Name\tSeries Count")
		fmt.Fprintln(w, "-----------\t------------")
		limit := 10
		if len(status.Data.SeriesCountByMetricName) < limit {
			limit = len(status.Data.SeriesCountByMetricName)
		}
		for i := 0; i < limit; i++ {
			item := status.Data.SeriesCountByMetricName[i]
			fmt.Fprintf(w, "%s\t%d\n", item.Name, item.Value)
		}
		w.Flush()
	}
}

func handleMetrics(client *VMClient, keyword string, format string) {
	names, err := client.GetMetricNames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve metric names: %v\n", err)
		os.Exit(1)
	}

	var filtered []string
	if keyword != "" {
		kwLower := strings.ToLower(keyword)
		for _, name := range names {
			if strings.Contains(strings.ToLower(name), kwLower) {
				filtered = append(filtered, name)
			}
		}
	} else {
		filtered = names
	}

	if format == "json" {
		out := map[string]any{
			"total":   len(names),
			"matched": len(filtered),
			"metrics": filtered,
		}
		json.NewEncoder(os.Stdout).Encode(out)
		return
	}

	fmt.Printf("Total Metrics: %d, Matched: %d\n", len(names), len(filtered))
	for _, name := range filtered {
		fmt.Println(name)
	}
}

func handleQuery(client *VMClient, query string, format string) {
	resp, err := client.QueryInstant(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Error: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(resp)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Metric\tValue")
	fmt.Fprintln(w, "------\t-----")
	for _, r := range resp.Data.Result {
		var parts []string
		for k, v := range r.Metric {
			parts = append(parts, fmt.Sprintf("%s=%q", k, v))
		}
		metricStr := strings.Join(parts, ", ")
		valStr := ""
		if len(r.Value) >= 2 {
			valStr = fmt.Sprintf("%v", r.Value[1])
		}
		fmt.Fprintf(w, "{%s}\t%s\n", metricStr, valStr)
	}
	w.Flush()
}

func handleQueryRange(client *VMClient, query, start, end, step, format string) {
	resp, err := client.QueryRange(query, start, end, step)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Range Error: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(resp)
		return
	}

	for _, r := range resp.Data.Result {
		var parts []string
		for k, v := range r.Metric {
			parts = append(parts, fmt.Sprintf("%s=%q", k, v))
		}
		fmt.Printf("--- Metric: {%s} ---\n", strings.Join(parts, ", "))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Timestamp\tValue")
		for _, val := range r.Values {
			if len(val) >= 2 {
				ts := fmt.Sprintf("%v", val[0])
				v := fmt.Sprintf("%v", val[1])
				fmt.Fprintf(w, "%s\t%s\n", ts, v)
			}
		}
		w.Flush()
	}
}

func handleNodes(client *VMClient, format string) {
	nodes, err := CollectClusterResources(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to collect cluster resources: %v\n", err)
		os.Exit(1)
	}

	if format == "json" {
		json.NewEncoder(os.Stdout).Encode(nodes)
		return
	}

	if len(nodes) == 0 {
		fmt.Println("No node metrics found. Verify if node_exporter is sending data to VictoriaMetrics.")
		return
	}

	fmt.Println("==========================================================================================")
	fmt.Printf(" Cluster Nodes Resource Status (%s)\n", client.BaseURL)
	fmt.Println("==========================================================================================")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Instance\tHostname\tCPU Used (%)\tCPU Idle (%)\tRAM Free\tRAM Total\tDisk Free\tLoad(1m)")
	fmt.Fprintln(w, "--------\t--------\t------------\t------------\t--------\t---------\t---------\t--------")

	for _, n := range nodes {
		freeRAMGB := float64(n.MemAvailBytes) / (1024 * 1024 * 1024)
		totRAMGB := float64(n.MemTotalBytes) / (1024 * 1024 * 1024)
		freeDiskGB := float64(n.DiskAvailBytes) / (1024 * 1024 * 1024)

		fmt.Fprintf(w, "%s\t%s\t%.1f%%\t%.1f%%\t%.1f GB\t%.1f GB\t%.1f GB\t%.2f\n",
			n.Instance, n.Hostname, n.CPUUsagePct, n.CPUIdlePct, freeRAMGB, totRAMGB, freeDiskGB, n.Load1)
	}
	w.Flush()
}

func handleSuggest(client *VMClient, format string) {
	nodes, err := CollectClusterResources(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to collect cluster resources: %v\n", err)
		os.Exit(1)
	}

	suggestions := GenerateDeploymentSuggestions(nodes)

	if format == "json" {
		summary := ClusterSummary{
			Endpoint:    client.BaseURL,
			CheckedAt:   time.Now(),
			TotalNodes:  len(nodes),
			Nodes:       nodes,
			Suggestions: suggestions,
		}
		json.NewEncoder(os.Stdout).Encode(summary)
		return
	}

	fmt.Println("==========================================================================================")
	fmt.Printf(" Compute Resource Capacity Analysis & Deployment Suggestions (%s)\n", client.BaseURL)
	fmt.Println("==========================================================================================")

	if len(suggestions) == 0 {
		fmt.Println("No active node metrics detected in VictoriaMetrics.")
		fmt.Println("Tip: Ensure node_exporter or system agents are running and scraping into VictoriaMetrics.")
		return
	}

	for _, s := range suggestions {
		fmt.Printf("\n[Rank #%d] Target Node: %s (%s)\n", s.CapacityRank, s.Instance, s.Hostname)
		fmt.Printf("  - Capacity Score:  %.1f / 100\n", s.Score)
		fmt.Printf("  - CPU Idle:        %.1f%%\n", s.CPUIdlePct)
		fmt.Printf("  - Available RAM:   %.1f GB\n", s.AvailRAMGB)
		fmt.Printf("  - Free Disk (/):   %.1f GB\n", s.AvailDiskGB)
		fmt.Printf("  - Suitability:     %s\n", s.Suitability)
		if len(s.WorkloadTypes) > 0 {
			fmt.Printf("  - Best For:        %s\n", strings.Join(s.WorkloadTypes, ", "))
		}
		fmt.Printf("  - Analysis:        %s\n", s.Reasoning)
		fmt.Printf("  - Action:          %s\n", s.Recommendation)
	}
	fmt.Println("\n==========================================================================================")
}
