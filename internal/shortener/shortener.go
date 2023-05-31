package shortener

import (
	"math/rand"
	"net/url"
	"strings"
)

type SimpleShortener struct{}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (s SimpleShortener) MakeShortPath() string {
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
