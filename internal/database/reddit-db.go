package database

import (
	"fmt"

	"github.com/samratjha96/hecate/internal/reddit"
)

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
