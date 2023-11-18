package spotifyclient

import (
	"bytes"
	"encoding/json"
	"net/url"
)

// UserService provides access to the Spotify Web API's user endpoints.
type UserService service

func (u *UserService) Me() (*Me, error) {
	Me := &Me{}
	err := u.client.get("v1", "/me", nil, Me)
	return Me, err
}

func (u *UserService) Devices() (*Devices, error) {
	Devices := &Devices{}
	err := u.client.get("v1", "/me/player/devices", nil, Devices)
	return Devices, err
}

func (u *UserService) Play(deviceID string, body *SetPlay) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)
	err = u.client.put("v1", "/me/player/play", url.Values{"device_id": {deviceID}}, reader)
	return err
}

func (u *UserService) Next(deviceID string) error {
	err := u.client.post("v1", "/me/player/next", url.Values{"device_id": {deviceID}}, nil, nil)
	return err
}

func (u *UserService) Previous(deviceID string) error {
	err := u.client.post("v1", "/me/player/previous", url.Values{"device_id": {deviceID}}, nil, nil)
	return err
}

func (u *UserService) Pause(deviceID string) error {
	err := u.client.put("v1", "/me/player/pause", url.Values{"device_id": {deviceID}}, nil)
	return err
}
