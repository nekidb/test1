package main

import (
	"net/http"
	"strings"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to best URL shortener!"))
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSuffix(strings.TrimPrefix(r.URL.RawQuery, "url=\""), "\"")
	w.Write([]byte(url))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/short":
		shortenerHandler(w, r)
	default:
		http.NotFound(w, r)
	}

}

// func generateURL() string {
// 	host := "localhost:8080/"
//
// 	letters := []rune("abcdefgABCDEFG")
// 	rnd := make([]rune, 5)
// 	for i := range rnd {
// 		rnd[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return host + string(rnd)
// }

func main() {
	server := NewServer()

	http.ListenAndServe(":8080", server)
}
