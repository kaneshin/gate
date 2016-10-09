package gate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kaneshin/gate/gate/facebook"
)

type (
	// FacebookService is a slack incoming webhook service.
	FacebookService struct {
		*service
	}
)

// NewFacebookService returns a new FacebookService.
func NewFacebookService(config *Config) *FacebookService {
	svc := &FacebookService{
		service: newService(config).withBaseURL(facebook.SendAPIURL),
	}

	q := svc.baseURL.Query()
	q.Set("access_token", config.AccessToken)
	svc.baseURL.RawQuery = q.Encode()

	return svc
}

// NewPayload returns a new Payload.
func (s FacebookService) NewPayload(id, text string) facebook.Payload {
	p := facebook.Payload{
		Recipient: facebook.Recipient{
			ID: id,
		},
		Message: facebook.Message{
			Text: text,
		},
	}

	return p
}

// Post posts data to Facebook.
func (s FacebookService) Post(v interface{}) (*http.Response, error) {
	var body io.Reader
	switch v := v.(type) {
	case io.Reader:
		body = v

	case facebook.Payload, *facebook.Payload:
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(v); err != nil {
			return nil, err
		}
		body = &buf

	}

	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyTypeJSON)

	return s.config.HTTPClient.Do(req)
}