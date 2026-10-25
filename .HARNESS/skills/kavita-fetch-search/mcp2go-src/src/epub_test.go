package main

import (
	"archive/zip"
	"bytes"
	"mcp2go/epub"
	"testing"
)

func TestStripTags(t *testing.T) {
	input := "<html><body><h1>Hello, World!</h1><p>This is a <b>test</b>.</p></body></html>"
	expected := "Hello, World!This is a test."
	result := epub.StripTags(input)
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

func TestParseEPUBBytes(t *testing.T) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	f, err := w.Create("OEBPS/text01.html")
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	_, err = f.Write([]byte("<html><body><p>Hello from EPUB!</p></body></html>"))
	if err != nil {
		t.Fatalf("Failed to write to zip file: %v", err)
	}

	w.Close()

	parsed, err := epub.ParseEPUBBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to parse EPUB bytes: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("Expected 1 parsed file, but got %d", len(parsed))
	}

	if parsed[0].Name != "OEBPS/text01.html" {
		t.Errorf("Expected file name 'OEBPS/text01.html', but got %q", parsed[0].Name)
	}

	if parsed[0].Text != "Hello from EPUB!" {
		t.Errorf("Expected text 'Hello from EPUB!', but got %q", parsed[0].Text)
	}
}
