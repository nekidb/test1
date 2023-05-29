package main

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

func (s SimpleShortener) Short(URL string) string {
	host := "localhost:8080/"

	str := generateString(5)

	return host + str
}
