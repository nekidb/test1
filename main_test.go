package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const srcURL = `https://github.com/nekidb`

func TestServer(t *testing.T) {
	storage := &StubStorage{nil}
	shortener := StubShortener{}
	server := NewServer(storage, shortener)

	t.Run("server returns data when good request", func(t *testing.T) {
		reqBody := createRequestBody(t, srcURL)
		request := createPutRequest("/", reqBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), srcURL)
	})
}

// func TestServer(t *testing.T) {
// 	storage := &StubStorage{nil}
// 	shortener := StubShortener{}
// 	server := NewServer(storage, shortener)
//
// 	t.Run("not found page", func(t *testing.T) {
// 		request := createGetRequest("/somepage")
// 		response := httptest.NewRecorder()
//
// 		server.ServeHTTP(response, request)
//
// 		assertStatusCode(t, response.Code, http.StatusNotFound)
// 		assertResponseBody(t, response.Body.String(), "404 page not found\n")
// 	})
// }

// func TestServerRedirecting(t *testing.T) {
// 	shortPath, srcURL := "/shortpath", "https://github.com/nekidb"
// 	storage := &StubStorage{data: make(map[string]string)}
// 	storage.Put(shortPath, srcURL)
//
// 	shortener := StubShortener{}
// 	server := NewServer(storage, shortener)
//
// 	request := createGetRequest(shortPath)
// 	response := httptest.NewRecorder()
//
// 	server.ServeHTTP(response, request)
//
// 	assertStatusCode(t, response.Code, http.StatusFound)
// 	assertLocation(t, response.Header().Get("Location"), srcURL)
// }

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

func createRequestBody(t *testing.T, srcURL string) io.Reader {
	t.Helper()

	inputData := InputData{srcURL}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(inputData)
	if err != nil {
		t.Fatal("marshal test input data: ", err)
	}
	return buf
}

func createPutRequest(path string, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPost, path, body)
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
