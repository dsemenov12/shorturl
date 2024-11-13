package main

import (
	"fmt"
    "net/http"
    "net/url"
    "github.com/go-chi/chi/v5"
    "github.com/dsemenov12/shorturl/internal/handlers"
    "github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/logger"
	"github.com/dsemenov12/shorturl/internal/middlewares/gzipmiddleware"
	"go.uber.org/zap"
)



func main() {
    if error := run(); error != nil {
        fmt.Println(error)
    }
}

func run() error {
	config.ParseFlags()

    baseURL, error := url.Parse(config.FlagBaseAddr)
    if error != nil {
        return error
    }
    
    router := chi.NewRouter()

	if error = logger.Initialize(config.FlagLogLevel); error != nil {
        return error
    }
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

	router.Post("/api/shorten", logger.RequestLogger(gzipmiddleware.GzipMiddleware(handlers.ShortenPost)))
    router.Post("/", logger.RequestLogger(handlers.PostURL))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(handlers.Redirect))

	error = http.ListenAndServe(config.FlagRunAddr, router)
    if error != nil {
        return error
    }

	return nil
}