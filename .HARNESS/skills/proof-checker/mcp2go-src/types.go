package main

import (
	"go/ast"
	"go/token"
)

// Severity represents the diagnostic severity level
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Category represents the rule category based on CODE_RULE.md
type Category string

const (
	CategoryErrorHandling Category = "ERR_HANDLING"
	CategoryNesting       Category = "NESTING_FLOW"
	CategoryResourceRAII  Category = "RESOURCE_RAII"
	CategoryEffectPurity  Category = "EFFECT_PURITY"
	CategoryNamingLength  Category = "NAMING_LENGTH"
	CategoryConcurrency   Category = "CONCURRENCY"
)

// Diagnostic represents a single finding by a proof-check rule
type Diagnostic struct {
	RuleID     string   `json:"rule_id"`
	Category   Category `json:"category"`
	FilePath   string   `json:"file_path"`
	Line       int      `json:"line"`
	Column     int      `json:"column"`
	Message    string   `json:"message"`
	Severity   Severity `json:"severity"`
	Suggestion string   `json:"suggestion,omitempty"`
}

// Rule defines the interface for an AST validation rule
type Rule interface {
	ID() string
	Category() Category
	Description() string
	CheckFile(fset *token.FileSet, file *ast.File, filePath string) []Diagnostic
}

// Config holds runtime options for the proof checker
type Config struct {
	TargetDir  string
	Strict     bool
	IncludeVet bool
	JSONOutput bool
}
