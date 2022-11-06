package spotifyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Playlists service

func (p *Playlists) List() ([]*Playlist, error) {
	playlistPage := &struct {
		PagingMeta
		Items []*Playlist `json:"items"`
	}{}

	// TODO: Iterate over all pages of playlists

	err := p.client.get("v1", "/me/playlists", nil, playlistPage)
	return playlistPage.Items, err
}

func (p *Playlists) Create(userID, name string, public, collaborative bool, description string) (*Playlist, error) {
	query := make(url.Values)
	query.Add("user_id", userID)

	body := &struct {
		Name          string `json:"name"`
		Public        bool   `json:"public"`
		Collaborative bool   `json:"collaborative"`
		Description   string `json:"description"`
	}{
		Name:          name,
		Public:        public,
		Collaborative: collaborative,
		Description:   description,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	playlist := new(Playlist)
	err = p.client.post("v1", fmt.Sprintf("/users/%s/playlists", userID), query, bytes.NewReader(data), playlist)

	return playlist, err
}

func (p *Playlists) Fetch(id string) (*Playlist, error) {
	playlist := new(Playlist)
	err := p.client.get("v1", fmt.Sprintf("/playlists/%s", id), nil, playlist)
	return playlist, err
}

func (p *Playlists) Update(id, name, description string, public, collaborative bool) (*Playlist, error) {
	body := &struct {
		Name          string `json:"name"`
		Public        bool   `json:"public"`
		Collaborative bool   `json:"collaborative"`
		Description   string `json:"description"`
	}{
		Name:          name,
		Public:        public,
		Collaborative: collaborative,
		Description:   description,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	playlist := new(Playlist)
	err = p.client.put("v1", fmt.Sprintf("/playlists/%s", id), nil, bytes.NewReader(data))
	return playlist, err
}
