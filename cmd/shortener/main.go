package main

import (
	"io"
	"net/http"
)

type ShortURLListMap map[string] string

var ShortURLList ShortURLListMap

const shortKey = "EwHXdJfB"

func postRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	ShortURLList = make(ShortURLListMap, 1)

	shortURL := "http://localhost:8080/EwHXdJfB"

	body, err := io.ReadAll(req.Body)
	if (err != nil) {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	
	ShortURLList[shortKey] = string(body);

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length", "30")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func redirect(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(res, req, ShortURLList[shortKey], http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc(`/`, postRequest)
	mux.HandleFunc(`/EwHXdJfB`, redirect)

	err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}
