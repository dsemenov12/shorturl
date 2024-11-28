package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/logger"
	"github.com/dsemenov12/shorturl/internal/middlewares/gziphandler"
	"github.com/dsemenov12/shorturl/internal/storage/pg"
	"github.com/go-chi/chi/v5"
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

	conn, err := sql.Open("pgx", config.FlagDatabaseDSN)
    if err != nil {
        return err
    }

	handlers.Storage = pg.NewStorage(conn)

    if err = handlers.Storage.Bootstrap(context.TODO()); err != nil {
        return err
    }
    
    router := chi.NewRouter()

	if err = logger.Initialize(config.FlagLogLevel); err != nil {
        return err
    }
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

    router.Post("/", logger.RequestLogger(handlers.PostURL))
	router.Get("/ping", logger.RequestLogger(handlers.Ping))
	router.Post("/api/shorten", logger.RequestLogger(handlers.ShortenPost))
	router.Post("/api/shorten/batch", logger.RequestLogger(handlers.ShortenBatchPost))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(handlers.Redirect))

	err = http.ListenAndServe(config.FlagRunAddr, gziphandler.GzipHandle(router))
    if err != nil {
        return err
    }

	return nil
}