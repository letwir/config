package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LlamaClient は llama.cpp (OpenAI互換) REST API クライアントですわ。
type LlamaClient struct {
	BaseURL string
	HTTPCli *http.Client
}

// NewLlamaClient は新しい LlamaClient を生成しますわ。
func NewLlamaClient(baseURL string) *LlamaClient {
	return &LlamaClient{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTPCli: &http.Client{},
	}
}

// Message はチャットメッセージ構造体ですわ。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest は OpenAI 互換のチャットリクエスト構造体ですわ。
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`
}

// ChatResponse は OpenAI 互換のチャットレスポンス構造体ですわ。
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// ChatChunk はストリーミング時のチャンク構造体ですわ。
type ChatChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// Generate は一括でコードを生成しますわ。
func (c *LlamaClient) Generate(ctx context.Context, model, systemPrompt, userPrompt string, temp float64, maxTokens int) (string, error) {
	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: temp,
		Stream:      false,
	}
	if maxTokens > 0 {
		reqBody.MaxTokens = maxTokens
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPCli.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// Stream はレスポンスをストリーミング出力しますわ。
// out に逐次トークンを書き込みます。
func (c *LlamaClient) Stream(ctx context.Context, model, systemPrompt, userPrompt string, temp float64, maxTokens int, out io.Writer) error {
	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: temp,
		Stream:      true,
	}
	if maxTokens > 0 {
		reqBody.MaxTokens = maxTokens
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPCli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("error reading stream: %w", err)
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// SSE format check: "data: ..."
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			dataStr := strings.TrimPrefix(line, "data: ")
			if dataStr == "[DONE]" {
				return nil
			}

			var chunk ChatChunk
			if err := json.Unmarshal([]byte(dataStr), &chunk); err != nil {
				// パースエラーは無視して次に進むか、エラーハンドリングを行いますわ
				continue
			}

			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				if content != "" {
					_, _ = io.WriteString(out, content)
				}
			}
		}
	}
}
