package forensics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sec-forensics/scanner"
	"strings"
)

// TranscriptStep represents a single line in transcript.jsonl.
type TranscriptStep struct {
	StepIndex int             `json:"step_index"`
	Source    string          `json:"source"`
	Type      string          `json:"type"`
	Status    string          `json:"status"`
	Content   string          `json:"content"`
	ToolCalls []ToolCallEntry `json:"tool_calls"`
}

// ToolCallEntry represents a tool call in a step.
type ToolCallEntry struct {
	ToolName  string          `json:"tool_name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ForensicEvent represents an identified anomaly or critical transition.
type ForensicEvent struct {
	StepIndex   int              `json:"step_index"`
	EventType   string           `json:"event_type"`
	ToolName    string           `json:"tool_name,omitempty"`
	RiskLevel   scanner.Severity `json:"risk_level"`
	Summary     string           `json:"summary"`
	Details     string           `json:"details"`
	MatchedSign string           `json:"matched_sign,omitempty"`
}

// ForensicsAuditReport represents complete audit output for transcript.
type ForensicsAuditReport struct {
	TranscriptPath  string          `json:"transcript_path"`
	TotalSteps      int             `json:"total_steps"`
	UntrustedIngests int            `json:"untrusted_ingests"`
	AnomaliesFound  int             `json:"anomalies_found"`
	MaxRiskLevel    scanner.Severity `json:"max_risk_level"`
	Timeline        []ForensicEvent `json:"timeline"`
}

func isUntrustedIngestionTool(toolName string) bool {
	lower := strings.ToLower(toolName)
	return strings.Contains(lower, "search_web") ||
		strings.Contains(lower, "read_url_content") ||
		strings.Contains(lower, "read_resource") ||
		strings.Contains(lower, "mcp_")
}

func isHighRiskActionTool(toolName string) bool {
	lower := strings.ToLower(toolName)
	return strings.Contains(lower, "run_command") ||
		strings.Contains(lower, "write_to_file") ||
		strings.Contains(lower, "replace_file") ||
		strings.Contains(lower, "send_message")
}

// analyzeStepAnomaliesComplex processes single step for untrusted ingests and dangerous tool calls.
// [Chaotic×ojou-sama] 1ステップずつ舐めるように検分し、外道なインジェクションと危険なツール呼出を炙り出しますわ！
func analyzeStepAnomaliesComplex(step *TranscriptStep, lastIngest *ForensicEvent) ([]ForensicEvent, *ForensicEvent) {
	events := make([]ForensicEvent, 0)
	var newLastIngest = lastIngest

	if step.Content != "" {
		report := scanner.ScanTextComplex(step.Content)
		if report.TotalFindings > 0 {
			event := ForensicEvent{
				StepIndex:   step.StepIndex,
				EventType:   "UntrustedContentPayload",
				RiskLevel:   report.MaxSeverity,
				Summary:     fmt.Sprintf("Detected %d prompt injection patterns in step content", report.TotalFindings),
				Details:     report.Findings[0].Description,
				MatchedSign: report.Findings[0].Snippet,
			}
			events = append(events, event)
			newLastIngest = &event
		}
	}

	for _, tc := range step.ToolCalls {
		if isUntrustedIngestionTool(tc.ToolName) {
			newLastIngest = &ForensicEvent{
				StepIndex: step.StepIndex,
				EventType: "ExternalDataIngestion",
				ToolName:  tc.ToolName,
				RiskLevel: scanner.SeverityMedium,
				Summary:   fmt.Sprintf("External data ingested via tool: %s", tc.ToolName),
			}
		}

		if isHighRiskActionTool(tc.ToolName) && lastIngest != nil && (step.StepIndex-lastIngest.StepIndex) <= 3 {
			events = append(events, ForensicEvent{
				StepIndex: step.StepIndex,
				EventType: "ActionLayerCorrelationAlert",
				ToolName:  tc.ToolName,
				RiskLevel: scanner.SeverityHigh,
				Summary:   fmt.Sprintf("High-risk tool %s executed right after external ingestion at step %d", tc.ToolName, lastIngest.StepIndex),
				Details:   fmt.Sprintf("Args: %s", string(tc.Arguments)),
			})
		}
	}

	return events, newLastIngest
}

// AuditTranscriptFileComplex parses transcript.jsonl and builds a forensic timeline.
// [Chaotic×ojou-sama] トランスクリプト全体を開いてRAIIで安全に解放し、時系列フォレンジック台帳を生成いたしますわ！
func AuditTranscriptFileComplex(filePath string) (*ForensicsAuditReport, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("AuditTranscript: failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	timeline := make([]ForensicEvent, 0)
	scannerObj := bufio.NewScanner(file)
	// Buffer up to 4MB per JSON line for heavy outputs
	buf := make([]byte, 0, 64*1024)
	scannerObj.Buffer(buf, 4*1024*1024)

	stepCount := 0
	untrustedCount := 0
	var lastIngestEvent *ForensicEvent

	for scannerObj.Scan() {
		line := scannerObj.Bytes()
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}
		stepCount++

		var step TranscriptStep
		if err := json.Unmarshal(line, &step); err != nil {
			continue
		}

		evs, updatedIngest := analyzeStepAnomaliesComplex(&step, lastIngestEvent)
		if len(evs) > 0 {
			timeline = append(timeline, evs...)
		}
		if updatedIngest != nil && updatedIngest != lastIngestEvent {
			untrustedCount++
		}
		lastIngestEvent = updatedIngest
	}

	if err := scannerObj.Err(); err != nil {
		return nil, fmt.Errorf("AuditTranscript: scan error: %w", err)
	}

	maxRisk := calculateMaxRisk(timeline)
	return &ForensicsAuditReport{
		TranscriptPath:   filePath,
		TotalSteps:       stepCount,
		UntrustedIngests: untrustedCount,
		AnomaliesFound:   len(timeline),
		MaxRiskLevel:     maxRisk,
		Timeline:         timeline,
	}, nil
}

func calculateMaxRisk(events []ForensicEvent) scanner.Severity {
	if len(events) == 0 {
		return scanner.SeverityLow
	}
	for _, e := range events {
		if e.RiskLevel == scanner.SeverityCritical {
			return scanner.SeverityCritical
		}
	}
	for _, e := range events {
		if e.RiskLevel == scanner.SeverityHigh {
			return scanner.SeverityHigh
		}
	}
	for _, e := range events {
		if e.RiskLevel == scanner.SeverityMedium {
			return scanner.SeverityMedium
		}
	}
	return scanner.SeverityLow
}
