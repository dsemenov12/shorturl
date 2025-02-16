package gziphandler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipWriter оборачивает http.ResponseWriter, позволяя записывать данные в сжатом формате GZIP.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write записывает данные в gzip-формате.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandle является middleware-функцией, которая обрабатывает сжатие и распаковку данных в формате GZIP.
// Если запрос клиента поддерживает сжатие GZIP, то данные, отправляемые сервером, будут сжаты в формат GZIP.
// Если запрос клиента уже использует GZIP для передачи данных, то данные будут распакованы перед обработкой.
//
// next: следующий обработчик в цепочке, который будет вызван после выполнения логики сжатия или распаковки.
//
// Возвращаемое значение: возвращает обработчик HTTP-запроса, который проверяет поддержку GZIP-сжатия в запросе,
// выполняет сжатие ответа и передает управление следующему обработчику.
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		sendsGzip := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if sendsGzip {
			cr, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
