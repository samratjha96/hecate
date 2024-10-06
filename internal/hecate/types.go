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

type SubscribeFrontendRequest struct {
	Subscription RedditSubscription `json:"subreddit"`
}
