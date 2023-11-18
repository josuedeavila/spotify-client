package spotifyclient

import (
	secure "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

type HREF string

func (h *HREF) Get(c *httpClient, obj interface{}) error {
	url, err := h.URL()
	if err != nil {
		return err
	}

	idx := strings.Index(url.Path, "/")
	version := url.Path[:idx]
	endpoint := url.Path[idx:]

	// return c.Get(version, endpoint, url.Query(), obj)
	return c.get(version, endpoint, url.Query(), obj)
}

func (h *HREF) URL() (*url.URL, error) {
	return url.Parse(string(*h))
}

func (im *PagingMeta) Get(c *httpClient, obj interface{}) error {
	return im.HREF.Get(c, obj)
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
