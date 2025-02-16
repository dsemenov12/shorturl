package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"net/http/httptest"

	"github.com/dsemenov12/shorturl/internal/models"
	"github.com/dsemenov12/shorturl/internal/storage/memory"
	"github.com/dsemenov12/shorturl/internal/rand"
)

// Бенчмарк для ShortenPost
func BenchmarkShortenPost(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Пример данных
	inputData := models.InputData{
		URL: "https://example.com",
	}

	// Подготовка запроса
	reqBody, _ := json.Marshal(inputData)
	req := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.ShortenPost(res, req)
	}
}

// Бенчмарк для ShortenBatchPost
func BenchmarkShortenBatchPost(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Пример данных
	batch := []models.BatchItem{
		{CorrelationID: "id1", OriginalURL: "https://example.com"},
		{CorrelationID: "id2", OriginalURL: "https://another.com"},
	}

	// Подготовка запроса
	reqBody, _ := json.Marshal(batch)
	req := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.ShortenBatchPost(res, req)
	}
}

// Бенчмарк для PostURL
func BenchmarkPostURL(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Пример данных
	inputData := "https://example.com"

	// Подготовка запроса
	reqBody := []byte(inputData)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.PostURL(res, req)
	}
}

// Бенчмарк для Redirect
func BenchmarkRedirect(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Записываем данные в хранилище
	shortKey := rand.RandStringBytes(8)
	storage.Set(context.Background(), shortKey, "https://example.com")

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", shortKey), nil)
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.Redirect(res, req)
	}
}

// Бенчмарк для UserUrls
func BenchmarkUserUrls(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Записываем несколько данных
	for i := 0; i < 100; i++ {
		storage.Set(context.Background(), rand.RandStringBytes(8), "https://example.com")
	}

	req := httptest.NewRequest("GET", "/api/user/urls", nil)
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.UserUrls(res, req)
	}
}

// Бенчмарк для DeleteUserUrls
func BenchmarkDeleteUserUrls(b *testing.B) {
	storage := memory.NewStorage()
	app := NewApp(storage)

	// Пример данных
	shortKeys := []string{
		"short1",
		"short2",
		"short3",
	}

	// Записываем данные в хранилище
	for _, key := range shortKeys {
		storage.Set(context.Background(), key, "https://example.com")
	}

	reqBody, _ := json.Marshal(shortKeys)
	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewReader(reqBody))
	res := httptest.NewRecorder()

	// Запуск бенчмарка
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.DeleteUserUrls(res, req)
	}
}