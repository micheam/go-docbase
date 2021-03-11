package docbase

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex
var defaultClient *Client

func init() {
	defaultClient = &Client{
		token:  os.Getenv("DOCBASE_TOKEN"),
		Client: http.DefaultClient,
	}
}

type Client struct {
	token string
	*http.Client
}

const baseURL string = "https://api.docbase.io"

func buildURL(paths ...string) string {
	return strings.Join(append([]string{baseURL}, paths...), "/")
}

func (c *Client) NewRequest(ctx context.Context, method, _url string, body io.Reader, param *url.Values) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, _url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Api-Version", "2")
	req.Header.Add("X-DocBaseToken", c.token)
	req.Header.Add("Content-Type", "application/json")
	if param != nil {
		req.URL.RawQuery = param.Encode()
	}
	return req, nil
}
