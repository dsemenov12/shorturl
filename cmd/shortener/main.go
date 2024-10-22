package main

import (
	"net/http"
	"github.com/dsemenov12/shorturl/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
    mux.HandleFunc(`/`, handlers.PostURL)
	mux.HandleFunc(`/EwHXdJfB`, handlers.Redirect)

	err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}
