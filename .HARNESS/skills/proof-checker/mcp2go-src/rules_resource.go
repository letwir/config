package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ResourceCloseRule enforces Acquire ⇒ defer Resource.Close() ∧ RAII
type ResourceCloseRule struct{}

func (r *ResourceCloseRule) ID() string {
	return "RES-001"
}

func (r *ResourceCloseRule) Category() Category {
	return CategoryResourceRAII
}

func (r *ResourceCloseRule) Description() string {
	return "Verifies that acquired resources (os.Open, sql.Open, resp.Body, etc.) are properly closed with defer (RAII)"
}

func (r *ResourceCloseRule) CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic {
	var diagnostics []Diagnostic

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		checkResourceRAII(fset, fn, filePath, &diagnostics)
	}

	return diagnostics
}

func checkResourceRAII(fset *token.FileSet, fn *ast.FuncDecl, filePath string, diags *[]Diagnostic) {
	// Track acquired resource variables
	type resourceAcquire struct {
		varName string
		pos     token.Pos
		kind    string
	}

	var acquires []resourceAcquire
	closedVars := make(map[string]bool)

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for i, rhs := range stmt.Rhs {
				call, ok := rhs.(*ast.CallExpr)
				if !ok {
					continue
				}

				funcName := getCallTargetName(call)
				if isResourceAcquireCall(funcName) {
					if i < len(stmt.Lhs) {
						ident, ok := stmt.Lhs[i].(*ast.Ident)
						if ok && ident.Name != "_" {
							acquires = append(acquires, resourceAcquire{
								varName: ident.Name,
								pos:     ident.Pos(),
								kind:    funcName,
							})
						}
					}
				}
			}

		case *ast.DeferStmt:
			// Check what is being deferred
			ast.Inspect(stmt.Call, func(dn ast.Node) bool {
				if sel, ok := dn.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "Close" {
						if ident, ok := sel.X.(*ast.Ident); ok {
							closedVars[ident.Name] = true
						}
						// handle resp.Body.Close()
						if subSel, ok := sel.X.(*ast.SelectorExpr); ok {
							if subIdent, ok := subSel.X.(*ast.Ident); ok {
								closedVars[subIdent.Name] = true
							}
						}
					}
				}
				return true
			})
		}
		return true
	})

	for _, acq := range acquires {
		// Ignore if resource is returned by the function
		if isVariableReturned(fn, acq.varName) {
			continue
		}

		if !closedVars[acq.varName] {
			pos := fset.Position(acq.pos)
			*diags = append(*diags, Diagnostic{
				RuleID:     "RES-001",
				Category:   CategoryResourceRAII,
				FilePath:   filePath,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    fmt.Sprintf("Resource '%s' acquired via '%s' might not have a corresponding 'defer %s.Close()'.", acq.varName, acq.kind, acq.varName),
				Severity:   SeverityWarning,
				Suggestion: fmt.Sprintf("Add 'defer %s.Close()' immediately after successful acquisition and error check.", acq.varName),
			})
		}
	}
}

func getCallTargetName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		if pkg, ok := fun.X.(*ast.Ident); ok {
			return pkg.Name + "." + fun.Sel.Name
		}
		return fun.Sel.Name
	case *ast.Ident:
		return fun.Name
	}
	return ""
}

func isResourceAcquireCall(name string) bool {
	acquireCalls := []string{
		"os.Open",
		"os.OpenFile",
		"os.Create",
		"sql.Open",
		"zip.OpenReader",
		"zip.NewReader",
	}
	for _, ac := range acquireCalls {
		if name == ac {
			return true
		}
	}
	return false
}

func isVariableReturned(fn *ast.FuncDecl, varName string) bool {
	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			for _, res := range ret.Results {
				if ident, ok := res.(*ast.Ident); ok && ident.Name == varName {
					found = true
				}
				if sel, ok := res.(*ast.SelectorExpr); ok {
					if subIdent, ok := sel.X.(*ast.Ident); ok && subIdent.Name == varName {
						found = true
					}
				}
			}
		}
		return true
	})
	return found
}
