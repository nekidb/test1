package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nekidb/test1/internal/config"
	"github.com/nekidb/test1/internal/server"
	"github.com/nekidb/test1/internal/shortener"
)

type SimpleStorage struct {
	data map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		data: make(map[string]string),
	}
}

func (s *SimpleStorage) Save(str1, str2 string) {
	s.data[str1] = str2
}

func (s *SimpleStorage) GetShortPath(str string) string {
	for k, v := range s.data {
		if v == str {
			return k
		}
	}

	return ""
}

func (s *SimpleStorage) GetSourceURL(str string) string {
	return s.data[str]
}

func (s *SimpleStorage) DeleteSourceURL(str string) {
	for k, v := range s.data {
		if v == str {
			delete(s.data, k)
		}
	}
}

func main() {
	config, err := config.Get(os.DirFS("."), "config.json")
	if err != nil {
		log.Fatal(err)
	}

	storage := NewSimpleStorage()
	shortener, _ := shortener.NewShortenerService(storage)

	server := server.NewServer(config.Host, config.Port, shortener)

	log.Printf("Starting server on %s%s", config.Host, config.Port)
	log.Println(http.ListenAndServe(config.Port, server))
}
