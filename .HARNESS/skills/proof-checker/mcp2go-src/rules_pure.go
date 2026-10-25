package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// EffectPurityRule enforces Sep(IO_Effect, PureDomain) ∧ SideEffectFn[VerbPrefix]
type EffectPurityRule struct{}

func (r *EffectPurityRule) ID() string {
	return "PUR-001"
}

func (r *EffectPurityRule) Category() Category {
	return CategoryEffectPurity
}

func (r *EffectPurityRule) Description() string {
	return "Enforces separation of IO_Effect and PureDomain with appropriate verb prefixes for effectful functions"
}

func (r *EffectPurityRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic

	// Identify whether file is a test file or main package
	isTestFile := strings.HasSuffix(filePath, "_test.go")

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		if isTestFile && strings.HasPrefix(fn.Name.Name, "Test") {
			continue
		}
		if fn.Name.Name == "main" || fn.Name.Name == "init" {
			continue
		}

		checkFunctionPurity(fset, fn, filePath, &diagnostics)
	}

	return diagnostics
}

func checkFunctionPurity(fset *token.FileSet, fn *ast.FuncDecl, filePath string, diags *[]Diagnostic) {
	hasIO := false
	var ioTarget string

	ioPackages := map[string]bool{
		"http": true, "net": true, "sql": true, "os": true, "exec": true, "pgx": true,
	}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			target := getCallTargetName(call)
			parts := strings.Split(target, ".")
			if len(parts) == 2 {
				pkg := parts[0]
				// Avoid false positive on os.Getenv (read-only configuration)
				if pkg == "os" && parts[1] == "Getenv" {
					return true
				}
				if ioPackages[pkg] {
					hasIO = true
					ioTarget = target
				}
			}
		}
		return true
	})

	if hasIO {
		// Must have a verb prefix according to SideEffectFn[VerbPrefix]
		name := fn.Name.Name
		if !hasVerbPrefix(name) {
			pos := fset.Position(fn.Pos())
			*diags = append(*diags, Diagnostic{
				RuleID:     "PUR-001",
				Category:   CategoryEffectPurity,
				FilePath:   filePath,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    fmt.Sprintf("Function '%s' performs I/O side-effects (e.g. '%s') but lacks a standard verb prefix.", name, ioTarget),
				Severity:   SeverityInfo,
				Suggestion: "Use a descriptive verb prefix for effectful functions (e.g., fetch, send, execute, save, load, read, write).",
			})
		}
	}
}

func hasVerbPrefix(name string) bool {
	commonVerbs := []string{
		"get", "set", "fetch", "send", "load", "save", "write", "read",
		"exec", "execute", "run", "init", "sync", "ingest", "compact",
		"close", "open", "search", "connect", "query", "update", "delete",
		"insert", "parse", "extract", "check", "verify", "build", "format",
		"handle", "create", "new", "do", "log", "start", "stop", "post", "put",
	}

	lowerName := strings.ToLower(name)
	for _, v := range commonVerbs {
		if strings.HasPrefix(lowerName, v) {
			return true
		}
	}
	return false
}

// FunctionLengthRule checks fn[len>50L]: name_affix(Large∨Heavy∨Complex∨Batch∨Mega)
type FunctionLengthRule struct {
	MaxLines int
}

func (r *FunctionLengthRule) ID() string {
	return "LEN-001"
}

func (r *FunctionLengthRule) Category() Category {
	return CategoryNamingLength
}

func (r *FunctionLengthRule) Description() string {
	return "Detects excessively long functions (>60 lines) that should be decomposed or explicitly affixed"
}

func (r *FunctionLengthRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic
	maxL := r.MaxLines
	if maxL <= 0 {
		maxL = 60
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		startPos := fset.Position(fn.Pos())
		endPos := fset.Position(fn.End())
		lineCount := endPos.Line - startPos.Line + 1

		if lineCount > maxL {
			// Check if name has affix
			name := fn.Name.Name
			if !hasMegaAffix(name) && name != "main" {
				diagnostics = append(diagnostics, Diagnostic{
					RuleID:     r.ID(),
					Category:   r.Category(),
					FilePath:   filePath,
					Line:       startPos.Line,
					Column:     startPos.Column,
					Message:    fmt.Sprintf("Function '%s' is %d lines long (> %d lines). Consider refactoring into smaller pure morphisms.", name, lineCount, maxL),
					Severity:   SeverityInfo,
					Suggestion: "Decompose complex logic into modular sub-functions or pure domain transforms.",
				})
			}
		}
	}

	return diagnostics
}

func hasMegaAffix(name string) bool {
	affixes := []string{"Batch", "Mega", "Complex", "Heavy", "Pipeline", "Handler", "Main"}
	for _, a := range affixes {
		if strings.Contains(name, a) {
			return true
		}
	}
	return false
}
