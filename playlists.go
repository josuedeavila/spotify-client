package spotifyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func (a *API) GetPlaylists() ([]*Playlist, error) {
	playlistPage := &struct {
		PagingMeta
		Items []*Playlist `json:"items"`
	}{}

	// TODO: Iterate over all pages of playlists

	err := a.get("v1", "/me/playlists", nil, playlistPage)
	return playlistPage.Items, err
}

func (a *API) CreatePlaylist(userID, name string, public, collaborative bool, description string) (*Playlist, error) {
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
	err = a.post("v1", fmt.Sprintf("/users/%s/playlists", userID), query, bytes.NewReader(data), playlist)

	return playlist, err
}

func (a *API) RetrivePlaylist(id string) (*Playlist, error) {
	playlist := new(Playlist)
	err := a.get("v1", fmt.Sprintf("/playlists/%s", id), nil, playlist)
	return playlist, err
}

func (a *API) DeletePlaylist(id string) (*Playlist, error) {
	playlist := new(Playlist)
	err := a.delete("v1", fmt.Sprintf("/playlists/%s", id), nil)
	return playlist, err
}
