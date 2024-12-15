package hecate

import (
	"context"
	"fmt"
	"log"

	"github.com/samratjha96/hecate/internal/database"
	"github.com/samratjha96/hecate/internal/reddit"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"
)

// IngestAllSubreddit ingests posts from all subreddits in the database
func IngestAllSubreddit(ctx context.Context, db *database.DB, sortBy string) error {
	subreddits, err := db.GetAllSubreddits()
	if err != nil {
		return fmt.Errorf("failed to fetch subreddits: %w", err)
	}

	log.Printf("Starting ingestion for %d subreddits", len(subreddits))
	for _, subreddit := range subreddits {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Printf("Ingesting subreddit: %s", subreddit.Name)
			if _, err := IngestSubreddit(ctx, db, RedditSubscription{Name: subreddit.Name, SortBy: sortBy}); err != nil {
				log.Printf("Error ingesting subreddit %s: %v", subreddit.Name, err)
				// Continue with the next subreddit instead of returning the error
				continue
			}
		}
	}
	log.Printf("Completed ingestion for all subreddits")
	return nil
}

// IngestSubreddit ingests posts from a single subreddit
func IngestSubreddit(ctx context.Context, db *database.DB, subreddit RedditSubscription) (reddit.Subreddit, error) {
	log.Printf("Fetching data for subreddit: %s (Sort: %s)", subreddit.Name, subreddit.SortBy)

	client := reddit.NewClient(userAgent)

	response, err := client.DescribeSubreddit(ctx, subreddit.Name, subreddit.SortBy)
	if err != nil {
		return response, fmt.Errorf("failed to fetch subreddit posts: %w", err)
	}

	log.Printf("Successfully fetched %d posts for subreddit: %s", len(response.Posts), subreddit.Name)

	if err := upsertSubredditAndPosts(ctx, db, response, subreddit.Name, subreddit.SortBy); err != nil {
		return response, fmt.Errorf("failed to upsert subreddit and posts: %w", err)
	}

	return response, nil
}

// upsertSubredditAndPosts handles database operations for subreddit and its posts
func upsertSubredditAndPosts(ctx context.Context, db *database.DB, response reddit.Subreddit, subredditName, sortBy string) error {
	if _, err := db.UpsertSubreddit(response.Name, response.NumberOfSubscribers); err != nil {
		return fmt.Errorf("failed to upsert subreddit: %w", err)
	}

	log.Printf("Upserting %d %s posts for r/%s", len(response.Posts), sortBy, subredditName)

	for _, post := range response.Posts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := db.UpsertPost(post, subredditName); err != nil {
				log.Printf("Error upserting post %s: %v", post.PostId, err)
				// Continue with the next post instead of returning the error
				continue
			}
		}
	}
	log.Printf("Successfully upserted all posts for r/%s", subredditName)
	return nil
}

// GetAllSubreddits retrieves all subreddits from the database
func GetAllSubreddits(db *database.DB) ([]SubredditFrontendResponse, error) {
	fetchedSubredditDaos, err := db.GetAllSubreddits()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subreddits: %w", err)
	}

	responses := convertToSubredditResponses(fetchedSubredditDaos)
	log.Printf("Retrieved %d subreddits", len(responses))
	return responses, nil
}

// convertToSubredditResponses converts database objects to frontend response objects
func convertToSubredditResponses(daos []database.SubredditDao) []SubredditFrontendResponse {
	responses := make([]SubredditFrontendResponse, len(daos))
	for i, dao := range daos {
		responses[i] = SubredditFrontendResponse{
			Name:                dao.Name,
			NumberOfSubscribers: dao.NumberOfSubscribers,
		}
	}
	return responses
}

// GetAllPostsForSubreddit retrieves all posts for a specific subreddit
func GetAllPostsForSubreddit(db *database.DB, subredditName string) ([]SubredditPostFrontendResponse, error) {
	fetchedPosts, err := db.GetSubredditPosts(subredditName)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts for subreddit %s: %w", subredditName, err)
	}

	responses := convertToPostResponses(fetchedPosts)
	log.Printf("Retrieved %d posts for subreddit: %s", len(responses), subredditName)
	return responses, nil
}

// convertToPostResponses converts database objects to frontend response objects
func convertToPostResponses(daos []database.SubredditPostDao) []SubredditPostFrontendResponse {
	responses := make([]SubredditPostFrontendResponse, len(daos))
	for i, dao := range daos {
		responses[i] = SubredditPostFrontendResponse{
			Title:         dao.Title,
			Content:       dao.Content,
			DiscussionURL: dao.DiscussionURL,
			CommentCount:  dao.CommentCount,
			Upvotes:       dao.Upvotes,
		}
	}
	return responses
}
