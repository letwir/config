package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Engine manages and executes all verification rules
type Engine struct {
	rules []Rule
	fset  *token.FileSet
}

// NewEngine constructs a new proof-check engine with default rules
func NewEngine() *Engine {
	return &Engine{
		fset: token.NewFileSet(),
		rules: []Rule{
			&NoErrorSwallowRule{},
			&NestingFlowRule{MaxAllowedDepth: 3},
			&ResourceCloseRule{},
			&EffectPurityRule{},
			&FunctionLengthRule{MaxLines: 60},
			&ConcurrencyRule{},
		},
	}
}

// CheckDirectory recursively scans and validates Go source files in the target directory
func (e *Engine) CheckDirectory(cfg Config) ([]Diagnostic, error) {
	var allDiags []Diagnostic

	err := filepath.Walk(cfg.TargetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor, git, venv, and hidden directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		diags, parseErr := e.CheckFile(path)
		if parseErr != nil {
			allDiags = append(allDiags, Diagnostic{
				RuleID:   "SYNTAX-001",
				Category: "SYNTAX",
				FilePath: path,
				Line:     1,
				Column:   1,
				Message:  fmt.Sprintf("Failed to parse Go source: %v", parseErr),
				Severity: SeverityError,
			})
			return nil
		}

		allDiags = append(allDiags, diags...)
		return nil
	})

	if err != nil {
		return allDiags, err
	}

	// Optionally run go vet
	if cfg.IncludeVet {
		vetDiags := runGoVet(cfg.TargetDir)
		allDiags = append(allDiags, vetDiags...)
	}

	return allDiags, nil
}

// CheckFile parses and executes rules on a single Go file
func (e *Engine) CheckFile(filePath string) ([]Diagnostic, error) {
	var fileDiags []Diagnostic

	node, err := parser.ParseFile(e.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, rule := range e.rules {
		diags := rule.CheckFile(e.fset, node, filePath)
		fileDiags = append(fileDiags, diags...)
	}

	return fileDiags, nil
}

// runGoVet executes 'go vet ./...' in the target directory
func runGoVet(dir string) []Diagnostic {
	var diags []Diagnostic

	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	_ = cmd.Run() // Exit code ignored; we parse stderr output

	output := stderr.String()
	if strings.TrimSpace(output) == "" {
		return diags
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, ":", 4)
		if len(parts) >= 3 {
			filePath := parts[0]
			msg := trimmed
			if len(parts) >= 4 {
				msg = strings.TrimSpace(parts[3])
			}

			diags = append(diags, Diagnostic{
				RuleID:   "GO-VET",
				Category: "STATIC_ANALYSIS",
				FilePath: filePath,
				Line:     1,
				Column:   1,
				Message:  msg,
				Severity: SeverityError,
			})
		}
	}

	return diags
}
