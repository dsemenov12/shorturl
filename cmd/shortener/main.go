package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/logger"
	"github.com/dsemenov12/shorturl/internal/middlewares/gziphandler"
	"github.com/dsemenov12/shorturl/internal/storage/pg"
	storageMemory "github.com/dsemenov12/shorturl/internal/storage/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
    if error := run(); error != nil {
        fmt.Println(error)
    }
}

func run() error {
	var storage handlers.Storage

	config.ParseFlags()
	
    baseURL, err := url.Parse(config.FlagBaseAddr)
    if err != nil {
        return err
    }

	storage = storageMemory.NewStorage()
    if config.FlagDatabaseDSN != "" {
		conn, err := sql.Open("pgx", config.FlagDatabaseDSN)
		if err != nil {
			return err
		}

		storage = pg.NewStorage(conn)
    }

	if err = storage.Bootstrap(context.TODO()); err != nil {
		return err
	}

	app := handlers.NewApp(storage)
    
    router := chi.NewRouter()

	if err = logger.Initialize(config.FlagLogLevel); err != nil {
        return err
    }
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

    router.Post("/", logger.RequestLogger(app.PostURL))
	router.Get("/ping", logger.RequestLogger(app.Ping))
	router.Post("/api/shorten", logger.RequestLogger(app.ShortenPost))
	router.Post("/api/shorten/batch", logger.RequestLogger(app.ShortenBatchPost))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(app.Redirect))

	err = http.ListenAndServe(config.FlagRunAddr, gziphandler.GzipHandle(router))
    if err != nil {
        return err
    }

	return nil
}