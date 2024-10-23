package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/go-chi/chi/v5"
)

type ShortURLListMap map[string] string

var ShortURLList ShortURLListMap


func PostURL(res http.ResponseWriter, req *http.Request) {
	ShortURLList = make(ShortURLListMap, 100)

	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	body, err := io.ReadAll(req.Body)
	if (err != nil) {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	
	ShortURLList[shortKey] = string(body);

	res.Header().Set("Content-Type", "text/plain")
	res.Header().Set("Content-Length",  strconv.Itoa(len(shortURL)))
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	http.Redirect(res, req, ShortURLList[shortKey], http.StatusTemporaryRedirect)
}