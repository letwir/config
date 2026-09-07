package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"google.golang.org/genai"
)

// リダイレクト先の真のURLを引っこ抜くメソッドですわ！
func resolveURL(rawURL string) string {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 自動でリダイレクトを追従させず、302レスポンスのまま返させますの！
			return http.ErrUseLastResponse
		},
	}

	// GETではなくHEADリクエストでヘッダだけ頂戴しますわ
	resp, err := client.Head(rawURL)
	if err != nil {
		return rawURL // エラーが起きたら仕方なく元の不細工なURLを返しますの
	}
	defer resp.Body.Close()

	// 300番台のステータスコードならLocationヘッダを覗き見しますわよ
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if loc, err := resp.Location(); err == nil && loc.String() != "" {
			return loc.String() // これがCurl可能な真のURLですわ！
		}
	}
	return rawURL
}

// Qiita URLからアイテムIDを抽出するメソッドですわ！
func extractQiitaItemID(rawURL string) string {
	re := regexp.MustCompile(`qiita\.com/(?:[^/]+/)?items/([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(rawURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// 別バイナリ qiita-search.exe をフック呼び出しするメソッドですわ！
func hookQiitaItem(itemID string) string {
	exePath := "qiita-search.exe"
	if execDir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		candidate := filepath.Join(execDir, "qiita-search.exe")
		if _, err := os.Stat(candidate); err == nil {
			exePath = candidate
		}
	}

	cmd := exec.Command(exePath, "item", itemID)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func main() {
	apiKey := os.Getenv("GEMINI_GROUNDING_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_GROUNDING_API_KEY environment variable is not set")
	}

	// Parse command-line flags
	searchQuery := flag.String("search", "", "Search query to ground with Gemini")
	searchQueryEnv := os.Getenv("GEMINI_SEARCH_QUERY")
	flag.Parse()

	q := *searchQuery
	if q == "" && searchQueryEnv != "" {
		q = searchQueryEnv
	}
	if q == "" {
		log.Fatal("Usage: grounding_search.exe --search \"検索キーワード\" or set GEMINI_SEARCH_QUERY env var")
	}
	log.Printf("DEBUG: Query received: %s", q)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		// Note: genai.Client does not have a Close method in this SDK version
		_ = client
	}()

	modelName := "gemini-2.5-flash"
	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}
	config := &genai.GenerateContentConfig{Tools: tools}

	result, err := client.Models.GenerateContent(ctx, modelName, genai.Text(*searchQuery), config)
	if err != nil {
		log.Fatalf("GenerateContent error: %v", err)
	}

	if len(result.Candidates) == 0 {
		log.Fatal("No candidates returned")
	}

	fmt.Println("=== Gemini Grounding Search Results ===\n")
	fmt.Println("【回答】")
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		}
	}
	fmt.Println()

	fmt.Println("【参照情報源】")
	var detectedQiitaIDs []string
	qiitaSeen := make(map[string]bool)

	if result.Candidates[0].GroundingMetadata != nil {
		for _, query := range result.Candidates[0].GroundingMetadata.WebSearchQueries {
			fmt.Printf("  • Search query: %s\n", query)
		}
		for _, chunk := range result.Candidates[0].GroundingMetadata.GroundingChunks {
			if chunk.Web != nil && chunk.Web.URI != "" {
				title := ""
				if chunk.Web.Title != "" {
					title = chunk.Web.Title
				}
				// わたくしたちの特製関数を通しますわ！
				realURI := resolveURL(chunk.Web.URI)
				fmt.Printf("  • %s\n    %s\n\n", title, realURI)

				// Qiita URLの自動検知
				if itemID := extractQiitaItemID(realURI); itemID != "" {
					if !qiitaSeen[itemID] {
						qiitaSeen[itemID] = true
						detectedQiitaIDs = append(detectedQiitaIDs, itemID)
					}
				} else if strings.Contains(realURI, "qiita.com") {
					// 記事個別URLでない場合も検知フラグ用
				}
			}
		}
	} else {
		fmt.Println("  （参照情報源なし）")
	}

	// 自動検出されたQiita記事が存在する場合、別バイナリ qiita-search.exe をフック呼び出しいたしますわ！
	if len(detectedQiitaIDs) > 0 {
		fmt.Println("=== 【自動検出】Qiita API 追補詳細情報 ===")
		for _, itemID := range detectedQiitaIDs {
			info := hookQiitaItem(itemID)
			if info != "" {
				fmt.Println(info)
			}
		}
	}
}

