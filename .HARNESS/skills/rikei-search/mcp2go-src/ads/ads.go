package ads

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

type ADSDoc struct {
	Bibcode  string   `json:"bibcode"`
	Title    []string `json:"title"`
	Author   []string `json:"author"`
	Year     string   `json:"year"`
	DOI      []string `json:"doi"`
	Abstract string   `json:"abstract"`
}

type ADSResponseBody struct {
	NumFound int      `json:"numFound"`
	Docs     []ADSDoc `json:"docs"`
}

type ADSResponse struct {
	Response ADSResponseBody `json:"response"`
}

func SearchADS(query string) (string, error) {
	token := os.Getenv("ADS_API_TOKEN")
	if token == "" {
		token = os.Getenv("NASA_ADS_TOKEN")
	}

	if token == "" {
		return "⚠️ [NASA ADS API] ADS_API_TOKEN が設定されておりませんわ。\n" +
			"利用には https://ui.adsabs.harvard.edu/ で無料アカウントを作成し、設定ページから API Token を取得して $env:ADS_API_TOKEN に設定してくださいまし。", nil
	}

	reqURL := fmt.Sprintf("https://api.adsabs.harvard.edu/v1/search/query?q=%s&fl=bibcode,title,author,year,doi,abstract&rows=10", url.QueryEscape(query))

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("NASA ADS API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("⚠️ [NASA ADS API Error %d]: %s", resp.StatusCode, string(body)), nil
	}

	var adsResp ADSResponse
	if err := json.Unmarshal(body, &adsResp); err != nil {
		return "", fmt.Errorf("Failed to parse NASA ADS JSON response: %w", err)
	}

	docs := adsResp.Response.Docs
	if len(docs) == 0 {
		return fmt.Sprintf("🔍 NASA ADS: '%s' に一致する物理学・天文学文献は見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### 🔭 NASA ADS (Astrophysics & Physics) 検索結果 (Query: %s, 全 %d 件)\n\n", query, adsResp.Response.NumFound))
	for i, doc := range docs {
		titleStr := "No Title"
		if len(doc.Title) > 0 {
			titleStr = doc.Title[0]
		}
		sb.WriteString(fmt.Sprintf("#### %d. %s (%s)\n", i+1, titleStr, doc.Year))
		
		if len(doc.Author) > 0 {
			authors := strings.Join(doc.Author, ", ")
			if len(authors) > 150 {
				authors = authors[:150] + " et al."
			}
			sb.WriteString(fmt.Sprintf("- **著者**: %s\n", authors))
		}
		sb.WriteString(fmt.Sprintf("- **Bibcode**: `%s`\n", doc.Bibcode))
		if len(doc.DOI) > 0 {
			sb.WriteString(fmt.Sprintf("- **DOI**: [`https://doi.org/%s`](https://doi.org/%s)\n", doc.DOI[0], doc.DOI[0]))
		}
		sb.WriteString(fmt.Sprintf("- **NASA ADS リンク**: [https://ui.adsabs.harvard.edu/abs/%s/abstract](https://ui.adsabs.harvard.edu/abs/%s/abstract)\n", doc.Bibcode, doc.Bibcode))
		if doc.Abstract != "" {
			abs := doc.Abstract
			if len(abs) > 250 {
				abs = abs[:250] + "..."
			}
			sb.WriteString(fmt.Sprintf("- **抄録**: %s\n", abs))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
