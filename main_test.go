package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"encoding/json"
)

const jsonString =  `{"url":"https://github.com/nekidb"}`

func TestReturnJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(jsonString))
	response := httptest.NewRecorder()

	storage := NewStubStorage()
	server := NewServer(storage)

	server.ServeHTTP(response, request)

	pair := &URLPair{}
	err := json.Unmarshal(response.Body.Bytes(), pair)
	if err != nil {
		t.Fatal(err)
	}

	got := pair.SrcURL
	want := "https://github.com/nekidb"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestStorageWrite(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(jsonString))
	response := httptest.NewRecorder()

	storage := NewStubStorage()
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

func NewStubStorage() *StubStorage {
	return &StubStorage{
		data: make(map[string]string),
	}
}

func (s *StubStorage) Write(srcURL, shortURL string) {
	s.data[srcURL] = shortURL
}

func (s *StubStorage) Get(srcURL string) (string, bool) {
	shortURL, ok := s.data[srcURL]
	return shortURL, ok 
}
