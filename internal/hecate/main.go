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

func Subscribe(db *database.DB, subreddits []RedditSubscription) ([]reddit.Subreddit, error) {
	fmt.Printf("Subscribing to %s", subreddits)

	client := reddit.NewClient("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")

	responses := make([]reddit.Subreddit, 0)

	// loop through subreddits
	for _, subreddit := range subreddits {

		response, err := client.DescribeSubreddit(subreddit.Name, subreddit.SortBy)

		// Upsert subreddit into database
		_, err = db.UpsertSubreddit(response.Name, response.NumberOfSubscribers)
		if err != nil {
			log.Printf("Failed to insert subreddit: %v", err)
			return responses, err
		}

		if err != nil {
			log.Printf("Failed to fetch subreddit posts: %v", err)
			return responses, err
		}

		fmt.Printf("Inserting %d %s posts from r/%s\n", len(response.Posts), subreddit.SortBy, subreddit.Name)

		for _, post := range response.Posts {
			err := db.UpsertPost(post, subreddit.Name)
			if err != nil {
				log.Printf("Failed to insert post: %v", err)
				return responses, err
			}
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func GetAllSubreddits(db *database.DB) ([]database.SubredditDao, error) {
	var subredditDaos []database.SubredditDao
	nextPage := 1

	for true {
		fetchedSubredditDaos, newNextPage, err := db.GetSubreddits(database.Paginate{
			Page:  nextPage,
			Limit: 10,
		})
		if err != nil {
			log.Printf("Failed to fetch subreddits: %v", err)
			return nil, err
		}

		subredditDaos = append(subredditDaos, fetchedSubredditDaos...)
		if newNextPage == nextPage {
			break
		}

		nextPage = newNextPage

	}
	return subredditDaos, nil
}
