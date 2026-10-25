package main

import (
	"fmt"
)

// Level represents the linter diagnostic level.
type Level int

const (
	Level0 Level = 0 // Lexical / Syntax
	Level1 Level = 1 // Structural / Typing
	Level2 Level = 2 // Semantic / Referential
	Level3 Level = 3 // Cross-Domain Diff / Bi-Simulation
)

// Severity represents diagnostic severity.
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARN"
	SeverityInfo    Severity = "INFO"
)

// Diagnostic represents a single compiler/linter diagnostic issue.
type Diagnostic struct {
	Level    Level    `json:"level"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Col      int      `json:"col"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	PIDGIN   string   `json:"pidgin"`
	Severity Severity `json:"severity"`
}

func (d Diagnostic) String() string {
	if d.PIDGIN != "" {
		return d.PIDGIN
	}
	return fmt.Sprintf("[%s] %s:%d: %s (%s)", d.Severity, d.File, d.Line, d.Message, d.Code)
}

// LRFRecord represents a parsed LRF record (R or L).
type LRFRecord struct {
	LineNum  int    `json:"line_num"`
	Type     string `json:"type"`     // "R" or "L"
	ID       string `json:"id"`       // e.g. "intent.ro"
	Guard    string `json:"guard"`    // e.g. "task:read-only" or "*"
	Modality string `json:"modality"` // e.g. "MUST", "MUST_NOT", "MAY"
	Effect   string `json:"effect"`   // e.g. "RO_LOCAL", "LW_SCOPE"
	Payload  string `json:"payload"`  // e.g. "answer/diagnose/review=inspect+evidence; repair=false"
	RawLine  string `json:"raw_line"`
}

// PIDGINCommand represents a CLI command documented in SKILL.md line 5-10 PIDGIN notation.
type PIDGINCommand struct {
	LineNum  int    `json:"line_num"`
	Emoji    string `json:"emoji"`
	Name     string `json:"name"`
	Notation string `json:"notation"`
	Syntax   string `json:"syntax"`
}

// ASTCommand represents a CLI command extracted from Go AST.
type ASTCommand struct {
	Name     string   `json:"name"`
	Flags    []string `json:"flags"`
	File     string   `json:"file"`
	LineNum  int      `json:"line_num"`
}
