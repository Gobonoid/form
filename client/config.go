package client

import "net/http"

//Config holds clients configuration values
type Config struct {
	c      *http.Client
	scheme string
}
