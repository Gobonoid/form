package client

import "net/http"

//Option definition for DefaultClient
type Option func(p *Config)

//WithHTTPClient allows injecting http.Client, if not used http.DefaultClient is being used
func WithHTTPClient(c *http.Client) Option {
	return func(p *Config) { p.c = c }
}

//WithHTTPS sets scheme to be used as HTTPS
func WithHTTPS() Option {
	return func(p *Config) { p.scheme = "https" }
}
