package main

import (
	"net/http"
)

func main() {
	storage := NewSimpleStorage()
	storage.Put("/lolkek", "https://github.com/nekidb")
	shortener := SimpleShortener{}
	server := NewServer("localhost:8080", storage, shortener)
	http.ListenAndServe(":8080", server)
}
