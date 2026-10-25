package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		targetPath = flag.String("path", ".", "Target file or directory to lint")
		maxLevel   = flag.Int("level", 3, "Maximum validation level (0, 1, 2, or 3)")
		jsonOutput = flag.Bool("json", false, "Output diagnostics in JSON format")
		strict     = flag.Bool("strict", false, "Treat warnings as errors")
		skillDir   = flag.String("skill", "", "Path to skill directory for Level 3 AST vs PIDGIN check")
	)
	flag.Parse()

	var allDiags []Diagnostic
	scannedFiles := 0

	// 1. Lint LRF and SEQ files (Level 0, Level 1, Level 2)
	fi, err := os.Stat(*targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 [FATAL] Path not found: %s\n", *targetPath)
		os.Exit(1)
	}

	if fi.IsDir() {
		err = filepath.WalkDir(*targetPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				name := strings.ToLower(d.Name())
				if strings.HasSuffix(name, ".lrf") {
					diags := lintLRFFile(path, *maxLevel)
					allDiags = append(allDiags, diags...)
					scannedFiles++
				} else if strings.HasSuffix(name, ".seq") {
					diags := lintSEQFile(path, *maxLevel)
					allDiags = append(allDiags, diags...)
					scannedFiles++
				}
			}
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "💥 [FATAL] Walk error: %v\n", err)
			os.Exit(1)
		}
	} else {
		name := strings.ToLower(fi.Name())
		if strings.HasSuffix(name, ".lrf") {
			allDiags = append(allDiags, lintLRFFile(*targetPath, *maxLevel)...)
			scannedFiles++
		} else if strings.HasSuffix(name, ".seq") {
			allDiags = append(allDiags, lintSEQFile(*targetPath, *maxLevel)...)
			scannedFiles++
		}
	}

	// 2. Level 3 AST vs PIDGIN check if skillDir or applicable target
	if *maxLevel >= 3 {
		targetSkill := *skillDir
		if targetSkill == "" && fi.IsDir() {
			// Check if targetPath itself contains SKILL.md
			if _, err := os.Stat(filepath.Join(*targetPath, "SKILL.md")); err == nil {
				targetSkill = *targetPath
			}
		}

		if targetSkill != "" {
			l3Diags := lintSkillLevel3(targetSkill)
			allDiags = append(allDiags, l3Diags...)
			scannedFiles++
		}
	}

	// 3. Output results
	errCount := 0
	warnCount := 0
	for _, d := range allDiags {
		if d.Severity == SeverityError {
			errCount++
		} else if d.Severity == SeverityWarning {
			warnCount++
		}
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(allDiags)
	} else {
		if len(allDiags) > 0 {
			fmt.Println("────────────────────────────────────────────────────────────")
			fmt.Println("📋 HARNESS COMPILER DIAGNOSTICS")
			fmt.Println("────────────────────────────────────────────────────────────")
			for _, d := range allDiags {
				fmt.Println(d.String())
			}
			fmt.Println("────────────────────────────────────────────────────────────")
		}

		fmt.Printf("🎯 [SUMMARY] Scanned: %d targets | Errors: %d | Warnings: %d | MaxLevel: %d\n",
			scannedFiles, errCount, warnCount, *maxLevel)

		if errCount == 0 && (!*strict || warnCount == 0) {
			fmt.Println("✨ [ALL_PASS] Γ ⊢ All Invariants Satisfied (Zero-Defect Soundness) ⊣ ⊤")
		} else {
			fmt.Println("💥 [DIAGNOSTIC_FAILURE] Invariants violated! Fix issues indicated by PIDGIN diagnostics ⊥")
		}
	}

	if errCount > 0 || (*strict && warnCount > 0) {
		os.Exit(1)
	}
}

func lintLRFFile(path string, maxLevel int) []Diagnostic {
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

	// Level 0: Lexical
	l0Diags := CheckLevel0(path, content)
	diags = append(diags, l0Diags...)

	if maxLevel < 1 {
		return diags
	}

	// Level 1: Structural Parsing & Typing
	records, l1Diags := ParseLevel1(path, content)
	diags = append(diags, l1Diags...)

	if maxLevel < 2 {
		return diags
	}

	// Level 2: Semantic & Referential
	baseDir := filepath.Dir(path)
	l2Diags := CheckLevel2(path, records, baseDir, RealFS{})
	diags = append(diags, l2Diags...)

	return diags
}

func lintSkillLevel3(skillDir string) []Diagnostic {
	var diags []Diagnostic
	skillName := filepath.Base(filepath.Clean(skillDir))

	// 1. Resolve SKILL.md location (handle separated structures)
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if _, err := os.Stat(skillPath); err != nil {
		candidatePaths := []string{
			filepath.Join("LTWskills", "skills", skillName, "SKILL.md"),
			filepath.Join("src", skillName, "SKILL.md"),
			filepath.Join("..", "LTWskills", "skills", skillName, "SKILL.md"),
			filepath.Join("C:", "Users", "letwir", ".harness", "skills", skillName, "SKILL.md"),
		}
		found := false
		for _, cp := range candidatePaths {
			if _, err := os.Stat(cp); err == nil {
				skillPath = cp
				found = true
				break
			}
		}
		if !found {
			diags = append(diags, Diagnostic{
				Level:    Level3,
				File:     skillDir,
				Line:     1,
				Col:      1,
				Code:     "L3:MISSING_SKILL_DOC",
				Severity: SeverityError,
				Message:  fmt.Sprintf("SKILL.md not found in %s or candidate paths", skillDir),
				PIDGIN:   fmt.Sprintf("💥 [L3:MISSING_SKILL_DOC] %s ⊢ SKILL.md not found! ⇒ SKILL.md を配置せよ ⊥", skillDir),
			})
			return diags
		}
	}

	skillBytes, err := os.ReadFile(skillPath)
	if err != nil {
		diags = append(diags, Diagnostic{
			Level:    Level3,
			File:     skillPath,
			Line:     1,
			Col:      1,
			Code:     "L3:READ_ERROR",
			Severity: SeverityError,
			Message:  err.Error(),
			PIDGIN:   fmt.Sprintf("💥 [L3:READ_ERROR] %s ⊢ cannot read SKILL.md: %v ⊥", skillPath, err),
		})
		return diags
	}

	pidginCmds := ExtractCommandsFromPIDGIN(skillBytes)

	// 2. Resolve Go source directory
	var searchDirs []string
	searchDirs = append(searchDirs, skillDir)
	srcCandidate := filepath.Join("src", skillName)
	if srcCandidate != skillDir {
		if _, err := os.Stat(srcCandidate); err == nil {
			searchDirs = append(searchDirs, srcCandidate)
		}
	}
	// Also check if skillDir is inside LTWskills and src is in repo root
	repoSrcCandidate := filepath.Join("..", "..", "src", skillName)
	if _, err := os.Stat(repoSrcCandidate); err == nil {
		searchDirs = append(searchDirs, repoSrcCandidate)
	}

	var astCmds []string
	goFilesFound := 0
	for _, sDir := range searchDirs {
		_ = filepath.WalkDir(sDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") && !strings.HasSuffix(d.Name(), "_test.go") {
				goFilesFound++
				goBytes, err := os.ReadFile(path)
				if err == nil {
					cmds, err := ExtractCommandsFromGoAST(path, goBytes)
					if err == nil {
						astCmds = append(astCmds, cmds...)
					}
				}
			}
			return nil
		})
	}

	if goFilesFound == 0 {
		diags = append(diags, Diagnostic{
			Level:    Level3,
			File:     skillDir,
			Line:     1,
			Col:      1,
			Code:     "L3:MISSING_GO_SRC",
			Severity: SeverityError,
			Message:  fmt.Sprintf("no Go source files found in %s or src/%s", skillDir, skillName),
			PIDGIN:   fmt.Sprintf("💥 [L3:MISSING_GO_SRC] %s ⊢ Go source code not found in %s or src/%s! ⇒ Go実装ファイルを配置せよ ⊥", skillDir, skillDir, skillName),
		})
	}

	// 3. Resolve USAGE.md
	var usageCmds []string
	usagePath := filepath.Join(filepath.Dir(skillPath), "USAGE.md")
	if usageBytes, err := os.ReadFile(usagePath); err == nil {
		usageCmds = ExtractCommandsFromUSAGE(usageBytes)
	}

	// 4. Perform Level 3 Diff Assertion: Run if EITHER exists
	if len(astCmds) > 0 || len(pidginCmds) > 0 {
		l3Diags := CheckLevel3Diff(skillName, astCmds, pidginCmds, usageCmds)
		diags = append(diags, l3Diags...)
	}

	return diags
}
