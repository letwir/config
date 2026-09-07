package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	ServerURL string `toml:"server_url"`
	APIKey    string `toml:"api_key"`
	OPDSURL   string `toml:"opds_url"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *GeminiInlineData `json:"inline_data,omitempty"`
}

type GeminiInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func loadAPIKey() string {
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		return apiKey
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		configPaths := []string{
			filepath.Join(exeDir, "config.toml"),
			"config.toml",
		}
		for _, p := range configPaths {
			data, err := os.ReadFile(p)
			if err == nil {
				var cfg Config
				if err := toml.Unmarshal(data, &cfg); err == nil && cfg.APIKey != "" {
					return cfg.APIKey
				}
			}
		}
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: epub-ocr-convert.exe <input_path (epub/image/dir)> [output_txt_path]")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := "書籍-convert.txt"
	if len(os.Args) >= 3 {
		outputPath = os.Args[2]
	}

	apiKey := loadAPIKey()
	if apiKey == "" {
		fmt.Println("Failed: GEMINI_API_KEY not found in environment variables or config.toml")
		os.Exit(1)
	}

	images, err := extractImages(inputPath)
	if err != nil {
		fmt.Printf("Failed to extract images: %v\n", err)
		os.Exit(1)
	}

	if len(images) == 0 {
		fmt.Println("Failed: No valid image files (.png, .jpg, .jpeg, .webp) found.")
		os.Exit(1)
	}

	fmt.Printf("Found %d image pages. Starting Gemini Vision OCR (detecting vertical/horizontal layout)...\n", len(images))

	var ocrTexts []string
	client := &http.Client{Timeout: 60 * time.Second}

	for i, img := range images {
		fmt.Printf("Processing page %d/%d (%s)...\n", i+1, len(images), img.Name)
		text, err := performOCR(client, apiKey, img.MimeType, img.Data)
		if err != nil {
			fmt.Printf("Warning: Failed to OCR page %d (%s): %v\n", i+1, img.Name, err)
			continue
		}
		if text != "" {
			ocrTexts = append(ocrTexts, fmt.Sprintf("--- Page %d (%s) ---\n%s", i+1, img.Name, text))
		}
	}

	var outputContent strings.Builder
	outputContent.WriteString("[変換内容]\n\n")
	outputContent.WriteString(strings.Join(ocrTexts, "\n\n"))

	err = os.WriteFile(outputPath, []byte(outputContent.String()), 0644)
	if err != nil {
		fmt.Printf("Failed: Failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %d pages and saved to '%s'!\n", len(ocrTexts), outputPath)
}

type ImageItem struct {
	Name     string
	MimeType string
	Data     []byte
}

func extractImages(inputPath string) ([]ImageItem, error) {
	fi, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	var items []ImageItem

	if !fi.IsDir() && strings.HasSuffix(strings.ToLower(inputPath), ".epub") {
		// EPUB Zip archive
		r, err := zip.OpenReader(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open EPUB as zip: %w", err)
		}
		defer r.Close()

		for _, f := range r.File {
			mime := getMimeType(f.Name)
			if mime == "" {
				continue
			}

			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			items = append(items, ImageItem{
				Name:     f.Name,
				MimeType: mime,
				Data:     data,
			})
		}
	} else if fi.IsDir() {
		// Directory of images
		entries, err := os.ReadDir(inputPath)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			mime := getMimeType(e.Name())
			if mime == "" {
				continue
			}
			p := filepath.Join(inputPath, e.Name())
			data, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			items = append(items, ImageItem{
				Name:     e.Name(),
				MimeType: mime,
				Data:     data,
			})
		}
	} else {
		// Single image file
		mime := getMimeType(inputPath)
		if mime != "" {
			data, err := os.ReadFile(inputPath)
			if err == nil {
				items = append(items, ImageItem{
					Name:     filepath.Base(inputPath),
					MimeType: mime,
					Data:     data,
				})
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return ""
	}
}

func performOCR(client *http.Client, apiKey, mimeType string, imgData []byte) (string, error) {
	b64Data := base64.StdEncoding.EncodeToString(imgData)

	prompt := "この画像は書籍・解説書のページです。縦書き・横書きを自動検出・判別し、正しい日本語の読み順で正確に文字起こし（OCR）を行ってください。余計な解説や装飾は含めず、抽出された本文テキストのみを出力してください。"

	reqPayload := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
					{
						InlineData: &GeminiInlineData{
							MimeType: mimeType,
							Data:     b64Data,
						},
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	// Try gemini-2.5-flash or fallback gemini-1.5-flash
	models := []string{"gemini-2.5-flash", "gemini-1.5-flash"}
	var lastErr error

	for _, model := range models {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
			continue
		}

		var gResp GeminiResponse
		if err := json.Unmarshal(respBody, &gResp); err != nil {
			lastErr = err
			continue
		}

		if gResp.Error != nil {
			lastErr = fmt.Errorf("Gemini API error: %s", gResp.Error.Message)
			continue
		}

		if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
			return strings.TrimSpace(gResp.Candidates[0].Content.Parts[0].Text), nil
		}
	}

	return "", fmt.Errorf("all model attempts failed: %v", lastErr)
}
