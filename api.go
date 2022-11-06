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

const (
	host = "api.spotify.com"
	scheme  = "https"
)

type service struct {
	client *Client
}

type Client struct {
	token  string
	host   string
	scheme string

	Playlists *Playlists
}

func NewAPI(token string) *Client {
	return &Client{
		token:  token,
		host:   host,
		scheme: scheme,
	}
}

func (c *Client) get(apiVersion, endpoint string, query url.Values, res interface{}) error {
	return c.Do(http.MethodGet, apiVersion, endpoint, query, nil, res)
}

func (c *Client) post(apiVersion, endpoint string, query url.Values, body io.Reader, res interface{}) error {
	return c.Do(http.MethodPost, apiVersion, endpoint, query, body, res)
}

func (c *Client) put(apiVersion, endpoint string, query url.Values, body io.Reader) error {
	return c.Do(http.MethodPut, apiVersion, endpoint, query, body, nil)
}

func (c *Client) delete(apiVersion, endpoint string, query url.Values) error {
	return c.Do(http.MethodDelete, apiVersion, endpoint, query, nil, nil)
}

func (c *Client) Do(method, apiVersion, endpoint string, query url.Values, body io.Reader, result interface{}) error {
	url := url.URL{
		Host:     c.host,
		Path:     path.Join(apiVersion, endpoint),
		RawQuery: query.Encode(),
		Scheme:   c.scheme,
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

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
