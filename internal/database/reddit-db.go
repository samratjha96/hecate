package database

import (
	"fmt"
	"time"

	"github.com/samratjha96/hecate/internal/reddit"
)

const (
	DefaultLimit = 10
)

type Paginate struct {
	Limit int
	Page  int
}

type SubredditDao struct {
	Name                string
	NumberOfSubscribers int
}

func (db *DB) GetAllSubreddits() ([]SubredditDao, error) {
	fetcher := func(page, limit int) (PaginatedResult[SubredditDao], error) {
		subreddits, nextPage, err := db.getSubredditsWithPagination(Paginate{
			Page:  page,
			Limit: limit,
		})
		return PaginatedResult[SubredditDao]{
			Items:    subreddits,
			NextPage: nextPage,
		}, err
	}

	return FetchAll(fetcher, 1, 10)
}

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
		return nil, nextPage, err
	}
	defer rows.Close()

	var subreddits []SubredditDao
	for rows.Next() {
		var s SubredditDao
		var id int64
		var createdAt time.Time
		err := rows.Scan(&id, &s.Name, &s.NumberOfSubscribers, &createdAt)
		if err != nil {
			return subreddits, nextPage, err
		}
		subreddits = append(subreddits, s)
	}

	if len(subreddits) > 0 {
		nextPage = pagination.Page + 1
	}

	return subreddits, nextPage, nil

}

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
		return 0, fmt.Errorf("failed to upsert subreddit: %v", err)
	}
	return id, nil
}

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
		return err
	}

	return nil
}
