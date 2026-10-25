package ieee

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

type IEEEArticle struct {
	Title            string   `json:"title"`
	Authors          Authors  `json:"authors"`
	PublicationTitle string   `json:"publication_title"`
	PublicationYear  string   `json:"publication_year"`
	DOI              string   `json:"doi"`
	Abstract         string   `json:"abstract"`
	HTMLURL          string   `json:"html_url"`
}

type Authors struct {
	Authors []AuthorItem `json:"authors"`
}

type AuthorItem struct {
	FullName string `json:"full_name"`
}

type IEEEResponse struct {
	TotalRecords int           `json:"total_records"`
	Articles     []IEEEArticle `json:"articles"`
}

func SearchIEEE(query string) (string, error) {
	apiKey := os.Getenv("IEEE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("IEEE_XPLORE_API_KEY")
	}

	if apiKey == "" {
		return "⚠️ [IEEE Xplore API] IEEE_API_KEY が設定されておりませんわ。\n" +
			"利用には https://developer.ieee.org で開発者登録を行い、API Key を取得して $env:IEEE_API_KEY に設定してくださいまし。", nil
	}

	reqURL := fmt.Sprintf("https://ieeexploreapi.ieee.org/api/v1/search/articles?apikey=%s&querytext=%s&max_records=10", apiKey, url.QueryEscape(query))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("IEEE Xplore API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("⚠️ [IEEE Xplore API Error %d]: %s", resp.StatusCode, string(body)), nil
	}

	var ieeeResp IEEEResponse
	if err := json.Unmarshal(body, &ieeeResp); err != nil {
		return "", fmt.Errorf("Failed to parse IEEE JSON response: %w", err)
	}

	if len(ieeeResp.Articles) == 0 {
		return fmt.Sprintf("🔍 IEEE Xplore: '%s' に一致する論文・技術文書は見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### ⚡ IEEE Xplore 検索結果 (Query: %s, 全 %d 件)\n\n", query, ieeeResp.TotalRecords))
	for i, item := range ieeeResp.Articles {
		sb.WriteString(fmt.Sprintf("#### %d. %s (%s)\n", i+1, item.Title, item.PublicationYear))
		
		var authorNames []string
		for _, a := range item.Authors.Authors {
			authorNames = append(authorNames, a.FullName)
		}
		if len(authorNames) > 0 {
			sb.WriteString(fmt.Sprintf("- **著者**: %s\n", strings.Join(authorNames, ", ")))
		}
		sb.WriteString(fmt.Sprintf("- **掲載誌/会議**: %s\n", item.PublicationTitle))
		if item.DOI != "" {
			sb.WriteString(fmt.Sprintf("- **DOI**: [`https://doi.org/%s`](https://doi.org/%s)\n", item.DOI, item.DOI))
		}
		if item.HTMLURL != "" {
			sb.WriteString(fmt.Sprintf("- **URL**: [IEEE Xplore Direct](%s)\n", item.HTMLURL))
		}
		if item.Abstract != "" {
			abs := item.Abstract
			if len(abs) > 250 {
				abs = abs[:250] + "..."
			}
			sb.WriteString(fmt.Sprintf("- **抄録**: %s\n", abs))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
