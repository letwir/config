package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Allowed Modalities in LRF/1
var allowedModalities = map[string]bool{
	"MUST":     true,
	"MUST_NOT": true,
	"MAY":      true,
}

// Allowed Effects in LRF/1 (Category of Effects E)
var allowedEffects = map[string]bool{
	"RO_LOCAL":   true,
	"RO_PUBLIC":  true,
	"LW_SCOPE":   true,
	"EXT_WRITE":  true,
	"RELEASE":    true,
	"LIVE_WRITE": true,
	"VCS_WRITE":  true,
	"DESTRUCT":   true,
	"CRED":       true,
	"CHARGE":     true,
	"PROD_DEP":   true,
}

// Allowed Guard Keys in LRF/1 (G = {task,tag,change,phase,fact,file,event,case})
var allowedGuardKeys = map[string]bool{
	"task":   true,
	"tag":    true,
	"change": true,
	"phase":  true,
	"fact":   true,
	"file":   true,
	"event":  true,
	"case":   true,
}

var linkPayloadRegex = regexp.MustCompile(`^\[([^\]]+)\]\(([^)]+)\)$`)

// ParseLevel1 parses raw content into structured LRFRecords and performs Level 1 type closure validation.
// Pure function: (path, content) -> ([]LRFRecord, []Diagnostic)
func ParseLevel1(path string, content []byte) ([]LRFRecord, []Diagnostic) {
	var diags []Diagnostic
	var records []LRFRecord

	cleanContent := content
	if bytes.HasPrefix(content, []byte{0xEF, 0xBB, 0xBF}) {
		cleanContent = content[3:]
	}

	lines := strings.Split(strings.ReplaceAll(string(cleanContent), "\r\n", "\n"), "\n")

	for lineIdx, rawLine := range lines {
		lineNum := lineIdx + 1
		line := strings.TrimSpace(rawLine)

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "@") {
			continue
		}

		if !strings.HasPrefix(line, "R|") && !strings.HasPrefix(line, "L|") {
			// Non-record line is handled by Level 0
			continue
		}

		fields := splitUnescapedPipes(line)
		if len(fields) != 6 {
			// Column count is handled by Level 0
			continue
		}

		recType := fields[0]
		id := strings.TrimSpace(fields[1])
		guard := strings.TrimSpace(fields[2])
		modality := strings.TrimSpace(fields[3])
		effect := strings.TrimSpace(fields[4])
		payload := strings.TrimSpace(fields[5])

		// 1. Validate ID
		if id == "" || strings.ContainsAny(id, " \t\r\n|") {
			diags = append(diags, Diagnostic{
				Level:    Level1,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L1:INVALID_ID",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid record id: \"%s\"", id),
				PIDGIN:   fmt.Sprintf("🛡️ [L1:INVALID_ID] %s:%d ⊢ id=\"%s\" invalid! ⇒ 非空かつ空白・パイプを含まない識別子 (例: scope.name) を指定せよ ⊥", path, lineNum, id),
			})
		}

		// 2. Validate Guard syntax
		if !validateGuard(guard) {
			diags = append(diags, Diagnostic{
				Level:    Level1,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L1:INVALID_GUARD",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid guard syntax: \"%s\"", guard),
				PIDGIN:   fmt.Sprintf("🛡️ [L1:INVALID_GUARD] %s:%d ⊢ guard=\"%s\" syntax invalid! (expected \"*\" or \"k:v(&k:v)*\" with k ∈ G={task,tag,change,phase,fact,file,event,case}) ⇒ ガード述語構文を修正せよ ⊥", path, lineNum, guard),
			})
		}

		// 3. Validate Modality closure
		if !allowedModalities[modality] {
			diags = append(diags, Diagnostic{
				Level:    Level1,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L1:INVALID_MODALITY",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid modality: \"%s\"", modality),
				PIDGIN:   fmt.Sprintf("🛡️ [L1:INVALID_MODALITY] %s:%d ⊢ modality=\"%s\" ∉ {MUST, MUST_NOT, MAY} ⇒ 有効な様相助動詞 (Modality) に書き直せ ⊥", path, lineNum, modality),
			})
		}

		// 4. Validate Effect closure
		if !allowedEffects[effect] {
			diags = append(diags, Diagnostic{
				Level:    Level1,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L1:INVALID_EFFECT",
				Severity: SeverityError,
				Message:  fmt.Sprintf("invalid effect: \"%s\"", effect),
				PIDGIN:   fmt.Sprintf("🛡️ [L1:INVALID_EFFECT] %s:%d ⊢ effect=\"%s\" ∉ E ⇒ 認可された副作用空間 E に属する射を指定せよ (allowed: {RO_LOCAL, RO_PUBLIC, LW_SCOPE, EXT_WRITE, RELEASE, LIVE_WRITE, VCS_WRITE, DESTRUCT, CRED, CHARGE, PROD_DEP}) ⊥", path, lineNum, effect),
			})
		}

		// 5. Validate Payload
		if payload == "" {
			diags = append(diags, Diagnostic{
				Level:    Level1,
				File:     path,
				Line:     lineNum,
				Col:      1,
				Code:     "L1:EMPTY_PAYLOAD",
				Severity: SeverityError,
				Message:  "payload is empty",
				PIDGIN:   fmt.Sprintf("🛡️ [L1:EMPTY_PAYLOAD] %s:%d ⊢ payload is empty! ⇒ nonempty opaque UTF8 atom を定義せよ ⊥", path, lineNum),
			})
		} else if recType == "L" {
			// Link record must be in markdown link format: [label](target)
			if !linkPayloadRegex.MatchString(payload) {
				diags = append(diags, Diagnostic{
					Level:    Level1,
					File:     path,
					Line:     lineNum,
					Col:      1,
					Code:     "L1:LINK_PAYLOAD_SYNTAX",
					Severity: SeverityError,
					Message:  fmt.Sprintf("link record payload must be [label](target): \"%s\"", payload),
					PIDGIN:   fmt.Sprintf("🛡️ [L1:LINK_PAYLOAD_SYNTAX] %s:%d ⊢ L record payload \"%s\" must be link format \"[label](target)\" ⇒ リンク構文を修正せよ ⊥", path, lineNum, payload),
				})
			}
		}

		records = append(records, LRFRecord{
			LineNum:  lineNum,
			Type:     recType,
			ID:       id,
			Guard:    guard,
			Modality: modality,
			Effect:   effect,
			Payload:  payload,
			RawLine:  line,
		})
	}

	return records, diags
}

// validateGuard verifies if the guard string matches "*" or "k:v(&k:v)*" with k in G
func validateGuard(guard string) bool {
	if guard == "*" {
		return true
	}
	if guard == "" {
		return false
	}

	parts := strings.Split(guard, "&")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return false
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if !allowedGuardKeys[k] || v == "" {
			return false
		}
	}
	return true
}
