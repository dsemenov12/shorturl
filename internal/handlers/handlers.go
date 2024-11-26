package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/storage/pg"
	"github.com/dsemenov12/shorturl/internal/structs/storage"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Storage *pg.Storage

func Ping(res http.ResponseWriter, req *http.Request) {
	if err := Storage.Ping(); err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
    }

	res.WriteHeader(http.StatusOK)
}

func ShortenPost(res http.ResponseWriter, req *http.Request) {
	var inputDataValue models.InputData

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
	if err = json.Unmarshal(body, &inputDataValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	storage.StorageObj.Set(shortKey, inputDataValue.URL)

	var result = models.ResultJSON{
		Result: shortURL,
	}

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	if config.FlagDatabaseDSN != "" {
		Storage.Insert(req.Context(), shortKey, inputDataValue.URL)
	} else {
		filestorage.Save(storage.StorageObj.Data)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func ShortenBatchPost(res http.ResponseWriter, req *http.Request) {
	var batch []models.BatchItem
	var result []models.BatchResultItem

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &batch); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	for _, batchItem := range batch {
		shortURL := config.FlagBaseAddr + "/" + batchItem.CorrelationId

		Storage.Insert(req.Context(), batchItem.CorrelationId, batchItem.OriginalURL)

		result = append(result, models.BatchResultItem{
			CorrelationId: batchItem.CorrelationId,
			ShortURL: shortURL,
		})
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

	if config.FlagDatabaseDSN != "" {
		Storage.Insert(req.Context(), shortKey, string(body))
	} else {
		storage.StorageObj.Set(shortKey, string(body))
		filestorage.Save(storage.StorageObj.Data)
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	var redirectLink string
	var err error

	if config.FlagDatabaseDSN != "" {
		redirectLink, err = Storage.Get(req.Context(), shortKey)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
	} else {
		redirectLink, err = storage.StorageObj.Get(shortKey)
		if err != nil {
			http.Error(res, "redirect not found", 404)
		}
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}