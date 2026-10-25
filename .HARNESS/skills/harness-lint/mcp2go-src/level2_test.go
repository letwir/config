package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// MockFS simulates file system existence for testing.
type MockFS struct {
	existingFiles map[string]bool
}

type dummyFileInfo struct {
	name string
}

func (d dummyFileInfo) Name() string       { return d.name }
func (d dummyFileInfo) Size() int64        { return 100 }
func (d dummyFileInfo) Mode() os.FileMode  { return 0644 }
func (d dummyFileInfo) ModTime() time.Time { return time.Now() }
func (d dummyFileInfo) IsDir() bool        { return false }
func (d dummyFileInfo) Sys() interface{}   { return nil }

func (m MockFS) Stat(path string) (os.FileInfo, error) {
	// Normalize path
	cleanPath := filepath.Clean(path)
	for k := range m.existingFiles {
		if filepath.Clean(k) == cleanPath {
			return dummyFileInfo{name: filepath.Base(path)}, nil
		}
	}
	return nil, errors.New("file not found")
}

func TestCheckLevel2_Valid(t *testing.T) {
	baseDir := "C:/rules"
	fs := MockFS{
		existingFiles: map[string]bool{
			filepath.Join(baseDir, "DIARY.lrf"): true,
		},
	}

	records := []LRFRecord{
		{
			LineNum:  1,
			Type:     "R",
			ID:       "intent.ro",
			Guard:    "task:read-only",
			Modality: "MUST",
			Effect:   "RO_LOCAL",
			Payload:  "answer/diagnose/review=inspect+evidence; repair=false",
		},
		{
			LineNum:  2,
			Type:     "L",
			ID:       "load.diary",
			Guard:    "*",
			Modality: "MUST",
			Effect:   "RO_LOCAL",
			Payload:  "[diary](DIARY.lrf)",
		},
		{
			LineNum:  3,
			Type:     "R",
			ID:       "intent.rw",
			Guard:    "task:change",
			Modality: "MUST",
			Effect:   "LW_SCOPE",
			Payload:  "requested workspace edit+checks=true; append=true",
		},
	}

	diags := CheckLevel2("rules.lrf", records, baseDir, fs)
	if len(diags) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %v", len(diags), diags)
	}
}

func TestCheckLevel2_BrokenLink(t *testing.T) {
	baseDir := "C:/rules"
	fs := MockFS{
		existingFiles: map[string]bool{}, // Empty
	}

	records := []LRFRecord{
		{
			LineNum:  5,
			Type:     "L",
			ID:       "load.missing",
			Guard:    "*",
			Modality: "MUST",
			Effect:   "RO_LOCAL",
			Payload:  "[missing](not_found.lrf)",
		},
	}

	diags := CheckLevel2("rules.lrf", records, baseDir, fs)
	if len(diags) != 1 || diags[0].Code != "L2:BROKEN_LINK" {
		t.Fatalf("expected L2:BROKEN_LINK, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "🔗") || !strings.Contains(diags[0].PIDGIN, "not_found.lrf") {
		t.Errorf("expected PIDGIN broken link format, got: %s", diags[0].PIDGIN)
	}
}

func TestCheckLevel2_DuplicateID(t *testing.T) {
	baseDir := "C:/rules"
	fs := MockFS{}

	records := []LRFRecord{
		{
			LineNum:  2,
			Type:     "R",
			ID:       "rule.dup",
			Guard:    "*",
			Modality: "MUST",
			Effect:   "RO_LOCAL",
			Payload:  "first occurrence",
		},
		{
			LineNum:  10,
			Type:     "R",
			ID:       "rule.dup",
			Guard:    "*",
			Modality: "MAY",
			Effect:   "RO_LOCAL",
			Payload:  "second occurrence",
		},
	}

	diags := CheckLevel2("rules.lrf", records, baseDir, fs)
	if len(diags) != 1 || diags[0].Code != "L2:DUPLICATE_ID" {
		t.Fatalf("expected L2:DUPLICATE_ID, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "🆔") || !strings.Contains(diags[0].PIDGIN, "line 2") {
		t.Errorf("expected PIDGIN duplicate ID format with first occurrence line, got: %s", diags[0].PIDGIN)
	}
}

func TestCheckLevel2_EffectMismatch(t *testing.T) {
	baseDir := "C:/rules"
	fs := MockFS{}

	records := []LRFRecord{
		{
			LineNum:  7,
			Type:     "R",
			ID:       "bad.write",
			Guard:    "*",
			Modality: "MUST",
			Effect:   "RO_LOCAL", // Read-only!
			Payload:  "write=file.txt; mutate database state",
		},
	}

	diags := CheckLevel2("rules.lrf", records, baseDir, fs)
	if len(diags) != 1 || diags[0].Code != "L2:EFFECT_MISMATCH" {
		t.Fatalf("expected L2:EFFECT_MISMATCH, got: %v", diags)
	}
	if !strings.Contains(diags[0].PIDGIN, "⚡") || !strings.Contains(diags[0].PIDGIN, "LW_SCOPE または LIVE_WRITE に昇格せよ") {
		t.Errorf("expected PIDGIN effect mismatch format, got: %s", diags[0].PIDGIN)
	}
}
