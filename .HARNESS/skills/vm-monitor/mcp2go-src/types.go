package main

import "time"

// QueryResponse represents the Prometheus/MetricsQL instant query response
type QueryResponse struct {
	Status    string    `json:"status"`
	Data      QueryData `json:"data"`
	ErrorType string    `json:"errorType,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// QueryData contains query results
type QueryData struct {
	ResultType string        `json:"resultType"`
	Result     []QueryResult `json:"result"`
}

// QueryResult represents a single metric series result
type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []any             `json:"value,omitempty"`
	Values [][]any           `json:"values,omitempty"`
}

// LabelValuesResponse represents response from /api/v1/label/<name>/values
type LabelValuesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// TSDBStatusResponse represents response from /api/v1/status/tsdb
type TSDBStatusResponse struct {
	Status string   `json:"status"`
	Data   TSDBData `json:"data"`
}

// TSDBData represents TSDB status details
type TSDBData struct {
	TotalSeries                int64             `json:"totalSeries,omitempty"`
	TotalLabelValuePairs       int64             `json:"totalLabelValuePairs,omitempty"`
	SeriesCountByMetricName    []SeriesMetricRef `json:"seriesCountByMetricName,omitempty"`
	LabelValueCountByLabelName []LabelCountRef   `json:"labelValueCountByLabelName,omitempty"`
}

// SeriesMetricRef metric series count reference
type SeriesMetricRef struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// LabelCountRef label value count reference
type LabelCountRef struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// NodeResource represents aggregated resource metrics for a single node / instance
type NodeResource struct {
	Instance       string    `json:"instance"`
	Hostname       string    `json:"hostname"`
	CPUUsagePct    float64   `json:"cpu_usage_pct"`
	CPUIdlePct     float64   `json:"cpu_idle_pct"`
	CPUCores       int       `json:"cpu_cores"`
	MemTotalBytes  uint64    `json:"mem_total_bytes"`
	MemAvailBytes  uint64    `json:"mem_avail_bytes"`
	MemUsagePct    float64   `json:"mem_usage_pct"`
	DiskTotalBytes uint64    `json:"disk_total_bytes"`
	DiskAvailBytes uint64    `json:"disk_avail_bytes"`
	DiskUsagePct   float64   `json:"disk_usage_pct"`
	Load1          float64   `json:"load_1"`
	Load5          float64   `json:"load_5"`
	Load15         float64   `json:"load_15"`
	Timestamp      time.Time `json:"timestamp"`
}

// DeploymentSuggestion represents a recommendation for workload placement
type DeploymentSuggestion struct {
	Instance       string   `json:"instance"`
	Hostname       string   `json:"hostname"`
	Score          float64  `json:"score"`
	CapacityRank   int      `json:"capacity_rank"`
	CPUIdlePct     float64  `json:"cpu_idle_pct"`
	AvailRAMGB     float64  `json:"avail_ram_gb"`
	AvailDiskGB    float64  `json:"avail_disk_gb"`
	Suitability    string   `json:"suitability"`
	WorkloadTypes  []string `json:"workload_types"`
	Reasoning      string   `json:"reasoning"`
	Recommendation string   `json:"recommendation"`
}

// ClusterSummary represents high-level metrics across all nodes
type ClusterSummary struct {
	Endpoint    string                 `json:"endpoint"`
	CheckedAt   time.Time              `json:"checked_at"`
	TotalNodes  int                    `json:"total_nodes"`
	TotalSeries int64                  `json:"total_series"`
	Nodes       []NodeResource         `json:"nodes"`
	Suggestions []DeploymentSuggestion `json:"suggestions"`
}
