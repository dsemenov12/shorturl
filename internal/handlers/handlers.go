package handlers

import (
	"io"
	"net/http"
)

type ShortURLListMap map[string] string

var ShortURLList ShortURLListMap

const shortKey = "EwHXdJfB"

func PostURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	ShortURLList = make(ShortURLListMap, 1)

	shortURL := "http://" + req.Host + "/" + shortKey

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

	http.Redirect(res, req, ShortURLList[shortKey], http.StatusTemporaryRedirect)
}