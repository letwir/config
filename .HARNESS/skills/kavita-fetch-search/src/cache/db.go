package cache

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"time"
)

type CachedSeries struct {
	SeriesID  int
	Name      string
	Author    string
	UpdatedAt time.Time
}

type CachedFile struct {
	ID           int
	ChapterID    int
	SeriesID     int
	SeriesTitle  string
	ChapterTitle string
	FilePath     string
	Content      string
	CachedAt     time.Time
}

type DB struct {
	conn *sql.DB
}

func Open(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS series (
		series_id INTEGER PRIMARY KEY,
		name TEXT,
		author TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_series_name ON series(name);
	CREATE INDEX IF NOT EXISTS idx_series_author ON series(author);

	CREATE TABLE IF NOT EXISTS cached_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chapter_id INTEGER,
		series_id INTEGER,
		series_title TEXT,
		chapter_title TEXT,
		file_path TEXT,
		content TEXT,
		cached_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_cached_files_chapter_id ON cached_files(chapter_id);
	CREATE INDEX IF NOT EXISTS idx_cached_files_series_id ON cached_files(series_id);
	`
	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

// GetSeriesCount returns the number of series in the local cache
func (db *DB) GetSeriesCount() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM series").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SaveSeriesCache inserts or updates a series metadata in local cache
func (db *DB) SaveSeriesCache(seriesID int, name, author string) error {
	query := `
	INSERT OR REPLACE INTO series (series_id, name, author, updated_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := db.conn.Exec(query, seriesID, name, author)
	return err
}

// SearchSeriesLocal searches local cache for series matching name or author
func (db *DB) SearchSeriesLocal(query string) ([]CachedSeries, error) {
	sqlQuery := `
	SELECT series_id, name, author, updated_at
	FROM series
	WHERE name LIKE ? OR author LIKE ?
	`
	likePattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, likePattern, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []CachedSeries
	for rows.Next() {
		var cs CachedSeries
		var updatedAtStr string
		err := rows.Scan(&cs.SeriesID, &cs.Name, &cs.Author, &updatedAtStr)
		if err != nil {
			return nil, err
		}
		cs.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		list = append(list, cs)
	}

	return list, nil
}

// IsChapterCached checks if a chapter's files are already cached
func (db *DB) IsChapterCached(chapterID int) (bool, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM cached_files WHERE chapter_id = ?", chapterID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ClearChapterCache deletes cached files for a chapter
func (db *DB) ClearChapterCache(chapterID int) error {
	_, err := db.conn.Exec("DELETE FROM cached_files WHERE chapter_id = ?", chapterID)
	return err
}

// SaveFileCache saves a single file content of a chapter
func (db *DB) SaveFileCache(chapterID int, seriesID int, seriesTitle, chapterTitle, filePath, content string) error {
	query := `
	INSERT INTO cached_files (chapter_id, series_id, series_title, chapter_title, file_path, content)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query, chapterID, seriesID, seriesTitle, chapterTitle, filePath, content)
	return err
}

// SearchContent searches for text content across cached files for a specific series
func (db *DB) SearchContent(seriesID int, query string) ([]CachedFile, error) {
	sqlQuery := `
	SELECT c.id, c.chapter_id, c.series_id, COALESCE(NULLIF(c.series_title, ''), s.name, 'Series ID ' || c.series_id), c.chapter_title, c.file_path, c.content, c.cached_at
	FROM cached_files c
	LEFT JOIN series s ON c.series_id = s.series_id
	WHERE c.content LIKE ? AND c.series_id = ?
	`
	likePattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, likePattern, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CachedFile
	for rows.Next() {
		var cf CachedFile
		var cachedAtStr string
		err := rows.Scan(
			&cf.ID,
			&cf.ChapterID,
			&cf.SeriesID,
			&cf.SeriesTitle,
			&cf.ChapterTitle,
			&cf.FilePath,
			&cf.Content,
			&cachedAtStr,
		)
		if err != nil {
			return nil, err
		}
		
		cf.CachedAt, _ = time.Parse("2006-01-02 15:04:05", cachedAtStr)
		results = append(results, cf)
	}

	return results, nil
}

// SearchAllContent searches for text content across all cached files in all series
func (db *DB) SearchAllContent(query string) ([]CachedFile, error) {
	sqlQuery := `
	SELECT c.id, c.chapter_id, c.series_id, COALESCE(NULLIF(c.series_title, ''), s.name, 'Series ID ' || c.series_id), c.chapter_title, c.file_path, c.content, c.cached_at
	FROM cached_files c
	LEFT JOIN series s ON c.series_id = s.series_id
	WHERE c.content LIKE ?
	`
	likePattern := "%" + query + "%"
	rows, err := db.conn.Query(sqlQuery, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CachedFile
	for rows.Next() {
		var cf CachedFile
		var cachedAtStr string
		err := rows.Scan(
			&cf.ID,
			&cf.ChapterID,
			&cf.SeriesID,
			&cf.SeriesTitle,
			&cf.ChapterTitle,
			&cf.FilePath,
			&cf.Content,
			&cachedAtStr,
		)
		if err != nil {
			return nil, err
		}
		
		cf.CachedAt, _ = time.Parse("2006-01-02 15:04:05", cachedAtStr)
		results = append(results, cf)
	}

	return results, nil
}

