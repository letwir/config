package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VMClient handles communication with VictoriaMetrics API
type VMClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewVMClient creates a new VMClient with default timeout
func NewVMClient(rawURL string) (*VMClient, error) {
	cleanURL := strings.TrimRight(strings.TrimSpace(rawURL), "/")
	if cleanURL == "" {
		cleanURL = "http://docker.tigris-tailor.ts.net:8428"
	}
	if !strings.HasPrefix(cleanURL, "http://") && !strings.HasPrefix(cleanURL, "https://") {
		cleanURL = "http://" + cleanURL
	}

	parsed, err := url.Parse(cleanURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	return &VMClient{
		BaseURL: parsed.String(),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// HealthCheck verifies if VictoriaMetrics endpoint is accessible
func (c *VMClient) HealthCheck() (bool, error) {
	reqURL := fmt.Sprintf("%s/health", c.BaseURL)
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return false, fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("health check returned status %d", resp.StatusCode)
	}
	return true, nil
}

// GetTSDBStatus retrieves TSDB statistics
func (c *VMClient) GetTSDBStatus() (*TSDBStatusResponse, error) {
	reqURL := fmt.Sprintf("%s/api/v1/status/tsdb", c.BaseURL)
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("tsdb status request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tsdb response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tsdb status returned status %d: %s", resp.StatusCode, string(body))
	}

	var statusResp TSDBStatusResponse
	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tsdb status json: %w", err)
	}
	return &statusResp, nil
}

// GetMetricNames retrieves all registered metric names
func (c *VMClient) GetMetricNames() ([]string, error) {
	reqURL := fmt.Sprintf("%s/api/v1/label/__name__/values", c.BaseURL)
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("metric names request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metric names body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metric names query returned %d: %s", resp.StatusCode, string(body))
	}

	var labelResp LabelValuesResponse
	err = json.Unmarshal(body, &labelResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metric names json: %w", err)
	}
	return labelResp.Data, nil
}

// QueryInstant executes an instant MetricsQL/PromQL query
func (c *VMClient) QueryInstant(query string) (*QueryResponse, error) {
	params := url.Values{}
	params.Set("query", query)

	reqURL := fmt.Sprintf("%s/api/v1/query?%s", c.BaseURL, params.Encode())
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read query response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query returned status %d: %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	err = json.Unmarshal(body, &queryResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query json response: %w", err)
	}
	return &queryResp, nil
}

// QueryRange executes a range MetricsQL/PromQL query
func (c *VMClient) QueryRange(query string, start, end, step string) (*QueryResponse, error) {
	params := url.Values{}
	params.Set("query", query)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	if step != "" {
		params.Set("step", step)
	}

	reqURL := fmt.Sprintf("%s/api/v1/query_range?%s", c.BaseURL, params.Encode())
	resp, err := c.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("query range execution failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read query range response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query range returned status %d: %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	err = json.Unmarshal(body, &queryResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query range json: %w", err)
	}
	return &queryResp, nil
}
