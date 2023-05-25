package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Storage interface {
	Write(srcURL, shortURL string)
	Get(srcURL string) (string, bool)
}

type SourceURL struct {
	URL string `json:"url"`
}

type Server struct {
	storage Storage
}

func NewServer(storage Storage) *Server {
	return &Server{storage}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	jsonURL := SourceURL{}
	err = json.Unmarshal(requestBody, &jsonURL)
	if err != nil {
		panic(err)
	}

	s.storage.Write(jsonURL.URL, "")
}

func main() {
}
