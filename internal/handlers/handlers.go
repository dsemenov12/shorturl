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

type app struct {
    storageDB *pg.StorageDB
	storage *storage.Storage
}

func NewApp(storageDB *pg.StorageDB, storage *storage.Storage) *app {
    return &app{storageDB: storageDB, storage: storage}
}

func (a *app) Ping(res http.ResponseWriter, req *http.Request) {
	if err := a.storageDB.Ping(); err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
    }

	res.WriteHeader(http.StatusOK)
}

func (a *app) ShortenPost(res http.ResponseWriter, req *http.Request) {
	var inputDataValue models.InputData

	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey
	status := http.StatusCreated

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

	if config.FlagDatabaseDSN != "" {
		shortKeyResult, err := a.storageDB.Insert(req.Context(), shortKey, inputDataValue.URL)
		if err != nil {
			shortURL = config.FlagBaseAddr + "/" + shortKeyResult
			status = http.StatusConflict
		}
	} else {
		a.storage.Set(shortKey, inputDataValue.URL)
		filestorage.Save(a.storage.Data)
	}

	var result = models.ResultJSON{
		Result: shortURL,
	}

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write(resp)
}

func (a *app) ShortenBatchPost(res http.ResponseWriter, req *http.Request) {
	var batch []models.BatchItem
	var result []models.BatchResultItem

	status := http.StatusCreated

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
		if batchItem.CorrelationID == "" || batchItem.OriginalURL == "" {
			continue
		}

		shortURL := config.FlagBaseAddr + "/" + batchItem.CorrelationID

		shortKeyResult, err := a.storageDB.Insert(req.Context(), batchItem.CorrelationID, batchItem.OriginalURL)
		if err != nil {
			shortURL = config.FlagBaseAddr + "/" + shortKeyResult
			status = http.StatusConflict
		}

		result = append(result, models.BatchResultItem{
			CorrelationID: batchItem.CorrelationID,
			ShortURL: shortURL,
		})
	}

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write(resp)
}

func (a *app) PostURL(res http.ResponseWriter, req *http.Request) {
	shortKey := util.RandStringBytes(8)
	shortURL := config.FlagBaseAddr + "/" + shortKey
	status := http.StatusCreated

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
		shortKeyResult, err := a.storageDB.Insert(req.Context(), shortKey, string(body))
		if err != nil {
			shortURL = config.FlagBaseAddr + "/" + shortKeyResult
			status = http.StatusConflict
		}
	} else {
		a.storage.Set(shortKey, string(body))
		filestorage.Save(a.storage.Data)
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(status)
	res.Write([]byte(shortURL))
}

func (a *app) Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	var redirectLink string
	var err error

	if config.FlagDatabaseDSN != "" {
		redirectLink, err = a.storageDB.Get(req.Context(), shortKey)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
		}
	} else {
		redirectLink, err = a.storage.Get(shortKey)
		if err != nil {
			http.Error(res, "redirect not found", 404)
		}
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}