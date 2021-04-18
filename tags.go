package docbase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Tag struct {
	Name string `json:"name"`
}

func ListTags(ctx context.Context, domain string) ([]Tag, error) {
	return defaultClient.ListTags(ctx, domain)
}

func (c *Client) ListTags(ctx context.Context, domain string) ([]Tag, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, buildURL("teams", domain, "tags"), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	tags := make([]Tag, 0)
	b, err := ioutil.ReadAll(resp.Body) // TOOD(micheam): change to io.ReadAll
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
