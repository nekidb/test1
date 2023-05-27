package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	shortener := StubShortener{}
	server := NewServer(shortener)

	t.Run("request home page", func(t *testing.T) {
		request := createGetRequest("/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "Welcome to best URL shortener!")
	})
	t.Run("not found page", func(t *testing.T) {
		request := createGetRequest("/somepage")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
		assertResponseBody(t, response.Body.String(), "404 page not found\n")
	})
}

func TestServerShortener(t *testing.T) {
	request := createGetRequest("/short?url=\"https://github.com/nekidb\"")
	response := httptest.NewRecorder()

	shortener := StubShortener{}
	server := NewServer(shortener)

	server.ServeHTTP(response, request)

	assertResponseBody(t, response.Body.String(), shortener.Short("https://githbub.com/nekidb"))
}

func TestServerRedirecting(t *testing.T) {
	shortener := StubShortener{}
	server := NewServer(shortener)

	request := createGetRequest("/short?url=\"https://github.com/nekidb\"")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	shortURL := response.Body.String()
	shortPath := strings.TrimPrefix(shortURL, "localhost:8080")

	request = createGetRequest(shortPath)
	response = httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatusCode(t, response.Code, http.StatusFound)
	assertLocation(t, response.Header().Get("Location"), "url=\"https://github.com/nekidb\"")
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

func createGetRequest(url string) *http.Request {
	return httptest.NewRequest(http.MethodGet, url, nil)
}

type StubShortener struct{}

func (s StubShortener) Short(URL string) string {
	return "localhost:8080/shorted"
}
