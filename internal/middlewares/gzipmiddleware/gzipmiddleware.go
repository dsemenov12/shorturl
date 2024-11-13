package gzipmiddleware

import (
	"strings"
    "net/http"
	"github.com/dsemenov12/shorturl/internal/compress/gzip"
)

func GzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ow := w

        acceptEncoding := r.Header.Get("Accept-Encoding")
        supportsGzip := strings.Contains(acceptEncoding, "gzip")
        if supportsGzip {
			if r.Header.Get("Content-Type") != "application/json" && r.Header.Get("Content-Type") != "text/html" {
				h.ServeHTTP(w, r)
				return
			}

            cw := gzip.NewCompressWriter(w)
            ow = cw

            defer cw.Close()
        }

        contentEncoding := r.Header.Get("Content-Encoding")
        sendsGzip := strings.Contains(contentEncoding, "gzip")
        if sendsGzip {
            cr, err := gzip.NewCompressReader(r.Body)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            r.Body = cr

            defer cr.Close()
        }

        h.ServeHTTP(ow, r)
    }
}