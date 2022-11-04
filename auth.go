package spotifyclient

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	accountsBaseURL                = "https://accounts.spotify.com"
	ScopePlaylistModifyPublic      = "playlist-modify-public"
	ScopePlaylistModifyPrivate     = "playlist-modify-private"
	ScopePlaylistReadPrivate       = "playlist-read-private"
	ScopePlaylistReadCollaborative = "playlist-read-collaborative"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func CreatePKCEVerifierAndChallenge() (string, string, error) {
	verifier, err := generateRandomVerifier()
	if err != nil {
		return "", "", err
	}

	challenge := calculateChallenge(verifier)

	return string(verifier), challenge, nil
}

func GenerateRandomState() (string, error) {
	buf := make([]byte, 7)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	state := hex.EncodeToString(buf)
	return state, nil
}

func BuildPKCEAuthURI(clientID, redirectURI, challenge, state string, scopes ...string) string {
	q := url.Values{}
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", redirectURI)
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge", challenge)
	q.Add("state", state)
	q.Add("scope", strings.Join(scopes, " "))

	return accountsBaseURL + "/authorize?" + q.Encode()
}

func RequestPKCEToken(clientID, code, redirectURI, verifier string) (*Token, error) {
	query := make(url.Values)
	query.Set("client_id", clientID)
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_uri", redirectURI)
	query.Set("code_verifier", verifier)
	body := strings.NewReader(query.Encode())

	return postToken(body)
}

func RefreshPKCEToken(refreshToken, clientID string) (*Token, error) {
	query := make(url.Values)
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken)
	query.Set("client_id", clientID)
	body := strings.NewReader(query.Encode())

	return postToken(body)
}

func postToken(body io.Reader) (*Token, error) {
	res, err := http.Post(accountsBaseURL+"/api/token", "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	token := new(Token)
	err = json.NewDecoder(res.Body).Decode(token)

	return token, err
}
