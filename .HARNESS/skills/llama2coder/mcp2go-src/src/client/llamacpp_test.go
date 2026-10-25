package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLlamaClient_Generate(t *testing.T) {
	// Mock サーバーの作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// パスの確認
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("expected path /v1/chat/completions, got %s", r.URL.Path)
		}

		// メソッドの確認
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		// リクエストボディの確認
		var chatReq ChatRequest
		err := json.NewDecoder(r.Body).Decode(&chatReq)
		if err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if chatReq.Temperature != 0.5 {
			t.Errorf("expected temperature 0.5, got %f", chatReq.Temperature)
		}
		if chatReq.Messages[0].Content != "System Prompt" {
			t.Errorf("expected system message, got %s", chatReq.Messages[0].Content)
		}

		// レスポンスの返却
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := ChatResponse{
			Choices: []struct {
				Message Message `json:"message"`
			}{
				{
					Message: Message{
						Role:    "assistant",
						Content: "print('Mock output')",
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cli := NewLlamaClient(server.URL)
	got, err := cli.Generate(context.Background(), "test-model", "System Prompt", "User Prompt", 0.5, 100)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := "print('Mock output')"
	if got != want {
		t.Errorf("Generate() = %q; want %q", got, want)
	}
}
