package main

import (
	"log"
	"net/http"

	"github.com/nekidb/test1/internal/server"
	"github.com/nekidb/test1/internal/shortener"
	"github.com/nekidb/test1/internal/storage"
)

func main() {
	const (
		host     = "localhost"
		port     = ":8080"
		database = "urls.db"
	)
	storage, err := storage.NewBoltStorage(database)
	if err != nil {
		log.Fatal(err)
	}
	shortener := shortener.SimpleShortener{}
	server := server.NewServer(host, port, storage, shortener)

	log.Printf("Starting server on %s%s", host, port)
	log.Println(http.ListenAndServe(port, server))
}
