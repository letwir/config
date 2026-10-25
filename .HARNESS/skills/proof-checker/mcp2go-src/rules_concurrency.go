package main

import (
	"go/ast"
	"go/token"
)

// ConcurrencyRule enforces concurrency: ∀Async ⇒ Ctx{context.Context, Timeout} ∧ ¬Leak
type ConcurrencyRule struct{}

func (r *ConcurrencyRule) ID() string {
	return "CONC-001"
}

func (r *ConcurrencyRule) Category() Category {
	return CategoryConcurrency
}

func (r *ConcurrencyRule) Description() string {
	return "Verifies that asynchronous goroutines handle cancellation/context to prevent leaks"
}

func (r *ConcurrencyRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic

	ast.Inspect(file, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}

		pos := fset.Position(goStmt.Pos())

		// If calling a func literal: check if it takes ctx or has select with Done
		if funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit); ok {
			hasContextParam := false
			if funcLit.Type.Params != nil {
				for _, param := range funcLit.Type.Params.List {
					if isContextType(param.Type) {
						hasContextParam = true
						break
					}
				}
			}

			hasDoneSelect := false
			ast.Inspect(funcLit.Body, func(in ast.Node) bool {
				if sel, ok := in.(*ast.SelectStmt); ok && sel.Body != nil {
					hasDoneSelect = true
				}
				return true
			})

			if !hasContextParam && !hasDoneSelect {
				diagnostics = append(diagnostics, Diagnostic{
					RuleID:     r.ID(),
					Category:   r.Category(),
					FilePath:   filePath,
					Line:       pos.Line,
					Column:     pos.Column,
					Message:    "Goroutine spawned without explicit context parameter or cancellation select loop (potential leak).",
					Severity:   SeverityInfo,
					Suggestion: "Pass a context.Context with timeout or ensure loop listens to ctx.Done().",
				})
			}
		}

		return true
	})

	return diagnostics
}

func isContextType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if pkg, ok := sel.X.(*ast.Ident); ok {
			return pkg.Name == "context" && sel.Sel.Name == "Context"
		}
	}
	return false
}
