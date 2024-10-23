package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/config"
)

func main() {
	config.ParseFlags()

	router := chi.NewRouter()

	router.Post("/", handlers.PostURL)
	router.Get("/{id}", handlers.Redirect)

	err := http.ListenAndServe(config.FlagRunAddr, router)
    if err != nil {
        panic(err)
    }
}
