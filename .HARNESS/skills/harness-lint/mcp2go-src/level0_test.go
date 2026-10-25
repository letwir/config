package main

import (
	"strings"
	"testing"
)

func TestCheckLevel0_Valid(t *testing.T) {
	validLRF := `@lrf=1|aud=GPT-5.6|scope=test
# Comment line
R|intent.ro|task:read-only|MUST|RO_LOCAL|answer/diagnose/review=inspect+evidence; repair=false
L|load.diary|*|MUST|RO_LOCAL|[diary](DIARY.lrf)
R|escaped.pipe|*|MUST|RO_LOCAL|payload with escaped \| pipe inside
`
	diags := CheckLevel0("test.lrf", []byte(validLRF))
	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %v", len(diags), diags)
	}
}

func TestCheckLevel0_BOM(t *testing.T) {
	// BOM + valid content
	content := append([]byte{0xEF, 0xBB, 0xBF}, []byte("@lrf=1|aud=test|scope=test\nR|id|*|MUST|RO_LOCAL|payload\n")...)
	diags := CheckLevel0("bom.lrf", content)
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic (BOM only), got %d: %v", len(diags), diags)
	}
	if diags[0].Code != "L0:UTF8_BOM" {
		t.Errorf("expected L0:UTF8_BOM, got %s", diags[0].Code)
	}
}

func TestCheckLevel0_BOM_With_InvalidHeader(t *testing.T) {
	// BOM + invalid header -> both BOM and HEADER_SYNTAX should be caught
	content := append([]byte{0xEF, 0xBB, 0xBF}, []byte("@lrf=2|invalid\nR|id|*|MUST|RO_LOCAL|payload\n")...)
	diags := CheckLevel0("bom_bad.lrf", content)
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics (BOM and HEADER_SYNTAX), got %d: %v", len(diags), diags)
	}
	codes := map[string]bool{diags[0].Code: true, diags[1].Code: true}
	if !codes["L0:UTF8_BOM"] || !codes["L0:HEADER_SYNTAX"] {
		t.Errorf("expected L0:UTF8_BOM and L0:HEADER_SYNTAX, got: %v", codes)
	}
}

func TestCheckLevel0_MissingHeader(t *testing.T) {
	noHeader := "R|id|*|MUST|RO_LOCAL|payload\n"
	diags := CheckLevel0("no_header.lrf", []byte(noHeader))
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic (MISSING_HEADER), got %d: %v", len(diags), diags)
	}
	if diags[0].Code != "L0:MISSING_HEADER" {
		t.Errorf("expected L0:MISSING_HEADER, got %s", diags[0].Code)
	}
}

func TestCheckLevel0_DuplicateHeader(t *testing.T) {
	dupHeader := "@lrf=1|aud=test|scope=test\n@lrf=1|aud=test|scope=test\nR|id|*|MUST|RO_LOCAL|payload\n"
	diags := CheckLevel0("dup.lrf", []byte(dupHeader))
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic (HEADER_DUPLICATE), got %d: %v", len(diags), diags)
	}
	if diags[0].Code != "L0:HEADER_DUPLICATE" {
		t.Errorf("expected L0:HEADER_DUPLICATE, got %s", diags[0].Code)
	}
}

func TestCheckLevel0_HeaderOrder(t *testing.T) {
	badOrder := "R|id|*|MUST|RO_LOCAL|payload\n@lrf=1|aud=test|scope=test\n"
	diags := CheckLevel0("order.lrf", []byte(badOrder))
	foundOrder := false
	for _, d := range diags {
		if d.Code == "L0:HEADER_ORDER" {
			foundOrder = true
		}
	}
	if !foundOrder {
		t.Fatalf("expected L0:HEADER_ORDER, got: %v", diags)
	}
}

func TestCheckLevel0_HeaderSyntax(t *testing.T) {
	// Missing aud
	noAud := "@lrf=1|scope=test\nR|id|*|MUST|RO_LOCAL|payload\n"
	diags := CheckLevel0("no_aud.lrf", []byte(noAud))
	if len(diags) != 1 || diags[0].Code != "L0:HEADER_SYNTAX" {
		t.Errorf("expected L0:HEADER_SYNTAX for missing aud, got: %v", diags)
	}

	// Extra pipes
	extraPipes := "@lrf=1|aud=test|scope=test|extra\nR|id|*|MUST|RO_LOCAL|payload\n"
	diags = CheckLevel0("extra.lrf", []byte(extraPipes))
	if len(diags) != 1 || diags[0].Code != "L0:HEADER_SYNTAX" {
		t.Errorf("expected L0:HEADER_SYNTAX for extra pipes, got: %v", diags)
	}
}

func TestCheckLevel0_UnknownSyntax_Prose(t *testing.T) {
	prose := `@lrf=1|aud=test|scope=test
Here is a natural language sentence.
- Bullet point
R|id|*|MUST|RO_LOCAL|payload
`
	diags := CheckLevel0("prose.lrf", []byte(prose))
	if len(diags) != 2 {
		t.Fatalf("expected 2 UNKNOWN_SYNTAX diagnostics, got %d: %v", len(diags), diags)
	}
	for _, d := range diags {
		if d.Code != "L0:UNKNOWN_SYNTAX" {
			t.Errorf("expected L0:UNKNOWN_SYNTAX, got %s", d.Code)
		}
	}
}

func TestCheckLevel0_InvalidRecordPrefix(t *testing.T) {
	badPrefix := `@lrf=1|aud=test|scope=test
M|id|*|MUST|RO_LOCAL|payload
X|id|*|MUST|RO_LOCAL|payload
r|id|*|MUST|RO_LOCAL|payload
`
	diags := CheckLevel0("bad_prefix.lrf", []byte(badPrefix))
	if len(diags) != 3 {
		t.Fatalf("expected 3 RECORD_PREFIX diagnostics, got %d: %v", len(diags), diags)
	}
	for _, d := range diags {
		if d.Code != "L0:RECORD_PREFIX" {
			t.Errorf("expected L0:RECORD_PREFIX, got %s", d.Code)
		}
	}
}

func TestCheckLevel0_ColumnCount(t *testing.T) {
	// 5 fields (missing payload)
	tooFew := "@lrf=1|aud=test|scope=test\nR|id|*|MUST|RO_LOCAL\n"
	diags := CheckLevel0("few.lrf", []byte(tooFew))
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Code != "L0:INVALID_COLUMN_COUNT" {
		t.Errorf("expected L0:INVALID_COLUMN_COUNT, got %s", diags[0].Code)
	}

	// 7 fields (raw pipe collision in payload)
	collision := "@lrf=1|aud=test|scope=test\nR|id|*|MUST|RO_LOCAL|payload|with unescaped pipe\n"
	diags = CheckLevel0("collision.lrf", []byte(collision))
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Code != "L0:PIPE_COLLISION" {
		t.Errorf("expected L0:PIPE_COLLISION, got %s", diags[0].Code)
	}
	if !strings.Contains(diags[0].PIDGIN, "delimiter collision") {
		t.Errorf("expected delimiter collision in PIDGIN, got: %s", diags[0].PIDGIN)
	}
}
