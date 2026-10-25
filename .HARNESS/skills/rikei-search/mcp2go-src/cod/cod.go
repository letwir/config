package cod

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CODResultItem struct {
	File       string `json:"file"`
	Formula    string `json:"formula"`
	Spacegroup string `json:"sg"`
	A          string `json:"a"`
	B          string `json:"b"`
	C          string `json:"c"`
	Alpha      string `json:"alpha"`
	Beta       string `json:"beta"`
	Gamma      string `json:"gamma"`
}

func SearchCOD(query string) (string, error) {
	reqURL := fmt.Sprintf("https://www.crystallography.net/cod/result?formula=%s&format=json", url.QueryEscape(query))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("COD API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("🔍 COD: 結晶データ '%s' の取得でステータス %d が返されましたわ。", query, resp.StatusCode), nil
	}

	var results []CODResultItem
	if err := json.Unmarshal(body, &results); err != nil {
		// Single or custom array check
		return fmt.Sprintf("ℹ️ COD レスポンス: %s", string(body)), nil
	}

	if len(results) == 0 {
		return fmt.Sprintf("🔍 COD (Crystallography Open Database): '%s' に一致する結晶構造は見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 💎 Crystallography Open Database (COD) 検索結果 (Query: %s)\n\n", query))
	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		item := results[i]
		sb.WriteString(fmt.Sprintf("#### %d. COD ID: `%s` (組成: %s)\n", i+1, item.File, item.Formula))
		if item.Spacegroup != "" {
			sb.WriteString(fmt.Sprintf("- **空間群**: `%s`\n", item.Spacegroup))
		}
		sb.WriteString(fmt.Sprintf("- **格子定数**: a=%s Å, b=%s Å, c=%s Å (α=%s°, β=%s°, γ=%s°)\n", item.A, item.B, item.C, item.Alpha, item.Beta, item.Gamma))
		sb.WriteString(fmt.Sprintf("- **CIFファイル**: [Download CIF](https://www.crystallography.net/cod/%s.cif)\n", item.File))
		sb.WriteString(fmt.Sprintf("- **Web詳細**: [https://www.crystallography.net/cod/%s.html](https://www.crystallography.net/cod/%s.html)\n\n", item.File, item.File))
	}

	return sb.String(), nil
}
