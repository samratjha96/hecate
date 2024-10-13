package hecate

type SubredditFrontendResponse struct {
	Name                string `json:"name"`
	NumberOfSubscribers int    `json:"numberOfSubscribers"`
}

type SubredditPostFrontendResponse struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	DiscussionURL string `json:"discussionUrl"`
	CommentCount  int    `json:"commentCount"`
	Upvotes       int    `json:"upvotes"`
}

type RedditSubscription struct {
	Name   string `json:"name"`
	SortBy string `json:"sortBy"`
}

type SubscribeFrontendRequest struct {
	Subscription RedditSubscription `json:"subreddit"`
}

type IngestAllFrontendRequest struct {
	SortBy string `json:"sortBy"`
}
