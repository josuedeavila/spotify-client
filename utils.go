package spotifyclient

import (
	secure "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	ms, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	d.Duration = time.Duration(ms * 1000000)
	return nil
}

type HREF string

func (h *HREF) Get(api *API, obj interface{}) error {
	url, err := h.URL()
	if err != nil {
		return err
	}

	idx := strings.Index(url.Path, "/")
	version := url.Path[:idx]
	endpoint := url.Path[idx:]

	return api.get(version, endpoint, url.Query(), obj)
}

func (h *HREF) URL() (*url.URL, error) {
	return url.Parse(string(*h))
}

func (im *PagingMeta) Get(api *API, obj interface{}) error {
	return im.HREF.Get(api, obj)
}

func generateRandomVerifier() ([]byte, error) {
	seed, err := func() (int64, error) {
		buf := make([]byte, 8)
		_, err := secure.Read(buf)
		if err != nil {
			return 0, err
		}

		seed := int64(binary.BigEndian.Uint64(buf))
		return seed, nil
	}()

	if err != nil {
		return nil, err
	}
	rand.Seed(seed)

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_.-~"

	verifier := make([]byte, 128)
	for i := 0; i < len(verifier); i++ {
		idx := rand.Intn(len(chars))
		verifier[i] = chars[idx]
	}

	return verifier, nil
}

func calculateChallenge(verifier []byte) string {
	hash := sha256.Sum256(verifier)
	challenge := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(challenge, "=")
}
