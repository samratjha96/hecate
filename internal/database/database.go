package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "hecate.db"
	dirPerms   = 0755
)

type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB() (*DB, error) {
	dataDir := os.Getenv("DB_DIRECTORY")
	if dataDir == "" {
		return nil, fmt.Errorf("DB_DIRECTORY environment variable is not set")
	}

	if err := os.MkdirAll(dataDir, dirPerms); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, dbFileName)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database at %s", dbPath)
	return &DB{db}, nil
}

// CreateTables creates the necessary tables in the database
func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS subreddits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			num_subscribers INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id TEXT UNIQUE NOT NULL,
			subreddit_name TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			discussion_url TEXT,
			comment_count INTEGER,
			upvotes INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			parent_comment_id INTEGER,
			content TEXT NOT NULL,
			comment_id TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
		)`,
	}

	for i, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table (query %d): %w", i+1, err)
		}
	}

	log.Println("Successfully created all necessary tables")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("Database connection closed successfully")
	return nil
}
