package spotifyclient

func (a *API) Me() (*Me, error) {
	Me := &Me{}
	err := a.get("v1", "/me", nil, Me)
	return Me, err
}
