package handlers

import (
	"io"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/dsemenov12/shorturl/internal/config"
)

type ShortURLListMap map[string] string

var ShortURLList ShortURLListMap


func PostURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	ShortURLList = make(ShortURLListMap, 100)

	shortKey := util.RandStringBytes(8)
	shortURL := "http://" + config.FlagRunAddr + "/" + shortKey
	if (config.FlagBaseAddr != "") {
		shortURL = config.FlagBaseAddr + "/" + shortKey
	}

	body, err := io.ReadAll(req.Body)
	if (err != nil) {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	
	ShortURLList[shortKey] = string(body);

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length", "30")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func Redirect(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortKey := chi.URLParam(req, "id")

	http.Redirect(res, req, ShortURLList[shortKey], http.StatusTemporaryRedirect)
}