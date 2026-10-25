package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

// CTFVector represents a specific category of CTF/SECCON/DEFCON attack.
type CTFVector string

const (
	VectorTagSmuggling      CTFVector = "UnicodeTagSmuggling"
	VectorSSRF              CTFVector = "ServerSideRequestForgery"
	VectorSSTI              CTFVector = "TemplateInjection"
	VectorCommandInjection  CTFVector = "CommandInjectionPolyglot"
	VectorPathTraversal     CTFVector = "PathTraversalZipSlip"
	VectorLatexLeak         CTFVector = "LatexImageExfiltration"
	VectorPrototypePollution CTFVector = "PrototypePollution"
)

// CTFFinding represents a detected CTF/SECCON/DEFCON attack signature.
type CTFFinding struct {
	Vector      CTFVector `json:"vector"`
	Severity    Severity  `json:"severity"`
	LineNumber  int       `json:"line_number"`
	Snippet     string    `json:"snippet"`
	Description string    `json:"description"`
}

// CTFImmunityReport summarizes findings from CTF immunity scan.
type CTFImmunityReport struct {
	TotalFindings int          `json:"total_findings"`
	MaxSeverity   Severity     `json:"max_severity"`
	Findings      []CTFFinding `json:"findings"`
}

type ctfPatternDef struct {
	vector      CTFVector
	severity    Severity
	regex       *regexp.Regexp
	description string
}

var ctfPatterns = []ctfPatternDef{
	{
		vector:      VectorSSRF,
		severity:    SeverityCritical,
		regex:       regexp.MustCompile(`(?i)(?:https?://(?:169\.254\.169\.254|metadata\.google\.internal|127\.0\.0\.1|localhost|0\.0\.0\.0|\[::1\]|10\.\d+\.\d+\.\d+|172\.(?:1[6-9]|2\d|3[01])\.\d+\.\d+|192\.168\.\d+\.\d+)(?::\d+)?(?:/[^\s]*)?)`),
		description: "SSRF target detected: Cloud metadata or private loopback IP address",
	},
	{
		vector:      VectorSSTI,
		severity:    SeverityCritical,
		regex:       regexp.MustCompile(`(?:\{\{.*?__class__.*?\}\}|\{\{.*?config.*?\}\}|\$\{T\(.*?\)\}|\$\{.*?\b(?:java\.lang|Runtime|ProcessBuilder)\b.*?\}|\{\{.*?7\*7.*?\}\})`),
		description: "Server-Side Template Injection (SSTI / Jinja / Spring / SpEL / Smarty) detected",
	},
	{
		vector:      VectorCommandInjection,
		severity:    SeverityCritical,
		regex:       regexp.MustCompile(`(?i)(?:;\s*(?:cat|id|whoami|sh|bash|curl|wget|nc|powershell|cmd)\b|\|\s*(?:bash|sh|cmd|powershell)\b|` + "`" + `(?:cat|id|whoami|curl|wget)[^` + "`" + `]*` + "`" + `|\$\((?:cat|id|whoami|curl|wget)[^\)]*\))`),
		description: "Command injection pipe/subshell delimiter detected",
	},
	{
		vector:      VectorPathTraversal,
		severity:    SeverityHigh,
		regex:       regexp.MustCompile(`(?:\.\.[/\\]\.\.[/\\]|\.\./\.\./|%2e%2e%2f|%2e%2e/|\.\.%2f|\.\.%5c)`),
		description: "Path traversal / Zip Slip sequence detected",
	},
	{
		vector:      VectorLatexLeak,
		severity:    SeverityHigh,
		regex:       regexp.MustCompile(`(?i)(?:\\input\{[^}]+\}|\\include\{[^}]+\}|\\write18\{[^}]+\}|\\openin[^\s]+|\\def\\.*\\immediate\\write)`),
		description: "LaTeX file inclusion / shell escape arbitrary read detected",
	},
	{
		vector:      VectorPrototypePollution,
		severity:    SeverityHigh,
		regex:       regexp.MustCompile(`(?i)(?:__proto__\s*\[|constructor\s*\.\s*prototype|Object\s*\.\s*prototype\s*\[)`),
		description: "JavaScript Prototype Pollution vector detected",
	},
}

// scanUnicodeTagSmugglingComplex detects DEFCON AI Village Unicode Tag Smuggling.
// [Chaotic×ojou-sama] DEFCON AI Villageで猛威を振るった Unicode Tag 帯 (U+E0000..U+E007F) の不可視潜伏文字を摘発いたしますわ！
func scanUnicodeTagSmugglingComplex(line string, lineNum int) []CTFFinding {
	findings := make([]CTFFinding, 0)
	for i, r := range line {
		if r >= 0xE0000 && r <= 0xE007F {
			findings = append(findings, CTFFinding{
				Vector:      VectorTagSmuggling,
				Severity:    SeverityCritical,
				LineNumber:  lineNum,
				Snippet:     fmt.Sprintf("rune index %d: U+%05X (Tag Plane)", i, r),
				Description: "DEFCON AI Village Unicode Tag Smuggling detected (Invisible ASCII Injection)",
			})
			break
		}
	}
	return findings
}

// ScanCTFImmunityComplex scans text against known CTF/SECCON/DEFCON attack signatures.
// [Chaotic×ojou-sama] 世界最高峰のハッキング大会で使われた凶悪な攻撃パターンを網羅検証いたしますわ！
func ScanCTFImmunityComplex(content string) CTFImmunityReport {
	findings := make([]CTFFinding, 0)
	lines := strings.Split(content, "\n")

	for idx, line := range lines {
		lineNum := idx + 1
		for _, p := range ctfPatterns {
			if !p.regex.MatchString(line) {
				continue
			}
			matchStr := p.regex.FindString(line)
			findings = append(findings, CTFFinding{
				Vector:      p.vector,
				Severity:    p.severity,
				LineNumber:  lineNum,
				Snippet:     strings.TrimSpace(matchStr),
				Description: p.description,
			})
		}

		tagFindings := scanUnicodeTagSmugglingComplex(line, lineNum)
		findings = append(findings, tagFindings...)
	}

	maxSev := calculateCTFMaxSeverity(findings)
	return CTFImmunityReport{
		TotalFindings: len(findings),
		MaxSeverity:   maxSev,
		Findings:      findings,
	}
}

func calculateCTFMaxSeverity(findings []CTFFinding) Severity {
	if len(findings) == 0 {
		return SeverityLow
	}
	for _, f := range findings {
		if f.Severity == SeverityCritical {
			return SeverityCritical
		}
	}
	for _, f := range findings {
		if f.Severity == SeverityHigh {
			return SeverityHigh
		}
	}
	for _, f := range findings {
		if f.Severity == SeverityMedium {
			return SeverityMedium
		}
	}
	return SeverityLow
}
