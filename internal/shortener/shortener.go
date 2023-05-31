package shortener

import "math/rand"

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
