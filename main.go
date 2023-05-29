package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Shortener interface {
	MakeShortPath() string
}

type Storage interface {
	Put(shortURL, srcURL string)
	GetSrcURL(shortURL string) (string, bool)
	GetShortURL(srcURL string) (string, bool)
}

type Server struct {
	storage   Storage
	shortener Shortener
}

func NewServer(storage Storage, shortener Shortener) *Server {
	return &Server{
		storage:   storage,
		shortener: shortener,
	}
}

func (s *Server) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting srcURL: ", r.URL.RawQuery)
	srcURL := strings.TrimSuffix(strings.TrimPrefix(r.URL.RawQuery, "url=%22"), "%22")
	shortPath := s.shortener.MakeShortPath()

	s.storage.Put(shortPath, srcURL)
	fmt.Println("writed to db: ", shortPath, srcURL)

	w.Write([]byte(shortPath))
}

func (s *Server) redirectHandler(shortPath string, fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srcURL, ok := s.storage.GetSrcURL(shortPath)
		// srcURL, ok := s.urls[shortPath]
		fmt.Println("trying to redirect: ", shortPath, srcURL, ok)
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
	fmt.Println("requested: ", path)
	router.HandleFunc("/", s.redirectHandler(path, http.NotFound))
	router.HandleFunc("/short", s.shortenerHandler)

	router.ServeHTTP(w, r)
}

func main() {
	// shortener := SimpleShortener{}
	// server := NewServer(shortener)
	// http.ListenAndServe(":8080", server)
}
