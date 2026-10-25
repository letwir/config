package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

var (
	// PIDGIN command pattern in SKILL.md: emoji + optional JA label + command name
	pidginIndexRegex = regexp.MustCompile(`[\x{1F300}-\x{1FAFF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}](?:\x{FE0F})?(?:\s*[^\x00-\x7F\s]+)?\s*([a-zA-Z0-9_\-]+)`)
	pidginLineRegex  = regexp.MustCompile(`(?:^|\n)[\x{1F300}-\x{1FAFF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}](?:\x{FE0F})?(?:\s*[^\x00-\x7F\s]+)?\s*([a-zA-Z0-9_\-]+)\s*⊢`)
	usageCmdRegex    = regexp.MustCompile(`(?:\.\/|\.\\|[a-zA-Z0-9_\-]+\.exe\s+)([a-zA-Z0-9_\-]+)`)
)

// ExtractCommandsFromGoAST extracts subcommand names from Go source code AST.
// Pure function: (filename, src) -> ([]string, error)
func ExtractCommandsFromGoAST(filename string, src []byte) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	cmdMap := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		// Detect switch statements: switch os.Args[1], switch subcmd, switch cmd
		case *ast.SwitchStmt:
			for _, stmt := range x.Body.List {
				if cc, ok := stmt.(*ast.CaseClause); ok {
					for _, expr := range cc.List {
						if bl, ok := expr.(*ast.BasicLit); ok && bl.Kind == token.STRING {
							val := strings.Trim(bl.Value, `"`+"`")
							// Ignore flags or empty
							if val != "" && !strings.HasPrefix(val, "-") {
								cmdMap[val] = true
							}
						}
					}
				}
			}

		// Detect flag.NewFlagSet("subcommand", ...) and flag definitions
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				// Detect flag.NewFlagSet("subcommand", ...)
				if sel.Sel.Name == "NewFlagSet" && len(x.Args) >= 1 {
					if bl, ok := x.Args[0].(*ast.BasicLit); ok && bl.Kind == token.STRING {
						val := strings.Trim(bl.Value, `"`+"`")
						if val != "" && !strings.HasPrefix(val, "-") {
							cmdMap[val] = true
						}
					}
				}

				// Detect flag.String("name", ...), flag.Int("name", ...), flag.Bool("name", ...)
				if (sel.Sel.Name == "String" || sel.Sel.Name == "Int" || sel.Sel.Name == "Bool" ||
					sel.Sel.Name == "StringVar" || sel.Sel.Name == "IntVar" || sel.Sel.Name == "BoolVar") && len(x.Args) >= 1 {
					argIdx := 0
					if strings.HasSuffix(sel.Sel.Name, "Var") {
						argIdx = 1
					}
					if len(x.Args) > argIdx {
						if bl, ok := x.Args[argIdx].(*ast.BasicLit); ok && bl.Kind == token.STRING {
							val := strings.Trim(bl.Value, `"`+"`")
							if val != "" && !strings.HasPrefix(val, "-") {
								cmdMap[val] = true
							}
						}
					}
				}
			}
		}
		return true
	})

	var cmds []string
	for cmd := range cmdMap {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	return cmds, nil
}

// ExtractCommandsFromPIDGIN extracts documented CLI commands from SKILL.md PIDGIN sections.
func ExtractCommandsFromPIDGIN(content []byte) []string {
	text := string(content)
	cmdMap := make(map[string]bool)

	// 1. Check line 5 index matches (emoji+command)
	matches := pidginIndexRegex.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			cmd := strings.TrimSpace(m[1])
			if isValidCommandName(cmd) {
				cmdMap[cmd] = true
			}
		}
	}

	// 2. Check lines 6-10 definition matches (emoji command ⊢)
	lineMatches := pidginLineRegex.FindAllStringSubmatch(text, -1)
	for _, m := range lineMatches {
		if len(m) >= 2 {
			cmd := strings.TrimSpace(m[1])
			if isValidCommandName(cmd) {
				cmdMap[cmd] = true
			}
		}
	}

	var cmds []string
	for cmd := range cmdMap {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	return cmds
}

// ExtractCommandsFromUSAGE extracts commands demonstrated in USAGE.md.
func ExtractCommandsFromUSAGE(content []byte) []string {
	text := string(content)
	cmdMap := make(map[string]bool)

	matches := usageCmdRegex.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			cmd := strings.TrimSpace(m[1])
			if isValidCommandName(cmd) {
				cmdMap[cmd] = true
			}
		}
	}

	var cmds []string
	for cmd := range cmdMap {
		cmds = append(cmds, cmd)
	}
	sort.Strings(cmds)
	return cmds
}

func isValidCommandName(s string) bool {
	if s == "" || strings.HasPrefix(s, "-") {
		return false
	}
	// Common non-command words
	ignored := map[string]bool{
		"help": true, "version": true, "exe": true, "go": true, "pwsh": true,
	}
	return !ignored[strings.ToLower(s)]
}

// CheckLevel3Diff verifies bi-simulation between Go AST, PIDGIN documentation, and USAGE.
// Pure function: computes symmetric differences and emits compiler diagnostics.
func CheckLevel3Diff(targetName string, astCmds, pidginCmds, usageCmds []string) []Diagnostic {
	var diags []Diagnostic

	astSet := toSet(astCmds)
	pidginSet := toSet(pidginCmds)
	usageSet := toSet(usageCmds)

	// 1. Doc Drift: Implemented in AST but missing from PIDGIN
	var missingInPIDGIN []string
	for cmd := range astSet {
		if !pidginSet[cmd] {
			missingInPIDGIN = append(missingInPIDGIN, cmd)
		}
	}
	sort.Strings(missingInPIDGIN)

	if len(missingInPIDGIN) > 0 {
		diffStr := strings.Join(missingInPIDGIN, ", ")
		diags = append(diags, Diagnostic{
			Level:    Level3,
			File:     targetName,
			Line:     1,
			Col:      1,
			Code:     "L3:DOC_DRIFT",
			Severity: SeverityError,
			Message:  fmt.Sprintf("commands in AST missing from PIDGIN contract: {%s}", diffStr),
			PIDGIN:   fmt.Sprintf("⚖️ [L3:DOC_DRIFT] %s ⊢ C_code ∖ C_pidgin = {%s} ≠ ∅ ⇒ CLIに実装済みなのにSKILL.md(PIDGIN)に未記載! Append to PIDGIN contract ⊥", targetName, diffStr),
		})
	}

	// 2. Dead Doc: Documented in PIDGIN but not implemented in AST
	var deadInPIDGIN []string
	for cmd := range pidginSet {
		if !astSet[cmd] {
			deadInPIDGIN = append(deadInPIDGIN, cmd)
		}
	}
	sort.Strings(deadInPIDGIN)

	if len(deadInPIDGIN) > 0 {
		diffStr := strings.Join(deadInPIDGIN, ", ")
		diags = append(diags, Diagnostic{
			Level:    Level3,
			File:     targetName,
			Line:     1,
			Col:      1,
			Code:     "L3:DEAD_DOC",
			Severity: SeverityError,
			Message:  fmt.Sprintf("commands in PIDGIN missing from Go AST: {%s}", diffStr),
			PIDGIN:   fmt.Sprintf("⚖️ [L3:DEAD_DOC] %s ⊢ C_pidgin ∖ C_code = {%s} ≠ ∅ ⇒ SKILL.mdに記載があるがGoソースコードに未実装! 実装するかドキュメントを削除せよ ⊥", targetName, diffStr),
		})
	}

	// 3. USAGE Drift: Implemented in AST but missing from USAGE.md (Warning)
	if len(usageCmds) > 0 {
		var missingInUsage []string
		for cmd := range astSet {
			if !usageSet[cmd] {
				missingInUsage = append(missingInUsage, cmd)
			}
		}
		sort.Strings(missingInUsage)

		if len(missingInUsage) > 0 {
			diffStr := strings.Join(missingInUsage, ", ")
			diags = append(diags, Diagnostic{
				Level:    Level3,
				File:     targetName,
				Line:     1,
				Col:      1,
				Code:     "L3:USAGE_DRIFT",
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("commands in AST missing from USAGE.md: {%s}", diffStr),
				PIDGIN:   fmt.Sprintf("⚖️ [L3:USAGE_DRIFT] %s ⊢ C_code ∖ C_usage = {%s} ≠ ∅ ⇒ CLIに実装済みなのにUSAGE.mdに記載なし! Update USAGE.md ⊥", targetName, diffStr),
			})
		}
	}

	return diags
}

func toSet(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range slice {
		if s != "" {
			m[s] = true
		}
	}
	return m
}
