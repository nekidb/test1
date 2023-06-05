package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/nekidb/test1/internal/shortener"
)

type Server struct {
	host, port string
	shortener  *shortener.ShortenerService
}

func NewServer(host, port string, shortener *shortener.ShortenerService) *Server {
	return &Server{
		host:      host,
		port:      port,
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

	ok, err := validateURL(srcURL)
	if err != nil {
		s.internalErrorHandler(w, r)
		return
	}
	if !ok {
		s.badRequestHandler(w, r)
		return
	}

	shortPath, err := s.shortener.GetShortPath(srcURL)
	if err != nil {
		s.internalErrorHandler(w, r)
	}

	shortURL := path.Join(s.host+s.port, shortPath)

	output, err := makeOutputJSON(shortURL)
	if err != nil {
		s.internalErrorHandler(w, r)
		return
	}
	w.Write(output)
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		s.notFoundHandler(w, r)
		return
	}

	trimmed := strings.TrimPrefix(path, "/")
	srcURL, err := s.shortener.GetSourceURL(trimmed)
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

func validateURL(str string) (bool, error) {
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
