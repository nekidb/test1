package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SourceURL struct {
	URL string `json:"url"`
}

func Server(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	jsonURL := SourceURL{}
	err = json.Unmarshal(requestBody, &jsonURL)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(jsonURL.URL))
}

func main() {
}
