package server

import (
	"encoding/json"
	"net/http"
	"path"
)

type Shortener interface {
	MakeShortPath() string
	ValidateURL(url string) (bool, error)
}

type Storage interface {
	Put(shortPath, srcURL string) error
	GetSrcURL(shortPath string) (string, error)
	GetShortPath(srcURL string) (string, error)
}

type Server struct {
	host, port string
	storage    Storage
	shortener  Shortener
}

func NewServer(host, port string, storage Storage, shortener Shortener) *Server {
	return &Server{
		host:      host,
		port:      port,
		storage:   storage,
		shortener: shortener,
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

type inputData struct {
	URL string `json:"url"`
}

type outputData struct {
	ShortURL string `json:"shortURL"`
}

func (s *Server) shortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.badRequestHandler(w, r)
		return
	}

	inputData := &inputData{}
	err := json.NewDecoder(r.Body).Decode(inputData)
	if err != nil {
		s.internalErrorHandler(w, r)
		return
	}

	srcURL := inputData.URL

	ok, err := s.shortener.ValidateURL(srcURL)
	if err != nil {
		s.internalErrorHandler(w, r)
		return
	}
	if !ok {
		s.badRequestHandler(w, r)
		return
	}

	// Check if source URL exists in dabase
	shortPath, err := s.storage.GetShortPath(srcURL)
	if err != nil {
		s.internalErrorHandler(w, r)
	}
	if shortPath != "" {
		// If found, then return corresponding short URL
		shortURL := path.Join(s.host+s.port, shortPath)

		output, err := makeOutputJSON(shortURL)
		if err != nil {
			s.internalErrorHandler(w, r)
			return
		}
		w.Write(output)
	} else {
		// If not found in database, then make new short path
		shortPath = s.shortener.MakeShortPath()

		s.storage.Put(shortPath, srcURL)

		shortURL := path.Join(s.host+s.port, shortPath)

		output, err := makeOutputJSON(shortURL)
		if err != nil {
			s.internalErrorHandler(w, r)
			return
		}
		w.Write(output)

	}
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		s.notFoundHandler(w, r)
		return
	}

	srcURL, err := s.storage.GetSrcURL(path)
	if err != nil {
		s.internalErrorHandler(w, r)
	}

	if srcURL == "" {
		s.notFoundHandler(w, r)
		return
	}

	http.Redirect(w, r, srcURL, http.StatusFound)
}

func (s *Server) internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page not found"))
}

func makeOutputJSON(url string) ([]byte, error) {
	outputData := outputData{url}
	out, err := json.Marshal(outputData)
	if err != nil {
		return nil, err
	}

	return out, nil
}
