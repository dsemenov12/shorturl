package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/storage"
	"github.com/dsemenov12/shorturl/internal/util"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type dataToFile struct {
    Data map[string]string
}

type app struct {
	storage storage.Storage
}

func NewApp(storage storage.Storage) *app {
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

	storageData := make(map[string]string)
	storageData[shortKey] = inputDataValue.URL

	filestorage.Save(storageData)

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

	redirectLink, shortKeyRes, isDeleted, err := a.storage.Get(req.Context(), shortKey)

	fmt.Println("get redirect")
	fmt.Println(shortKeyRes)

	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
        return
	}
	if isDeleted {
		http.Error(res, "", http.StatusGone)
        return
	}

	http.Redirect(res, req, redirectLink, http.StatusTemporaryRedirect)
}

func (a *app) UserUrls(res http.ResponseWriter, req *http.Request) {
	var result []models.ShortURLItem
	var shortKey string
	var originalURL string

	rows, err := a.storage.GetUserURL(req.Context())
	if err != nil {
        http.Error(res, err.Error(), http.StatusNoContent)
		return
    }
	defer rows.Close()

	countRows := 0
	for rows.Next() {
        err = rows.Scan(&shortKey, &originalURL)
        if err != nil {
           continue
        }

		result = append(result, models.ShortURLItem{
			OriginalURL: originalURL,
			ShortURL: config.FlagBaseAddr + "/" + shortKey,
		})

		countRows++
    }

	if countRows == 0 {
        http.Error(res, "", http.StatusNoContent)
		return
    }

	resp, err := json.MarshalIndent(result, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func (a *app) DeleteUserUrls(res http.ResponseWriter, req *http.Request) {
	var shortKeys []string

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &shortKeys); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := generator(doneCh, shortKeys)

	resultCh := a.delete(req.Context(), doneCh, inputCh)

	for res := range resultCh {
        fmt.Println(res)
    }

	res.WriteHeader(http.StatusAccepted)
}

func (a *app) delete(ctx context.Context, doneCh chan struct{}, inputCh chan string) chan string {
	deleteRes := make(chan string)

	go func() {
		defer close(deleteRes)

		for data := range inputCh {
			a.storage.Delete(ctx, data)

			select {
			case <-doneCh:
				return
			case deleteRes <- "ok":
			}
		}
	}()

	return deleteRes
}

func generator(doneCh chan struct{}, input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}