package hecate

import (
    "fmt"
    "log"
    "github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/reddit"
)

func Main() int {
    fmt.Println("Hello from Hecate")
    // Initialize database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	_, err = db.InsertSubreddit("solotravel", "A place for solo travelers to share experiences")
	if err != nil {
		log.Fatalf("Failed to insert subreddit: %v", err)
	}
    return 0
}
