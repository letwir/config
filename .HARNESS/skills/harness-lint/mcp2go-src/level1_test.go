package main

import (
	"strings"
	"testing"
)

func TestParseLevel1_Valid(t *testing.T) {
	content := `@lrf=1|aud=GPT-5.6|scope=test
R|intent.ro|task:read-only|MUST|RO_LOCAL|answer/diagnose/review=inspect+evidence; repair=false
L|load.diary|*|MUST|RO_LOCAL|[diary](DIARY.lrf)
R|multi.guard|task:code&tag:subagent|MAY|LW_SCOPE|composite guard payload
R|cred.check|event:auth|MUST_NOT|CRED|never print credentials
`
	records, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %v", len(diags), diags)
	}
	if len(records) != 4 {
		t.Fatalf("expected 4 records, got %d", len(records))
	}

	r0 := records[0]
	if r0.ID != "intent.ro" || r0.Guard != "task:read-only" || r0.Modality != "MUST" || r0.Effect != "RO_LOCAL" {
		t.Errorf("record 0 mismatch: %+v", r0)
	}

	r1 := records[1]
	if r1.Type != "L" || r1.Payload != "[diary](DIARY.lrf)" {
		t.Errorf("record 1 mismatch: %+v", r1)
	}

	r2 := records[2]
	if r2.Guard != "task:code&tag:subagent" {
		t.Errorf("record 2 guard mismatch: %s", r2.Guard)
	}
}

func TestParseLevel1_InvalidModality(t *testing.T) {
	content := "@lrf=1|aud=test|scope=test\nR|rule.1|*|SHOULD|RO_LOCAL|payload\n"
	_, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 1 || diags[0].Code != "L1:INVALID_MODALITY" {
		t.Fatalf("expected L1:INVALID_MODALITY, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "modality=\"SHOULD\" ∉ {MUST, MUST_NOT, MAY}") {
		t.Errorf("expected PIDGIN notation for modality, got: %s", diags[0].PIDGIN)
	}
}

func TestParseLevel1_InvalidEffect(t *testing.T) {
	content := "@lrf=1|aud=test|scope=test\nR|rule.1|*|MUST|RW_DISK|payload\n"
	_, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 1 || diags[0].Code != "L1:INVALID_EFFECT" {
		t.Fatalf("expected L1:INVALID_EFFECT, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "effect=\"RW_DISK\" ∉ E") {
		t.Errorf("expected PIDGIN notation for effect, got: %s", diags[0].PIDGIN)
	}
}

func TestParseLevel1_InvalidGuard(t *testing.T) {
	cases := []struct {
		name  string
		guard string
	}{
		{"empty", ""},
		{"unknown_key", "unknown:val"},
		{"missing_colon", "task"},
		{"empty_value", "task:"},
		{"double_amp", "task:a&&tag:b"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			content := "@lrf=1|aud=test|scope=test\nR|rule.1|" + c.guard + "|MUST|RO_LOCAL|payload\n"
			_, diags := ParseLevel1("test.lrf", []byte(content))
			found := false
			for _, d := range diags {
				if d.Code == "L1:INVALID_GUARD" {
					found = true
				}
			}
			if !found {
				t.Fatalf("expected L1:INVALID_GUARD for guard '%s', got: %v", c.guard, diags)
			}
		})
	}
}

func TestParseLevel1_InvalidID(t *testing.T) {
	content := "@lrf=1|aud=test|scope=test\nR|id with space|*|MUST|RO_LOCAL|payload\n"
	_, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 1 || diags[0].Code != "L1:INVALID_ID" {
		t.Fatalf("expected L1:INVALID_ID, got: %v", diags)
	}
}

func TestParseLevel1_EmptyPayload(t *testing.T) {
	content := "@lrf=1|aud=test|scope=test\nR|rule.1|*|MUST|RO_LOCAL|   \n"
	_, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 1 || diags[0].Code != "L1:EMPTY_PAYLOAD" {
		t.Fatalf("expected L1:EMPTY_PAYLOAD, got: %v", diags)
	}
}

func TestParseLevel1_LinkPayloadSyntax(t *testing.T) {
	// Raw path without markdown link syntax
	content := "@lrf=1|aud=test|scope=test\nL|load.diary|*|MUST|RO_LOCAL|DIARY.lrf\n"
	_, diags := ParseLevel1("test.lrf", []byte(content))
	if len(diags) != 1 || diags[0].Code != "L1:LINK_PAYLOAD_SYNTAX" {
		t.Fatalf("expected L1:LINK_PAYLOAD_SYNTAX, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "L record payload \"DIARY.lrf\" must be link format") {
		t.Errorf("expected link format error in PIDGIN, got: %s", diags[0].PIDGIN)
	}
}
