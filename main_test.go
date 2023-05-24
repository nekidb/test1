package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	json := `{url:"https://github.com/nekidb"}`
	body := &bytes.Buffer{}
	body.WriteString(json)
	request := httptest.NewRequest(http.MethodGet, "/", body)
	response := httptest.NewRecorder()

	Server(response, request)

	got := response.Body.String()

	want := `{url:"https://github.com/nekidb"}`
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
