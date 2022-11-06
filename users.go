package spotifyclient

type Users service

func (u *Users) Me() (*Me, error) {
	Me := &Me{}
	err := u.client.get("v1", "/me", nil, Me)
	return Me, err
}
