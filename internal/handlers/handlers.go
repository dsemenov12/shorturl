package handlers

import (
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/go-chi/chi/v5"
)

type storage struct {
	mx sync.Mutex
    data map[string]string
}

func (s *storage) Get(key string) (string, error) {
	s.mx.Lock()
    defer s.mx.Unlock()
    return s.data[key], nil
}

func (s *storage) Set(key string, value string) {
	s.mx.Lock()
    defer s.mx.Unlock()
	s.data[key] = value
}

var storageObj = storage{data: make(map[string]string)}

func PostURL(res http.ResponseWriter, req *http.Request) {
	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	storageObj.Set(shortKey, string(body))

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length",  strconv.Itoa(len(shortURL)))
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")
	redirectLink, err := storageObj.Get(shortKey)
	if err != nil {
		http.Error(res, "redirect not found", 404)
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}