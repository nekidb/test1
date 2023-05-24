package main

import (
	"io/ioutil"
	"net/http"
)

func Server(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	w.Write(requestBody)
}

func main() {
}
