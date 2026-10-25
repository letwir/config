package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

	configPaths := []string{
		filepath.Join(exeDir, "config.toml"),
		"config.toml",
	}

	for _, p := range configPaths {
		cfg, err := loadConfigFromPath(p)
		if err == nil {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("could not load config.toml from the executable directory or current directory")
}

func loadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config.toml: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.toml: %w", err)
	}
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg Config) error {
	if err := validateHTTPURL("server_url", cfg.ServerURL); err != nil {
		return err
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return fmt.Errorf("config.toml api_key must not be blank")
	}
	if err := validateHTTPURL("opds_url", cfg.OPDSURL); err != nil {
		return err
	}
	return nil
}

func validateHTTPURL(field, value string) error {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return fmt.Errorf("config.toml %s must be an absolute http or https URL", field)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed: Error loading config: %v\n", err)
		os.Exit(1)
	}

	client := kavita.NewClient(cfg.ServerURL)

	// Try opening local SQLite cache DB
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed: Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "kavita_epub_cache.db")
	fmt.Printf("Opening cache database at %s...\n", dbPath)
	db, err := cache.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed: Failed to open cache database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	switch subcommand {
	case "sync":
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			fmt.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}
		if err := syncDatabase(client, db); err != nil {
			fmt.Printf("Failed: %v\n", err)
			os.Exit(1)
		}

	case "search":
		// Handle search command flags
		searchFlags := flag.NewFlagSet("search", flag.ExitOnError)
		forceSync := searchFlags.Bool("sync", false, "Force sync with Kavita server before search")
		searchFlags.Parse(os.Args[2:])

		if searchFlags.NArg() < 1 {
			fmt.Println("Usage: kavita-fetch-search.exe search [--sync] <query>")
			os.Exit(1)
		}
		searchQuery := searchFlags.Arg(0)

		// Trigger automatic initial sync if local DB is completely empty
		seriesCount, err := db.GetSeriesCount()
		if err == nil && seriesCount == 0 {
			fmt.Println("Local cache is empty. Triggering automatic initial sync...")
			if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
				fmt.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db); err != nil {
				fmt.Printf("Failed: %v\n", err)
				os.Exit(1)
			}
		} else if *forceSync {
			fmt.Println("Force sync flag specified. Syncing metadata...")
			if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
				fmt.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db); err != nil {
				fmt.Printf("Failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Perform local search
		results, err := db.SearchSeriesLocal(searchQuery)
		if err != nil {
			fmt.Printf("Failed: Local metadata search failed: %v\n", err)
			os.Exit(1)
		}

		// If nothing found and we haven't synced just now, try sync once and search again (Self-Healing ETL)
		if len(results) == 0 && !*forceSync && (seriesCount > 0) {
			fmt.Println("No series found in local cache. Syncing with Kavita server...")
			if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
				fmt.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db); err == nil {
				results, _ = db.SearchSeriesLocal(searchQuery)
			} else {
				fmt.Printf("Failed: Sync failed during automatic search recovery: %v\n", err)
				os.Exit(1)
			}
		}

		if len(results) == 0 {
			fmt.Println("No series found matching the query.")
			return
		}

		// Build OPDS base URL (trim trailing slash for clean concatenation)
		opdsBase := strings.TrimRight(cfg.OPDSURL, "/")

		fmt.Printf("Found %d series in local cache:\n", len(results))
		for _, s := range results {
			opdsURL := fmt.Sprintf("%s/series/%d", opdsBase, s.SeriesID)
			fmt.Printf("  ID: %d | Title: %s | Author: %s\n", s.SeriesID, s.Name, s.Author)
			fmt.Printf("    OPDS: %s\n", opdsURL)
			if s.Summary != "" {
				// Truncate very long summaries for display (max ~300 chars)
				summary := s.Summary
				if len(summary) > 300 {
					summary = summary[:300] + "..."
				}
				fmt.Printf("    Summary: %s\n", summary)
			}
		}

	case "query":
		if len(os.Args) < 4 {
			fmt.Println("Usage: kavita-fetch-search.exe query <series_id> <search_query>")
			os.Exit(1)
		}

		seriesID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Failed: Invalid series_id: %s\n", os.Args[2])
			os.Exit(1)
		}
		searchQuery := os.Args[3]

		// Authenticate first, required to check and fetch new EPUB files
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			fmt.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}

		// Fetch chapters from server
		chapters, err := client.GetChapters(seriesID)
		if err != nil {
			fmt.Printf("Failed: Failed to fetch chapters from server: %v\n", err)
			os.Exit(1)
		}

		// Sync texts for untracked chapters (individual ETL)
		for _, chapter := range chapters {
			chapterTitle := chapter.Title
			if chapterTitle == "" {
				chapterTitle = fmt.Sprintf("Chapter ID %d", chapter.ID)
			}

			isCached, err := db.IsChapterCached(chapter.ID)
			if err != nil {
				fmt.Printf("Warning: Failed to check cache for chapter %d: %v\n", chapter.ID, err)
				continue
			}

			if isCached {
				continue
			}

			fmt.Printf("Downloading chapter '%s' (ID: %d) to cache...\n", chapterTitle, chapter.ID)
			epubBytes, err := client.DownloadChapter(chapter.ID)
			if err != nil {
				fmt.Printf("Failed: Failed to download chapter %d: %v\n", chapter.ID, err)
				os.Exit(1) // Connection error during download should fail
			}

			parsedFiles, err := epub.ParseEPUBBytes(epubBytes)
			if err != nil {
				fmt.Printf("Warning: Failed to parse EPUB for chapter %d: %v\n", chapter.ID, err)
				continue
			}

			db.ClearChapterCache(chapter.ID)
			for _, pf := range parsedFiles {
				// We don't have target series title here, retrieve it or pass empty since search query outputs will show series title from cached series DB table join
				err := db.SaveFileCache(chapter.ID, seriesID, "", chapterTitle, pf.Name, pf.Text)
				if err != nil {
					fmt.Printf("Warning: Failed to save file cache for %s: %v\n", pf.Name, err)
				}
			}
		}

		// Perform full-text search on cached DB
		fmt.Printf("Searching text for '%s'...\n", searchQuery)
		results, err := db.SearchContent(seriesID, searchQuery)
		if err != nil {
			fmt.Printf("Failed: Full-text search failed: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Println("No matches found in the cached content.")
			return
		}

		// Get Series title from DB series table for display
		seriesTitle := fmt.Sprintf("Series ID %d", seriesID)
		localSeries, err := db.SearchSeriesLocal(strconv.Itoa(seriesID))
		if err == nil && len(localSeries) > 0 {
			for _, ls := range localSeries {
				if ls.SeriesID == seriesID {
					seriesTitle = ls.Name
					break
				}
			}
		}

		fmt.Printf("Found matches in %d files:\n", len(results))
		for _, res := range results {
			lines := strings.Split(res.Content, "\n")
			var matchedLinesCount int

			for i, line := range lines {
				if strings.Contains(strings.ToLower(line), strings.ToLower(searchQuery)) {
					if matchedLinesCount == 0 {
						fmt.Printf("\n--- [%s] - %s (File: %s) ---\n", seriesTitle, res.ChapterTitle, res.FilePath)
					}
					matchedLinesCount++

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

					highlighted := highlightMatch(strings.TrimSpace(line), searchQuery)
					fmt.Printf("    * %s\n", highlighted)

					if nextLine != "" {
						fmt.Printf("    - %s\n", nextLine)
					}
					fmt.Println()
				}
			}
		}

	case "dl":
		if len(os.Args) < 3 {
			fmt.Println("Usage: kavita-fetch-search.exe dl <series_id>")
			os.Exit(1)
		}

		seriesID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Failed: Invalid series_id: %s\n", os.Args[2])
			os.Exit(1)
		}

		// Authenticate to access REST API
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			fmt.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}

		opdsBase := strings.TrimRight(cfg.OPDSURL, "/")

		// Fetch volumes with chapters and file info
		volumes, err := client.GetSeriesVolumes(seriesID)
		if err != nil {
			fmt.Printf("Failed: Failed to fetch volumes: %v\n", err)
			os.Exit(1)
		}

		if len(volumes) == 0 {
			fmt.Println("No volumes found for this series.")
			return
		}

		fmt.Printf("Series ID: %d — EPUB Download URLs (curl-ready):\n\n", seriesID)

		for _, vol := range volumes {
			volName := vol.Name
			if volName == "" {
				volName = fmt.Sprintf("Volume %d", vol.ID)
			}

			for _, ch := range vol.Chapters {
				chTitle := ch.Title
				if chTitle == "" {
					chTitle = fmt.Sprintf("Chapter %d", ch.ID)
				}

				if len(ch.Files) == 0 {
					// If no file info returned, construct URL with generic filename
					dlURL := kavita.BuildOPDSDownloadURL(opdsBase, seriesID, vol.ID, ch.ID, "book.epub")
					fmt.Printf("  [%s / %s]\n", volName, chTitle)
					fmt.Printf("    curl -o \"%s.epub\" \"%s\"\n\n", sanitizeFilename(chTitle), dlURL)
				} else {
					for _, f := range ch.Files {
						dlURL := kavita.BuildOPDSDownloadURL(opdsBase, seriesID, vol.ID, ch.ID, f.FilePath)
						fmt.Printf("  [%s / %s]\n", volName, chTitle)
						fmt.Printf("    curl -o \"%s\" \"%s\"\n\n", sanitizeFilename(filepath.Base(f.FilePath)), dlURL)
					}
				}
			}
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  kavita-fetch-search.exe search [--sync] <query>  : Search series/author locally (triggers ETL if empty)")
	fmt.Println("                                                     Results include OPDS URL, summary, and author info.")
	fmt.Println("  kavita-fetch-search.exe query <series_id> <q>    : Perform full-text search in a series (pulls missing EPUBs)")
	fmt.Println("  kavita-fetch-search.exe dl <series_id>           : Show curl-ready OPDS download URLs for all EPUBs in a series")
	fmt.Println("  kavita-fetch-search.exe sync                     : Force ETL sync of all metadata from Kavita server")
}

func syncDatabase(client *kavita.Client, db *cache.DB) error {
	fmt.Println("Syncing Kavita metadata to local database...")
	seriesList, err := client.GetAllSeries()
	if err != nil {
		return fmt.Errorf("failed to get all series from Kavita: %w", err)
	}

	fmt.Printf("Syncing %d series and fetching metadata details...\n", len(seriesList))
	for _, series := range seriesList {
		meta, err := client.GetSeriesMetadata(series.ID)
		author := "Unknown"
		summary := ""
		if err == nil {
			if len(meta.Writers) > 0 {
				var writers []string
				for _, w := range meta.Writers {
					writers = append(writers, w.Name)
				}
				author = strings.Join(writers, ", ")
			}
			summary = meta.Summary
		}

		err = db.SaveSeriesCache(series.ID, series.Name, author, summary)
		if err != nil {
			fmt.Printf("Warning: failed to save cache for series %s: %v\n", series.Name, err)
		}
	}
	fmt.Println("Sync completed successfully!")
	return nil
}

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

		result.WriteString("\033[33m")
		result.WriteString(line[actualMatchIdx : actualMatchIdx+len(query)])
		result.WriteString("\033[0m")

		idx = actualMatchIdx + len(query)
	}

	return result.String()
}

// sanitizeFilename replaces characters unsafe for filenames
func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}
