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
	apiHost = "api.spotify.com"
	scheme  = "https"
)

type API struct {
	token  string
	host   string
	scheme string
}

func NewAPI(token string) *API {
	return &API{
		token:  token,
		host:   apiHost,
		scheme: scheme,
	}
}

func (a *API) get(apiVersion, endpoint string, query url.Values, res interface{}) error {
	return a.Do(http.MethodGet, apiVersion, endpoint, query, nil, res)
}

func (a *API) post(apiVersion, endpoint string, query url.Values, body io.Reader, res interface{}) error {
	return a.Do(http.MethodPost, apiVersion, endpoint, query, body, res)
}

func (a *API) put(apiVersion, endpoint string, query url.Values, body io.Reader) error {
	return a.Do(http.MethodPut, apiVersion, endpoint, query, body, nil)
}

func (a *API) delete(apiVersion, endpoint string, query url.Values) error {
	return a.Do(http.MethodDelete, apiVersion, endpoint, query, nil, nil)
}

func (a *API) Do(method, apiVersion, endpoint string, query url.Values, body io.Reader, result interface{}) error {
	url := url.URL{
		Host:     a.host,
		Path:     path.Join(apiVersion, endpoint),
		RawQuery: query.Encode(),
		Scheme:   a.scheme,
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.token))

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
