package client

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

//DefaultClient that implements form HttpClient interface and behaves as DI container
type DefaultClient struct {
	baseURL string
	conf    *Config
}

//NewDefaultClient behaves as a constructor
func NewDefaultClient(baseURL string, opts ...Option) (*DefaultClient, error) {
	if baseURL == "" {
		return nil, errors.New("empty baseURL")
	}
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse baseURL: %s", baseURL)
	}
	conf := &Config{}
	for _, opt := range opts {
		opt(conf)
	}
	if conf.c == nil {
		conf.c = http.DefaultClient
	}
	if conf.scheme == "" {
		conf.scheme = "http"
	}

	return &DefaultClient{
		conf:    conf,
		baseURL: baseURL,
	}, nil
}

//Get request method implementation
func (client *DefaultClient) Get(ctx context.Context, path string) (*http.Response, error) {
	u := &url.URL{
		Scheme: client.conf.scheme,
		Host:   client.baseURL,
		Path:   path,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GET requestWithContext")
	}
	resp, err := client.conf.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	return resp, nil
}

//Post request method implementation
func (client *DefaultClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	u := &url.URL{
		Scheme: client.conf.scheme,
		Host:   client.baseURL,
		Path:   path,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new POST RequestWithContext")
	}
	resp, err := client.conf.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	return resp, nil
}

//DeleteWithQueryParams is a DELETE request method implementation with extra query params
func (client *DefaultClient) DeleteWithQueryParams(ctx context.Context, path string, q url.Values) (*http.Response, error) {
	u := &url.URL{
		Scheme:   client.conf.scheme,
		Host:     client.baseURL,
		Path:     path,
		RawQuery: q.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new DELETE RequestWithContext")
	}
	resp, err := client.conf.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	return resp, nil
}
