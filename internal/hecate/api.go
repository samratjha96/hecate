package hecate

import (
	"fmt"
	"log"

	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/reddit"
)

type RedditSubscription struct {
	Name   string `json:"name"`
	SortBy string `json:"sortBy"`
}

func IngestSubreddit(db *database.DB, subreddit RedditSubscription) (reddit.Subreddit, error) {
	fmt.Printf("Subscribing to %s", subreddit)

	client := reddit.NewClient("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")

	response, err := client.DescribeSubreddit(subreddit.Name, subreddit.SortBy)

	// Upsert subreddit into database
	_, err = db.UpsertSubreddit(response.Name, response.NumberOfSubscribers)
	if err != nil {
		log.Printf("Failed to insert subreddit: %v", err)
		return response, err
	}

	if err != nil {
		log.Printf("Failed to fetch subreddit posts: %v", err)
		return response, err
	}

	fmt.Printf("Inserting %d %s posts from r/%s\n", len(response.Posts), subreddit.SortBy, subreddit.Name)

	for _, post := range response.Posts {
		err := db.UpsertPost(post, subreddit.Name)
		if err != nil {
			log.Printf("Failed to insert post: %v", err)
			return response, err
		}
	}
	return response, nil
}

func GetAllSubreddits(db *database.DB) ([]SubredditFrontendResponse, error) {
	fetchedSubredditDaos, err := db.GetAllSubreddits()
	if err != nil {
		log.Printf("Failed to fetch subreddits: %v", err)
		return nil, err
	}

	var subredditResponses []SubredditFrontendResponse
	for _, dao := range fetchedSubredditDaos {
		subredditResponses = append(
			subredditResponses,
			SubredditFrontendResponse{
				Name:                dao.Name,
				NumberOfSubscribers: dao.NumberOfSubscribers,
			},
		)
	}
	return subredditResponses, nil
}

func GetAllPostsForSubreddit(db *database.DB, subredditName string) ([]SubredditPostFrontendResponse, error) {
	fetchedPosts, err := db.GetSubredditPosts(subredditName)
	if err != nil {
		log.Printf("Failed to get posts for subreddit %s with error: %v", subredditName, err)
		return nil, err
	}
	var postsResponses []SubredditPostFrontendResponse
	for _, dao := range fetchedPosts {
		postsResponses = append(
			postsResponses,
			SubredditPostFrontendResponse{
				Title:         dao.Title,
				Content:       dao.Content,
				DiscussionURL: dao.DiscussionURL,
				CommentCount:  dao.CommentCount,
				Upvotes:       dao.Upvotes,
			},
		)
	}
	return postsResponses, nil
}
