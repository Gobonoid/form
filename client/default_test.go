package client

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	validTestBaseURL = "test.com"
)

func TestNewClient(t *testing.T) {
	var optionalClient = &http.Client{}

	tests := []struct {
		name             string
		expectErrMessage string
		baseURL          string
		options          []Option
	}{
		{
			name:             "empty baseURL",
			expectErrMessage: "empty baseURL",
			baseURL:          "",
		},
		{
			name:             "incorrect url",
			expectErrMessage: "failed to parse baseURL.*",
			baseURL:          "www.432234324234-'12#31#23#1#65$%^&*)_+/co.uk########################",
		},
		{
			name:    "default options",
			baseURL: validTestBaseURL,
		},
		{
			name:    "custom client",
			baseURL: validTestBaseURL,
			options: []Option{WithHTTPClient(optionalClient), WithHTTPS()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewDefaultClient(tt.baseURL, tt.options...)
			if tt.expectErrMessage != "" {
				assert.Error(t, err)
				assert.Nil(t, c)
				assert.Regexp(t, tt.expectErrMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, c)
				if len(tt.options) == 0 {
					assert.Equal(t, c.conf.c, http.DefaultClient)
					assert.Equal(t, c.conf.scheme, "http")
				} else {
					assert.Equal(t, c.conf.c, optionalClient)
					assert.Equal(t, c.conf.scheme, "https")
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	fakeErr := errors.New("fake error")

	c, err := NewDefaultClient(validTestBaseURL)
	require.NoError(t, err)
	httpmock.Activate()
	defer httpmock.Deactivate()

	tests := []struct {
		name             string
		expectErrMessage string
		ctx              context.Context
		path             string
		setup            func(path, respBody string)
		respBody         string
	}{
		{
			name:             "failed to create request",
			expectErrMessage: "failed to create new GET requestWithContext",
			ctx:              nil,
		},
		{
			name:             "no response",
			ctx:              ctx,
			path:             "/IWontRespondHere",
			expectErrMessage: "request failed.*",
			setup: func(path, respBody string) {
				httpmock.RegisterNoResponder(httpmock.NewErrorResponder(fakeErr))
			},
		},
		{
			name: "success",
			ctx:  ctx,
			path: "/something",
			setup: func(path string, respBody string) {
				httpmock.RegisterResponder(http.MethodGet, c.conf.scheme+"://"+validTestBaseURL+path,
					httpmock.NewStringResponder(http.StatusOK, respBody))
			},
			respBody: `{"test": "blablab"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.path, tt.respBody)
			}
			resp, err := c.Get(tt.ctx, tt.path)

			if tt.expectErrMessage != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Regexp(t, tt.expectErrMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestPost(t *testing.T) {
	ctx := context.Background()
	fakeErr := errors.New("fake error")

	c, err := NewDefaultClient(validTestBaseURL)
	require.NoError(t, err)
	httpmock.Activate()
	defer httpmock.Deactivate()

	tests := []struct {
		name             string
		expectErrMessage string
		ctx              context.Context
		path             string
		setup            func(path, respBody string)
		reqBody          string
	}{
		{
			name:             "failed to create request",
			expectErrMessage: "failed to create new POST RequestWithContext",
			ctx:              nil,
		},
		{
			name:             "no response",
			ctx:              ctx,
			path:             "/IWontRespondHere",
			expectErrMessage: "request failed.*",
			setup: func(path, respBody string) {
				httpmock.RegisterNoResponder(httpmock.NewErrorResponder(fakeErr))
			},
		},
		{
			name: "success",
			ctx:  ctx,
			path: "/something",
			setup: func(path string, reqBody string) {
				httpmock.RegisterResponder(http.MethodPost, c.conf.scheme+"://"+validTestBaseURL+path,
					func(req *http.Request) (*http.Response, error) {
						p, err := io.ReadAll(req.Body)
						req.Body.Close()
						assert.NoError(t, err)
						assert.Equal(t, string(p), reqBody)
						return &http.Response{
							StatusCode: http.StatusCreated,
						}, nil
					})
			},
			reqBody: "fakeString",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.path, tt.reqBody)
			}
			resp, err := c.Post(tt.ctx, tt.path, strings.NewReader(tt.reqBody))
			if tt.expectErrMessage != "" {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Regexp(t, tt.expectErrMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
