package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileSystem provides an abstraction over file system inspection for pure/mockable testing.
type FileSystem interface {
	Stat(path string) (os.FileInfo, error)
}

// RealFS implements FileSystem using standard os.Stat.
type RealFS struct{}

func (RealFS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

var (
	linkExtractRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	mutationKeywords = []string{
		"write", "mutate", "delete", "append", "edit", "patch", "create", "commit", "destruct", "overwrite", "rewrite", "repair",
	}
)

// CheckLevel2 performs Level 2 semantic, referential integrity, and effect soundness checks.
func CheckLevel2(path string, records []LRFRecord, baseDir string, fs FileSystem) []Diagnostic {
	var diags []Diagnostic
	if fs == nil {
		fs = RealFS{}
	}

	seenIDs := make(map[string]int)

	for _, rec := range records {
		lineNum := rec.LineNum

		// 1. ID Uniqueness check (Verifier Advisory)
		if prevLine, exists := seenIDs[rec.ID]; exists {
			diags = append(diags, Diagnostic{
				Level:    Level2,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L2:DUPLICATE_ID",
				Severity: SeverityError,
				Message:  fmt.Sprintf("duplicate record id: \"%s\" (first seen at line %d)", rec.ID, prevLine),
				PIDGIN:   fmt.Sprintf("🆔 [L2:DUPLICATE_ID] %s:%d ⊢ record id=\"%s\" is duplicate (first seen at line %d) ⇒ 一意なID (Unique ID) を付与せよ ⊥", path, lineNum, rec.ID, prevLine),
			})
		} else {
			seenIDs[rec.ID] = lineNum
		}

		// 2. Link target existence check for L records
		if rec.Type == "L" {
			matches := linkExtractRegex.FindStringSubmatch(rec.Payload)
			if len(matches) == 3 {
				target := matches[2]
				// Resolve relative target against baseDir
				var resolvedPath string
				if filepath.IsAbs(target) {
					resolvedPath = target
				} else {
					resolvedPath = filepath.Join(baseDir, target)
				}

				if _, err := fs.Stat(resolvedPath); err != nil {
					diags = append(diags, Diagnostic{
						Level:    Level2,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L2:BROKEN_LINK",
						Severity: SeverityError,
						Message:  fmt.Sprintf("linked target file not found: \"%s\" (resolved: \"%s\")", target, resolvedPath),
						PIDGIN:   fmt.Sprintf("🔗 [L2:BROKEN_LINK] %s:%d ⊢ link target=\"%s\" (resolved=\"%s\") ⇒ file not found on disk! ∃! relative target required ⊥", path, lineNum, target, resolvedPath),
					})
				}
			}
		}

		// 3. Effect Soundness check: Detect mutation payload paired with read-only effect
		isReadOnlyEffect := (rec.Effect == "RO_LOCAL" || rec.Effect == "RO_PUBLIC")
		if isReadOnlyEffect {
			lowerPayload := strings.ToLower(rec.Payload)
			for _, kw := range mutationKeywords {
				// Match keyword as a word or key=val pattern (e.g. "write=", "append=", "create=")
				if containsMutationWord(lowerPayload, kw) {
					// Check if action is negated (e.g., "repair=false", "no-write", "no patch", "no rule patch", "never ...")
					if strings.Contains(lowerPayload, kw+"=false") ||
						strings.Contains(lowerPayload, "no-"+kw) ||
						strings.Contains(lowerPayload, "no "+kw) ||
						strings.Contains(lowerPayload, "no rule "+kw) ||
						strings.Contains(lowerPayload, "never "+kw) ||
						strings.Contains(lowerPayload, "without "+kw) {
						continue
					}
					diags = append(diags, Diagnostic{
						Level:    Level2,
						File:     path,
						Line:     lineNum,
						Col:      1,
						Code:     "L2:EFFECT_MISMATCH",
						Severity: SeverityError,
						Message:  fmt.Sprintf("mutation action \"%s\" found in payload but effect is read-only \"%s\"", kw, rec.Effect),
						PIDGIN:   fmt.Sprintf("⚡ [L2:EFFECT_MISMATCH] %s:%d ⊢ action \"%s\" in payload implies mutation, but effect=\"%s\" is read-only! ⇒ effect を LW_SCOPE または LIVE_WRITE に昇格せよ ⊥", path, lineNum, kw, rec.Effect),
					})
					break
				}
			}
		}
	}

	return diags
}

// containsMutationWord checks if keyword is present as a meaningful token, avoiding substring false positives.
func containsMutationWord(payload, kw string) bool {
	patterns := []string{
		kw + "=",
		kw + ":",
		" " + kw + " ",
		" " + kw + ";",
		" " + kw + ",",
	}
	for _, p := range patterns {
		if strings.Contains(payload, p) {
			return true
		}
	}
	// Check prefix / suffix
	if strings.HasPrefix(payload, kw+" ") || strings.HasPrefix(payload, kw+"=") || strings.HasPrefix(payload, kw+":") {
		return true
	}
	if strings.HasSuffix(payload, " "+kw) {
		return true
	}
	return false
}
