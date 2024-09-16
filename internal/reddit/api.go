package reddit

import (
	"fmt"
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

type subredditResponseJson struct {
	Data struct {
		Children []struct {
			Data struct {
				Id                   string  `json:"id"`
				Title                string  `json:"title"`
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
