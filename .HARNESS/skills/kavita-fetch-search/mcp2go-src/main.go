package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mcp2go/cache"
	"mcp2go/epub"
	"mcp2go/kavita"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	ServerURL string `toml:"server_url"`
	APIKey    string `toml:"api_key"`
	OPDSURL   string `toml:"opds_url"`
}

func loadConfig() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	// Try reading config.toml in executable directory first, then fallback to current directory
	configPaths := []string{
		filepath.Join(exeDir, "config.toml"),
		"config.toml",
	}

	var data []byte
	var readErr error
	for _, p := range configPaths {
		data, readErr = os.ReadFile(p)
		if readErr == nil {
			break
		}
	}
	if readErr != nil {
		return nil, fmt.Errorf("could not read config.toml: %w", readErr)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.toml: %w", err)
	}

	return &cfg, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: kavita_api.exe <series_name> <search_query>")
		os.Exit(1)
	}

	seriesName := os.Args[1]
	searchQuery := os.Args[2]

	fmt.Println("Loading config.toml...")
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Kavita Client
	fmt.Printf("Connecting to Kavita at %s...\n", cfg.ServerURL)
	client := kavita.NewClient(cfg.ServerURL)
	if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Authenticated successfully!")

	// Initialize SQLite Cache
	dbPath := "kavita_epub_cache.db"
	fmt.Printf("Opening cache database at %s...\n", dbPath)
	db, err := cache.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed to open cache database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Search for series
	fmt.Printf("Searching for series matching '%s'...\n", seriesName)
	seriesList, err := client.SearchSeries(seriesName)
	if err != nil {
		fmt.Printf("Failed to search series: %v\n", err)
		os.Exit(1)
	}

	if len(seriesList) == 0 {
		fmt.Println("No series found matching the query.")
		os.Exit(1)
	}

	fmt.Printf("Found %d series. Using the first match:\n", len(seriesList))
	targetSeries := seriesList[0]
	fmt.Printf("  Target Series: %s (ID: %d)\n", targetSeries.Name, targetSeries.SeriesID)

	// Fetch chapters
	fmt.Println("Fetching chapters...")
	chapters, err := client.GetChapters(targetSeries.SeriesID)
	if err != nil {
		fmt.Printf("Failed to get chapters: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d chapters.\n", len(chapters))

	// Download and cache chapters
	for _, chapter := range chapters {
		chapterTitle := chapter.Title
		if chapterTitle == "" {
			chapterTitle = fmt.Sprintf("Chapter ID %d", chapter.ID)
		}

		isCached, err := db.IsChapterCached(chapter.ID)
		if err != nil {
			fmt.Printf("Error checking cache for chapter %d: %v\n", chapter.ID, err)
			continue
		}

		if isCached {
			fmt.Printf("Chapter '%s' (ID: %d) is cached.\n", chapterTitle, chapter.ID)
			continue
		}

		fmt.Printf("Downloading chapter '%s' (ID: %d)...\n", chapterTitle, chapter.ID)
		epubBytes, err := client.DownloadChapter(chapter.ID)
		if err != nil {
			fmt.Printf("Failed to download chapter %d: %v\n", chapter.ID, err)
			continue
		}

		fmt.Printf("Parsing EPUB content for '%s'...\n", chapterTitle)
		parsedFiles, err := epub.ParseEPUBBytes(epubBytes)
		if err != nil {
			fmt.Printf("Failed to parse EPUB %d: %v\n", chapter.ID, err)
			continue
		}

		// Clear just in case before saving
		db.ClearChapterCache(chapter.ID)

		for _, pf := range parsedFiles {
			err := db.SaveFileCache(chapter.ID, targetSeries.SeriesID, targetSeries.Name, chapterTitle, pf.Name, pf.Text)
			if err != nil {
				fmt.Printf("Failed to save cache for file %s: %v\n", pf.Name, err)
			}
		}
		fmt.Printf("Cached %d files for chapter '%s'.\n", len(parsedFiles), chapterTitle)
	}

	// Perform full-text search on cached DB
	fmt.Printf("\nSearching database for '%s'...\n", searchQuery)
	results, err := db.SearchContent(targetSeries.SeriesID, searchQuery)
	if err != nil {
		fmt.Printf("Database search failed: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No matches found in the cached content.")
		return
	}

	fmt.Printf("Found matches in %d files:\n", len(results))

	// Display matches with context
	for _, res := range results {
		lines := strings.Split(res.Content, "\n")
		var matchedLinesCount int

		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(searchQuery)) {
				if matchedLinesCount == 0 {
					fmt.Printf("\n--- [%s] - %s (File: %s) ---\n", res.SeriesTitle, res.ChapterTitle, res.FilePath)
				}
				matchedLinesCount++

				// Context lines
				prevLine := ""
				if i > 0 {
					prevLine = strings.TrimSpace(lines[i-1])
				}
				nextLine := ""
				if i < len(lines)-1 {
					nextLine = strings.TrimSpace(lines[i+1])
				}

				if prevLine != "" {
					fmt.Printf("    - %s\n", prevLine)
				}

				// Highlight matching query
				highlighted := highlightMatch(strings.TrimSpace(line), searchQuery)
				fmt.Printf("    * %s\n", highlighted)

				if nextLine != "" {
					fmt.Printf("    - %s\n", nextLine)
				}
				fmt.Println()
			}
		}
	}
}

// highlightMatch wraps matched substring with yellow color ANSI code
func highlightMatch(line, query string) string {
	lowerLine := strings.ToLower(line)
	lowerQuery := strings.ToLower(query)

	var result strings.Builder
	idx := 0

	for {
		matchIdx := strings.Index(lowerLine[idx:], lowerQuery)
		if matchIdx == -1 {
			result.WriteString(line[idx:])
			break
		}

		actualMatchIdx := idx + matchIdx
		result.WriteString(line[idx:actualMatchIdx])
		
		// Yellow color ANSI sequence
		result.WriteString("\033[33m")
		result.WriteString(line[actualMatchIdx : actualMatchIdx+len(query)])
		result.WriteString("\033[0m")

		idx = actualMatchIdx + len(query)
	}

	return result.String()
}
