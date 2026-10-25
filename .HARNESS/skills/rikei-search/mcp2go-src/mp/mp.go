package mp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type MPSummaryItem struct {
	MaterialID          string    `json:"material_id"`
	FormulaPretty       string    `json:"formula_pretty"`
	Symmetry            Symmetry  `json:"symmetry"`
	BandGap             float64   `json:"band_gap"`
	EnergyAboveHull     float64   `json:"energy_above_hull"`
	FormationEnergyPerAtom float64 `json:"formation_energy_per_atom"`
	IsStable            bool      `json:"is_stable"`
	Volume              float64   `json:"volume"`
}

type Symmetry struct {
	Symbol        string `json:"symbol"`
	Number        int    `json:"number"`
	PointGroup    string `json:"point_group"`
	CrystalSystem string `json:"crystal_system"`
}

type MPResponse struct {
	Data []MPSummaryItem `json:"data"`
}

func SearchMP(query string) (string, error) {
	apiKey := os.Getenv("MP_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("MATERIALSPROJECT_API_KEY")
	}

	if apiKey == "" {
		return "⚠️ [Materials Project API] MP_API_KEY が設定されておりませんわ。\n" +
			"利用には https://next-gen.materialsproject.org でアカウント作成し、ダッシュボードから無料API Keyを取得して $env:MP_API_KEY に設定してくださいまし。", nil
	}

	reqURL := fmt.Sprintf("https://api.materialsproject.org/materials/summary/?formula=%s&_limit=10", url.QueryEscape(query))
	if strings.Contains(query, ",") {
		// elements query
		reqURL = fmt.Sprintf("https://api.materialsproject.org/materials/summary/?elements=%s&_limit=10", url.QueryEscape(query))
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Materials Project API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("⚠️ [Materials Project API Error %d]: %s", resp.StatusCode, string(body)), nil
	}

	var mpResp MPResponse
	if err := json.Unmarshal(body, &mpResp); err != nil {
		return "", fmt.Errorf("Failed to parse MP JSON response: %w", err)
	}

	if len(mpResp.Data) == 0 {
		return fmt.Sprintf("🔍 Materials Project: '%s' に一致する物質データは見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 🧪 Materials Project 検索結果 (Query: %s)\n\n", query))
	for i, item := range mpResp.Data {
		sb.WriteString(fmt.Sprintf("#### %d. %s (ID: `%s`)\n", i+1, item.FormulaPretty, item.MaterialID))
		sb.WriteString(fmt.Sprintf("- **結晶系 (Crystal System)**: %s (%s, No.%d)\n", item.Symmetry.CrystalSystem, item.Symmetry.Symbol, item.Symmetry.Number))
		sb.WriteString(fmt.Sprintf("- **バンドギャップ (Band Gap)**: `%.3f eV`\n", item.BandGap))
		sb.WriteString(fmt.Sprintf("- **生成エネルギー (Formation Energy)**: `%.3f eV/atom`\n", item.FormationEnergyPerAtom))
		sb.WriteString(fmt.Sprintf("- **Hullエネルギー (Energy Above Hull)**: `%.3f eV` (安定度: %v)\n", item.EnergyAboveHull, item.IsStable))
		sb.WriteString(fmt.Sprintf("- **単位胞体積 (Volume)**: `%.2f Å³`\n", item.Volume))
		sb.WriteString(fmt.Sprintf("- **詳細Web**: [https://next-gen.materialsproject.org/materials/%s](https://next-gen.materialsproject.org/materials/%s)\n\n", item.MaterialID, item.MaterialID))
	}

	return sb.String(), nil
}
