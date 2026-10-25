package main

import (
	"go/ast"
	"go/token"
	"strings"
)

// NoErrorSwallowRule enforces ¬SilentSwallow ∧ ¬_ = err
type NoErrorSwallowRule struct{}

func (r *NoErrorSwallowRule) ID() string {
	return "ERR-001"
}

func (r *NoErrorSwallowRule) Category() Category {
	return CategoryErrorHandling
}

func (r *NoErrorSwallowRule) Description() string {
	return "Prohibits silencing errors via '_ = err' or empty error-handling blocks (¬SilentSwallow)"
}

func (r *NoErrorSwallowRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			// Detect _ = err
			for i, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Name != "_" {
					continue
				}
				if i < len(node.Rhs) {
					rhsIdent, isRhsIdent := node.Rhs[i].(*ast.Ident)
					if isRhsIdent && (rhsIdent.Name == "err" || strings.HasSuffix(rhsIdent.Name, "Err") || strings.HasPrefix(rhsIdent.Name, "err")) {
						pos := fset.Position(ident.Pos())
						diagnostics = append(diagnostics, Diagnostic{
							RuleID:     r.ID(),
							Category:   r.Category(),
							FilePath:   filePath,
							Line:       pos.Line,
							Column:     pos.Column,
							Message:    "Silent error discard detected: '_ = " + rhsIdent.Name + "'. Errors must be wrapped or handled.",
							Severity:   SeverityError,
							Suggestion: "Handle or return the error explicitly: if err != nil { return fmt.Errorf(\"context: %w\", err) }",
						})
					}
				}
			}

		case *ast.IfStmt:
			// Detect if err != nil {} with empty body
			if isErrorNilCheck(node.Cond) {
				if node.Body == nil || len(node.Body.List) == 0 {
					pos := fset.Position(node.Pos())
					diagnostics = append(diagnostics, Diagnostic{
						RuleID:     r.ID(),
						Category:   r.Category(),
						FilePath:   filePath,
						Line:       pos.Line,
						Column:     pos.Column,
						Message:    "Empty error handling block detected: 'if err != nil {}'. Error is silently swallowed.",
						Severity:   SeverityError,
						Suggestion: "Provide explicit error handling, wrapping, or logging within the block.",
					})
				}
			}
		}
		return true
	})

	return diagnostics
}

func isErrorNilCheck(expr ast.Expr) bool {
	binary, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	if binary.Op != token.NEQ {
		return false
	}

	isErr := func(e ast.Expr) bool {
		ident, ok := e.(*ast.Ident)
		return ok && (ident.Name == "err" || strings.HasSuffix(ident.Name, "Err") || strings.HasPrefix(ident.Name, "err"))
	}

	isNil := func(e ast.Expr) bool {
		ident, ok := e.(*ast.Ident)
		return ok && ident.Name == "nil"
	}

	return (isErr(binary.X) && isNil(binary.Y)) || (isNil(binary.X) && isErr(binary.Y))
}
