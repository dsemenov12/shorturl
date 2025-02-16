package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dsemenov12/shorturl/internal/config"
	"github.com/dsemenov12/shorturl/internal/filestorage"
	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/rand"
	"github.com/dsemenov12/shorturl/internal/storage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// dataToFile представляет структуру для хранения данных в файле.
type dataToFile struct {
	Data map[string]string
}

// app представляет основное приложение, которое взаимодействует с хранилищем.
type app struct {
	storage storage.Storage
}

// NewApp создает новый экземпляр приложения.
func NewApp(storage storage.Storage) *app {
	return &app{storage: storage}
}

// ShortenPost обрабатывает запрос на сокращение URL в формате JSON.
func (a *app) ShortenPost(res http.ResponseWriter, req *http.Request) {
	var inputDataValue models.InputData

	shortKey := rand.RandStringBytes(8)
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

// ShortenBatchPost обрабатывает пакетное сокращение URL.
// Ожидает массив объектов с исходными URL и возвращает массив с сокращёнными ссылками.
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
			ShortURL:      shortURL,
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

// PostURL обрабатывает сокращение URL, переданного в теле запроса.
// Принимает URL в виде текста, сокращает его и возвращает короткую ссылку.
func (a *app) PostURL(res http.ResponseWriter, req *http.Request) {
	shortKey := rand.RandStringBytes(8)
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

// Redirect обрабатывает перенаправление по короткому URL.
func (a *app) Redirect(res http.ResponseWriter, req *http.Request) {
	shortKey := chi.URLParam(req, "id")

	var redirectLink string
	var err error

	redirectLink, _, isDeleted, err := a.storage.Get(req.Context(), shortKey)

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

// UserUrls возвращает список URL, сохраненных пользователем.
func (a *app) UserUrls(res http.ResponseWriter, req *http.Request) {
	result, err := a.storage.GetUserURL(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusNoContent)
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

// DeleteUserUrls обрабатывает запрос на удаление списка сокращенных URL-адресов,
// полученного в теле запроса в формате JSON.
// Функция читает список коротких URL, передает их в канал для обработки
// и запускает процесс удаления в фоне.
// В случае успешного запуска удаления возвращает HTTP статус 202 (Accepted).
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
			splitData := strings.Split(data, "/")
			code := splitData[len(splitData)-1]

			a.storage.Delete(ctx, code)

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
