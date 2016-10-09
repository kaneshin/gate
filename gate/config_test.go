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

	conf.WithHTTPClient(http.DefaultClient)
	assert.NotNil(conf.HTTPClient)
}
