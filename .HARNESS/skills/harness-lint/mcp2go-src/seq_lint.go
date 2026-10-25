package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	seqSectionRegex = regexp.MustCompile(`^[\x{1F300}-\x{1FAFF}\x{2600}-\x{26FF}\x{2700}-\x{27BF}\p{So}](?:\x{FE0F})?\s+([A-Za-z0-9_\-]+)(?:\s+@[a-zA-Z0-9_/\-]+)?\s*$`)
	seqIdentRegex   = regexp.MustCompile(`^[A-Za-z0-9_\-]+$`)
)

// lintSEQFile reads and performs DSL syntax linting on a .seq file.
func lintSEQFile(path string, maxLevel int) []Diagnostic {
	var diags []Diagnostic

	content, err := os.ReadFile(path)
	if err != nil {
		return append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:READ_ERROR",
			Severity: SeverityError,
			Message:  err.Error(),
			PIDGIN:   fmt.Sprintf("💥 [L0:READ_ERROR] %s ⊢ cannot read file: %v ⊥", path, err),
		})
	}

	if maxLevel >= 0 {
		diags = append(diags, CheckSEQ(path, content)...)
	}

	return diags
}

// CheckSEQ performs syntax integrity checks on raw Harness Sequence DSL (.seq) content.
// Pure function: (path, content) -> []Diagnostic
func CheckSEQ(path string, content []byte) []Diagnostic {
	var diags []Diagnostic

	// 1. Check for UTF-8 BOM and strip
	cleanContent := content
	if bytes.HasPrefix(content, []byte{0xEF, 0xBB, 0xBF}) {
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:UTF8_BOM",
			Severity: SeverityError,
			Message:  "UTF-8 BOM detected",
			PIDGIN:   fmt.Sprintf("💥 [L0:UTF8_BOM] %s ⊢ UTF-8 BOM detected! ⇒ BOMを除去して pure UTF-8-noBOM に変換せよ ⊥", path),
		})
		cleanContent = content[3:]
	}

	lines := strings.Split(strings.ReplaceAll(string(cleanContent), "\r\n", "\n"), "\n")

	// Empty file check
	if len(cleanContent) == 0 || (len(lines) == 1 && strings.TrimSpace(lines[0]) == "") {
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:SEQ_DIALECT_MISSING",
			Severity: SeverityError,
			Message:  "missing required first-line dialect marker @dialect harness-seq/1",
			PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_DIALECT_MISSING] %s:1 ⊢ missing first-line dialect marker! (expected: @dialect harness-seq/1) ⇒ ファイル先頭行に \"@dialect harness-seq/1\" を定義せよ ⊥", path),
		}, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:SEQ_VERSION_MISSING",
			Severity: SeverityError,
			Message:  "missing required @seq version directive (@seq 1 or @seq 2)",
			PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_VERSION_MISSING] %s:1 ⊢ @seq version directive missing! ⇒ \"@seq 1\" または \"@seq 2\" を定義せよ ⊥", path),
		}, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:SEQ_ENTRY_MISSING",
			Severity: SeverityError,
			Message:  "missing required @entry directive",
			PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_ENTRY_MISSING] %s:1 ⊢ @entry directive missing! ⇒ \"@entry <label>\" を定義せよ ⊥", path),
		})
		return diags
	}

	// 2. Validate first-line dialect marker
	firstLine := strings.TrimSpace(lines[0])
	if firstLine != "@dialect harness-seq/1" {
		if strings.HasPrefix(firstLine, "@dialect") {
			diags = append(diags, Diagnostic{
				Level:    Level0,
				File:     path,
				Line:     1,
				Col:      1,
				Code:     "L0:SEQ_DIALECT_INVALID",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid first-line dialect marker: %q (expected @dialect harness-seq/1)", firstLine),
				PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_DIALECT_INVALID] %s:1 ⊢ invalid first-line dialect marker %q (expected: @dialect harness-seq/1) ⇒ \"@dialect harness-seq/1\" に修正せよ ⊥", path, firstLine),
			})
		} else {
			diags = append(diags, Diagnostic{
				Level:    Level0,
				File:     path,
				Line:     1,
				Col:      1,
				Code:     "L0:SEQ_DIALECT_MISSING",
				Severity: SeverityError,
				Message:  "missing required first-line dialect marker @dialect harness-seq/1",
				PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_DIALECT_MISSING] %s:1 ⊢ missing first-line dialect marker! (expected: @dialect harness-seq/1) ⇒ ファイル先頭行に \"@dialect harness-seq/1\" を定義せよ ⊥", path),
			})
		}
	}

	seqFound := false
	entryFound := false
	entryValid := false
	entryLabel := ""
	entryLine := 1

	seenSections := make(map[string]int)

	for lineIdx, rawLine := range lines {
		lineNum := lineIdx + 1
		line := strings.TrimSpace(rawLine)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Subsequent duplicate @dialect check
		if lineIdx > 0 && strings.HasPrefix(line, "@dialect") {
			diags = append(diags, Diagnostic{
				Level:    Level0,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L0:SEQ_DIALECT_DUPLICATE",
				Severity: SeverityError,
				Message:  fmt.Sprintf("duplicate dialect marker: %s", line),
				PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_DIALECT_DUPLICATE] %s:%d ⊢ duplicate dialect marker detected ⇒ dialectマーカーは第1行に厳密に1つのみ定義せよ ⊥", path, lineNum),
			})
			continue
		}

		// First-line dialect marker itself requires no further directive processing
		if lineIdx == 0 && line == "@dialect harness-seq/1" {
			continue
		}

		// @seq version directive check
		if line == "@seq" || strings.HasPrefix(line, "@seq ") || strings.HasPrefix(line, "@seq\t") {
			if seqFound {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:SEQ_VERSION_DUPLICATE",
					Severity: SeverityError,
					Message:  fmt.Sprintf("duplicate @seq version directive: %s", line),
					PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_VERSION_DUPLICATE] %s:%d ⊢ duplicate @seq directive detected ⇒ @seq バージョン指示子は厳密に1つのみ定義せよ ⊥", path, lineNum),
				})
			} else {
				seqFound = true
				fields := strings.Fields(line)
				if len(fields) != 2 || fields[0] != "@seq" {
					diags = append(diags, Diagnostic{
						Level:    Level0,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L0:SEQ_VERSION_SYNTAX",
						Severity: SeverityError,
						Message:  fmt.Sprintf("invalid @seq directive syntax: %s (expected @seq 1 or @seq 2)", line),
						PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_VERSION_SYNTAX] %s:%d ⊢ invalid @seq directive syntax %q (expected: @seq 1 または @seq 2) ⇒ 構文を修正せよ ⊥", path, lineNum, line),
					})
				} else if fields[1] != "1" && fields[1] != "2" {
					diags = append(diags, Diagnostic{
						Level:    Level0,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L0:SEQ_VERSION_UNSUPPORTED",
						Severity: SeverityError,
						Message:  fmt.Sprintf("unsupported @seq version: %s (supported: 1, 2)", fields[1]),
						PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_VERSION_UNSUPPORTED] %s:%d ⊢ unsupported @seq version %q ⇒ バージョンは 1 または 2 を指定せよ ⊥", path, lineNum, fields[1]),
					})
				}
			}
			continue
		}

		// @entry directive check
		if line == "@entry" || strings.HasPrefix(line, "@entry ") || strings.HasPrefix(line, "@entry\t") {
			if entryFound {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:SEQ_ENTRY_DUPLICATE",
					Severity: SeverityError,
					Message:  fmt.Sprintf("duplicate @entry directive: %s", line),
					PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_ENTRY_DUPLICATE] %s:%d ⊢ duplicate @entry directive detected ⇒ @entry は厳密に1つのみ定義せよ ⊥", path, lineNum),
				})
			} else {
				entryFound = true
				fields := strings.Fields(line)
				if len(fields) != 2 || fields[0] != "@entry" || !seqIdentRegex.MatchString(fields[1]) {
					diags = append(diags, Diagnostic{
						Level:    Level0,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L0:SEQ_ENTRY_SYNTAX",
						Severity: SeverityError,
						Message:  fmt.Sprintf("invalid @entry directive syntax: %s (expected @entry <label>)", line),
						PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_ENTRY_SYNTAX] %s:%d ⊢ invalid @entry syntax %q ⇒ @entry <label> 形式で単一の識別子を指定せよ ⊥", path, lineNum, line),
					})
				} else {
					entryValid = true
					entryLabel = fields[1]
					entryLine = lineNum
				}
			}
			continue
		}

		// Section header check
		if matches := seqSectionRegex.FindStringSubmatch(line); len(matches) >= 2 {
			label := matches[1]
			if prevLine, exists := seenSections[label]; exists {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:SEQ_DUPLICATE_SECTION",
					Severity: SeverityError,
					Message:  fmt.Sprintf("duplicate section label: %q (first seen at line %d)", label, prevLine),
					PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_DUPLICATE_SECTION] %s:%d ⊢ duplicate section label %q (first seen at line %d) ⇒ セクションラベルは一意でなければならない ⊥", path, lineNum, label, prevLine),
				})
			} else {
				seenSections[label] = lineNum
			}
		}
	}

	// 3. Final presence and referential checks
	if !seqFound {
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:SEQ_VERSION_MISSING",
			Severity: SeverityError,
			Message:  "missing required @seq version directive (@seq 1 or @seq 2)",
			PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_VERSION_MISSING] %s:1 ⊢ @seq version directive missing! ⇒ \"@seq 1\" または \"@seq 2\" を定義せよ ⊥", path),
		})
	}

	if !entryFound {
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:SEQ_ENTRY_MISSING",
			Severity: SeverityError,
			Message:  "missing required @entry directive",
			PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_ENTRY_MISSING] %s:1 ⊢ @entry directive missing! ⇒ \"@entry <label>\" を定義せよ ⊥", path),
		})
	} else if entryValid {
		if _, exists := seenSections[entryLabel]; !exists {
			diags = append(diags, Diagnostic{
				Level:    Level0,
				File:     path,
				Line:     entryLine,
				Col:      1,
				Code:     "L0:SEQ_ENTRY_UNMATCHED",
				Severity: SeverityError,
				Message:  fmt.Sprintf("entry label %q does not match any section label in file", entryLabel),
				PIDGIN:   fmt.Sprintf("💥 [L0:SEQ_ENTRY_UNMATCHED] %s:%d ⊢ entry label %q does not match any section label in file ⇒ @entry はファイル内に存在するセクションラベルを指定せよ ⊥", path, entryLine, entryLabel),
			})
		}
	}

	return diags
}
