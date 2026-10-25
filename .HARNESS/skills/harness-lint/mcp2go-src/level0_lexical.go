package main

import (
	"bytes"
	"fmt"
	"strings"
)

// CheckLevel0 performs Level 0 lexical & syntax integrity checks on raw file content.
// Pure function: (path, content) -> []Diagnostic
func CheckLevel0(path string, content []byte) []Diagnostic {
	var diags []Diagnostic

	// 1. Check for UTF-8 BOM and strip for subsequent parsing
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
	headerFound := false
	recordFound := false
	firstNonEmptyNonComment := true

	for lineIdx, rawLine := range lines {
		lineNum := lineIdx + 1
		line := strings.TrimSpace(rawLine)

		if line == "" {
			continue
		}

		// Comment line
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Header line (@lrf=...)
		if strings.HasPrefix(line, "@lrf=") {
			if headerFound {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:HEADER_DUPLICATE",
					Severity: SeverityError,
					Message:  fmt.Sprintf("duplicate header: %s", line),
					PIDGIN:   fmt.Sprintf("💥 [L0:HEADER_DUPLICATE] %s:%d ⊢ duplicate header detected! ⇒ ヘッダーはファイル冒頭に厳密に1つのみ定義せよ ⊥", path, lineNum),
				})
				continue
			}

			if !firstNonEmptyNonComment || recordFound {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:HEADER_ORDER",
					Severity: SeverityError,
					Message:  "header must be the first non-empty non-comment line",
					PIDGIN:   fmt.Sprintf("💥 [L0:HEADER_ORDER] %s:%d ⊢ header must be the first line! ⇒ ヘッダーをレコードより前に配置せよ ⊥", path, lineNum),
				})
			}

			headerFound = true
			firstNonEmptyNonComment = false

			fields := splitUnescapedPipes(line)
			if len(fields) != 3 || fields[0] != "@lrf=1" || !strings.HasPrefix(fields[1], "aud=") || len(fields[1]) <= 4 || !strings.HasPrefix(fields[2], "scope=") || len(fields[2]) <= 6 {
				diags = append(diags, Diagnostic{
					Level:    Level0,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L0:HEADER_SYNTAX",
					Severity: SeverityError,
					Message:  fmt.Sprintf("invalid header syntax: %s", line),
					PIDGIN:   fmt.Sprintf("💥 [L0:HEADER_SYNTAX] %s:%d ⊢ header \"%s\" invalid! (expected: @lrf=1|aud=<aud>|scope=<scope>) ⇒ 厳密に3カラムのヘッダー構文を修復せよ ⊥", path, lineNum, line),
				})
			}
			continue
		}

		firstNonEmptyNonComment = false

		// Record line checks (R| or L|)
		if strings.HasPrefix(line, "R|") || strings.HasPrefix(line, "L|") {
			recordFound = true
			fields := splitUnescapedPipes(line)
			if len(fields) != 6 {
				if len(fields) > 6 {
					diags = append(diags, Diagnostic{
						Level:    Level0,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L0:PIPE_COLLISION",
						Severity: SeverityError,
						Message:  fmt.Sprintf("delimiter collision: %d fields detected (expected 6)", len(fields)),
						PIDGIN:   fmt.Sprintf("💥 [L0:PIPE_COLLISION] %s:%d ⊢ raw \"|\" delimiter collision detected! fields=%d ≠ 6 ⇒ payload/guard内のパイプは \"\\|\" でエスケープするか除去せよ ⊥", path, lineNum, len(fields)),
					})
				} else {
					diags = append(diags, Diagnostic{
						Level:    Level0,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L0:INVALID_COLUMN_COUNT",
						Severity: SeverityError,
						Message:  fmt.Sprintf("invalid column count: %d (expected 6)", len(fields)),
						PIDGIN:   fmt.Sprintf("💥 [L0:INVALID_COLUMN_COUNT] %s:%d ⊢ field count %d ≠ 6 (expected: R|id|guard|modality|effect|payload) ⇒ カラム数を厳密に6個に修正せよ ⊥", path, lineNum, len(fields)),
					})
				}
			}
			continue
		}

		// Check for malformed record prefix (e.g., M|, X|, r|, l|, Rule|)
		if strings.Contains(line, "|") {
			prefix := strings.SplitN(line, "|", 2)[0]
			diags = append(diags, Diagnostic{
				Level:    Level0,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L0:RECORD_PREFIX",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid record prefix: %s", prefix),
				PIDGIN:   fmt.Sprintf("💥 [L0:RECORD_PREFIX] %s:%d ⊢ invalid record prefix \"%s\" ⇒ 行頭は \"R|\" または \"L|\" でなければならない ⊥", path, lineNum, prefix),
			})
			continue
		}

		// Any other prose, markdown bullet, or unknown syntax
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     lineNum,
			Col:      1,
			Code:     "L0:UNKNOWN_SYNTAX",
			Severity: SeverityError,
			Message:  fmt.Sprintf("unknown syntax line: %s", line),
			PIDGIN:   fmt.Sprintf("💥 [L0:UNKNOWN_SYNTAX] %s:%d ⊢ unknown syntax \"%s\" ⇒ 自然言語散文・不正構文は禁止 ⊥ (LRF規格行のみ許可)", path, lineNum, line),
		})
	}

	// Header required check
	if !headerFound {
		diags = append(diags, Diagnostic{
			Level:    Level0,
			File:     path,
			Line:     1,
			Col:      1,
			Code:     "L0:MISSING_HEADER",
			Severity: SeverityError,
			Message:  "missing required @lrf=1 header",
			PIDGIN:   fmt.Sprintf("💥 [L0:MISSING_HEADER] %s ⊢ header missing! ⇒ 冒頭に \"@lrf=1|aud=<aud>|scope=<scope>\" を定義せよ ⊥", path),
		})
	}

	return diags
}

// splitUnescapedPipes splits a line by raw '|' not preceded by '\'.
func splitUnescapedPipes(s string) []string {
	var fields []string
	var current strings.Builder
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			current.WriteByte(c)
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			current.WriteByte(c)
			continue
		}
		if c == '|' {
			fields = append(fields, current.String())
			current.Reset()
			continue
		}
		current.WriteByte(c)
	}
	fields = append(fields, current.String())
	return fields
}
