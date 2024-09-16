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

	subreddits := []string{"solotravel", "travelhacks", "wanderlust"}

	// loop through subreddits
	for _, subreddit := range subreddits {
		client := reddit.NewClient("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")
		subreddit, err := client.FetchSubreddit(subreddit)
		if err != nil {
			log.Fatalf("Failed to fetch subreddit: %v", err)
		}

		// Upsert subreddit into database
		_, err = db.UpsertSubreddit(subreddit.Name, subreddit.NumberOfSubscribers)
		if err != nil {
			log.Fatalf("Failed to insert subreddit: %v", err)
		}
	}
	return 0
}
