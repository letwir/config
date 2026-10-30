package kavita

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// In-memory transports only: no listeners, sockets, process environment or API.
type fixtureTransport func(*http.Request) (*http.Response, error)

func (f fixtureTransport) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func useFixtureTransport(t *testing.T, transport fixtureTransport) {
	t.Helper()
	previous := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: transport}
	t.Cleanup(func() { http.DefaultClient = previous })
}

type fixtureBody struct {
	io.Reader
	closed bool
}

func (b *fixtureBody) Close() error { b.closed = true; return nil }

type failingBody struct{}

func (failingBody) Read([]byte) (int, error) {
	return 0, errors.New("fixture-key-placeholder https://example.test/api/opds/fixture-key-placeholder")
}
func (failingBody) Close() error { return nil }

func TestClientErrorsNeverExposeRequestOrResponseSecrets(t *testing.T) {
	const key = "fixture-key-placeholder"
	const base = "https://example.test/api/opds/" + key
	operations := []struct {
		name string
		call func(*Client) error
	}{
		{"authenticate", func(c *Client) error { return c.Authenticate(key, "test") }},
		{"search", func(c *Client) error { _, err := c.SearchSeries(key); return err }},
		{"all series", func(c *Client) error { _, err := c.GetAllSeries(); return err }},
		{"metadata", func(c *Client) error { _, err := c.GetSeriesMetadata(1); return err }},
		{"chapters", func(c *Client) error { _, err := c.GetChapters(1); return err }},
		{"volumes", func(c *Client) error { _, err := c.GetSeriesVolumes(1); return err }},
		{"download", func(c *Client) error { _, err := c.DownloadChapter(1); return err }},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, mode := range []string{"transport", "status", "decode or read", "request construction"} {
				t.Run(mode, func(t *testing.T) {
					var body *fixtureBody
					useFixtureTransport(t, func(req *http.Request) (*http.Response, error) {
						if mode == "transport" {
							return nil, &url.Error{Op: "fixture", URL: req.URL.String(), Err: errors.New(key)}
						}
						status := http.StatusOK
						var reader io.Reader = failingBody{}
						if mode == "status" {
							status = http.StatusUnauthorized
							reader = strings.NewReader(key + " " + req.URL.String())
						}
						body = &fixtureBody{Reader: reader}
						return &http.Response{StatusCode: status, Body: body, Header: make(http.Header), Request: req}, nil
					})
					client := NewClient(base)
					client.Token = key
					if mode == "request construction" {
						client.BaseURL = "https://example.test/%zz/" + key
					}
					err := operation.call(client)
					if err == nil {
						t.Fatal("expected a safe error")
					}
					for _, secret := range []string{key, base, "https://"} {
						if strings.Contains(err.Error(), secret) {
							t.Fatal("client error exposed a credential or request URL")
						}
					}
					if mode == "status" && !strings.Contains(err.Error(), "401") {
						t.Fatal("safe HTTP status context was lost")
					}
					if body != nil && !body.closed {
						t.Fatal("response body was not closed")
					}
				})
			}
		})
	}
}

func TestAuthenticateEscapesDecodedAPIKey(t *testing.T) {
	const key = "fixture/key +?&#placeholder"
	var body *fixtureBody
	useFixtureTransport(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("apiKey") != key || req.URL.Query().Get("pluginName") != "fixture plugin" {
			t.Fatal("authentication query did not preserve the decoded key or plugin name")
		}
		body = &fixtureBody{Reader: strings.NewReader(`{"token":"fixture-bearer-placeholder"}`)}
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header), Request: req}, nil
	})
	client := NewClient("https://example.test")
	if err := client.Authenticate(key, "fixture plugin"); err != nil {
		t.Fatal("fixture authentication failed")
	}
	if client.Token != "fixture-bearer-placeholder" || body == nil || !body.closed {
		t.Fatal("authentication did not retain the token and close the body")
	}
}

func TestResponseDecodeErrorsAreSecretSafe(t *testing.T) {
	useFixtureTransport(t, func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"token":{"fixture-key-placeholder":"https://example.test/api/opds/fixture-key-placeholder"}}`)), Header: make(http.Header), Request: req}, nil
	})
	err := NewClient("https://example.test").Authenticate("fixture-key-placeholder", "test")
	if err == nil || strings.Contains(err.Error(), "fixture-key-placeholder") || strings.Contains(err.Error(), "https://") {
		t.Fatal("JSON decode failure was not secret-safe")
	}
}
