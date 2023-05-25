package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

const jsonString =  `{"url":"https://github.com/nekidb"}`

// func TestServer(t *testing.T) {
// 	request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(jsonString))
// 	response := httptest.NewRecorder()
// 
// 	server := NewServer()
// 
// 	server.ServeHTTP(response, request)
// 
// 	got := response.Body.String()
// 	want := "https://github.com/nekidb"
// 
// 	if got != want {
// 		t.Errorf("got %v, want %v", got, want)
// 	}
// }

func TestStorageWrite(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(jsonString))
	response := httptest.NewRecorder()

	storage := &StubStorage{
		data: make(map[string]string),
	}
	server := NewServer(storage)

	server.ServeHTTP(response, request)

	_, got := storage.Get("https://github.com/nekidb")
	want := true
	if got != want{
		t.Errorf("got %v, want %v", got, want)
	}
}

type StubStorage struct {
	data map[string]string
}

func (s *StubStorage) Write(srcURL, shortURL string) {
	s.data[srcURL] = shortURL
}

func (s *StubStorage) Get(srcURL string) (string, bool) {
	shortURL, ok := s.data[srcURL]
	return shortURL, ok 
}
