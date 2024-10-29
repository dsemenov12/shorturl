package main

import (
	"net/http"
	"net/url"
	"github.com/go-chi/chi/v5"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/config"
)

func main() {
    config.ParseFlags()

    baseURL, error := url.Parse(config.FlagBaseAddr)
    if error != nil {
        panic(error)
    }
    
    router := chi.NewRouter()

    router.Post("/", handlers.PostURL)
    router.Get(baseURL.Path + "/{id}", handlers.Redirect)

	errorServe := http.ListenAndServe(config.FlagRunAddr, router)
    if errorServe != nil {
        panic(errorServe)
    }
}
