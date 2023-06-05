package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nekidb/test1/internal/config"
	"github.com/nekidb/test1/internal/server"
	"github.com/nekidb/test1/internal/shortener"
	"github.com/nekidb/test1/internal/storage"
)

func main() {
	config, err := config.Get(os.DirFS("."), "config.json")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewBoltStorage(config.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	shortener, _ := shortener.NewShortenerService(storage)

	server := server.NewServer(config.Host, config.Port, shortener)

	log.Printf("Starting server on %s%s", config.Host, config.Port)
	log.Println(http.ListenAndServe(config.Port, server))
}
