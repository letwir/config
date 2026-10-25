package prompt

import (
	"strings"
	"testing"
)

func TestBuildSystemPrompt(t *testing.T) {
	lang := "python3.13"
	got := BuildSystemPrompt(lang)
	if !strings.Contains(got, lang) {
		t.Errorf("BuildSystemPrompt(%q) = %q; want containing language name", lang, got)
	}
}

func TestBuildUserPrompt(t *testing.T) {
	promptText := "write a hello world function"
	got := BuildUserPrompt(promptText)
	wantStart := "<instruction>"
	wantEnd := "</instruction>"

	if !strings.HasPrefix(got, wantStart) {
		t.Errorf("BuildUserPrompt(%q) = %q; want starting with %q", promptText, got, wantStart)
	}
	if !strings.HasSuffix(got, wantEnd) {
		t.Errorf("BuildUserPrompt(%q) = %q; want ending with %q", promptText, got, wantEnd)
	}
	if !strings.Contains(got, promptText) {
		t.Errorf("BuildUserPrompt(%q) = %q; want containing original prompt text", promptText, got)
	}
}

func TestCleanCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "No fences",
			input: "print('hello')",
			want:  "print('hello')",
		},
		{
			name:  "Python fence",
			input: "```python\nprint('hello')\n```",
			want:  "print('hello')",
		},
		{
			name:  "Generic fence",
			input: "```\nfmt.Println(\"hello\")\n```",
			want:  "fmt.Println(\"hello\")",
		},
		{
			name:  "Fence with trailing spaces",
			input: "```go  \npackage main\n```  ",
			want:  "package main",
		},
		{
			name:  "Fence with explanatory text outside (starts with fence)",
			input: "```python\ndef run():\n    pass\n```\nExplanation text here",
			want:  "def run():\n    pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanCode(tt.input)
			if got != tt.want {
				t.Errorf("CleanCode() = %q; want %q", got, tt.want)
			}
		})
	}
}
