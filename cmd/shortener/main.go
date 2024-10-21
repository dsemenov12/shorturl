package main

import (
	"net/http"
)

func postRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortUrl := "http://localhost:8080/EwHXdJfB"

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length", "30")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortUrl))
}

func redirect(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(res, req, "https://practicum.yandex.ru/", http.StatusMovedPermanently)
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
