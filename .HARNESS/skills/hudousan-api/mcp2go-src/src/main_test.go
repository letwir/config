package main

import (
	"compress/gzip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// モック用のテスト
func TestGetAPIKey(t *testing.T) {
	// 一時的に環境変数をセットして検証しますの
	orig := os.Getenv("KOKKOU_HUDOUSAN_API")
	defer os.Setenv("KOKKOU_HUDOUSAN_API", orig)

	os.Setenv("KOKKOU_HUDOUSAN_API", "dummy-key-desu-wa")
	key, err := getAPIKey()
	if err != nil {
		t.Fatalf("getAPIKey() err: %v", err)
	}
	if key != "dummy-key-desu-wa" {
		t.Errorf("expected dummy-key-desu-wa, got %s", key)
	}
}

func TestGzipDecoding(t *testing.T) {
	// モックレスポンス（gzip圧縮されたJSON）を作成しますわ
	rawJSON := `{"status":"OK","data":[{"id":"13101","name":"千代田区"}]}`
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = io.WriteString(zw, rawJSON)
	_ = zw.Close()

	// モックサーバーの作成
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf.Bytes())
	}))
	defer ts.Close()

	// HTTPリクエストのテスト
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var reader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	decoded, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read decoded body: %v", err)
	}

	if string(decoded) != rawJSON {
		t.Errorf("expected %s, got %s", rawJSON, string(decoded))
	}
}
