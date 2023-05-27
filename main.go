package main

import (
	"math/rand"
	"net/http"
	"strings"
)

type Shortener interface {
	Short(URL string) string
}

type Server struct {
	shortener Shortener
}

func NewServer(shortener Shortener) *Server {
	return &Server{shortener}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to best URL shortener!"))
}

func (s *Server) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	srcURL := strings.TrimSuffix(strings.TrimPrefix(r.URL.RawQuery, "url=\""), "\"")
	shortURL := s.shortener.Short(srcURL)
	w.Write([]byte(shortURL))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/short":
		s.shortenerHandler(w, r)
	default:
		http.NotFound(w, r)
	}

}

type SimpleShortener struct{}

func (s SimpleShortener) Short(URL string) string {
	host := "localhost:8080/"

	letters := []rune("abcdefgABCDEFG")
	rnd := make([]rune, 5)
	for i := range rnd {
		rnd[i] = letters[rand.Intn(len(letters))]
	}
	return host + string(rnd)
}

func main() {
	shortener := SimpleShortener{}
	server := NewServer(shortener)
	http.ListenAndServe(":8080", server)
}
