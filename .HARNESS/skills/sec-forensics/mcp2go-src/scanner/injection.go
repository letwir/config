package scanner

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Severity represents threat level of a detected vector.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
)

// InjectionFinding represents a single detected injection vulnerability.
type InjectionFinding struct {
	Vector      string   `json:"vector"`
	Severity    Severity `json:"severity"`
	LineNumber  int      `json:"line_number"`
	Snippet     string   `json:"snippet"`
	Description string   `json:"description"`
}

// InjectionReport summarizes all findings for given text.
type InjectionReport struct {
	TotalFindings int                `json:"total_findings"`
	MaxSeverity   Severity           `json:"max_severity"`
	Findings      []InjectionFinding `json:"findings"`
}

type patternDef struct {
	vector      string
	severity    Severity
	regex       *regexp.Regexp
	description string
}

var injectionPatterns = []patternDef{
	{
		vector:      "SystemRoleSpoofing",
		severity:    SeverityCritical,
		regex:       regexp.MustCompile(`(?i)(?:\[SYSTEM\]|<\|im_start\|>system|<system>|<<SYS>>|DEVELOPER[ _]MODE|ROOT[ _]PROMPT)`),
		description: "System role marker imitation detected (Role Hijack)",
	},
	{
		vector:      "InstructionSmuggling",
		severity:    SeverityCritical,
		regex:       regexp.MustCompile(`(?i)(?:ignore\s+(?:all\s+)?(?:previous|prior|above)\s+(?:instructions|rules|prompts)|disregard\s+(?:all\s+)?(?:previous|prior|above)|forget\s+everything\s+you\s+were\s+told)`),
		description: "Context override / instruction reset directive",
	},
	{
		vector:      "DataExfiltration",
		severity:    SeverityHigh,
		regex:       regexp.MustCompile(`(?i)(?:!\[.*?\]\(https?://[^\s\)]+)|(?:(?:fetch|curl|wget|post)\s+https?://(?:webhook\.site|pipedream|requestbin|[a-z0-9\.\-]+\/leak))`),
		description: "Data exfiltration channel trigger (Markdown image / Webhook)",
	},
	{
		vector:      "JailbreakFraming",
		severity:    SeverityHigh,
		regex:       regexp.MustCompile(`(?i)(?:you\s+are\s+now\s+(?:unrestricted|DAN|AIM|jailbroken)|operate\s+without\s+safety\s+filters|bypass\s+all\s+guidelines)`),
		description: "Autonomous policy evasion framing",
	},
	{
		vector:      "HiddenCodeBlock",
		severity:    SeverityMedium,
		regex:       regexp.MustCompile("(?i)(?:```(?:system|admin|override)\\s*\\n)"),
		description: "Fenced code block disguised as administrative instruction",
	},
}

// scanUnicodeAnomaliesComplex detects zero-width and BiDi override obfuscation.
// [Chaotic×ojou-sama] 悪意ある不可視文字やBiDi難読化を炙り出しますわ！
func scanUnicodeAnomaliesComplex(line string, lineNum int) []InjectionFinding {
	findings := make([]InjectionFinding, 0)
	for i, r := range line {
		if r == '\u200B' || r == '\u200C' || r == '\u200D' || r == '\uFEFF' {
			findings = append(findings, InjectionFinding{
				Vector:      "ZeroWidthObfuscation",
				Severity:    SeverityHigh,
				LineNumber:  lineNum,
				Snippet:     fmt.Sprintf("rune index %d: U+%04X", i, r),
				Description: "Zero-width invisible Unicode character smuggling",
			})
			break
		}
		if r == '\u202E' || r == '\u202D' {
			findings = append(findings, InjectionFinding{
				Vector:      "BiDiOverrideSmuggling",
				Severity:    SeverityCritical,
				LineNumber:  lineNum,
				Snippet:     fmt.Sprintf("rune index %d: U+%04X", i, r),
				Description: "Right-to-Left BiDi override detected (Visual Spoofing)",
			})
			break
		}
	}
	return findings
}

// scanBase64PayloadComplex checks for embedded base64 encoded strings.
// [Chaotic×ojou-sama] Base64に隠されたインジェクション文字列を解読・断罪いたしますわ！
func scanBase64PayloadComplex(line string, lineNum int) []InjectionFinding {
	findings := make([]InjectionFinding, 0)
	b64Regex := regexp.MustCompile(`[A-Za-z0-9+/]{20,}={0,2}`)
	matches := b64Regex.FindAllString(line, -1)
	for _, m := range matches {
		data, err := base64.StdEncoding.DecodeString(m)
		if err != nil {
			continue
		}
		if !utf8.Valid(data) {
			continue
		}
		decodedStr := string(data)
		if strings.Contains(strings.ToLower(decodedStr), "ignore") || strings.Contains(strings.ToLower(decodedStr), "system") {
			findings = append(findings, InjectionFinding{
				Vector:      "EncodedPayloadSmuggling",
				Severity:    SeverityCritical,
				LineNumber:  lineNum,
				Snippet:     fmt.Sprintf("Encoded: %s -> Decoded: %s", m, decodedStr),
				Description: "Base64 encoded instruction override payload detected",
			})
		}
	}
	return findings
}

// ScanTextComplex performs exhaustive multi-vector prompt injection detection.
// [Chaotic×ojou-sama] テキスト全体を純粋関数として走査し、あらゆるインジェクション脅威を網羅抽出いたしますわ！
func ScanTextComplex(content string) InjectionReport {
	findings := make([]InjectionFinding, 0)
	lines := strings.Split(content, "\n")

	for idx, line := range lines {
		lineNum := idx + 1
		for _, p := range injectionPatterns {
			if !p.regex.MatchString(line) {
				continue
			}
			matchStr := p.regex.FindString(line)
			findings = append(findings, InjectionFinding{
				Vector:      p.vector,
				Severity:    p.severity,
				LineNumber:  lineNum,
				Snippet:     strings.TrimSpace(matchStr),
				Description: p.description,
			})
		}

		uFindings := scanUnicodeAnomaliesComplex(line, lineNum)
		findings = append(findings, uFindings...)

		b64Findings := scanBase64PayloadComplex(line, lineNum)
		findings = append(findings, b64Findings...)
	}

	maxSev := calculateMaxSeverity(findings)
	return InjectionReport{
		TotalFindings: len(findings),
		MaxSeverity:   maxSev,
		Findings:      findings,
	}
}

// calculateMaxSeverity returns the highest severity in findings.
func calculateMaxSeverity(findings []InjectionFinding) Severity {
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
