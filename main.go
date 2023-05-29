package main

import (
	"fmt"
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

func (s *Server) redirectHandler(shortURL string, fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srcURL, ok := s.urls[shortURL]
		fmt.Println(shortURL, srcURL, ok)
		if ok {
			http.Redirect(w, r, srcURL, http.StatusFound)
			return
		}
		fallback(w, r)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	path := r.URL.Path
	shortURL := "localhost:8080" + path

	router.HandleFunc("/", s.redirectHandler(shortURL, http.NotFound))
	router.HandleFunc("/short", s.shortenerHandler)

	router.ServeHTTP(w, r)
}

func main() {
	shortener := SimpleShortener{}
	server := NewServer(shortener)
	http.ListenAndServe(":8080", server)
}
