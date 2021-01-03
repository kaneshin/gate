package gate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path"
)

// pixelaAPIURL is an API URL for Pixela.
const pixelaAPIURL = "https://pixe.la/v1"

type (
	// PixelaService is a slack incoming webhook service.
	PixelaService struct {
		*service
	}
)

// NewPixelaService returns a new PixelaService.
func NewPixelaService(config *Config) *PixelaService {
	return &PixelaService{
		service: newService(config).withBaseURL(pixelaAPIURL),
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

// Put updates a graph.
func (s PixelaService) Put(id, suffix string) (*http.Response, error) {
	u := *s.baseURL
	u.Path = path.Join(u.Path, "users", s.config.ID, "graphs", id, suffix)
	req, err := http.NewRequest("PUT", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Length", "0")
	req.Header.Add("X-USER-TOKEN", s.service.config.AccessToken)
	return s.config.HTTPClient.Do(req)
}

// Increment increments quantity of a graph.
func (s PixelaService) Increment(payload GraphPayload) (*http.Response, error) {
	return s.Put(payload.ID, "increment")
}

// Decrement decrements quantity of a graph.
func (s PixelaService) Decrement(payload GraphPayload) (*http.Response, error) {
	return s.Put(payload.ID, "decrement")
}
