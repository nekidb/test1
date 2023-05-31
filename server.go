package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
)

type Shortener interface {
	MakeShortPath() string
}

type Storage interface {
	Put(shortURL, srcURL string)
	GetSrcURL(shortPath string) (string, bool)
	GetShortPath(srcURL string) (string, bool)
}

type InputData struct {
	URL string `json:"url"`
}

type OutputData struct {
	ShortURL string `json:"shortURL"`
}

type Server struct {
	host      string
	storage   Storage
	shortener Shortener
}

func NewServer(host string, storage Storage, shortener Shortener) *Server {
	return &Server{
		host:      host,
		storage:   storage,
		shortener: shortener,
	}
}

func (s *Server) badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page not found"))
}

func (s *Server) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.badRequestHandler(w, r)
	}

	inputData := &InputData{}
	err := json.NewDecoder(r.Body).Decode(inputData)
	if err != nil {
		s.badRequestHandler(w, r)
		return
	}
	if inputData.URL == "" {
		s.badRequestHandler(w, r)
		return
	}

	srcURL := inputData.URL
	shortPath, _ := s.storage.GetShortPath(srcURL)
	shortURL := path.Join(s.host, shortPath)

	output, err := makeOutputJSON(shortURL)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(output)
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		s.notFoundHandler(w, r)
		return
	}

	srcURL, ok := s.storage.GetSrcURL(path)
	if ok {
		http.Redirect(w, r, srcURL, http.StatusFound)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.shortenerHandler(w, r)
	case http.MethodGet:
		s.redirectHandler(w, r)
	}
}

func makeOutputJSON(url string) ([]byte, error) {
	outputData := OutputData{url}
	out, err := json.Marshal(outputData)
	if err != nil {
		return nil, err
	}

	return out, nil
}
