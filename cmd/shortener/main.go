package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/handlers"
	"github.com/dsemenov12/shorturl/internal/middlewares/authcookiehandler"
	"github.com/dsemenov12/shorturl/internal/middlewares/authhandler"
	"github.com/dsemenov12/shorturl/internal/middlewares/gziphandler"
	"github.com/dsemenov12/shorturl/internal/middlewares/logger"
	"github.com/dsemenov12/shorturl/internal/storage"
	"github.com/dsemenov12/shorturl/internal/storage/memory"
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
	var storage storage.Storage

	config.ParseFlags()

	baseURL, err := url.Parse(config.FlagBaseAddr)
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	go func() {
		logger.Log.Info("Starting pprof", zap.String("address", ":6060"))
		if err := http.ListenAndServe(":6060", nil); err != nil {
			logger.Log.Error("pprof server failed", zap.Error(err))
		}
	}()

	storage = memory.NewStorage()
	if config.FlagDatabaseDSN != "" {
		conn, err := sql.Open("pgx", config.FlagDatabaseDSN)
		if err != nil {
			return err
		}

		router.Get("/ping", logger.RequestLogger(func(res http.ResponseWriter, req *http.Request) {
			if err := conn.Ping(); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}

			res.WriteHeader(http.StatusOK)
		}))

		storage = pg.NewStorage(conn)
	}

	if err = storage.Bootstrap(context.TODO()); err != nil {
		return err
	}

	app := handlers.NewApp(storage)

	if err = logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

	router.Post("/", logger.RequestLogger(authhandler.AuthHandle(app.PostURL)))
	router.Post("/api/shorten", logger.RequestLogger(authhandler.AuthHandle(app.ShortenPost)))
	router.Post("/api/shorten/batch", logger.RequestLogger(authhandler.AuthHandle(app.ShortenBatchPost)))
	router.Get(baseURL.Path+"/{id}", logger.RequestLogger(app.Redirect))
	router.Get("/api/user/urls", logger.RequestLogger(authhandler.AuthHandle(app.UserUrls)))
	router.Delete("/api/user/urls", logger.RequestLogger(authcookiehandler.AuthCookieHandle(app.DeleteUserUrls)))

	err = http.ListenAndServe(
		config.FlagRunAddr,
		gziphandler.GzipHandle(router),
	)
	if err != nil {
		return err
	}

	return nil
}
