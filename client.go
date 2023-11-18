package spotifyclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Default Spotify API endpoint and scheme.
const (
	defaultHost   = "api.spotify.com"
	defaultScheme = "https"
)

type service struct {
	client *httpClient
}

// hyypConfig provides access to the Spotify Web API.
type httpClient struct {
	Host   string
	Scheme string
	token  string
}

// Client provides access to the Spotify Web API.
type Client struct {
	User     *UserService
	Playlist *PlaylistService
}

// NewClient creates a new Spotify Web API client.
func NewClient(token string, client *httpClient) *Client {
	if client == nil || client.Host == "" || client.Scheme == "" {
		client = &httpClient{
			Host:   defaultHost,
			Scheme: defaultScheme,
		}
	}
	client.token = token

	return &Client{
		User:     &UserService{client: client},
		Playlist: &PlaylistService{client: client},
	}
}

// Get sends a GET request to the Spotify Web API.
func (h *httpClient) get(apiVersion, endpoint string, query url.Values, res interface{}) error {
	return h.do(http.MethodGet, apiVersion, endpoint, query, nil, res)
}

// Post sends a POST request to the Spotify Web API.
func (h *httpClient) post(apiVersion, endpoint string, query url.Values, body io.Reader, res interface{}) error {
	return h.do(http.MethodPost, apiVersion, endpoint, query, body, res)
}

// Put sends a PUT request to the Spotify Web API.
func (h *httpClient) put(apiVersion, endpoint string, query url.Values, body io.Reader) error {
	return h.do(http.MethodPut, apiVersion, endpoint, query, body, nil)
}

// Delete sends a DELETE request to the Spotify Web API.
func (h *httpClient) delete(apiVersion, endpoint string, query url.Values) error {
	return h.do(http.MethodDelete, apiVersion, endpoint, query, nil, nil)
}

func (h *httpClient) do(method, apiVersion, endpoint string, query url.Values, body io.Reader, result interface{}) error {
	url := url.URL{
		Host:     h.Host,
		Path:     path.Join(apiVersion, endpoint),
		RawQuery: query.Encode(),
		Scheme:   h.Scheme,
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Success
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if result != nil {
			if err := json.NewDecoder(res.Body).Decode(result); err != nil {
				return err
			}
		}
		return nil
	}

	// Error
	spotifyErr := &struct {
		Error HttpError `json:"error"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(spotifyErr); err != nil {
		return err
	}

	return errors.New(spotifyErr.Error.Message)
}
