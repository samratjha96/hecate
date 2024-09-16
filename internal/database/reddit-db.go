package database

import "fmt"

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
