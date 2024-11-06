package main

import (
    "net/http"
    "net/url"
    "github.com/go-chi/chi/v5"
    "github.com/dsemenov12/shorturl/internal/handlers"
    "github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/logger"
	"go.uber.org/zap"
)

func main() {
    config.ParseFlags()

    baseURL, error := url.Parse(config.FlagBaseAddr)
    if error != nil {
        panic(error)
    }
    
    router := chi.NewRouter()

	if errorLoger := logger.Initialize(config.FlagLogLevel); errorLoger != nil {
        panic(errorLoger)
    }
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

    router.Post("/", logger.RequestLogger(handlers.PostURL))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(handlers.Redirect))

	errorServe := http.ListenAndServe(config.FlagRunAddr, router)
    if errorServe != nil {
        panic(errorServe)
    }
}
