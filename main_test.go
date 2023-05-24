package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	json := `{"url":"https://github.com/nekidb"}`
	body := &bytes.Buffer{}
	body.WriteString(json)
	request := httptest.NewRequest(http.MethodGet, "/", body)
	response := httptest.NewRecorder()

	server := NewServer()

	server.ServeHTTP(response, request)

	got := response.Body.String()
	want := "https://github.com/nekidb"

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
