package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
	Bootstrap(ctx context.Context) error
	Set(ctx context.Context, shortKey string, url string) (string, error)
	Get(ctx context.Context, shortKey string) (string, error)
}

type dataToFile struct {
    Data map[string]string
}

type app struct {
	storage Storage
}

func NewApp(storage Storage) *app {
    return &app{storage: storage}
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

	shortKeyResult, err := a.storage.Set(req.Context(), shortKey, inputDataValue.URL)
	if err != nil {
		shortURL = config.FlagBaseAddr + "/" + shortKeyResult
		status = http.StatusConflict
	}

	data := dataToFile{Data: make(map[string]string)}
	data.Data[shortKey] = inputDataValue.URL
	filestorage.Save(data.Data)

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

		shortKeyResult, err := a.storage.Set(req.Context(), batchItem.CorrelationID, batchItem.OriginalURL)
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

	shortKeyResult, err := a.storage.Set(req.Context(), shortKey, string(body))
	if err != nil {
		shortURL = config.FlagBaseAddr + "/" + shortKeyResult
		status = http.StatusConflict
	}

	data := dataToFile{Data: make(map[string]string)}
	data.Data[shortKey] = string(body)
	filestorage.Save(data.Data)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(status)
	res.Write([]byte(shortURL))
}

func (a *app) Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	var redirectLink string
	var err error

	redirectLink, err = a.storage.Get(req.Context(), shortKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}