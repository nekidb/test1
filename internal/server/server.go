package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nekidb/test1/internal/shortener"
)

type Server struct {
	server    *http.Server
	shortener *shortener.ShortenerService
	router    *chi.Mux
}

func NewServer(shortener *shortener.ShortenerService) *Server {
	srv := Server{
		server:    &http.Server{},
		shortener: shortener,
		router:    chi.NewRouter(),
	}

	srv.initRouter()

	return &srv
}

func (s *Server) Serve(ln net.Listener) error {
	srv := &http.Server{
		Handler: s.router,
	}
	return srv.Serve(ln)
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) initRouter() {
	shortenRouter := chi.NewRouter()
	shortenRouter.Use(s.validationMiddleware)
	shortenRouter.Post("/", s.shortenHandler)
	shortenRouter.Delete("/", s.deleteSourceHandler)
	s.router.Mount("/api/shorten", shortenRouter)

	s.router.Get("/s/{short}", s.redirectHandler)

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

func (s *Server) validationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		c := context.WithValue(r.Context(), "srcURL", srcURL)
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	srcURL, ok := r.Context().Value("srcURL").(string)
	// srcURL, ok := "https://github.com/xenking", true
	if srcURL == "" || !ok {
		s.handleError(w, http.StatusBadRequest)
	}

	shortPath, err := s.shortener.GetShortPath(srcURL)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	// todo: need to store somewhere host and port
	shortURL := path.Join("localhost:8080", "s", shortPath)

	output, err := makeOutputJSON(shortURL)
	if err != nil {
		s.handleError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

func (s *Server) deleteSourceHandler(w http.ResponseWriter, r *http.Request) {
	srcURL, ok := r.Context().Value("srcURL").(string)
	if srcURL == "" || !ok {
		s.handleError(w, http.StatusBadRequest)
	}

	if err := s.shortener.DeleteSourceURL(srcURL); err != nil {
		s.handleError(w, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
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
