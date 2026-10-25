package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sec-forensics/scanner"
	"time"
)

// DependabotAlert represents a GitHub Dependabot alert structure.
type DependabotAlert struct {
	Number      int                    `json:"number"`
	State       string                 `json:"state"`
	HTMLURL     string                 `json:"html_url"`
	CreatedAt   string                 `json:"created_at"`
	FixedAt     *string                `json:"fixed_at,omitempty"`
	DismissedAt *string                `json:"dismissed_at,omitempty"`
	SecurityAdv SecurityAdvisoryDetail `json:"security_advisory"`
	Vulnerable  SecurityVulnerability  `json:"security_vulnerability"`
}

// SecurityAdvisoryDetail contains advisory specifics.
type SecurityAdvisoryDetail struct {
	GHSAID      string            `json:"ghsa_id"`
	CVEID       *string           `json:"cve_id,omitempty"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Severity    string            `json:"severity"`
	CVSS        CVSSDetail        `json:"cvss"`
	CWEs        []CWEInfo         `json:"cwes"`
}

// SecurityVulnerability details affected package and patched version.
type SecurityVulnerability struct {
	Package struct {
		Ecosystem string `json:"ecosystem"`
		Name      string `json:"name"`
	} `json:"package"`
	Severity              string  `json:"severity"`
	VulnerableRange       string  `json:"vulnerable_version_range"`
	FirstPatchedVer       *struct {
		Identifier string `json:"identifier"`
	} `json:"first_patched_version,omitempty"`
}

// CVSSDetail holds CVSS score and vector.
type CVSSDetail struct {
	Score        float64 `json:"score"`
	VectorString *string `json:"vector_string,omitempty"`
}

// CWEInfo holds CWE id and name.
type CWEInfo struct {
	CWEID string `json:"cwe_id"`
	Name  string `json:"name"`
}

// DependabotTriageReport summarizes repo security alerts.
type DependabotTriageReport struct {
	Repo           string            `json:"repo"`
	TotalAlerts    int               `json:"total_alerts"`
	OpenAlerts     int               `json:"open_alerts"`
	FixedAlerts    int               `json:"fixed_alerts"`
	MaxSeverity    scanner.Severity  `json:"max_severity"`
	Items          []AlertSummary    `json:"items"`
}

// AlertSummary is a condensed view of an alert for security report.
type AlertSummary struct {
	Number       int      `json:"number"`
	State        string   `json:"state"`
	Package      string   `json:"package"`
	Ecosystem    string   `json:"ecosystem"`
	Severity     string   `json:"severity"`
	CVE          string   `json:"cve"`
	GHSA         string   `json:"ghsa"`
	Summary      string   `json:"summary"`
	PatchedVer   string   `json:"patched_version"`
	URL          string   `json:"url"`
}

// FetchDependabotAlertsComplex retrieves and parses alerts via gh CLI.
// [Chaotic×ojou-sama] GitHub CLIを呼び出し、依存関係の脆弱性を一切合切摘み上げますわ！
func FetchDependabotAlertsComplex(repo string) (*DependabotTriageReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gh", "api", fmt.Sprintf("repos/%s/dependabot/alerts", repo))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("FetchDependabotAlerts: gh api failed for %s: %w", repo, err)
	}

	var rawAlerts []DependabotAlert
	if err := json.Unmarshal(out, &rawAlerts); err != nil {
		return nil, fmt.Errorf("FetchDependabotAlerts: json decode failed: %w", err)
	}

	openCount := 0
	fixedCount := 0
	items := make([]AlertSummary, 0, len(rawAlerts))

	for _, a := range rawAlerts {
		if a.State == "open" {
			openCount++
		} else if a.State == "fixed" {
			fixedCount++
		}

		cve := ""
		if a.SecurityAdv.CVEID != nil {
			cve = *a.SecurityAdv.CVEID
		}

		patched := "None"
		if a.Vulnerable.FirstPatchedVer != nil {
			patched = a.Vulnerable.FirstPatchedVer.Identifier
		}

		items = append(items, AlertSummary{
			Number:     a.Number,
			State:      a.State,
			Package:    a.Vulnerable.Package.Name,
			Ecosystem:  a.Vulnerable.Package.Ecosystem,
			Severity:   a.SecurityAdv.Severity,
			CVE:        cve,
			GHSA:       a.SecurityAdv.GHSAID,
			Summary:    a.SecurityAdv.Summary,
			PatchedVer: patched,
			URL:        a.HTMLURL,
		})
	}

	maxSev := scanner.SeverityLow
	for _, it := range items {
		if it.State == "open" && it.Severity == "critical" {
			maxSev = scanner.SeverityCritical
			break
		}
		if it.State == "open" && it.Severity == "high" {
			maxSev = scanner.SeverityHigh
		}
	}

	return &DependabotTriageReport{
		Repo:        repo,
		TotalAlerts: len(rawAlerts),
		OpenAlerts:  openCount,
		FixedAlerts: fixedCount,
		MaxSeverity: maxSev,
		Items:       items,
	}, nil
}
