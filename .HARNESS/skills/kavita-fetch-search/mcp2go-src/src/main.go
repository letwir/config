package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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

	return loadConfigFromPaths(configPaths, os.LookupEnv)
}

// Resolve explicit overrides first, so invalid environment settings cannot
// silently fall back to legacy credentials. Complete environment settings need no file.
func loadConfigFromPaths(paths []string, lookup func(string) (string, bool)) (*Config, error) {
	envConfig, err := applyConfigEnv(Config{}, lookup)
	if err != nil {
		return nil, err
	}
	if validateConfig(envConfig) == nil {
		return &envConfig, nil
	}
	for _, path := range paths {
		cfg, err := readConfigFromPath(path)
		if err != nil {
			continue
		}
		merged, err := applyConfigEnv(*cfg, lookup)
		if err != nil {
			return nil, err
		}
		if validateConfig(merged) == nil {
			return &merged, nil
		}
	}
	return nil, fmt.Errorf("could not load valid Kavita settings from environment or config.toml")
}

func applyConfigEnv(cfg Config, lookup func(string) (string, bool)) (Config, error) {
	if value, set := lookup("KAVITA_URL"); set {
		if err := validateHTTPURL("KAVITA_URL", value); err != nil {
			return Config{}, err
		}
		cfg.ServerURL = strings.TrimSpace(value)
	}
	if value, set := lookup("KAVITA_OPDS"); set {
		key, err := opdsAPIKey(value)
		if err != nil {
			return Config{}, err
		}
		cfg.OPDSURL = strings.TrimSpace(value)
		cfg.APIKey = key
	}
	return cfg, nil
}

func opdsAPIKey(value string) (string, error) {
	invalid := fmt.Errorf("KAVITA_OPDS must be an absolute http or https URL ending in /api/opds/<API key>")
	if validateHTTPURL("KAVITA_OPDS", value) != nil {
		return "", invalid
	}
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil {
		return "", invalid
	}
	if parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.User != nil {
		return "", invalid
	}
	parts := strings.Split(strings.TrimRight(parsed.EscapedPath(), "/"), "/")
	if len(parts) < 4 || parts[len(parts)-3] != "api" || parts[len(parts)-2] != "opds" {
		return "", invalid
	}
	key, err := url.PathUnescape(parts[len(parts)-1])
	if err != nil || strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n\x00") {
		return "", invalid
	}
	return key, nil
}

func loadConfigFromPath(path string) (*Config, error) {
	cfg, err := readConfigFromPath(path)
	if err != nil {
		return nil, err
	}
	if err := validateConfig(*cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Read unvalidated fields so environment overrides can replace invalid values.
func readConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config.toml")
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		// Parser errors may quote a credential-bearing source line.
		return nil, fmt.Errorf("failed to parse config.toml")
	}
	return &cfg, nil
}

func validateConfig(cfg Config) error {
	if err := validateHTTPURL("server_url", cfg.ServerURL); err != nil {
		return err
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return fmt.Errorf("Kavita API key must not be blank")
	}
	if err := validateHTTPURL("opds_url", cfg.OPDSURL); err != nil {
		return err
	}
	return nil
}

func validateHTTPURL(field, value string) error {
	parsed, err := url.Parse(strings.TrimSpace(value))
	// url.Parse does not retain an empty fragment delimiter, so check it directly.
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.ForceQuery || strings.Contains(value, "#") {
		return fmt.Errorf("%s must be an absolute http or https URL without userinfo, query, or fragment", field)
	}
	return nil
}

// safeOutput is the final display boundary for metadata, local errors and flag
// diagnostics. Transport errors are separately sanitized in the HTTP client.
type safeOutput struct {
	writer   io.Writer
	replacer *strings.Replacer
}

func newSafeOutput(writer io.Writer, cfg Config) *safeOutput {
	secrets := []string{cfg.OPDSURL, strings.TrimRight(cfg.OPDSURL, "/"), cfg.ServerURL}
	keys := []string{cfg.APIKey}
	if key, err := opdsAPIKey(cfg.OPDSURL); err == nil {
		keys = append(keys, key)
		if parsed, err := url.Parse(cfg.OPDSURL); err == nil {
			parts := strings.Split(strings.TrimRight(parsed.EscapedPath(), "/"), "/")
			keys = append(keys, parts[len(parts)-1])
		}
	}
	for _, key := range keys {
		secrets = append(secrets, key, url.PathEscape(key), url.QueryEscape(key))
	}
	// Replace whole URLs before shorter credentials contained inside them.
	sort.SliceStable(secrets, func(i, j int) bool { return len(secrets[i]) > len(secrets[j]) })
	var pairs []string
	for _, secret := range secrets {
		if secret != "" {
			pairs = append(pairs, secret, "[REDACTED]")
		}
	}
	return &safeOutput{writer: writer, replacer: strings.NewReplacer(pairs...)}
}

// Alternate percent encodings may occur in echoed URLs or metadata. Detect
// those too, without converting arbitrary search text into decoded output.
func (out *safeOutput) redact(value string) string {
	redacted := out.replacer.Replace(value)
	for _, candidate := range []string{redacted, strings.ReplaceAll(redacted, "+", " ")} {
		decoded := decodeDisplayEscapes(candidate)
		if out.replacer.Replace(decoded) != decoded {
			return "[REDACTED]"
		}
	}
	return redacted
}

func decodeDisplayEscapes(value string) string {
	var decoded strings.Builder
	for i := 0; i < len(value); i++ {
		if value[i] == '%' && i+2 < len(value) {
			if part, err := url.PathUnescape(value[i : i+3]); err == nil {
				decoded.WriteString(part)
				i += 2
				continue
			}
		}
		decoded.WriteByte(value[i])
	}
	return decoded.String()
}

func (out *safeOutput) Write(data []byte) (int, error) {
	_, err := io.WriteString(out.writer, out.redact(string(data)))
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (out *safeOutput) Printf(format string, args ...any) {
	fmt.Fprintf(out, format, args...)
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

	out := newSafeOutput(os.Stdout, *cfg)
	client := kavita.NewClient(cfg.ServerURL)

	// Try opening local SQLite cache DB
	exePath, err := os.Executable()
	if err != nil {
		out.Printf("Failed: Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "kavita_epub_cache.db")
	out.Printf("Opening cache database at %s...\n", dbPath)
	db, err := cache.Open(dbPath)
	if err != nil {
		out.Printf("Failed: Failed to open cache database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	switch subcommand {
	case "sync":
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			out.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}
		if err := syncDatabase(client, db, out); err != nil {
			out.Printf("Failed: %v\n", err)
			os.Exit(1)
		}

	case "search":
		// Handle search command flags
		searchFlags := flag.NewFlagSet("search", flag.ExitOnError)
		searchFlags.SetOutput(newSafeOutput(os.Stderr, *cfg))
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
				out.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db, out); err != nil {
				out.Printf("Failed: %v\n", err)
				os.Exit(1)
			}
		} else if *forceSync {
			fmt.Println("Force sync flag specified. Syncing metadata...")
			if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
				out.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db, out); err != nil {
				out.Printf("Failed: %v\n", err)
				os.Exit(1)
			}
		}

		// Perform local search
		results, err := db.SearchSeriesLocal(searchQuery)
		if err != nil {
			out.Printf("Failed: Local metadata search failed: %v\n", err)
			os.Exit(1)
		}

		// If nothing found and we haven't synced just now, try sync once and search again (Self-Healing ETL)
		if len(results) == 0 && !*forceSync && (seriesCount > 0) {
			fmt.Println("No series found in local cache. Syncing with Kavita server...")
			if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
				out.Printf("Failed: Authentication failed: %v\n", err)
				os.Exit(1)
			}
			if err := syncDatabase(client, db, out); err == nil {
				results, _ = db.SearchSeriesLocal(searchQuery)
			} else {
				out.Printf("Failed: Sync failed during automatic search recovery: %v\n", err)
				os.Exit(1)
			}
		}

		if len(results) == 0 {
			fmt.Println("No series found matching the query.")
			return
		}

		fmt.Printf("Found %d series in local cache:\n", len(results))
		for _, s := range results {
			out.Printf("  ID: %d | Title: %s | Author: %s\n", s.SeriesID, s.Name, s.Author)
			fmt.Println("    OPDS: [credential-bearing URL omitted]")
			if s.Summary != "" {
				// Truncate very long summaries for display (max ~300 chars)
				summary := out.redact(s.Summary)
				if len(summary) > 300 {
					summary = summary[:300] + "..."
				}
				out.Printf("    Summary: %s\n", summary)
			}
		}

	case "query":
		if len(os.Args) < 4 {
			fmt.Println("Usage: kavita-fetch-search.exe query <series_id> <search_query>")
			os.Exit(1)
		}

		seriesID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			out.Printf("Failed: Invalid series_id: %s\n", os.Args[2])
			os.Exit(1)
		}
		searchQuery := os.Args[3]

		// Authenticate first, required to check and fetch new EPUB files
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			out.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}

		// Fetch chapters from server
		chapters, err := client.GetChapters(seriesID)
		if err != nil {
			out.Printf("Failed: Failed to fetch chapters from server: %v\n", err)
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
				out.Printf("Warning: Failed to check cache for chapter %d: %v\n", chapter.ID, err)
				continue
			}

			if isCached {
				continue
			}

			out.Printf("Downloading chapter '%s' (ID: %d) to cache...\n", chapterTitle, chapter.ID)
			epubBytes, err := client.DownloadChapter(chapter.ID)
			if err != nil {
				out.Printf("Failed: Failed to download chapter %d: %v\n", chapter.ID, err)
				os.Exit(1) // Connection error during download should fail
			}

			parsedFiles, err := epub.ParseEPUBBytes(epubBytes)
			if err != nil {
				out.Printf("Warning: Failed to parse EPUB for chapter %d: %v\n", chapter.ID, err)
				continue
			}

			db.ClearChapterCache(chapter.ID)
			for _, pf := range parsedFiles {
				// We don't have target series title here, retrieve it or pass empty since search query outputs will show series title from cached series DB table join
				err := db.SaveFileCache(chapter.ID, seriesID, "", chapterTitle, pf.Name, pf.Text)
				if err != nil {
					out.Printf("Warning: Failed to save file cache for %s: %v\n", pf.Name, err)
				}
			}
		}

		// Perform full-text search on cached DB
		out.Printf("Searching text for '%s'...\n", searchQuery)
		results, err := db.SearchContent(seriesID, searchQuery)
		if err != nil {
			out.Printf("Failed: Full-text search failed: %v\n", err)
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
						out.Printf("\n--- [%s] - %s (File: %s) ---\n", seriesTitle, res.ChapterTitle, res.FilePath)
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
						out.Printf("    - %s\n", prevLine)
					}

					// Redact before ANSI highlighting can split a credential.
					highlighted := highlightMatch(out.redact(strings.TrimSpace(line)), out.redact(searchQuery))
					out.Printf("    * %s\n", highlighted)

					if nextLine != "" {
						out.Printf("    - %s\n", nextLine)
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
			out.Printf("Failed: Invalid series_id: %s\n", os.Args[2])
			os.Exit(1)
		}

		// Authenticate to access REST API
		if err := client.Authenticate(cfg.APIKey, "KavitaGoSearchCLI"); err != nil {
			out.Printf("Failed: Authentication failed: %v\n", err)
			os.Exit(1)
		}

		// Fetch volumes with chapters and file info
		volumes, err := client.GetSeriesVolumes(seriesID)
		if err != nil {
			out.Printf("Failed: Failed to fetch volumes: %v\n", err)
			os.Exit(1)
		}

		if len(volumes) == 0 {
			fmt.Println("No volumes found for this series.")
			return
		}

		out.Printf("Series ID: %d — EPUB files (credential-bearing download URLs omitted):\n\n", seriesID)

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
					out.Printf("  [%s / %s] (chapter ID: %d)\n", volName, chTitle, ch.ID)
					out.Printf("    File: %s.epub | Download URL: [omitted]\n\n", sanitizeFilename(chTitle))
				} else {
					for _, f := range ch.Files {
						out.Printf("  [%s / %s] (chapter ID: %d)\n", volName, chTitle, ch.ID)
						out.Printf("    File: %s | Download URL: [omitted]\n\n", sanitizeFilename(filepath.Base(f.FilePath)))
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
	fmt.Println("                                                     Results include series IDs, summary, and author info; credential-bearing URLs are omitted.")
	fmt.Println("  kavita-fetch-search.exe query <series_id> <q>    : Perform full-text search in a series (pulls missing EPUBs)")
	fmt.Println("  kavita-fetch-search.exe dl <series_id>           : Show EPUB file information; credential-bearing download URLs are omitted")
	fmt.Println("  kavita-fetch-search.exe sync                     : Force ETL sync of all metadata from Kavita server")
}

func syncDatabase(client *kavita.Client, db *cache.DB, out *safeOutput) error {
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
			out.Printf("Warning: failed to save cache for series %s: %v\n", series.Name, err)
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
