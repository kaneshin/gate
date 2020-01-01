package gate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlack(t *testing.T) {

	t.Run("SlackIncomingService", func(t *testing.T) {
		t.Parallel()

		svc := NewSlackIncomingService(NewConfig())
		assert.NotNil(t, svc)
		assert.Nil(t, svc.baseURL)
	})
}
