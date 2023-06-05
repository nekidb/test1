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
	router := chi.NewRouter()

	server := &Server{
		host:      host,
		port:      port,
		shortener: shortener,
		router:    router,
	}

	server.initRouter()

	return server
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *Server) initRouter() {
	s.router.Get("/s/{short}", s.errorWrapper(s.redirectHandler))
	s.router.Post("/api/shorten", s.errorWrapper(s.shortenHandler))
	s.router.NotFound(s.notFoundErrorHandler)
}

type inputData struct {
	URL string `json:"url"`
}

type outputData struct {
	ShortURL string `json:"shortURL"`
}

func (s *Server) errorWrapper(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		switch err {
		case errService:
			s.internalErrorHandler(w, r)
			// 	case errBadRequest:
			// 		s.badRequestErrorHandler(w, r)
		case errNotFound:
			s.notFoundErrorHandler(w, r)
		default:
			s.badRequestErrorHandler(w, r)
		}
	}
}

func (s *Server) internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) badRequestErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) notFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) error {
	inputData := &inputData{}
	if err := json.NewDecoder(r.Body).Decode(inputData); err != nil {
		return errBadRequest
	}

	srcURL := inputData.URL

	err := validateURL(srcURL)
	if err != nil {
		return errBadRequest
	}

	shortPath, err := s.shortener.GetShortPath(srcURL)
	if err != nil {
		return errService
	}

	shortURL := path.Join(s.host+s.port, "s", shortPath)

	output, err := makeOutputJSON(shortURL)
	if err != nil {
		return errService
	}

	w.Write(output)

	return nil
}

func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) error {
	shortPath := chi.URLParam(r, "short")

	srcURL, err := s.shortener.GetSourceURL(shortPath)
	if err != nil {
		return errService
	}

	if srcURL == "" {
		return errNotFound
	}

	http.Redirect(w, r, srcURL, http.StatusFound)
	return nil
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
