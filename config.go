package gate

import "net/http"

type (
	// Config provides service configuration for service.
	Config struct {
		// The HTTP client to use when sending requests.
		HTTPClient *http.Client

		AccessToken string
	}
)

var (
	defaultConfig = *(NewConfig().WithHTTPClient(http.DefaultClient))
)

// NewConfig returns a pointer of new Config objects.
func NewConfig() *Config {
	return &Config{}
}

// WithHTTPClient sets a config HTTPClient value returning a Config pointer
// for chaining.
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// WithAccessToken sets a access token value to verify service returning
// a Config pointer for chaining.
func (c *Config) WithAccessToken(token string) *Config {
	c.AccessToken = token
	return c
}
