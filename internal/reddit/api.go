package reddit

import (
	"fmt"
	"html"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
}

type Subreddit struct {
	Name                string
	NumberOfSubscribers int
}

type RedditPost struct {
	PostId        string
	Title         string
	Content       string
	DiscussionUrl string
	CommentCount  int
	Upvotes       int
	TimePosted    time.Time
}

type RedditPosts []RedditPost

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

func NewClient(userAgent string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Second * 30},
		userAgent:  userAgent,
	}
}

func (c *Client) FetchSubreddit(subreddit string) (Subreddit, error) {
	requestUrl := fmt.Sprintf("https://www.reddit.com/r/%s/hot.json", subreddit)

	request, err := http.NewRequest("GET", requestUrl, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")

	if err != nil {
		return Subreddit{}, err
	}

	responseJson, err := decodeJsonFromRequest[subredditResponseJson](c.httpClient, request)

	return Subreddit{
		Name:                subreddit,
		NumberOfSubscribers: responseJson.Data.Children[0].Data.SubredditSubscribers,
	}, err
}

func (c *Client) GetTopPosts(subreddit string) (RedditPosts, error) {
	requestUrl := fmt.Sprintf("https://www.reddit.com/r/%s/top.json", subreddit)
	request, err := http.NewRequest("GET", requestUrl, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0")

	if err != nil {
		return make(RedditPosts, 0), err
	}
	responseJson, err := decodeJsonFromRequest[subredditResponseJson](c.httpClient, request)

	if len(responseJson.Data.Children) == 0 {
		return nil, fmt.Errorf("no posts found")
	}

	posts := make([]RedditPost, 0, len(responseJson.Data.Children))

	for i := range responseJson.Data.Children {
		post := &responseJson.Data.Children[i].Data

		forumPost := RedditPost{
			PostId:        post.Id,
			Title:         html.UnescapeString(post.Title),
			Content:       html.UnescapeString(post.SelfText),
			DiscussionUrl: post.Url,
			CommentCount:  post.CommentsCount,
			Upvotes:       post.Upvotes,
			TimePosted:    time.Unix(int64(post.Time), 0),
		}

		posts = append(posts, forumPost)
	}

	return posts, err
}
