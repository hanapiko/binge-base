package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}

	// Initialize tables
	if err := database.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Printf("✅ Database initialized successfully at %s", dbPath)
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// initTables creates the necessary tables if they don't exist
func (d *Database) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tmdb_id INTEGER UNIQUE NOT NULL,
			title TEXT NOT NULL,
			overview TEXT,
			poster_path TEXT,
			backdrop_path TEXT,
			release_date TEXT,
			vote_average REAL,
			vote_count INTEGER,
			popularity REAL,
			runtime INTEGER,
			status TEXT,
			tagline TEXT,
			budget INTEGER,
			revenue INTEGER,
			imdb_rating TEXT,
			rotten_tomatoes_rating TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tv_shows (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tmdb_id INTEGER UNIQUE NOT NULL,
			name TEXT NOT NULL,
			overview TEXT,
			poster_path TEXT,
			backdrop_path TEXT,
			first_air_date TEXT,
			last_air_date TEXT,
			vote_average REAL,
			vote_count INTEGER,
			popularity REAL,
			number_of_seasons INTEGER,
			number_of_episodes INTEGER,
			status TEXT,
			type TEXT,
			imdb_rating TEXT,
			rotten_tomatoes_rating TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS watchlist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			content_id INTEGER NOT NULL,
			content_type TEXT NOT NULL CHECK(content_type IN ('movie', 'tv')),
			is_watched BOOLEAN DEFAULT FALSE,
			added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			watched_at DATETIME,
			UNIQUE(user_id, content_id, content_type)
		)`,
		`CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS movie_genres (
			movie_id INTEGER,
			genre_id INTEGER,
			PRIMARY KEY (movie_id, genre_id),
			FOREIGN KEY (movie_id) REFERENCES movies(id),
			FOREIGN KEY (genre_id) REFERENCES genres(id)
		)`,
		`CREATE TABLE IF NOT EXISTS tv_genres (
			tv_id INTEGER,
			genre_id INTEGER,
			PRIMARY KEY (tv_id, genre_id),
			FOREIGN KEY (tv_id) REFERENCES tv_shows(id),
			FOREIGN KEY (genre_id) REFERENCES genres(id)
		)`,
	}

	for _, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// InsertMovie inserts a movie into the database
func (d *Database) InsertMovie(movie interface{}) error {
	// This will be implemented when we have the actual movie data structure
	// For now, it's a placeholder
	return nil
}

// InsertTVShow inserts a TV show into the database
func (d *Database) InsertTVShow(tvShow interface{}) error {
	// This will be implemented when we have the actual TV show data structure
	// For now, it's a placeholder
	return nil
}

// GetWatchlist retrieves watchlist items for a user
func (d *Database) GetWatchlist(userID string) ([]interface{}, error) {
	query := `
		SELECT content_id, content_type, is_watched, added_at, watched_at
		FROM watchlist
		WHERE user_id = ?
		ORDER BY added_at DESC
	`

	rows, err := d.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query watchlist: %w", err)
	}
	defer rows.Close()

	var items []interface{}
	for rows.Next() {
		var contentID int
		var contentType string
		var isWatched bool
		var addedAt string
		var watchedAt *string
		if err := rows.Scan(&contentID, &contentType, &isWatched, &addedAt, &watchedAt); err != nil {
			continue
		}
		item := map[string]interface{}{
			"content_id":   contentID,
			"content_type": contentType,
			"is_watched":   isWatched,
			"added_at":     addedAt,
			"watched_at":   watchedAt,
		}
		items = append(items, item)
	}
	return items, nil
}

// AddToWatchlist adds an item to the user's watchlist
func (d *Database) AddToWatchlist(userID string, contentID int, contentType string) error {
	query := `
		INSERT OR REPLACE INTO watchlist (user_id, content_id, content_type, is_watched, added_at)
		VALUES (?, ?, ?, FALSE, CURRENT_TIMESTAMP)
	`

	_, err := d.DB.Exec(query, userID, contentID, contentType)
	if err != nil {
		return fmt.Errorf("failed to add to watchlist: %w", err)
	}

	return nil
}

// RemoveFromWatchlist removes an item from the user's watchlist
func (d *Database) RemoveFromWatchlist(userID string, contentID int, contentType string) error {
	query := `
		DELETE FROM watchlist
		WHERE user_id = ? AND content_id = ? AND content_type = ?
	`

	_, err := d.DB.Exec(query, userID, contentID, contentType)
	if err != nil {
		return fmt.Errorf("failed to remove from watchlist: %w", err)
	}

	return nil
}

// MarkAsWatched marks an item as watched in the user's watchlist
func (d *Database) MarkAsWatched(userID string, contentID int, contentType string, watched bool) error {
	var query string
	if watched {
		query = `
			UPDATE watchlist
			SET is_watched = TRUE, watched_at = CURRENT_TIMESTAMP
			WHERE user_id = ? AND content_id = ? AND content_type = ?
		`
	} else {
		query = `
			UPDATE watchlist
			SET is_watched = FALSE, watched_at = NULL
			WHERE user_id = ? AND content_id = ? AND content_type = ?
		`
	}

	_, err := d.DB.Exec(query, userID, contentID, contentType)
	if err != nil {
		return fmt.Errorf("failed to mark as watched: %w", err)
	}

	return nil
}
