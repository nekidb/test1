package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Shortener interface {
	MakeShortPath() string
}

type Storage interface {
	Put(shortURL, srcURL string)
	GetSrcURL(shortPath string) (string, bool)
	GetShortPath(srcURL string) (string, bool)
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
	if r.URL.Path == "/" {
		inputData := &InputData{}
		err := json.NewDecoder(r.Body).Decode(inputData)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(inputData.URL))
	}
}

func (s *Server) redirectHandler(shortPath string, fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srcURL, ok := s.storage.GetSrcURL(shortPath)
		if ok {
			http.Redirect(w, r, srcURL, http.StatusFound)
			return
		}
		fallback(w, r)
	}
}

type InputData struct {
	URL string `json:"url"`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()
	router.HandleFunc("/", s.shortenerHandler)

	router.ServeHTTP(w, r)
}

func main() {
	storage := NewSimpleStorage()
	shortener := SimpleShortener{}
	server := NewServer(storage, shortener)
	http.ListenAndServe(":8080", server)
}
