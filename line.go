package gate

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kaneshin/gate/line"
)

type (
	// LINENotifyService is a slack incoming webhook service.
	LINENotifyService struct {
		*service
	}
)

// NewLINENotifyService returns a new LINENotifyService.
func NewLINENotifyService(config *Config) *LINENotifyService {
	svc := &LINENotifyService{
		service: newService(config).withBaseURL(line.NotifyAPIURL),
	}
	return svc
}

// NewPayload returns a new Payload.
func (s LINENotifyService) NewPayload(text string) line.Payload {
	p := line.Payload{
		Message: text,
	}
	return p
}

// Post posts data to LINE.
func (s LINENotifyService) Post(v interface{}) (*http.Response, error) {
	var body io.Reader
	switch v := v.(type) {
	case io.Reader:
		body = v

	case string:
		buf := bytes.NewBufferString("message=" + v)
		body = buf

	case line.Payload:
		buf := bytes.NewBufferString("message=" + v.Message)
		body = buf

	case *line.Payload:
		buf := bytes.NewBufferString("message=" + v.Message)
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
