package database

import (
	"fmt"
	"log"
	"time"

	"github.com/samratjha96/hecate/internal/reddit"
)

const (
	DefaultLimit = 10
	DefaultPage  = 1
)

type Paginate struct {
	Limit int
	Page  int
}

type SubredditDao struct {
	Name                string
	NumberOfSubscribers int
}

type SubredditPostDao struct {
	Title         string
	Content       string
	DiscussionURL string
	CommentCount  int
	Upvotes       int
	SubredditName string
}

// GetAllSubreddits retrieves all subreddits from the database
func (db *DB) GetAllSubreddits() ([]SubredditDao, error) {
	fetcher := func(page, limit int) (PaginatedResult[SubredditDao], error) {
		subreddits, nextPage, err := db.getSubredditsWithPagination(Paginate{
			Page:  page,
			Limit: limit,
		})
		if err != nil {
			return PaginatedResult[SubredditDao]{}, fmt.Errorf("failed to fetch subreddits: %w", err)
		}
		return PaginatedResult[SubredditDao]{
			Items:    subreddits,
			NextPage: nextPage,
		}, nil
	}

	return FetchAll(fetcher, DefaultPage, DefaultLimit)
}

// getSubredditsWithPagination retrieves a paginated list of subreddits
func (db *DB) getSubredditsWithPagination(pagination Paginate) ([]SubredditDao, int, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	nextPage := pagination.Page

	query := `
        SELECT id, name, num_subscribers, created_at
        FROM subreddits
        ORDER BY id
        LIMIT $1
        OFFSET $2
    `

	rows, err := db.Query(query, pagination.Limit, offset)
	if err != nil {
		return nil, nextPage, fmt.Errorf("failed to query subreddits: %w", err)
	}
	defer rows.Close()

	var subreddits []SubredditDao
	for rows.Next() {
		var s SubredditDao
		var id int64
		var createdAt time.Time
		if err := rows.Scan(&id, &s.Name, &s.NumberOfSubscribers, &createdAt); err != nil {
			return nil, nextPage, fmt.Errorf("failed to scan subreddit row: %w", err)
		}
		subreddits = append(subreddits, s)
	}

	if err := rows.Err(); err != nil {
		return nil, nextPage, fmt.Errorf("error iterating subreddit rows: %w", err)
	}

	if len(subreddits) > 0 {
		nextPage = pagination.Page + 1
	}

	return subreddits, nextPage, nil
}

// UpsertSubreddit inserts or updates a subreddit in the database
func (db *DB) UpsertSubreddit(name string, numberOfSubscribers int) (int, error) {
	var id int
	query := `
        INSERT INTO subreddits (name, num_subscribers)
        VALUES ($1, $2)
        ON CONFLICT (name) 
        DO UPDATE SET num_subscribers = EXCLUDED.num_subscribers
        RETURNING id
    `
	err := db.QueryRow(query, name, numberOfSubscribers).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert subreddit: %w", err)
	}
	log.Printf("Upserted subreddit: %s with %d subscribers", name, numberOfSubscribers)
	return id, nil
}

// UpsertPost inserts or updates a post in the database
func (db *DB) UpsertPost(post reddit.RedditPost, subredditName string) error {
	query := `
        INSERT INTO posts (subreddit_name, post_id, title, content, discussion_url, comment_count, upvotes, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (post_id) DO UPDATE SET
            title = EXCLUDED.title,
            content = EXCLUDED.content,
            discussion_url = EXCLUDED.discussion_url,
            comment_count = EXCLUDED.comment_count,
            upvotes = EXCLUDED.upvotes,
            created_at = EXCLUDED.created_at,
            updated_at = CURRENT_TIMESTAMP
    `

	_, err := db.Exec(query, subredditName, post.PostId, post.Title, post.Content, post.DiscussionUrl, post.CommentCount, post.Upvotes, post.TimePosted)
	if err != nil {
		return fmt.Errorf("failed to upsert post: %w", err)
	}
	log.Printf("Upserted post: %s for subreddit: %s", post.Title, subredditName)
	return nil
}

// GetSubredditPosts retrieves all posts for a given subreddit
func (db *DB) GetSubredditPosts(subredditName string) ([]SubredditPostDao, error) {
	fetcher := func(page, limit int) (PaginatedResult[SubredditPostDao], error) {
		posts, nextPage, err := db.getSubredditPostsWithPagination(subredditName, Paginate{
			Page:  page,
			Limit: limit,
		})
		if err != nil {
			return PaginatedResult[SubredditPostDao]{}, fmt.Errorf("failed to fetch posts for subreddit %s: %w", subredditName, err)
		}
		return PaginatedResult[SubredditPostDao]{
			Items:    posts,
			NextPage: nextPage,
		}, nil
	}

	return FetchAll(fetcher, DefaultPage, DefaultLimit)
}

// SearchPosts searches for posts across all subreddits
func (db *DB) SearchPosts(query string) ([]SubredditPostDao, error) {
	sqlQuery := `
		SELECT p.title, p.content, p.discussion_url, p.comment_count, p.upvotes, p.subreddit_name
		FROM posts p
		WHERE p.title ILIKE $1 OR p.content ILIKE $1
		ORDER BY p.created_at DESC
		LIMIT 100
	`
	searchPattern := "%" + query + "%"

	rows, err := db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}
	defer rows.Close()

	var posts []SubredditPostDao
	for rows.Next() {
		var p SubredditPostDao
		if err := rows.Scan(&p.Title, &p.Content, &p.DiscussionURL, &p.CommentCount, &p.Upvotes, &p.SubredditName); err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}

	return posts, nil
}

// getSubredditPostsWithPagination retrieves a paginated list of posts for a given subreddit
func (db *DB) getSubredditPostsWithPagination(subredditName string, pagination Paginate) ([]SubredditPostDao, int, error) {
	offset := (pagination.Page - 1) * pagination.Limit
	nextPage := pagination.Page

	query := `
        SELECT title, content, discussion_url, comment_count, upvotes
        FROM posts
        WHERE subreddit_name = $1
        ORDER BY created_at DESC
        LIMIT $2
        OFFSET $3
    `

	rows, err := db.Query(query, subredditName, pagination.Limit, offset)
	if err != nil {
		return nil, nextPage, fmt.Errorf("failed to query posts for subreddit %s: %w", subredditName, err)
	}
	defer rows.Close()

	var posts []SubredditPostDao
	for rows.Next() {
		var p SubredditPostDao
		if err := rows.Scan(&p.Title, &p.Content, &p.DiscussionURL, &p.CommentCount, &p.Upvotes); err != nil {
			return nil, nextPage, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, nextPage, fmt.Errorf("error iterating post rows: %w", err)
	}

	if len(posts) > 0 {
		nextPage = pagination.Page + 1
	}

	return posts, nextPage, nil
}
