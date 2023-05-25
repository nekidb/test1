package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type URLPair struct {
	SrcURL   string `json:"srcURL"`
	ShortURL string `json:"shortURL"`
}

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

	shortURL := generateURL()
	pair := &URLPair{
		SrcURL:   jsonURL.URL,
		ShortURL: shortURL,
	}

	resultJSON, err := json.Marshal(pair)
	if err != nil {
		panic(err)
	}

	w.Write(resultJSON)
	s.storage.Write(jsonURL.URL, "")
}

func generateURL() string {
	host := "localhost:8080/"

	letters := []rune("abcdefgABCDEFG")
	rnd := make([]rune, 5)
	for i := range rnd {
		rnd[i] = letters[rand.Intn(len(letters))]
	}
	return host + string(rnd)
}

func main() {
}
