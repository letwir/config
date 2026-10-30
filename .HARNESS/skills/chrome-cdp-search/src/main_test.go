package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		engine  string
		count   int
	}{
		{"defaults", []string{"-q", "go"}, false, "bing", 10},
		{"all options", []string{"--query", "a b", "--engine", "google", "--count", "4", "--browser", "/usr/bin/chromium", "--timeout", "2m"}, false, "google", 4},
		{"missing query", nil, true, "", 0},
		{"engine invalid", []string{"-q", "x", "--engine", "yahoo"}, true, "", 0},
		{"count low", []string{"-q", "x", "--count", "0"}, true, "", 0},
		{"count high", []string{"-q", "x", "--count", "51"}, true, "", 0},
		{"timeout invalid", []string{"-q", "x", "--timeout", "soon"}, true, "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := parseArgs(tt.args)
			if (e != nil) != tt.wantErr {
				t.Fatalf("parseArgs error=%v", e)
			}
			if e == nil && (got.Engine != tt.engine || got.Count != tt.count) {
				t.Fatalf("got engine=%q count=%d", got.Engine, got.Count)
			}
		})
	}
}
func TestSearchURL(t *testing.T) {
	for _, engine := range []string{"bing", "google", "duckduckgo"} {
		got := searchURL(engine, "two words & more")
		if !strings.Contains(got, "two+words+%26+more") {
			t.Errorf("%s query not escaped: %s", engine, got)
		}
		if !strings.HasPrefix(got, "https://") {
			t.Errorf("unexpected scheme: %s", got)
		}
	}
}
func TestResultExpression(t *testing.T) {
	for _, engine := range []string{"bing", "google", "duckduckgo"} {
		expr := resultExpression(engine, 7)
		for _, want := range []string{"querySelectorAll", "rank:out.length+1", "https?:", "title:", "snippet:"} {
			if !strings.Contains(expr, want) {
				t.Errorf("%s expression missing %q", engine, want)
			}
		}
		if !strings.Contains(expr, "limit=7") {
			t.Errorf("%s count missing", engine)
		}
	}
}
func TestParseTimeout(t *testing.T) {
	o, e := parseArgs([]string{"-q", "q", "--timeout", "1500ms"})
	if e != nil || o.Timeout != 1500*time.Millisecond {
		t.Fatalf("timeout=%v err=%v", o.Timeout, e)
	}
}
