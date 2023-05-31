package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	host      = "localhost"
	port      = ":8080"
	srcURL    = "https://github.com/nekidb"
	shortPath = "/shorted"
)

func TestBadRequest(t *testing.T) {
	storage := &StubStorage{nil}
	shortener := StubShortener{}
	server := NewServer(host, port, storage, shortener)

	t.Run("POST request with bad input data", func(t *testing.T) {
		request := createPostRequest(t, "/", "")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
	t.Run("POST request with wrong path", func(t *testing.T) {
		request := createPostRequest(t, "/somepage", srcURL)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
	t.Run("empty GET request", func(t *testing.T) {
		request := createGetRequest("/")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
		assertResponseBody(t, response.Body.String(), "Page not found")
	})

}

func TestServer(t *testing.T) {
	storage := NewStubStorage()
	storage.Put(shortPath, srcURL)

	shortener := StubShortener{}
	server := NewServer(host, port, storage, shortener)

	t.Run("server returns correct shortURL", func(t *testing.T) {
		request := createPostRequest(t, "/", srcURL)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := createExpectedOutput(t, host, port, shortPath)
		assertStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), want)
	})
	t.Run("redirect to source URL", func(t *testing.T) {
		request := createGetRequest(shortPath)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusFound)
		assertLocation(t, response.Header().Get("Location"), srcURL)

	})
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

func createRequestBody(t *testing.T, inputData any) io.Reader {
	t.Helper()

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(inputData)
	if err != nil {
		t.Fatal("marshal test input data: ", err)
	}
	return buf
}

func createPostRequest(t *testing.T, path, srcURL string) *http.Request {
	inputData := InputData{srcURL}
	reqBody := createRequestBody(t, inputData)
	return httptest.NewRequest(http.MethodPost, path, reqBody)
}

func createGetRequest(path string) *http.Request {
	return httptest.NewRequest(http.MethodGet, path, nil)
}

func createExpectedOutput(t *testing.T, host, port, shortPath string) string {
	t.Helper()

	v := OutputData{host + port + shortPath}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

type StubShortener struct{}

func (s StubShortener) MakeShortPath() string {
	return "/shorted"
}

type StubStorage struct {
	data map[string]string
}

func NewStubStorage() *StubStorage {
	return &StubStorage{
		data: make(map[string]string),
	}
}

func (s *StubStorage) Put(shortURL, srcURL string) error {
	s.data[shortURL] = srcURL
	return nil
}

func (s *StubStorage) GetSrcURL(shortPath string) (string, error) {
	srcURL, _ := s.data[shortPath]
	return srcURL, nil
}

func (s *StubStorage) GetShortPath(srcURL string) (string, error) {
	for k, v := range s.data {
		if v == srcURL {
			return k, nil
		}
	}

	return "", nil
}
