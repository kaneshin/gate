package gate

import (
	"bytes"
	"io"
	"net/http"
)

// lineNotifyAPIURL is Notify API URL for LINE.
const lineNotifyAPIURL = "https://notify-api.line.me/api/notify"

type (
	// LINENotifyService is a slack incoming webhook service.
	LINENotifyService struct {
		*service
	}
)

// NewLINENotifyService returns a new LINENotifyService.
func NewLINENotifyService(config *Config) *LINENotifyService {
	svc := &LINENotifyService{
		service: newService(config).withBaseURL(lineNotifyAPIURL),
	}
	return svc
}

// Post posts a payload to LINE.
func (s LINENotifyService) Post(contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", s.baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+s.service.config.AccessToken)
	return s.config.HTTPClient.Do(req)
}

// PostMessagePayload posts a message payload to LINE.
func (s LINENotifyService) PostMessagePayload(payload MessagePayload) (*http.Response, error) {
	buf := bytes.NewBufferString("message=" + payload.Message)
	return s.Post(bodyTypeURLEncoded, buf)
}
