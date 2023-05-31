package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		host     = "localhost"
		port     = ":8080"
		database = "urls.db"
	)
	storage, err := NewBoltStorage(database)
	if err != nil {
		log.Fatal(err)
	}
	shortener := SimpleShortener{}
	server := NewServer(host, port, storage, shortener)

	log.Printf("Starting server on %s%s", host, port)
	log.Println(http.ListenAndServe(port, server))
}
