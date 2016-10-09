package gate

import (
	"net/http"
	"net/url"
)

const (
	bodyTypeJSON       = "application/json; charset=utf-8"
	bodyTypeURLEncoded = "application/x-www-form-urlencoded"
)

type (
	service struct {
		config  *Config
		baseURL *url.URL
	}
)

// newService returns a new service. If a nil Config is
// provided, DefaultConfig will be used.
func newService(config *Config) *service {
	if config == nil {
		c := defaultConfig
		config = &c
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	return &service{
		config: config,
	}
}

// withBaseURL sets a base url value returning a service pointer
// for chaining.
func (s *service) withBaseURL(baseURL string) *service {
	s.baseURL, _ = url.Parse(baseURL)
	return s
}
