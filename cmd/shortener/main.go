package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/dsemenov12/shorturl/internal/handlers"
)

func main() {
	router := chi.NewRouter()

	router.Post("/", handlers.PostURL)
	router.Get("/{id}", handlers.Redirect)

	err := http.ListenAndServe(`:8080`, router)
    if err != nil {
        panic(err)
    }
}
