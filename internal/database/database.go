package database

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	connStr := "host=localhost port=5432 user=admin password=password dbname=hecate sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS subreddits (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			subreddit_id INTEGER REFERENCES subreddits(id),
			title VARCHAR(300) NOT NULL,
			content TEXT,
			post_id VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			metadata JSONB
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER REFERENCES posts(id),
			parent_comment_id INTEGER REFERENCES comments(id),
			content TEXT NOT NULL,
			comment_id VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			metadata JSONB
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}


func (db *DB) InsertSubreddit(name, description string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO subreddits (name, description) VALUES ($1, $2) RETURNING id", name, description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert subreddit: %v", err)
	}
	return id, nil
}
