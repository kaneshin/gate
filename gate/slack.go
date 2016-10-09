package gate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
func (s SlackIncomingService) NewPayload(channel, text string, atts ...Attachment) Payload {
	p := Payload{
		Channel: channel,
		Text:    text,
	}

	if len(atts) > 0 {
		p.Attachments = atts
	}
	return p
}

// Post posts data to slack.
func (s SlackIncomingService) Post(v interface{}) (*http.Response, error) {

	var body io.Reader
	switch v := v.(type) {
	case io.Reader:
		body = v

	case Payload, *Payload:
		buf := bytes.NewBufferString("payload=")
		if err := json.NewEncoder(buf).Encode(v); err != nil {
			return nil, err
		}
		body = buf

	}

	return s.config.HTTPClient.Post(s.baseURL.String(), bodyTypeURLEncoded, body)
}
