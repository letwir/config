package main

import (
	"strings"
	"testing"
)

func TestExtractCommandsFromGoAST(t *testing.T) {
	goCode := `package main

import (
	"flag"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	switch os.Args[1] {
	case "show":
		// show impl
	case "list":
		// list impl
	case "submit":
		// submit impl
	default:
		flag.NewFlagSet("export", flag.ExitOnError)
	}
}
`
	cmds, err := ExtractCommandsFromGoAST("main.go", []byte(goCode))
	if err != nil {
		t.Fatalf("unexpected error parsing Go AST: %v", err)
	}

	expected := map[string]bool{"show": true, "list": true, "submit": true, "export": true}
	if len(cmds) != 4 {
		t.Fatalf("expected 4 commands, got %d: %v", len(cmds), cmds)
	}
	for _, c := range cmds {
		if !expected[c] {
			t.Errorf("unexpected command extracted: %s", c)
		}
	}
}

func TestExtractCommandsFromPIDGIN(t *testing.T) {
	skillContent := `---
name: llm-memory
description: test
---
[SPR/PIDGIN::ρ→max] 🔍show ∘ 📊list ∘ 📥submit
🔍 show ⊢ -limit N ⇒ 直近N件を表示
📊 list ⊢ -json ⇒ 全件リスト表示
📥 submit ⊢ -file F ⇒ アンケート提出
`
	cmds := ExtractCommandsFromPIDGIN([]byte(skillContent))
	expected := map[string]bool{"show": true, "list": true, "submit": true}

	if len(cmds) != 3 {
		t.Fatalf("expected 3 commands from PIDGIN, got %d: %v", len(cmds), cmds)
	}
	for _, c := range cmds {
		if !expected[c] {
			t.Errorf("unexpected command in PIDGIN: %s", c)
		}
	}
}

func TestCheckLevel3Diff_Matched(t *testing.T) {
	astCmds := []string{"list", "show", "submit"}
	pidginCmds := []string{"list", "show", "submit"}
	usageCmds := []string{"list", "show", "submit"}

	diags := CheckLevel3Diff("llm-memory", astCmds, pidginCmds, usageCmds)
	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %v", len(diags), diags)
	}
}

func TestCheckLevel3Diff_SurveyShowOmission(t *testing.T) {
	// Replicating: "survey show を実装したのにSKILL.mdへ書き忘れた"
	astCmds := []string{"list", "show", "submit"}
	pidginCmds := []string{"list", "submit"} // "show" omitted!
	usageCmds := []string{"list", "submit"}

	diags := CheckLevel3Diff("llm-memory", astCmds, pidginCmds, usageCmds)

	foundDocDrift := false
	for _, d := range diags {
		if d.Code == "L3:DOC_DRIFT" {
			foundDocDrift = true
			if !strings.Contains(d.PIDGIN, "C_code ∖ C_pidgin = {show}") {
				t.Errorf("expected missing command 'show' in diff, got: %s", d.PIDGIN)
			}
			if !strings.Contains(d.PIDGIN, "⚖️") || !strings.Contains(d.PIDGIN, "⊥") {
				t.Errorf("expected PIDGIN symbols, got: %s", d.PIDGIN)
			}
		}
	}
	if !foundDocDrift {
		t.Fatalf("expected L3:DOC_DRIFT diagnostic for missing 'show', got: %v", diags)
	}
}

func TestCheckLevel3Diff_DeadDoc(t *testing.T) {
	// Documented in PIDGIN but not implemented in AST
	astCmds := []string{"list"}
	pidginCmds := []string{"list", "ghost_cmd"}
	usageCmds := []string{"list"}

	diags := CheckLevel3Diff("test-skill", astCmds, pidginCmds, usageCmds)

	foundDeadDoc := false
	for _, d := range diags {
		if d.Code == "L3:DEAD_DOC" {
			foundDeadDoc = true
			if !strings.Contains(d.PIDGIN, "C_pidgin ∖ C_code = {ghost_cmd}") {
				t.Errorf("expected dead command 'ghost_cmd' in diff, got: %s", d.PIDGIN)
			}
		}
	}
	if !foundDeadDoc {
		t.Fatalf("expected L3:DEAD_DOC diagnostic for dead 'ghost_cmd', got: %v", diags)
	}
}
