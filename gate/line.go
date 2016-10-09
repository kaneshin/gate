package gate

import (
	"bytes"
	"io"
	"net/http"
)

type (
	// LINEService is a slack incoming webhook service.
	LINEService struct {
		*service
	}
)

// NewLINEService returns a new LINEService.
func NewLINEService(config *Config) *LINEService {
	const notifyAPIURL = "https://notify-api.line.me/api/notify"
	svc := &LINEService{
		service: newService(config).withBaseURL(notifyAPIURL),
	}
	return svc
}

// Post posts data to LINE.
func (s LINEService) Post(v interface{}) (*http.Response, error) {

	var body io.Reader
	switch v := v.(type) {
	case io.Reader:
		body = v

	case string:
		buf := bytes.NewBufferString("message=" + v)
		body = buf

	}

	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyTypeURLEncoded)
	req.Header.Add("Authorization", "Bearer "+s.service.config.AccessToken)

	return s.config.HTTPClient.Do(req)
}
