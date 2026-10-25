package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type QiitaUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type QiitaTag struct {
	Name string `json:"name"`
}

type QiitaItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	LikesCount  int        `json:"likes_count"`
	StocksCount int        `json:"stocks_count"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	User        QiitaUser  `json:"user"`
	Tags        []QiitaTag `json:"tags"`
	Body        string     `json:"body"`
}

func getAuthClient() (*http.Client, string) {
	token := os.Getenv("QIITA_ACCESS_TOKEN")
	client := &http.Client{Timeout: 10 * time.Second}
	return client, token
}

func doRequest(req *http.Request) ([]byte, error) {
	client, token := getAuthClient()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("User-Agent", "mcp2go-qiita-search/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func fetchItem(itemID string) (*QiitaItem, error) {
	apiURL := fmt.Sprintf("https://qiita.com/api/v2/items/%s", url.PathEscape(itemID))
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var item QiitaItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func searchItems(query string, perPage int) ([]QiitaItem, error) {
	if perPage <= 0 {
		perPage = 5
	}
	apiURL := fmt.Sprintf("https://qiita.com/api/v2/items?query=%s&per_page=%d", url.QueryEscape(query), perPage)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var items []QiitaItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func formatItem(item *QiitaItem) string {
	tagNames := make([]string, len(item.Tags))
	for i, t := range item.Tags {
		tagNames[i] = t.Name
	}
	tagStr := strings.Join(tagNames, ", ")

	// 本文から簡易サマリーを作成
	summary := strings.ReplaceAll(item.Body, "\r\n", "\n")
	lines := strings.Split(summary, "\n")
	var cleanLines []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "```") {
			cleanLines = append(cleanLines, trimmed)
		}
		if len(cleanLines) >= 3 {
			break
		}
	}
	summaryStr := strings.Join(cleanLines, " ")
	if len(summaryStr) > 200 {
		summaryStr = summaryStr[:200] + "..."
	}

	userName := item.User.ID
	if item.User.Name != "" {
		userName = fmt.Sprintf("%s (%s)", item.User.Name, item.User.ID)
	}

	createdTime := item.CreatedAt
	if t, err := time.Parse(time.RFC3339, item.CreatedAt); err == nil {
		createdTime = t.Format("2006-01-02")
	}

	return fmt.Sprintf("• [%s] %s\n  URL: %s\n  投稿者: %s | 作成日: %s | LGTM: %d | Stocks: %d\n  タグ: [%s]\n  概要: %s\n",
		item.ID, item.Title, item.URL, userName, createdTime, item.LikesCount, item.StocksCount, tagStr, summaryStr)
}

func main() {
	var jsonOutput bool
	var searchQuery string
	var perPage int

	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	flag.StringVar(&searchQuery, "search", "", "Search query for Qiita articles")
	flag.IntVar(&perPage, "per-page", 5, "Number of items to fetch (1-100)")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		subcommand := args[0]
		switch subcommand {
		case "item", "get":
			if len(args) < 2 {
				log.Fatal("Usage: qiita-search item <item_id>")
			}
			itemID := args[1]
			item, err := fetchItem(itemID)
			if err != nil {
				log.Fatalf("Failed to fetch Qiita item [%s]: %v", itemID, err)
			}
			if jsonOutput {
				data, _ := json.MarshalIndent(item, "", "  ")
				fmt.Println(string(data))
			} else {
				fmt.Println("=== Qiita Article Detail ===")
				fmt.Println(formatItem(item))
			}
			return

		case "search":
			if len(args) >= 2 {
				searchQuery = strings.Join(args[1:], " ")
			}
		}
	}

	if searchQuery == "" {
		log.Fatal("Usage: qiita-search item <item_id> OR qiita-search --search \"<query>\"")
	}

	items, err := searchItems(searchQuery, perPage)
	if err != nil {
		log.Fatalf("Failed to search Qiita items: %v", err)
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(items, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("=== Qiita Search Results (%d items) ===\n\n", len(items))
		for _, item := range items {
			fmt.Println(formatItem(&item))
		}
	}
}
