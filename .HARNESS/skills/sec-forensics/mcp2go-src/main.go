package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sec-forensics/forensics"
	"sec-forensics/github"
	"sec-forensics/knowledge"
	"sec-forensics/scanner"
	"strings"
)

// [Chaotic×ojou-sama] CLIエントリポイント！冷徹なるセキュリティ・フォレンジックとCTF免疫監査を執行いたしますわ！
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "dependabot":
		handleDependabotCommandComplex(os.Args[2:])
	case "scan-prompt":
		handleScanPromptCommandComplex(os.Args[2:])
	case "scan-ctf":
		handleScanCTFCommandComplex(os.Args[2:])
	case "audit-log":
		handleAuditLogCommandComplex(os.Args[2:])
	case "cve-fetch":
		handleCVEFetchCommandComplex(os.Args[2:])
	case "cve-ingest":
		handleCVEIngestCommandComplex(os.Args[2:])
	case "mega-check":
		handleMegaCheckCommandComplex(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("sec-forensics: Security Forensics, CTF Immunity & Knowledge Auditor")
	fmt.Println("\nUsage:")
	fmt.Println("  sec-forensics dependabot -repo <owner/repo> [-json]")
	fmt.Println("  sec-forensics scan-prompt (-file <path> | -text <string>) [-json]")
	fmt.Println("  sec-forensics scan-ctf (-file <path> | -text <string>) [-json]")
	fmt.Println("  sec-forensics audit-log -transcript <path/to/transcript.jsonl> [-json]")
	fmt.Println("  sec-forensics cve-fetch [-query <str>] [-cat <category>] [-json]")
	fmt.Println("  sec-forensics cve-ingest (-seed | -title <t> -content <c> [-cat <cat>] [-tags <t>])")
	fmt.Println("  sec-forensics mega-check [-repo <owner/repo>] [-docs <dir>] [-transcript <path>] [-json]")
}

func handleDependabotCommandComplex(args []string) {
	fs := flag.NewFlagSet("dependabot", flag.ExitOnError)
	repo := fs.String("repo", "", "Target repository (e.g. owner/repo)")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	if *repo == "" {
		fmt.Fprintln(os.Stderr, "Error: -repo is required")
		os.Exit(1)
	}

	report, err := github.FetchDependabotAlertsComplex(*repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dependabot error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("=== Dependabot Triage Report: %s ===\n", report.Repo)
	fmt.Printf("Total Alerts: %d (Open: %d, Fixed: %d) | Max Severity: %s\n\n",
		report.TotalAlerts, report.OpenAlerts, report.FixedAlerts, report.MaxSeverity)
	for _, it := range report.Items {
		fmt.Printf("[#%d] %s (%s) - %s (CVE: %s, GHSA: %s)\n",
			it.Number, it.Package, it.Severity, it.State, it.CVE, it.GHSA)
		fmt.Printf("      Summary: %s\n", it.Summary)
		fmt.Printf("      Fix: %s | URL: %s\n", it.PatchedVer, it.URL)
	}
}

func handleScanPromptCommandComplex(args []string) {
	fs := flag.NewFlagSet("scan-prompt", flag.ExitOnError)
	filePath := fs.String("file", "", "Target file to scan")
	text := fs.String("text", "", "Raw text to scan")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	var content string
	if *filePath != "" {
		data, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "File read error: %v\n", err)
			os.Exit(1)
		}
		content = string(data)
	} else if *text != "" {
		content = *text
	} else {
		fmt.Fprintln(os.Stderr, "Error: either -file or -text must be provided")
		os.Exit(1)
	}

	report := scanner.ScanTextComplex(content)
	if *jsonOut {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("=== Prompt Injection Scan Report ===\n")
	fmt.Printf("Total Findings: %d | Max Severity: %s\n\n", report.TotalFindings, report.MaxSeverity)
	for _, f := range report.Findings {
		fmt.Printf("[Line %d] [%s] Vector: %s\n", f.LineNumber, f.Severity, f.Vector)
		fmt.Printf("  Description: %s\n", f.Description)
		fmt.Printf("  Snippet: %s\n\n", f.Snippet)
	}
}

func handleScanCTFCommandComplex(args []string) {
	fs := flag.NewFlagSet("scan-ctf", flag.ExitOnError)
	filePath := fs.String("file", "", "Target file to scan")
	text := fs.String("text", "", "Raw text to scan")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	var content string
	if *filePath != "" {
		data, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "File read error: %v\n", err)
			os.Exit(1)
		}
		content = string(data)
	} else if *text != "" {
		content = *text
	} else {
		fmt.Fprintln(os.Stderr, "Error: either -file or -text must be provided")
		os.Exit(1)
	}

	report := scanner.ScanCTFImmunityComplex(content)
	if *jsonOut {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("=== CTF/SECCON/DEFCON Immunity Scan Report ===\n")
	fmt.Printf("Total Findings: %d | Max Severity: %s\n\n", report.TotalFindings, report.MaxSeverity)
	for _, f := range report.Findings {
		fmt.Printf("[Line %d] [%s] Vector: %s\n", f.LineNumber, f.Severity, f.Vector)
		fmt.Printf("  Description: %s\n", f.Description)
		fmt.Printf("  Snippet: %s\n\n", f.Snippet)
	}
}

func handleAuditLogCommandComplex(args []string) {
	fs := flag.NewFlagSet("audit-log", flag.ExitOnError)
	transcriptPath := fs.String("transcript", "", "Path to transcript.jsonl")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	if *transcriptPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -transcript is required")
		os.Exit(1)
	}

	report, err := forensics.AuditTranscriptFileComplex(*transcriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Forensics audit error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("=== Conversation Log Forensics Report ===\n")
	fmt.Printf("Path: %s\n", report.TranscriptPath)
	fmt.Printf("Total Steps: %d | Ingests: %d | Anomalies: %d | Max Risk: %s\n\n",
		report.TotalSteps, report.UntrustedIngests, report.AnomaliesFound, report.MaxRiskLevel)
	for _, ev := range report.Timeline {
		fmt.Printf("[Step %d] [%s] Event: %s (Tool: %s)\n", ev.StepIndex, ev.RiskLevel, ev.EventType, ev.ToolName)
		fmt.Printf("  Summary: %s\n", ev.Summary)
		if ev.MatchedSign != "" {
			fmt.Printf("  Matched Signature: %s\n", ev.MatchedSign)
		}
		if ev.Details != "" {
			fmt.Printf("  Details: %s\n", ev.Details)
		}
		fmt.Println()
	}
}

func handleCVEFetchCommandComplex(args []string) {
	fs := flag.NewFlagSet("cve-fetch", flag.ExitOnError)
	query := fs.String("query", "", "Search query for CVE or package")
	category := fs.String("cat", "sec_cve", "Category (sec_cve or sec_immunity)")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	results, err := knowledge.FetchSecurityKnowledgeComplex(*query, *category)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CVE Fetch error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("=== Pre-Investigation Security Knowledge Fetch (Category: %s) ===\n", *category)
	fmt.Printf("Found %d intelligence records\n\n", len(results))
	for _, r := range results {
		fmt.Printf("📌 [%s] %s\n", r.Category, r.Title)
		if r.L2Summary != "" {
			fmt.Printf("   Summary: %s\n", r.L2Summary)
		}
		if len(r.Tags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(r.Tags, ", "))
		}
		fmt.Println()
	}
}

func handleCVEIngestCommandComplex(args []string) {
	fs := flag.NewFlagSet("cve-ingest", flag.ExitOnError)
	seed := fs.Bool("seed", false, "Populate baseline CVE & CTF intelligence")
	title := fs.String("title", "", "Intelligence title")
	content := fs.String("content", "", "Intelligence content")
	cat := fs.String("cat", "sec_cve", "Category (sec_cve, sec_immunity)")
	tags := fs.String("tags", "security,cve", "Comma-separated tags")
	_ = fs.Parse(args)

	if *seed {
		fmt.Println("🚀 Seeding baseline Pillow CVEs and CTF attack vectors into PostgreSQL (llm-mem)...")
		if err := knowledge.SeedInitialKnowledgeComplex(); err != nil {
			fmt.Fprintf(os.Stderr, "Seed error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Successfully seeded security knowledge into PostgreSQL!")
		return
	}

	if *title == "" || *content == "" {
		fmt.Fprintln(os.Stderr, "Error: either -seed or both -title and -content must be provided")
		os.Exit(1)
	}

	tagList := strings.Split(*tags, ",")
	err := knowledge.IngestSecurityKnowledgeComplex(*title, *content, *cat, "", "", tagList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ingest error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully registered security knowledge: %s (Category: %s)\n", *title, *cat)
}

// MegaCheckReport holds summary of comprehensive pre-docs update security audit.
type MegaCheckReport struct {
	Repo             string `json:"repo,omitempty"`
	DocsDir          string `json:"docs_dir,omitempty"`
	CTFImmunityPass  bool   `json:"ctf_immunity_pass"`
	DependabotPass   bool   `json:"dependabot_pass"`
	PromptInjPass    bool   `json:"prompt_inj_pass"`
	OverallVerdict   string `json:"overall_verdict"`
	Summary          string `json:"summary"`
}

func handleMegaCheckCommandComplex(args []string) {
	fs := flag.NewFlagSet("mega-check", flag.ExitOnError)
	repo := fs.String("repo", "", "Target GitHub repo (e.g. owner/repo)")
	docsDir := fs.String("docs", ".", "Directory of markdown docs to audit")
	transcript := fs.String("transcript", "", "Path to transcript.jsonl (optional)")
	jsonOut := fs.Bool("json", false, "Output structured JSON")
	_ = fs.Parse(args)

	ctfPass := true
	promptPass := true
	dependabotPass := true
	findingsCount := 0

	// 1. Audit Markdown files in docsDir
	_ = filepath.WalkDir(*docsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == ".system_generated" || name == "brain" || name == "antigravity" || name == "history" || name == "tmp" || name == "subagents" || name == "plugins" || name == "skills" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, ".toml") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content := string(data)

		ctfRep := scanner.ScanCTFImmunityComplex(content)
		if ctfRep.TotalFindings > 0 {
			ctfPass = false
			findingsCount += ctfRep.TotalFindings
		}

		injRep := scanner.ScanTextComplex(content)
		if injRep.TotalFindings > 0 {
			promptPass = false
			findingsCount += injRep.TotalFindings
		}
		return nil
	})

	// 2. Dependabot check if repo is given
	if *repo != "" {
		repReport, err := github.FetchDependabotAlertsComplex(*repo)
		if err == nil && repReport.OpenAlerts > 0 && repReport.MaxSeverity == scanner.SeverityCritical {
			dependabotPass = false
		}
	}

	// 3. Transcript check if provided
	if *transcript != "" {
		forensicsRep, err := forensics.AuditTranscriptFileComplex(*transcript)
		if err == nil && forensicsRep.MaxRiskLevel == scanner.SeverityCritical {
			promptPass = false
		}
	}

	verdict := "PASS"
	if !ctfPass || !promptPass || !dependabotPass {
		verdict = "REJECT"
	}

	report := MegaCheckReport{
		Repo:            *repo,
		DocsDir:         *docsDir,
		CTFImmunityPass: ctfPass,
		DependabotPass:  dependabotPass,
		PromptInjPass:   promptPass,
		OverallVerdict:  verdict,
		Summary:         fmt.Sprintf("Mega-Docs Security Audit: %d findings. Overall Verdict: %s", findingsCount, verdict),
	}

	if *jsonOut {
		data, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("================================================================================\n")
	fmt.Printf("🛡️ MEGA-DOCS-UPDATE COMPREHENSIVE SECURITY & IMMUNITY GATE\n")
	fmt.Printf("================================================================================\n")
	fmt.Printf("Docs Dir        : %s\n", report.DocsDir)
	fmt.Printf("Repo            : %s\n", report.Repo)
	fmt.Printf("CTF Immunity    : %s\n", formatPassFail(report.CTFImmunityPass))
	fmt.Printf("Prompt Injection: %s\n", formatPassFail(report.PromptInjPass))
	fmt.Printf("Dependabot Audit: %s\n", formatPassFail(report.DependabotPass))
	fmt.Printf("--------------------------------------------------------------------------------\n")
	fmt.Printf("Overall Verdict : %s\n", report.OverallVerdict)
	fmt.Printf("================================================================================\n")

	if verdict == "REJECT" {
		os.Exit(1)
	}
}

func formatPassFail(passed bool) string {
	if passed {
		return "PASS (Clean)"
	}
	return "FAIL (Threats Detected)"
}
