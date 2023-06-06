package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nekidb/test1/internal/shortener"
)

var (
	errService    = errors.New("service error")
	errBadRequest = errors.New("bad request")
	errNotFound   = errors.New("not found")
)

type Server struct {
	host, port string
	shortener  *shortener.ShortenerService
	router     *chi.Mux
}

func NewServer(host, port string, shortener *shortener.ShortenerService) *Server {
	server := Server{
		host:      host,
		port:      port,
		shortener: shortener,
		router:    chi.NewRouter(),
	}

	server.initRouter()

	return &server
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *Server) initRouter() {
	s.router.Get("/s/{short}", s.redirectHandler)
	s.router.Post("/api/shorten", s.shortenHandler)

	fn := func(w http.ResponseWriter, r *http.Request) {
		s.handleError(w, http.StatusNotFound)
	}
	s.router.NotFound(http.HandlerFunc(fn))
}

type inputData struct {
	URL string `json:"url"`
}

type outputData struct {
	ShortURL string `json:"shortURL"`
}

func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	inputData := &inputData{}
	if err := json.NewDecoder(r.Body).Decode(inputData); err != nil {
		s.handleError(w, http.StatusBadRequest)
		return
	}

	srcURL := inputData.URL

	err := validateURL(srcURL)
	if err != nil {
		s.handleError(w, http.StatusBadRequest)
		return
	}

	shortPath, err := s.shortener.GetShortPath(srcURL)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	shortURL := path.Join(s.host+s.port, "s", shortPath)

	output, err := makeOutputJSON(shortURL)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortPath := chi.URLParam(r, "short")

	srcURL, err := s.shortener.GetSourceURL(shortPath)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	if srcURL == "" {
		s.handleError(w, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, srcURL, http.StatusFound)
}

func (s *Server) handleError(w http.ResponseWriter, code int) {
	customError := struct {
		Code int    `json:"code"`
		Msg  string `json:"error"`
	}{}

	customError.Code = code

	switch code {
	case http.StatusNotFound:
		customError.Msg = "Not found"
	case http.StatusInternalServerError:
		customError.Msg = "Internal server error"
	case http.StatusBadRequest:
		customError.Msg = "Bad request"
	}

	data, err := json.Marshal(&customError)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func makeOutputJSON(url string) ([]byte, error) {
	outputData := outputData{url}
	out, err := json.Marshal(outputData)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func validateURL(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}

	if u.Scheme == "" || !isGoodHost(u.Host) {
		return err
	}
	return nil
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
