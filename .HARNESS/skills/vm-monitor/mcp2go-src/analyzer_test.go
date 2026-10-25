package main

import (
	"testing"
)

func TestGenerateDeploymentSuggestions(t *testing.T) {
	nodes := []NodeResource{
		{
			Instance:      "100.84.48.65:9100",
			Hostname:      "docker",
			CPUIdlePct:    85.0,
			CPUUsagePct:   15.0,
			MemTotalBytes: 32 * 1024 * 1024 * 1024,
			MemAvailBytes: 24 * 1024 * 1024 * 1024, // 24GB free
			DiskTotalBytes: 500 * 1024 * 1024 * 1024,
			DiskAvailBytes: 300 * 1024 * 1024 * 1024,
		},
		{
			Instance:      "100.70.36.130:9100",
			Hostname:      "letwir-main",
			CPUIdlePct:    20.0,
			CPUUsagePct:   80.0,
			MemTotalBytes: 16 * 1024 * 1024 * 1024,
			MemAvailBytes: 1 * 1024 * 1024 * 1024, // 1GB free
			DiskTotalBytes: 500 * 1024 * 1024 * 1024,
			DiskAvailBytes: 50 * 1024 * 1024 * 1024,
		},
		{
			Instance:      "100.118.42.33:9100",
			Hostname:      "pi-hole",
			CPUIdlePct:    70.0,
			CPUUsagePct:   30.0,
			MemTotalBytes: 4 * 1024 * 1024 * 1024,
			MemAvailBytes: 2 * 1024 * 1024 * 1024, // 2GB free
			DiskTotalBytes: 64 * 1024 * 1024 * 1024,
			DiskAvailBytes: 30 * 1024 * 1024 * 1024,
		},
	}

	suggestions := GenerateDeploymentSuggestions(nodes)
	if len(suggestions) != 3 {
		t.Fatalf("expected 3 suggestions, got %d", len(suggestions))
	}

	// First rank should be the docker node with 24GB RAM and 85% CPU idle
	if suggestions[0].Hostname != "docker" {
		t.Errorf("expected top rank to be docker, got %s", suggestions[0].Hostname)
	}

	if suggestions[0].CapacityRank != 1 {
		t.Errorf("expected rank 1, got %d", suggestions[0].CapacityRank)
	}

	if suggestions[2].Hostname != "letwir-main" {
		t.Errorf("expected lowest rank to be letwir-main, got %s", suggestions[2].Hostname)
	}
}
