package docbase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type (
	PostID int
	Scope  string
)

func (p PostID) Int() int {
	return int(p)
}

const (
	ScopeEveryone Scope = "everyone"
	ScopeGroup    Scope = "group"
	ScopePrivate  Scope = "private"
)

func ParsePostID(s string) (PostID, error) {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return PostID(i), nil
}

func (p PostID) String() string {
	return fmt.Sprintf("%d", p)
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
	Scope      Scope         `json:"scope"`
	SharingURL string        `json:"sharing_url"`
	Tags       []Tag         `json:"tags"`
	User       User          `json:"user"`
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

func (c *Client) ListPosts(ctx context.Context, domain string, param url.Values) ([]Post, *Meta, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, buildURL("teams", domain, "posts"), nil, &param)
	if err != nil {
		return nil, nil, err
	}
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

func (c *Client) GetPost(ctx context.Context, domain string, id PostID) (*Post, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, buildURL("teams", domain, "posts", fmt.Sprint(id)), nil, nil)
	if err != nil {
		return nil, err
	}
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

type PostOption struct {
	Draft  *bool    `json:"draft,omitempty"`
	Notice *bool    `json:"notice,omitempty"`
	Tags   []string `json:"tags"`
	Scope  string   `json:"scope,omitempty"` // TODO(micheam): 指定可能な値を明示する everyon (default), group, private
	Groups []int    `json:"groups"`          // require on scope:groups
}

func NewPost(ctx context.Context, domain string, title string, body io.Reader, option PostOption) (*Post, error) {
	return defaultClient.NewPost(ctx, domain, title, body, option)
}

func (c *Client) NewPost(ctx context.Context, domain string, title string, body io.Reader, option PostOption) (*Post, error) {
	if domain == "" {
		return nil, errors.New("`domain` must not be empty")
	}
	if title == "" {
		return nil, errors.New("`title` must not be empty")
	}

	var _body = new(bytes.Buffer)
	{
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}

		type RequestBody struct {
			Title string `json:"title"`
			Body  string `json:"body"`

			PostOption
		}
		bb, err := json.Marshal(RequestBody{title, string(b), option})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		_body = bytes.NewBuffer(bb)
	}

	req, err := c.NewRequest(ctx, http.MethodPost, buildURL("teams", domain, "posts"), _body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http reqest: %w", err)
	}
	defer resp.Body.Close()
	if 300 <= resp.StatusCode {
		// TODO(micheam): handle error object.
		//   エラーの詳細情報がBodyで返却されるので、ちゃんと扱う
		return nil, fmt.Errorf("docbase api returns NG: %s", resp.Status)
	}
	var (
		created = new(Post)
		bytes   = []byte{}
	)
	if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read Response body: %w", err)
	}
	if err := json.Unmarshal(bytes, created); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return created, nil
}

// UpdateFields は、更新対象のフィールドを保持します。
type UpdateFields struct {
	Title  *string   `json:"title,omitempty"`
	Draft  *bool     `json:"draft,omitempty"`
	Notice *bool     `json:"notice,omitempty"`
	Tags   *[]string `json:"tags,omitempty"`
	Scope  *string   `json:"scope,omitempty"`
	Groups *[]int    `json:"groups,omitempty"`
}

func UpdatePost(ctx context.Context, domain string, id PostID, body io.Reader, fields UpdateFields) (*Post, error) {
	return defaultClient.UpdatePost(ctx, domain, id, body, fields)
}

func (c *Client) UpdatePost(ctx context.Context, domain string, id PostID, body io.Reader, fields UpdateFields) (*Post, error) {
	if domain == "" {
		return nil, errors.New("`domain` must not be empty")
	}

	var _body *bytes.Buffer
	{
		rb := struct {
			Body *string `json:"body,omitempty"`
			UpdateFields
		}{
			UpdateFields: fields,
		}
		if body != nil {
			b, err := ioutil.ReadAll(body)
			if err != nil {
				return nil, fmt.Errorf("failed to read body: %w", err)
			}
			rb.Body = stringPtr(string(b))
		}
		bb, err := json.Marshal(rb)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		_body = bytes.NewBuffer(bb)
	}

	req, err := c.NewRequest(ctx, http.MethodPatch, buildURL("teams", domain, "posts", id.String()), _body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to patch request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http reqest: %w", err)
	}
	defer resp.Body.Close()
	if 300 <= resp.StatusCode {
		// TODO(micheam): handle error object.
		//   エラーの詳細情報がBodyで返却されるので、ちゃんと扱う
		return nil, fmt.Errorf("docbase api returns NG: %s", resp.Status)
	}
	var (
		created = new(Post)
		bytes   = []byte{}
	)
	if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read Response body: %w", err)
	}
	if err := json.Unmarshal(bytes, created); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return created, nil
}
