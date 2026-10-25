//go:build !windows
package main

// getEnvFromRegistry は環境変数が設定されていない場合に空文字列を返します
func getEnvFromRegistry(name string) string {
	return ""
}
