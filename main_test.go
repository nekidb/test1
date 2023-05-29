package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	storage := &StubStorage{nil}
	shortener := StubShortener{}
	server := NewServer(storage, shortener)

	t.Run("not found page", func(t *testing.T) {
		request := createGetRequest("/somepage")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
		assertResponseBody(t, response.Body.String(), "404 page not found\n")
	})
}

// func TestServerShortener(t *testing.T) {
// 	request := createGetRequest("/short?url=\"https://github.com/nekidb\"")
// 	response := httptest.NewRecorder()
//
// 	shortener := StubShortener{}
// 	server := NewServer(shortener)
//
// 	server.ServeHTTP(response, request)
//
// 	assertResponseBody(t, response.Body.String(), shortener.Short("https://githbub.com/nekidb"))
// }

func TestServerRedirecting(t *testing.T) {
	shortPath, srcURL := "/shortpath", "https://github.com/nekidb"
	storage := &StubStorage{data: make(map[string]string)}
	storage.Put(shortPath, srcURL)

	shortener := StubShortener{}
	server := NewServer(storage, shortener)

	request := createGetRequest(shortPath)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatusCode(t, response.Code, http.StatusFound)
	assertLocation(t, response.Header().Get("Location"), srcURL)
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatusCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertLocation(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got location %q, want %q", got, want)
	}
}

func createGetRequest(path string) *http.Request {
	return httptest.NewRequest(http.MethodGet, path, nil)
}

type StubShortener struct{}

func (s StubShortener) MakeShortPath() string {
	return "/shorted"
}

type StubStorage struct {
	data map[string]string
}

func (s *StubStorage) Put(shortURL, srcURL string) {
	s.data[shortURL] = srcURL
}

func (s *StubStorage) GetSrcURL(shortPath string) (string, bool) {
	srcURL, ok := s.data[shortPath]
	return srcURL, ok
}

func (s *StubStorage) GetShortPath(srcURL string) (string, bool) {
	for k, v := range s.data {
		if v == srcURL {
			return k, true
		}
	}

	return "", false
}
