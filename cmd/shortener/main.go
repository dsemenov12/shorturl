package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os/signal"
	"syscall"
	"time"

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

// Глобальные переменные для информации о сборке
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildData()

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

	// Контекст с отменой
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

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

	router.Post("/", logger.RequestLogger(authhandler.AuthHandle(app.PostURL)))
	router.Post("/api/shorten", logger.RequestLogger(authhandler.AuthHandle(app.ShortenPost)))
	router.Post("/api/shorten/batch", logger.RequestLogger(authhandler.AuthHandle(app.ShortenBatchPost)))
	router.Get(baseURL.Path+"/{id}", logger.RequestLogger(app.Redirect))
	router.Get("/api/user/urls", logger.RequestLogger(authhandler.AuthHandle(app.UserUrls)))
	router.Delete("/api/user/urls", logger.RequestLogger(authcookiehandler.AuthCookieHandle(app.DeleteUserUrls)))

	server := &http.Server{
		Addr:    config.FlagRunAddr,
		Handler: gziphandler.GzipHandle(router),
	}

	go func() {
		logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

		var err error
		if config.FlagEnableHTTPS {
			certFile := "cert.pem"
			keyFile := "key.pem"
			err = server.ListenAndServeTLS(certFile, keyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	logger.Log.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if config.FlagEnableHTTPS {
		certFile := "cert.pem"
		keyFile := "key.pem"

		logger.Log.Info("Starting HTTPS server", zap.String("address", config.FlagRunAddr))
		return http.ListenAndServeTLS(config.FlagRunAddr, certFile, keyFile, gziphandler.GzipHandle(router))
	}

	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))

	err = http.ListenAndServe(
		config.FlagRunAddr,
		gziphandler.GzipHandle(router),
	)
	if err != nil {
		return err
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server shutdown failed", zap.Error(err))
	} else {
		logger.Log.Info("Server exited properly")
	}

	return nil
}

// printBuildData - вывод информации о сборке.
func printBuildData() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}
