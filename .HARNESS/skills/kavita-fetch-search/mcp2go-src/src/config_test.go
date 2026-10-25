package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	const token = "test-token-not-a-secret"
	content := "server_url = \"https://kavita.example.test\"\n" +
		"api_key = \"" + token + "\"\n" +
		"opds_url = \"https://kavita.example.test/api/opds\"\n"
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfigFromPath(configPath)
	if err != nil {
		t.Fatalf("loadConfigFromPath() error = %v", err)
	}
	if cfg.APIKey != token {
		t.Fatal("loadConfigFromPath() did not use the TOML api_key")
	}
}

func TestValidateConfigRejectsInvalidSettings(t *testing.T) {
	valid := Config{
		ServerURL: "https://kavita.example.test",
		APIKey:    "test-token-not-a-secret",
		OPDSURL:   "https://kavita.example.test/api/opds",
	}

	tests := []struct {
		name   string
		config Config
	}{
		{name: "blank server URL", config: Config{APIKey: valid.APIKey, OPDSURL: valid.OPDSURL}},
		{name: "unsupported server URL scheme", config: Config{ServerURL: "ftp://kavita.example.test", APIKey: valid.APIKey, OPDSURL: valid.OPDSURL}},
		{name: "blank API key", config: Config{ServerURL: valid.ServerURL, APIKey: " \t", OPDSURL: valid.OPDSURL}},
		{name: "malformed OPDS URL", config: Config{ServerURL: valid.ServerURL, APIKey: valid.APIKey, OPDSURL: "not a URL"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if err == nil {
				t.Fatal("validateConfig() error = nil")
			}
			if strings.Contains(err.Error(), valid.APIKey) {
				t.Fatal("validateConfig() exposed the API key")
			}
		})
	}
}
