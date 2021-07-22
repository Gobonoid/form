package form

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

//HTTPClient definition that is used in form API client
type HTTPClient interface {
	Get(ctx context.Context, path string) (*http.Response, error)
	Post(ctx context.Context, path string, body io.Reader) (*http.Response, error)
	DeleteWithQueryParams(ctx context.Context, path string, q url.Values) (*http.Response, error)
}
