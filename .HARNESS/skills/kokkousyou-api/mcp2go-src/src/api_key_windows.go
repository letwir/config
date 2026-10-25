//go:build windows
package main

import (
	"golang.org/x/sys/windows/registry"
)

// getEnvFromRegistry は環境変数が設定されていない場合に Windows レジストリから取得を試みます
func getEnvFromRegistry(name string) string {
	// ユーザー環境変数 (HKCU)
	k, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err == nil {
		val, _, err := k.GetStringValue(name)
		k.Close()
		if err == nil && val != "" {
			return val
		}
	}
	// システム環境変数 (HKLM)
	k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	if err == nil {
		val, _, err := k.GetStringValue(name)
		k.Close()
		if err == nil && val != "" {
			return val
		}
	}
	return ""
}
