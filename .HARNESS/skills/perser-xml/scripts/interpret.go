package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PageInfo 抽出した普遍的なページ情報
type PageInfo struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	SiteName    string            `json:"site_name"`
	Image       string            `json:"image"`
	Type        string            `json:"type"`
	Meta        map[string]string `json:"meta"`
	JSONLD      []interface{}     `json:"json_ld"`
	OGP         map[string]string `json:"ogp"`
	// 特定サイト向けの拡張情報
	Product *ProductInfo `json:"product,omitempty"`
}

// ProductInfo 特定の製品/作品情報（DMM/DLsiteなどの互換用）
type ProductInfo struct {
	Title   string `json:"title"`
	Maker   string `json:"maker"`
	Release string `json:"release"`
	Genre   string `json:"genre"`
	SKU     string `json:"sku"`
	Price   string `json:"price"`
	Source  string `json:"source"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: interpret.exe <html_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	html := string(data)
	info := interpretHTML(html)

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func interpretHTML(html string) PageInfo {
	info := PageInfo{
		Meta: make(map[string]string),
		OGP:  make(map[string]string),
	}

	// 1. タイトル
	info.Title = extractTitle(html)

	// 2. メタタグの抽出
	extractMetaTags(html, &info)

	// 3. JSON-LDの抽出 (複数対応)
	info.JSONLD = extractAllJSONLD(html)

	// 4. 特定サイト向けの解析 (DMM/FANZA)
	if strings.Contains(html, "dmm.co.jp") || strings.Contains(html, "dmm.com") || strings.Contains(info.Title, "FANZA") {
		product := extractDMMInfo(html, info.Title)
		info.Product = &product
	}

	return info
}

func extractTitle(html string) string {
	re := regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)
	if m := re.FindStringSubmatch(html); len(m) > 1 {
		return cleanHTMLText(m[1])
	}
	return ""
}

func extractMetaTags(html string, info *PageInfo) {
	re := regexp.MustCompile(`(?i)<meta\s+[^>]*(?:name|property)=["']([^"']+)["']\s+content=["']([^"']*)["'][^>]*>`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, m := range matches {
		key := strings.ToLower(m[1])
		val := cleanHTMLText(m[2])

		if strings.HasPrefix(key, "og:") {
			info.OGP[key] = val
			// 代表的なものはトップレベルにも
			switch key {
			case "og:title":
				if info.Title == "" {
					info.Title = val
				}
			case "og:description":
				info.Description = val
			case "og:image":
				info.Image = val
			case "og:site_name":
				info.SiteName = val
			case "og:type":
				info.Type = val
			case "og:url":
				info.URL = val
			}
		} else {
			info.Meta[key] = val
			if key == "description" && info.Description == "" {
				info.Description = val
			}
		}
	}
}

func extractAllJSONLD(html string) []interface{} {
	var result []interface{}
	re := regexp.MustCompile(`(?is)<script[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		var ld interface{}
		if err := json.Unmarshal([]byte(match[1]), &ld); err == nil {
			result = append(result, ld)
		}
	}
	return result
}

func cleanHTMLText(s string) string {
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return strings.TrimSpace(s)
}

// --- DMM/FANZA Specific Logic (既存コードから移植・整理) ---

func extractDMMInfo(html string, rawTitle string) ProductInfo {
	p := ProductInfo{Source: "DMM-Specialized"}

	// タイトルからサークル名抽出
	if strings.Contains(rawTitle, "(") {
		start := strings.Index(rawTitle, "(")
		end := strings.Index(rawTitle, ")")
		if start != -1 && end != -1 && end > start {
			p.Maker = strings.TrimSpace(rawTitle[start+1 : end])
			p.Title = strings.TrimSpace(rawTitle[:start])
		}
	}
	if p.Title == "" {
		p.Title = rawTitle
	}

	// dl/dt/dd 抽出
	makerRe := regexp.MustCompile(`(?is)<dt[^>]*>\s*(?:作者|サークル名)\s*</dt>.*?<dd[^>]*>(.*?)</dd>`)
	if m := makerRe.FindStringSubmatch(html); m != nil {
		p.Maker = cleanTags(m[1])
	}

	dateRe := regexp.MustCompile(`(?is)<dt[^>]*>\s*配信開始日\s*</dt>.*?<dd[^>]*>(.*?)</dd>`)
	if m := dateRe.FindStringSubmatch(html); m != nil {
		p.Release = cleanTags(m[1])
	}

	skuRe := regexp.MustCompile(`(?is)<dt[^>]*>\s*(?:品番|商品番号)\s*</dt>.*?<dd[^>]*>(.*?)</dd>`)
	if m := skuRe.FindStringSubmatch(html); m != nil {
		p.SKU = cleanTags(m[1])
	}

	return p
}

func cleanTags(s string) string {
	tagRe := regexp.MustCompile(`<[^>]+>`)
	s = tagRe.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}
