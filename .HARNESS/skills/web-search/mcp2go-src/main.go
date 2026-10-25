package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Result struct {
	Rank    int    `json:"rank"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet,omitempty"`
}

type Response struct {
	Status  string   `json:"status"`
	Query   string   `json:"query,omitempty"`
	URL     string   `json:"url,omitempty"`
	Results []Result `json:"results,omitempty"`
	Content string   `json:"content,omitempty"`
	Error   string   `json:"error,omitempty"`
}

var tagPattern = regexp.MustCompile(`(?s)<[^>]+>`)
var linkPattern = regexp.MustCompile(`(?is)<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`)

func client() *http.Client { return &http.Client{Timeout: 20 * time.Second} }

func cleanHTML(s string) string {
	s = tagPattern.ReplaceAllString(s, " ")
	s = strings.NewReplacer("&amp;", "&", "&quot;", "\"", "&#39;", "'", "&lt;", "<", "&gt;", ">", "\n", " ").Replace(s)
	return strings.Join(strings.Fields(s), " ")
}

func search(q string, limit int) ([]Result, error) {
	u := "https://www.google.com/search?q=" + url.QueryEscape(q) + "&hl=en"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; mcp2go-web-search/1.0)")
	resp, err := client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("google search: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google search: HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	results := make([]Result, 0, limit)
	for _, m := range linkPattern.FindAllStringSubmatch(string(body), -1) {
		href := m[1]
		if !strings.HasPrefix(href, "/url?q=") {
			continue
		}
		parsed, err := url.ParseQuery(strings.TrimPrefix(href, "/url?"))
		if err != nil {
			continue
		}
		link := parsed.Get("q")
		if link == "" || !strings.HasPrefix(link, "http") || seen[link] {
			continue
		}
		seen[link] = true
		results = append(results, Result{Rank: len(results) + 1, Title: cleanHTML(m[2]), URL: link})
		if len(results) >= limit {
			break
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("google search returned no parseable results")
	}
	return results, nil
}

func fetch(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("invalid URL: %s", raw)
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; mcp2go-web-search/1.0)")
	resp, err := client().Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", raw, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("fetch %s: HTTP %d", raw, resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", err
	}
	return cleanHTML(string(b)), nil
}

func main() {
	q := flag.String("search", "", "Google search query")
	rawURL := flag.String("url", "", "Fetch and extract a URL")
	limit := flag.Int("limit", 10, "Maximum search results")
	jsonOut := flag.Bool("json", false, "JSON output")
	flag.Parse()
	var out Response
	var err error
	switch {
	case *rawURL != "":
		out.URL = *rawURL
		out.Content, err = fetch(*rawURL)
	case *q != "":
		out.Query = *q
		out.Results, err = search(*q, max(1, min(*limit, 50)))
	default:
		out.Error = "use --search <query> or --url <url>"
	}
	if err != nil {
		out.Error = err.Error()
	}
	if out.Error != "" {
		out.Status = "error"
	} else {
		out.Status = "success"
	}
	if *jsonOut {
		if err := json.NewEncoder(log.Writer()).Encode(out); err != nil {
			log.Fatal(err)
		}
	} else if out.Error != "" {
		log.Fatal(out.Error)
	} else if out.Content != "" {
		fmt.Println(out.Content)
	} else {
		for _, r := range out.Results {
			fmt.Printf("%d. %s\n   %s\n", r.Rank, r.Title, r.URL)
		}
	}
}
