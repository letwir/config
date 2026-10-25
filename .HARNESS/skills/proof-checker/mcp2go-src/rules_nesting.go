package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// NestingFlowRule enforces EarlyReturn ∧ ¬NestedIfElse ∧ GuardClause(exit_fast)
type NestingFlowRule struct {
	MaxAllowedDepth int
}

func (r *NestingFlowRule) ID() string {
	return "NEST-001"
}

func (r *NestingFlowRule) Category() Category {
	return CategoryNesting
}

func (r *NestingFlowRule) Description() string {
	return "Enforces EarlyReturn, GuardClause, and maximum nesting depth (max depth <= 3)"
}

func (r *NestingFlowRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic
	maxDepth := r.MaxAllowedDepth
	if maxDepth <= 0 {
		maxDepth = 3
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		checkFuncNesting(fset, fn, filePath, maxDepth, &diagnostics)
	}

	return diagnostics
}

func checkFuncNesting(fset *token.FileSet, fn *ast.FuncDecl, filePath string, maxDepth int, diags *[]Diagnostic) {
	var walkBlock func(stmts []ast.Stmt, currentDepth int)

	walkBlock = func(stmts []ast.Stmt, currentDepth int) {
		for _, stmt := range stmts {
			switch s := stmt.(type) {
			case *ast.IfStmt:
				newDepth := currentDepth + 1
				if newDepth > maxDepth {
					pos := fset.Position(s.Pos())
					*diags = append(*diags, Diagnostic{
						RuleID:     "NEST-001",
						Category:   CategoryNesting,
						FilePath:   filePath,
						Line:       pos.Line,
						Column:     pos.Column,
						Message:    fmt.Sprintf("Excessive nesting depth (%d > %d) in function '%s'. Violates EarlyReturn rule.", newDepth, maxDepth, fn.Name.Name),
						Severity:   SeverityWarning,
						Suggestion: "Refactor with guard clauses (if cond { return }) or extract nested logic into separate pure functions.",
					})
				}
				if s.Body != nil {
					walkBlock(s.Body.List, newDepth)
				}
				if s.Else != nil {
					switch elseStmt := s.Else.(type) {
					case *ast.BlockStmt:
						walkBlock(elseStmt.List, newDepth)
					case *ast.IfStmt:
						walkBlock([]ast.Stmt{elseStmt}, currentDepth) // Else-if chain
					}
				}

			case *ast.ForStmt:
				newDepth := currentDepth + 1
				if newDepth > maxDepth {
					pos := fset.Position(s.Pos())
					*diags = append(*diags, Diagnostic{
						RuleID:     "NEST-001",
						Category:   CategoryNesting,
						FilePath:   filePath,
						Line:       pos.Line,
						Column:     pos.Column,
						Message:    fmt.Sprintf("Excessive nesting in loop (%d > %d) in function '%s'.", newDepth, maxDepth, fn.Name.Name),
						Severity:   SeverityWarning,
						Suggestion: "Use guard clauses (continue / break early) or extract inner loop processing.",
					})
				}
				if s.Body != nil {
					walkBlock(s.Body.List, newDepth)
				}

			case *ast.RangeStmt:
				newDepth := currentDepth + 1
				if newDepth > maxDepth {
					pos := fset.Position(s.Pos())
					*diags = append(*diags, Diagnostic{
						RuleID:     "NEST-001",
						Category:   CategoryNesting,
						FilePath:   filePath,
						Line:       pos.Line,
						Column:     pos.Column,
						Message:    fmt.Sprintf("Excessive nesting in range-loop (%d > %d) in function '%s'.", newDepth, maxDepth, fn.Name.Name),
						Severity:   SeverityWarning,
						Suggestion: "Use guard clauses or helper functions.",
					})
				}
				if s.Body != nil {
					walkBlock(s.Body.List, newDepth)
				}

			case *ast.SwitchStmt:
				if s.Body != nil {
					for _, cc := range s.Body.List {
						if clause, ok := cc.(*ast.CaseClause); ok {
							walkBlock(clause.Body, currentDepth+1)
						}
					}
				}

			case *ast.TypeSwitchStmt:
				if s.Body != nil {
					for _, cc := range s.Body.List {
						if clause, ok := cc.(*ast.CaseClause); ok {
							walkBlock(clause.Body, currentDepth+1)
						}
					}
				}

			case *ast.BlockStmt:
				walkBlock(s.List, currentDepth)
			}
		}
	}

	walkBlock(fn.Body.List, 0)
}
