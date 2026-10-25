//go:build !windows
package main

import (
	"fmt"
	"os"
)

// getAPIKey は環境変数からAPIキーを取得します
func getAPIKey() (string, error) {
	key := os.Getenv("KOKKOU_HUDOUSAN_API")
	if key != "" {
		return key, nil
	}
	return "", fmt.Errorf("KOKKOU_HUDOUSAN_API 環境変数が設定されていません")
}
