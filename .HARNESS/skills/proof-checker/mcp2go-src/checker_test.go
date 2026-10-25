package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func parseSnippet(t *testing.T, src string) (*token.FileSet, *token.FileSet, *RuleTestCase) {
	// helper
	return nil, nil, nil
}

type RuleTestCase struct {
	Name           string
	Code           string
	ExpectedDiags  int
	ExpectedRuleID string
}

func TestNoErrorSwallowRule(t *testing.T) {
	rule := &NoErrorSwallowRule{}

	tests := []RuleTestCase{
		{
			Name: "Valid error handling",
			Code: `package test
import "fmt"
func run() error {
	err := doSomething()
	if err != nil {
		return fmt.Errorf("failed: %w", err)
	}
	return nil
}`,
			ExpectedDiags: 0,
		},
		{
			Name: "Silent discard _ = err",
			Code: `package test
func run() {
	err := doSomething()
	_ = err
}`,
			ExpectedDiags:  1,
			ExpectedRuleID: "ERR-001",
		},
		{
			Name: "Empty error block if err != nil {}",
			Code: `package test
func run() {
	err := doSomething()
	if err != nil {
	}
}`,
			ExpectedDiags:  1,
			ExpectedRuleID: "ERR-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.Code, 0)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			diags := rule.CheckFile(fset, file, "test.go")
			if len(diags) != tt.ExpectedDiags {
				t.Errorf("Expected %d diags, got %d", tt.ExpectedDiags, len(diags))
			}
			if tt.ExpectedDiags > 0 && len(diags) > 0 {
				if diags[0].RuleID != tt.ExpectedRuleID {
					t.Errorf("Expected RuleID %s, got %s", tt.ExpectedRuleID, diags[0].RuleID)
				}
			}
		})
	}
}

func TestNestingFlowRule(t *testing.T) {
	rule := &NestingFlowRule{MaxAllowedDepth: 3}

	tests := []RuleTestCase{
		{
			Name: "Shallow code with early return",
			Code: `package test
func process(x int) int {
	if x <= 0 {
		return 0
	}
	for i := 0; i < x; i++ {
		if i%2 == 0 {
			continue
		}
	}
	return x
}`,
			ExpectedDiags: 0,
		},
		{
			Name: "Deeply nested code (depth 4)",
			Code: `package test
func deep(x int) {
	if x > 0 {
		for i := 0; i < x; i++ {
			if i > 2 {
				for j := 0; j < i; j++ {
					print(j)
				}
			}
		}
	}
}`,
			ExpectedDiags:  1,
			ExpectedRuleID: "NEST-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.Code, 0)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			diags := rule.CheckFile(fset, file, "test.go")
			if len(diags) != tt.ExpectedDiags {
				t.Errorf("Expected %d diags, got %d", tt.ExpectedDiags, len(diags))
			}
		})
	}
}

func TestResourceCloseRule(t *testing.T) {
	rule := &ResourceCloseRule{}

	tests := []RuleTestCase{
		{
			Name: "Resource opened and deferred Close",
			Code: `package test
import "os"
func readFile() error {
	f, err := os.Open("data.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}`,
			ExpectedDiags: 0,
		},
		{
			Name: "Resource opened without defer Close",
			Code: `package test
import "os"
func leakFile() error {
	f, err := os.Open("data.txt")
	if err != nil {
		return err
	}
	buf := make([]byte, 100)
	f.Read(buf)
	return nil
}`,
			ExpectedDiags:  1,
			ExpectedRuleID: "RES-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.Code, 0)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			diags := rule.CheckFile(fset, file, "test.go")
			if len(diags) != tt.ExpectedDiags {
				t.Errorf("Expected %d diags, got %d", tt.ExpectedDiags, len(diags))
			}
		})
	}
}

func TestEffectPurityRule(t *testing.T) {
	rule := &EffectPurityRule{}

	tests := []RuleTestCase{
		{
			Name: "Verb prefix on effectful function",
			Code: `package test
import "net/http"
func fetchUserData(url string) (*http.Response, error) {
	return http.Get(url)
}`,
			ExpectedDiags: 0,
		},
		{
			Name: "Missing verb prefix on effectful function",
			Code: `package test
import "net/http"
func userData(url string) (*http.Response, error) {
	return http.Get(url)
}`,
			ExpectedDiags:  1,
			ExpectedRuleID: "PUR-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.Code, 0)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			diags := rule.CheckFile(fset, file, "test.go")
			if len(diags) != tt.ExpectedDiags {
				t.Errorf("Expected %d diags, got %d", tt.ExpectedDiags, len(diags))
			}
		})
	}
}
