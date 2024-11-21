package main

import (
	"fmt"
    "net/http"
    "net/url"

    "github.com/go-chi/chi/v5"
    "github.com/dsemenov12/shorturl/internal/handlers"
    "github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/logger"
	"github.com/dsemenov12/shorturl/internal/middlewares/gziphandler"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"go.uber.org/zap"
)

func main() {
    if error := run(); error != nil {
        fmt.Println(error)
    }
}

func run() error {
	config.ParseFlags()
	filestorage.Load()

    baseURL, err := url.Parse(config.FlagBaseAddr)
    if err != nil {
        return err
    }
    
    router := chi.NewRouter()

	if err = logger.Initialize(config.FlagLogLevel); err != nil {
        return err
    }
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

	router.Get("/ping", logger.RequestLogger(handlers.Ping))
	router.Post("/api/shorten", logger.RequestLogger(handlers.ShortenPost))
    router.Post("/", logger.RequestLogger(handlers.PostURL))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(handlers.Redirect))

	err = http.ListenAndServe(config.FlagRunAddr, gziphandler.GzipHandle(router))
    if err != nil {
        return err
    }

	return nil
}