package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "fixture.toml")
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

// Test lookups never inspect the process environment or private configuration.
func envLookup(values map[string]string) func(string) (string, bool) {
	return func(name string) (string, bool) {
		value, ok := values[name]
		return value, ok
	}
}

func TestApplyConfigEnv(t *testing.T) {
	legacy := Config{ServerURL: "https://legacy.example.test", APIKey: "legacy-placeholder", OPDSURL: "https://legacy.example.test/api/opds"}
	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{"legacy fallback", nil, legacy},
		{"URL only", map[string]string{"KAVITA_URL": "https://env.example.test"}, Config{ServerURL: "https://env.example.test", APIKey: legacy.APIKey, OPDSURL: legacy.OPDSURL}},
		{"OPDS only", map[string]string{"KAVITA_OPDS": "https://env.example.test/api/opds/new-placeholder/"}, Config{ServerURL: legacy.ServerURL, APIKey: "new-placeholder", OPDSURL: "https://env.example.test/api/opds/new-placeholder/"}},
		{"both overrides", map[string]string{"KAVITA_URL": "https://env.example.test", "KAVITA_OPDS": "https://opds.example.test/prefix/api/opds/a%2Fb%20c%2B"}, Config{ServerURL: "https://env.example.test", APIKey: "a/b c+", OPDSURL: "https://opds.example.test/prefix/api/opds/a%2Fb%20c%2B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyConfigEnv(legacy, envLookup(tt.env))
			if err != nil || got != tt.want {
				t.Fatal("environment precedence did not produce the expected settings")
			}
		})
	}
	_, err := applyConfigEnv(legacy, func(name string) (string, bool) {
		if name != "KAVITA_URL" && name != "KAVITA_OPDS" {
			t.Fatal("configuration looked up an unrelated environment variable")
		}
		return "", false
	})
	if err != nil {
		t.Fatal("absent environment overrides failed")
	}
}

func TestOPDSAPIKey(t *testing.T) {
	tests := []struct {
		name, value, key string
	}{
		{"valid", "https://example.test/api/opds/test-placeholder", "test-placeholder"},
		{"trailing slash", "https://example.test/api/opds/test-placeholder/", "test-placeholder"},
		{"multiple trailing slashes", "https://example.test/api/opds/test-placeholder///", "test-placeholder"},
		{"reverse proxy and encoded token", "https://example.test/books/api/opds/a%2Fb%3Fc%23d%25e%2Bf%20g/", "a/b?c#d%e+f g"},
		{"lowercase escapes", "http://example.test/api/opds/a%2fb", "a/b"},
		{"empty", "", ""},
		{"blank", " \t", ""},
		{"relative", "/api/opds/test-placeholder", ""},
		{"wrong scheme", "ftp://example.test/api/opds/test-placeholder", ""},
		{"missing token", "https://example.test/api/opds", ""},
		{"missing token trailing slash", "https://example.test/api/opds/", ""},
		{"root", "https://example.test/", ""},
		{"wrong route", "https://example.test/opds/test-placeholder", ""},
		{"extra route", "https://example.test/api/opds/test-placeholder/series", ""},
		{"malformed escape", "https://example.test/api/opds/test-placeholder%zz", ""},
		{"blank encoded token", "https://example.test/api/opds/%20%09", ""},
		{"control token", "https://example.test/api/opds/test-placeholder%0A", ""},
		{"query", "https://example.test/api/opds/test-placeholder?key=test-placeholder", ""},
		{"fragment", "https://example.test/api/opds/test-placeholder#test-placeholder", ""},
		{"userinfo", "https://user:test-placeholder@example.test/api/opds/test-placeholder", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := opdsAPIKey(tt.value)
			if tt.key == "" {
				if err == nil {
					t.Fatal("invalid OPDS setting was accepted")
				}
				if strings.Contains(err.Error(), "test-placeholder") || (tt.value != "" && strings.Contains(err.Error(), tt.value)) {
					t.Fatal("OPDS validation exposed the supplied value")
				}
				return
			}
			if err != nil || key != tt.key {
				t.Fatal("OPDS final path component was not decoded correctly")
			}
		})
	}
}

func TestLoadConfigEnvironmentAndFallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fixture.toml")
	legacy := "server_url = \"https://legacy.example.test\"\napi_key = \"legacy-placeholder\"\nopds_url = \"https://legacy.example.test/api/opds\"\n"
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfigFromPaths([]string{path}, envLookup(nil))
	if err != nil || cfg.APIKey != "legacy-placeholder" || cfg.OPDSURL != "https://legacy.example.test/api/opds" {
		t.Fatal("legacy configuration fallback failed")
	}
	for _, invalid := range []string{"", "https://example.test/api/opds/", "https://example.test/api/opds/test-placeholder%zz"} {
		cfg, err := loadConfigFromPaths([]string{path}, envLookup(map[string]string{"KAVITA_OPDS": invalid}))
		if err == nil || cfg != nil || strings.Contains(err.Error(), "test-placeholder") {
			t.Fatal("explicit invalid OPDS override did not fail closed safely")
		}
	}
	if _, err := loadConfigFromPaths([]string{path}, envLookup(map[string]string{"KAVITA_URL": ""})); err == nil {
		t.Fatal("explicit empty server override fell back to legacy settings")
	}
	complete := map[string]string{"KAVITA_URL": "https://env.example.test", "KAVITA_OPDS": "https://env.example.test/api/opds/env-placeholder"}
	cfg, err = loadConfigFromPaths(nil, envLookup(complete))
	if err != nil || cfg.APIKey != "env-placeholder" {
		t.Fatal("complete environment settings required a file")
	}
	if err := os.WriteFile(path, []byte("server_url = \"\"\napi_key = \"\"\nopds_url = \"https://legacy.example.test/api/opds\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = loadConfigFromPaths([]string{path}, envLookup(complete))
	if err != nil || cfg.ServerURL != complete["KAVITA_URL"] || cfg.APIKey != "env-placeholder" {
		t.Fatal("environment did not override invalid legacy fields")
	}
	// Partial OPDS overrides must replace an invalid legacy key and URL.
	if err := os.WriteFile(path, []byte("server_url = \"https://legacy.example.test\"\napi_key = \"\"\nopds_url = \"\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = loadConfigFromPaths([]string{path}, envLookup(map[string]string{"KAVITA_OPDS": complete["KAVITA_OPDS"]}))
	if err != nil || cfg.APIKey != "env-placeholder" || cfg.ServerURL != "https://legacy.example.test" {
		t.Fatal("OPDS-only override was validated too early")
	}
	// Partial overrides must also be applied before validating file fields.
	if err := os.WriteFile(path, []byte("server_url = \"\"\napi_key = \"legacy-placeholder\"\nopds_url = \"https://legacy.example.test/api/opds\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = loadConfigFromPaths([]string{path}, envLookup(map[string]string{"KAVITA_URL": complete["KAVITA_URL"]}))
	if err != nil || cfg.APIKey != "legacy-placeholder" {
		t.Fatal("partial override was validated too early")
	}
	if err := os.WriteFile(path, []byte("api_key = \"test-placeholder\" malformed"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = loadConfigFromPath(path)
	if err == nil || strings.Contains(err.Error(), "test-placeholder") {
		t.Fatal("TOML parse error was not secret-safe")
	}
}

func TestSafeOutput(t *testing.T) {
	const key = "test/key +placeholder"
	cfg := Config{ServerURL: "https://example.test", APIKey: key, OPDSURL: "https://example.test/api/opds/" + url.PathEscape(key) + "/"}
	var stdout, stderr bytes.Buffer
	out := newSafeOutput(&stdout, cfg)
	out.Printf("Failed: %v; OPDS: %s/series/1; Download: %s; metadata: %s %s\n", errors.New(key), cfg.OPDSURL, cfg.OPDSURL+"series/1/volume/2/chapter/3/download/book.epub", url.PathEscape(key), url.QueryEscape(key))
	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.SetOutput(newSafeOutput(&stderr, cfg))
	if err := flags.Parse([]string{"--" + key}); err == nil {
		t.Fatal("unknown flag was accepted")
	}
	for _, output := range []string{stdout.String(), stderr.String()} {
		for _, secret := range []string{key, url.PathEscape(key), url.QueryEscape(key), cfg.OPDSURL} {
			if strings.Contains(output, secret) {
				t.Fatal("stdout/stderr boundary exposed a credential or key-bearing URL")
			}
		}
		if !strings.Contains(output, "[REDACTED]") {
			t.Fatal("redaction marker missing")
		}
	}
	if !strings.Contains(stdout.String(), "Failed:") {
		t.Fatal("redaction discarded useful safe context")
	}
	// Sanitization precedes both truncation and ANSI highlighting.
	line := out.redact("prefix " + key + " suffix")
	out.Printf("%s", highlightMatch(line, "key"))
	if strings.Contains(stdout.String(), key) {
		t.Fatal("highlighted content exposed a credential")
	}
	// Even legacy OPDS URLs with a different embedded key must not be printed.
	legacy := newSafeOutput(&stdout, Config{APIKey: "different-placeholder", OPDSURL: "https://example.test/api/opds/embedded-placeholder"})
	legacy.Printf("%s", "embedded-placeholder")
	if strings.Contains(stdout.String(), "embedded-placeholder") {
		t.Fatal("legacy embedded OPDS credential was not redacted")
	}
}

func TestSafeOutputAlternateEncodings(t *testing.T) {
	const key = "fixture/key +placeholder"
	var encoded strings.Builder
	for i := 0; i < len(key); i++ {
		fmt.Fprintf(&encoded, "%%%02x", key[i])
	}
	tests := []string{
		"https://other.example.test/api/opds/" + encoded.String(),
		"%zz " + encoded.String(), // An unrelated malformed escape must not disable redaction.
		strings.ReplaceAll(url.QueryEscape(key), "%2F", "%2f"),
	}
	for _, text := range tests {
		var output bytes.Buffer
		out := newSafeOutput(&output, Config{APIKey: key})
		out.Printf("%s", text)
		if strings.Contains(output.String(), key) || output.String() == text || !strings.Contains(output.String(), "[REDACTED]") {
			t.Fatal("alternate credential encoding was not redacted")
		}
	}
	var output bytes.Buffer
	newSafeOutput(&output, Config{APIKey: key}).Printf("safe metadata 100%% %s", "normal title")
	if output.String() != "safe metadata 100% normal title" {
		t.Fatal("safe diagnostics were unnecessarily changed")
	}
}

func TestSafeOutputPropagatesWriteFailure(t *testing.T) {
	out := newSafeOutput(failingWriter{}, Config{})
	if _, err := fmt.Fprint(out, "safe text"); err == nil {
		t.Fatal("writer failure was swallowed")
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("fixture write failure") }

func TestValidateHTTPURLPreservesPaths(t *testing.T) {
	for _, value := range []string{
		"http://example.test",
		"https://example.test/books/proxy/",
		"https://example.test/books/a%3Fb%23c",
	} {
		if err := validateHTTPURL("fixture_url", value); err != nil {
			t.Fatal("absolute HTTP(S) URL with a proxy prefix or escaped path was rejected")
		}
	}
}

func TestURLSettingsRejectUnsafeComponents(t *testing.T) {
	tests := []struct {
		name, value string
	}{
		{"userinfo", "https://fixture-user:component-placeholder@example.test/books"},
		{"username only", "https://component-placeholder@example.test/books"},
		{"empty userinfo", "https://@example.test/books"},
		{"query", "https://example.test/books?key=component-placeholder"},
		{"empty forced query", "https://example.test/books?"},
		{"fragment", "https://example.test/books#component-placeholder"},
		{"empty fragment", "https://example.test/books#"},
	}
	for _, field := range []string{"KAVITA_URL", "server_url", "opds_url"} {
		for _, tt := range tests {
			t.Run(field+"/"+tt.name, func(t *testing.T) {
				cfg := Config{
					ServerURL: "https://example.test/books",
					APIKey:    "fixture-placeholder",
					OPDSURL:   "https://example.test/books/api/opds",
				}
				var err error
				switch field {
				case "KAVITA_URL":
					_, err = applyConfigEnv(cfg, envLookup(map[string]string{field: tt.value}))
				case "server_url":
					cfg.ServerURL = tt.value
					err = validateConfig(cfg)
				case "opds_url":
					cfg.OPDSURL = tt.value
					err = validateConfig(cfg)
				}
				// Exact field-only diagnostics must never echo any supplied URL component.
				want := field + " must be an absolute http or https URL without userinfo, query, or fragment"
				if err == nil || err.Error() != want {
					t.Fatal("unsafe URL was accepted or its diagnostic was not field-only")
				}
			})
		}
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
