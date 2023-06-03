package shortener

import (
	"math/rand"
	"net/url"
	"strings"
)

type SimpleShortener struct{}

func (s SimpleShortener) GetShortPath() string {
	str := generateString(5)

	return "/" + str
}

func (s SimpleShortener) ValidateURL(str string) (bool, error) {
	u, err := url.Parse(str)
	if err != nil {
		return false, err
	}

	if u.Scheme == "" || !isGoodHost(u.Host) {
		return false, nil
	}
	return true, nil
}

func isGoodHost(host string) bool {
	if host == "" {
		return false
	}
	if !strings.Contains(host, ".") {
		return false
	}
	return true
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateString(n int) string {
	var sb strings.Builder
	sb.Grow(n)

	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}

	return sb.String()
}
