package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

type Shortener interface {
	Short(URL string) string
}

type Server struct {
	urls      map[string]string
	shortener Shortener
}

func NewServer(shortener Shortener) *Server {
	return &Server{
		urls:      make(map[string]string),
		shortener: shortener,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to best URL shortener!"))
}

func (s *Server) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	srcURL := strings.TrimSuffix(strings.TrimPrefix(r.URL.RawQuery, "url=\\%22"), "\\%22")
	shortURL := s.shortener.Short(srcURL)

	s.urls[shortURL] = srcURL

	fmt.Println("writed: ", shortURL, srcURL)

	w.Write([]byte(shortURL))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch path {
	case "/":
		homeHandler(w, r)
	case "/short":
		s.shortenerHandler(w, r)
	default:
		shortURL := "localhost:8080" + path
		srcURL, ok := s.urls[shortURL]
		fmt.Println(shortURL, srcURL, ok)
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, srcURL, http.StatusFound)
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
