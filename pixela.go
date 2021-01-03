package gate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/budougumi0617/pixela"
)

// pixelaAPIURL is an API URL for Pixela.
const pixelaAPIURL = "https://pixe.la/v1"

type (
	// PixelaService is a slack incoming webhook service.
	PixelaService struct {
		*service
		client *pixela.Client
	}
)

// NewPixelaService returns a new PixelaService.
func NewPixelaService(config *Config) *PixelaService {
	return &PixelaService{
		service: newService(config).withBaseURL(pixelaAPIURL),
		client:  pixela.New(config.ID, config.AccessToken),
	}
}

// Post posts a payload to pixela.
func (s PixelaService) Post(id string, body io.Reader) (*http.Response, error) {
	u := *s.baseURL
	u.Path = path.Join(u.Path, "users", s.config.ID, "graphs", id)
	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyTypeJSON)
	req.Header.Add("X-USER-TOKEN", s.service.config.AccessToken)
	return s.config.HTTPClient.Do(req)
}

// PostGraphPayload posts a graph payload to pixela.
func (s PixelaService) PostGraphPayload(payload GraphPayload) (*http.Response, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return s.Post(payload.ID, buf)
}
