package handlers

import (
	"io"
	"net/http"
	"strconv"
	"bytes"
	"encoding/json"
    
    "github.com/dsemenov12/shorturl/internal/structs/storage"
	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/go-chi/chi/v5"
)

var storageObj = storage.Storage{Data: make(map[string]string)}

func ShortenPost(res http.ResponseWriter, req *http.Request) {
	var inputDataValue models.InputData
	var buf bytes.Buffer

	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &inputDataValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	storageObj.Set(shortKey, inputDataValue.URL)

	var result = models.ResultJSON{
		Result: shortURL,
	}

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

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