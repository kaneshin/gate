package gate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig()
	assert.NotNil(conf)
	assert.Nil(conf.HTTPClient)
	assert.Empty(conf.AccessToken)

	conf.WithHTTPClient(http.DefaultClient).
		WithAccessToken("access-token")

	assert.Equal(http.DefaultClient, conf.HTTPClient)
	assert.Equal("access-token", conf.AccessToken)
}
