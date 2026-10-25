//go:build windows
package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// getAPIKey は環境変数またはWindowsレジストリからAPIキーを取得します
func getAPIKey() (string, error) {
	key := os.Getenv("KOKKOU_HUDOUSAN_API")
	if key != "" {
		return key, nil
	}

	// Windowsレジストリからの取得を試みる
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		val, _, err := k.GetStringValue("KOKKOU_HUDOUSAN_API")
		if err == nil && val != "" {
			return val, nil
		}
	}

	return "", fmt.Errorf("KOKKOU_HUDOUSAN_API 環境変数またはレジストリキーが設定されていません")
}
