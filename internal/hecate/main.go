package hecate

import (
	"fmt"
	"log"

	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/reddit"
)

type RedditSubscription struct {
	Name   string `json:"name"`
	SortBy string `json:"sort_by"`
}

func Subscribe(db *database.DB, subreddits []RedditSubscription) error {
	fmt.Println("Hello from Hecate")

	client := reddit.NewClient("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")

	// loop through subreddits
	for _, subreddit := range subreddits {

		response, err := client.DescribeSubreddit(subreddit.Name, subreddit.SortBy)

		// Upsert subreddit into database
		_, err = db.UpsertSubreddit(response.Name, response.NumberOfSubscribers)
		if err != nil {
			log.Printf("Failed to insert subreddit: %v", err)
			return err
		}

		if err != nil {
			log.Printf("Failed to fetch subreddit posts: %v", err)
			return err
		}

		fmt.Printf("Inserting %d %s posts from r/%s\n", len(response.Posts), subreddit.SortBy, subreddit.Name)

		for _, post := range response.Posts {
			err := db.UpsertPost(post, subreddit.Name)
			if err != nil {
				log.Printf("Failed to insert post: %v", err)
				return err
			}
		}
	}
	return nil
}
