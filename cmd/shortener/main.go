package main

import (
	"io"
	"fmt"
    "net/http"
    "net/url"
	"strings"
    "github.com/go-chi/chi/v5"
    "github.com/dsemenov12/shorturl/internal/handlers"
    "github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/logger"
	"go.uber.org/zap"
	"compress/gzip"
)

type gzipWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
} 

func gzipHandle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            // если gzip не поддерживается, передаём управление
            // дальше без изменений
            next.ServeHTTP(w, r)
            return
        }

        // создаём gzip.Writer поверх текущего w
        gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
        if err != nil {
            io.WriteString(w, err.Error())
            return
        }
        defer gz.Close()

        w.Header().Set("Content-Encoding", "gzip")
        // передаём обработчику страницы переменную типа gzipWriter для вывода данных
        next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
    })
}

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

	router.Post("/api/shorten", logger.RequestLogger(handlers.ShortenPost))
    router.Post("/", logger.RequestLogger(handlers.PostURL))
    router.Get(baseURL.Path + "/{id}", logger.RequestLogger(handlers.Redirect))

	error = http.ListenAndServe(config.FlagRunAddr, gzipHandle(router))
    if error != nil {
        return error
    }

	return nil
}