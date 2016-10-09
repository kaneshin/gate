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

		payload := svc.NewPayload("#general", "foobar")
		assert.Equal(t, "#general", payload.Channel)
		assert.Equal(t, "foobar", payload.Text)
	})
}
