package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ExtractMetricFloat extracts float value from QueryResult
func ExtractMetricFloat(result QueryResult) (float64, error) {
	if len(result.Value) < 2 {
		return 0, fmt.Errorf("result value array is too short")
	}

	rawVal, ok := result.Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("result value is not string")
	}

	val, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse metric float %q: %w", rawVal, err)
	}
	return val, nil
}

// CollectClusterResources fetches and aggregates metrics across all discovered nodes
func CollectClusterResources(client *VMClient) ([]NodeResource, error) {
	nodeMap := make(map[string]*NodeResource)

	ensureNode := func(instance, hostname string) *NodeResource {
		key := instance
		if key == "" {
			key = hostname
		}
		if key == "" {
			key = "unknown"
		}
		nr, exists := nodeMap[key]
		if !exists {
			nr = &NodeResource{
				Instance:  instance,
				Hostname:  hostname,
				Timestamp: time.Now(),
			}
			nodeMap[key] = nr
		}
		if nr.Hostname == "" && hostname != "" {
			nr.Hostname = hostname
		}
		return nr
	}

	// 1. CPU Idle % query
	cpuQuery := `100 - (avg by (instance, nodename, host) (rate(node_cpu_seconds_total{mode!="idle"}[5m])) * 100)`
	if cpuResp, err := client.QueryInstant(cpuQuery); err == nil && cpuResp.Status == "success" {
		for _, res := range cpuResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			if host == "" {
				host = res.Metric["host"]
			}
			val, err := ExtractMetricFloat(res)
			if err == nil {
				node := ensureNode(inst, host)
				node.CPUIdlePct = val
				node.CPUUsagePct = 100.0 - val
				if node.CPUUsagePct < 0 {
					node.CPUUsagePct = 0
				}
			}
		}
	}

	// 2. Memory Available query
	memAvailQuery := `node_memory_MemAvailable_bytes`
	if memResp, err := client.QueryInstant(memAvailQuery); err == nil && memResp.Status == "success" {
		for _, res := range memResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			val, err := ExtractMetricFloat(res)
			if err == nil && val >= 0 {
				node := ensureNode(inst, host)
				node.MemAvailBytes = uint64(val)
			}
		}
	}

	// 3. Memory Total query
	memTotalQuery := `node_memory_MemTotal_bytes`
	if memTotResp, err := client.QueryInstant(memTotalQuery); err == nil && memTotResp.Status == "success" {
		for _, res := range memTotResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			val, err := ExtractMetricFloat(res)
			if err == nil && val > 0 {
				node := ensureNode(inst, host)
				node.MemTotalBytes = uint64(val)
				if node.MemAvailBytes > 0 && node.MemTotalBytes >= node.MemAvailBytes {
					used := node.MemTotalBytes - node.MemAvailBytes
					node.MemUsagePct = (float64(used) / float64(node.MemTotalBytes)) * 100.0
				}
			}
		}
	}

	// 4. Disk Free query (root mount)
	diskAvailQuery := `node_filesystem_avail_bytes{mountpoint="/"}`
	if diskResp, err := client.QueryInstant(diskAvailQuery); err == nil && diskResp.Status == "success" {
		for _, res := range diskResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			val, err := ExtractMetricFloat(res)
			if err == nil && val >= 0 {
				node := ensureNode(inst, host)
				node.DiskAvailBytes = uint64(val)
			}
		}
	}

	// 5. Disk Total query (root mount)
	diskTotQuery := `node_filesystem_size_bytes{mountpoint="/"}`
	if diskTotResp, err := client.QueryInstant(diskTotQuery); err == nil && diskTotResp.Status == "success" {
		for _, res := range diskTotResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			val, err := ExtractMetricFloat(res)
			if err == nil && val > 0 {
				node := ensureNode(inst, host)
				node.DiskTotalBytes = uint64(val)
				if node.DiskAvailBytes > 0 && node.DiskTotalBytes >= node.DiskAvailBytes {
					used := node.DiskTotalBytes - node.DiskAvailBytes
					node.DiskUsagePct = (float64(used) / float64(node.DiskTotalBytes)) * 100.0
				}
			}
		}
	}

	// 6. Load Average (1m)
	load1Query := `node_load1`
	if loadResp, err := client.QueryInstant(load1Query); err == nil && loadResp.Status == "success" {
		for _, res := range loadResp.Data.Result {
			inst := res.Metric["instance"]
			host := res.Metric["nodename"]
			val, err := ExtractMetricFloat(res)
			if err == nil {
				node := ensureNode(inst, host)
				node.Load1 = val
			}
		}
	}

	// 7. CPU Cores count
	coreQuery := `count by (instance) (node_cpu_seconds_total{mode="idle"})`
	if coreResp, err := client.QueryInstant(coreQuery); err == nil && coreResp.Status == "success" {
		for _, res := range coreResp.Data.Result {
			inst := res.Metric["instance"]
			val, err := ExtractMetricFloat(res)
			if err == nil && val > 0 {
				node := ensureNode(inst, "")
				node.CPUCores = int(val)
			}
		}
	}

	nodes := make([]NodeResource, 0, len(nodeMap))
	for _, n := range nodeMap {
		nodes = append(nodes, *n)
	}

	// Sort nodes alphabetically by instance name
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Instance < nodes[j].Instance
	})

	return nodes, nil
}

// GenerateDeploymentSuggestions analyzes node resources and generates optimal deployment suggestions
func GenerateDeploymentSuggestions(nodes []NodeResource) []DeploymentSuggestion {
	var suggestions []DeploymentSuggestion

	for _, n := range nodes {
		availRAMGB := float64(n.MemAvailBytes) / (1024 * 1024 * 1024)
		availDiskGB := float64(n.DiskAvailBytes) / (1024 * 1024 * 1024)

		// Calculate composite capacity score (0 ~ 100)
		// 45% RAM availability, 40% CPU Idle, 15% Disk availability
		ramScore := 0.0
		if n.MemTotalBytes > 0 {
			ramScore = (float64(n.MemAvailBytes) / float64(n.MemTotalBytes)) * 100.0
		}
		cpuScore := n.CPUIdlePct
		diskScore := 0.0
		if n.DiskTotalBytes > 0 {
			diskScore = (float64(n.DiskAvailBytes) / float64(n.DiskTotalBytes)) * 100.0
		}

		score := (ramScore * 0.45) + (cpuScore * 0.40) + (diskScore * 0.15)
		if n.MemTotalBytes == 0 && n.CPUIdlePct == 0 {
			// Fallback if detailed node_exporter metrics are missing
			score = 50.0
		}

		var workloads []string
		var suitability string
		var reasoning string
		var recommendation string

		if cpuScore < 15.0 || (n.MemTotalBytes > 0 && ramScore < 10.0) {
			suitability = "Constrained / Highly Loaded"
			reasoning = fmt.Sprintf("High utilization detected (CPU Idle: %.1f%%, Avail RAM: %.1f GB). Additional workload may cause thrashing.", cpuScore, availRAMGB)
			recommendation = "Do NOT deploy new services. Consider scaling out or migrating existing loads."
		} else if availRAMGB >= 8.0 && cpuScore >= 50.0 {
			suitability = "High Capacity (Heavy Workloads)"
			workloads = []string{"LLM Inference / Ollama", "Local Embedding Model", "Heavy Build Pipeline", "Database / Vector Search"}
			reasoning = fmt.Sprintf("Abundant memory (%.1f GB free) and CPU idle capacity (%.1f%% idle).", availRAMGB, cpuScore)
			recommendation = "Ideal for hosting local AI/LLM models, embedding services, or intensive batch processing."
		} else if availRAMGB >= 3.0 && cpuScore >= 40.0 {
			suitability = "Medium Capacity (Standard Services)"
			workloads = []string{"Backend API Servers", "Web Applications", "CI/CD Runners", "Message Queues / Redis"}
			reasoning = fmt.Sprintf("Balanced resources available (%.1f GB RAM free, %.1f%% CPU idle).", availRAMGB, cpuScore)
			recommendation = "Well-suited for containerized microservices, web apps, or auxiliary workers."
		} else {
			suitability = "Light Capacity (Low Overhead Services)"
			workloads = []string{"Lightweight Webhook Listener", "DNS / Proxy Service", "Cron Jobs / Scheduled Checkers", "Prometheus Exporters"}
			reasoning = fmt.Sprintf("Modest headroom (%.1f GB RAM free, %.1f%% CPU idle).", availRAMGB, cpuScore)
			recommendation = "Recommended for lightweight Go/Rust micro-daemons or edge utilities with low memory footprints."
		}

		host := n.Hostname
		if host == "" {
			host = strings.Split(n.Instance, ":")[0]
		}

		suggestions = append(suggestions, DeploymentSuggestion{
			Instance:       n.Instance,
			Hostname:       host,
			Score:          score,
			CPUIdlePct:     n.CPUIdlePct,
			AvailRAMGB:     availRAMGB,
			AvailDiskGB:    availDiskGB,
			Suitability:    suitability,
			WorkloadTypes:  workloads,
			Reasoning:      reasoning,
			Recommendation: recommendation,
		})
	}

	// Sort suggestions by Score descending (highest available capacity first)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	for i := range suggestions {
		suggestions[i].CapacityRank = i + 1
	}

	return suggestions
}
