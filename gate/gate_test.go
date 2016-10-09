package gate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {

	t.Run("Service with config", func(t *testing.T) {
		t.Parallel()

		svc := newService(NewConfig())
		assert.NotNil(t, svc)
		assert.NotNil(t, svc.config)
		assert.NotNil(t, svc.config.HTTPClient)
	})

	t.Run("Service with nil config", func(t *testing.T) {
		t.Parallel()

		svc := newService(nil)
		assert.NotNil(t, svc)
		assert.NotNil(t, svc.config)
		assert.NotNil(t, svc.config.HTTPClient)
	})
}
