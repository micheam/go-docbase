package docbase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type PostID int

func ParsePostID(s string) (PostID, error) {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return PostID(i), nil
}

type Post struct {
	ID         PostID        `json:"id"`
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	Draft      bool          `json:"draft"`
	Archived   bool          `json:"archived"`
	URL        string        `json:"url"`
	CreatedAt  string        `json:"created_at"` // ISO 8601
	UpdatedAt  string        `json:"updated_at"` // ISO 8601
	Scope      string        `json:"scope"`
	SharingURL string        `json:"sharing_url"`
	Tags       []interface{} `json:"tags"`
	User       interface{}   `json:"user"`
	Stars      int           `json:"stars_count"`
	GoodJob    int           `json:"good_jobs_count"`
	Comments   []interface{} `json:"comments"`
	Groups     []interface{} `json:"groups"`
}

type Meta struct {
	PreviousPageURL string `json:"previous_page"`
	NextPageURL     string `json:"next_page"`
	Total           int    `json:"total"`
}

func ListPosts(ctx context.Context, domain string, param url.Values) ([]Post, *Meta, error) {
	return defaultClient.ListPosts(ctx, domain, param)
}

func (c *APIClient) ListPosts(ctx context.Context, domain string, param url.Values) ([]Post, *Meta, error) {
	url := fmt.Sprintf("https://api.docbase.io/teams/%s/posts", domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("X-DocBaseToken", c.token)
	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = param.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	data := struct {
		Posts []Post `json:"posts"`
		Meta  *Meta  `json:"meta"`
	}{}
	b, err := ioutil.ReadAll(resp.Body) // TOOD(micheam): change to io.ReadAll
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, nil, err
	}
	return data.Posts, data.Meta, nil
}

func GetPost(ctx context.Context, domain string, id PostID) (*Post, error) {
	return defaultClient.GetPost(ctx, domain, id)
}

func (c *APIClient) GetPost(ctx context.Context, domain string, id PostID) (*Post, error) {
	_url := fmt.Sprintf("https://api.docbase.io/teams/%s/posts/%d", domain, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, _url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("X-DocBaseToken", c.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	post := new(Post)
	b, err := ioutil.ReadAll(resp.Body) // TOOD(micheam): change to io.ReadAll
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}