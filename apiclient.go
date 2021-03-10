package docbase

import (
	"net/http"
	"os"
	"sync"
)

var Client interface {
	Get(url string) (resp *http.Response, err error)
	// Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

var lock sync.Mutex
var defaultClient *APIClient

func init() {
	defaultClient = &APIClient{
		token:  os.Getenv("DOCBASE_TOKEN"),
		Client: http.DefaultClient,
	}
}

type APIClient struct {
	token string
	*http.Client
}
