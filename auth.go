package spotifyclient

import (
	"encoding/base64"
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
	ScopeUserReadPlaybackState     = "user-read-playback-state"
	ScopeUserModifyPlaybackState   = "user-modify-playback-state"
)

// Token represents an OAuth2 token.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// CreatePKCEVerifierAndChallenge creates a PKCE verifier and challenge.
func CreatePKCEVerifierAndChallenge() (string, string, error) {
	verifier, err := generateRandomVerifier()
	if err != nil {
		return "", "", err
	}

	challenge := calculateChallenge(verifier)

	return string(verifier), challenge, nil
}

// GenerateRandomState generates a random state string.
func GenerateRandomState() (string, error) {
	buf := make([]byte, 7)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	state := hex.EncodeToString(buf)
	return state, nil
}

// BuildAuthURI builds a authorization URI.
func BuildAuthURI(clientID, redirectURI, state string, scopes ...string) string {
	q := url.Values{}
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", redirectURI)
	q.Add("state", state)
	q.Add("scope", strings.Join(scopes, " "))

	return accountsBaseURL + "/authorize?" + q.Encode()
}

// BuildPKCEAuthURI builds a PKCE authorization URI.
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

func RequestToken(clientID, clientSecret, code, redirectURI string) (*Token, error) {
	query := make(url.Values)
	query.Set("client_id", clientID)
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_uri", redirectURI)
	body := strings.NewReader(query.Encode())

	return postToken(body, clientID, clientSecret)

}

// RequestPKCEToken requests a token using the PKCE flow.
func RequestPKCEToken(clientID, clientSecret, code, redirectURI, verifier string) (*Token, error) {
	query := make(url.Values)
	query.Set("client_id", clientID)
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_uri", redirectURI)
	query.Set("code_verifier", verifier)
	body := strings.NewReader(query.Encode())

	return postToken(body, clientID, clientSecret)
}

// RefreshPKCEToken refreshes a token using the PKCE flow.
func RefreshPKCEToken(refreshToken, clientID, clientSecret string) (*Token, error) {
	query := make(url.Values)
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken)
	query.Set("client_id", clientID)
	body := strings.NewReader(query.Encode())

	return postToken(body, clientID, clientSecret)

}

func postToken(body io.Reader, clientID, clientSecret string) (*Token, error) {
	req, err := http.NewRequest("POST", accountsBaseURL+"/api/token", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	token := new(Token)
	err = json.NewDecoder(res.Body).Decode(token)

	return token, err
}
