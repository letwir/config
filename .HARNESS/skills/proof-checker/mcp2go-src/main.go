package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	var cfg Config
	flag.StringVar(&cfg.TargetDir, "path", ".", "Target directory or file path to check")
	flag.BoolVar(&cfg.Strict, "strict", false, "Strict mode: treat warnings as errors")
	flag.BoolVar(&cfg.IncludeVet, "vet", false, "Include 'go vet' static analysis")
	flag.BoolVar(&cfg.JSONOutput, "json", false, "Output results in JSON format")
	flag.Parse()

	engine := NewEngine()
	diagnostics, err := engine.CheckDirectory(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory '%s': %v\n", cfg.TargetDir, err)
		os.Exit(2)
	}

	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, d := range diagnostics {
		switch d.Severity {
		case SeverityError:
			errorCount++
		case SeverityWarning:
			warningCount++
		case SeverityInfo:
			infoCount++
		}
	}

	if cfg.JSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]interface{}{
			"target":      cfg.TargetDir,
			"diagnostics": diagnostics,
			"summary": map[string]int{
				"errors":   errorCount,
				"warnings": warningCount,
				"info":     infoCount,
			},
			"verdict": calculateVerdict(errorCount, warningCount, cfg.Strict),
		})
	} else {
		printFormattedReport(cfg.TargetDir, diagnostics, errorCount, warningCount, infoCount, cfg.Strict)
	}

	if errorCount > 0 || (cfg.Strict && warningCount > 0) {
		os.Exit(1)
	}
	os.Exit(0)
}

func calculateVerdict(errors, warnings int, strict bool) string {
	if errors > 0 || (strict && warnings > 0) {
		return "REJECT"
	}
	return "PASS"
}

func printFormattedReport(target string, diagnostics []Diagnostic, errors, warnings, infos int, strict bool) {
	fmt.Println("================================================================================")
	fmt.Printf(" [PROOF-CHECKER] AST & Category Theory Soundness Verification Gate\n")
	fmt.Printf(" Target: %s | Strict: %v\n", target, strict)
	fmt.Println("================================================================================")

	if len(diagnostics) == 0 {
		fmt.Println(" [PASS] All source files satisfy CODE_RULE.md AST & Categorical Invariants!")
		fmt.Println(" Verdict: PASS (0 errors, 0 warnings, 0 infos)")
		fmt.Println("================================================================================")
		return
	}

	for _, d := range diagnostics {
		icon := "[INFO]"
		if d.Severity == SeverityError {
			icon = "[ERROR]"
		} else if d.Severity == SeverityWarning {
			icon = "[WARN]"
		}

		fmt.Printf("%s [%s] %s:%d:%d\n", icon, d.RuleID, d.FilePath, d.Line, d.Column)
		fmt.Printf("  Message:    %s\n", d.Message)
		if d.Suggestion != "" {
			fmt.Printf("  Suggestion: %s\n", d.Suggestion)
		}
		fmt.Println("--------------------------------------------------------------------------------")
	}

	verdict := calculateVerdict(errors, warnings, strict)
	fmt.Printf("Summary: %d Errors, %d Warnings, %d Infos\n", errors, warnings, infos)
	if verdict == "PASS" {
		fmt.Printf("Verdict: PASS (Proof checking succeeded)\n")
	} else {
		fmt.Printf("Verdict: REJECT (Violations must be fixed before proceeding)\n")
	}
	fmt.Println("================================================================================")
}
