package kavita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL string
	Token   string
}

type UserDto struct {
	Token string `json:"token"`
}

type SearchResultGroupDto struct {
	Series []SearchResultDto `json:"series"`
}

type SearchResultDto struct {
	SeriesID int    `json:"seriesId"`
	Name     string `json:"name"`
}

type ChapterDto struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type SeriesDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PersonDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SeriesMetadataDto struct {
	ID      int         `json:"id"`
	Summary string      `json:"summary"`
	Writers []PersonDto `json:"writers"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
	}
}

// Authenticate plugin with API key
func (c *Client) Authenticate(apiKey string, pluginName string) error {
	authURL := fmt.Sprintf("%s/api/Plugin/authenticate?apiKey=%s&pluginName=%s", c.BaseURL, url.QueryEscape(apiKey), url.QueryEscape(pluginName))
	
	resp, err := http.Post(authURL, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to make authenticate request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var user UserDto
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	if user.Token == "" {
		return fmt.Errorf("received empty token in auth response")
	}

	c.Token = user.Token
	return nil
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

// SearchSeries searches for series matching query (Kavita search API)
func (c *Client) SearchSeries(query string) ([]SearchResultDto, error) {
	searchPath := fmt.Sprintf("/api/Search/search?queryString=%s", url.QueryEscape(query))
	req, err := c.newRequest("GET", searchPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned status %d", resp.StatusCode)
	}

	var result SearchResultGroupDto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return result.Series, nil
}

// GetAllSeries fetches all series from library (using all-v2 API with large PageSize)
func (c *Client) GetAllSeries() ([]SeriesDto, error) {
	allPath := "/api/Series/all-v2?PageNumber=1&PageSize=10000"
	
	// Post requires a filter DTO body. Passing empty JSON object {} represents no filters.
	req, err := c.newRequest("POST", allPath, bytes.NewBufferString("{}"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get all series request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get all series returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var series []SeriesDto
	if err := json.NewDecoder(resp.Body).Decode(&series); err != nil {
		return nil, fmt.Errorf("failed to decode series response: %w", err)
	}

	return series, nil
}

// GetSeriesMetadata fetches detailed metadata including writers for a series
func (c *Client) GetSeriesMetadata(seriesID int) (*SeriesMetadataDto, error) {
	metaPath := fmt.Sprintf("/api/Series/metadata?seriesId=%d", seriesID)
	req, err := c.newRequest("GET", metaPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get series metadata request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get series metadata returned status %d", resp.StatusCode)
	}

	var meta SeriesMetadataDto
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("failed to decode metadata response: %w", err)
	}

	return &meta, nil
}

// GetChapters returns all chapters for a given series
func (c *Client) GetChapters(seriesID int) ([]ChapterDto, error) {
	chaptersPath := fmt.Sprintf("/api/Search/chapters-by-series?seriesId=%d", seriesID)
	req, err := c.newRequest("GET", chaptersPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get chapters request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get chapters returned status %d", resp.StatusCode)
	}

	var chapters []ChapterDto
	if err := json.NewDecoder(resp.Body).Decode(&chapters); err != nil {
		return nil, fmt.Errorf("failed to decode chapters response: %w", err)
	}

	return chapters, nil
}

// DownloadChapter downloads epub file for a given chapter
func (c *Client) DownloadChapter(chapterID int) ([]byte, error) {
	dlPath := fmt.Sprintf("/api/Download/chapter?chapterId=%d", chapterID)
	req, err := c.newRequest("GET", dlPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download chapter request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download chapter returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}
