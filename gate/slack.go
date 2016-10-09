package gate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kaneshin/gate/gate/slack"
)

type (
	// SlackIncomingService is a slack incoming webhook service.
	SlackIncomingService struct {
		*service
	}
)

// NewSlackIncomingService returns a new SlackIncomingService.
func NewSlackIncomingService(config *Config) *SlackIncomingService {
	return &SlackIncomingService{
		service: newService(config),
	}
}

// WithBaseURL sets a base url value returning a service pointer
// for chaining.
func (s *SlackIncomingService) WithBaseURL(baseURL string) *SlackIncomingService {
	s.service.withBaseURL(baseURL)
	return s
}

// NewPayload returns a new Payload.
func (s SlackIncomingService) NewPayload(channel, text string, atts ...slack.Attachment) slack.Payload {
	p := slack.Payload{
		Channel: channel,
		Text:    text,
	}

	p.Attachments = atts
	return p
}

// Post posts data to slack.
func (s SlackIncomingService) Post(v interface{}) (*http.Response, error) {
	var body io.Reader
	switch v := v.(type) {
	case io.Reader:
		body = v

	case slack.Payload, *slack.Payload:
		buf := bytes.NewBufferString("payload=")
		if err := json.NewEncoder(buf).Encode(v); err != nil {
			return nil, err
		}
		body = buf

	}

	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyTypeURLEncoded)

	return s.config.HTTPClient.Do(req)
}
