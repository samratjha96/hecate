package reddit

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL         = "https://www.reddit.com/r/%s/top.json?t=%s"
	defaultTimeout  = 30 * time.Second
	userAgentString = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"
)

// Client represents a Reddit API client
type Client struct {
	httpClient *http.Client
	userAgent  string
}

// Subreddit represents a Reddit subreddit
type Subreddit struct {
	Name                string
	NumberOfSubscribers int
	Posts               RedditPosts
}

// RedditPost represents a single Reddit post
type RedditPost struct {
	PostId        string
	Title         string
	Content       string
	DiscussionUrl string
	CommentCount  int
	Upvotes       int
	TimePosted    time.Time
}

// RedditPosts is a slice of RedditPost
type RedditPosts []RedditPost

// subredditResponseJson represents the JSON structure of a Reddit API response
type subredditResponseJson struct {
	Data struct {
		Children []struct {
			Data struct {
				Id                   string  `json:"id"`
				Title                string  `json:"title"`
				SelfText             string  `json:"selftext"`
				Upvotes              int     `json:"ups"`
				Url                  string  `json:"url"`
				Time                 float64 `json:"created"`
				CommentsCount        int     `json:"num_comments"`
				Domain               string  `json:"domain"`
				Permalink            string  `json:"permalink"`
				Stickied             bool    `json:"stickied"`
				Pinned               bool    `json:"pinned"`
				IsSelf               bool    `json:"is_self"`
				Thumbnail            string  `json:"thumbnail"`
				Flair                string  `json:"link_flair_text"`
				SubredditSubscribers int     `json:"subreddit_subscribers"`
				ParentList           []struct {
					Id        string `json:"id"`
					Subreddit string `json:"subreddit"`
					Permalink string `json:"permalink"`
				} `json:"crosspost_parent_list"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// NewClient creates a new Reddit API client
func NewClient(userAgent string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  userAgent,
	}
}

// makeRequest creates a new HTTP request for the Reddit API
func makeRequest(subreddit, sort string) (*http.Request, error) {
	requestUrl := fmt.Sprintf(baseURL, subreddit, strings.ToLower(sort))
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("User-Agent", userAgentString)
	return request, nil
}

// DescribeSubreddit fetches and describes a subreddit
func (c *Client) DescribeSubreddit(ctx context.Context, subreddit, sort string) (Subreddit, error) {
	if subreddit == "" || sort == "" {
		return Subreddit{}, fmt.Errorf("subreddit and sort cannot be empty")
	}

	request, err := makeRequest(subreddit, sort)
	if err != nil {
		return Subreddit{}, err
	}

	responseJson, err := decodeJSONFromRequest[subredditResponseJson](ctx, c.httpClient, request)
	if err != nil {
		return Subreddit{}, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(responseJson.Data.Children) == 0 {
		return Subreddit{}, fmt.Errorf("no posts found for subreddit: %s", subreddit)
	}

	posts := make([]RedditPost, 0, len(responseJson.Data.Children))
	for _, child := range responseJson.Data.Children {
		post := child.Data
		forumPost := RedditPost{
			PostId:        post.Id,
			Title:         html.UnescapeString(post.Title),
			Content:       html.UnescapeString(post.SelfText),
			DiscussionUrl: fmt.Sprintf("https://reddit.com%s", post.Permalink),
			CommentCount:  post.CommentsCount,
			Upvotes:       post.Upvotes,
			TimePosted:    time.Unix(int64(post.Time), 0),
		}
		posts = append(posts, forumPost)
	}

	return Subreddit{
		Name:                subreddit,
		NumberOfSubscribers: responseJson.Data.Children[0].Data.SubredditSubscribers,
		Posts:               posts,
	}, nil
}
