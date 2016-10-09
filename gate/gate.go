package gate

import "net/http"

const (
	bodyTypeJSON       = "application/json; charset=utf-8"
	bodyTypeURLEncoded = "application/x-www-form-urlencoded"
)

type (
	service struct {
		config *Config
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
