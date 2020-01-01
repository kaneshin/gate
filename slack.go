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

// Post posts data to slack.
func (s SlackIncomingService) Post(contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return s.config.HTTPClient.Do(req)
}

type TextData struct {
	Text string `json:"text"`
}

// PostTextData posts text data to slack.
func (s SlackIncomingService) PostTextData(data TextData) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return s.Post(bodyTypeJSON, buf)
}
