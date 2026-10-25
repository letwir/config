package aflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AFLUXRecord map[string]interface{}

func SearchAFLOW(query string) (string, error) {
	// Query parsing: species or compound
	afluxQuery := fmt.Sprintf("species(%s),paging(1)", query)
	if strings.Contains(query, "=") || strings.Contains(query, "(") {
		afluxQuery = query
	}

	reqURL := fmt.Sprintf("http://aflow.org/API/aflux/?%s", url.PathEscape(afluxQuery))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("AFLOW API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("⚠️ [AFLOW API Error %d]: %s", resp.StatusCode, string(body)), nil
	}

	var results []AFLUXRecord
	if err := json.Unmarshal(body, &results); err != nil {
		// If returned string or error
		return fmt.Sprintf("ℹ️ AFLOW Response: %s", string(body)), nil
	}

	if len(results) == 0 {
		return fmt.Sprintf("🔍 AFLOW: '%s' に一致するデータは見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### ⚛️ AFLOW / AFLUX 検索結果 (Query: %s)\n\n", query))
	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		rec := results[i]
		compound, _ := rec["compound"].(string)
		spacegroup, _ := rec["spacegroup_relax"].(string)
		egap, _ := rec["Egap"].(float64)
		auid, _ := rec["auid"].(string)
		enthalpy, _ := rec["enthalpy_formation_atom"].(float64)

		sb.WriteString(fmt.Sprintf("#### %d. 化合物: `%s`\n", i+1, compound))
		if spacegroup != "" {
			sb.WriteString(fmt.Sprintf("- **空間群 (Space Group)**: `%s`\n", spacegroup))
		}
		if egap != 0 {
			sb.WriteString(fmt.Sprintf("- **バンドギャップ (Band Gap)**: `%.3f eV`\n", egap))
		}
		if enthalpy != 0 {
			sb.WriteString(fmt.Sprintf("- **生成エンタルピー**: `%.4f eV/atom`\n", enthalpy))
		}
		if auid != "" {
			sb.WriteString(fmt.Sprintf("- **AUID**: `%s` (Web: [http://aflow.org/material/?id=%s](http://aflow.org/material/?id=%s))\n", auid, auid, auid))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
