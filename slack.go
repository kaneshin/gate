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

// Post posts a payload to slack.
func (s SlackIncomingService) Post(contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return s.config.HTTPClient.Do(req)
}

// PostTextPayload posts a text payload to slack.
func (s SlackIncomingService) PostTextPayload(payload TextPayload) (*http.Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return s.Post(bodyTypeJSON, buf)
}
