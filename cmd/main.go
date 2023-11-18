package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	spotify "github.com/josuerosadeavila/spotify-client"
)

const (
	ClientID     = "236940d9feb44ab685210db8eae09d4f"
	ClientSecret = "7207ce2a2ba84a44968b60302cd7a9f1"
	RedirectURI  = "http://localhost:1024/callback"
)

func Login() string {
	state, err := spotify.GenerateRandomState()
	if err != nil {
		panic(err)
	}

	scopes := []string{
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
	}
	uri := spotify.BuildAuthURI(ClientID, RedirectURI, state, scopes...)

	// if err := browser.OpenURL(uri); err != nil {
	// 	panic(err)
	// }

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		http.Get(RedirectURI)
		return nil
	}

	response, err := client.Do(req)
	if err == nil {
		if response.StatusCode == http.StatusFound { //status code 302
			fmt.Println(response.Location())
		}
	}

	code, err := ListenForCode(state)
	if err != nil {
		panic(err)
	}

	token, err := spotify.RequestToken(ClientID, ClientSecret, code, RedirectURI)
	if err != nil {
		panic(err)
	}

	return token.AccessToken
}

func ListenForCode(state string) (code string, err error) {
	server := &http.Server{Addr: ":1024"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state || r.URL.Query().Get("error") != "" {
			err = errors.New("authorization failed")
			fmt.Fprintln(w, "Failure.")
		} else {
			code = r.URL.Query().Get("code")
			fmt.Fprintln(w, "Success!")
		}

		// Use a separate thread so browser doesn't show a "No Connection" message
		go func() {
			server.Shutdown(context.Background())
		}()
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return "", err
	}

	return
}

func main() {
	token := Login()
	api := spotify.NewClient(token, nil)

	d, err := api.User.Devices()
	if err != nil {
		panic(err)
	}

	body := &spotify.SetPlay{
		ContextURI: "spotify:album:5ht7ItJgpBH7W6vJ5BqpPr",
		PositionMs: 0,
		Offset: struct {
			Position int "json:\"position\""
		}{
			Position: 0,
		},
	}

	err = api.User.Play(d.Devices[0].ID, body)
	if err != nil {
		panic(err)
	}

}
